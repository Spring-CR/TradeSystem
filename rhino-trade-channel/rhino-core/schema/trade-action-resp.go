package schema

import (
	"log"
	"rhino-common/utils/byteutils"

	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type TradeActionResp struct {
	ID                    int64
	OrderID               string `sql:"size: 188"`
	ClOrdID               string `sql:"unique: pk_tar, size: 188"`
	OrigClOrdID           string `sql:"size: 188"`
	ExecID                string `sql:"index: pk_tar, size: 188"`
	ExecRefID             string `sql:"size: 188"`
	ExecTransType         string `sql:"size: 2"`
	ExecType              string `sql:"size: 2"`
	OrdStatus             string `sql:"size: 2"`
	OrdRejReason          string
	CxlRejResponseTo      string
	ExecRestatementReason string
	Account               string `sql:"size: 64"`
	Symbol                string `sql:"size: 64"`
	SymbolSfx             string `sql:"size: 8"`
	SecurityID            string `sql:"size: 64"`
	IDSource              string `sql:"size: 2"`
	SecurityType          string `sql:"size: 16"`
	Side                  string `sql:"size: 2"`
	OpenClose             string `sql:"size: 2"`
	OrderQty              float64
	CashOrderQty          float64
	OrdType               string `sql:"size: 2"`
	Price                 float64
	Currency              string `sql:"size: 4"`
	EffectiveTime         string `sql:"size: 64"`
	ExpireTime            string `sql:"size: 64"`
	LastShares            int64
	LastPx                float64
	LeavesQty             int64
	CumQty                int64
	AvgPx                 float64
	TransactTime          int64
	ExchangeTradeDate     string `sql:"size: 16"`
	MsgTime               int64
	DBInsertTime          int64
	MsgSeq                int64
	ChannelCode           string `sql:"index: pk_tar, size: 32"`
	RawMsg                string `sql:"-"`
	AppOrdID              string `sql:"-"`
	ExtendAttr            string `sql:"type: MEDIUMTEXT"`
	ExtendAttrMap         map[string]interface{} `sql:"-"`
}

func (o *TradeActionResp) GetCacheKey() string {
	// 因ClOrdID、ExecID都以数字结尾，直接对接不一定唯一
	return o.ClOrdID + o.ChannelCode + o.ExecID
}

func (o *TradeActionResp) RecoverExtendAttrMap() {
	if len(o.ExtendAttr) == 0 {
		o.ExtendAttrMap = make(map[string]interface{})
	} else {
		m := make(map[string]interface{})
		err := json.Unmarshal(byteutils.GetZeroCopyBytes(o.ExtendAttr), &m)
		if err == nil {
			o.ExtendAttrMap = m
		} else {
			log.Println("fail to unmarshal json for resp, data="+o.ExtendAttr)
			o.ExtendAttrMap = make(map[string]interface{})
		}
	}
}

/* Todo:
1、applicationconfig，定义resp的exterAttr
2、resp同步之后，将exterAttr插入sqlite(建表sql、插入记录sql、构造插入记录的传参)
3、归档时将resp的exterAttr插入当日及历史的extend_res(建表sql、插入记录sql、构造插入记录的传参)
*/