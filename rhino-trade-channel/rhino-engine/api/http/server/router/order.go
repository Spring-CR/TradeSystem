package router

import (
	"rhino-engine/api/http/api_const"
	"rhino-engine/api/http/server/handlers"

	"github.com/gin-gonic/gin"
)

func setOrderRouter(e *gin.Engine, handler *handlers.OrderHandler) {

	r := e.Group(api_const.RootRouteOrder)
	{
		r.POST(api_const.SubRouteSaveDraft, handler.SaveOrderDraft)
		r.POST(api_const.SubRouteDeleteDraft, handler.DeleteOrderDraft)
		r.GET(api_const.SubRouteGetDraft, handler.GetOrderDraft)
		r.POST(api_const.SubRouteExecOrder, handler.ExecOrder)
		r.POST(api_const.SubRouteExecOrderById, handler.ExecOrderById)
		r.POST(api_const.SubRouteCancelOrder, handler.CancelOrder)
		r.POST(api_const.SubmitOrderForReview, handler.SubmitOrderForReview)
		r.POST(api_const.SubRouteApproveAndExecOrderById, handler.ApproveAndExecOrderById)
		r.POST(api_const.SubRouteDisapproveOrderById, handler.DisapproveOrderById)
		r.POST(api_const.SubRouteCancelReviewingOrder, handler.CancelReviewingOrderById)
		r.GET(api_const.SubForceArchiving, handler.ForceArchiving)
		r.GET(api_const.SubForcePurging, handler.ForcePurging)
		r.GET(api_const.SubDump, handler.Dump)
		r.GET(api_const.SubRoutePositions, handler.Positions)
		r.POST(api_const.SubRoutePositions+"/"+api_const.SubRoutePositionAdj, handler.AdjustPosition)
		r.POST(api_const.SubRouteUpdateOrderAttribute, handler.UpdateOrderAttribute)
	}
}
