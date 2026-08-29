package fix

// 作者：林春泉

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/jsonutil"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/app_store"
	"rhino-core/types"
	"rhino-trade-channel/adapter/config"
	"rhino-trade-channel/adapter/store/fix_store"
	"strconv"
	"strings"
	"time"

	"github.com/manucorporat/try"
	"github.com/quickfixgo/tag"

	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
)

/*
*
1、每天需要重新建一次session
2、开市期间启动，需要检查是否有断层消息（灾难恢复）
*/
type GenericFIXChannel struct {
	cfg                  *domain_cfg.TradeChannelCfg
	msgSeqGen            *domain_cfg.MsgSeqGen
	senderCompID         string
	targetCompID         string
	initiator            *quickfix.Initiator
	tradeClient          *TradeClient
	onTradeActionResp    func(tradeActionResp *schema.TradeActionResp)
	cleanTradeActionResp func(tradeActionResp *schema.TradeActionResp) (ignore bool)
	refineOrderCancelID  func(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string)
	tradeActionRespBuf   chan *schema.TradeActionResp
}

func NewGenericFIXChannel(cfg *domain_cfg.TradeChannelCfg, onTradeActionResp func(tradeActionResp *schema.TradeActionResp), cleanTradeActionResp func(tradeActionResp *schema.TradeActionResp) (ok bool), refineOrderCancelID func(rawTargetClOrdID string, rawClOrdID string, order *schema.TradeOrder) (targetClOrdID, clOrdID string)) (fixChannel *GenericFIXChannel, de *domain_error.Error) {

	var configFileContent []byte
	configFileContent, de = cfg.GetTradeChannelCfgAdapter().ToAppConfig()
	if de != nil {
		return
	}

	msgSeqGen := cfg.GetMsgSeqGen()
	senderSeqNum, _ := msgSeqGen.GetMsgSeqNum()
	if senderSeqNum == 0 {
		// 重置序号
		log.Printf("在允许交易的范围之外建立FIX连接：重置序号!")
		configFileContent = bytes.ReplaceAll(configFileContent, []byte("ResetOnLogon=N"), []byte("ResetOnLogon=Y"))
	}

	// step1: 生成fix客户端的配置文件（为了调试而输出文件，实际上启动fix客户端只需要直接读取配置内容即可）
	configDir := strings.TrimSpace(cfg.GetTradeChannel().ConfigDir)
	if configDir == "" {
		configDir = config.DefaultTradeChannelConfigDir
	}
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}
	configPath := filepath.Join(configDir, "tradeclient.cfg")
	err = os.WriteFile(configPath, configFileContent, 0644)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}

	// step2: 启动fix客户端
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(configFileContent))
	if err != nil {
		de = domain_error.Build(domain_error.FIX_PARSE_CLI_CFG_ERR_CODE, err)
		return
	}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_INIT_LOG_ERR_CODE, err)
		return
	}

	log.Printf("Config onTradeActionResp for channel %s\n", cfg.GetChannelCode())
	fixChannel = &GenericFIXChannel{onTradeActionResp: onTradeActionResp}

	app, err := NewTradeClient(cfg, fixChannel)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_INIT_LOG_ERR_CODE, err)
		return
	}

	// step3: 创建storeFactory，从数据库取数并设置目标消息序号的最大值
	// Todo...
	storeFactory := fix_store.NewAdvanceMemoryStoreFactory(msgSeqGen)

	// step4: 启动iniator
	initiator, err := quickfix.NewInitiator(app, storeFactory, appSettings, fileLogFactory)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_INIT_CLI_ERR_CODE, err)
		return
	}
	err = initiator.Start()
	if err != nil {
		de = domain_error.Build(domain_error.FIX_START_CLI_ERR_CODE, err)
		return
	}

	//fixChannel = &GenericFIXChannel{initiator: initiator, tradeClient: app, cfg: cfg, msgSeqGen: msgSeqGen, tradeActionRespBuf: make(chan *schema.TradeActionResp, 4096*2)}
	fixChannel.initiator = initiator
	fixChannel.tradeClient = app
	fixChannel.cfg = cfg
	fixChannel.msgSeqGen = msgSeqGen
	fixChannel.tradeActionRespBuf = make(chan *schema.TradeActionResp, 4096*2)

	// 传入函数
	fixChannel.tradeClient.processExecutionReport = fixChannel.processExecutionReport
	fixChannel.tradeClient.processOrderCancelReject = fixChannel.processOrderCancelReject
	fixChannel.cleanTradeActionResp = cleanTradeActionResp
	fixChannel.refineOrderCancelID = refineOrderCancelID

	// 开启回报消费协程
	fixChannel.startProcessExecutionReportFromBuf()

	for _, cfgItem := range cfg.GetTradeChannelCfgItems() {
		if cfgItem.ConfigItemName == "TargetCompID" {
			fixChannel.targetCompID = cfgItem.ConfigItemValue
		} else if cfgItem.ConfigItemName == "SenderCompID" {
			fixChannel.senderCompID = cfgItem.ConfigItemValue
		}
	}

	// step5: 等待session就绪
	for !app.IsSessionReady() {
		configData, _ := os.ReadFile(configPath)
		log.Printf("FIX session is not ready, sleep 5 seconds, config file:%s\n", configData)
		time.Sleep(5 * time.Second)
	}

	return
}

