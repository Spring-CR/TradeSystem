package order_position

import (
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/quota"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"sync"
	"time"

	"github.com/manucorporat/try"
)

type PositionCalculator struct {
	lock                       *sync.RWMutex                                                               // 读写锁
	applicationCfg             *domain_cfg.ApplicationCfg                                                  // 全局应用配置
	runInTradeEngine           bool                                                                        // 是否运行在TradeEngine内部，运行在order report，将产生更多的计算任务
	positionAdapter            PositionAdapter                                                             // 持仓适配器
	positionMap                map[string]*quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp] // 持仓K-V映射，key由适配器根据order数据结果计算指定
	tradeRespCh                chan *types.TradeActionRespReturn                                           // 监听执行回报的channel
	processPositionChangeEvent func(events []*PositionChangeEvent)                                         // 处理持仓数据变更的函数，在OrderReport要实现的
}

var (
	FillExecType = map[string]bool{
		"F": true,
		"1": true,
		"2": true,
	}
)

func NewPositionCalculator(applicationCfg *domain_cfg.ApplicationCfg, runInTradeEngine bool, processPositionChangeEvent func(events []*PositionChangeEvent)) (pc *PositionCalculator, de *domain_error.Error) {

	var positionAdapter PositionAdapter
	adapterPath := applicationCfg.GetOrdPositionAdapterPath()
	if adapterPath != "" {
		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_positionAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, applicationCfg)
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

	pc = &PositionCalculator{
		lock:                       &sync.RWMutex{},
		applicationCfg:             applicationCfg,
		runInTradeEngine:           runInTradeEngine,
		positionAdapter:            positionAdapter,
		positionMap:                make(map[string]*quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp]),
		tradeRespCh:                make(chan *types.TradeActionRespReturn, 10240),
		processPositionChangeEvent: processPositionChangeEvent,
	}

	pc.startTradeRespLister()
	pc.listeningDataSyncEvent()
	pc.waitForPositionBaseDataReady()

	return
}

func (pc *PositionCalculator) AcquireOrderQuota(force bool, order *schema.TradeOrder) (de *domain_error.Error) {

	positionAdapter := pc.positionAdapter

	if positionAdapter.WouldOrderPositionIncrease(order) {
		return
	}

	keyList, forceList := pc.getPositionKeyList(order)

	for i, positionKey := range keyList {
		force := forceList[i]

		// 如果不是首次请求额度，则忽略；引对重复下单时，发生的重复锁定额度问题
		_, qc, newCreate := pc.getOrCreateQuotaControl(order, positionKey)
		if qc.HasAcquireQuto(order.AppOrdID) {
			continue
		}

		// 非force模式，需要检查是否额度充足
		ok := qc.AcquireQuota(order.AppOrdID, quota.NewQuotaAcquire[*schema.TradeOrder](order, pc.positionAdapter.GetOrderQuota(order)), force, false)
		if !ok {
			de = domain_error.Build(domain_error.QTY_NOT_ENOUGH_ERR_CODE, nil, pc.positionAdapter.GetQtyNotEnoughErrMsgPrefix(order))
			return de
		}

		events := positionAdapter.UpdatePositionAfterAcquireOrderQuota(order, positionKey, qc, pc.runInTradeEngine, newCreate, i)
		if len(events) > 0 {
			pc.processPositionChangeEvent(events)
		}
	}

	return
}

func (pc PositionCalculator) getPositionKeyList(order *schema.TradeOrder) (keyList []string, forceList []bool) {
	positionAdapter := pc.positionAdapter
	keyList, forceList = positionAdapter.GetPositionKeyList(order)
	if len(forceList) < len(keyList) {
		for i := len(forceList); i < len(keyList); i++ {
			forceList = append(forceList, false)
		}
	}
	return
}

func (pc *PositionCalculator) getOrCreateQuotaControl(order *schema.TradeOrder, positionKey string) (string, *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp], bool) {

	//positionAdapter := pc.positionAdapter
	newCreate := false

	pc.lock.RLock()
	qc, ok := pc.positionMap[positionKey]
	newCreate = !ok
	pc.lock.RUnlock()

	if !ok {
		log.Printf("%s not in positionKey, start to GetBaseQuantity\n", positionKey)
		baseQuota, positionBaseRecord := pc.positionAdapter.GetBaseQuantity(order, positionKey)
		log.Printf("positionKey:%s, baseQuota:%v, positionBaseRecord:%v\n", positionKey, baseQuota, positionBaseRecord)

		// 创建qc
		pc.lock.Lock()
		// 在判断一次，确保获得锁，且qc不存在时才创建
		qc, ok = pc.positionMap[positionKey]
		newCreate = !ok
		if !ok {
			qc = quota.NewQuotaControl[*schema.TradeOrder, *schema.TradeActionResp](pc.positionAdapter.GetQuotaMetadata(positionBaseRecord), baseQuota, 10*time.Minute, 2*time.Minute, func(q *quota.QuotaAcquire[*schema.TradeOrder]) (releasedQuota float64) {
				releasedQuota = pc.positionAdapter.ComputeReleaseQuota(q)
				return
			})
			// 加入map
			pc.positionMap[positionKey] = qc
		}
		pc.lock.Unlock()
	}

	return positionKey, qc, newCreate
}

