package order_status

import (
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-core/types"
	"strings"
)

func (s *OrderStatusReplica) memInitByOrderTopology() {

	log.Printf("start to init memdb by order topology")

	s.orderCache.FilterOrderByFunction(func(order *types.TraceableTradeOrder) bool{

		orderBasic := order.GetBasicInfo()

		err := s.insertOrder(s.dbWrite, orderBasic)
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Printf("===> init memdb by order topology, ignore duplicate order, clOrdID=%s, appOrdID=%s, id=%d, err=%v\n", orderBasic.ClOrdID, orderBasic.AppOrdID, orderBasic.ID, err)
				return true
			}
			domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while insert order, appOrdIDp=%s, error=%v\n", orderBasic.AppOrdID, err))
		}

		log.Printf("===> init memdb by order topology, insert order, appOrdID=%s, id=%d, err=%v\n", orderBasic.AppOrdID, orderBasic.ID, err)

		_, _, tradeActionResps := types.ExtractTraceableTradeOrder(order)
		for _, tradeActionResp := range tradeActionResps{
			err = s.insertOrderResp(s.dbWrite, orderBasic, tradeActionResp)
			if err != nil {
				if strings.Contains(err.Error(), "UNIQUE constraint failed") {
					log.Printf("===> init memdb by order topology, ignore duplicate order resp, clOrdID=%s, appOrdID=%s, execID=%s\n", orderBasic.ClOrdID, orderBasic.AppOrdID, tradeActionResp.ExecID)
					continue
				}
				domain_error.ProcessSevereError(false, 0, nil, err, fmt.Sprintf("error occurs while insert order resp, appOrdIDp=%s, error=%v\n", orderBasic.AppOrdID, err))
				continue
			}
			log.Printf("===> init memdb by order topology, insert order resp, appOrdID=%s, execID=%s\n", orderBasic.AppOrdID, tradeActionResp.ExecID)
		}

		return true
	})
}
