package ficc

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/types"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MarginDescPrefix = "TITANS_FICC|"
)

var (
	CapitalController *CapitalControl
)

type HandleData struct {
	Limit           bool    `json:"limit"`
	KeyCapAcctId    int     `json:"keyCapAcctId"`
	Amount          float64 `json:"amount"`
	ActionType      string  `json:"actionType"`
	MarginDesc      string  `json:"marginDesc"`
	InternalTradeId string  `json:"internalTradeId"`
	Currency        string  `json:"currency"`
}

type CapitalEvent struct {
	AccountEvent   string        `json:"accountEvent"`
	Rollback       bool          `json:"rollback"`
	HandleDataList []*HandleData `json:"handleDataList"`
}

// 处理资金控制相关逻辑
// 1、开仓验资冻结|校验失败，解冻回滚
// 2、开仓冻资（分初始冻资和结束冻资）
// 3、平仓释放挂帐金额 -- ok
// 4、收市全部释放冻结

type CapitalControl struct {
	configMap                         map[string]*schema.ApplicationCfgItem
	kafkaProducer                     *kafka.TenaciousProducerWithBuffer
	hisReturnCapitalMap               map[string]bool // 在单线程中操作，不需要lock
	hisFreezeCapitalForTradeStartMap  map[string]bool // 订单首次冻结，在api中响应，非单线程，
	hisFreezeCapitalForTradeEndMap    map[string]bool // 订单结束时最后一次冻结
	hisRollbackFreezedCapitalMap      map[string]bool // 订单出错时或者订单结束时回滚冻结资金
	hisDFDUnfreezeCapitalMap          map[string]bool // DONEFODAY时解冻所有当前冻结的资金
	lockFreezeCapitalForTradeStartMap *sync.RWMutex   // 确保不能重复冻结资金
	lockRollbackFreezedCapitalMap     *sync.RWMutex   // 确保不能重复解冻资金
}

func NewCapitalControl(applicationCfg *domain_cfg.ApplicationCfg) {
	configMap := applicationCfg.GetApplicationCfgItemMap()
	brokers := strings.Split(configMap["TitansKafkaBroker"].ConfigItemValue, ",")
	topic := configMap["TitansKafkaCapitalTopic"].ConfigItemValue
	p, err := kafka.NewTenaciousProducerWithBuffer("/tmp/capital-control-kafka-msg.log", true, 0, topic, brokers, 64)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to create kafka producer")
	}

	prefixLen := len(MarginDescPrefix)
	hisReturnCapitalMap := map[string]bool{}
	hisFreezeCapitalForTradeStartMap := map[string]bool{}
	hisFreezeCapitalForTradeEndMap := map[string]bool{}
	hisRollbackFreezedCapitalMap := map[string]bool{}
	hisDFDUnfreezeCapitalMap := map[string]bool{}
	_, _, messages, err := kafka.GetHistoricMessages(brokers, topic)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to get hisReturnCapitalMap")
	}
	for _, message := range messages {
		log.Printf("===>kafka msg:%s\n", message)
		capitalEvent := &CapitalEvent{}
		err = json.Unmarshal(message, capitalEvent)
		if err != nil {
			continue
		}
		if capitalEvent.AccountEvent != "TRADING_ORDER" {
			continue
		}
		for _, handleData := range capitalEvent.HandleDataList {
			if !strings.HasPrefix(handleData.MarginDesc, MarginDescPrefix) {
				continue
			}
			switch handleData.ActionType {
			case "SELL_DEAL":
				log.Printf("===>SELL_DEAL map, handleData.MarginDesc[prefixLen:]=%s", handleData.MarginDesc[prefixLen:])
				hisReturnCapitalMap[handleData.MarginDesc[prefixLen:]] = true

			case "FREEZE":

				switch handleData.MarginDesc[prefixLen:] {
				case "FIRST":
					hisFreezeCapitalForTradeStartMap[handleData.InternalTradeId] = true
				case "LAST":
					hisFreezeCapitalForTradeEndMap[handleData.InternalTradeId] = true
				}

			case "UNFREEZE":

				switch handleData.MarginDesc[prefixLen:] {
				case "ROLLBACK":
					hisRollbackFreezedCapitalMap[handleData.InternalTradeId] = true
				case "DFD":
					hisDFDUnfreezeCapitalMap[handleData.InternalTradeId] = true
				}
			}
		}
	}

	CapitalController = &CapitalControl{
		configMap:                         configMap,
		kafkaProducer:                     p,
		hisReturnCapitalMap:               hisReturnCapitalMap,
		hisFreezeCapitalForTradeStartMap:  hisFreezeCapitalForTradeStartMap,
		hisFreezeCapitalForTradeEndMap:    hisFreezeCapitalForTradeEndMap,
		hisRollbackFreezedCapitalMap:      hisRollbackFreezedCapitalMap,
		hisDFDUnfreezeCapitalMap:          hisDFDUnfreezeCapitalMap,
		lockFreezeCapitalForTradeStartMap: &sync.RWMutex{},
		lockRollbackFreezedCapitalMap:     &sync.RWMutex{},
	}
}

