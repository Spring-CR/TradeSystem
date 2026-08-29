package ficc_v2

import (
	"fmt"
	"log"
	"math"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"time"
)

type PositionRecord struct {
	// 标的属性
	Account           int
	CounterpartyID    int
	Counterparty      string
	Symbol            string
	SymbolName        string
	Currency          string
	PlanCode          string
	UltraContractCode string
	SecurityExchange  string
	SecurityType      string
	ParValue          float64
	LongMarginRatio   float64
	ShortMarginRatio  float64
	// 统计属性
	InitNetPositionT0          float64 // T+0清算速度的初始净持仓
	InitNetPositionT1          float64 // T+1清算速度的初始净持仓
	NetPositionT0              float64 // T+0清算速度的净持仓
	NetPositionT1              float64 // T+1清算速度的净持仓
	LongAvailablePositionT0    float64 // T+0多头可用持仓
	LongAvailablePositionT1    float64 // T+1多头可用持仓
	ShortAvailablePositionT0   float64 // T+0空头可用持仓
	ShortAvailablePositionT1   float64 // T+1空头可用持仓
	LongCleanPriceCost         float64 // 多头净价持仓成本
	LongDirtyPriceCost         float64 // 多头全价持仓成本
	LongDirtyPriceWithFeeCost  float64 // 多头全价（含费）持仓成本
	ShortCleanPriceCost        float64 // 空头净价持仓成本
	ShortDirtyPriceCost        float64 // 空头全价持仓成本
	ShortDirtyPriceWithFeeCost float64 // 空头全价（含费）持仓成本
	LongAvgCleanPrice          float64 // 多头净价持仓均价
	LongAvgDirtyPrice          float64 // 多头全价持仓均价
	LongAvgDirtyPriceWithFee   float64 // 多头全价（含费）持仓均价
	ShortAvgCleanPrice         float64 // 空头净价持仓均价
	ShortAvgDirtyPrice         float64 // 空头全价持仓均价
	ShortAvgDirtyPriceWithFee  float64 // 空头全价（含费）持仓均价
	MaxLongMarginOccupancy     float64 // 最大多头保证金占用
	MaxShortMarginOccupancy    float64 // 最大空头保证金占用
	MaxMarginOccupancy         float64 // 最大单边保证金占用
}

type PositionBaseRecord struct {
	quantity      float64 // 单位：张，代表证券数量
	baseCashQty   float64 // 单位：元，代表全面总额
	baseDirtyCost float64
	baseCleanCost float64
	baseInitCost  float64
}

func (a *TitansFiccOrderPositionAdapter) loadOrConstructPositionRecord(tradeOrder *schema.TradeOrder) (positionRecord *PositionRecord) {

	log.Printf("loadOrConstructPositionRecord for order :%s\n", tradeOrder.AppOrdID)

	account, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "account", enum.AttrValueType_INT)
	counterparty, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "counterparty", enum.AttrValueType_STRING)
	symbol := tradeOrder.Symbol
	symbolName, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "symbolName", enum.AttrValueType_STRING)
	currency := tradeOrder.Currency
	planCode, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "planCode", enum.AttrValueType_STRING)
	ultraContractCode, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "ultraContractCode", enum.AttrValueType_STRING)
	securityExchange := tradeOrder.SecurityExchange
	parValue, _, _ := attrutil.GetAttrValue(tradeOrder.ExtendAttrMap, "parValue", enum.AttrValueType_FLOAT)

	positionRecord = &PositionRecord{
		Account:           account.(int),
		CounterpartyID:    account.(int),
		Counterparty:      counterparty.(string),
		Symbol:            symbol,
		SymbolName:        symbolName.(string),
		Currency:          currency,
		PlanCode:          planCode.(string),
		UltraContractCode: ultraContractCode.(string),
		SecurityExchange:  securityExchange,
		SecurityType:      "BOND",
		ParValue:          parValue.(float64),
	}

	a.refinePositionRecord(positionRecord)

	a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

	return
}

var (
	yearsHours = 24.0 * 365.0
)

