package order_position_manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-common/utils/logger"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"strings"
	"sync"
	"time"

	"github.com/manucorporat/try"
)

type PositionManager struct {
	lock             *sync.RWMutex
	applicationCfg   *domain_cfg.ApplicationCfg // 全局应用配置
	runInTradeEngine bool
	calculatorMap    map[string]*PositionCalculator
	positionAdapter  PositionAdapter                   // 持仓适配器
	tradeRespCh      chan *types.TradeActionRespReturn // 监听执行回报的channel
	orderLog         *logger.OrderLog
	persisFunc       func(evt *PositionChangeEvent)
	mockTradeOrders  map[string][]*schema.TradeOrder // 用于调仓的虚拟订单
}

func NewPositionManager(applicationCfg *domain_cfg.ApplicationCfg, runInTradeEngine bool, persisFunc func(evt *PositionChangeEvent)) (*PositionManager, *domain_error.Error) {
	var positionAdapter PositionAdapter
	adapterPath := applicationCfg.GetOrdPositionAdapterPath()
	os.RemoveAll("/opt/log/position/")

	var orderLog *logger.OrderLog
	//Todo: 后面改成只在order-report中启用日志
	orderLog = logger.NewOrderLog("/opt/log/position/", 4096, 10*time.Second, func(tradeOrder *schema.TradeOrder) (key string) {
		key = tradeOrder.Account + "_" + tradeOrder.Symbol
		log.Printf("===> key for order[%s] %s: %s\n", tradeOrder.Symbol, tradeOrder.AppOrdID, key)
		return
	}, !runInTradeEngine)

	if adapterPath != "" {
		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_positionAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg, orderLog)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct PositionAdapter from "+adapterPath)
		}
		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct PositionAdapter from "+adapterPath)
		}
		// 获取apiAdapter
		positionAdapter = _positionAdapter.(PositionAdapter)
		log.Printf("finish get PositionAdapter:%s\n", adapterPath)
	}

	pm := &PositionManager{
		lock:             &sync.RWMutex{},
		applicationCfg:   applicationCfg,
		runInTradeEngine: runInTradeEngine,
		calculatorMap:    make(map[string]*PositionCalculator),
		positionAdapter:  positionAdapter,
		tradeRespCh:      make(chan *types.TradeActionRespReturn, 10240),
		orderLog:         orderLog,
		persisFunc:       persisFunc,
		mockTradeOrders:  make(map[string][]*schema.TradeOrder),
	}

	pm.startTradeRespLister()
	pm.listeningDataSyncEvent()
	pm.waitForPositionBaseDataReady()

	return pm, nil
}

func (pm *PositionManager) GetTradeRespCh() chan *types.TradeActionRespReturn {
	return pm.tradeRespCh
}

func (pm *PositionManager) startTradeRespLister() {
	log.Printf("PositionManager.startTradeRespLister...")
	go func() {
		for {
			tradeResp, ok := <-pm.tradeRespCh
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, errors.New("tradeRespCh in position-manager closed"), "tradeRespCh in position-manager closed")
				if tradeResp.WaitGroup != nil {
					log.Println("WaitGroup.Done() in PositionManager1")
					tradeResp.WaitGroup.Done()
				}
				break
			}
			try.This(func() {
				log.Printf("start to updatePositionByTradeResp1")
				// 20260828发现这里阻塞！！！
				pm.updatePositionByTradeResp(tradeResp)
				log.Printf("finish updatePositionByTradeResp1")
			}).Catch(func(err try.E) {
				log.Printf("error occur while run updatePositionByTradeResp! error:%v\n", err)
				de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while run updatePositionByTradeResp! error:%v", err))
				domain_error.ProcessSevereError(false, 0, de, nil, fmt.Sprintf("error occur while run updatePositionByTradeResp! error:%v\n", err))
			})
			if tradeResp.WaitGroup != nil {
				log.Println("WaitGroup.Done() in PositionManager2")
				tradeResp.WaitGroup.Done()
			}
		}
	}()
}

func (pm *PositionManager) GetQuotaNotEnoughHandler() func(tradeOrder *schema.TradeOrder, de *domain_error.Error) {
	return pm.positionAdapter.GetQuotaNotEnoughHandler()
}

