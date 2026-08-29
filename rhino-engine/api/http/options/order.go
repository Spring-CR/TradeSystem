package options

// import (
// 	"log"
// 	"rhino-core/schema"

// 	jsoniter "github.com/json-iterator/go"
// )

// var (
// 	json = jsoniter.ConfigCompatibleWithStandardLibrary
// )

// type ApiTradeOrder struct {
// 	*schema.TradeOrder
// 	ExtendAttr map[string]interface{}
// 	AlgParams  map[string]interface{}
// }

// func (o *ApiTradeOrder) ToTradeOrder() *schema.TradeOrder {
// 	if o.ExtendAttr != nil {
// 		jsExtendAttr, err := json.Marshal(o.ExtendAttr)
// 		if err != nil {
// 			log.Printf("ToTradeOrder error:%v, o.ExtendAttr:%v\n", err, o.ExtendAttr)
// 		}
// 		o.TradeOrder.ExtendAttr = string(jsExtendAttr)
// 	}
// 	if o.AlgParams != nil {
// 		jsAlgParams, err := json.Marshal(o.AlgParams)
// 		if err != nil {
// 			log.Printf("ToTradeOrder error:%v, o.AlgParams:%v\n", err, o.AlgParams)
// 		}
// 		o.TradeOrder.AlgParams = string(jsAlgParams)
// 	}
// 	return o.TradeOrder
// }

// func NewApiTradeOrder(rawOrder *schema.TradeOrder) *ApiTradeOrder {
// 	apiTradeOrder := &ApiTradeOrder{
// 		TradeOrder: rawOrder,
// 	}
// 	apiTradeOrder.ExtendAttr = strToMap(rawOrder.ExtendAttr)
// 	apiTradeOrder.AlgParams = strToMap(rawOrder.AlgParams)
// 	return apiTradeOrder
// }

// func strToMap(str string) map[string]interface{} {
// 	m := make(map[string]interface{})
// 	json.Unmarshal([]byte(str), m)
// 	return m
// }

type ExecOrderByIdOption struct {
	AppOrdID    string
	OrdExecUser string
	Metadata    map[string]interface{}
}

type ReviewOrderByIdOption struct {
	AppOrdID string
	Reviewer string
	Metadata map[string]interface{}
}

type DeleteOrderDraftByIdOption struct {
	AppOrdID        string
	OrdDraftDelUser string
}

type CancelOrderOption struct {
	AppOrdID      string
	ActionKey     string
	OrdCancelUser string
}
