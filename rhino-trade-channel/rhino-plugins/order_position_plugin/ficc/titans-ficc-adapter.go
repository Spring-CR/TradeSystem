package ficc

import (
	"fmt"
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/quota"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-data/datamap"
	statusficc "rhino-plugins/order_status_plugin/ficc"
	"strconv"
	"strings"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

var (
	InitMarginRatio  = 0.05
	DirtyPriceRation = 1.05
)

type TitansFiccOrderPositionAdapter struct {
	applicationCfg *domain_cfg.ApplicationCfg
}

func NewTitansFiccOrderPositionAdapter(applicationCfg *domain_cfg.ApplicationCfg) (adapter *TitansFiccOrderPositionAdapter, de *domain_error.Error) {
	log.Printf("construct TitansFiccOrderPositiondapter...")

	// 从配置文件提取初保比率
	configMap := applicationCfg.GetApplicationCfgItemMap()
	configItem, ok := configMap["InitMarginRatio"]
	if ok && configItem.ConfigItemValue != "" {
		val, err := strconv.ParseFloat(configItem.ConfigItemValue, 64)
		if err == nil && val > 0 {
			InitMarginRatio = val / 100.0
			log.Printf("Reset InitMarginRatio=%f\n", InitMarginRatio)
		}
	}

	adapter = &TitansFiccOrderPositionAdapter{applicationCfg: applicationCfg}
	NewCapitalControl(applicationCfg)
	return
}

func (a *TitansFiccOrderPositionAdapter) GetPositionKeyList(order *schema.TradeOrder) (keyList []string, forceList []bool) {

	counterpartyID, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "account", enum.AttrValueType_INT)
	planCode, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "planCode", enum.AttrValueType_STRING)
	symbol, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "symbol", enum.AttrValueType_STRING)
	side, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "side", enum.AttrValueType_STRING)
	if side.(string) == "" || side.(string) == "1" {
		side = "LONG"
	} else {
		side = "SHORT"
	}
	_settlType, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "settlType", enum.AttrValueType_STRING)
	settlType := strings.ToUpper(strings.ReplaceAll(_settlType.(string), "+", ""))

	keyList = append(keyList, fmt.Sprintf("%v-%v-%v-%v-%v", counterpartyID, planCode, symbol, side, settlType))
	if settlType == "T0" {
		settlType = "T1"
	} else {
		settlType = "T0"
	}
	// 反转settlType，再加一次
	keyList = append(keyList, fmt.Sprintf("%v-%v-%v-%v-%v", counterpartyID, planCode, symbol, side, settlType))

	forceList = []bool{false, true}

	return
}

func (a *TitansFiccOrderPositionAdapter) WouldOrderPositionIncrease(order *schema.TradeOrder) bool {
	return order.OpenClose == "O"
}

func (a *TitansFiccOrderPositionAdapter) WouldOrderPositionDecrease(order *schema.TradeOrder) bool {
	return order.OpenClose == "C"
}

func (a *TitansFiccOrderPositionAdapter) IsOrderFinished(tradeActionResp *schema.TradeActionResp) bool {
	return statusficc.EndStatus[tradeActionResp.OrdStatus]
}

func (a *TitansFiccOrderPositionAdapter) IsTradeActionRespInBeginStatus(tradeActionResp *schema.TradeActionResp) bool {
	return tradeActionResp.OrdStatus == string(enum.OrdStatus_PendingNew)
}

func (a *TitansFiccOrderPositionAdapter) GetQtyNotEnoughErrMsgPrefix(order *schema.TradeOrder) string {
	settlType, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "settlType", enum.AttrValueType_STRING)
	return fmt.Sprintf("%v可用面额", settlType)
}