// 多头开仓首次冻结资金
func (c *CapitalControl) FreezeCapitalForTradeStart(force bool, order *schema.TradeOrder) (availableAmt, freezedAmount float64, ok bool, err error) {

	c.lockFreezeCapitalForTradeStartMap.RLock()
	key := c.getInternalKeyForOrder(order)
	if c.hisFreezeCapitalForTradeStartMap[key] {
		c.lockFreezeCapitalForTradeStartMap.RUnlock()
		log.Printf("开仓订单%s已完成订单执行前的资金冻结操作，不能重复资金冻结操作\n", key)
		return
	}
	c.lockFreezeCapitalForTradeStartMap.RUnlock()

	c.lockFreezeCapitalForTradeStartMap.Lock()
	defer c.lockFreezeCapitalForTradeStartMap.Unlock()

	if c.hisFreezeCapitalForTradeStartMap[key] {
		log.Printf("开仓订单%s已完成订单执行前的资金冻结操作，不能重复资金冻结操作\n", key)
		return
	}

	accountID := getCapitalAcctIDByOrder(order)
	url := strings.TrimRight(c.configMap["CapitalServiceUrl"].ConfigItemValue, "/")
	availableAmt, err = getAvailableAmountWithRetry(url, accountID, c.configMap["ServiceAppID"].ConfigItemValue, c.configMap["ServiceAppSecret"].ConfigItemValue, 3)

	if err != nil {
		return
	}

	// 券面总额(万元)*10000/债券面值*意向全价*标的初保比例
	freezedAmount = getOrderFreezeAmountOnStart(order)
	if !force && freezedAmount > availableAmt {
		return
	}

	capitalEvent := &CapitalEvent{
		AccountEvent: "TRADING_ORDER",
		HandleDataList: []*HandleData{
			{
				KeyCapAcctId:    accountID,
				Amount:          freezedAmount,
				ActionType:      "FREEZE",
				MarginDesc:      MarginDescPrefix + "FIRST",
				InternalTradeId: key,
				Currency:        getCurrencyByOrder(order),
			},
		},
	}

	jsData, _ := json.Marshal(capitalEvent)
	log.Printf("【订单开单冻结金额：%v】【订单：%s】publish capitalEvent:%s\n", freezedAmount, order.AppOrdID, jsData)

	c.kafkaProducer.SendMessage(jsData)

	ok = true
	c.hisFreezeCapitalForTradeStartMap[key] = true

	return
}

func getOrderFreezeAmountOnStart(order *schema.TradeOrder) float64 {
	freezedAmount := order.OrderQty / getParValueByOrder(order) * getDirtyPriceByOrder(order) * InitMarginRatio
	return freezedAmount
}

func getOrderFreezeAmountOnEnd(order *schema.TradeOrder) float64 {
	freezedAmount := float64(order.CumQty) / getParValueByOrder(order) * getDirtyPriceByOrder(order) * InitMarginRatio
	return freezedAmount
}

