package dbutil

import (
	"database/sql"
	"log"
	"rhino-common/domain_error"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

const (
	WorkerInputBuffer = 128
	AutoTxInputBuffer = 1024
)

type txTask struct {
	txFunc   TxFunc // 事务执行函数
	preCond  string // 前件，即该事务执行时，需要满足的前提条件
	postCond string // 后件，即存在其他的事物，需要等待该事务执行的执行结果
}

type txTaskWorker struct {
	num            int
	parent         *ConcurrentAutoTx
	flushTxCount   int                      // 如果缓存的事务执行任务满足flushTxCount所指定的数量，则强制执行缓存中的所有任务
	txTasks        []*txTask                // 用于缓存事务执行任务
	inputTxTask    chan *txTask             // 接收事务执行任务的channel
	inputTick      <-chan time.Time         // 接收周期性的时钟信号，当接收到该信号时，如果任务缓存不为空，则强制执行缓存中的所有任务
	flushSignal    chan bool                // 强制提交tx
	outputPostCond chan string              // 对于定义了后件的事务执行任务，执行完毕之后需要往该channel发送后件
	outputError    chan *domain_error.Error // 在worker执行任务的过程中，如果发生异常，需要往该channel报送。注意：违反唯一性约束的插入操作，可以不必进行报送
}

func newTxTaskWorker(i int, parent *ConcurrentAutoTx) *txTaskWorker {
	worker := &txTaskWorker{
		num:            i,
		parent:         parent,
		flushTxCount:   parent.flushTxCount,
		inputTxTask:    make(chan *txTask, WorkerInputBuffer),
		inputTick:      time.NewTicker(parent.tick).C,
		flushSignal:    make(chan bool),
		outputPostCond: parent.workerOutputPostCond,
		outputError:    parent.outputError,
	}
	worker.start()
	return worker
}

func (w *txTaskWorker) process(task *txTask) {
	w.inputTxTask <- task
}

func (w *txTaskWorker) start() {
	go func() {
		for {
			select {
			case <-w.inputTick: // 时钟周期到达，强制执行缓存中的全部任务
				if len(w.txTasks) > 0 {
					w.commit()
				}
			case txTask := <-w.inputTxTask: // 缓存满，强制执行缓存中的全部任务
				w.txTasks = append(w.txTasks, txTask)
				if len(w.txTasks) >= w.flushTxCount {
					w.commit()
				}
			case <-w.flushSignal:
				if len(w.txTasks) > 0 {
					w.commit()
				}
			}
		}
	}()
}

func (w *txTaskWorker) commit() {
	var postConds []string
	lastDB := w.parent.getDB()

	tx, de := BeginTx(lastDB)
	if de != nil {
		select {
		case w.outputError <- de:
		default:
			// 非阻塞
			log.Println("Error sending to outputError channel:", de) // 添加日志
		}
		return
	}

	n := len(w.txTasks)

	for _, f := range w.txTasks {
		de := f.txFunc(tx)
		if de != nil {
			select {
			case w.outputError <- de:
			default:
				// 非阻塞
			}
			RollbackTx(tx)
			return
		}
		if len(f.postCond) > 0 { // 累积后件不为空的记录
			postConds = append(postConds, f.postCond)
		}
	}
	// 提交事务
	de = CommitTx(tx)
	if de != nil {
		select {
		case w.outputError <- de:
		default:
			// 非阻塞
		}
		RollbackTx(tx)
		return
	} else {
		// 重置任务队列
		w.txTasks = w.txTasks[:0]
		// 发送全部后件
		for _, postCond := range postConds {
			w.outputPostCond <- postCond
		}
	}
	atomic.AddInt64(&totalTxCount, int64(n))
	//log.Printf("finish commit %d tx. totalTxCount=%d", n, totalTxCount)
}

func (w *txTaskWorker) flush() {
	w.flushSignal <- true
}

type ConcurrentAutoTx struct {
	tick                 time.Duration            // 时钟周期
	flushTxCount         int                      // 指定事务执行器用于缓存事务执行任务的数目，如果缓存的事务执行任务满足flushTxCount所指定的数量，则强制执行缓存中的所有任务
	lastDB               unsafe.Pointer           // 最新的DB
	postCondMap          *sync.Map                // 后件哈希表，指关注键
	pendingTasks         *sync.Map                // 待处理任务的哈希表，键是事务任务的前件，值是一个待处理的任务列表
	workers              []*txTaskWorker          // 事务执行器，注意，执行器的数量必须是2的幂
	workerOutputPostCond chan string              // 用于接收事务执行器发送的后件
	inputTxTask          chan *txTask             // 用于接收输入的事务执行任务
	outputError          chan *domain_error.Error // 用于向调用方反馈异常
}

type pendingTaskList struct {
	lock        *sync.RWMutex
	pendingTask []*txTask
}

func newPendingTaskList() *pendingTaskList {
	return &pendingTaskList{lock: &sync.RWMutex{}}
}

func (l *pendingTaskList) size() int {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return len(l.pendingTask)
}

func (l *pendingTaskList) add(task *txTask) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.pendingTask = append(l.pendingTask, task)
}

