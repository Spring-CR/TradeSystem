package stream

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-core/order_domain"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-engine/api/api_adapter"
)

// 如果业务要求下单的有序性，建议使用stream API；如果使用rpc，因客户可以即时获得response，应在应用层保证有序性。
// 支持自定义sharding到指定的worker。
type streamAPIWorker struct {
	apiAdapter   api_adapter.APIAdapter
	engine       *order_domain.OrderEngine
	client       StreamClient
	jobChan      chan *IngressMessage
	systemCode   string
	businessCode string
}

func newStreamApIWorker(apiAdapter api_adapter.APIAdapter, engine *order_domain.OrderEngine, client StreamClient, jobBuffer int) *streamAPIWorker {
	w := &streamAPIWorker{
		apiAdapter: apiAdapter,
		engine:     engine,
		client:     client,
		jobChan:    make(chan *IngressMessage, jobBuffer),
	}
	sys, busi := w.engine.GetSystemAndBusinessCodes()
	w.systemCode = sys
	w.businessCode = busi
	w.start()
	return w
}

func (w *streamAPIWorker) start() {
	go func() {
		for {
			msg := <-w.jobChan
			//log.Printf("=====> Consumed message In worker: %s, offset: %d, at:%d\n", msg.Data, msg.MsgSeq, timeutil.ConvertTimeToMilliseconds(time.Now()))
			msgType, msgProps, msgObj, de := w.apiAdapter.DistinguishIngressMessage(msg.Data)
			if de != nil {
				domain_error.ReportIfErrorHappen(de)
			} else {
				w.onMessageReceived(msgType, msgProps, msg.Data, msgObj, msg.MsgSeq, msg.WorkerAffinity)
			}
			//log.Printf("finish process message, offset:%d\n", msg.MsgSeq)
		}
	}()
}

func (w *streamAPIWorker) onMessageReceived(msgType enum.ApiMessageType, msgProps map[string]interface{}, rawMsg []byte, msgObj interface{}, msgSeq int64, workerAffinity int) {
	switch msgType {
	case enum.ApiMessageType_NewOrderSingle:
		var tradeOrder *schema.TradeOrder
		var de *domain_error.Error
		var ok bool
		if msgObj == nil {
			tradeOrder, de = w.apiAdapter.ConvertNewOrderSingleMessage(rawMsg, msgProps)
			if w.processsNewOrderSingleError(tradeOrder, de) {
				return
			}
		} else {
			tradeOrder, ok = msgObj.(*schema.TradeOrder)
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("bad NewOrderSingleMessage: %s", rawMsg), "errors occur while parse NewOrderSingleMessage")
				return
			}
		}

		tradeOrder.SystemCode = w.systemCode
		tradeOrder.BusinessCode = w.businessCode
		tradeOrder.MsgSeq = msgSeq
		tradeOrder.WorkerAffinity = workerAffinity

		// RefineAndValidate往后移，因为该方法可能依赖MsgSeq等更多的属性
		de = w.apiAdapter.RefineAndValidate(tradeOrder, true)
		if w.processsNewOrderSingleError(tradeOrder, de) {
			return
		}

		//jsonData, _ := json.MarshalIndent(tradeOrder, "", "  ")
		//log.Printf("jsonData of tradeOrder:%s\n", jsonData)
		//log.Printf("=====> AcceptNewOrderSingleRequest of tradeOrder:%s, at:%d\n", tradeOrder.AppOrdID, timeutil.ConvertTimeToMilliseconds(time.Now()))
		_, de = w.engine.GetOrderAcceptor().AcceptNewOrderSingleRequest(tradeOrder)
		if w.processsNewOrderSingleError(tradeOrder, de) {
			return
		}
	case enum.ApiMessageType_OrderCancelRequest:
		var ordCxlReq *types.ApplicationOrderCancelRequest
		var de *domain_error.Error
		var ok bool
		if msgObj == nil {
			ordCxlReq, de = w.apiAdapter.ConvertOrderCancelRequestMessage(rawMsg, msgProps)
			if w.processsOrderCancelRequestError(ordCxlReq, de) {
				return
			}
		} else {
			ordCxlReq, ok = msgObj.(*types.ApplicationOrderCancelRequest)
			if !ok {
				domain_error.ProcessSevereError(false, 0, nil, fmt.Errorf("bad OrderCancelRequestMessage: %s", rawMsg), "errors occur while parse OrderCancelRequestMessage")
				return
			}
		}
		ordCxlReq.StreamInputMsgSeq = msgSeq
		_, de = w.engine.GetOrderAcceptor().AcceptOrderCancelRequest(ordCxlReq)
		if w.processsOrderCancelRequestError(ordCxlReq, de) {
			return
		}
	}
}

func (w *streamAPIWorker) processsNewOrderSingleError(tradeOrder *schema.TradeOrder, de *domain_error.Error) (errHappen bool) {
	if de != nil {
		errHappen = true
		// 拒单，判定tradeOrder不为空才能发送消息，否则会panic
		if tradeOrder !=nil {
			respMsg := w.apiAdapter.ProcessNewOrderSingleError(tradeOrder, de)
			if len(respMsg) == 0 {
				return
			}
			msgSeq, err := w.client.SendMessage(respMsg)
			if err != nil {
				log.Printf("fail to send response message for NewOrderSingle, respMsg:%s, err:%v, msgSeq:%d\n", respMsg, err, msgSeq)
			}
		}
		return
	}
	return
}

func (w *streamAPIWorker) processsOrderCancelRequestError(ordCxlReq *types.ApplicationOrderCancelRequest, de *domain_error.Error) (errHappen bool) {
	if de != nil {
		errHappen = true
		// 拒单
		log.Printf("ordCxlReq:%v\n", ordCxlReq)
		order, _ := w.engine.GetOrderByAppOrdID(ordCxlReq.AppOrdID)
		respMsg := w.apiAdapter.ProcesssOrderCancelRequestError(ordCxlReq, order, de)
		if len(respMsg) == 0 {
			return
		}
		msgSeq, err := w.client.SendMessage(respMsg)
		if err != nil {
			log.Printf("fail to send response message for OrderCancelRequest, respMsg:%s, err:%v, msgSeq:%d\n", respMsg, err, msgSeq)
		}
		return
	}
	return
}

func (w *streamAPIWorker) process(msg *IngressMessage) {
	w.jobChan <- msg
}