func (c *CapitalControl) RollbackFreezedOrderCapitalForTradeError(order *schema.TradeOrder) {

	key := c.getInternalKeyForOrder(order)
	c.lockRollbackFreezedCapitalMap.RLock()
	if c.hisRollbackFreezedCapitalMap[key] {
		c.lockRollbackFreezedCapitalMap.RUnlock()
		log.Printf("订单%s冻结资金已经回滚，不能重复资金解冻回滚操作\n", key)
		return
	}
	c.lockRollbackFreezedCapitalMap.RUnlock()

	c.lockFreezeCapitalForTradeStartMap.RLock()
	if !c.hisFreezeCapitalForTradeStartMap[key] {
		c.lockFreezeCapitalForTradeStartMap.RUnlock()
		log.Printf("找不到订单%s的资金冻结记录，不能执行资金解冻回滚操作\n", key)
		return
	}
	c.lockFreezeCapitalForTradeStartMap.RUnlock()

	c.lockRollbackFreezedCapitalMap.Lock()
	defer c.lockRollbackFreezedCapitalMap.Unlock()

	if c.hisRollbackFreezedCapitalMap[key] {
		log.Printf("订单%s冻结资金已经回滚，不能重复资金解冻回滚操作\n", key)
		return
	}

	// 券面总额(万元)*10000/债券面值*意向全价*标的初保比例
	freezedAmountOnStart := getOrderFreezeAmountOnStart(order)

	capitalEvent := &CapitalEvent{
		AccountEvent: "TRADING_ORDER",
		HandleDataList: []*HandleData{
			{
				KeyCapAcctId:    getCapitalAcctIDByOrder(order),
				Amount:          freezedAmountOnStart,
				ActionType:      "UNFREEZE",
				MarginDesc:      MarginDescPrefix + "ROLLBACK",
				InternalTradeId: key,
				Currency:        getCurrencyByOrder(order),
			},
		},
	}

	jsData, _ := json.Marshal(capitalEvent)
	log.Printf("【回滚冻结金额：%v】【订单：%s】publish capitalEvent:%s\n", freezedAmountOnStart, order.AppOrdID, jsData)

	c.kafkaProducer.SendMessage(jsData)

	c.hisRollbackFreezedCapitalMap[key] = true
}

func (c *CapitalControl) RollbackAndReFreezeCapitalForTradeEnd(order *schema.TradeOrder) {

	var handleDataList []*HandleData
	// 先解冻回滚
	key := c.getInternalKeyForOrder(order)

	c.lockRollbackFreezedCapitalMap.Lock()
	if c.hisRollbackFreezedCapitalMap[key] {
		log.Printf("订单%s冻结资金已经回滚，不能重复资金解冻回滚操作\n", key)
		
		c.lockRollbackFreezedCapitalMap.Unlock()
	} else {
		c.hisRollbackFreezedCapitalMap[key] = true
		c.lockRollbackFreezedCapitalMap.Unlock()

		c.lockFreezeCapitalForTradeStartMap.RLock()
		if c.hisFreezeCapitalForTradeStartMap[key] {
			c.lockFreezeCapitalForTradeStartMap.RUnlock()
			handleDataList = append(handleDataList, &HandleData{
				KeyCapAcctId:    getCapitalAcctIDByOrder(order),
				Amount:          getOrderFreezeAmountOnStart(order),
				ActionType:      "UNFREEZE",
				MarginDesc:      MarginDescPrefix + "ROLLBACK",
				InternalTradeId: key,
				Currency:        getCurrencyByOrder(order),
			})
		} else {
			c.lockFreezeCapitalForTradeStartMap.RUnlock()
		}
	}
	// 再重新冻结
	amount := getOrderFreezeAmountOnEnd(order)
	if c.hisFreezeCapitalForTradeEndMap[key] {
		log.Printf("订单%s结束时的冻结操作已经执行，不能重复操作\n", key)
	} else if amount > 0 { // 需要冻结的金额大于0时，才需要操作
		handleDataList = append(handleDataList, &HandleData{
			KeyCapAcctId:    getCapitalAcctIDByOrder(order),
			Amount:          amount,
			ActionType:      "FREEZE",
			MarginDesc:      MarginDescPrefix + "LAST",
			InternalTradeId: key,
			Currency:        getCurrencyByOrder(order),
		})
	}

	if len(handleDataList) == 0 {
		return
	}

	capitalEvent := &CapitalEvent{
		AccountEvent:   "TRADING_ORDER",
		HandleDataList: handleDataList,
	}

	jsData, _ := json.Marshal(capitalEvent)
	log.Printf("【最终冻结金额：%v】【回滚/冻结订单：%s】publish capitalEvent:%s\n", amount, order.AppOrdID, jsData)

	c.kafkaProducer.SendMessage(jsData)
}

