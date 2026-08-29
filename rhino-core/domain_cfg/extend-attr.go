package domain_cfg

import (
	"rhino-common/enum"
	"rhino-core/schema"
)

// 调整属性，除了要修改本文件 extend-attr.go，还需要同步调整 order-achive-table.go 、 mem_db_orm.go， 可以搜索关键字‘增加了四个用户’
func ConfigExtendAttrItemsForTradeOrderExtending(extendAttrItems []*schema.ExtendAttrItem) []*schema.ExtendAttrItem {

	// 添加状态属性
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_ord_status_update_time",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
		Index:         true,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_ord_status",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  2,
		Index:         true,
	})

	// 添加新的动态属性请从这里开始向下加入 =====================>
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_db_insert_time",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
		Index:         true,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{ // 预存指令再执行时可重复设置，所以是一个动态属性
		AttrName:      "f_transact_time",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
		Index:         true,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_reviewer",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  128,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_approve_status",
		AttrValueType: int(enum.AttrValueType_INT),
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_last_shares",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_last_px",
		AttrValueType: int(enum.AttrValueType_FLOAT),
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_leaves_qty",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_cum_qty",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_avg_px",
		AttrValueType: int(enum.AttrValueType_FLOAT),
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_ord_rej_reason",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  1024,
	})
	// 添加新的动态属性请从这里开始向上加入 <=====================

	// id属性，用于关联
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_app_ord_id",
		AttrValueType: int(enum.AttrValueType_STRING),
		Index:         true,
	})
	// 添加新的静态属性请从这里开始向下加入 =====================>
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_cl_ord_id",
		AttrValueType: int(enum.AttrValueType_STRING),
		Index:         true,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_ord_create_time",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_alg_params",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  2048,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_trade_date",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  8,
	})
	//增加：增加了四个用户, f_ord_creator、f_ord_draft_update_user、f_ord_draft_del_user、f_ord_exec_user
	// extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
	// 	AttrName:      "f_ord_creator", 
	// 	AttrValueType: int(enum.AttrValueType_STRING),
	// 	AttrValueLen:  64,
	// })
	// extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
	// 	AttrName:      "f_ord_draft_update_user", 
	// 	AttrValueType: int(enum.AttrValueType_STRING),
	// 	AttrValueLen:  64,
	// })
	// extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
	// 	AttrName:      "f_ord_draft_del_user", 
	// 	AttrValueType: int(enum.AttrValueType_STRING),
	// 	AttrValueLen:  64,
	// })
	// extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
	// 	AttrName:      "f_ord_exec_user", 
	// 	AttrValueType: int(enum.AttrValueType_STRING),
	// 	AttrValueLen:  64,
	// })

	return extendAttrItems
}

// 注意：这个函数应该在ConfigExtendAttrItemsForTradeOrderExtending之后执行
func ConfigExtendAttrItemsForTradeActionRespExtending(extendAttrItems []*schema.ExtendAttrItem) []*schema.ExtendAttrItem {

	// 去掉Order表里的唯一性约束
	for _, extendAttrItem := range extendAttrItems {
		if extendAttrItem.Unique {
			extendAttrItem.Unique = false
			extendAttrItem.Index = true
		}
	}

	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_orig_cl_ord_id",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  188,
		Index:         true,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_exec_id",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  188,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_exec_ref_id",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  188,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_exec_trans_type",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  2,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_exec_type",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  2,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_msg_time",
		AttrValueType: int(enum.AttrValueType_INT),
		AttrValueLen:  20,
	})
	extendAttrItems = append(extendAttrItems, &schema.ExtendAttrItem{
		AttrName:      "f_channel_code",
		AttrValueType: int(enum.AttrValueType_STRING),
		AttrValueLen:  32,
	})

	return extendAttrItems
}
