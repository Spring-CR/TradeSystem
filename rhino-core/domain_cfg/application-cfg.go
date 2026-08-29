package domain_cfg

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-common/utils/bean"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-data/datamap"
	"strings"
	"time"
)

type ApplicationCfg struct {
	application              *schema.Application
	appCfgItems              []*schema.ApplicationCfgItem
	extendAttrItems          []*schema.ExtendAttrItem
	positionAttrItems        []*schema.PositionAttrItem
	tradeActionRespAttrItems []*schema.TradeActionRespAttrItem
	tradeChannels            []*TradeChannelDetails
	appDB                    *sql.DB
	autoTx                   *dbutil.ConcurrentAutoTx
	centralDB                *sql.DB
	autoTxOuputErrChan       chan *domain_error.Error
	autoSyncRepo             *datamap.AutoSyncRepo
	dataSyncEventChan        chan *datamap.DataChangeEvent
}

func NewApplicationCfg(application *schema.Application, appCfgItems []*schema.ApplicationCfgItem, extendAttrItems []*schema.ExtendAttrItem, positionAttrItems []*schema.PositionAttrItem, tradeActionRespAttrItems []*schema.TradeActionRespAttrItem, tradeChannels []*TradeChannelDetails) (cfg *ApplicationCfg, de *domain_error.Error) {
	cfg = &ApplicationCfg{application: application, appCfgItems: appCfgItems, extendAttrItems: extendAttrItems, positionAttrItems: positionAttrItems, tradeActionRespAttrItems: tradeActionRespAttrItems, tradeChannels: tradeChannels}

	// ---------------------- 设置租户数据库 ----------------------
	db, err := sql.Open("mysql", application.DatabaseUrl)
	if err != nil {
		de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		return
	}

	// 参考：https://juejin.cn/post/6844904087427776519
	// chatgpt：https://chatgpt.com/c/671072a7-d9e0-800c-b356-d1e0ba2363e0
	db.SetMaxOpenConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(128) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	cfg.appDB = db
	//cfg.autoTx, cfg.autoTxOuputErrChan = dbutil.NewAutoTx(1*time.Second, 500)
	cfg.autoTx, cfg.autoTxOuputErrChan = dbutil.NewConcurrentAutoTx(16, 1*time.Second, 256)
	//cfg.autoTx, cfg.autoTxOuputErrChan = dbutil.NewConcurrentAutoTx(1, 1*time.Second, 256) // 改为1试试 // 结果太慢了，每秒下单和回报的tps只有900了

	cfg.startAutoTx()
	// ---------------------------------------------------------

	// ---------------------- 设置中央数据库 ----------------------
	db, err = sql.Open("mysql", application.CentralDatabaseUrl)
	if err != nil {
		de = domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		return
	}
	db.SetMaxOpenConns(64)
	db.SetMaxIdleConns(4)
	db.SetConnMaxLifetime(time.Second * 600)

	cfg.centralDB = db
	// ---------------------------------------------------------

	cfg.initAutoSyncRepo()

	return
}

func (c *ApplicationCfg) startAutoTx() {
	go func() {
		for {
			de := <-c.autoTxOuputErrChan
			if de != nil {
				log.Printf("Receive error from autoTx:%s\n", de.ErrorString())
				domain_error.ProcessSevereError(false, 0, de, nil, "error occurs in autoTx of ApplicationCfg")
			}
		}
	}()
	c.autoTx.Start()
}

func (c *ApplicationCfg) GetAppDB() *sql.DB {
	return c.appDB
}

func (c *ApplicationCfg) GetCentralDB() *sql.DB {
	return c.centralDB
}

func (c *ApplicationCfg) GetCentralDBUrl() string {
	return c.application.CentralDatabaseUrl
}

//	func (c *TradeChannelCfg) GetAutoTx() *dbutil.AutoTx {
//		return c.autoTx
//	}
func (c *ApplicationCfg) GetAutoTx() *dbutil.ConcurrentAutoTx {
	return c.autoTx
}

func (c *ApplicationCfg) GetKafkaBrokers() []string {
	var brokers []string
	strs := strings.Split(c.application.KafkaBrokers, ",")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		brokers = append(brokers, str)
	}
	return brokers
}

func (c *ApplicationCfg) GetWorkingDir() string {
	if c.application.WorkingDir == "" {
		//return "/tmp"
		domain_error.ProcessSevereError(true, 5, nil, fmt.Errorf("working dir must not be empty"), "working dir must not be empty")
	}
	return c.application.WorkingDir
}

func (c *ApplicationCfg) GetSystemAndBusinessCodes() (string, string) {
	return c.application.SystemCode, c.application.BusinessCode
}

func (c *ApplicationCfg) GetTradeChannels() []*TradeChannelDetails {
	return c.tradeChannels
}

func (c *ApplicationCfg) GetApiAdapterPath() string {
	return c.application.ApiAdapterPath
}

func (c *ApplicationCfg) GetFixServerAdapterPath() string {
	return c.application.FixServerAdapterPath
}

func (c *ApplicationCfg) GetOrdStatusAdapterPath() string {
	return c.application.OrdStatusAdapterPath
}

