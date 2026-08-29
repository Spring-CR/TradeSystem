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
				pm.updatePositionByTradeResp(tradeResp)
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
	pm.lock.RLock()
	pc, ok := pm.calculatorMap[key]
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
	log.Printf("get autoSyncRepo:%v\n", autoSyncRepo)

	for !autoSyncRepo.IsCollectionReady("PositionBase") {
		log.Printf("PositionBase is not ready, go to sleep 10 seconds...")
		time.Sleep(10 * time.Second)
	}
	m, de := autoSyncRepo.GetMapData("PositionBase")
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, nil, "PositionBase fail")
	}
	log.Printf("get PositionBase, m=%d\n", len(m))

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
			_, ok := <-evtCh
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, errors.New("DataSyncEventChan closed"), "DataSyncEventChan closed")
				break
			}
			// Todo， 暂时没啥用处
		}
	}()
}
