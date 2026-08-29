package handlers

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/server/middleware"
	"rhino-common/utils/bean"
	"rhino-common/utils/request"
	"rhino-common/utils/timeutil"
	"rhino-core/order_domain"
	"rhino-core/schema"
	"rhino-core/types"
	"rhino-engine/api/api_adapter"
	"rhino-engine/api/http/api_const"
	"rhino-engine/api/http/options"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	taskCount  int64
	json       = jsoniter.ConfigCompatibleWithStandardLibrary
	reviewLock = &sync.Mutex{}
	execLock   = &sync.Mutex{}
)

type OrderHandler struct {
	apiAdapter     api_adapter.APIAdapter
	engine         *order_domain.OrderEngine
	workerCount    int
	workerSharding func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int)
	systemCode     string
	businessCode   string
}

func NewOrderHandler(apiAdapter api_adapter.APIAdapter, engine *order_domain.OrderEngine, workerCount int) *OrderHandler {
	inst := &OrderHandler{apiAdapter: apiAdapter, engine: engine, workerCount: workerCount}
	workSharding := apiAdapter.GetWorkerSharding(workerCount, engine.GetOrderByAppOrdID)
	n := workerCount - 1
	if workSharding == nil {
		workSharding = func(rawMsg []byte, decodeOrder *schema.TradeOrder, cumTaskCount int) (workerIndex int) {
			return cumTaskCount & n // 使用按位与选择 worker
		}
	}
	inst.workerSharding = workSharding
	sys, busi := engine.GetSystemAndBusinessCodes()
	inst.systemCode = sys
	inst.businessCode = busi
	return inst
}

func (h *OrderHandler) SaveOrderDraft(c *gin.Context) {

	msgProps := map[string]interface{}{}
	if !middleware.BindInputOption(c, &msgProps) {
		return
	}

	msgPropsJs, err := json.Marshal(msgProps)
	log.Printf("======>msgPropsJs:%s, err:%v\n", msgPropsJs, err)

	order, de := h.apiAdapter.ConvertNewOrderSingleMessage(nil, msgProps)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	order.SystemCode = h.systemCode
	order.BusinessCode = h.businessCode

	// 订单属性精华和校验
	de = h.apiAdapter.RefineAndValidate(order, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	// 转json文本
	if len(order.ExtendAttrMap) > 0 && len(order.ExtendAttr) == 0 {
		data, _ := json.Marshal(order.ExtendAttrMap)
		order.ExtendAttr = string(data)
	}

	if len(order.AlgParamsMap) > 0 && len(order.AlgParamsMap) == 0 {
		data, _ := json.Marshal(order.AlgParamsMap)
		order.AlgParams = string(data)
	}

	// 查找缓存的draft，合并extendAttr
	found := false
	cachedOrder, ok := h.engine.GetOrderByAppOrdID(order.AppOrdID)
	if ok {
		for k, v := range cachedOrder.ExtendAttrMap {
			if _, ok := order.ExtendAttrMap[k]; !ok {
				found = true
				order.ExtendAttrMap[k] = v
			}
			log.Printf("cached key:%v, cached value:%v\n", k, v)
		}
	} else {
		log.Printf("no cached order for id %s\n", order.AppOrdID)
	}
	if found {
		data, _ := json.Marshal(order.ExtendAttrMap)
		order.ExtendAttr = string(data)
	}

	de = h.engine.GetOrderAcceptor().AcceptOrderDraft(order, enum.ActionType_Draft)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}

func (h *OrderHandler) SubmitOrderForReview(c *gin.Context) {

	msgProps := map[string]interface{}{}
	if !middleware.BindInputOption(c, &msgProps) {
		return
	}

	msgPropsJs, err := json.Marshal(msgProps)
	log.Printf("======> SubmitOrderForReview msgPropsJs:%s, err:%v\n", msgPropsJs, err)

	order, de := h.apiAdapter.ConvertNewOrderSingleMessage(nil, msgProps)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	order.SystemCode = h.systemCode
	order.BusinessCode = h.businessCode

	// 订单属性精华和校验
	de = h.apiAdapter.RefineAndValidate(order, false)
	if de != nil && h.apiAdapter.ErrorCouldBeIgnoreAfterReview(de) {
		de.Refine(domain_error.WARNING, order)
		de = nil
	}
	if middleware.ProcessDomainError(de, c) {
		return
	}

	// 额度检查
	pm := h.engine.GetOrderOrchestrator().GetPositionManager()
	if pm != nil && pm.GetQuotaNotEnoughHandler() != nil {
		sufficient, de := pm.HasSufficientQuota(order)
		if de != nil && !sufficient {
			pm.GetQuotaNotEnoughHandler()(order, de)
		}
	}

	// 转json文本
	if len(order.ExtendAttrMap) > 0 && len(order.ExtendAttr) == 0 {
		data, _ := json.Marshal(order.ExtendAttrMap)
		order.ExtendAttr = string(data)
	}

	if len(order.AlgParamsMap) > 0 && len(order.AlgParamsMap) == 0 {
		data, _ := json.Marshal(order.AlgParamsMap)
		order.AlgParams = string(data)
	}

	de = h.engine.GetOrderAcceptor().AcceptOrderDraft(order, enum.ActionType_SubmitForReview)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	domain_error.BuildWithDetails(domain_error.WARNING, order, domain_error.ORDER_REVIEW_WARING_CODE, nil, order.AppOrdID)
}

func (h *OrderHandler) DeleteOrderDraft(c *gin.Context) {

	opt := &options.DeleteOrderDraftByIdOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.OrdDraftDelUser == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "OrdDraftDelUser")
		middleware.ProcessDomainError(de, c)
		return
	}

	orderDraftDeletion := &types.ApplicationOrderDraftDeleteRequest{
		ActionUser: opt.OrdDraftDelUser,
		ActionTime: timeutil.ConvertTimeToMilliseconds(time.Now()),
		AppOrdID:   opt.AppOrdID,
	}

	de := h.engine.GetOrderAcceptor().AcceptOrderDraftDeletion(orderDraftDeletion)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}

