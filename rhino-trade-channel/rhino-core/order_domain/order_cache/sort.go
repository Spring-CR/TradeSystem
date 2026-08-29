package order_cache

import (
	"rhino-common/enum"
	"rhino-core/schema"
	"rhino-core/types"
	"slices"
	"sort"
)

type orderOrResp struct {
	tradeOrder      *schema.TradeOrder
	tradeActionResp *schema.TradeActionResp
}

func (c *OrderCache) sortOrderAndResps(directOrderMap map[string]*types.TraceableTradeOrder, tradeActionResps []*schema.TradeActionResp) (orderOrResps []*orderOrResp) {

	var tradeOrders []*schema.TradeOrder
	for _, tradeOrder := range directOrderMap {
		if tradeOrder.GetBasicInfo().OrdStatus == string(enum.OrdStatus_Submit) {
			tradeOrders = append(tradeOrders, tradeOrder.GetBasicInfo())
		}
	}

	sort.Slice(tradeOrders, func(i, j int) bool {
		return tradeOrders[i].QuotaValidateTime < tradeOrders[j].QuotaValidateTime
	})

	for _, tradeOrder := range tradeOrders {
		orderOrResps = append(orderOrResps, &orderOrResp{tradeOrder: tradeOrder})
	}

	i := 0
	j := 0
	for _, tradeActionResp := range tradeActionResps {
		n := len(orderOrResps)
		for i = j; i < n; i++ {

			if orderOrResps[i].tradeActionResp != nil {
				continue
			}

			tradeOrder := orderOrResps[i].tradeOrder
			if tradeActionResp.DBInsertTime < tradeOrder.DBInsertTime {
				break
			}
		}
		orderOrResps = slices.Insert(orderOrResps, i, &orderOrResp{tradeActionResp: tradeActionResp})
		j = i + 1
	}

	return
}
