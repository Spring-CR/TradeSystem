package domain_order_validate

import (
	"database/sql"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"strings"
	"sync/atomic"
	"time"
)

var (
	supportHandlInst = map[string]bool{
		string(enum.OrderHandlInst_DMA):    true,
		string(enum.OrderHandlInst_DSA):    true,
		string(enum.OrderHandlInst_MANUAL): true,
	}
	supportSide = map[string]bool{
		string(enum.Side_Buy):             true,
		string(enum.Side_Sell):            true,
		string(enum.Side_BuyMinus):        true,
		string(enum.Side_SellPlus):        true,
		string(enum.Side_SellShort):       true,
		string(enum.Side_SellShortExempt): true,
		string(enum.Side_Undisclosed):     true,
		string(enum.Side_Cross):           true,
		string(enum.Side_CrossShort):      true,
	}
	supportOrdType = map[string]bool{
		string(enum.OrdType_MARKET): true,
		string(enum.OrdType_LIMIT):  true,
	}
	supportCurrency = map[string]bool{
		"":                        true, // 允许留空
		string(enum.Currency_CNY): true,
		string(enum.Currency_HKD): true,
		string(enum.Currency_USD): true,
		string(enum.Currency_EUR): true,
		string(enum.Currency_GBP): true,
		string(enum.Currency_JPY): true,
		string(enum.Currency_CAD): true,
		string(enum.Currency_AUD): true,
		string(enum.Currency_CHF): true,
		string(enum.Currency_SGD): true,
		string(enum.Currency_KRW): true,
		string(enum.Currency_INR): true,
		string(enum.Currency_BRL): true,
		string(enum.Currency_RUB): true,
		string(enum.Currency_ZAR): true,
		string(enum.Currency_SAR): true,
		string(enum.Currency_AED): true,
		string(enum.Currency_QAR): true,
	}
	supportIDSource = map[string]bool{
		"":                                    true,
		string(enum.IDSource_CUSIP):           true,
		string(enum.IDSource_SEDOL):           true,
		string(enum.IDSource_QUIK):            true,
		string(enum.IDSource_ISIN_NUMBER):     true,
		string(enum.IDSource_RIC_CODE):        true,
		string(enum.IDSource_EXCHANGE_SYMBOL): true,
	}
	supportOpenClose = map[string]bool{
		string(enum.OpenClose_CLOSE): true,
		string(enum.OpenClose_OPEN):  true,
		string(enum.OpenClose_ROLLED): true,
		string(enum.OpenClose_FIFO): true,
	}

	// 这里是用于分析延迟
	totalInsertDBDuration   int64
	totalTxInsertDBDuration int64
	totalUpdateMemDuration  int64

	// 缓存一些枚举类型的string转换值，提高程序效率
	strActionType_New      = string(enum.ActionType_New)
	strOrdStatus_Submit    = string(enum.OrdStatus_Submit)
	strOrdFillStatus_Empty = string(enum.OrdFillStatus_Empty)
)