func (pm *PositionManager) HasSufficientQuota(tradeOrder *schema.TradeOrder) (sufficient bool, de *domain_error.Error) {
	key := pm.positionAdapter.GetPositionCalculatorKey(tradeOrder)
	log.Printf("HasSufficientQuota, tradeOrder.AppOrdID=%s, key=%v\n", tradeOrder.AppOrdID, key)
	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
	pm.lock.RUnlock()
	if ok {
		return pc.HasSufficientQuota(tradeOrder)
	}
	return
}

func (pm *PositionManager) FreezeQuota(force bool, tradeOrder *schema.TradeOrder) (sufficient bool, de *domain_error.Error) {
	key := pm.positionAdapter.GetPositionCalculatorKey(tradeOrder)
	log.Printf("FreezeQuota, force=%v, tradeOrder.AppOrdID=%s, key=%v\n", force, tradeOrder.AppOrdID, key)
	pm.orderLog.Printf(tradeOrder, nil, "[FreezeQuota] force=%v, key=%v", force, key)

	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
	pm.lock.RUnlock()
	if !ok {
		pm.lock.Lock()
		pc, ok = pm.calculatorMap[key]
		if !ok {
			pc = NewPositionCalculator(key, tradeOrder, pm.positionAdapter, pm.orderLog, pm.persisFunc)
			pm.calculatorMap[key] = pc
		}
		pm.lock.Unlock()
	}
	return pc.FreezeQuota(force, tradeOrder)
}

func (pm *PositionManager) RollbackFreezeQuota(tradeOrder *schema.TradeOrder) {
	key := pm.positionAdapter.GetPositionCalculatorKey(tradeOrder)
	log.Printf("RollbackFreezeQuota, tradeOrder.AppOrdID=%s, key=%v\n", tradeOrder.AppOrdID, key)
	pm.orderLog.Printf(tradeOrder, nil, "[RollbackFreezeQuota] key=%v", key)

	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
	pm.lock.RUnlock()
	if !ok {
		return
	}
	pc.RollbackFreezeQuota(tradeOrder)
}

func (pm *PositionManager) updatePositionByTradeResp(tradeResp *types.TradeActionRespReturn) {
	tradeOrder := tradeResp.GetTradeOrder()
	key := pm.positionAdapter.GetPositionCalculatorKey(tradeOrder)
	log.Printf("start to GetPositionCalculatorKey by key:%s\n", key)
	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
	log.Printf("get pc, ok:%v\n", ok)
	pm.lock.RUnlock()
	if !ok {
		return
	}
	
	pc.UpdatePositionByTradeResp(tradeResp)
}

// Todo, 1、因为持仓计算设计PositionBase，因此要等到数据上场载入内存成功之后，才能执行系统恢复
// 等待PositionBase数据就绪
func (pm *PositionManager) waitForPositionBaseDataReady() {
	log.Println("waitForPositionBaseDataReady...")
	autoSyncRepo := pm.applicationCfg.GetAutoSyncRepo()
	// 更通用，不用指定PositionBase数据集的具体名字
	for !autoSyncRepo.IsAllPositionBaseCollectionsReady() {
		log.Printf("PositionBase collection is not ready, go to sleep 10 seconds...")
		time.Sleep(10 * time.Second)
	}

	params := pm.positionAdapter.LoadInitPositionRecords(pm.applicationCfg)
	for _, param := range params {
		pc := NewPositionCalculatorByConstructParam(param, pm.positionAdapter, pm.orderLog, pm.persisFunc)
		pm.calculatorMap[param.Key] = pc

		if pm.persisFunc != nil {
			positionData := pm.positionAdapter.GeneralizePositionRecord(param.Metadata, true)
			pm.persisFunc(&PositionChangeEvent{InsertOrUpdate: 0, PositionData: positionData})
		}
	}

	log.Println("finish waitForPositionBaseDataReady")
}

func (pm *PositionManager) ReloadPositionRecordsForPurgingTask(purgingLog *schema.DataPurgingLog) {
	log.Printf("start ReloadPositionRecordsForPurgingTask")
	params := pm.positionAdapter.ReloadPositionRecordsForPurgingTask(pm.applicationCfg, purgingLog)
	for _, param := range params {
		pc := NewPositionCalculatorByConstructParam(param, pm.positionAdapter, pm.orderLog, pm.persisFunc)
		pm.calculatorMap[param.Key] = pc

		if pm.persisFunc != nil {
			positionData := pm.positionAdapter.GeneralizePositionRecord(param.Metadata, true)
			pm.persisFunc(&PositionChangeEvent{InsertOrUpdate: 0, PositionData: positionData})
		}
	}
	log.Printf("finish ReloadPositionRecordsForPurgingTask, %d records\n", len(params))
}

