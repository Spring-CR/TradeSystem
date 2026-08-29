package ficc_v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/logger"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	"rhino-data/datamap"
	"strconv"
)

var (
	//json         = jsoniter.ConfigCompatibleWithStandardLibrary
	sideSell     = string(enum.Side_Sell)
	sideBuy      = string(enum.Side_Buy)
	fillExecType = map[string]bool{
		"F": true,
		"1": true,
		"2": true,
	}
)

const (
	settleTypeT0            = "T+0"
	settleTypeT1            = "T+1"
	positionUnitNameLongT0  = "L_T0"
	positionUnitNameLongT1  = "L_T1"
	positionUnitNameShortT0 = "S_T0"
	positionUnitNameShortT1 = "S_T1"
)

type TitansFiccOrderPositionAdapter struct {
	applicationCfg        *domain_cfg.ApplicationCfg
	orderLog              *logger.OrderLog
	marginExposureManager *MarginExposureManager
	strictCtpyMap         map[string]bool
}

func NewTitansFiccOrderPositionAdapter(applicationCfg *domain_cfg.ApplicationCfg, orderLog *logger.OrderLog) (adapter *TitansFiccOrderPositionAdapter, de *domain_error.Error) {
	log.Printf("construct TitansFiccOrderPositiondapter v2...")

	var initMarginRatio float64
	// 从配置文件提取初保比率
	configMap := applicationCfg.GetApplicationCfgItemMap()
	configItem, ok := configMap["InitMarginRatio"]
	if ok && configItem.ConfigItemValue != "" {
		val, err := strconv.ParseFloat(configItem.ConfigItemValue, 64)
		if err == nil && val > 0 {
			initMarginRatio = val / 100.0
			log.Printf("Reset InitMarginRatio=%f\n", initMarginRatio)
		}
	}

	marginExposureManager := NewMarginExposureManager(initMarginRatio, initMarginRatio)

	adapter = &TitansFiccOrderPositionAdapter{applicationCfg: applicationCfg, orderLog: orderLog, marginExposureManager: marginExposureManager}
	adapter.initStrictCtpyMap()

	return
}

func (a *TitansFiccOrderPositionAdapter) HasSufficientQuota(positionUnits []*order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (sufficient bool, de *domain_error.Error) {

	side := tradeOrder.Side
	if side != sideSell { // 买入无需验券，直接返回
		sufficient = true
		return
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		sufficient = true
		return
	}

	if metadata == nil {
		sufficient = true
		return
	}

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	val := tradeOrder.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)

	switch settleType {
	case settleTypeT0:
		longAvailablePositionT0 := positionRecord.LongAvailablePositionT0
		acquirePositionT0 := tradeOrder.OrderQty
		if acquirePositionT0 > longAvailablePositionT0 {
			de = domain_error.Build(domain_error.QTY_NOT_ENOUGH_ERR_CODE, nil, "T+0多头可用面额")
			de.Msg = de.Msg + ": 当前标的可用持仓不足，若要开仓空头交易，请联系业务人员申请借券"
			return
		}
	case settleTypeT1:
		longAvailablePositionT1 := positionRecord.LongAvailablePositionT1
		acquirePositionT1 := tradeOrder.OrderQty
		if acquirePositionT1 > longAvailablePositionT1 { // 出现T+1超卖
			symbol := tradeOrder.Symbol
			allowOverSold, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "allowOverSold", enum.AttrValueType_BOOL)
			if !allowOverSold.(bool) {
				de = domain_error.Build(domain_error.SHORT_SELL_SIMBOL_ERR_CODE, nil, symbol)
				return
			}
		}
	}

	sufficient = true
	return
}