func CheckAndBaseSetDomainOrder(order *schema.TradeOrder, cfg *domain_cfg.TradeChannelCfg, validateOrder func(order *schema.TradeOrder)(channelOrder interface{}, de *domain_error.Error), afterTradeOrderUpsert func(tradeOrder *schema.TradeOrder)) (duplicatedOrder bool, channelOrder interface{}, de *domain_error.Error) {

	// 如果TransactTime未设置，则用当前时间来设置
	if order.TransactTime <= 0 {
		order.TransactTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	}

	if !supportHandlInst[order.HandlInst] {
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_HANDINST_ERR_CODE, nil, order.HandlInst)
		return
	}

	if order.AppOrdID == "" {
		de = domain_error.Build(domain_error.DATA_CONVERT_APPORDID_EMPTY_ERR_CODE, nil)
		return
	}

	symbolDomain := order.Symbol
	if order.IDSource == "" || enum.IDSource(order.IDSource) == enum.IDSource_EXCHANGE_SYMBOL {
		dotIdx := strings.LastIndex(symbolDomain, ".")
		if dotIdx > 0 {
			symbolDomain = symbolDomain[:dotIdx]
			order.SymbolSfx = order.Symbol[dotIdx+1:]
		}

		// 中信QFII和GFGFIX都没有用到这个字段，如果原始值为空，可以不进行设置
		// if order.SecurityID == "" {
		// 	order.SecurityID = symbolDomain
		// }
		// if order.IDSource == "" {
		// 	order.IDSource = string(enum.IDSource_EXCHANGE_SYMBOL)
		// }
	}
	order.Symbol2 = symbolDomain

	if order.Symbol == "" {
		de = domain_error.Build(domain_error.DATA_CONVERT_SYMBOL_EMPTY_ERR_CODE, nil)
		return
	}

	if !supportSide[order.Side] {
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_TRADE_SIDE_ERR_CODE, nil, order.Side)
		return
	}

	if !supportOrdType[order.OrdType] {
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_ORDER_TYPE_ERR_CODE, nil, order.OrdType)
		return
	}

	if !supportCurrency[order.Currency] {
		var supports []string
		for k := range supportCurrency {
			if k == "" {
				continue
			}
			supports = append(supports, k)
		}
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_CURRENCY_ERR_CODE, nil, strings.Join(supports, "、"), order.Currency)
		return
	}

	if !supportIDSource[order.IDSource] {
		var supports = []string{"1(美国和加拿大的证券识别代码-CUSIP)", "2(英国的证券识别代码-SEDOL)", "3(俄罗斯的证券识别代码-QUIK)", "4(国际证券识别号码-ISIN)", "5(路透证券识别代码-RIC)", "8(交易所发布的证券代码)"}
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_ID_SOURCE_ERR_CODE, nil, strings.Join(supports, "、"), order.IDSource)
		return
	}

	if order.OrderQty <= 0 && order.CashOrderQty <= 0 {
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_QTY_ERR_CODE, nil, order.OrderQty, order.CashOrderQty)
		return
	}

	// if order.SecurityID == "" {
	// 	de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_QTY_ERR_CODE, nil, order.OrderQty, order.CashOrderQty)
	// 	return
	// }

	if order.OpenClose != "" && !supportOpenClose[order.OpenClose] {
		var supports = []string{"C(平仓-Close)", "O(开仓-Open)", "R(ROLLED)", "F(FIFO)"}
		de = domain_error.Build(domain_error.DATA_CONVERT_ILLEGAL_OPEN_CLOSE_ERR_CODE, nil, strings.Join(supports, "、"), order.OpenClose)
		return
	}

	// 验证时间有效期性
	curr := timeutil.ConvertTimeToMilliseconds(time.Now())
	if order.ExpireTime > 0 {
		if curr > order.ExpireTime {
			de = domain_error.Build(domain_error.DATA_CONVERT_ORDER_EXPIRED_ERR_CODE, nil, timeutil.ConvertMillisecondsToTime(curr), timeutil.ConvertMillisecondsToTime(order.ExpireTime))
			return
		}
	} else if order.ExpireDate > 0 {
		if curr > order.ExpireDate {
			de = domain_error.Build(domain_error.DATA_CONVERT_ORDER_EXPIRED_ERR_CODE, nil, timeutil.ConvertMillisecondsToTime(curr), timeutil.ConvertMillisecondsToTime(order.ExpireDate))
			return
		}
	}

	// 检查订单channel是否匹配
	if order.ChannelCode != "" && order.ChannelCode != cfg.GetTradeChannel().ChannelCode {
		de = domain_error.Build(domain_error.DATA_CONVERT_TRADE_CHANNEL_UNMATCH_ERR_CODE, nil, order.ChannelCode, cfg.GetTradeChannel().ChannelCode)
		return
	}

	// 执行客制化订单校验
	channelOrder, de = validateOrder(order)
	if de != nil {
		return
	}

	// 要先保存，才有ID
	if order.ClOrdID == "" {

		beginInsertDB := time.Now()
		// 获取当前时间戳
		nowTime := timeutil.ConvertTimeToMilliseconds(time.Now())
		// 该订单是否首次插入数据库的标识
		orderInsert := false
		// 新建委托的ActionType
		newOrderActionType := strActionType_New
		// 设置状态更新时间
		order.OrdStatusUpdateTime = nowTime
		// 设置订单动作为新建委托
		order.LatestActionType = newOrderActionType
		// 设置为提交状态
		order.OrdStatus = strOrdStatus_Submit
		// 设置为未成交状态
		order.OrdFillStatus = strOrdFillStatus_Empty
		// 剩余数量设置为订单的目标数量
		leavesQty := int64(order.OrderQty)
		if leavesQty <= 0 {
			leavesQty = int64(order.CashOrderQty)
		}
		order.LeavesQty = leavesQty

		// 登记入库时间。处理订单草稿首次入库使用DBInsertTime，draft被提交执行也是用DBInsertTime记录。
		order.DBInsertTime = nowTime

		if order.ID <= 0 {

			// 订单创建时间
			order.OrdCreateTime = nowTime
			order.OrdCreator = order.OrdExecUser

			// 在channel保存订单，可以提高整体效率
			err := app_store.InsertTradeOrder(cfg.GetAppDB(), order)
			if err != nil {
				if !dbutil.IsMysqlDuplicateEntryError(err) { // 过滤违反唯一性约束的错误
					de = domain_error.Build(domain_error.DATA_CONVERT_SAVE_NEW_ORDER_ERR_CODE, err)
					return
				} else {
					duplicatedOrder = true
					return
				}
			}
			orderInsert = true

			//log.Printf("=====> order %s insert into db at:%d\n", order.AppOrdID, timeutil.ConvertTimeToMilliseconds(time.Now()))

			// 记录延时
			insertDBDuration := int64(time.Since(beginInsertDB))

			log.Printf("=====> order %s insert into db cost:%d ns\n", order.AppOrdID, insertDBDuration)
			
			atomic.AddInt64(&totalInsertDBDuration, insertDBDuration)
		} else {
			// 在channel保存订单，可以提高整体效率
			err := app_store.UpdateTradeOrderById(cfg.GetAppDB(), order)
			if err != nil {
				de = domain_error.Build(domain_error.DATA_CONVERT_SAVE_NEW_ORDER_ERR_CODE, err)
				return
			}
			orderInsert = false

			//log.Printf("=====> order %s update into db at:%d\n", order.AppOrdID, timeutil.ConvertTimeToMilliseconds(time.Now()))

			// 记录延时
			insertDBDuration := int64(time.Since(beginInsertDB))
			atomic.AddInt64(&totalInsertDBDuration, insertDBDuration)
		}

		// dateStr, err := timeutil.GetDateStrForTimeZone(order.TransactTime, cfg.GetTradeChannel().TimeZone)
		// if err != nil {
		// 	de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		// 	return
		// }
		// order.TradeDate = timeutil.Parse8BitDateStrToNum(dateStr)
		// order.ClOrdID = fmt.Sprintf("%s-%s-%s-%d", dateStr, order.SystemCode, order.BusinessCode, order.ID)
		dateNum := cfg.GetTradeChannelDetails().GetCurrentExchangeDate()
		order.TradeDate = dateNum
		order.ClOrdID = fmt.Sprintf("%d-%s-%s-%d", dateNum, order.SystemCode, order.BusinessCode, order.ID)
		order.DBInsertOnOrdExec = orderInsert

		// 首次写入数据库，要把关联的表也设置好
		//if orderInsert { // 不能区分是否OrderInsert了，都要加入tradeActionLatestResp以及更新内存模型

		if afterTradeOrderUpsert != nil {
			beginUpdateMemDuration := time.Now()
			afterTradeOrderUpsert(order) // 新记录插入
			atomic.AddInt64(&totalUpdateMemDuration, int64(time.Since(beginUpdateMemDuration)))
		}

		beginTxInsertDBDuration := time.Now()

		cfg.GetAutoTx().Input(cfg.GetAppDB(), func(tx *sql.Tx) (de *domain_error.Error) {
			tradeActionLatestResp := &schema.TradeActionLatestResp{
				ActionUser:        order.OrdExecUser,
				ActionTime:        nowTime,
				ActionMsgTime:     nowTime,
				ActionType:        order.LatestActionType,
				ActionKey:         order.AppOrdID, // 对于下单委托，action key直接取AppOrdID
				AppOrdID:          order.AppOrdID,
				RootClOrdID:       order.ClOrdID,
				ClOrdID:           order.ClOrdID,
				ChannelCode:       order.ChannelCode,
				StreamInputMsgSeq: order.MsgSeq,
			}
			err := app_store.InsertTradeActionLatestResp(tx, tradeActionLatestResp)
			if err != nil && !dbutil.IsMysqlDuplicateEntryError(err) { // 排除重复插入的错误
				de = domain_error.Build(domain_error.DATA_CONVERT_INSERT_TRADE_ACTION_LATEST_RESP_ERR_CODE, err, order.ClOrdID)
				return
			}
			return
		}, "", "")

		atomic.AddInt64(&totalTxInsertDBDuration, int64(time.Since(beginTxInsertDBDuration)))
		//}
	}

	if order.ClOrdID == "" {
		de = domain_error.Build(domain_error.DATA_CONVERT_CLORDID_EMPTY_ERR_CODE, nil)
		return
	}

	return
}

func GetAllBDuration() (int64, int64, int64) {
	return totalInsertDBDuration, totalTxInsertDBDuration, totalUpdateMemDuration
}
