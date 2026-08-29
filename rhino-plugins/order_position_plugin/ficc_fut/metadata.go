package ficc_fut

import (
	"fmt"
	"log"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-core/schema"
	"strconv"
	"strings"
)

type PositionRecord struct {
	// 标的属性
	Key                string  `json:"key"`
	Account            int     `json:"account"`
	CounterpartyID     int     `json:"counterpartyID"`
	Counterparty       string  `json:"counterparty"`
	Symbol2            string  `json:"symbol2"`
	SymbolName         string  `json:"symbolName"`
	Currency           string  `json:"currency"`
	PlanCode           string  `json:"planCode"`
	UltraContractCode  string  `json:"ultraContractCode"`
	SecurityExchange   string  `json:"securityExchange"`
	SecurityType       string  `json:"securityType"`
	ProductCode        string  `json:"productCode"`
	ContractMultiplier float64 `json:"contractMultiplier"`
	ExchangeRateCNY    float64 `json:"exchangeRateCNY"`
	LongMarginRatio    float64 `json:"longMarginRatio"`
	ShortMarginRatio   float64 `json:"shortMarginRatio"`
	ExchangeArea       string  `json:"exchangeArea"`
	ContractBaseDate   string  `json:"contractBaseDate"` // 底仓上场日
	// 统计属性
	InitNetPosition           float64 `json:"initNetPosition"`           // 合约底仓+T-1历史订单累计计算的持仓
	InitLongPriceCost         float64 `json:"initLongPriceCost"`         // 初始多头持仓成本
	InitLongPriceWithFeeCost  float64 `json:"initLongPriceWithFeeCost"`  // 初始多头（含费）持仓成本
	InitShortPriceCost        float64 `json:"initShortPriceCost"`        // 初始空头持仓成本
	InitShortPriceWithFeeCost float64 `json:"initShortPriceWithFeeCost"` // 初始空头（含费）持仓成本
	NetPosition               float64 `json:"netPosition"`               // 当前净持仓
	LongAvailablePosition     float64 `json:"longAvailablePosition"`     // 多头可用持仓
	ShortAvailablePosition    float64 `json:"shortAvailablePosition"`    // 空头可用持仓
	LongPriceCost             float64 `json:"longPriceCost"`             // 多头持仓成本
	LongPriceWithFeeCost      float64 `json:"longPriceWithFeeCost"`      // 多头（含费）持仓成本
	ShortPriceCost            float64 `json:"shortPriceCost"`            // 空头持仓成本
	ShortPriceWithFeeCost     float64 `json:"shortPriceWithFeeCost"`     // 空头（含费）持仓成本
	LongAvgPrice              float64 `json:"longAvgPrice"`              // 多头持仓均价
	LongAvgPriceWithFee       float64 `json:"longAvgPriceWithFee"`       // 多头（含费）持仓均价
	ShortAvgPrice             float64 `json:"shortAvgPrice"`             // 空持仓均价
	ShortAvgPriceWithFee      float64 `json:"shortAvgPriceWithFee"`      // 空头（含费）持仓均价
	BuyOrderLeftQty           float64 `json:"buyOrderLeftQty"`           // 多头持仓手数
	SellOrderLeftQty          float64 `json:"sellOrderLeftQty"`          // 空头持仓手数
	BuyOrderLeftCost          float64 `json:"buyOrderLeftCost"`          // 买单剩余委托规模（本币）
	SellOrderLeftCost         float64 `json:"sellOrderLeftCost"`         // 卖单剩余委托规模（本币）
	LongPriceCNYWithFeeCost   float64 `json:"longPriceCNYWithFeeCost"`   // 多头（含费, 本币）持仓成本
	ShortPriceCNYWithFeeCost  float64 `json:"shortPriceCNYWithFeeCost"`  // 空头（含费, 本币）持仓成本
}

