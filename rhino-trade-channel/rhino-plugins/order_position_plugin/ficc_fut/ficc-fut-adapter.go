package ficc_fut

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/logger"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	"rhino-data/datamap"
	"strconv"
	"strings"
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
	positionUnitNameLong  = "L"
	positionUnitNameShort = "S"
)

type FiccFutOrderPositionAdapter struct {
	applicationCfg        *domain_cfg.ApplicationCfg
	autoSyncRepo          *datamap.AutoSyncRepo
	orderLog              *logger.OrderLog
	marginExposureManager *MarginExposureManager
}

func NewFiccFutOrderPositionAdapter(applicationCfg *domain_cfg.ApplicationCfg, orderLog *logger.OrderLog) (adapter *FiccFutOrderPositionAdapter, de *domain_error.Error) {

	log.Printf("construct FiccFutOrderPositionAdapter...")
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

	// Todo, 未计算保证金费率
	marginExposureManager := NewMarginExposureManager(initMarginRatio, initMarginRatio, orderLog)

	adapter = &FiccFutOrderPositionAdapter{applicationCfg: applicationCfg, orderLog: orderLog, marginExposureManager: marginExposureManager, autoSyncRepo: applicationCfg.GetAutoSyncRepo()}

	return
}

func (a *FiccFutOrderPositionAdapter) HasSufficientQuota(positionUnits []*order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (sufficient bool, de *domain_error.Error) {

	sufficient = true
	var errType []string
	var errMsg []string

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	longQty, shortQty, buyOrderLeftQty, sellOrderLeftQty := a.marginExposureManager.GetQty(positionRecord)
	longCost, shortCost, buyOrderLeftCost, sellOrderLeftCost := a.marginExposureManager.GetCost(positionRecord)

	limitLongQty, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "limitLongQty", enum.AttrValueType_FLOAT)
	limitShortQty, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "limitShortQty", enum.AttrValueType_FLOAT)
	limitNotional, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "limitNotional", enum.AttrValueType_FLOAT)

	log.Printf("HasSufficientQuota, appOrdID:%s, side:%v, limitLongQty:%v, limitShortQty:%v, limitNotional:%v, longQty:%v, shortQty:%v, buyOrderLeftQty:%v, sellOrderLeftQty:%v, longCost:%v, shortCost:%v, buyOrderLeftCost:%v, sellOrderLeftCost:%v\n",
		tradeOrder.AppOrdID, tradeOrder.Side, limitLongQty, limitShortQty, limitNotional, longQty, shortQty, buyOrderLeftQty, sellOrderLeftQty, longCost, shortCost, buyOrderLeftCost, sellOrderLeftCost)

	if tradeOrder.Side == "1" && tradeOrder.OrderQty+longQty+buyOrderLeftQty > limitLongQty.(float64) {
		sufficient = false
		errType = append(errType, "1")
		errMsg = append(errMsg, "多头持仓手数超过期货品种最大交易手数限额")
	}

	if tradeOrder.Side == "2" && tradeOrder.OrderQty+shortQty+sellOrderLeftQty > limitShortQty.(float64) {
		sufficient = false
		errType = append(errType, "2")
		errMsg = append(errMsg, "空头持仓手数超过期货品种最大交易手数限额")
	}

	priceCNY, ok, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "priceCNYWithFee", enum.AttrValueType_FLOAT)
	if !ok {
		domain_error.ProcessSevereError(false, 0, nil, errors.New("fail to get priceCNA for order "+tradeOrder.AppOrdID), "fail to get priceCNA")
		priceCNY = tradeOrder.Price
	}

	contractMultiplier, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "contractMultiplier", enum.AttrValueType_FLOAT)
	if contractMultiplier == 0 {
		domain_error.ProcessSevereError(false, 0, nil, errors.New("fail to get contractMultiplier for order "+tradeOrder.AppOrdID), "fail to get contractMultiplier")
		contractMultiplier = 1.0
	}

	ordertNotional := priceCNY.(float64) * contractMultiplier.(float64) * tradeOrder.OrderQty

	if tradeOrder.Side == "1" && ordertNotional+longCost+buyOrderLeftCost > limitNotional.(float64) { // 买单
		sufficient = false
		errType = append(errType, "3")
		errMsg = append(errMsg, "多头名义本金超过产品最大交易名义本金限额")
	}

	if tradeOrder.Side == "2" && ordertNotional+shortCost+sellOrderLeftCost > limitNotional.(float64) { // 卖单
		sufficient = false
		errType = append(errType, "4")
		errMsg = append(errMsg, "空头名义本金超过产品最大交易名义本金限额")
	}

	if !sufficient {
		de = domain_error.Build(domain_error.QUOTA_LIMIT_EXCEEDED_ERR_CODE, errors.New(strings.Join(errType, ",")), strings.Join(errMsg, ","))
	}

	return
}