// func (c *GenericFIXChannel) AddTradeActionRespListener(onTradeActionResp func(tradeActionResp *schema.TradeActionResp)) {
// 	c.onTradeActionResp = onTradeActionResp
// }

// Todo 后续改成基于原子操作动态获取配置
func (c *GenericFIXChannel) GetChannelConfig() *domain_cfg.TradeChannelCfg {
	return c.cfg
}

// 处理回包
func (c *GenericFIXChannel) processExecutionReport(msg *quickfix.Message) {

	// 在多并发中，需要考虑消息插入的序号递增问题，避免旧消息（后到）覆盖新消息（先到）
	m := msg.Body

	orderID, _ := m.GetString(tag.OrderID)
	clOrdID, _ := m.GetString(tag.ClOrdID)
	origClOrdID, _ := m.GetString(tag.OrigClOrdID)
	execID, _ := m.GetString(tag.ExecID)
	execRefID, _ := m.GetString(tag.ExecRefID)
	execTransType, _ := m.GetString(tag.ExecTransType)
	execType, _ := m.GetString(tag.ExecType)
	ordStatus, _ := m.GetString(tag.OrdStatus)
	ordRejReason, _ := m.GetString(tag.OrdRejReason)
	text, _ := m.GetString(tag.Text)
	execRestatementReason, _ := m.GetString(tag.ExecRestatementReason)
	account, _ := m.GetString(tag.Account)
	symbol, _ := m.GetString(tag.Symbol)
	symbolSfx, _ := m.GetString(tag.SymbolSfx)
	securityID, _ := m.GetString(tag.SecurityID)
	idSource, _ := m.GetString(tag.IDSource)
	securityType, _ := m.GetString(tag.SecurityType)
	side, _ := m.GetString(tag.Side)
	openClose, _ := m.GetString(tag.OpenClose)
	_orderQty, _ := m.GetString(tag.OrderQty)
	_cashOrderQty, _ := m.GetString(tag.CashOrderQty)
	ordType, _ := m.GetString(tag.OrdType)
	_price, _ := m.GetString(tag.Price)
	currency, _ := m.GetString(tag.Currency)
	effectiveTime, _ := m.GetString(tag.EffectiveTime)
	expireTime, _ := m.GetString(tag.ExpireTime)

	_lastShares, err := m.GetString(tag.LastShares) // 因为无对于fix4.2还是fix4.4，LastShares和LastQty的tag是一样的，都是32
	if err != nil {
		log.Printf("no LastShares in the ExecutionReport, clOrdID=%v, orderID=%v, execID=%v, skip!\n", clOrdID, orderID, execID)
		//return
	}
	_lastPx, _ := m.GetString(tag.LastPx)
	_leavesQty, _ := m.GetString(tag.LeavesQty)
	_cumQty, _ := m.GetString(tag.CumQty)
	_avgPx, _ := m.GetString(tag.AvgPx)
	// log.Printf("_lastShares:%s, _lastPx:%s, leavesQty:%s, _cumQty:%s, _avgPx:%s\n", _lastShares, _lastPx, _leavesQty, _cumQty, _avgPx)

	transactTime, errTime := m.GetTime(tag.TransactTime)

	exchangeTradeDate, _ := m.GetString(tag.TradeDate)

	if ordStatus == string(enum.OrdStatus_Rejected) && len(text) > 0 {
		if len(ordRejReason) > 0 {
			ordRejReason += "-" + text
		} else {
			ordRejReason = text
		}
	}

	h := msg.Header
	if errTime != nil {
		transactTime, _ = h.GetTime(tag.SendingTime)
	}
	msgSeq, _ := h.GetInt(tag.MsgSeqNum)

	tradeActionResp := &schema.TradeActionResp{
		OrderID:               orderID,
		ClOrdID:               clOrdID,
		OrigClOrdID:           origClOrdID,
		ExecID:                execID,
		ExecRefID:             execRefID,
		ExecTransType:         execTransType,
		ExecType:              execType,
		OrdStatus:             ordStatus,
		OrdRejReason:          ordRejReason,
		ExecRestatementReason: execRestatementReason,
		Account:               account,
		Symbol:                symbol,
		SymbolSfx:             symbolSfx,
		SecurityID:            securityID,
		IDSource:              idSource,
		SecurityType:          securityType,
		Side:                  side,
		OpenClose:             openClose,
		OrderQty:              parseFloat64(_orderQty),
		CashOrderQty:          parseFloat64(_cashOrderQty),
		OrdType:               ordType,
		Price:                 parseFloat64(_price),
		Currency:              currency,
		EffectiveTime:         effectiveTime,
		ExpireTime:            expireTime,
		LastShares:            parseInt64(_lastShares),
		LastPx:                parseFloat64(_lastPx),
		LeavesQty:             parseInt64(_leavesQty),
		CumQty:                parseInt64(_cumQty),
		AvgPx:                 parseFloat64(_avgPx),
		TransactTime:          timeutil.ConvertTimeToMilliseconds(transactTime),
		ExchangeTradeDate:     exchangeTradeDate,
		MsgTime:               timeutil.ConvertTimeToMilliseconds(time.Now()),
		MsgSeq:                int64(msgSeq),
		ChannelCode:           c.cfg.GetTradeChannel().ChannelCode,
		RawMsg:                msg.String(),
	}

	// 对成交回报进行清洗
	ignore := c.cleanTradeActionResp(tradeActionResp)
	if ignore {
		log.Printf("======> ignore tradeActionResp: %s\n", msg.Bytes())
		return
	}

	jsonutil.PrintSimple("======> received tradeActionResp:\n", tradeActionResp)
	c.tradeActionRespBuf <- tradeActionResp
}