func (a *TitansFiccOrderPositionAdapter) CalculateOrderFreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (freezeQuota float64, freezeQuotaBuf float64, ok bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}

	val := tradeOrder.ExtendAttrMap["settlType"]
	settleType, _ := val.(string)

	a.orderLog.Printf(tradeOrder, nil, "[CalculateOrderFreezeQuotaInPositionUnit=%v] SettlType=%v, OrderQty=%v, LongAvailablePositionT0=%v, LongAvailablePositionT1=%v, ShortAvailablePositionT0=%v, ShortAvailablePositionT1=%v", positionUnit.GetName(), settleType, tradeOrder.OrderQty, positionRecord.LongAvailablePositionT0, positionRecord.LongAvailablePositionT1, positionRecord.ShortAvailablePositionT0, positionRecord.ShortAvailablePositionT1)

	// 按持仓类型，计算需冻结持仓
	switch positionUnit.GetName() {

	// 计算T+0多头可用持仓
	case positionUnitNameLongT0:

		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.LongAvailablePositionT0 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.LongAvailablePositionT0
				ok = true
				return
			}
		case settleTypeT1:
			if positionRecord.LongAvailablePositionT0 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.LongAvailablePositionT0
				ok = true
				return
			}
		}

	// 计算T+1多头可用持仓
	case positionUnitNameLongT1:

		if tradeOrder.Side != sideSell {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.LongAvailablePositionT1 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.LongAvailablePositionT1
				ok = true
				return
			}
		case settleTypeT1:
			if positionRecord.LongAvailablePositionT1 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.LongAvailablePositionT1
				ok = true
				return
			}
		}

	// 计算T+0空头可用持仓
	case positionUnitNameShortT0:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.ShortAvailablePositionT0 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.ShortAvailablePositionT0
				ok = true
				return
			}
		case settleTypeT1:
			freezeQuotaBuf = math.Max(0, positionRecord.ShortAvailablePositionT1-positionRecord.ShortAvailablePositionT0)
			if freezeQuotaBuf >= tradeOrder.OrderQty {
				freezeQuota = 0
				ok = true
				return
			} else {
				freezeQuota = tradeOrder.OrderQty - freezeQuotaBuf
				if freezeQuota >= positionRecord.ShortAvailablePositionT0 {
					freezeQuota = positionRecord.ShortAvailablePositionT0
				}
			}
		}

	// 计算T+1空头可用持仓
	case positionUnitNameShortT1:

		if tradeOrder.Side != sideBuy {
			return
		}

		switch settleType {
		case settleTypeT0:
			if positionRecord.ShortAvailablePositionT1 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.ShortAvailablePositionT1
				ok = true
				return
			}
		case settleTypeT1:
			if positionRecord.ShortAvailablePositionT1 > tradeOrder.OrderQty {
				freezeQuota = tradeOrder.OrderQty
				ok = true
				return
			} else {
				freezeQuota = positionRecord.ShortAvailablePositionT1
				ok = true
				return
			}
		}
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) AfterFreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, freezeQty float64, tradeOrder *schema.TradeOrder, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker, lastPositionUnit bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLongT0:

		positionRecord.LongAvailablePositionT0 -= freezeQty

	case positionUnitNameLongT1:

		positionRecord.LongAvailablePositionT1 -= freezeQty

	case positionUnitNameShortT0:

		positionRecord.ShortAvailablePositionT0 -= freezeQty

	case positionUnitNameShortT1:

		positionRecord.ShortAvailablePositionT1 -= freezeQty
	}

	if lastPositionUnit {
		a.marginExposureManager.AddOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

		//a.orderLog.Printf(tradeOrder, nil, "[AfterFreezeQuotaInPositionUnit] LongAvailablePositionT0=%v, LongAvailablePositionT1=%v, ShortAvailablePositionT0=%v, ShortAvailablePositionT1=%v, MaxLongMarginOccupancy=%v, MaxShortMarginOccupancy=%v, MaxMarginOccupancy=%v", positionRecord.LongAvailablePositionT0, positionRecord.LongAvailablePositionT1, positionRecord.ShortAvailablePositionT0, positionRecord.ShortAvailablePositionT1, positionRecord.MaxLongMarginOccupancy, positionRecord.MaxShortMarginOccupancy, positionRecord.MaxMarginOccupancy)
		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterFreezeQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

func (a *TitansFiccOrderPositionAdapter) GetPositionCalculatorKey(tradeOrder *schema.TradeOrder) (key string) {
	key = tradeOrder.Account + "_" + tradeOrder.Symbol
	return
}

func (a *TitansFiccOrderPositionAdapter) GetPositionCalculatorConstructParams(tradeOrder *schema.TradeOrder) (param *order_position_manager.PositionCalculatorConstructParam) {
	var metadata interface{}
	var positionUnits []*order_position_manager.PositionUnit
	positionRecord := a.loadOrConstructPositionRecord(tradeOrder)
	metadata = positionRecord
	positionUnits = []*order_position_manager.PositionUnit{
		order_position_manager.NewPositionUnit(positionUnitNameLongT0, positionRecord.LongAvailablePositionT0, a.orderLog),
		order_position_manager.NewPositionUnit(positionUnitNameLongT1, positionRecord.LongAvailablePositionT1, a.orderLog),
		order_position_manager.NewPositionUnit(positionUnitNameShortT0, positionRecord.ShortAvailablePositionT0, a.orderLog),
		order_position_manager.NewPositionUnit(positionUnitNameShortT1, positionRecord.ShortAvailablePositionT1, a.orderLog),
	}
	param = &order_position_manager.PositionCalculatorConstructParam{
		Key:           a.GetPositionCalculatorKey(tradeOrder),
		Metadata:      metadata,
		PositionUnits: positionUnits,
	}
	return
}

func (a *TitansFiccOrderPositionAdapter) LoadInitPositionRecords(applicationConfig *domain_cfg.ApplicationCfg) (params []*order_position_manager.PositionCalculatorConstructParam) {

	m, _ := applicationConfig.GetAutoSyncRepo().GetMapData("PositionBase")
	if len(m) == 0 {
		return
	}

	visiMap := make(map[string]bool)

	for _, valList := range m {
		n := len(valList)
		if n == 0 {
			continue
		}
		record := valList[n-1]

		account, _, _ := attrutil.GetAttrValue(record, "CounterpartyID", enum.AttrValueType_INT)
		counterparty, _, _ := attrutil.GetAttrValue(record, "Counterparty", enum.AttrValueType_STRING)
		symbol, _, _ := attrutil.GetAttrValue(record, "Symbol", enum.AttrValueType_STRING)
		symbolName, _, _ := attrutil.GetAttrValue(record, "SecurityName", enum.AttrValueType_STRING)
		currency, _, _ := attrutil.GetAttrValue(record, "Currency", enum.AttrValueType_STRING)
		planCode, _, _ := attrutil.GetAttrValue(record, "PlanCode", enum.AttrValueType_STRING)
		ultraContractCode, _, _ := attrutil.GetAttrValue(record, "UltraContractCode", enum.AttrValueType_STRING)
		securityExchange, _, _ := attrutil.GetAttrValue(record, "SecurityExchange", enum.AttrValueType_STRING)
		parValue, _, _ := attrutil.GetAttrValue(record, "ParValue", enum.AttrValueType_FLOAT)

		key := fmt.Sprintf("%v_%v", account, symbol)
		if visiMap[key] {
			continue
		}

		positionRecord := &PositionRecord{
			Account:           account.(int),
			CounterpartyID:    account.(int),
			Counterparty:      counterparty.(string),
			Symbol:            symbol.(string),
			SymbolName:        symbolName.(string),
			Currency:          currency.(string),
			PlanCode:          planCode.(string),
			UltraContractCode: ultraContractCode.(string),
			SecurityExchange:  securityExchange.(string),
			SecurityType:      "BOND",
			ParValue:          parValue.(float64),
		}

		a.refinePositionRecord(positionRecord)

		a.marginExposureManager.CalculateInitMarginExposure(positionRecord)

		positionUnits := []*order_position_manager.PositionUnit{
			order_position_manager.NewPositionUnit(positionUnitNameLongT0, positionRecord.LongAvailablePositionT0, a.orderLog),
			order_position_manager.NewPositionUnit(positionUnitNameLongT1, positionRecord.LongAvailablePositionT1, a.orderLog),
			order_position_manager.NewPositionUnit(positionUnitNameShortT0, positionRecord.ShortAvailablePositionT0, a.orderLog),
			order_position_manager.NewPositionUnit(positionUnitNameShortT1, positionRecord.ShortAvailablePositionT1, a.orderLog),
		}
		param := &order_position_manager.PositionCalculatorConstructParam{
			Key:           key,
			Metadata:      positionRecord,
			PositionUnits: positionUnits,
		}

		params = append(params, param)

		visiMap[key] = true
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) GeneralizePositionRecord(metadata interface{}, insert bool) (positionData map[string]interface{}) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		return
	}

	if insert {
		positionData = map[string]interface{}{
			"key":                        fmt.Sprintf("%v_%v", positionRecord.Account, positionRecord.Symbol),
			"account":                    positionRecord.Account,
			"counterpartyID":             positionRecord.CounterpartyID,
			"counterparty":               positionRecord.Counterparty,
			"symbol":                     positionRecord.Symbol,
			"symbolName":                 positionRecord.SymbolName,
			"currency":                   positionRecord.Currency,
			"planCode":                   positionRecord.PlanCode,
			"ultraContractCode":          positionRecord.UltraContractCode,
			"securityExchange":           positionRecord.SecurityExchange,
			"securityType":               positionRecord.SecurityType,
			"parValue":                   positionRecord.ParValue,
			"initNetPositionT0":          positionRecord.InitNetPositionT0,
			"initNetPositionT1":          positionRecord.InitNetPositionT1,
			"netPositionT0":              positionRecord.NetPositionT0,
			"netPositionT1":              positionRecord.NetPositionT1,
			"longAvailablePositionT0":    positionRecord.LongAvailablePositionT0,
			"longAvailablePositionT1":    positionRecord.LongAvailablePositionT1,
			"shortAvailablePositionT0":   positionRecord.ShortAvailablePositionT0,
			"shortAvailablePositionT1":   positionRecord.ShortAvailablePositionT1,
			"longCleanPriceCost":         positionRecord.LongCleanPriceCost,
			"longDirtyPriceCost":         positionRecord.LongDirtyPriceCost,
			"longDirtyPriceWithFeeCost":  positionRecord.LongDirtyPriceWithFeeCost,
			"shortCleanPriceCost":        positionRecord.ShortCleanPriceCost,
			"shortDirtyPriceCost":        positionRecord.ShortDirtyPriceCost,
			"shortDirtyPriceWithFeeCost": positionRecord.ShortDirtyPriceWithFeeCost,
			"longAvgCleanPrice":          positionRecord.LongAvgCleanPrice,
			"longAvgDirtyPrice":          positionRecord.LongAvgDirtyPrice,
			"longAvgDirtyPriceWithFee":   positionRecord.LongAvgDirtyPriceWithFee,
			"shortAvgCleanPrice":         positionRecord.ShortAvgCleanPrice,
			"shortAvgDirtyPrice":         positionRecord.ShortAvgDirtyPrice,
			"shortAvgDirtyPriceWithFee":  positionRecord.ShortAvgDirtyPriceWithFee,
			"maxLongMarginOccupancy":     positionRecord.MaxLongMarginOccupancy,
			"maxShortMarginOccupancy":    positionRecord.MaxShortMarginOccupancy,
			"maxMarginOccupancy":         positionRecord.MaxMarginOccupancy,
		}
	} else {
		positionData = map[string]interface{}{
			"key":                        fmt.Sprintf("%v_%v", positionRecord.Account, positionRecord.Symbol),
			"netPositionT0":              positionRecord.NetPositionT0,
			"netPositionT1":              positionRecord.NetPositionT1,
			"longAvailablePositionT0":    positionRecord.LongAvailablePositionT0,
			"longAvailablePositionT1":    positionRecord.LongAvailablePositionT1,
			"shortAvailablePositionT0":   positionRecord.ShortAvailablePositionT0,
			"shortAvailablePositionT1":   positionRecord.ShortAvailablePositionT1,
			"longCleanPriceCost":         positionRecord.LongCleanPriceCost,
			"longDirtyPriceCost":         positionRecord.LongDirtyPriceCost,
			"longDirtyPriceWithFeeCost":  positionRecord.LongDirtyPriceWithFeeCost,
			"shortCleanPriceCost":        positionRecord.ShortCleanPriceCost,
			"shortDirtyPriceCost":        positionRecord.ShortDirtyPriceCost,
			"shortDirtyPriceWithFeeCost": positionRecord.ShortDirtyPriceWithFeeCost,
			"longAvgCleanPrice":          positionRecord.LongAvgCleanPrice,
			"longAvgDirtyPrice":          positionRecord.LongAvgDirtyPrice,
			"longAvgDirtyPriceWithFee":   positionRecord.LongAvgDirtyPriceWithFee,
			"shortAvgCleanPrice":         positionRecord.ShortAvgCleanPrice,
			"shortAvgDirtyPrice":         positionRecord.ShortAvgDirtyPrice,
			"shortAvgDirtyPriceWithFee":  positionRecord.ShortAvgDirtyPriceWithFee,
			"maxLongMarginOccupancy":     positionRecord.MaxLongMarginOccupancy,
			"maxShortMarginOccupancy":    positionRecord.MaxShortMarginOccupancy,
			"maxMarginOccupancy":         positionRecord.MaxMarginOccupancy,
		}
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) GetQuotaNotEnoughHandler() func(tradeOrder *schema.TradeOrder, de *domain_error.Error) {
	return nil
}

func (a *TitansFiccOrderPositionAdapter) PreparePositionAdjustmentParams(tradeOrder *schema.TradeOrder) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error) {
	return
}

