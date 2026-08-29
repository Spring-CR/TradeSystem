package order_cache

import (
	"bytes"
	"os"
	"rhino-core/schema"
	"sort"
	"time"
)

func (c *OrderCache) Dump() {

	var tradeActionResps[]*schema.TradeActionResp
	var tradeActionLatestResps[]*schema.TradeActionLatestResp
	var tradeOrders[]*schema.TradeOrder
	c.rootOrdersLock.Lock()
	c.tradeActionRespMapLock.Lock()
	for _, resp := range c.tradeActionRespMap {
		tradeActionResps = append(tradeActionResps, resp.GetTradeActionRespList()...)
		tradeActionLatestResps = append(tradeActionLatestResps, resp.GetTradeActionLatestResp())
	}

	for _, order := range c.rootOrders {
		tradeOrders = append(tradeOrders, order.GetBasicInfo())
	}

	c.tradeActionRespMapLock.Unlock()
	c.rootOrdersLock.Unlock()
	

	sort.Slice(tradeActionResps, func(i, j int) bool {
		if tradeActionResps[i].MsgTime == tradeActionResps[j].MsgTime {
			return tradeActionResps[i].MsgSeq < tradeActionResps[j].MsgSeq
		}
		return tradeActionResps[i].MsgTime < tradeActionResps[j].MsgTime
	})

	sort.Slice(tradeActionLatestResps, func(i, j int) bool {
		return tradeActionLatestResps[i].ActionTime < tradeActionLatestResps[j].ActionTime
	})

	buf := bytes.NewBufferString("")

	buf.WriteString("TradeOrders:\n")
	for _, v := range tradeOrders {
		lineData, _ := json.Marshal(v)
		buf.Write(lineData)
		buf.WriteByte('\n')
	}

	buf.WriteString("\nTradeActionLatestResps:\n")
	for _, v := range tradeActionLatestResps {
		lineData, _ := json.Marshal(v)
		buf.Write(lineData)
		buf.WriteByte('\n')
	}

	buf.WriteString("\nTradeActionResps:\n")
	for _, v := range tradeActionResps {
		lineData, _ := json.Marshal(v)
		buf.Write(lineData)
		buf.WriteByte('\n')
	}

	now := time.Now().Format("20060102150405.000")
	os.WriteFile("/tmp/cachedump-"+now+".txt", buf.Bytes(), 0644)
}