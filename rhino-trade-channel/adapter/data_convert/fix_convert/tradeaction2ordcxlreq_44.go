package fix_convert

import (
	"rhino-common/utils/timeutil"
	"rhino-core/schema"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	ordercancelrequest44 "github.com/quickfixgo/fix44/ordercancelrequest"
	"github.com/quickfixgo/quickfix"
)

func DomainTradeAction2OrderCancelRequest44(cancelAction *schema.TradeActionLatestResp) (orderCancelRequest ordercancelrequest44.OrderCancelRequest) {

	origClOrdId := field.NewOrigClOrdID(cancelAction.OrigClOrdID)
	clOrdId := field.NewClOrdID(cancelAction.ClOrdID)
	//symbol := field.NewSymbol(cancelAction.Symbol)
	side := field.NewSide(enum.Side(cancelAction.Side))
	transactTime := field.TransactTimeField{
		FIXUTCTimestamp: quickfix.FIXUTCTimestamp{
			Time:      timeutil.ConvertMillisecondsToTime(cancelAction.TransactTime),
			Precision: quickfix.Millis,
		},
	}

	orderCancelRequest = ordercancelrequest44.New(origClOrdId, clOrdId, side, transactTime)
	
	return
}