func (a *TitansFiccOrderPositionAdapter) GetPositionCalculatorKeysForPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, mockTradeOrders map[string][]*schema.TradeOrder, purgingLog *schema.DataPurgingLog) (positionCalculatorKeysToPurge []string) {
	return
}

func (a *TitansFiccOrderPositionAdapter) ReloadPositionRecordsForPurgingTask(applicationConfig *domain_cfg.ApplicationCfg, purgingLog *schema.DataPurgingLog) (params []*order_position_manager.PositionCalculatorConstructParam) {
	return
}

func (a *TitansFiccOrderPositionAdapter) ProcessPositionBaseDataSyncEvent(event *datamap.DataChangeEvent) {

}

func (a *TitansFiccOrderPositionAdapter) ParsePositionRecordFromReposRecord(extendAttrMap map[string]interface{}) (key string, positionRecord interface{}) {
	return
}

func (a *TitansFiccOrderPositionAdapter) GetPositionCalculatorConstructParamForPositionRecord(positionRecord interface{}) (param *order_position_manager.PositionCalculatorConstructParam) {
	return
}

func (a *TitansFiccOrderPositionAdapter) PreparePositionAdjustmentParamsByPositionBaseDiff(positionRecordBase, positionRecordCurr interface{}) (mockTradeOrder *schema.TradeOrder, mockTradeActionResp *schema.TradeActionResp, de *domain_error.Error) {
	return
}

func (a *TitansFiccOrderPositionAdapter) PrepareForRecover(positionRecord interface{}) {

}

func (a *TitansFiccOrderPositionAdapter) UpdatePositionBaseDynamically()(dynamicallyUpdate bool) {
	return
}