func (pc *PositionCalculator) ReleaseOrderQuota(order *schema.TradeOrder) {

	if pc.positionAdapter.WouldOrderPositionIncrease(order) {
		return
	}

	keyList, _ := pc.getPositionKeyList(order)
	for _, positionKey := range keyList {
		pc.lock.RLock()
		qc, ok := pc.positionMap[positionKey]
		pc.lock.RUnlock()
		if !ok {
			continue
		}
		qc.ReleaseQuota(order.AppOrdID)
	}
}

var (
	ordTradingStatus = map[string]bool{
		string(enum.OrdStatus_PartiallyFilled): true,
		string(enum.OrdStatus_Filled):          true,
	}
)

func (pc *PositionCalculator) updatePositionByTradeResp(tradeResp *types.TradeActionRespReturn) {
	log.Printf("=======> updatePositionByTradeResp, AppOrdID=%s, ClOrdID=%s, OrdStatus=%s, LastShares=%v, LastPx=%v, CumQty=%v, LeavesQty=%v, Symbol=%v\n", tradeResp.GetTradeOrder().AppOrdID, tradeResp.GetTradeOrder().ClOrdID, tradeResp.CurrentTradeActionResp.OrdStatus, tradeResp.CurrentTradeActionResp.LastShares,  tradeResp.CurrentTradeActionResp.LastPx,  tradeResp.CurrentTradeActionResp.CumQty,  tradeResp.CurrentTradeActionResp.LeavesQty,  tradeResp.GetTradeOrder().Symbol)
	order := tradeResp.TraceableTradeActionResp.GetTradeOrder()
	keyList, _ := pc.getPositionKeyList(order)
	positionAdapter := pc.positionAdapter
	for i, positionKey := range keyList {
		positionKey, qc, newCreate := pc.getOrCreateQuotaControl(order, positionKey)
		v1, v2 := qc.GetQuota()
		wouldOrderPositionDecrease := positionAdapter.WouldOrderPositionDecrease(order)
		isTradeActionRespInBeginStatus := positionAdapter.IsTradeActionRespInBeginStatus(tradeResp.CurrentTradeActionResp)
		log.Printf("positionKey: %s, newCreate:%v, baseQuota:%v, quota:%v, symbol:%v, order.appOrdID:%v, order.clOrdID:%v, OrdStatus:%v, WouldOrderPositionDecrease:%v, isTradeActionRespInBeginStatus:%v\n", positionKey, newCreate, v1, v2, order.Symbol, order.AppOrdID, order.ClOrdID, tradeResp.CurrentTradeActionResp.OrdStatus, wouldOrderPositionDecrease, isTradeActionRespInBeginStatus)
		if wouldOrderPositionDecrease {
			// 卖单，持仓数量将减少
			//if newCreate {
			if isTradeActionRespInBeginStatus && !qc.HasAcquireQuto(order.AppOrdID) { // 只有在orderCache才运行，否则，将会重复reqirequota【但是重启是就会有问题！就少了reqirequota】
				orderQuota := pc.positionAdapter.GetOrderQuota(order)
				log.Printf("AcquireQuota by TradeActionRespReturn, appOrdID:=%s, clOrdID=%s, ordStatus=%v, orderQuota=%v, positionKey=%v\n", tradeResp.GetTradeOrder().AppOrdID, tradeResp.GetTradeOrder().ClOrdID, tradeResp.CurrentTradeActionResp.OrdStatus, orderQuota, positionKey)
				qc.AcquireQuota(order.AppOrdID, quota.NewQuotaAcquire[*schema.TradeOrder](order, orderQuota), true, false)
			} else {
				// 判断是否终态
				log.Printf("===> ordStatus:%v\n", tradeResp.CurrentTradeActionResp.OrdStatus)
				// 多头平仓，部分成交时，返回的额度的规则：
				// 1、T0: T0、T1的持仓都返回增加
				// 2、T1: T1持仓直接返回增加，T0=实际冻结数量 - min(实际冻结数量,成交数量)
				if positionAdapter.IsOrderFinished(tradeResp.CurrentTradeActionResp) {
					log.Printf("===> order is finished, go to ReleaseQuota")
					qc.ReleaseQuota(order.AppOrdID)
				}
			}

			events := positionAdapter.AfterUpdatePositionByTradeResp(order, tradeResp, positionKey, qc, true, pc.runInTradeEngine, newCreate, i)
			log.Printf("event1.Len=%d\n", len(events))
			if len(events) > 0 && pc.processPositionChangeEvent != nil {
				pc.processPositionChangeEvent(events)
			}
		} else {
			// 买单，持仓数量将增加
			if (ordTradingStatus[tradeResp.CurrentTradeActionResp.OrdStatus] || FillExecType[tradeResp.CurrentTradeActionResp.ExecType]) && pc.positionAdapter.CouldIncreaseQuota(order, tradeResp.CurrentTradeActionResp, positionKey) {
				qc.ReturnQuota(tradeResp.CurrentTradeActionResp.GetCacheKey(), quota.NewQuotaReturn[*schema.TradeActionResp](tradeResp.CurrentTradeActionResp, float64(tradeResp.CurrentTradeActionResp.LastShares)))
			}

			events := positionAdapter.AfterUpdatePositionByTradeResp(order, tradeResp, positionKey, qc, false, pc.runInTradeEngine, newCreate, i)
			log.Printf("event2.Len=%d\n", len(events))
			if len(events) > 0 && pc.processPositionChangeEvent != nil {
				pc.processPositionChangeEvent(events)
			}
		}
	}
	log.Println("<======= updatePositionByTradeResp")
}