func (a *TitansFiccOrderPositionAdapter) refinePositionRecord(positionRecord *PositionRecord) {

	t0LongPositionBaseRecord := a.getPositionBaseRecord(positionRecord, "LONG", "T0")
	t0ShortPositionBaseRecord := a.getPositionBaseRecord(positionRecord, "SHORT", "T0")
	t1LongPositionBaseRecord := a.getPositionBaseRecord(positionRecord, "LONG", "T1")
	t1ShortPositionBaseRecord := a.getPositionBaseRecord(positionRecord, "SHORT", "T1")

	netPositionT0, _, _, _, _, _, _ := a.nettingPositionBaseRecord(t0LongPositionBaseRecord, t0ShortPositionBaseRecord)
	netPositionT1, avgCleanPriceT1, avgDirtyPriceT1, avgDirtyPriceWithFeeT1, cleanPriceCostT1, dirtyPriceCostT1, dirtyPriceWithFeeCostT1 := a.nettingPositionBaseRecord(t1LongPositionBaseRecord, t1ShortPositionBaseRecord)

	if netPositionT0 > 0 && netPositionT1 < 0 || netPositionT0 < 0 && netPositionT1 > 0 {
		log.Printf("多空头砸差后，符号方向不一致，netPositionT0=%v, netPositionT1=%v\n", netPositionT0, netPositionT1)
		return
	}

	if netPositionT1 >= 0 { // 多头占优
		positionRecord.NetPositionT0 = netPositionT0
		positionRecord.NetPositionT1 = netPositionT1
		positionRecord.LongAvailablePositionT0 = netPositionT0
		positionRecord.LongAvailablePositionT1 = netPositionT1
		positionRecord.ShortAvailablePositionT0 = 0
		positionRecord.ShortAvailablePositionT1 = 0
		positionRecord.LongCleanPriceCost = cleanPriceCostT1
		positionRecord.LongDirtyPriceCost = dirtyPriceCostT1
		positionRecord.LongDirtyPriceWithFeeCost = dirtyPriceWithFeeCostT1
		positionRecord.ShortCleanPriceCost = 0
		positionRecord.ShortDirtyPriceCost = 0
		positionRecord.ShortDirtyPriceWithFeeCost = 0
		positionRecord.LongAvgCleanPrice = avgCleanPriceT1
		positionRecord.LongAvgDirtyPrice = avgDirtyPriceT1
		positionRecord.LongAvgDirtyPriceWithFee = avgDirtyPriceWithFeeT1
		positionRecord.ShortAvgCleanPrice = 0
		positionRecord.ShortAvgDirtyPrice = 0
		positionRecord.ShortAvgDirtyPriceWithFee = 0
	} else { // 空头占优
		positionRecord.NetPositionT0 = netPositionT0
		positionRecord.NetPositionT1 = netPositionT1
		positionRecord.LongAvailablePositionT0 = 0
		positionRecord.LongAvailablePositionT1 = 0
		positionRecord.ShortAvailablePositionT0 = -netPositionT0
		positionRecord.ShortAvailablePositionT1 = -netPositionT1
		positionRecord.LongCleanPriceCost = 0
		positionRecord.LongDirtyPriceCost = 0
		positionRecord.LongDirtyPriceWithFeeCost = 0
		positionRecord.ShortCleanPriceCost = cleanPriceCostT1
		positionRecord.ShortDirtyPriceCost = dirtyPriceCostT1
		positionRecord.ShortDirtyPriceWithFeeCost = dirtyPriceWithFeeCostT1
		positionRecord.LongAvgCleanPrice = 0
		positionRecord.LongAvgDirtyPrice = 0
		positionRecord.LongAvgDirtyPriceWithFee = 0
		positionRecord.ShortAvgCleanPrice = avgCleanPriceT1
		positionRecord.ShortAvgDirtyPrice = avgDirtyPriceT1
		positionRecord.ShortAvgDirtyPriceWithFee = avgDirtyPriceWithFeeT1
	}

	positionRecord.InitNetPositionT0 = positionRecord.NetPositionT0
	positionRecord.InitNetPositionT1 = positionRecord.NetPositionT1

	log.Printf("refinePositionRecord of margin ration for order :%s\n", positionRecord.PlanCode+"-"+positionRecord.Symbol)

	// 计算初保比率
	valList, _, _ := a.applicationCfg.GetAutoSyncRepo().Get("Security", positionRecord.Symbol)
	if len(valList) == 0 {
		log.Printf("cannot get Security for %s\n", positionRecord.Symbol)
		return
	}
	symbolData := valList[len(valList)-1]
	val, _, _ := attrutil.GetAttrValue(symbolData, "MaturityDate", enum.AttrValueType_STRING)
	if len(val.(string)) < 10 {
		log.Printf("cannot get MaturityDate of Security for %v\n", symbolData)
		return
	}
	maturityDate, err := time.ParseInLocation(time.DateOnly, val.(string)[:10], timeutil.CnTimeLocation)
	if err != nil {
		log.Printf("cannot get MaturityDate and ParseInLocation of Security for %v\n", symbolData)
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to parse maturityDate from "+val.(string))
		return
	}

	currDate, err := time.ParseInLocation(time.DateOnly, time.Now().In(timeutil.CnTimeLocation).Format(time.DateOnly), timeutil.CnTimeLocation)
	if err != nil {
		log.Printf("cannot get MaturityDate and ParseInLocation of Security for %v\n", symbolData)
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to parse currDate from "+val.(string))
		return
	}

	maturityYears := maturityDate.Sub(currDate).Hours() / yearsHours
	var maturityPhase int
	if maturityYears <= 2 {
		maturityPhase = 1
	} else if maturityYears > 2 && maturityYears <= 5 {
		maturityPhase = 2
	} else if maturityYears > 5 && maturityYears <= 10 {
		maturityPhase = 3
	} else if maturityYears > 10 && maturityYears <= 30 {
		maturityPhase = 4
	} else if maturityYears > 30 {
		maturityPhase = 5
	}

	bondType, ok, _ := attrutil.GetAttrValue(symbolData, "BondType", enum.AttrValueType_STRING)
	if !ok {
		log.Printf("cannot get bondType of Security for %v\n", symbolData)
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to parse bondType from %v", symbolData))
		return
	}
	planCode := positionRecord.PlanCode

	key := fmt.Sprintf("%v-%v-%v", planCode, bondType, maturityPhase)
	valList, _, _ = a.applicationCfg.GetAutoSyncRepo().Get("MarginThreshold", key)
	if len(valList) == 0 {
		log.Printf("cannot find out MarginThreshold for key: %v\n", key)
		return
	}
	marginThresholdData := valList[len(valList)-1]
	marginPercent, ok, _ := attrutil.GetAttrValue(marginThresholdData, "MarginPercent", enum.AttrValueType_FLOAT)
	if !ok {
		log.Printf("cannot get marginPercent of Security for %v\n", marginThresholdData)
		domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("fail to parse marginPercent from %v", marginThresholdData))
		return
	}
	log.Printf("marginPercent for key: %v, marginPercent:%v\n", key, marginPercent)

	positionRecord.LongMarginRatio = marginPercent.(float64)/100
	positionRecord.ShortMarginRatio = marginPercent.(float64)/100
}

