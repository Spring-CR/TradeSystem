package ficc

import (
	"rhino-common/enum"
	"rhino-core/schema"
)

func (a *TitansFiccOrderPositionAdapter) GetPositionAttributionConfigs() (attrItems []*schema.ExtendAttrItem) {
	attrItems = []*schema.ExtendAttrItem{
		{AttrName: "key", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 512, Unique: true},
		{AttrName: "account", AttrValueType: int(enum.AttrValueType_INT), AttrValueLen: 20, Index: true},
		{AttrName: "counterpartyID", AttrValueType: int(enum.AttrValueType_INT), AttrValueLen: 20, Index: true},
		{AttrName: "counterparty", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 500},
		{AttrName: "symbol", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 200, Index: true},
		{AttrName: "symbolName", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 200},
		{AttrName: "planCode", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 128, Index: true},
		{AttrName: "ultraContractCode", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 128, Index: true},
		{AttrName: "securityExchange", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 20},
		{AttrName: "securityType", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 30},
		{AttrName: "longShort", AttrValueType: int(enum.AttrValueType_STRING), AttrValueLen: 8},
		{AttrName: "baseCashQty", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "totalCashQty", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "baseNotional", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "totalNotional", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "parValue", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "averageCost", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "baseQuotaT0", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "quotaT0", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "baseQuotaT1", AttrValueType: int(enum.AttrValueType_FLOAT)},
		{AttrName: "quotaT1", AttrValueType: int(enum.AttrValueType_FLOAT)},
	}
	return attrItems
}
