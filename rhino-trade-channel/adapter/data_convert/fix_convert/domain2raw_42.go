package fix_convert

import (
	"rhino-common/utils/timeutil"
	"rhino-core/schema"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	newordersingle42 "github.com/quickfixgo/fix42/newordersingle"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

func DomainOrderToRawOrderFix42(order *schema.TradeOrder) (rawOrder newordersingle42.NewOrderSingle) {

	clOrdId := field.NewClOrdID(order.ClOrdID)
	handlInst := field.NewHandlInst(enum.HandlInst(order.HandlInst))
	symbol := field.NewSymbol(order.Symbol)
	side := field.NewSide(enum.Side(order.Side))
	transactTime := field.TransactTimeField{
		FIXUTCTimestamp: quickfix.FIXUTCTimestamp{
			Time:      timeutil.ConvertMillisecondsToTime(order.TransactTime),
			Precision: quickfix.Millis,
		},
	}
	ordType := field.NewOrdType(enum.OrdType(order.OrdType))

	// 构建fix4.2订单数据结构
	rawOrder = newordersingle42.New(clOrdId, handlInst, symbol, side, transactTime, ordType)

	// 设置交易账户, Tag 1
	if order.Account != "" {
		rawOrder.SetAccount(order.Account)
	}

	// 设置货币类型, Tag 15
	if order.Currency != "" {
		rawOrder.SetCurrency(order.Currency)
	}

	// 设置证券标识符类型, Tag Tag 22
	if order.IDSource != "" {
		rawOrder.SetIDSource(enum.IDSource(order.IDSource))
	}

	// 设置订单标的数量, Tag 38
	if order.OrderQty > 0 {
		rawOrder.SetOrderQty(decimal.NewFromFloat(order.OrderQty), 2)
	}

	// 设置订单金额数量, Tag 152
	if order.CashOrderQty > 0 {
		rawOrder.SetCashOrderQty(decimal.NewFromFloat(order.CashOrderQty), 2)
	}

	// 设置订单价格, Tag 44 或者 Tag 99
	if order.Price > 0 {
		if order.OrdType == string(enum.OrdType_LIMIT) {
			rawOrder.SetPrice(decimal.NewFromFloat(order.Price), 4)
		} else {
			rawOrder.SetStopPx(decimal.NewFromFloat(order.Price), 4)
		}
	}

	// 设置证券ID, 48
	if order.SecurityID != "" {
		rawOrder.SetSecurityID(order.SecurityID)
	}

	// 设置Text, Tag 58, 预留到GFGFIX适配器实现，用于设置算法参数

	// TimeInForce固定为DAY(0), Tag 59
	rawOrder.SetTimeInForce(enum.TimeInForce_DAY)

	// SymbolSfx没什么用，不用设置了..., Tag 65

	// 设置开平仓标记, Tag Tag 77
	if order.OpenClose != "" {
		rawOrder.SetOpenClose(enum.OpenClose(order.OpenClose))
	}

	// 设置最小执行成交量, Tag 110
	if order.MinQty > 0 {
		rawOrder.SetMinQty(decimal.NewFromFloat(order.MinQty), 0)
	}

	// 设置过期时间， Tag 126
	if order.ExpireTime > 0 {
		rawOrder.SetExpireTime(timeutil.ConvertMillisecondsToTime(order.ExpireTime))
	}

	// 设置证券类型（SetSecurityType）, Tag 167。这里很难完全枚举，和具体的业务场景、柜台都有关系，QG、FIX各版本协议都无法覆盖对方。一个可选的方案是通过securitylib来设置，即通过标的的库的类型来级联设置。

	// 设置生效时间, Tag 168
	if order.EffectiveTime > 0 {
		rawOrder.SetEffectiveTime(timeutil.ConvertMillisecondsToTime(order.EffectiveTime))
	}

	// 设置交易所, Tag 207, 目前没有标识，很难统一。而且有些柜台是基于国家地区标识来设置，并不是严格按交易所。

	// 设置合约乘数, Tag 231
	if order.ContractMultiplier > 0 {
		rawOrder.SetContractMultiplier(decimal.NewFromFloat(order.ContractMultiplier), 2)
	}

	return
}