func (a *TitansFiccOrderPositionAdapter) GetQuotaMetadata(positionBaseRecord map[string]interface{}) map[string]interface{} {

	metadata := make(map[string]interface{})
	metadata["account"] = positionBaseRecord["CounterpartyID"]
	metadata["counterpartyID"] = positionBaseRecord["CounterpartyID"]
	metadata["counterparty"] = positionBaseRecord["Counterparty"]
	metadata["symbol"] = positionBaseRecord["Symbol"]
	metadata["symbolName"] = positionBaseRecord["SecurityName"]
	metadata["currency"] = positionBaseRecord["Currency"]
	metadata["planCode"] = positionBaseRecord["PlanCode"]
	metadata["ultraContractCode"] = positionBaseRecord["UltraContractCode"]
	metadata["securityExchange"] = positionBaseRecord["SecurityExchange"]
	metadata["securityType"] = positionBaseRecord["SecurityType"]
	metadata["longShort"] = positionBaseRecord["LongShort"]
	metadata["baseCashQty"] = positionBaseRecord["BaseCashQty"]
	metadata["totalCashQty"] = positionBaseRecord["BaseCashQty"]
	metadata["baseNotional"] = positionBaseRecord["BaseNotional"]
	metadata["totalNotional"] = positionBaseRecord["BaseNotional"]
	// 新增全价总资金成本和净价总资金成本
	metadata["baseDirtyCost"] = positionBaseRecord["BaseDirtyCost"]
	metadata["totalDirtyCost"] = positionBaseRecord["BaseDirtyCost"]
	metadata["baseCleanCost"] = positionBaseRecord["BaseCleanCost"]
	metadata["totalCleanCost"] = positionBaseRecord["BaseCleanCost"]

	parValue, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "ParValue", enum.AttrValueType_FLOAT)
	if !ok {
		jsData, _ := json.Marshal(positionBaseRecord)
		domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("ParValue not fond in %s", jsData), fmt.Sprintf("ParValue not fond in %s\n", jsData))
	}
	if parValue.(float64) <= 0 {
		parValue = 100.0
	}
	metadata["parValue"] = parValue
	calculateAverageCost(metadata)
	return metadata
}

func getDirtyPriceByOrder(order *schema.TradeOrder) float64 {
	// 后面累总名义本金，需要用到订单价格参数
	price, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "price", enum.AttrValueType_FLOAT)
	dirtyPrice, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	if dirtyPrice.(float64) <= 0 {
		dirtyPrice = price.(float64) * DirtyPriceRation
	}
	return dirtyPrice.(float64)
}

func getPriceByOrder(order *schema.TradeOrder) (float64, float64) {
	// 后面累总名义本金，需要用到订单价格参数
	price, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "price", enum.AttrValueType_FLOAT)
	dirtyPrice, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "dirtyPrice", enum.AttrValueType_FLOAT)
	if dirtyPrice.(float64) <= 0 {
		dirtyPrice = price.(float64) * DirtyPriceRation
	}
	return price.(float64), dirtyPrice.(float64)
}

func getParValueByOrder(order *schema.TradeOrder) float64 {
	// 后面累总名义本金，需要用到订单价格参数
	parValue, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "parValue", enum.AttrValueType_FLOAT)
	if parValue == 0 {
		parValue = 100.0
	}
	return parValue.(float64)
}

func getCapitalAcctIDByOrder(order *schema.TradeOrder) int {
	capitalAcctID, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "capitalAcctID", enum.AttrValueType_INT)
	return capitalAcctID.(int)
}

func getCurrencyByOrder(order *schema.TradeOrder) string {
	currency, _, _ := attrutil.GetAttrValue(order.ExtendAttrMap, "currency", enum.AttrValueType_STRING)
	return currency.(string)
}

func calculateAverageCost(metadata map[string]interface{}) {
	totalCashQty, _ := metadata["totalCashQty"].(float64)
	if totalCashQty <= 0 {
		metadata["averageCost"] = 0.0
	} else {
		parValue, _ := metadata["parValue"].(float64)
		if parValue == 0 {
			parValue = 100.0
		}

		/* 从基于名义本金，改成基于全价、净价的总资本成本
		totalNotional, _ := metadata["totalNotional"].(float64)
		metadata["averageCost"] = totalNotional / (totalCashQty / parValue)
		*/
		amount := totalCashQty / parValue
		totalDirtyCost, _ := metadata["totalDirtyCost"].(float64)
		metadata["averageCost"] = totalDirtyCost / amount
		totalCleanCost, _ := metadata["totalCleanCost"].(float64)
		metadata["averageCleanCost"] = totalCleanCost / amount
	}
}

func (a *TitansFiccOrderPositionAdapter) ComputeReleaseQuota(quotaAcquire *quota.QuotaAcquire[*schema.TradeOrder]) (releaseQuota float64) {

	order := quotaAcquire.GetSource()
	ordStatus := order.OrdStatus

	if statusficc.DuringTradingStatus[ordStatus] {
		return
	}

	if statusficc.EndStatus[order.OrdStatus] {
		// 释放数量 = 实际冻结数量 - min(实际冻结数量,成交数量)
		acquiredQuota := quotaAcquire.GetAcquiredQuota()
		releaseQuota := acquiredQuota - math.Min(acquiredQuota, float64(order.CumQty))
		log.Printf("selling order %s is finished, release quota %v\n", order.ClOrdID, releaseQuota)
		return releaseQuota
	}

	releaseQuota = order.OrderQty
	log.Printf("selling order %s is timeout for entering trading phase, go to release quota %v\n", order.ClOrdID, releaseQuota)

	return
}