func (h *OrderHandler) GetOrderDraft(c *gin.Context) {

	id, de := request.GetQueryAsString(c, api_const.ParamId, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	tradeOrder, ok := h.engine.GetOrderByAppOrdID(id)
	if !ok {
		de = domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, id)
		middleware.ProcessDomainError(de, c)
		return
	}

	if tradeOrder.OrdStatus != string(enum.OrdStatus_Draft) {
		de = domain_error.Build(domain_error.ORDER_ISNOT_OF_DRAFT_ERR_CODE, nil, id)
		middleware.ProcessDomainError(de, c)
		return
	}

	tradeOrder.ExtendAttrMap = strToMap(tradeOrder.ExtendAttr)
	tradeOrder.AlgParamsMap = strToMap(tradeOrder.AlgParams)

	appOrd, de := h.apiAdapter.ConvertTradeOrderParams(tradeOrder)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	middleware.ResponseJson(c, appOrd)
}

func (h *OrderHandler) ExecOrder(c *gin.Context) {

	msgProps := map[string]interface{}{}
	if !middleware.BindInputOption(c, &msgProps) {
		return
	}

	order, de := h.apiAdapter.ConvertNewOrderSingleMessage(nil, msgProps)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	if order.AppOrdID == "" {
		de := domain_error.Build(domain_error.APPORDID_EMPTY_ERR_CODE, nil)
		middleware.ProcessDomainError(de, c)
		return
	}
	// 表明是从预存订单发起，检查是否重复提交
	if order.ID > 0 {

		// 执行已存在的未执行订单，需要加锁
		execLock.Lock()
		defer execLock.Unlock()

		tradeOrder, ok := h.engine.GetTraceableOrderByAppOrdID(order.AppOrdID)
		if !ok {
			de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, order.AppOrdID)
			middleware.ProcessDomainError(de, c)
			return
		}

		if tradeOrder.GetBasicInfo().OrdStatus != string(enum.OrdStatus_Draft) {
			de := domain_error.Build(domain_error.ORDER_ISNOT_OF_DRAFT_ERR_CODE, nil, tradeOrder.GetBasicInfo().AppOrdID)
			middleware.ProcessDomainError(de, c)
			return
		}

		// 在这里更新draftTime
		order.OrdDraftUpdateTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	}

	h.doExecOrder(c, order)
}