func (a *FiccFutOrderPositionAdapter) CalculateOrderFreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, tradeOrder *schema.TradeOrder, metadata interface{}) (freezeQuota float64, freezeQuotaBuf float64, ok bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	if len(tradeOrder.ExtendAttrMap) == 0 {
		return
	}

	a.orderLog.Printf(tradeOrder, nil, "[CalculateOrderFreezeQuotaInPositionUnit=%v] OrderQty=%v, LongAvailablePosition=%v, ShortAvailablePosition=%v", positionUnit.GetName(), tradeOrder.OrderQty, positionRecord.LongAvailablePosition, positionRecord.ShortAvailablePosition)

	// 按持仓类型，计算需冻结持仓
	switch positionUnit.GetName() {

	// 计算多头可用持仓
	case positionUnitNameLong:

		if tradeOrder.Side != sideSell {
			return
		}

		if positionRecord.LongAvailablePosition > tradeOrder.OrderQty {
			freezeQuota = tradeOrder.OrderQty
			ok = true
			return
		} else {
			freezeQuota = positionRecord.LongAvailablePosition
			ok = true
			return
		}

	// 计算空头可用持仓
	case positionUnitNameShort:

		if tradeOrder.Side != sideBuy {
			return
		}

		if positionRecord.ShortAvailablePosition > tradeOrder.OrderQty {
			freezeQuota = tradeOrder.OrderQty
			ok = true
			return
		} else {
			freezeQuota = positionRecord.ShortAvailablePosition
			ok = true
			return
		}
	}

	return
}

func (a *FiccFutOrderPositionAdapter) AfterFreezeQuotaInPositionUnit(positionUnit *order_position_manager.PositionUnit, freezeQty float64, tradeOrder *schema.TradeOrder, metadata interface{}, quotaLocker *order_position_manager.QuotaLocker, lastPositionUnit bool) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		errMsg := "metadata is not type of *PositionRecord"
		domain_error.ProcessSevereError(true, 5, nil, errors.New(errMsg), errMsg)
	}

	switch positionUnit.GetName() {

	case positionUnitNameLong:

		positionRecord.LongAvailablePosition -= freezeQty

	case positionUnitNameShort:

		positionRecord.ShortAvailablePosition -= freezeQty
	}

	if lastPositionUnit {
		log.Printf("AfterFreezeQuotaInPositionUnit, AddOrder:%s\n", tradeOrder.AppOrdID)
		a.marginExposureManager.AddOrder(tradeOrder)
		a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)
		js, _ := json.MarshalIndent(positionRecord, "", "  ")
		a.orderLog.Printf(tradeOrder, nil, "[AfterFreezeQuotaInPositionUnit] PositionRecord=%s", js)
	}
}

func (a *FiccFutOrderPositionAdapter) GetPositionCalculatorKey(tradeOrder *schema.TradeOrder) (key string) {
	key = fmt.Sprintf("%v-%v", tradeOrder.Account, tradeOrder.ExtendAttrMap["symbol2"])
	return
}

func (a *FiccFutOrderPositionAdapter) GetPositionCalculatorConstructParams(tradeOrder *schema.TradeOrder) (param *order_position_manager.PositionCalculatorConstructParam) {
	var metadata interface{}
	var positionUnits []*order_position_manager.PositionUnit
	positionRecord := a.loadOrConstructPositionRecord(tradeOrder)
	metadata = positionRecord
	positionUnits = []*order_position_manager.PositionUnit{
		order_position_manager.NewPositionUnit(positionUnitNameLong, positionRecord.LongAvailablePosition, a.orderLog),
		order_position_manager.NewPositionUnit(positionUnitNameShort, positionRecord.ShortAvailablePosition, a.orderLog),
	}
	param = &order_position_manager.PositionCalculatorConstructParam{
		Key:           a.GetPositionCalculatorKey(tradeOrder),
		Metadata:      metadata,
		PositionUnits: positionUnits,
	}
	return
}

