package handlers

import (
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/server/middleware"
	"rhino-common/utils/request"
	"rhino-common/utils/timeutil"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"rhino-data-sync/api/api_const"
	"rhino-data-sync/api/options"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)

type DataSyncNotifyHandler struct {
	dataSyncCfg *domain_cfg.DataSyncConfig
}

func NewDataSyncNotifyHandler(dataSyncCfg *domain_cfg.DataSyncConfig) *DataSyncNotifyHandler {
	return &DataSyncNotifyHandler{dataSyncCfg}
}

func (h *DataSyncNotifyHandler) Notify(c *gin.Context) {
	opt := &options.DataSyncNotifyOption{}
	if !middleware.BindInputOption(c, opt) {
		return
	}

	systemCode, businessCode := h.dataSyncCfg.GetSystemAndBusinessCode()

	// 参数校验
	if systemCode == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "systemCode")
		middleware.ProcessDomainError(de, c)
		return
	}
	if businessCode == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "businessCode")
		middleware.ProcessDomainError(de, c)
		return
	}
	if opt.SyncDate == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "syncDate")
		middleware.ProcessDomainError(de, c)
		return
	}
	if opt.TableName == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "tableName")
		middleware.ProcessDomainError(de, c)
		return
	}
	if len(opt.SyncParams) == 0 {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, "syncParams")
		middleware.ProcessDomainError(de, c)
		return
	}

	syncParams, _ := json.Marshal(opt.SyncParams)

	// 构造DataSyncLog
	dsl := &schema.DataSyncLog{
		SystemCode:       systemCode,
		BusinessCode:     businessCode,
		SyncDate:         opt.SyncDate,
		TableName:        opt.TableName,
		SyncType:         opt.SyncType,
		SyncParams:       string(syncParams),
		ReportTime:       timeutil.ConvertTimeToMilliseconds(time.Now()),
		FirstSyncTime:    0,
		CurrentSyncTime:  0,
		CompleteSyncTime: 0,
		ExecCount:        0,
		SyncPhase:        int(enum.DataSyncLogPhase_New),
		FailLog:          "",
	}

	err := admin_store.InsertDataSyncLog(h.dataSyncCfg.GetCentralDB(), dsl)
	if err != nil {
		de := domain_error.Build(domain_error.DATABASE_OPERATION_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}

	middleware.ResponseJson(c, dsl)
}

func (h *DataSyncNotifyHandler) GetSyncLog(c *gin.Context) {
	id, de := request.GetQueryAsInt(c, api_const.ParamID, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	syncLog, err := admin_store.GetDataSyncLogById(h.dataSyncCfg.GetCentralDB(), int64(id))
	if err != nil {
		de = domain_error.Build(domain_error.DATABASE_QUERY_ERROR, err)
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}
	if syncLog == nil {
		de = domain_error.Build(domain_error.DATABASE_RECORD_EMPTY_ERR_CODE, nil)
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}
	middleware.ResponseJson(c, syncLog)
}