func (a *TitansFiccOrderPositionAdapter) GetBaseQuantity(order *schema.TradeOrder, positionKey string) (baseQuota float64, positionBaseRecord map[string]interface{}) {

	valList, _, _ := a.applicationCfg.GetAutoSyncRepo().Get("PositionBase", positionKey)
	if len(valList) > 0 {
		positionBaseRecord = valList[len(valList)-1]
		// 获取初始底仓
		val, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "Quantity", enum.AttrValueType_FLOAT)
		if ok {
			baseQuota = val.(float64)
		}
	} else {
		log.Printf("cannot GetBaseQuantity for %s\n, try to create one", positionKey)

		positionBaseRecord = map[string]interface{}{
			"CounterpartyID":    order.ExtendAttrMap["account"],
			"Counterparty":      order.ExtendAttrMap["counterparty"],
			"Symbol":            order.ExtendAttrMap["symbol"],
			"SecurityName":      order.ExtendAttrMap["symbolName"],
			"Currency":          order.ExtendAttrMap["currency"],
			"PlanCode":          order.ExtendAttrMap["planCode"],
			"UltraContractCode": order.ExtendAttrMap["ultraContractCode"],
			"SecurityExchange":  order.ExtendAttrMap["securityExchange"],
			"SecurityType":      "BOND",
			"LongShort":         order.ExtendAttrMap["side"],
			"BaseCashQty":       0.0,
			"BaseNotional":      0.0,
			"BaseCleanCost":     0.0,
			"BaseDirtyCost":     0.0,
			"ParValue":          order.ExtendAttrMap["parValue"],
		}

		if positionBaseRecord["SecurityExchange"] == "" || positionBaseRecord["SecurityExchange"] == "8" {
			positionBaseRecord["SecurityExchange"] = "NIB"
		}
		if positionBaseRecord["LongShort"] == "" || positionBaseRecord["LongShort"] == "1" {
			positionBaseRecord["LongShort"] = "LONG"
		} else {
			positionBaseRecord["LongShort"] = "SHORT"
		}
	}

	positionBaseRecordJson, _ := json.Marshal(positionBaseRecord)
	log.Printf("in GetBaseQuantity, positionBaseRecordJson:%s\n", positionBaseRecordJson)

	parValue, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "ParValue", enum.AttrValueType_FLOAT)
	if !ok {
		jsData, _ := json.Marshal(positionBaseRecord)
		domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("ParValue not fond in %s", jsData), fmt.Sprintf("ParValue not fond in %s\n", jsData))
	}
	if parValue.(float64) <= 0 {
		parValue = 100.0
	}

	baseQuota = parValue.(float64) * baseQuota
	positionBaseRecord["Quantity"] = baseQuota

	log.Printf("parValue:%v, baseQuota for %s:%v\n", parValue, positionKey, baseQuota)

	return
}

// 如果是T0订单的成交回报的，T0、T1的持仓都需要增加；如果是T1订单的成交回报，只有T1的持仓需要增加
func (a *TitansFiccOrderPositionAdapter) CouldIncreaseQuota(order *schema.TradeOrder, tradeActionResp *schema.TradeActionResp, positionKey string) bool {
	val := order.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)
	if settleType == "T+1" && strings.Contains(positionKey, "T0") {
		return false
	}
	return true
}

func (a *TitansFiccOrderPositionAdapter) GetOrderQuota(order *schema.TradeOrder) float64 {
	return order.OrderQty
}