func (h *OrderHandler) ExecOrderById(c *gin.Context) {

	opt := &options.ExecOrderByIdOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.OrdExecUser == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "OrdExecUser")
		middleware.ProcessDomainError(de, c)
		return
	}

	// 执行已存在的未执行订单，需要加锁
	execLock.Lock()
	defer execLock.Unlock()

	// 从预存订单发起，检查是否重复提交
	tradeOrder, ok := h.engine.GetTraceableOrderByAppOrdID(opt.AppOrdID)
	if !ok {
		de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, opt.AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	if tradeOrder.GetBasicInfo().OrdStatus != string(enum.OrdStatus_Draft) {
		de := domain_error.Build(domain_error.ORDER_ISNOT_OF_DRAFT_ERR_CODE, nil, tradeOrder.GetBasicInfo().AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	tradeOrder.GetBasicInfo().OrdExecUser = opt.OrdExecUser

	h.doExecOrder(c, tradeOrder.GetBasicInfo())
}

func (h *OrderHandler) doExecOrder(c *gin.Context, order *schema.TradeOrder) {

	log.Printf("enter doExecOrder")

	order.SystemCode = h.systemCode
	order.BusinessCode = h.businessCode

	workerIndex := h.workerSharding(nil, order, int(taskCount))
	// 提前设置好workerIndex，后续如果有streamAPI的调用，将绑定到特定的worker处理
	order.WorkerAffinity = workerIndex

	// 订单属性精化和校验
	de := h.apiAdapter.RefineAndValidate(order, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	// 转json文本
	if len(order.ExtendAttrMap) > 0 && len(order.ExtendAttr) == 0 {
		data, _ := json.Marshal(order.ExtendAttrMap)
		order.ExtendAttr = string(data)
	}

	if len(order.AlgParamsMap) > 0 && len(order.AlgParamsMap) == 0 {
		data, _ := json.Marshal(order.AlgParamsMap)
		order.AlgParams = string(data)
	}

	// 更新ExtendAttrMap
	_, de = h.apiAdapter.ConvertTradeOrderParams(order)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	duplicatedOrder, de := h.engine.GetOrderAcceptor().AcceptNewOrderSingleRequest(order)
	if de == nil && duplicatedOrder {
		de = domain_error.Build(domain_error.DUPLICATE_ORDER_ERR_CODE, nil, order.AppOrdID)
	}
	if middleware.ProcessDomainError(de, c) {
		return
	}
	// 任务数自增
	atomic.AddInt64(&taskCount, 1)
}

func (h *OrderHandler) ApproveAndExecOrderById(c *gin.Context) {

	reviewLock.Lock()
	defer reviewLock.Unlock()

	opt := &options.ExecOrderByIdOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.OrdExecUser == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "OrdExecUser")
		middleware.ProcessDomainError(de, c)
		return
	}

	// 从预存订单发起，检查是否重复提交
	tradeOrder, ok := h.engine.GetTraceableOrderByAppOrdID(opt.AppOrdID)
	if !ok {
		de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, opt.AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	if tradeOrder.GetBasicInfo().OrdStatus != string(enum.OrdStatus_PendingReview) {
		de := domain_error.Build(domain_error.ORDER_ISNOT_OF_PENDING_REVIEW_ERR_CODE, nil, tradeOrder.GetBasicInfo().AppOrdID, "审批通过")
		middleware.ProcessDomainError(de, c)
		return
	}

	tradeOrder.GetBasicInfo().TransactTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	tradeOrder.GetBasicInfo().OrdExecUser = opt.OrdExecUser
	tradeOrder.GetBasicInfo().Reviewer = opt.OrdExecUser
	tradeOrder.GetBasicInfo().ReviewTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	tradeOrder.GetBasicInfo().ApproveStatus = int(enum.ApproveStatus_Approved)
	if len(opt.Metadata) > 0 {
		for k, v := range opt.Metadata {
			tradeOrder.GetBasicInfo().ExtendAttrMap[k] = v
		}
	}

	h.doExecOrder(c, tradeOrder.GetBasicInfo())
}

func (h *OrderHandler) DisapproveOrderById(c *gin.Context) {

	reviewLock.Lock()
	defer reviewLock.Unlock()

	opt := &options.ReviewOrderByIdOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.Reviewer == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "Reviewer")
		middleware.ProcessDomainError(de, c)
		return
	}

	// 从预存订单发起，检查是否重复提交
	tradeOrder, ok := h.engine.GetTraceableOrderByAppOrdID(opt.AppOrdID)
	if !ok {
		de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, opt.AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	if tradeOrder.GetBasicInfo().OrdStatus != string(enum.OrdStatus_PendingReview) {
		log.Printf("===>tradeOrder.GetBasicInfo().OrdStatus:%v\n", tradeOrder.GetBasicInfo().OrdStatus)
		de := domain_error.Build(domain_error.ORDER_ISNOT_OF_PENDING_REVIEW_ERR_CODE, nil, tradeOrder.GetBasicInfo().AppOrdID, "审批拒绝")
		middleware.ProcessDomainError(de, c)
		return
	}

	tradeOrder.GetBasicInfo().OrdExecUser = ""
	tradeOrder.GetBasicInfo().Reviewer = opt.Reviewer
	tradeOrder.GetBasicInfo().ReviewTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	tradeOrder.GetBasicInfo().ApproveStatus = int(enum.ApproveStatus_Rejected)

	if len(opt.Metadata) > 0 {
		for k, v := range opt.Metadata {
			tradeOrder.GetBasicInfo().ExtendAttrMap[k] = v
		}
	}

	newOrder := &schema.TradeOrder{}
	bean.Copy(tradeOrder.GetBasicInfo()).To(newOrder)

	de := h.engine.GetOrderAcceptor().AcceptOrderDraft(newOrder, enum.ActionType_ReviewCompleted)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}

func (h *OrderHandler) CancelReviewingOrderById(c *gin.Context) {

	reviewLock.Lock()
	defer reviewLock.Unlock()

	opt := &options.CancelOrderOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.OrdCancelUser == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "OrdCancelUser")
		middleware.ProcessDomainError(de, c)
		return
	}

	// 从预存订单发起，检查是否重复提交
	tradeOrder, ok := h.engine.GetTraceableOrderByAppOrdID(opt.AppOrdID)
	if !ok {
		de := domain_error.Build(domain_error.CANNOT_FIND_ORDER_BY_APP_ORD_ID_ERR_CODE, nil, opt.AppOrdID)
		middleware.ProcessDomainError(de, c)
		return
	}

	if tradeOrder.GetBasicInfo().OrdStatus != string(enum.OrdStatus_PendingReview) {
		log.Printf("===>tradeOrder.GetBasicInfo().OrdStatus:%v\n", tradeOrder.GetBasicInfo().OrdStatus)
		de := domain_error.Build(domain_error.ORDER_ISNOT_OF_PENDING_REVIEW_ERR_CODE, nil, tradeOrder.GetBasicInfo().AppOrdID, "撤回待审批订单")
		middleware.ProcessDomainError(de, c)
		return
	}

	tradeOrder.GetBasicInfo().OrdExecUser = ""
	tradeOrder.GetBasicInfo().Reviewer = ""
	tradeOrder.GetBasicInfo().ReviewTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	tradeOrder.GetBasicInfo().ApproveStatus = int(enum.ApproveStatus_NotSubmit)

	newOrder := &schema.TradeOrder{}
	bean.Copy(tradeOrder.GetBasicInfo()).To(newOrder)

	de := h.engine.GetOrderAcceptor().AcceptOrderDraft(newOrder, enum.ActionType_CancelForReview)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}

