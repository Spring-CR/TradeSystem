package fix_store

import (
	"log"
	"rhino-core/domain_cfg"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
	"github.com/quickfixgo/quickfix"
)

var(
	FixStore quickfix.MessageStore
)

type advanceMemoryStore struct {
	msgSeqGen    *domain_cfg.MsgSeqGen
	lock         *sync.RWMutex
	creationTime time.Time
	messageMap   map[int][]byte
	resetCount   int32
}

func (store *advanceMemoryStore) NextSenderMsgSeqNum() int {
	return store.msgSeqGen.NextSenderMsgSeqNum()
}

func (store *advanceMemoryStore) NextTargetMsgSeqNum() int {
	return store.msgSeqGen.NextTargetMsgSeqNum()
}

func (store *advanceMemoryStore) IncrNextSenderMsgSeqNum() error {
	return store.msgSeqGen.IncrNextSenderMsgSeqNum()
}

func (store *advanceMemoryStore) IncrNextTargetMsgSeqNum() error {
	return store.msgSeqGen.IncrNextTargetMsgSeqNum()
}

func (store *advanceMemoryStore) SetNextSenderMsgSeqNum(nextSeqNum int) error {
	return store.msgSeqGen.SetNextSenderMsgSeqNum(nextSeqNum)
}
func (store *advanceMemoryStore) SetNextTargetMsgSeqNum(nextSeqNum int) error {
	return store.msgSeqGen.SetNextTargetMsgSeqNum(nextSeqNum)
}

func (store *advanceMemoryStore) CreationTime() time.Time {
	return store.creationTime
}

func (store *advanceMemoryStore) SetCreationTime(t time.Time) {
	store.creationTime = t
}

func (store *advanceMemoryStore) Reset() error {
	log.Println("reset...")
	defer func() {
		atomic.AddInt32(&store.resetCount, 1)
	}()

	if store.msgSeqGen.GetReachMaxRetryLogonFail() == 1 {
		log.Println("do nothing...")
		return nil
	}

	if atomic.LoadInt32(&store.resetCount) > 0 {
		store.msgSeqGen.Reset()
	}

	store.creationTime = time.Now()
	store.messageMap = nil

	return nil
}

func (store *advanceMemoryStore) Refresh() error {
	// NOP, nothing to refresh.
	return nil
}

func (store *advanceMemoryStore) Close() error {
	// NOP, nothing to close.
	return nil
}

func (store *advanceMemoryStore) SaveMessage(seqNum int, msg []byte) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	if store.messageMap == nil {
		store.messageMap = make(map[int][]byte)
	}

	store.messageMap[seqNum] = msg
	return nil
}

func (store *advanceMemoryStore) SaveMessageAndIncrNextSenderMsgSeqNum(seqNum int, msg []byte) error {
	err := store.SaveMessage(seqNum, msg)
	if err != nil {
		return err
	}
	return store.IncrNextSenderMsgSeqNum()
}

func (store *advanceMemoryStore) IterateMessages(beginSeqNum, endSeqNum int, cb func([]byte) error) error {
	for seqNum := beginSeqNum; seqNum <= endSeqNum; seqNum++ {
		store.lock.RLock() // 加锁
		if m, ok := store.messageMap[seqNum]; ok {
			if err := cb(m); err != nil {
				store.lock.RUnlock() // 解锁
				return err
			}
		}
		store.lock.RUnlock() // 解锁
	}
	return nil
}

func (store *advanceMemoryStore) GetMessages(beginSeqNum, endSeqNum int) ([][]byte, error) {
	var msgs [][]byte
	err := store.IterateMessages(beginSeqNum, endSeqNum, func(m []byte) error {
		msgs = append(msgs, m)
		return nil
	})
	return msgs, err
}

type advanceMemoryStoreFactory struct {
	msgSeqGen *domain_cfg.MsgSeqGen
}

func (f advanceMemoryStoreFactory) Create(s quickfix.SessionID) (quickfix.MessageStore, error) {
	m := new(advanceMemoryStore)
	m.lock = &sync.RWMutex{}
	m.msgSeqGen = f.msgSeqGen
	log.Printf("Create memory for session:%v-%v-%v,  NextSenderMsgSeqNum:%d, NextTargetMsgSeqNum:%d\n", s.BeginString, s.SenderCompID, s.TargetCompID, m.NextSenderMsgSeqNum(), m.NextTargetMsgSeqNum())
	if err := m.Reset(); err != nil {
		return m, errors.Wrap(err, "reset")
	}
	FixStore = m
	return m, nil
}

// NewMemoryStoreFactory returns a MessageStoreFactory instance that created in-memory MessageStores.
func NewAdvanceMemoryStoreFactory(msgSeqGen *domain_cfg.MsgSeqGen) quickfix.MessageStoreFactory { return advanceMemoryStoreFactory{msgSeqGen: msgSeqGen} }