func (a *TitansFiccOrderPositionAdapter) UpdatePositionAfterAcquireOrderQuota(order *schema.TradeOrder, positionKey string, qc *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp], runInTradeEngine bool, newCreate bool, keyIndex int) (events []*order_position.PositionChangeEvent) {
	if a.WouldOrderPositionIncrease(order) || runInTradeEngine {
		return
	}
	log.Printf("UpdatePositionAfterAcquireOrderQuota:: order.AppOrdID:%v, positionKey:%v, runInTradeEngine:%v, newCreate:%v, keyIndex:%v\n", order.AppOrdID, positionKey, runInTradeEngine, newCreate, keyIndex)
	isT0 := strings.Contains(positionKey, "T0")
	insertOrUpdat := 0
	if newCreate { // 新建position记录
		event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InPositionMap: true}

		metadata, lock := qc.GetMetadata()

		lock.RLock()
		// 拷贝源数据
		for k, v := range metadata {
			if k != "" {
				event.PositionData[k] = v
			}
		}
		lock.RUnlock()

		baseQuota, quota := qc.GetQuota()
		// 设置T0/T1的持仓

		if isT0 {
			event.PositionData["baseQuotaT0"] = baseQuota
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["baseQuotaT1"] = baseQuota
			event.PositionData["quotaT1"] = quota
		}

		if keyIndex == 0 {
			// 真正的新建
			if isT0 {
				event.PositionData["baseQuotaT1"] = 0.0
				event.PositionData["quotaT1"] = 0.0
			} else {
				event.PositionData["baseQuotaT0"] = 0.0
				event.PositionData["quotaT0"] = 0.0
			}
		} else {
			// 要转为update
			insertOrUpdat = 1
		}
		event.InsertOrUpdate = insertOrUpdat
		configEventKey(event)
		events = append(events, event)

		postionMetaJs, _ := json.Marshal(event.PositionData)
		log.Printf("postionMetaJs_1: %s\n", postionMetaJs)

	} else { // 更新
		insertOrUpdat = 1
		event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InsertOrUpdate: insertOrUpdat, InPositionMap: true}

		metadata, lock := qc.GetMetadata()
		lock.Lock()

		setKeyMetadata(event.PositionData, metadata)

		_, quota := qc.GetQuota()
		// 设置T0/T1的持仓
		isT0 := strings.Contains(positionKey, "T0")
		if isT0 {
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["quotaT1"] = quota
		}

		lock.Unlock()

		configEventKey(event)
		
		events = append(events, event)

		postionMetaJs, _ := json.Marshal(event.PositionData)
		log.Printf("postionMetaJs_2: %s\n", postionMetaJs)
	}
	return
}

