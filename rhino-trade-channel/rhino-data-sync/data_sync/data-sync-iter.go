package data_sync

import (
	"bytes"
	"encoding/csv"
	"errors"
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

func (s *DataSync) syncDataForTableAndSyncTypeByIterativeHttpHook(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {

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

	if len(syncParam.ArgsList) == 0 {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("argsList cannot be empty, HttpSyncParam: %s", dataSyncLog.SyncParams)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	var totalRecords []interface{}

	argSize := len(syncParam.ArgsList)
	iterSize := len(syncParam.ArgsList[0])

	log.Printf(">>>>>> argSize:%d, iterSize:%d\n", argSize, iterSize)

	for i := 0; i < iterSize; i++ {

		var args []interface{}
		for j := 0; j < argSize; j++ {
			args = append(args, syncParam.ArgsList[j][i])
		}
		log.Printf(">>>>>> args:%v\n", args)

		url := syncParam.Url
		if strings.Contains(url, "%v") {
			url = fmt.Sprintf(url, args...)
		}

		log.Printf("url:%s\n", url)

		bodyJson := make(map[string]interface{})
		
		if syncParam.BodyJson != nil && len(syncParam.BodyJson) > 0 {
			for k, v := range syncParam.BodyJson {
				str, ok := v.(string)
				if ok && strings.Contains(str, "%v") {
					bodyJson[k] = fmt.Sprintf(str, args)
				} else {
					bodyJson[k] = v
				}
			}
		}

		// 准备请求体
		var bodyReader io.Reader
		if len(bodyJson) > 0 && (syncParam.Method == "POST" || syncParam.Method == "PUT" || syncParam.Method == "PATCH") {
			data, _ := json.Marshal(bodyJson)
			bodyReader = bytes.NewBuffer(data)
		}

		// 创建请求
		req, err := http.NewRequest(syncParam.Method, url, bodyReader)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to create http request, err: %v", err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}

		// 设置请求头
		if syncParam.Headers != nil {
			for key, value := range syncParam.Headers {
				req.Header.Set(key, value)
			}
		}

		// 如果没有设置Content-Type且需要发送JSON body，则设置默认的Content-Type
		if len(bodyJson) > 0 && req.Header.Get("Content-Type") == "" {
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

		log.Printf("parse data:%s\n", data)

		result := map[string]interface{}{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to unmarshal http response, data: %s, err: %v", data, err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}

		var tmpValue interface{} = result
		// 提取本次返回的记录列表
		fields := strings.Split(syncParam.RecordFieldPath, ".")
		for i, field := range fields {
			log.Printf("field:%s\n", field)
			tmpValue = tmpValue.(map[string]interface{})[field]
			if i == len(fields)-1 {
				break
			}
		}
		records, ok := tmpValue.([]interface{})
		if !ok {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to get records in iterative reuest, type asset failed, go type=%T", tmpValue)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, errors.New(dataSyncLog.FailLog), dataSyncLog.FailLog)
			return
		}

		totalRecords = append(totalRecords, records...)
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