func (a *FiccFutOrderPositionAdapter) LoadInitPositionRecords(applicationConfig *domain_cfg.ApplicationCfg) (params []*order_position_manager.PositionCalculatorConstructParam) {

	params1 := a.loadInitPositionRecords(applicationConfig, "PositionBaseDms")
	if len(params1) > 0 {
		params = append(params, params1...)
	}

	params2 := a.loadInitPositionRecords(applicationConfig, "PositionBaseOvs")
	if len(params2) > 0 {
		params = append(params, params2...)
	}

	return
}

func (a *FiccFutOrderPositionAdapter) loadInitPositionRecords(applicationConfig *domain_cfg.ApplicationCfg, positionBaseCollectionName string) (params []*order_position_manager.PositionCalculatorConstructParam) {

	m, _ := applicationConfig.GetAutoSyncRepo().GetMapData(positionBaseCollectionName)
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

		_, positionRecord := a.parsePositionRecordFromReposRecord(record)

		key := positionRecord.Key
		if visiMap[key] {
			continue
		}

		/*
		a.marginExposureManager.CalculateInitMarginExposure(positionRecord)

		positionUnits := []*order_position_manager.PositionUnit{
			order_position_manager.NewPositionUnit(positionUnitNameLong, positionRecord.LongAvailablePosition, a.orderLog),
			order_position_manager.NewPositionUnit(positionUnitNameShort, positionRecord.ShortAvailablePosition, a.orderLog),
		}
		param := &order_position_manager.PositionCalculatorConstructParam{
			Key:           key,
			Metadata:      positionRecord,
			PositionUnits: positionUnits,
		}*/

		param := a.GetPositionCalculatorConstructParamForPositionRecord(positionRecord)
		if param == nil {
			log.Printf("fail to GetPositionCalculatorConstructParamForPositionRecord,ke:%s\n", positionRecord)
			continue
		}

		params = append(params, param)

		visiMap[key] = true
	}

	return
}

func (a *FiccFutOrderPositionAdapter) GetPositionCalculatorConstructParamForPositionRecord(_positionRecord interface{}) (param *order_position_manager.PositionCalculatorConstructParam) {

	positionRecord, ok := _positionRecord.(*PositionRecord)

	if !ok {
		return
	}

	// 注意：只有首次新增的时候，这时才调用CalculateInitMarginExposure
	a.marginExposureManager.CalculateInitMarginExposure(positionRecord)

	positionUnits := []*order_position_manager.PositionUnit{
		order_position_manager.NewPositionUnit(positionUnitNameLong, positionRecord.LongAvailablePosition, a.orderLog),
		order_position_manager.NewPositionUnit(positionUnitNameShort, positionRecord.ShortAvailablePosition, a.orderLog),
	}
	param = &order_position_manager.PositionCalculatorConstructParam{
		Key:           positionRecord.Key,
		Metadata:      positionRecord,
		PositionUnits: positionUnits,
	}

	return
}