func (c *CapitalControl) UnFreezeAllCapital() {
	freezedMap := make(map[string]float64)
	brokers := strings.Split(c.configMap["TitansKafkaBroker"].ConfigItemValue, ",")
	topic := c.configMap["TitansKafkaCapitalTopic"].ConfigItemValue
	_, _, messages, err := kafka.GetHistoricMessages(brokers, topic)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to get hisReturnCapitalMap")
	}
	for _, message := range messages {
		log.Printf("===>kafka msg:%s\n", message)
		capitalEvent := &CapitalEvent{}
		err = json.Unmarshal(message, capitalEvent)
		if err != nil {
			continue
		}
		if capitalEvent.AccountEvent != "TRADING_ORDER" {
			continue
		}
		for _, handleData := range capitalEvent.HandleDataList {

			if !strings.HasPrefix(handleData.MarginDesc, MarginDescPrefix) {
				continue
			}

			key := handleData.InternalTradeId
			val := freezedMap[key]
			switch handleData.ActionType {
			case "FREEZE":
				val += handleData.Amount
			case "UNFREEZE":
				val -= handleData.Amount
			}
			freezedMap[key] = val
		}
	}
	var handleDataList []*HandleData
	for k, v := range freezedMap {
		idx := strings.LastIndex(k, "|")
		if idx >= 0 {
			_capitalAcctID := k[idx+1:]
			capitalAcctID, err := strconv.Atoi(_capitalAcctID)
			if err != nil {
				domain_error.ProcessSevereError(false, 0, nil, err, "fail to parse capitalAcctID from "+k)
			}
			if capitalAcctID <= 0 {
				continue
			}

			if v <= 0 {
				continue
			}

			handleDataList = append(handleDataList, &HandleData{
				KeyCapAcctId:    capitalAcctID,
				Amount:          v,
				ActionType:      "UNFREEZE",
				MarginDesc:      MarginDescPrefix + "DFD",
				InternalTradeId: k,
				Currency:        "CNY",
			})
		}
	}
	capitalEvent := &CapitalEvent{
		AccountEvent:   "TRADING_ORDER",
		HandleDataList: handleDataList,
	}

	jsData, _ := json.Marshal(capitalEvent)
	log.Printf("【日终解冻全部被冻结金额】publish capitalEvent:%s\n", jsData)

	c.kafkaProducer.SendMessage(jsData)

}

func (c *CapitalControl) getInternalKeyForOrder(order *schema.TradeOrder) string {
	key := fmt.Sprintf("%v|%v|%v", order.AppOrdID, order.ExtendAttrMap["account"], order.ExtendAttrMap["capitalAcctID"])
	return key
}

