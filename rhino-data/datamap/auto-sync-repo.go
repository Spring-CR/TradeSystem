package datamap

import (
	"database/sql"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"sort"
	"time"

	"github.com/manucorporat/try"
)

type DataMapConfig struct {
	SystemCode         string
	BusinessCode       string
	CentralDatabaseUrl string
	DatabaseUrl        string
	EnableFuzzyMap     bool
	TableConfigs       []*TableConfigForMapping
}

type DataChangeEvent struct {
	SystemCode   string
	BusinessCode string
	TableAlias   string
	AddKeys      []string
	ChgKeys      []string
	DelKeys      []string
}

type AutoSyncRepo struct {
	config       *DataMapConfig
	centralDB    *sql.DB
	appDB        *sql.DB
	loopInterval time.Duration
	eventCh      chan *DataChangeEvent
	tableMaps    map[string]*DualDataMap
}

func NewAutoSyncRepo(config *DataMapConfig, centralDB *sql.DB, appDB *sql.DB, loopInterval time.Duration, eventCh chan *DataChangeEvent) *AutoSyncRepo {
	// 循环间隔不能低于20秒
	if loopInterval < 20*time.Second {
		loopInterval = 20 * time.Second
	}
	inst := &AutoSyncRepo{config: config, centralDB: centralDB, appDB: appDB, loopInterval: loopInterval, eventCh: eventCh, tableMaps: map[string]*DualDataMap{}}
	inst.initDB()
	inst.initTableMaps()
	return inst
}

func (r *AutoSyncRepo) initDB() {

	if r.centralDB == nil {

		db, err := sql.Open("mysql", r.config.CentralDatabaseUrl)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to create central db")
		}
		db.SetMaxOpenConns(4) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
		db.SetMaxIdleConns(4) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
		db.SetConnMaxLifetime(time.Second * 600)

		r.centralDB = db
	}

	if r.appDB == nil {

		db, err := sql.Open("mysql", r.config.DatabaseUrl)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to create app db")
		}
		db.SetMaxOpenConns(4) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
		db.SetMaxIdleConns(4) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
		db.SetConnMaxLifetime(time.Second * 600)

		r.appDB = db
	}
}

func (r *AutoSyncRepo) initTableMaps() {
	for _, tableConfig := range r.config.TableConfigs {
		tableMap := NewDualDataMap(r.appDB, tableConfig, r.config.EnableFuzzyMap)
		r.tableMaps[tableConfig.TableAlias] = tableMap
		r.tableMaps[tableConfig.TableName] = tableMap
	}
}

func (r *AutoSyncRepo) Start() {
	go func() {
		for {
			try.This(func() {

				log.Printf("Start to sync data...")

				for _, tableConfig := range r.config.TableConfigs {

					log.Printf("======>for table:%s\n", tableConfig.TableAlias)

					dataSysLogs, err := admin_store.FindDataSyncLogsBySystemCodeAndBusinessCodeAndTableNameAndSyncPhase(r.centralDB, r.config.SystemCode, r.config.BusinessCode, tableConfig.TableName, int(enum.DataSyncLogPhase_Complete))
					if dbutil.IsDbRecordEmptyError(err) {
						err = nil
					}
					log.Printf("dataSysLogs.Len=%d\n", len(dataSysLogs))

					if err != nil {
						domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to FindDataSyncLogsBySystemCodeAndBusinessCodeAndTableNameAndSyncPhase:%v, %v, %v, %v", r.config.SystemCode, r.config.BusinessCode, tableConfig.TableName, int(enum.DataSyncLogPhase_Complete)))
						continue
					}
					if len(dataSysLogs) == 0 {
						log.Printf("dataSysLogs is 0, continue...")
						continue
					}

					if len(dataSysLogs) > 1 {
						sort.Slice(dataSysLogs, func(i, j int) bool {
							return dataSysLogs[i].ReportTime > dataSysLogs[j].ReportTime
						})
					}
					dataSysLog := dataSysLogs[0]

					tableMap, ok := r.tableMaps[tableConfig.TableName]
					if !ok {
						domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to get tableMap for %s", tableConfig.TableName))
						continue
					}

					log.Printf("start to SyncData %s\n", dataSysLog.TableName)
					synced, addKeys, chgKeys, delKeys := tableMap.SyncData(dataSysLog)
					log.Printf("===> SyncData for %s, synced:%v, addKeys:%d, chgKeys:%d, delKeys:%d, dataSysLog-%s-%d\n", tableConfig.TableName, synced, len(addKeys), len(chgKeys), len(delKeys), dataSysLog.TableName, dataSysLog.ReportTime)
					if r.eventCh != nil && synced && (len(addKeys) > 0 || len(chgKeys) > 0 || len(delKeys) > 0) {
						evt := &DataChangeEvent{
							SystemCode:   r.config.SystemCode,
							BusinessCode: r.config.BusinessCode,
							TableAlias:   tableConfig.TableAlias,
							AddKeys:      addKeys,
							ChgKeys:      chgKeys,
							DelKeys:      delKeys,
						}
						r.eventCh <- evt
					}
				}

			}).Catch(func(err try.E) {
			})

			time.Sleep(r.loopInterval)
		}
	}()
}

func (r *AutoSyncRepo) Get(collection string, key string) (val []map[string]interface{}, ok bool, de *domain_error.Error) {
	t, ok := r.tableMaps[collection]
	if !ok {
		de = domain_error.Build(domain_error.COLLECTION_NOT_FOUND_ERR_CODE, nil, collection)
		return
	}

	val, ok = t.Get(key)

	return
}

func (r *AutoSyncRepo) GetByKeyList(collection string, keyList[]string) (val []map[string]interface{}, de *domain_error.Error) {
	t, ok := r.tableMaps[collection]
	if !ok {
		de = domain_error.Build(domain_error.COLLECTION_NOT_FOUND_ERR_CODE, nil, collection)
		return
	}

	for _, key := range keyList {
		tmpVal, ok := t.Get(key)
		if ok {
			val = append(val, tmpVal...)
		}
	}

	return
}

func  (r *AutoSyncRepo) GetMapData(collection string) (mapData map[string][]map[string]interface{}, de *domain_error.Error) {
	t, ok := r.tableMaps[collection]
	if !ok {
		de = domain_error.Build(domain_error.COLLECTION_NOT_FOUND_ERR_CODE, nil, collection)
		return
	}
	mapData = t.GetReadMap()
	return
}

func (r *AutoSyncRepo) GetByFuzzyKey(collection string, key string) (val []map[string]interface{}, ok bool, de *domain_error.Error) {
	t, ok := r.tableMaps[collection]
	if !ok {
		de = domain_error.Build(domain_error.COLLECTION_NOT_FOUND_ERR_CODE, nil, collection)
		return
	}

	val, ok = t.GetByFuzzyKey(key)

	return
}

func (r *AutoSyncRepo) GetSyncLogs() (syncLogs[]*schema.DataSyncLog) {
	for k, v := range r.tableMaps {
		log.Printf("======>GetSyncLogs, key=%s\n", k)
		syncLog, ok := v.GetSyncLog()
		if !ok || syncLog.SyncPhase != int(enum.DataSyncLogPhase_Complete) {
			continue
		}
		syncLogs = append(syncLogs, syncLog)
	}
	if syncLogs == nil {
		syncLogs = make([]*schema.DataSyncLog, 0)
	}
	return
}

func  (r *AutoSyncRepo) IsCollectionReady(collection string) bool {
	t, ok := r.tableMaps[collection]
	if !ok {
		return false
	}
	return t.IsReady()
}