func (c *GenericFIXChannel) processOrderCancelReject(msg *quickfix.Message) {

	// 在多并发中，需要考虑消息插入的序号递增问题，避免旧消息（后到）覆盖新消息（先到）
	m := msg.Body

	orderID, _ := m.GetString(tag.OrderID)
	clOrdID, _ := m.GetString(tag.ClOrdID)
	origClOrdID, _ := m.GetString(tag.OrigClOrdID)
	ordStatus, _ := m.GetString(tag.OrdStatus)
	account, _ := m.GetString(tag.Account)
	cxlRejReason, _ := m.GetString(tag.CxlRejReason)
	cxlRejResponseTo, _ := m.GetString(tag.CxlRejResponseTo)
	text, _ := m.GetString(tag.Text)
	transactTime, _ := m.GetTime(tag.TransactTime)

	h := msg.Header
	if transactTime.Unix() <= 0 {
		transactTime, _ = h.GetTime(tag.SendingTime)
	}
	msgSeq, _ := h.GetInt(tag.MsgSeqNum)

	tradeActionResp := &schema.TradeActionResp{
		OrderID:          orderID,
		ClOrdID:          clOrdID, // 这里的clOrdID，并不是订单号，而是OrderCancelRequest的请求ID，具有唯一性
		OrigClOrdID:      origClOrdID,
		OrdStatus:        ordStatus,
		OrdRejReason:     types.OrderCancelRejectPrefix + strings.TrimSpace(cxlRejReason) + " " + strings.TrimSpace(text),
		CxlRejResponseTo: cxlRejResponseTo,
		Account:          account,
		TransactTime:     timeutil.ConvertTimeToMilliseconds(transactTime),
		MsgTime:          timeutil.ConvertTimeToMilliseconds(time.Now()),
		MsgSeq:           int64(msgSeq),
		ChannelCode:      c.cfg.GetTradeChannel().ChannelCode,
		RawMsg:           msg.String(),
	}

	// 对成交回报进行清洗
	c.cleanTradeActionResp(tradeActionResp)

	c.tradeActionRespBuf <- tradeActionResp
}