// 多头平仓释放挂帐金额
func (c *CapitalControl) ReturnCapital(order *schema.TradeOrder, tradeResp *types.TradeActionRespReturn, averageCost float64, parValue float64) {
	key := tradeResp.CurrentTradeActionResp.GetCacheKey()
	if c.hisReturnCapitalMap[key] {
		log.Printf("平仓成交回报%s已经挂过账，不需要重复操作\n", key)
		return
	}
	// T日挂账金额 = 预计释放保证金 + 预计结算盈亏
	// 预计释放保证金 = 本次平仓数量 * 标的实时持仓均价 * 标的初保比例
	// 预计结算盈亏 = 本次平仓数量 *（平仓价格（全价） - 标的实时持仓均价）
	// 标的实时持仓均价：参考wiki2.4.2 债券互换多头持仓计算规则
	// 初保比例：按5%计算
	jsData, _ := json.Marshal(tradeResp.CurrentTradeActionResp)
	log.Printf("收到多头平仓的成交回报，开始计算挂帐金额，成交回报数据结构：%s\n", jsData)
	log.Printf("===>券面金额：%f\n", parValue)
	log.Printf("===>成交面额：%v\n", tradeResp.CurrentTradeActionResp.LastShares)
	qty := float64(tradeResp.CurrentTradeActionResp.LastShares) / parValue
	log.Printf("===>本次平仓数量：%v\n", qty)
	log.Printf("===>标的实时持仓均价：%v\n", averageCost)
	rate := InitMarginRatio
	log.Printf("===>标的初保比例（小时制）：%v\n", rate)
	dirtyPrice := getDirtyPriceByOrder(order)
	log.Printf("===>全价：%v\n", dirtyPrice)
	if averageCost <= 0 {
		log.Printf("===>实时的持仓均价为0, 全价代替实时持仓均价")
		averageCost = dirtyPrice
	}
	amount := qty*averageCost*rate + qty*(dirtyPrice-averageCost)
	log.Printf("===>挂账金额 ：%v\n", amount)

	capitalEvent := &CapitalEvent{
		AccountEvent: "TRADING_ORDER",
		HandleDataList: []*HandleData{
			{
				KeyCapAcctId:    getCapitalAcctIDByOrder(order),
				Amount:          amount,
				ActionType:      "SELL_DEAL",
				MarginDesc:      MarginDescPrefix + key,
				InternalTradeId: c.getInternalKeyForOrder(order),
				Currency:        getCurrencyByOrder(order),
			},
		},
	}

	jsData, _ = json.Marshal(capitalEvent)
	log.Printf("【挂帐金额：%v】【订单：%s】publish capitalEvent:%s\n", amount, order.AppOrdID, jsData)

	c.kafkaProducer.SendMessage(jsData)

	c.hisReturnCapitalMap[key] = true
}

// API响应结构体
type CapitalAccountResponse struct {
	ServiceID string `json:"serviceId"`
	ErrCode   struct {
		Code int    `json:"code"`
		Chs  string `json:"chs"`
		Eng  string `json:"eng"`
	} `json:"errCode"`
	ErrMsg interface{} `json:"errMsg"`
	Data   struct {
		AvailableAmt float64 `json:"availableAmt"`
	} `json:"data"`
	Timestamp int64 `json:"timestamp"`
}

// GetAvailableAmount 根据资金账户获取可用资金
// url: 服务基础URL，如 "http://titans-tst.gf.com.cn/api/titans/margin/1.0.0/capitalAccountManage"
// accountID: 资金账户ID
// appId: 应用ID
// appSecret: 应用密钥
func getAvailableAmount(url string, accountID int, appId, appSecret string) (float64, error) {
	// 构建完整的请求URL
	fullURL := fmt.Sprintf("%s/account/availableAmt/%d", url, accountID)

	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return 0, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("AppId", appId)
	req.Header.Set("AppSecret", appSecret)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("HTTP请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析JSON响应
	var apiResponse CapitalAccountResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return 0, fmt.Errorf("解析JSON响应失败: %v", err)
	}

	// 检查API错误码
	if apiResponse.ErrCode.Code != 200 {
		return 0, fmt.Errorf("API返回错误: %s (代码: %d)", apiResponse.ErrCode.Chs, apiResponse.ErrCode.Code)
	}

	log.Printf("getAvailableAmount, url=%s, accountID=%v, appId=%v, appSecret=%v, availableAmt=%v, body=%s\n", url, accountID, appId, appSecret, apiResponse.Data.AvailableAmt, body)

	return apiResponse.Data.AvailableAmt, nil
}

// 更健壮的版本，包含重试机制
func getAvailableAmountWithRetry(url string, accountID int, appId, appSecret string, maxRetries int) (float64, error) {
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		availableAmt, err := getAvailableAmount(url, accountID, appId, appSecret)
		if err == nil {
			return availableAmt, nil
		}
		lastErr = err

		// 可以在这里添加延迟重试逻辑
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	return 0, fmt.Errorf("经过 %d 次重试后仍然失败，最后错误: %v", maxRetries, lastErr)
}