func (pm *PositionManager) Dump() map[string]interface{} {
	m := make(map[string]interface{})
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	log.Printf("len(calculatorMap):%d\n", len(pm.calculatorMap))
	for k, c := range pm.calculatorMap {
		log.Printf("key=%s\n", k)
		m[k] = c.metadata
	}
	js, err := json.MarshalIndent(m, "", "  ")
	log.Printf("position data:%s, err:%v\n", js, err)
	return m
}

func (pm *PositionManager) listeningDataSyncEvent() {
	go func() {
		evtCh := pm.applicationCfg.GetDataSyncEventChan()
		log.Printf("start to listeningDataSyncEvent")
		for {
			event, ok := <-evtCh
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, errors.New("DataSyncEventChan closed"), "DataSyncEventChan closed")
				time.Sleep(10 * time.Second)
				continue
			}

			if !pm.positionAdapter.UpdatePositionBaseDynamically() {
				return
			}

			if !strings.HasPrefix(event.TableAlias, "PositionBase") {
				continue
			}
			//	pm.positionAdapter.ProcessPositionBaseDataSyncEvent(event, pm.calculatorMap)
			keyList := event.AddKeys
			keyList = append(keyList, event.ChgKeys...)
			visiMap := make(map[string]bool)
			for _, key := range keyList {
				valList, ok, _ := pm.applicationCfg.GetAutoSyncRepo().Get(event.TableAlias, key)
				if !ok || len(valList) == 0 {
					continue
				}
				positionData := valList[0]
				key, positionRecordBase := pm.positionAdapter.ParsePositionRecordFromReposRecord(positionData)
				if key == "" {
					continue
				}
				if visiMap[key] {
					continue
				}
				pm.lock.Lock()

				//positionRecordCurr, ok := pm.calculatorMap[key]
				calculatorParent, ok := pm.calculatorMap[key]
				
				if !ok {

					param := pm.positionAdapter.GetPositionCalculatorConstructParamForPositionRecord(positionRecordBase)
					// 直接插入
					pc := NewPositionCalculatorByConstructParam(param, pm.positionAdapter, pm.orderLog, pm.persisFunc)
					pm.calculatorMap[param.Key] = pc

					if pm.persisFunc != nil {
						positionData := pm.positionAdapter.GeneralizePositionRecord(param.Metadata, true)
						pm.persisFunc(&PositionChangeEvent{InsertOrUpdate: 0, PositionData: positionData})
					}

				} else {

					//positionRecordCurr := calculatorParent.metadata

					// // 调仓，构造调仓参数
					// mockTradeOrder, mockTradeActionResp, de := pm.positionAdapter.PreparePositionAdjustmentParamsByPositionBaseDiff(positionRecordBase, positionRecordCurr)
					// if de != nil {
					// 	log.Println("fail to PreparePositionAdjustmentParamsByPositionBaseDiff")
					// 	continue
					// }
					// if mockTradeOrder == nil || mockTradeActionResp == nil {
					// 	log.Println("cannot get mockTradeOrder and mockTradeActionResp")
					// 	continue
					// }
					// pm.AdjustPosition(mockTradeOrder, mockTradeActionResp)

					// 清理现场，准备重演
					pm.positionAdapter.PrepareForRecover(positionRecordBase)

					param := pm.positionAdapter.GetPositionCalculatorConstructParamForPositionRecord(positionRecordBase)

					calculatorParent.lock.Lock()
					
					successor := NewPositionCalculatorByConstructParam(param, pm.positionAdapter, pm.orderLog, pm.persisFunc)
					calculatorParent.SetSuccessor(successor)

					pm.calculatorMap[param.Key] = successor
					
					if pm.persisFunc != nil {
						positionData := pm.positionAdapter.GeneralizePositionRecord(param.Metadata, true)
						pm.persisFunc(&PositionChangeEvent{InsertOrUpdate: 0, PositionData: positionData})
					}

					// 重演开始
					log.Printf("开始重演：%v\n", key)
					successor.RecoverFromParent(calculatorParent)
					log.Printf("结束重演：%v\n", key)

					calculatorParent.lock.Unlock()
				}

				pm.lock.Unlock()
			}
		}
	}()
}

func (pm *PositionManager) GetCachePositionRecord() {

}
