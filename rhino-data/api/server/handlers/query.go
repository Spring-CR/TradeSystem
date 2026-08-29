package handlers

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/server/middleware"
	"rhino-data/api/api_const"
	"rhino-data/api/options"
	"rhino-data/datamap"
	"strings"
	"time"

	"rhino-common/utils/request"

	"github.com/gin-gonic/gin"
)

type DataQueryHandler struct {
	autoSyncRepo *datamap.AutoSyncRepo
}

func NewDataQueryHandler(autoSyncRepo *datamap.AutoSyncRepo) *DataQueryHandler {
	return &DataQueryHandler{autoSyncRepo: autoSyncRepo}
}

func (h *DataQueryHandler) Query(c *gin.Context) {
	begin := time.Now()
	collection, de := request.GetQueryAsString(c, api_const.ParamCollection, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	key, de := request.GetQueryAsString(c, api_const.ParamKey, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	fuzzy, de := request.GetQueryAsBool(c, api_const.ParamFuzzy, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	fields, de := request.GetQueryAsString(c, api_const.ParamFileds, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	maxRecords, de := request.GetQueryAsInt(c, api_const.ParamMaxRecords, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	keyListStr, de := request.GetQueryAsString(c, api_const.ParamKeyList, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}

	if key == "" && keyListStr == "" {
		de := domain_error.Build(domain_error.API_PARAM_NOT_ALLOW_EMPTY_ERR_CODE, nil, api_const.ParamKey+" or "+api_const.ParamKeyList)
		if middleware.ProcessDomainError(de, c) {
			return
		}
	}

	log.Printf("===>Query, collection:%s, key:%s, keyList:%s, fuzzy:%v, fields:%s, maxRecords:%d\n", collection, key, keyListStr, fuzzy, fields, maxRecords)

	if len(keyListStr) > 0 {
		keyList := strings.Split(keyListStr, ",")
		val, de := h.autoSyncRepo.GetByKeyList(collection, keyList)
		if middleware.ProcessDomainError(de, c) {
			return
		}
		n := len(val)
		queryResult := &options.QueryResult{Total: n, Data: val, DisplayLen: n}
		log.Printf("time cost:%v\n", time.Since(begin))
		middleware.ResponseJson(c, queryResult)
		return
	}

	if key == "ALL_MAP_DATA" {
		mapData, de := h.autoSyncRepo.GetMapData(collection)
		if middleware.ProcessDomainError(de, c) {
			return
		}
		middleware.ResponseJson(c, mapData)
		return
	}

	var val []map[string]interface{}
	var ok bool
	if fuzzy {
		val, ok, de = h.autoSyncRepo.GetByFuzzyKey(collection, key)
	} else {
		val, ok, de = h.autoSyncRepo.Get(collection, key)
	}
	if middleware.ProcessDomainError(de, c) {
		return
	}

	if !ok || val == nil {
		val = make([]map[string]interface{}, 0)
	}

	if len(fields) > 0 {
		n := len(val)
		tmp := make([]map[string]interface{}, n)
		fieldList := strings.Split(fields, ",")
		for i := 0; i < n; i++ {
			row := make(map[string]interface{})
			for _, field := range fieldList {
				field = strings.TrimSpace(field)
				if field == "" {
					continue
				}
				v, ok := val[i][field]
				if ok {
					row[field] = v
				}
			}
			tmp[i] = row
		}
		val = tmp
	}

	if maxRecords == 0 {
		maxRecords = 200
	}

	n := len(val)
	displayLen := n

	if len(val) > maxRecords {
		val = val[:maxRecords]
		displayLen = maxRecords
	}

	queryResult := &options.QueryResult{Total: n, Data: val, DisplayLen: displayLen}

	log.Printf("time cost:%v\n", time.Since(begin))

	middleware.ResponseJson(c, queryResult)
}
