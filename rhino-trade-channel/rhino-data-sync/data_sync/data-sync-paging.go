package data_sync

import (
	"bytes"
	"encoding/csv"
	"errors"
	"farm/util/request"
	"fmt"
	"io"
	"log"
	"net/http"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"strings"
	"time"
)

func (s *DataSync) syncDataForTableAndSyncTypeByPagingHttpHook(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {

	syncParam := &HttpSyncParam{}

	var err error
	dataSyncLog.SyncParams, _, err = refineParam(s.dataSyncCfg.GetAppDB(), dataSyncLog.SyncParams)
	if err != nil {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to refineParam, data: %s, err: %v", dataSyncLog.SyncParams, err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	err = json.Unmarshal([]byte(dataSyncLog.SyncParams), syncParam)
	if err != nil {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to unmarshal HttpSyncParam, data: %s, err: %v", dataSyncLog.SyncParams, err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	if syncParam.Url == "" {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("url cannot be empty, HttpSyncParam: %s", dataSyncLog.SyncParams)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	if syncParam.Method == "" {
		syncParam.Method = "GET" // 默认方法
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var totalRecords []interface{}

	for i := syncParam.PagingFrom; i < 10000; i++ {

		// 计算分页参数
		pageNum := i
		pageSize := syncParam.PagingSize

		// 准备请求体
		var bodyReader io.Reader
		if len(syncParam.BodyJson) > 0 && (syncParam.Method == "POST" || syncParam.Method == "PUT" || syncParam.Method == "PATCH") {

			if syncParam.PageArgType == 0 {
				syncParam.BodyJson[syncParam.PageNumField] = pageNum
				syncParam.BodyJson[syncParam.PageSizeField] = pageSize
			}

			data, _ := json.Marshal(syncParam.BodyJson)
			bodyReader = bytes.NewBuffer(data)
		}

		// 创建请求
		req, err := http.NewRequest(syncParam.Method, syncParam.Url, bodyReader)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to create http request, err: %v", err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}
		if syncParam.PageArgType == 1 {
			request.AddQueryParam(req, syncParam.PageNumField, pageNum)
			request.AddQueryParam(req, syncParam.PageSizeField, pageSize)
		}

		// 设置请求头
		if syncParam.Headers != nil {
			for key, value := range syncParam.Headers {
				req.Header.Set(key, value)
			}
		}

		// 如果没有设置Content-Type且需要发送JSON body，则设置默认的Content-Type
		if len(syncParam.BodyJson) > 0 && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		// 执行请求
		resp, err := client.Do(req)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to send http request, err: %v", err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)

		result := map[string]interface{}{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to unmarshal http response, data: %s, err: %v", data, err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}

		// 提取本次返回的记录数
		fields := strings.Split(syncParam.TotalFieldPath, ".")
		var tmpValue interface{} = result
		for i, field := range fields {
			tmpValue = tmpValue.(map[string]interface{})[field]
			if i == len(fields)-1 {
				break
			}
		}
		var total int
		switch v := tmpValue.(type) {
		case int:
			total = v
		case int8:
			total = int(v)
		case int16:
			total = int(v)
		case int32:
			total = int(v)
		case int64:
			total = int(v)
		case uint:
			total = int(v)
		case uint8:
			total = int(v)
		case uint16:
			total = int(v)
		case uint32:
			total = int(v)
		case uint64:
			total = int(v)
		case float32:
			total = int(v)
		case float64:
			total = int(v)
		default:
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("unsuuport value type for total value: %v", v)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, errors.New(dataSyncLog.FailLog), dataSyncLog.FailLog)
			return
		}

		// 提取本次返回的记录列表
		fields = strings.Split(syncParam.RecordFieldPath, ".")
		tmpValue = result
		for i, field := range fields {
			tmpValue = tmpValue.(map[string]interface{})[field]
			if i == len(fields)-1 {
				break
			}
		}
		records, ok := tmpValue.([]interface{})
		if !ok {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to get records in paging reuest, type asset failed, go type=%T", tmpValue)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, errors.New(dataSyncLog.FailLog), dataSyncLog.FailLog)
			return
		}

		totalRecords = append(totalRecords, records...)

		if len(totalRecords) >= total {
			break
		}
	}

	// 生成csv
	buf := bytes.NewBufferString(syncParam.CsvHeader + "\n")
	headers := strings.Split(syncParam.CsvHeader, ",")
	csvWriter := csv.NewWriter(buf)
	for _, _record := range totalRecords {
		record, ok := _record.(map[string]interface{})
		if !ok {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to get records in paging reuest, type asset failed, go type=%T", _record)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, errors.New(dataSyncLog.FailLog), dataSyncLog.FailLog)
			return
		}
		var csvRow []string
		for _, header := range headers {
			val, ok := record[header]
			if ok && val != nil {
				csvRow = append(csvRow, fmt.Sprintf("%v", val))
			} else {
				csvRow = append(csvRow, "")
			}
		}
		err = csvWriter.Write(csvRow)
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "fail to write csv record for application config")
		}
	}
	csvWriter.Flush()


	// 转化csv
	body, err := s.refineCsvContent(tableConfig, buf.Bytes())
	if err != nil {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while refine csv, err: %v", err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	// 将csv导入数据库
	err = dbutil.CreateTableFromCsv(s.dataSyncCfg.GetAppDB(), body, tableConfig, func() bool {
		dataSyncLog, _ := admin_store.GetDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog.ID)
		if dataSyncLog != nil && dataSyncLog.SyncPhase == int(enum.DataSyncLogPhase_Complete) {
			return false
		}
		return true
	})
	if err != nil {
		log.Println("csv 导入数据库失败:", err)
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while import csv file to db, err: %v", err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Complete)
	dataSyncLog.CompleteSyncTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	dataSyncLog.FailLog = ""
	admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
	log.Printf("======> data sync log for table %v is complete", dataSyncLog.TableName)
}
