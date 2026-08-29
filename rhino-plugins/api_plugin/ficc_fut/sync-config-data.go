package ficc_fut

import (
	"bytes"
	"encoding/csv"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"strings"
	"time"
)

func (a *FiccFurAPIAdapter) syncAppConfigData(applicationCfg *domain_cfg.ApplicationCfg, configMap map[string]*schema.ApplicationCfgItem) {

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

	endTimeStr := strings.Split(configMap["TradingTimeOvs"].ConfigItemValue, "-")[1]
	syncDate, err := computeSyncDate(configMap, endTimeStr)
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
}

const (
	dateOnly = "20060102"
	dateTime = "20060102-15:04"
)
// endTimeStr格式：HH:mm
func computeSyncDate(configMap map[string]*schema.ApplicationCfgItem, endTimeStr string) (syncDate string, err error) {

	var date time.Time

	// 当前时间
	nowTime := time.Now().In(timeutil.CnTimeLocation)
	// 当前日期
	nowDateStr := nowTime.Format(dateOnly)
	// 当日的结束时间
	nowDateEndTime, err := time.ParseInLocation(dateTime, nowDateStr+"-"+endTimeStr, timeutil.CnTimeLocation)
	if err != nil {
		return "", err
	}

	if nowTime.Before(nowDateEndTime) {
		date = nowTime.Add(-24 * time.Hour)
	} else {
		date = nowTime
	}

	syncDate = date.Format(dateOnly)

	log.Printf("======>computeSyncDate, syncDate:%s\n", syncDate)

	return
}