func (c *ApplicationCfg) GetOrdPositionAdapterPath() string {
	return c.application.OrdPositionAdapterPath
}

func (c *ApplicationCfg) GetOrdCapitalAdapterPath() string {
	return c.application.OrdCapitalAdapterPath
}

func (c *ApplicationCfg) GetOrdExecutorAdapterPath() string {
	return c.application.OrdExecutorAdapterPath
}

func (c *ApplicationCfg) GetScheduleAdapterPath() string {
	return c.application.ScheduleAdapterPath
}

func (c *ApplicationCfg) GetApiKafkaBrokers() string {
	return c.application.KafkaBrokers
}

func (c *ApplicationCfg) GetHttpAPIPort() int {
	return c.application.HttpAPIPort
}

func (c *ApplicationCfg) GetTradeChannelReqTopic() string {
	sys, busi := c.GetSystemAndBusinessCodes()
	return fmt.Sprintf("%s-%s-channel-req", sys, busi)
}

func (c *ApplicationCfg) GetExtendAttrItems() []*schema.ExtendAttrItem {
	var extendAttrItems []*schema.ExtendAttrItem
	for _, v := range c.extendAttrItems {
		newExtendAttrItem := &schema.ExtendAttrItem{}
		bean.Copy(v).To(newExtendAttrItem)
		extendAttrItems = append(extendAttrItems, newExtendAttrItem)
	}
	//return c.extendAttrItems
	return extendAttrItems
}

func (c *ApplicationCfg) GetPositionAttrItems() []*schema.PositionAttrItem {
	var positionAttrItems []*schema.PositionAttrItem
	for _, v := range c.positionAttrItems {
		newPositionAttrItem := &schema.PositionAttrItem{}
		bean.Copy(v).To(newPositionAttrItem)
		positionAttrItems = append(positionAttrItems, newPositionAttrItem)
	}
	//return c.extendAttrItems
	return positionAttrItems
}

func (c *ApplicationCfg) GetTradeActionRespAttrItems() []*schema.TradeActionRespAttrItem {
	var tradeActionRespAttrItems []*schema.TradeActionRespAttrItem
	for _, v := range c.tradeActionRespAttrItems {
		newAttrItem := &schema.TradeActionRespAttrItem{}
		bean.Copy(v).To(newAttrItem)
		tradeActionRespAttrItems = append(tradeActionRespAttrItems, newAttrItem)
	}
	//return c.extendAttrItems
	return tradeActionRespAttrItems
}

// 获取可以执行归档的开始时间和最晚可执行归档的时间
func (c *ApplicationCfg) GetTimeRangeForDataArchiving() (beginTimeSecond, endTimeSencond int, err error) {
	beginTimeSecond, err = timeutil.GetCumulativeSecondsFromSimpleTimeString(c.application.DataArchiveCnBeginTime)
	if err != nil {
		return
	}
	endTimeSencond, err = timeutil.GetCumulativeSecondsFromSimpleTimeString(c.application.DataArchiveCnLatestTime)
	if err != nil {
		return
	}

	if c.application.IsDSTSensitive {

		isDst := timeutil.IsDST()

		// 判断是否处于夏令时
		if isDst {
			log.Println("当前处于夏令时（DST）")
		} else {
			log.Println("当前处于冬令时（Standard Time）")
			beginTimeSecond += 3600
			endTimeSencond += 3600
		}
	}

	return
}

func (c *ApplicationCfg) GetApiToken() string {
	return c.application.ApiToken
}

func (c *ApplicationCfg) GetApplicationCfgItems() []*schema.ApplicationCfgItem {
	return c.appCfgItems
}

func (c *ApplicationCfg) GetApplicationCfgItemMap() map[string]*schema.ApplicationCfgItem {
	m := map[string]*schema.ApplicationCfgItem{}
	for _, item := range c.appCfgItems {
		m[item.ConfigItemName] = item
	}
	return m
}

func (c *ApplicationCfg) initAutoSyncRepo() {
	if c.application.DataRepoConfigPath == "" {
		return
	}
	data, err := os.ReadFile(c.application.DataRepoConfigPath)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to initAutoSyncRepo")
	}
	log.Printf("table s config data: %s\n", data)

	config := &datamap.DataMapConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to initAutoSyncRepo")
	}

	for _, tableConfig := range config.TableConfigs {
		data, _ := json.Marshal(tableConfig)
		log.Printf("======>table %s config:%s\n", tableConfig.TableAlias, data)
	}

	var eventCh chan *datamap.DataChangeEvent
	if c.application.PublishDataSyncEvent {
		eventCh = make(chan *datamap.DataChangeEvent, 1024000)
	}
	c.dataSyncEventChan = eventCh
	autoSyncRepo := datamap.NewAutoSyncRepo(config, c.GetCentralDB(), c.GetAppDB(), 0, eventCh)
	autoSyncRepo.Start()
	c.autoSyncRepo = autoSyncRepo
}

func (c *ApplicationCfg) GetAutoSyncRepo() *datamap.AutoSyncRepo {
	return c.autoSyncRepo
}

func (c *ApplicationCfg) GetDataSyncEventChan() chan *datamap.DataChangeEvent {
	return c.dataSyncEventChan
}
