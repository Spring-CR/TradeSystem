package quota

import (
	"log"
	"rhino-common/utils/atomicutil"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type QuotaAcquire[T any] struct {
	source        T
	acquiredQuota float64
	acquiredTime  time.Time
	deleteFlag    uint32
	quietChecked  uint32
}

func NewQuotaAcquire[T any](source T, acquiredQuota float64) *QuotaAcquire[T] {
	return &QuotaAcquire[T]{
		source:        source,
		acquiredQuota: acquiredQuota,
		acquiredTime:  time.Now(),
	}
}

func (qa *QuotaAcquire[T]) GetAcquiredQuota() float64 {
	return atomicutil.LoadFloat64(&qa.acquiredQuota)
}

func (qa *QuotaAcquire[T]) GetSource() T {
	return qa.source
}

type QuotaReturn[T any] struct {
	source      T
	returnQuota float64
}

func NewQuotaReturn[T any](source T, returnQuota float64) *QuotaReturn[T] {
	return &QuotaReturn[T]{
		source:      source,
		returnQuota: returnQuota,
	}
}

func (qa *QuotaReturn[T]) GetSource() T {
	return qa.source
}

type QuotaControl[T, C any] struct {
	lockAcquire      *sync.RWMutex                                    // quotaAcquiredMap锁
	lockReturn       *sync.RWMutex                                    // quotaReturnMap锁
	lockMetadata     *sync.RWMutex                                    // metadata锁
	metadata         map[string]interface{}                           // 元数据：订单基础数据+底仓的基础数据+派生计算结果
	baseQuota        float64                                          // 初始的额度，例如T-1日终的底仓
	quota            float64                                          // 可用额度
	lockedQuota      float64                                          // 锁定的额度
	quietPeriod      time.Duration                                    // 额度锁定静默期（在订单首次提交， 不发生任何状态变化的最大锁定时长）
	checkPeriod      time.Duration                                    // 检查订单静默期是否超限的循环等待时间
	quotaAcquiredMap map[string]*QuotaAcquire[T]                      // 额度请求map, key为appOrdID，对应卖单
	quotaReturnMap   map[string]*QuotaReturn[C]                       // 额度返回map, key为execID，对应是买单的成交回报
	releaseFunc      func(q *QuotaAcquire[T]) (releasedQuota float64) // 额度释放函数
}

func NewQuotaControl[T, C any](metadata map[string]interface{}, baseQuota float64, quietPeriod time.Duration, checkPeriod time.Duration, releaseFunc func(q *QuotaAcquire[T]) (releasedQuota float64)) *QuotaControl[T, C] {
	qc := &QuotaControl[T, C]{
		lockAcquire:      &sync.RWMutex{},
		lockReturn:       &sync.RWMutex{},
		lockMetadata:     &sync.RWMutex{},
		metadata:         metadata,
		baseQuota:        baseQuota,
		quota:            baseQuota,
		lockedQuota:      0,
		quietPeriod:      quietPeriod,
		checkPeriod:      checkPeriod,
		quotaAcquiredMap: map[string]*QuotaAcquire[T]{},
		quotaReturnMap:   map[string]*QuotaReturn[C]{},
		releaseFunc:      releaseFunc,
	}
	qc.starCheckQuietQuota()
	return qc
}

func (qc *QuotaControl[T, C]) starCheckQuietQuota() {
	go func() {
		ticker := time.NewTicker(qc.checkPeriod)
		for {
			timeNow := <-ticker.C
			// var delKeys[]string
			qc.lockAcquire.RLock()
			for _, quotaAcquired := range qc.quotaAcquiredMap {
				// 已经标识为删除的就忽略
				if atomic.LoadUint32(&quotaAcquired.deleteFlag) > 0 || atomic.LoadUint32(&quotaAcquired.quietChecked) > 0 {
					continue
				}
				timeAcquireQuota := quotaAcquired.acquiredTime
				if timeNow.Sub(timeAcquireQuota) > qc.quietPeriod {
					releasedQuota := qc.releaseFunc(quotaAcquired)
					if releasedQuota == 0 {
						// 如果订单处于可交易状态，应该收到返回值为0的可释放额度；如果订单刚好完全成交，也会返回0
						// 并且，标记静默期检查已经通过了
						atomic.StoreUint32(&quotaAcquired.quietChecked, 1)
					} else if releasedQuota == quotaAcquired.acquiredQuota && atomic.CompareAndSwapUint32(&quotaAcquired.deleteFlag, 0, 1) { // 如果额度和请求时一致，则标记删除
						// 返回额度
						atomicutil.AddFloat64(&qc.quota, releasedQuota)
						// 释放锁定量
						atomicutil.AddFloat64(&qc.lockedQuota, -releasedQuota)
						log.Printf("selling order is timeout for entering trading phase, go to release total quota acquired %v\n", quotaAcquired.acquiredQuota)
					}
				}
			}
			qc.lockAcquire.RUnlock()
		}
	}()
}

// 需要登记该key初始锁定的额度
func (qc *QuotaControl[T, C]) AcquireQuota(key string, q *QuotaAcquire[T], force bool, allowNegative bool) (ok bool) {

	newVal := atomicutil.AddFloat64(&qc.quota, -q.acquiredQuota)

	log.Printf("===> AcquireQuota, key:%v, force:%v, allowNegative:%v, newVal:%v\n", key, force, allowNegative, newVal)

	if !force && newVal < 0 {
		atomicutil.AddFloat64(&qc.quota, q.acquiredQuota)
		return
	}

	ok = true

	// 增加锁额
	atomicutil.AddFloat64(&qc.lockedQuota, q.acquiredQuota)
	if !allowNegative && newVal < 0 {
		// 确保最少可用额度为0
		atomicutil.AddFloat64(&qc.lockedQuota, newVal)
		atomicutil.AddFloat64(&qc.quota, -newVal)
		// 设置真实的锁定额度
		q.acquiredQuota += newVal
	}

	qc.lockAcquire.Lock()
	qc.quotaAcquiredMap[key] = q

	log.Printf("===> New Quota, key:%v, force:%v, allowNegative:%v, newVal:%v\n", key, force, allowNegative, qc.quota)

	qc.lockAcquire.Unlock()

	return
}

func (qc *QuotaControl[T, C]) ReturnQuota(key string, q *QuotaReturn[C]) (ok bool) {

	qc.lockReturn.Lock()
	defer qc.lockReturn.Unlock()

	_, exist := qc.quotaReturnMap[key]
	if exist {
		return false
	}

	qc.quotaReturnMap[key] = q
	//增加可用额度
	newQuota := atomicutil.AddFloat64(&qc.quota, q.returnQuota)
	atomicutil.AddFloat64(&qc.lockedQuota, -q.returnQuota)

	log.Printf("buying order return quota %v by %v, newQuota: %v\n", q.returnQuota, key, newQuota)

	return
}

func(qc*QuotaControl[T, C]) HasAcquireQuto(key string) bool {
	qc.lockAcquire.RLock()
	_, ok := qc.quotaAcquiredMap[key]
	qc.lockAcquire.RUnlock()
	return ok
}

func (qc *QuotaControl[T, C]) ReleaseQuota(key string) {
	qc.lockAcquire.RLock()
	quotaAcquired, ok := qc.quotaAcquiredMap[key]
	qc.lockAcquire.RUnlock()
	if !ok {
		log.Printf("cannot found key %s and release quota\n", key)
		var keys []string
		qc.lockAcquire.RLock()
		for k := range qc.quotaAcquiredMap {
			keys = append(keys, k)
		}
		log.Printf("current quotaAcquiredMap keys:%s\n", strings.Join(keys, ","))
		qc.lockAcquire.RUnlock()
		return
	}

	swapped := atomic.CompareAndSwapUint32(&quotaAcquired.deleteFlag, 0, 1)
	if !swapped {
		// 已经被标记为删除了，可以直接退出
		return
	}

	releasedQuota := qc.releaseFunc(quotaAcquired)
	if releasedQuota <= 0 {
		return
	}
	// 返回额度，标记删除
	newQuota := atomicutil.AddFloat64(&qc.quota, releasedQuota)
	// 释放锁定量
	atomicutil.AddFloat64(&qc.lockedQuota, -releasedQuota)

	log.Printf("release quota:%v, for key:%s, newQuota:%v\n", releasedQuota, key, newQuota)
}

func (qc *QuotaControl[T, C]) GetMetadata() (map[string]interface{}, *sync.RWMutex) {
	return qc.metadata, qc.lockMetadata
}

func (qc *QuotaControl[T, C]) GetQuota() (baseQuota float64, quota float64) {
	baseQuota = atomicutil.LoadFloat64(&qc.baseQuota)
	quota = atomicutil.LoadFloat64(&qc.quota)
	return
}

func (qc *QuotaControl[T, C]) UpdateQuota(newBaseQuota float64) (newQuota float64, changed bool) {
	// 发生变动则设置ok
	diff := newBaseQuota - qc.baseQuota
	if diff == 0 {
		return atomicutil.LoadFloat64(&qc.quota), false
	}
	qc.baseQuota = newBaseQuota
	return atomicutil.AddFloat64(&qc.quota, diff), true
}