func (pc *PositionCalculator) GetTradeRespCh() chan *types.TradeActionRespReturn {
	return pc.tradeRespCh
}

func (pc *PositionCalculator) startTradeRespLister() {
	go func() {
		for {
			tradeResp, ok := <-pc.tradeRespCh
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, errors.New("tradeRespCh in position-calculator closed"), "tradeRespCh in position-calculator closed")
				if tradeResp.WaitGroup != nil {
					tradeResp.WaitGroup.Done()
				}
				break
			}
			try.This(func() {
				pc.updatePositionByTradeResp(tradeResp)
			}).Catch(func(err try.E) {
				log.Printf("error occur while run updatePositionByTradeResp! error:%v\n", err)
				de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while run updatePositionByTradeResp! error:%v", err))
				domain_error.ProcessSevereError(false, 0, de, nil, fmt.Sprintf("error occur while run updatePositionByTradeResp! error:%v\n", err))
			})
			if tradeResp.WaitGroup != nil {
				tradeResp.WaitGroup.Done()
			}
		}
	}()
}

func (pc *PositionCalculator) listeningDataSyncEvent() {
	go func() {
		evtCh := pc.applicationCfg.GetDataSyncEventChan()
		log.Printf("start to listeningDataSyncEvent")
		for {
			evt, ok := <-evtCh
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, errors.New("DataSyncEventChan closed"), "DataSyncEventChan closed")
				break
			}
			events := pc.positionAdapter.ProcessDataSyncEvent(evt, pc.lock, pc.positionMap)
			if len(events) == 0 {
				continue
			}
			if pc.processPositionChangeEvent != nil {
				if len(events) > 0 {
					pc.processPositionChangeEvent(events)
				}
			}
		}
	}()
}

// Todo, 1、恢复的时候，要先执行买单，再执行卖单
//       2、数据上场载入内存成功之后，才能执行系统恢复
// func (pc *PositionCalculator) AdjPositionAfterSystemRecover() {

// }

// Todo, 1、因为持仓计算设计PositionBase，因此要等到数据上场载入内存成功之后，才能执行系统恢复
// 等待PositionBase数据就绪
func (pc *PositionCalculator) waitForPositionBaseDataReady() {
	log.Println("waitForPositionBaseDataReady...")
	autoSyncRepo := pc.applicationCfg.GetAutoSyncRepo()
	log.Printf("get autoSyncRepo:%v\n", autoSyncRepo)
	/*m, de := autoSyncRepo.GetMapData("PositionBase")
	if de !=nil{
		domain_error.ProcessSevereError(false, 0, de, nil, "PositionBase fail")
	}
	for len(m) == 0 {
		time.Sleep(10 * time.Second)
		m, de = autoSyncRepo.GetMapData("PositionBase")
		if de !=nil{
			domain_error.ProcessSevereError(false, 0, de, nil, "get PositionBase fail")
		}
		log.Printf("get PositionBase, m=%d\n", len(m))
	}*/
	for !autoSyncRepo.IsCollectionReady("PositionBase") {
		log.Printf("PositionBase is not ready, go to sleep 10 seconds...")
		time.Sleep(10 * time.Second)
	}
	m, de := autoSyncRepo.GetMapData("PositionBase")
	if de !=nil{
		domain_error.ProcessSevereError(false, 0, de, nil, "PositionBase fail")
	}
	log.Printf("get PositionBase, m=%d\n", len(m))
	
	log.Println("finish waitForPositionBaseDataReady")
}