func (a *FiccFutOrderPositionAdapter) loadOrConstructPositionRecord(tradeOrder *schema.TradeOrder) (positionRecord *PositionRecord) {

	log.Printf("loadOrConstructPositionRecord for order :%s\n", tradeOrder.AppOrdID)

	key := fmt.Sprintf("%v-%v", tradeOrder.Account, tradeOrder.ExtendAttrMap["symbol2"])

	var datasetSuffix string
	switch tradeOrder.ChannelCode {
	case "olts-fut":
		datasetSuffix = "Dms"
	case "stars-fut":
		datasetSuffix = "Ovs"
	}

	valList, _, _ := a.applicationCfg.GetAutoSyncRepo().Get("PositionBase"+datasetSuffix, key)

	if len(valList) == 0 {

		extendAttrMap := tradeOrder.ExtendAttrMap
		account, _, _ := attrutil.GetAttrValue(extendAttrMap, "account", enum.AttrValueType_INT)
		counterpartyID := account
		counterparty, _, _ := attrutil.GetAttrValue(extendAttrMap, "counterparty", enum.AttrValueType_STRING)
		symbol2, _, _ := attrutil.GetAttrValue(extendAttrMap, "symbol2", enum.AttrValueType_STRING)
		symbolName, _, _ := attrutil.GetAttrValue(extendAttrMap, "symbolName", enum.AttrValueType_STRING)
		currency, _, _ := attrutil.GetAttrValue(extendAttrMap, "currency", enum.AttrValueType_STRING)
		planCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "planCode", enum.AttrValueType_STRING)
		ultraContractCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "ultraContractCode", enum.AttrValueType_STRING)
		securityExchange, _, _ := attrutil.GetAttrValue(extendAttrMap, "securityExchange2", enum.AttrValueType_STRING)
		securityType, _, _ := attrutil.GetAttrValue(extendAttrMap, "futType", enum.AttrValueType_STRING)
		productCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "productCode", enum.AttrValueType_STRING)
		contractMultiplier, _, _ := attrutil.GetAttrValue(extendAttrMap, "contractMultiplier", enum.AttrValueType_FLOAT)
		exchangeRateCNY, _, _ := attrutil.GetAttrValue(extendAttrMap, "exchangeRateCNY", enum.AttrValueType_FLOAT)
		if exchangeRateCNY.(float64) == 0.0 {
			exchangeRateCNY = 1.0
		}
		longMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
		shortMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
		tradeDate, _, _ := attrutil.GetAttrValue(extendAttrMap, "tradeDate", enum.AttrValueType_INT)
		contractBaseDate := strconv.Itoa(tradeDate.(int))

		// priceCNYWithFee, ok, _ := attrutil.GetAttrValue(extendAttrMap, "priceCNYWithFee", enum.AttrValueType_FLOAT)
		// if !ok || priceCNYWithFee == 0 {
		// 	priceCNYWithFee = tradeOrder.Price
		// }

		positionRecord = &PositionRecord{
			Key:                key,
			Account:            account.(int),
			CounterpartyID:     counterpartyID.(int),
			Counterparty:       counterparty.(string),
			Symbol2:            symbol2.(string),
			SymbolName:         symbolName.(string),
			Currency:           currency.(string),
			PlanCode:           planCode.(string),
			UltraContractCode:  ultraContractCode.(string),
			SecurityExchange:   securityExchange.(string),
			SecurityType:       securityType.(string),
			ProductCode:        productCode.(string),
			ContractMultiplier: contractMultiplier.(float64),
			ExchangeRateCNY:    exchangeRateCNY.(float64),
			LongMarginRatio:    longMarginRatio.(float64),
			ShortMarginRatio:   shortMarginRatio.(float64),
			ExchangeArea:       strings.ToLower(datasetSuffix),
			ContractBaseDate:   contractBaseDate,
		}

		// 重演的时候会去自动计算，因此不需要在这里计算BuyOrderLeftQty、BuyOrderLeftCost、SellOrderLeftQty、SellOrderLeftCost
		// switch tradeOrder.Side {
		// case "1":
		// 	positionRecord.BuyOrderLeftQty = tradeOrder.OrderQty
		// 	positionRecord.BuyOrderLeftCost = priceCNYWithFee.(float64) * positionRecord.BuyOrderLeftQty * positionRecord.ContractMultiplier
		// case "2":
		// 	positionRecord.SellOrderLeftQty = tradeOrder.OrderQty
		// 	positionRecord.SellOrderLeftCost = priceCNYWithFee.(float64) * positionRecord.SellOrderLeftQty * positionRecord.ContractMultiplier
		// }

	} else {

		extendAttrMap := valList[0]

		_, positionRecord = a.parsePositionRecordFromReposRecord(extendAttrMap)

		if positionRecord.LongMarginRatio == 0 || positionRecord.ShortMarginRatio == 0 {
			extendAttrMap := tradeOrder.ExtendAttrMap
			longMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
			shortMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
			positionRecord.LongMarginRatio = longMarginRatio.(float64)
			positionRecord.ShortMarginRatio = shortMarginRatio.(float64)
		}

		if positionRecord.ExchangeArea == "" {
			positionRecord.ExchangeArea = strings.ToLower(datasetSuffix)
		}
	}

	a.marginExposureManager.CalculateMarginExposure(positionRecord, tradeOrder)

	return
}