func (l *pendingTaskList) getAndPopAll() []*txTask {
	if l.size() == 0 {
		return nil
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	tasks := l.pendingTask
	l.pendingTask = nil
	return tasks
}

func NewConcurrentAutoTx(workerCount int, tick time.Duration, flushTxCount int) (autoTx *ConcurrentAutoTx, outputError chan *domain_error.Error) {

	// 确保 size 是 2 的幂
	if workerCount&(workerCount-1) != 0 {
		panic("workerCount must be a power of 2")
	}

	autoTx = &ConcurrentAutoTx{
		tick:                 tick,
		flushTxCount:         flushTxCount,
		postCondMap:          &sync.Map{},
		pendingTasks:         &sync.Map{},
		inputTxTask:          make(chan *txTask, AutoTxInputBuffer),
		workerOutputPostCond: make(chan string, WorkerInputBuffer*workerCount*2),
		outputError:          make(chan *domain_error.Error, 1),
	}

	for i := 0; i < workerCount; i++ {
		autoTx.workers = append(autoTx.workers, newTxTaskWorker(i, autoTx))
	}

	outputError = autoTx.outputError

	go func() {
		for {
			total := atomic.LoadInt64(&totalTxCount)
			log.Printf("atx.inputTxTask len=%d, autoTx.workerOutputPostCond len=%d, totalTxCount=%d\n", len(autoTx.inputTxTask), len(autoTx.workerOutputPostCond), total)
			time.Sleep(30 * time.Second)
		}
	}()

	return
}

func (atx *ConcurrentAutoTx) Start() {
	// 开始分派任务
	atx.dispatch()
	// 持续处理后件
	atx.processPostCond()
}

func (atx *ConcurrentAutoTx) getDB() *sql.DB {
	lastDB := (*sql.DB)(atomic.LoadPointer(&atx.lastDB))
	return lastDB
}

func (atx *ConcurrentAutoTx) setDB(lastDB *sql.DB) {
	atomic.CompareAndSwapPointer(&atx.lastDB, atx.lastDB, unsafe.Pointer(lastDB))
	//atomic.StorePointer(&atx.lastDB, unsafe.Pointer(lastDB))
}

func (atx *ConcurrentAutoTx) Input(db *sql.DB, txFunc TxFunc, preCond, postCond string) {
	// 更新db
	atx.setDB(db)
	atx.inputTxTask <- &txTask{txFunc: txFunc, preCond: preCond, postCond: postCond}
}

func (atx *ConcurrentAutoTx) dispatch() {
	go func() {

		taskCount := 0 // 记录任务数量
		n := len(atx.workers) - 1

		for {
			txTask := <-atx.inputTxTask
			if len(txTask.preCond) == 0 { // 无前件依赖，分发到worker
				workerIndex := taskCount & n // 使用按位与选择 worker
				worker := atx.workers[workerIndex]
				taskCount++

				worker.process(txTask)

			} else { // 有前件依赖
				if _, ok := atx.postCondMap.Load(txTask.preCond); ok { // 前后件匹配
					workerIndex := taskCount & n // 使用按位与选择 worker
					worker := atx.workers[workerIndex]
					taskCount++

					worker.process(txTask)

				} else { // 加入到等待队列
					value, _ := atx.pendingTasks.LoadOrStore(txTask.preCond, newPendingTaskList())
					pendingTaskList := value.(*pendingTaskList)
					pendingTaskList.add(txTask)
				}
			}
		}
	}()
}

// 处理后件
func (atx *ConcurrentAutoTx) processPostCond() {
	ticker := time.NewTicker(atx.tick)
	go func() {
		for {
			select {
			case postCond := <-atx.workerOutputPostCond:
				_, ok := atx.postCondMap.Load(postCond)
				if !ok {
					atx.postCondMap.Store(postCond, true)
					// 检查是否有后件任务列表
					value, loaded := atx.pendingTasks.Load(postCond)
					if loaded {
						atx.processPendingTaskList(value)
					}
				}
			case <-ticker.C:
				var keyDelete []any
				atx.pendingTasks.Range(func(key any, value any) bool {
					if _, ok := atx.postCondMap.Load(key); ok { // 前后件匹配
						atx.processPendingTaskList(value)
						// 准备从等待列表中删除对应的键值对
						keyDelete = append(keyDelete, key)
					}
					return true
				})
				// 从等待列表中删除已经前后件匹配的任务列表
				for _, key := range keyDelete {
					atx.pendingTasks.Delete(key)
				}
			}
		}
	}()
}

func (atx *ConcurrentAutoTx) processPendingTaskList(value any) {
	pendingTaskList := value.(*pendingTaskList)
	tasks := pendingTaskList.getAndPopAll()
	for _, task := range tasks {
		atx.inputTxTask <- task
	}
}

// 强制提交全部tx
func (atx *ConcurrentAutoTx) Flush() {
	for _, w := range atx.workers {
		w.flush()
	}
}