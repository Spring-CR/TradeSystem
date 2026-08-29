package fix_convert

import (
	"rhino-common/utils/timeutil"
	"rhino-core/schema"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	ordercancelrequest42 "github.com/quickfixgo/fix42/ordercancelrequest"
	"github.com/quickfixgo/quickfix"
)

func DomainTradeAction2OrderCancelRequest42(cancelAction *schema.TradeActionLatestResp) (orderCancelRequest ordercancelrequest42.OrderCancelRequest) {

	origClOrdId := field.NewOrigClOrdID(cancelAction.OrigClOrdID)
	clOrdId := field.NewClOrdID(cancelAction.ClOrdID)
	symbol := field.NewSymbol(cancelAction.Symbol)
	side := field.NewSide(enum.Side(cancelAction.Side))
	transactTime := field.TransactTimeField{
		FIXUTCTimestamp: quickfix.FIXUTCTimestamp{
			Time:      timeutil.ConvertMillisecondsToTime(cancelAction.TransactTime),
			Precision: quickfix.Millis,
		},
	}

	orderCancelRequest = ordercancelrequest42.New(origClOrdId, clOrdId, symbol, side, transactTime)
	
	return
}