func (a *TitansFiccOrderPositionAdapter) AfterUpdatePositionByTradeResp(
	order *schema.TradeOrder, tradeResp *types.TradeActionRespReturn,
	positionKey string, qc *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp],
	isOrderPositionDecrease bool, runInTradeEngine bool, newCreate bool,
	keyIndex int) (events []*order_position.PositionChangeEvent) {

	log.Println("Enter AfterUpdatePositionByTradeResp..")

	// 如果是跑在交易引擎，不用计算，直接退出
	//if runInTradeEngine {
	//	return
	//}
	// 改成在交易引擎之中也要计算，因为验资的逻辑里涉及到了持仓均价，这里要计算才能获取实时的持仓均价

	insertOrUpdat := 0
	isT0 := strings.Contains(positionKey, "T0")
	// 如果是新建的 或者 跑在order_report&&平仓单&&PendingNew状态
	if newCreate /*|| (!runInTradeEngine && isOrderPositionDecrease && a.IsTradeActionRespInBeginStatus(tradeResp.CurrentTradeActionResp))*/ {

		event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InPositionMap: true}

		metadata, lock := qc.GetMetadata()

		postionMetaJs, _ := json.Marshal(metadata)
		log.Printf("postionMetaJs: %s\n", postionMetaJs)

		lock.RLock()
		// 拷贝源数据
		for k, v := range metadata {
			if k != "" {
				event.PositionData[k] = v
			}
		}
		lock.RUnlock()

		baseQuota, quota := qc.GetQuota()
		// 设置T0/T1的持仓

		if isT0 {
			event.PositionData["baseQuotaT0"] = baseQuota
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["baseQuotaT1"] = baseQuota
			event.PositionData["quotaT1"] = quota
		}

		if keyIndex == 0 {
			// 真正的新建
			if isT0 {
				event.PositionData["baseQuotaT1"] = 0.0
				event.PositionData["quotaT1"] = 0.0
			} else {
				event.PositionData["baseQuotaT0"] = 0.0
				event.PositionData["quotaT0"] = 0.0
			}
		} else {
			// 要转为update
			insertOrUpdat = 1
		}
		event.InsertOrUpdate = insertOrUpdat
		configEventKey(event)
		events = append(events, event)
	} else {
		insertOrUpdat = 1

		if !runInTradeEngine && isOrderPositionDecrease && a.IsTradeActionRespInBeginStatus(tradeResp.CurrentTradeActionResp) {
			event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InsertOrUpdate: insertOrUpdat, InPositionMap: true}

			metadata, lock := qc.GetMetadata()
			lock.Lock()

			setKeyMetadata(event.PositionData, metadata)

			_, quota := qc.GetQuota()
			// 设置T0/T1的持仓
			isT0 := strings.Contains(positionKey, "T0")
			if isT0 {
				event.PositionData["quotaT0"] = quota
			} else {
				event.PositionData["quotaT1"] = quota
			}

			lock.Unlock()

			configEventKey(event)
			
			events = append(events, event)

			postionMetaJs, _ := json.Marshal(event.PositionData)
			log.Printf("postionMetaJs_3: %s\n", postionMetaJs)
		}
	}

	// 如果是终态，更新T0/T1额度值（卖单会有额度释放，因此需要更新T0/T1持仓）
	if a.IsOrderFinished(tradeResp.CurrentTradeActionResp) {
		event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InsertOrUpdate: insertOrUpdat, InPositionMap: true}

		metadata, lock := qc.GetMetadata()

		lock.RLock()

		setKeyMetadata(event.PositionData, metadata)

		_, quota := qc.GetQuota()
		if isT0 {
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["quotaT1"] = quota
		}

		lock.RUnlock()
		configEventKey(event)
		events = append(events, event)

		evtJs, _ := json.Marshal(event)
		log.Printf("order finish, add event:%s\n", evtJs)

		// 多头开仓，买单，持仓将增加。因前面的AfterUpdatePositionByTradeResp可能会执行两次，不判断keyIndex将引入重复解冻的风险
		if runInTradeEngine && !isOrderPositionDecrease && keyIndex == 0 {
			CapitalController.RollbackAndReFreezeCapitalForTradeEnd(order)
		}
	}

	// 非交易类型的不考虑
	tradeActionResp := tradeResp.CurrentTradeActionResp
	if !statusficc.TradedStatus[tradeActionResp.OrdStatus] && !order_position.FillExecType[tradeActionResp.ExecType] {
		return
	}

	event := &order_position.PositionChangeEvent{PositionData: map[string]interface{}{}, InsertOrUpdate: insertOrUpdat, InPositionMap: true}
	cleanPrice, dirtyPrice := getPriceByOrder(order)
	if isOrderPositionDecrease {
		// 多头平仓：
		// 1、总持仓面额要减成交面额
		// 2、总名义本金要减成交面额/债券面值*订单全价（当订单全价为空的时候，例如FIX交易，订单全价 = 订单净价*1.05）
		metadata, lock := qc.GetMetadata()
		lock.Lock()

		setKeyMetadata(event.PositionData, metadata)

		totalCashQty, _ := metadata["totalCashQty"].(float64)
		totalCashQty -= float64(tradeActionResp.LastShares)
		metadata["totalCashQty"] = totalCashQty

		/* 从基于名义本金，改成基于全价、净价的总资本成本
		totalNotional, _ := metadata["totalNotional"].(float64)
		parValue, _ := metadata["parValue"].(float64)
		totalNotional -= float64(tradeActionResp.LastShares) / parValue * dirtyPrice
		metadata["totalNotional"] = totalNotional
		*/
		parValue, _ := metadata["parValue"].(float64)
		tmpVal := float64(tradeActionResp.LastShares) / parValue
		// 计算总全价资金成本
		averageCost, _, _ := attrutil.GetAttrValue(metadata, "averageCost", enum.AttrValueType_FLOAT)
		totalDirtyCost, _ := metadata["totalDirtyCost"].(float64)
		totalDirtyCost -= tmpVal * averageCost.(float64)
		metadata["totalDirtyCost"] = totalDirtyCost
		// 计算总净价资金成本
		averageCleanCost, _, _ := attrutil.GetAttrValue(metadata, "averageCleanCost", enum.AttrValueType_FLOAT)
		totalCleanCost, _ := metadata["totalCleanCost"].(float64)
		totalCleanCost -= tmpVal * averageCleanCost.(float64)
		metadata["totalCleanCost"] = totalCleanCost

		// 计算持仓均价
		calculateAverageCost(metadata)
		averageCost, _, _ = attrutil.GetAttrValue(metadata, "averageCost", enum.AttrValueType_FLOAT)
		event.PositionData["totalCashQty"] = metadata["totalCashQty"]
		//event.PositionData["totalNotional"] = metadata["totalNotional"]
		event.PositionData["totalDirtyCost"] = metadata["totalDirtyCost"]
		event.PositionData["totalCleanCost"] = metadata["totalCleanCost"]
		event.PositionData["averageCost"] = averageCost
		event.PositionData["averageCleanCost"] = metadata["averageCleanCost"]
		_, quota := qc.GetQuota()
		// 设置T0/T1的持仓
		isT0 := strings.Contains(positionKey, "T0")
		if isT0 {
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["quotaT1"] = quota
		}

		lock.Unlock()

		// 多头平仓，返回挂帐金额
		if runInTradeEngine {
			CapitalController.ReturnCapital(order, tradeResp, averageCost.(float64), parValue)
		}

	} else {
		// 多头开仓
		// 1、总持仓面额要加成交面额
		// 2、总名义本金要加成交面额/债券面值*订单全价（当订单全价为空的时候，例如FIX交易，订单全价 = 订单净价*1.05）
		metadata, lock := qc.GetMetadata()
		lock.Lock()

		setKeyMetadata(event.PositionData, metadata)

		totalCashQty, _ := metadata["totalCashQty"].(float64)
		totalCashQty += float64(tradeActionResp.LastShares)
		metadata["totalCashQty"] = totalCashQty

		/* 从基于名义本金，改成基于全价、净价的总资本成本
		totalNotional, _ := metadata["totalNotional"].(float64)
		parValue, _ := metadata["parValue"].(float64)
		totalNotional += float64(tradeActionResp.LastShares) / parValue * dirtyPrice
		metadata["totalNotional"] = totalNotional
		*/
		parValue, _ := metadata["parValue"].(float64)
		tmpVal := float64(tradeActionResp.LastShares) / parValue
		// 计算总全价资金成本
		totalDirtyCost, _ := metadata["totalDirtyCost"].(float64)
		totalDirtyCost += tmpVal * dirtyPrice
		metadata["totalDirtyCost"] = totalDirtyCost
		// 计算总净价资金成本
		totalCleanCost, _ := metadata["totalCleanCost"].(float64)
		totalCleanCost += tmpVal * cleanPrice
		metadata["totalCleanCost"] = totalCleanCost

		// 计算持仓均价
		calculateAverageCost(metadata)

		event.PositionData["totalCashQty"] = metadata["totalCashQty"]
		//event.PositionData["totalNotional"] = metadata["totalNotional"]
		event.PositionData["totalDirtyCost"] = metadata["totalDirtyCost"]
		event.PositionData["totalCleanCost"] = metadata["totalCleanCost"]
		event.PositionData["averageCost"] = metadata["averageCost"]
		event.PositionData["averageCleanCost"] = metadata["averageCleanCost"]
		_, quota := qc.GetQuota()
		// 设置T0/T1的持仓
		isT0 := strings.Contains(positionKey, "T0")
		if isT0 {
			event.PositionData["quotaT0"] = quota
		} else {
			event.PositionData["quotaT1"] = quota
		}

		lock.Unlock()
	}
	configEventKey(event)
	events = append(events, event)

	return events
}

