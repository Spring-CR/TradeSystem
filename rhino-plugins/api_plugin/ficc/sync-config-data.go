package ficc

import (
	"bytes"
	"encoding/csv"
	"log"
	"os"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"strings"
	"time"
)

func (a *TitansFiccAPIAdapter) syncAppConfigData(applicationCfg *domain_cfg.ApplicationCfg, configMap map[string]*schema.ApplicationCfgItem) {

	log.Printf("ApplicationCfgItems.Len=%d\n", len(applicationCfg.GetApplicationCfgItems()))

	buf := bytes.NewBufferString("ConfigItemName,ConfigItemValue\n")
	csvWriter := csv.NewWriter(buf)
	for _, item := range applicationCfg.GetApplicationCfgItems() {
		log.Printf("======> item:%s, value:%s\n", item.ConfigItemName, item.ConfigItemValue)
		err := csvWriter.Write([]string{item.ConfigItemName, item.ConfigItemValue})
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to write csv record for application config")
		}
	}
	csvWriter.Flush()

	systemCode, businessCode := applicationCfg.GetSystemAndBusinessCodes()

	endTimeStr := strings.Split(configMap["TradingTime"].ConfigItemValue, "-")[1]
	syncDate, _, err := computeSyncDate(configMap, endTimeStr)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to computeSyncDate")
	}

	syncLog := &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         syncDate,
		TableName:        "ApplicationConfig",
		SyncType:         int(enum.SyncType_CSV),
		SyncParams:       buf.String(),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}

	err = admin_store.InsertDataSyncLog(applicationCfg.GetCentralDB(), syncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to InsertDataSyncLog for application config")
	}

	a.syncMarginRatioConfigData(systemCode, businessCode, syncDate, applicationCfg, configMap)
}

func (a *TitansFiccAPIAdapter) syncMarginRatioConfigData(systemCode, businessCode, syncDate string, applicationCfg *domain_cfg.ApplicationCfg, configMap map[string]*schema.ApplicationCfgItem) {

	pathMarginConfig := configMap["PathMarginConfig"].ConfigItemValue
	data, err := os.ReadFile(pathMarginConfig)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to ReadFile from " + pathMarginConfig)
	}

	syncLog := &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         syncDate,
		TableName:        "MarginThreshold",
		SyncType:         int(enum.SyncType_CSV),
		SyncParams:       string(data),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}

	err = admin_store.InsertDataSyncLog(applicationCfg.GetCentralDB(), syncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to InsertDataSyncLog for MarginRatio")
	}
}