func (a *FiccFutOrderPositionAdapter) GeneralizePositionRecord(metadata interface{}, insert bool) (positionData map[string]interface{}) {

	positionRecord, ok := metadata.(*PositionRecord)
	if !ok {
		return
	}

	if insert {
		positionData = map[string]interface{}{
			// 标的属性
			"key":                positionRecord.Key,
			"account":            positionRecord.Account,
			"counterpartyID":     positionRecord.CounterpartyID,
			"counterparty":       positionRecord.Counterparty,
			"symbol2":            positionRecord.Symbol2,
			"symbolName":         positionRecord.SymbolName,
			"currency":           positionRecord.Currency,
			"planCode":           positionRecord.PlanCode,
			"ultraContractCode":  positionRecord.UltraContractCode,
			"securityExchange":   positionRecord.SecurityExchange,
			"securityType":       positionRecord.SecurityType,
			"productCode":        positionRecord.ProductCode,
			"contractMultiplier": positionRecord.ContractMultiplier,
			"exchangeRateCNY":    positionRecord.ExchangeRateCNY,
			"longMarginRatio":    positionRecord.LongMarginRatio,
			"shortMarginRatio":   positionRecord.ShortMarginRatio,
			"exchangeArea":       positionRecord.ExchangeArea,
			"contractBaseDate":   positionRecord.ContractBaseDate,
			// 统计属性
			"initNetPosition":           positionRecord.InitNetPosition,
			"initLongPriceCost":         positionRecord.InitLongPriceCost,
			"initLongPriceWithFeeCost":  positionRecord.InitLongPriceWithFeeCost,
			"initShortPriceCost":        positionRecord.InitShortPriceCost,
			"initShortPriceWithFeeCost": positionRecord.InitShortPriceWithFeeCost,
			"netPosition":               positionRecord.NetPosition,
			"longAvailablePosition":     positionRecord.LongAvailablePosition,
			"shortAvailablePosition":    positionRecord.ShortAvailablePosition,
			"longPriceCost":             positionRecord.LongPriceCost,
			"longPriceWithFeeCost":      positionRecord.LongPriceWithFeeCost,
			"longPriceCNYWithFeeCost":   positionRecord.LongPriceCNYWithFeeCost,
			"shortPriceCost":            positionRecord.ShortPriceCost,
			"shortPriceWithFeeCost":     positionRecord.ShortPriceWithFeeCost,
			"shortPriceCNYWithFeeCost":  positionRecord.ShortPriceCNYWithFeeCost,
			"longAvgPrice":              positionRecord.LongAvgPrice,
			"longAvgPriceWithFee":       positionRecord.LongAvgPriceWithFee,
			"shortAvgPrice":             positionRecord.ShortAvgPrice,
			"shortAvgPriceWithFee":      positionRecord.ShortAvgPriceWithFee,
			"buyOrderLeftQty":           positionRecord.BuyOrderLeftQty,
			"sellOrderLeftQty":          positionRecord.SellOrderLeftQty,
			"buyOrderLeftCost":          positionRecord.BuyOrderLeftCost,
			"sellOrderLeftCost":         positionRecord.SellOrderLeftCost,
		}
	} else {
		positionData = map[string]interface{}{
			"key":                      positionRecord.Key,
			"contractBaseDate":         positionRecord.ContractBaseDate,
			"netPosition":              positionRecord.NetPosition,
			"longAvailablePosition":    positionRecord.LongAvailablePosition,
			"shortAvailablePosition":   positionRecord.ShortAvailablePosition,
			"longPriceCost":            positionRecord.LongPriceCost,
			"longPriceWithFeeCost":     positionRecord.LongPriceWithFeeCost,
			"longPriceCNYWithFeeCost":  positionRecord.LongPriceCNYWithFeeCost,
			"shortPriceCost":           positionRecord.ShortPriceCost,
			"shortPriceWithFeeCost":    positionRecord.ShortPriceWithFeeCost,
			"shortPriceCNYWithFeeCost": positionRecord.ShortPriceCNYWithFeeCost,
			"longAvgPrice":             positionRecord.LongAvgPrice,
			"longAvgPriceWithFee":      positionRecord.LongAvgPriceWithFee,
			"shortAvgPrice":            positionRecord.ShortAvgPrice,
			"shortAvgPriceWithFee":     positionRecord.ShortAvgPriceWithFee,
			"buyOrderLeftQty":          positionRecord.BuyOrderLeftQty,
			"sellOrderLeftQty":         positionRecord.SellOrderLeftQty,
			"buyOrderLeftCost":         positionRecord.BuyOrderLeftCost,
			"sellOrderLeftCost":        positionRecord.SellOrderLeftCost,
		}
	}

	return
}

func (a *FiccFutOrderPositionAdapter) GetQuotaNotEnoughHandler() func(*schema.TradeOrder, *domain_error.Error) {

	return func(tradeOrder *schema.TradeOrder, de *domain_error.Error) {

		er := de.Err
		if er == nil {
			return
		}

		resultCode := er.Error()

		tradeOrder.ExtendAttrMap["limitCheckResult"] = resultCode

		de.Refine(domain_error.ERROR, tradeOrder)
	}
}