func setKeyMetadata(keyMetadata, metadata map[string]interface{}) {
	keyMetadata["counterpartyID"] = metadata["counterpartyID"]
	keyMetadata["symbol"] = metadata["symbol"]
	keyMetadata["planCode"] = metadata["planCode"]
	keyMetadata["longShort"] = metadata["longShort"]
}

func (a *TitansFiccOrderPositionAdapter) ProcessDataSyncEvent(evt *datamap.DataChangeEvent, mapLock *sync.RWMutex, positionMap map[string]*quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp]) (events []*order_position.PositionChangeEvent) {

	log.Printf("received DataSyncEvent, TableAlias:%s, AddKeys:%d, ChgKeys:%d, DelKeys:%d\n", evt.TableAlias, len(evt.AddKeys), len(evt.ChgKeys), len(evt.DelKeys))

	if evt.TableAlias != "PositionBase" {
		return
	}

	if len(evt.AddKeys)+len(evt.ChgKeys) <= 0 {
		return
	}

	// 15353-GFZQ-FBZG-MARGINCB-0001-2471187.IB-LONG-T0
	autoSyncRepo := a.applicationCfg.GetAutoSyncRepo()
	// 新增
	for _, positionKey := range evt.AddKeys {
		log.Printf("===>for positionKey:%s\n", positionKey)
		valList, _, _ := autoSyncRepo.Get("PositionBase", positionKey)
		if len(valList) == 0 {
			continue
		}
		positionBaseRecord := valList[len(valList)-1]
		mapLock.RLock()
		qc, ok := positionMap[positionKey]
		mapLock.RUnlock()
		if ok {
			// 如果在positionMap存在
			metadata, changed := a.updateQuotaIfNecessary(positionKey, qc, positionBaseRecord)
			if changed {
				event := &order_position.PositionChangeEvent{InsertOrUpdate: 1, PositionData: metadata, InPositionMap: true}
				configEventKey(event)
				events = append(events, event)
			}
		} else {
			// 如果在positionMap不存在
			metadata := a.GetQuotaMetadata(positionBaseRecord)
			quota := getQuota(positionBaseRecord)
			if positionBaseRecord["TradeDateFlag"] == "T0" {
				metadata["baseQuotaT0"] = quota
				metadata["quotaT0"] = quota
				//metadata["baseQuotaT1"] = 0.0
				//metadata["quotaT1"] = 0.0
			} else {
				//metadata["baseQuotaT0"] = 0
				//metadata["quotaT0"] = 0
				metadata["baseQuotaT1"] = quota
				metadata["quotaT1"] = quota
			}
			event := &order_position.PositionChangeEvent{InsertOrUpdate: 0, PositionData: metadata, InPositionMap: false}
			configEventKey(event)
			events = append(events, event)
		}
	}
	// 更新
	for _, positionKey := range evt.ChgKeys {
		valList, _, _ := autoSyncRepo.Get("PositionBase", positionKey)
		if len(valList) == 0 {
			continue
		}
		positionBaseRecord := valList[len(valList)-1]
		mapLock.RLock()
		qc, ok := positionMap[positionKey]
		mapLock.RUnlock()
		if ok {
			// 如果在positionMap存在
			metadata, changed := a.updateQuotaIfNecessary(positionKey, qc, positionBaseRecord)
			if changed {
				event := &order_position.PositionChangeEvent{InsertOrUpdate: 1, PositionData: metadata, InPositionMap: true}
				configEventKey(event)
				events = append(events, event)
			}
		} else {
			// 如果在positionMap不存在
			metadata := a.GetQuotaMetadata(positionBaseRecord)
			quota := getQuota(positionBaseRecord)
			if positionBaseRecord["TradeDateFlag"] == "T0" {
				metadata["baseQuotaT0"] = quota
				metadata["quotaT0"] = quota
				//metadata["baseQuotaT1"] = 0.0
				//metadata["quotaT1"] = 0.0
			} else {
				//metadata["baseQuotaT0"] = 0
				//metadata["quotaT0"] = 0
				metadata["baseQuotaT1"] = quota
				metadata["quotaT1"] = quota
			}
			event := &order_position.PositionChangeEvent{InsertOrUpdate: 1, PositionData: metadata, InPositionMap: false}
			configEventKey(event)
			events = append(events, event)
		}
	}

	return
}

