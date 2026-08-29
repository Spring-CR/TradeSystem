package domain_cfg

import (
	"database/sql"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"time"
)

type DataSyncConfig struct {
	systemCode          string
	businessCode        string
	mrkBeginTime        string
	mrkCloseTime        string
	mrkCloseTimeZone    string
	dataSyncAdapterPath string
	tradeDateServiceUrl string
	tradeDateServiceAID string
	tradeDateServiceSec string
	tableConfigs        []*dbutil.TableConfig
	centralDB           *sql.DB
	appDB               *sql.DB
}

func NewDataSyncConfig(systemCode string, businessCode string, centralDatabaseUrl string, appDatabaseUrl string, tableConfigs []*dbutil.TableConfig, mrkBeginTime, mrkCloseTime string, mrkCloseTimeZone string, dataSyncAdapterPath string, tradeDateServiceUrl string, tradeDateServiceAID string, tradeDateServiceSec string) *DataSyncConfig {

	c := &DataSyncConfig{
		systemCode:          systemCode,
		businessCode:        businessCode,
		tableConfigs:        tableConfigs,
		mrkBeginTime:        mrkBeginTime,
		mrkCloseTime:        mrkCloseTime,
		mrkCloseTimeZone:    mrkCloseTimeZone,
		dataSyncAdapterPath: dataSyncAdapterPath,
		tradeDateServiceUrl: tradeDateServiceUrl,
		tradeDateServiceAID: tradeDateServiceAID,
		tradeDateServiceSec: tradeDateServiceSec,
	}

	db, err := sql.Open("mysql", centralDatabaseUrl)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to create central db")
	}
	db.SetMaxOpenConns(16) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(16) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	c.centralDB = db

	db, err = sql.Open("mysql", appDatabaseUrl)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to create app db")
	}
	db.SetMaxOpenConns(16) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetMaxIdleConns(16) // 这两个值设置为一样会好一些，否则，有时会出现链接有问题
	db.SetConnMaxLifetime(time.Second * 600)

	c.appDB = db

	return c
}

func (c *DataSyncConfig) GetCentralDB() *sql.DB {
	return c.centralDB
}

func (c *DataSyncConfig) GetAppDB() *sql.DB {
	return c.appDB
}

func (c *DataSyncConfig) GetTableConfigs() []*dbutil.TableConfig {
	return c.tableConfigs
}

func (c *DataSyncConfig) GetSystemAndBusinessCode() (string, string) {
	return c.systemCode, c.businessCode
}

func (c *DataSyncConfig) GetMrkCloseTime() (mrkCloseTime string, mrkCloseTimeZone string) {
	if c.mrkCloseTimeZone == "" {
		c.mrkCloseTimeZone = timeutil.CnTimeZoneName
	}
	return c.mrkCloseTime, c.mrkCloseTimeZone
}

func (c *DataSyncConfig) GetMrkBeginTime() (mrkBeginTime string) {
	return c.mrkBeginTime
}


func (c *DataSyncConfig) GetDataSyncAdapterPath() string {
	return c.dataSyncAdapterPath
}

func (c *DataSyncConfig) GetTradeDateServiceUrl() string {
	return c.tradeDateServiceUrl
}

func (c *DataSyncConfig) GetSecretInfo() (appId string, appSecret string) {
	return c.tradeDateServiceAID, c.tradeDateServiceSec
}
