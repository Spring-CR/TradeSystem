package datamap

import (
	"database/sql"
	"rhino-common/enum"
	"rhino-common/utils/fuzzy"
	"rhino-core/schema"
	"sync"
	"sync/atomic"
)

type DualDataMap struct {
	db             *sql.DB
	lock           *sync.Mutex
	tableConfig    *TableConfigForMapping
	enableFuzzyMap bool
	dataSyncLog    *schema.DataSyncLog
	readMap        atomic.Pointer[map[string][]map[string]interface{}] // 存储 map 的指针
	writeMap       map[string][]map[string]interface{}
	fuzzyMap       atomic.Pointer[fuzzy.FuzzyMap[map[string]interface{}]]
}

func NewDualDataMap(db *sql.DB, tableConfig *TableConfigForMapping, enableFuzzyMap bool) *DualDataMap {

	inst := &DualDataMap{db: db, lock: &sync.Mutex{}, tableConfig: tableConfig, enableFuzzyMap: enableFuzzyMap, writeMap: make(map[string][]map[string]interface{})}

	if len(inst.tableConfig.ColumnsForUnique) == 0 {
		for _, v := range inst.tableConfig.TableConfig.Columns {
			inst.tableConfig.ColumnsForUnique = append(inst.tableConfig.ColumnsForUnique, v.ColumnAlias)
		}
	}

	return inst
}

func (m *DualDataMap) Get(key string) (val []map[string]interface{}, ok bool) {
	readPtr := m.readMap.Load()
	if readPtr == nil {
		return nil, false
	}
	readMap := *readPtr
	val, ok = readMap[key]
	// for k := range readMap {
	// 	log.Printf("======》 debug only ======> key:%s\n", k)
	// }
	return
}

func (m *DualDataMap) GetByFuzzyKey(key string) (val []map[string]interface{}, ok bool) {
	readPtr := m.fuzzyMap.Load()
	if readPtr == nil {
		return nil, false
	}
	fuzzyMap := *readPtr
	val, ok = fuzzyMap.Get(key)
	return
}

func (m *DualDataMap) GetReadMap() map[string][]map[string]interface{} {
	readPtr := m.readMap.Load()
	if readPtr == nil {
		return nil
	}
	readMap := *readPtr
	return readMap
}

// Reset 重置所有数据映射，清空数据但保留配置
func (m *DualDataMap) Reset() {
	m.lock.Lock()
	defer m.lock.Unlock()

	// 清空 writeMap 但保留映射结构
	if m.writeMap != nil {
		for k := range m.writeMap {
			delete(m.writeMap, k)
		}
	}

	// 获取当前的 readMap 指针，清空其内容但不替换指针
	currentReadMapPtr := m.readMap.Load()
	if currentReadMapPtr != nil {
		currentReadMap := *currentReadMapPtr
		for k := range currentReadMap {
			delete(currentReadMap, k)
		}
	} else {
		// 如果当前为空，创建一个新的空映射
		emptyReadMap := make(map[string][]map[string]interface{})
		m.readMap.Store(&emptyReadMap)
	}

	m.fuzzyMap.Store(nil)
	m.dataSyncLog = nil
}

func (m *DualDataMap) IsReady() bool {
	return m.dataSyncLog!=nil && m.dataSyncLog.SyncPhase == int(enum.DataSyncLogPhase_Complete)
}

func (m *DualDataMap) GetSyncLog() (syncLog *schema.DataSyncLog, ok bool) {
	if m.dataSyncLog == nil {
		return
	}
	syncLog = m.dataSyncLog
	ok = true
	return
}