func (c *GenericFIXChannel) startProcessExecutionReportFromBuf() {
	go func() {
		for {
			tradeActionResp := <-c.tradeActionRespBuf
			try.This(func() {
				// 插入tradeActionResp
				c.cfg.GetAutoTx().Input(c.cfg.GetAppDB(), func(tx *sql.Tx) (de *domain_error.Error) {

					// 设置插入时间
					tradeActionResp.DBInsertTime = timeutil.ConvertTimeToMilliseconds(time.Now())

					err := app_store.InsertTradeActionResp(tx, tradeActionResp)
					if err != nil && !dbutil.IsMysqlDuplicateEntryError(err) { // 排除重复插入的错误

						de = domain_error.Build(domain_error.TRADE_RESP_INSERT_TO_DB_ERR_CODE, err, tradeActionResp.RawMsg)
						return
					}
					return
				}, "", "")

				if c.onTradeActionResp != nil {
					c.onTradeActionResp(tradeActionResp)
				} else {
					domain_error.ProcessSevereError(true, 5, nil, nil, "onTradeActionResp function is not defined!")
				}
			}).Catch(func(err try.E) {
				log.Printf("error occur while processing execution report! error:%v\n", err)
				de := domain_error.Build(domain_error.GENERIC_ERR_CODE, fmt.Errorf("error occur while processing execution report! error:%v", err))
				domain_error.ReportIfErrorHappen(de)
			})
		}
	}()
}

func parseFloat64(str string) float64 {
	if len(str) == 0 {
		return 0.0
	}
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Printf("parseFloat64 error:%v\n", err)
	}
	return val
}

func parseInt64(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	idx := strings.Index(str, ".")
	if idx > 0 {
		str = str[:idx]
	}
	val, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Printf("parseInt64 error:%v\n", err)
		//domain_error.ProcessSevereError(false, 0, nil, err, "error in parseInt64")
	}
	return val
}

type header interface {
	Set(f quickfix.FieldWriter) *quickfix.FieldMap
}

func (c *GenericFIXChannel) QueryHeader(h header) {
	h.Set(field.NewTargetCompID(c.targetCompID))
	h.Set(field.NewSenderCompID(c.senderCompID))
}

func (c *GenericFIXChannel) IsSessionReady() bool {
	return c.tradeClient.IsSessionReady()
}