func getQuota(positionBaseRecord map[string]interface{}) float64 {
	quantity, _, _ := attrutil.GetAttrValue(positionBaseRecord, "Quantity", enum.AttrValueType_FLOAT)
	parValue, _, _ := attrutil.GetAttrValue(positionBaseRecord, "ParValue", enum.AttrValueType_FLOAT)
	return quantity.(float64) * math.Max(parValue.(float64), 100.0)
}

func (a *TitansFiccOrderPositionAdapter) updateQuotaIfNecessary(positionKey string, qc *quota.QuotaControl[*schema.TradeOrder, *schema.TradeActionResp], positionBaseRecord map[string]interface{}) (metadata map[string]interface{}, changed bool) {

	parValue, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "ParValue", enum.AttrValueType_FLOAT)
	if !ok {
		jsData, _ := json.Marshal(positionBaseRecord)
		domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("ParValue not fond in %s", jsData), fmt.Sprintf("ParValue not fond in %s\n", jsData))
	}
	if parValue.(float64) <= 0 {
		parValue = 100.0
	}
	baseQuota, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "Quantity", enum.AttrValueType_FLOAT)
	if !ok {
		jsData, _ := json.Marshal(positionBaseRecord)
		domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("error: Quantity not fond in %s", jsData), fmt.Sprintf("Quantity not fond in %s\n", jsData))
	}

	baseQuota = parValue.(float64) * baseQuota.(float64)

	metadata = make(map[string]interface{})
	metadata["counterpartyID"] = positionBaseRecord["CounterpartyID"]
	metadata["symbol"] = positionBaseRecord["Symbol"]
	metadata["planCode"] = positionBaseRecord["PlanCode"]
	metadata["longShort"] = positionBaseRecord["LongShort"]

	newQuota, quotaChanged := qc.UpdateQuota(baseQuota.(float64))
	// 更新T0/T1仓位
	if quotaChanged {
		isT0 := strings.Contains(positionKey, "T0")
		if isT0 {
			metadata["baseQuotaT0"] = baseQuota
			metadata["quotaT0"] = newQuota
		} else {
			metadata["baseQuotaT1"] = baseQuota
			metadata["quotaT1"] = newQuota
		}
	}

	// 更新总持仓
	totalCashQtyChanged := false
	baseCashQty, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "BaseCashQty", enum.AttrValueType_FLOAT)
	if ok {
		qcMetadata, lock := qc.GetMetadata()
		lock.Lock()
		currBaseCashQty, _, _ := attrutil.GetAttrValue(qcMetadata, "baseCashQty", enum.AttrValueType_FLOAT)
		diff := baseCashQty.(float64) - currBaseCashQty.(float64)
		if diff != 0 {
			qcMetadata["baseCashQty"] = baseCashQty
			totalCashQty, _, _ := attrutil.GetAttrValue(qcMetadata, "totalCashQty", enum.AttrValueType_FLOAT)
			qcMetadata["totalCashQty"] = totalCashQty.(float64) + diff
			totalCashQtyChanged = true
			metadata["baseCashQty"] = qcMetadata["baseCashQty"]
			metadata["totalCashQty"] = qcMetadata["totalCashQty"]
		}
		lock.Unlock()
	}

	// 更新总名本，改成资金总成本
	/*
		totalNotionalChanged := false
		baseNotional, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "BaseNotional", enum.AttrValueType_FLOAT)
		if ok {
			qcMetadata, lock := qc.GetMetadata()
			lock.Lock()
			currBaseNotional, _, _ := attrutil.GetAttrValue(qcMetadata, "baseNotional", enum.AttrValueType_FLOAT)
			diff := baseNotional.(float64) - currBaseNotional.(float64)
			if diff != 0 {
				qcMetadata["baseNotional"] = baseNotional
				totalNotional, _, _ := attrutil.GetAttrValue(qcMetadata, "totalNotional", enum.AttrValueType_FLOAT)
				qcMetadata["totalNotional"] = totalNotional.(float64) + diff
				totalNotionalChanged = true
				metadata["baseNotional"] = qcMetadata["baseNotional"]
				metadata["totalNotional"] = qcMetadata["totalNotional"]
			}
			lock.Unlock()
		}
	*/
	// 更新全价总资金成本
	totalDirtyCostChanged := false
	baseDirtyCost, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "BaseDirtyCost", enum.AttrValueType_FLOAT)
	if ok {
		qcMetadata, lock := qc.GetMetadata()
		lock.Lock()
		currBaseDirtyCost, _, _ := attrutil.GetAttrValue(qcMetadata, "baseDirtyCost", enum.AttrValueType_FLOAT)
		diff := baseDirtyCost.(float64) - currBaseDirtyCost.(float64)
		if diff != 0 {
			qcMetadata["baseDirtyCost"] = baseDirtyCost
			totalDirtyCost, _, _ := attrutil.GetAttrValue(qcMetadata, "totalDirtyCost", enum.AttrValueType_FLOAT)
			qcMetadata["totalDirtyCost"] = totalDirtyCost.(float64) + diff
			totalDirtyCostChanged = true
			metadata["baseDirtyCost"] = qcMetadata["baseDirtyCost"]
			metadata["totalDirtyCost"] = qcMetadata["totalDirtyCost"]
		}
		lock.Unlock()
	}
	// 更新净价总资金成本
	totalCleanCostChanged := false
	baseCleanCost, ok, _ := attrutil.GetAttrValue(positionBaseRecord, "BaseCleanCost", enum.AttrValueType_FLOAT)
	if ok {
		qcMetadata, lock := qc.GetMetadata()
		lock.Lock()
		currBaseCleanCost, _, _ := attrutil.GetAttrValue(qcMetadata, "baseCleanCost", enum.AttrValueType_FLOAT)
		diff := baseCleanCost.(float64) - currBaseCleanCost.(float64)
		if diff != 0 {
			qcMetadata["baseCleanCost"] = baseCleanCost
			totalCleanCost, _, _ := attrutil.GetAttrValue(qcMetadata, "totalCleanCost", enum.AttrValueType_FLOAT)
			qcMetadata["totalCleanCost"] = totalCleanCost.(float64) + diff
			totalCleanCostChanged = true
			metadata["baseCleanCost"] = qcMetadata["baseCleanCost"]
			metadata["totalCleanCost"] = qcMetadata["totalCleanCost"]
		}
		lock.Unlock()
	}

	// 更新持仓均价
	if totalCashQtyChanged || totalDirtyCostChanged || totalCleanCostChanged {
		metadata["parValue"] = parValue
		calculateAverageCost(metadata)
	}

	changed = quotaChanged || totalCashQtyChanged || totalDirtyCostChanged || totalCleanCostChanged

	return metadata, changed
}

func configEventKey(event *order_position.PositionChangeEvent) {
	event.PositionData["key"] = fmt.Sprintf("%v-%v-%v-%v", event.PositionData["counterpartyID"], event.PositionData["symbol"], event.PositionData["planCode"], event.PositionData["longShort"])
}
