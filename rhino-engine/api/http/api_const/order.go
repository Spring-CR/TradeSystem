package api_const

const (
	SubRouteSaveDraft               = "/save_draft"                   // 插入或者更新订单草稿
	SubmitOrderForReview            = "/approval_request"             // 提交审批单
	SubRouteGetDraft                = "/get_draft"                    // 根据AppOrdID获得订单草稿
	SubRouteDeleteDraft             = "/delete_draft"                 // 根据AppOrdID删除订单草稿
	SubRouteExecOrder               = "/exec_order"                   // 执行订单
	SubRouteExecOrderById           = "/exec_order_by_id"             // 根据订单id执行订单（比如在订单草稿列表中，选择某个订单直接执行，当中只传订单AppOrdID和用户ID）
	SubRouteApproveAndExecOrderById = "/approve_and_exec_order_by_id" // 审批通过并执行订单
	SubRouteDisapproveOrderById     = "/disapprove_order_by_id"       // 审批拒绝订单
	SubRouteCancelOrder             = "/cancel_order_by_id"           // 撤销订单
	SubRouteCancelReviewingOrder    = "/cancel_reviewing_order_by_id" // 撤销待审批订单
	SubForceArchiving               = "/force_archiving"              // 对当前的订单数据执行强制归档操作
	SubForcePurging                 = "/force_purging"                // 对当前的订单数据执行强制清理操作
	SubDump                         = "/dump"                         // 对当前的当日订单状态执行dump操作
	SubRoutePositions               = "/positions"                    // 查看交易引擎持仓记录
	SubRoutePositionAdj             = "/adjust"                       // 持仓调整
	SubRouteUpdateOrderAttribute    = "/update_order_attribute"       // 更新订单属性
)