func strToMap(str string) map[string]interface{} {
	m := make(map[string]interface{})
	if str == "" {
		return m
	}
	err := json.Unmarshal([]byte(str), &m)
	domain_error.ProcessSevereError(false, 0, nil, err, "fail to conver str to map")
	return m
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {

	opt := &options.CancelOrderOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	if opt.AppOrdID == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "AppOrdID")
		middleware.ProcessDomainError(de, c)
		return
	}

	if opt.OrdCancelUser == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "OrdCancelUser")
		middleware.ProcessDomainError(de, c)
		return
	}

	orderCancel := &types.ApplicationOrderCancelRequest{
		ActionUser: opt.OrdCancelUser,
		ActionTime: timeutil.ConvertTimeToMilliseconds(time.Now()),
		AppOrdID:   opt.AppOrdID,
		ActionKey:  opt.ActionKey,
	}

	_, de := h.engine.GetOrderAcceptor().AcceptOrderCancelRequest(orderCancel)
	if middleware.ProcessDomainError(de, c) {
		return
	}
}

func (h *OrderHandler) ForceArchiving(c *gin.Context) {
	h.engine.ForceArchiving()
}

func (h *OrderHandler) ForcePurging(c *gin.Context) {
	h.engine.ForcePurging()
}

func (h *OrderHandler) Dump(c *gin.Context) {
	h.engine.Dump()
}

func (h *OrderHandler) Positions(c *gin.Context) {

	log.Println("go to get Positions")

	positionManager := h.engine.GetOrderOrchestrator().GetPositionManager()
	if positionManager == nil {
		log.Println("positionManager is not config")
		return
	}
	result := positionManager.Dump()
	middleware.ResponseJson(c, result)
}