//a.marginExposureManager.CalculateInitMarginExposure(positionRecord)

func (a *TitansFiccOrderPositionAdapter) nettingPositionBaseRecord(longPositionBaseRecord *PositionBaseRecord, shortPositionBaseRecord *PositionBaseRecord) (
	netPosition, avgCleanPrice, avgDirtyPrice, avgDirtyPriceWithFee, cleanPriceCost, dirtyPriceCost, dirtyPriceWithFeeCost float64) {

	netPosition = longPositionBaseRecord.baseCashQty - shortPositionBaseRecord.baseCashQty
	totalQuantity := longPositionBaseRecord.quantity + shortPositionBaseRecord.quantity

	// 添加除零保护
	if math.Abs(totalQuantity) < 1e-10 { // 使用小的阈值判断
		avgCleanPrice = 0
		avgDirtyPrice = 0
		avgDirtyPriceWithFee = 0
	} else {
		avgCleanPrice = (longPositionBaseRecord.baseCleanCost + shortPositionBaseRecord.baseCleanCost) / totalQuantity
		avgDirtyPrice = (longPositionBaseRecord.baseDirtyCost + shortPositionBaseRecord.baseDirtyCost) / totalQuantity
		avgDirtyPriceWithFee = (longPositionBaseRecord.baseInitCost + shortPositionBaseRecord.baseInitCost) / totalQuantity
	}

	if netPosition > 0 {
		quantity := longPositionBaseRecord.quantity - shortPositionBaseRecord.quantity
		cleanPriceCost = avgCleanPrice * quantity
		dirtyPriceCost = avgDirtyPrice * quantity
		dirtyPriceWithFeeCost = avgDirtyPriceWithFee * quantity
	} else {
		quantity := shortPositionBaseRecord.quantity - longPositionBaseRecord.quantity
		cleanPriceCost = avgCleanPrice * quantity
		dirtyPriceCost = avgDirtyPrice * quantity
		dirtyPriceWithFeeCost = avgDirtyPriceWithFee * quantity
	}

	return
}

func (a *TitansFiccOrderPositionAdapter) getPositionBaseRecord(positionRecord *PositionRecord, longShort, settleType string) (positionBaseRecord *PositionBaseRecord) {

	positionBaseRecord = &PositionBaseRecord{}

	key := fmt.Sprintf("%v-%v-%v-%v-%v", positionRecord.Account, positionRecord.PlanCode, positionRecord.Symbol, longShort, settleType)
	valList, _, _ := a.applicationCfg.GetAutoSyncRepo().Get("PositionBase", key)
	if len(valList) == 0 {
		return
	}

	positionBaseRecordMap := valList[len(valList)-1]

	quantity, _, _ := attrutil.GetAttrValue(positionBaseRecordMap, "Quantity", enum.AttrValueType_FLOAT)
	baseCashQty, _, _ := attrutil.GetAttrValue(positionBaseRecordMap, "BaseCashQty", enum.AttrValueType_FLOAT)
	baseDirtyCost, _, _ := attrutil.GetAttrValue(positionBaseRecordMap, "BaseDirtyCost", enum.AttrValueType_FLOAT)
	baseCleanCost, _, _ := attrutil.GetAttrValue(positionBaseRecordMap, "BaseCleanCost", enum.AttrValueType_FLOAT)
	baseInitCost, _, _ := attrutil.GetAttrValue(positionBaseRecordMap, "BaseInitCost", enum.AttrValueType_FLOAT)

	positionBaseRecord.quantity = quantity.(float64)
	positionBaseRecord.baseCashQty = baseCashQty.(float64)
	positionBaseRecord.baseDirtyCost = baseDirtyCost.(float64)
	positionBaseRecord.baseCleanCost = baseCleanCost.(float64)
	positionBaseRecord.baseInitCost = baseInitCost.(float64)

	return
}