func (a *FiccFutOrderPositionAdapter) parsePositionRecordFromReposRecord(extendAttrMap map[string]interface{}) (k string, positionRecord *PositionRecord) {
	key, metadata := a.ParsePositionRecordFromReposRecord(extendAttrMap)
	if key == "" {
		return
	}
	return key, metadata.(*PositionRecord)
}

func (a *FiccFutOrderPositionAdapter) ParsePositionRecordFromReposRecord(extendAttrMap map[string]interface{}) (k string, positionRecord interface{}) {

	key, _, _ := attrutil.GetAttrValue(extendAttrMap, "Key", enum.AttrValueType_STRING)
	account, _, _ := attrutil.GetAttrValue(extendAttrMap, "Account", enum.AttrValueType_INT)
	counterpartyID := account
	counterparty, _, _ := attrutil.GetAttrValue(extendAttrMap, "Counterparty", enum.AttrValueType_STRING)
	symbol2, _, _ := attrutil.GetAttrValue(extendAttrMap, "Symbol2", enum.AttrValueType_STRING)
	symbolName, _, _ := attrutil.GetAttrValue(extendAttrMap, "SymbolName", enum.AttrValueType_STRING)
	currency, _, _ := attrutil.GetAttrValue(extendAttrMap, "Currency", enum.AttrValueType_STRING)
	planCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "PlanCode", enum.AttrValueType_STRING)
	ultraContractCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "UltraContractCode", enum.AttrValueType_STRING)
	securityExchange, _, _ := attrutil.GetAttrValue(extendAttrMap, "SecurityExchange", enum.AttrValueType_STRING)
	securityType, _, _ := attrutil.GetAttrValue(extendAttrMap, "SecurityType", enum.AttrValueType_STRING)
	productCode, _, _ := attrutil.GetAttrValue(extendAttrMap, "ProductCode", enum.AttrValueType_STRING)
	contractMultiplier, _, _ := attrutil.GetAttrValue(extendAttrMap, "ContractMultiplier", enum.AttrValueType_FLOAT)
	exchangeRateCNY, _, _ := attrutil.GetAttrValue(extendAttrMap, "ExchangeRateCNY", enum.AttrValueType_FLOAT)
	if exchangeRateCNY.(float64) == 0 {
		exchangeRateCNY = 1.0
	}
	longMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
	shortMarginRatio, _, _ := attrutil.GetAttrValue(extendAttrMap, "marginRatio", enum.AttrValueType_FLOAT)
	exchangeArea, _, _ := attrutil.GetAttrValue(extendAttrMap, "ExchangeArea", enum.AttrValueType_STRING)
	contractBaseDate, _, _ := attrutil.GetAttrValue(extendAttrMap, "ContractBaseDate", enum.AttrValueType_STRING)

	initNetPosition, _, _ := attrutil.GetAttrValue(extendAttrMap, "InitNetPosition", enum.AttrValueType_FLOAT)
	initLongPriceCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "InitLongPriceCost", enum.AttrValueType_FLOAT)
	initLongPriceWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "InitLongPriceWithFeeCost", enum.AttrValueType_FLOAT)
	initShortPriceCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "InitShortPriceCost", enum.AttrValueType_FLOAT)
	initShortPriceWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "InitShortPriceWithFeeCost", enum.AttrValueType_FLOAT)

	netPosition, _, _ := attrutil.GetAttrValue(extendAttrMap, "NetPosition", enum.AttrValueType_FLOAT)
	longAvailablePosition, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongAvailablePosition", enum.AttrValueType_FLOAT)
	shortAvailablePosition, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortAvailablePosition", enum.AttrValueType_FLOAT)
	longPriceCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongPriceCost", enum.AttrValueType_FLOAT)
	longPriceWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongPriceWithFeeCost", enum.AttrValueType_FLOAT)
	longPriceCNYWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongPriceCNYWithFeeCost", enum.AttrValueType_FLOAT)
	shortPriceCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortPriceCost", enum.AttrValueType_FLOAT)
	shortPriceWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortPriceWithFeeCost", enum.AttrValueType_FLOAT)
	shortPriceCNYWithFeeCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortPriceCNYWithFeeCost", enum.AttrValueType_FLOAT)
	longAvgPrice, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongAvgPrice", enum.AttrValueType_FLOAT)
	longAvgPriceWithFee, _, _ := attrutil.GetAttrValue(extendAttrMap, "LongAvgPriceWithFee", enum.AttrValueType_FLOAT)
	shortAvgPrice, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortAvgPrice", enum.AttrValueType_FLOAT)
	shortAvgPriceWithFee, _, _ := attrutil.GetAttrValue(extendAttrMap, "ShortAvgPriceWithFee", enum.AttrValueType_FLOAT)
	buyOrderLeftQty, _, _ := attrutil.GetAttrValue(extendAttrMap, "BuyOrderLeftQty", enum.AttrValueType_FLOAT)
	sellOrderLeftQty, _, _ := attrutil.GetAttrValue(extendAttrMap, "SellOrderLeftQty", enum.AttrValueType_FLOAT)
	buyOrderLeftCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "BuyOrderLeftCost", enum.AttrValueType_FLOAT)
	sellOrderLeftCost, _, _ := attrutil.GetAttrValue(extendAttrMap, "SellOrderLeftCost", enum.AttrValueType_FLOAT)

	if key == "" {
		key = fmt.Sprintf("%v-%v", account, symbol2)
	}

	positionRecord = &PositionRecord{
		Key:                key.(string),
		Account:            account.(int),
		CounterpartyID:     counterpartyID.(int),
		Counterparty:       counterparty.(string),
		Symbol2:            symbol2.(string),
		SymbolName:         symbolName.(string),
		Currency:           currency.(string),
		PlanCode:           planCode.(string),
		UltraContractCode:  ultraContractCode.(string),
		SecurityExchange:   securityExchange.(string),
		SecurityType:       securityType.(string),
		ProductCode:        productCode.(string),
		ContractMultiplier: contractMultiplier.(float64),
		ExchangeRateCNY:    exchangeRateCNY.(float64),
		LongMarginRatio:    longMarginRatio.(float64),
		ShortMarginRatio:   shortMarginRatio.(float64),
		ExchangeArea:       exchangeArea.(string),
		ContractBaseDate:   contractBaseDate.(string),

		InitNetPosition:           initNetPosition.(float64),
		InitLongPriceCost:         initLongPriceCost.(float64),
		InitLongPriceWithFeeCost:  initLongPriceWithFeeCost.(float64),
		InitShortPriceCost:        initShortPriceCost.(float64),
		InitShortPriceWithFeeCost: initShortPriceWithFeeCost.(float64),

		NetPosition:              netPosition.(float64),
		LongAvailablePosition:    longAvailablePosition.(float64),
		ShortAvailablePosition:   shortAvailablePosition.(float64),
		LongPriceCost:            longPriceCost.(float64),
		LongPriceWithFeeCost:     longPriceWithFeeCost.(float64),
		LongPriceCNYWithFeeCost:  longPriceCNYWithFeeCost.(float64),
		ShortPriceCost:           shortPriceCost.(float64),
		ShortPriceWithFeeCost:    shortPriceWithFeeCost.(float64),
		ShortPriceCNYWithFeeCost: shortPriceCNYWithFeeCost.(float64),
		LongAvgPrice:             longAvgPrice.(float64),
		LongAvgPriceWithFee:      longAvgPriceWithFee.(float64),
		ShortAvgPrice:            shortAvgPrice.(float64),
		ShortAvgPriceWithFee:     shortAvgPriceWithFee.(float64),
		BuyOrderLeftQty:          buyOrderLeftQty.(float64),
		SellOrderLeftQty:         sellOrderLeftQty.(float64),
		BuyOrderLeftCost:         buyOrderLeftCost.(float64),
		SellOrderLeftCost:        sellOrderLeftCost.(float64),
	}

	return key.(string), positionRecord
}
