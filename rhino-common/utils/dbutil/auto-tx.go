package dbutil

import (
	"database/sql"
	"rhino-common/domain_error"
	"sync"
	"sync/atomic"
	"time"
)

var(
	totalTxCount int64
)

type TxFunc func(*sql.Tx) (de *domain_error.Error)

type txFuncAndDB struct {
	txFunc TxFunc
	db     *sql.DB
}

type AutoTx struct {
	tick         time.Duration
	flushTxCount int
	txFuncs      []TxFunc
	lock         *sync.Mutex
	input        chan *txFuncAndDB
	outputError  chan *domain_error.Error
	lastDB       *sql.DB
}

func NewAutoTx(tick time.Duration, flushTxCount int) (autoTx *AutoTx, outputError chan *domain_error.Error) {
	autoTx = &AutoTx{tick: tick, flushTxCount: flushTxCount, lock: &sync.Mutex{}, input: make(chan *txFuncAndDB, 100000), outputError: make(chan *domain_error.Error, 1)}
	outputError = autoTx.outputError
	return
}

func (atx *AutoTx) StartTicker() {
	ticker := time.NewTicker(atx.tick)
	go func() {
		for {
			<-ticker.C
			if len(atx.txFuncs) > 0 {
				atx.lock.Lock()
				atx.commit()
				atx.lock.Unlock()
			}
		}
	}()

	go func ()  {
		for {
			data := <-atx.input
			atx.onAddAndRunTxFunction(data.db, data.txFunc)
		}
	}()
}

func (atx *AutoTx) commit() {
	
	if len(atx.txFuncs) == 0 {
		return
	}

	tx, de := BeginTx(atx.lastDB)
	if de != nil {
		select {
		case atx.outputError <- de:
		default:
			// 非阻塞
		}
		return
	}

	n := len(atx.txFuncs)

	for _, f := range atx.txFuncs {
		de := f(tx)
		if de != nil {
			select {
			case atx.outputError <- de:
			default:
				// 非阻塞
			}
			RollbackTx(tx)
			return
		}
	}
	// 提交事务
	de = CommitTx(tx)
	if de != nil {
		select {
		case atx.outputError <- de:
		default:
			// 非阻塞
		}
		RollbackTx(tx)
	} else {
		atx.txFuncs = atx.txFuncs[:0]
	}
	atomic.AddInt64(&totalTxCount, int64(n))
	//log.Printf("finish commit %d tx. totalTxCount=%d", n, totalTxCount)
}

func (atx *AutoTx) AddAndRunTxFunction(db *sql.DB, txFunc TxFunc) {
	atx.input <- &txFuncAndDB{txFunc: txFunc, db: db}
}

func (atx *AutoTx) onAddAndRunTxFunction(db *sql.DB, txFunc TxFunc) {
	atx.lock.Lock()
	defer atx.lock.Unlock()
	atx.lastDB = db
	atx.txFuncs = append(atx.txFuncs, txFunc)
	if len(atx.txFuncs) >= atx.flushTxCount {
		atx.commit()
	}
}

func (atx *AutoTx) UpdateDB(db *sql.DB) {
	atx.lock.Lock()
	defer atx.lock.Unlock()
	atx.lastDB = db
}

func (atx *AutoTx) Flush(db *sql.DB) {
	if len(atx.txFuncs) == 0 {
		return
	}
	atx.lock.Lock()
	defer atx.lock.Unlock()
	atx.lastDB = db
	atx.commit()
}
