package handlers

import (
	"rhino-common/server/middleware"
	"rhino-data/datamap"

	"github.com/gin-gonic/gin"
)

type SimpleSyncLog struct {
	TableName  string
	ReportTime int64
}

type SyncLogHandler struct {
	autoSyncRepo *datamap.AutoSyncRepo
}

func NewSyncLogHandler(autoSyncRepo *datamap.AutoSyncRepo) *SyncLogHandler {
	return &SyncLogHandler{autoSyncRepo: autoSyncRepo}
}

func (h *SyncLogHandler) GetSyncLogs(c *gin.Context) {
	syncLogs := h.autoSyncRepo.GetSyncLogs()
	if len(syncLogs) == 0 {
		middleware.ResponseJson(c, syncLogs)
		return
	}
	var simpleSyncLogs []*SimpleSyncLog
	visitMap := make(map[string]bool)
	for _, syncLog := range syncLogs {
		_, ok := visitMap[syncLog.TableName]
		if ok {
			continue
		}
		simpleSyncLogs = append(simpleSyncLogs, &SimpleSyncLog{TableName: syncLog.TableName, ReportTime: syncLog.ReportTime})
		visitMap[syncLog.TableName] = true
	}
	middleware.ResponseJson(c, simpleSyncLogs)
}
