package data_sync

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-core/adapter_registry"
	"rhino-core/domain_cfg"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
	"rhino-data-sync/data_sync_adapter"
	_ "rhino-plugins/data_sync_plugin"
	"sort"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	pahaseIgnoreMap = map[int]bool{
		int(enum.DataSyncLogPhase_Complete): true,
		int(enum.DataSyncLogPhase_Cancel):   true,
	}
	timeoutMillions = int64(30 * 60 * 1000)
	json            = jsoniter.ConfigCompatibleWithStandardLibrary
	csvQryClient    = &http.Client{
		Timeout: 5 * time.Second, // 全局超时，覆盖连接、传输等所有阶段
	}
	csvDownloadClient = &http.Client{
		Timeout: 5 * time.Minute, // 全局超时，覆盖连接、传输等所有阶段
	}
)

type DataSync struct {
	dataSyncCfg     *domain_cfg.DataSyncConfig
	dataSyncAdapter data_sync_adapter.DataSyncAdapter
}

func NewDataSync(dataSyncCfg *domain_cfg.DataSyncConfig) *DataSync {
	checkMrkCloseTime(dataSyncCfg)
	var dataSyncAdapter data_sync_adapter.DataSyncAdapter
	adapterPath := dataSyncCfg.GetDataSyncAdapterPath()
	if adapterPath != "" {
		// 从注册表获取适配器的构造函数（目前，对于apiAdpater，是无参数函数，有参函数需要根据特殊情况来处理了）
		_dataSyncAdapter, de, err := adapter_registry.CallAdapterFunction(adapterPath, dataSyncCfg)
		if err != nil {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct DataSyncAdapter")
		}
		if de != domain_error.NilDomainError {
			domain_error.ProcessSevereError(true, 5, nil, err, "fail to construct DataSyncAdapter")
		}
		// 获取apiAdapter
		dataSyncAdapter = _dataSyncAdapter.(data_sync_adapter.DataSyncAdapter)
		log.Println("finish get DataSyncAdapter")
	}

	return &DataSync{dataSyncCfg: dataSyncCfg, dataSyncAdapter: dataSyncAdapter}
}

func checkMrkCloseTime(dataSyncCfg *domain_cfg.DataSyncConfig) {
	mrkCloseTime, mrkCloseTimeZone := dataSyncCfg.GetMrkCloseTime()
	location := timeutil.GetTimeZone(mrkCloseTimeZone)
	dateStr := time.Now().In(location).Format(time.DateOnly)
	datetimeStr := dateStr + " " + mrkCloseTime
	_, err := time.ParseInLocation(time.DateTime, datetimeStr, location)
	if err != nil {
		domain_error.ProcessSevereError(true, 5, nil, err, "fail to parse mrkCloseTime")
	}
}

// 启动数据清理器
func (s *DataSync) Start() {
	go func() {
		for {
			s.syncData()
			time.Sleep(60 * time.Second)
		}
	}()
}

func (s *DataSync) syncData() {
	log.Printf("start sync data")
	for _, tableConfig := range s.dataSyncCfg.GetTableConfigs() {
		s.syncDataForTable(tableConfig)
	}
}

func (s *DataSync) syncDataForTable(tableConfig *dbutil.TableConfig) {

	log.Printf("===> star to sync data for table %s\n", tableConfig.TableName)

	mrkCloseTime, mrkCloseTimeZone := s.dataSyncCfg.GetMrkCloseTime()
	location := timeutil.GetTimeZone(mrkCloseTimeZone)
	dateStr := time.Now().In(location).Format(time.DateOnly)
	datetimeStr := dateStr + " " + mrkCloseTime

	closeTime, _ := time.ParseInLocation(time.DateTime, datetimeStr, location)
	currTime := time.Now().In(location)
	if currTime.After(closeTime) {
		dateStr = currTime.Add(24 * time.Hour).Format(time.DateOnly)
	}

	// 数据库中的日期没有中划线
	dateStr = strings.ReplaceAll(dateStr, "-", "")
	systemCode, businessCode := s.dataSyncCfg.GetSystemAndBusinessCode()

	dataSyncLogs, err := admin_store.FindDataSyncLogsBySystemCodeAndBusinessCodeAndSyncDateAndTableName(s.dataSyncCfg.GetCentralDB(), systemCode, businessCode, dateStr, tableConfig.TableName)

	log.Printf("dataSyncLogs len=%d for %s\n", len(dataSyncLogs), tableConfig.TableName)

	if err != nil && dbutil.IsDbRecordEmptyError(err) {
		err = nil
	}

	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to find data sync logs")
		return
	}

	// 如果没有数据，则直接返回
	if len(dataSyncLogs) == 0 {
		return
	}

	// 按记录日期从大到小排序
	sort.Slice(dataSyncLogs, func(i, j int) bool {
		return dataSyncLogs[i].ReportTime > dataSyncLogs[j].ReportTime
	})

	var dataSyncLog *schema.DataSyncLog

	for _, v := range dataSyncLogs {
		if v.SyncPhase == int(enum.DataSyncLogPhase_New) {
			dataSyncLog = v
			break
		} else if pahaseIgnoreMap[v.SyncPhase] { // if Complete, Cancel
			return
		} else if v.SyncPhase == int(enum.DataSyncLogPhase_Fail) {
			// 如果超时了，要改成取消状态，并且退出
			if v.FirstSyncTime+timeoutMillions < timeutil.ConvertTimeToMilliseconds(time.Now()) {
				v.SyncPhase = int(enum.DataSyncLogPhase_Cancel)
				v.FailLog = "fail to sync data, timeout"
				admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), v)
				return
			}
			// 如果没有超时，则继续处理
			dataSyncLog = v
			break
		} else if v.SyncPhase == int(enum.DataSyncLogPhase_Processing) {
			// 处理时间太长时，要重试
			if v.FirstSyncTime+timeoutMillions/2 < timeutil.ConvertTimeToMilliseconds(time.Now()) {
				log.Printf("data sync log for table %v is processing too long, retry", v.TableName)
				dataSyncLog = v
			} else { // 还在合理的处理时间范围内，直接退出
				return
			}
		}
	}

	// 无需处理
	if dataSyncLog == nil {
		log.Printf("no actived dataSyncLog for table %s\n", dataSyncLog.TableName)
		return
	}

	// 原始的同步状态，用于判断是否需要更新
	origSyncPhase := dataSyncLog.SyncPhase
	origExecCount := dataSyncLog.ExecCount

	// select for update加锁
	tx, de := dbutil.BeginTx(s.dataSyncCfg.GetCentralDB())
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, nil, "fail to begin tx")
		dbutil.RollbackTx(tx)
		return
	}

	// 通过select for update语句加行锁
	dataSyncLog, err = admin_store.SelectDataSyncLogByIdForUpdate(tx, dataSyncLog.ID)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to select data sync log for update")
		dbutil.RollbackTx(tx)
		return
	}

	// 当判断到dataSyncLog已在另一个进程/线程处理时，就直接退出
	// 1、 当前状态是完成、取消，则直接退出
	// 2、 当前状态和原来的状态不一致时，直接退出
	// 3、 当前状态是处理中，并且执行次数和原来的不一致时，直接退出
	if pahaseIgnoreMap[dataSyncLog.SyncPhase] || dataSyncLog.SyncPhase != origSyncPhase || (dataSyncLog.SyncPhase == int(enum.DataSyncLogPhase_Processing) && dataSyncLog.ExecCount != origExecCount) {
		log.Println(">>> 判断到dataSyncLog已在另一个进程/线程处理时，就直接退出")
		log.Printf("   1. pahaseIgnoreMap:%v, dataSyncLog.SyncPhase:%v\n", dataSyncLog, dataSyncLog.SyncPhase)
		log.Printf("   2. dataSyncLog.SyncPhase:%v, origSyncPhase:%v\n", dataSyncLog.SyncPhase, origSyncPhase)
		log.Printf("   3. dataSyncLog.ExecCount:%v, origExecCount:%v\n", dataSyncLog.ExecCount, origExecCount)
		dbutil.RollbackTx(tx)
		return
	}

	// 更新状态为处理中
	dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Processing)
	dataSyncLog.CurrentSyncTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	dataSyncLog.ExecCount += 1
	if dataSyncLog.FirstSyncTime <= 0 {
		dataSyncLog.FirstSyncTime = dataSyncLog.CurrentSyncTime
	}

	err = admin_store.UpdateDataSyncLogById(tx, dataSyncLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to update data sync log")
		dbutil.RollbackTx(tx)
		return
	}

	// 提交事务
	de = dbutil.CommitTx(tx)
	if de != nil {
		domain_error.ProcessSevereError(false, 0, de, nil, "fail to commit tx")
		dbutil.RollbackTx(tx)
		return
	}

	s.syncDataForTableAndSyncType(tableConfig, dataSyncLog)
}

func (s *DataSync) syncDataForTableAndSyncType(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {
	switch dataSyncLog.SyncType {
	case int(enum.SyncType_DSP):
		s.syncDataForTableAndSyncTypeByBigDataDsp(tableConfig, dataSyncLog)
	case int(enum.SyncType_HTTP_HOOK):
		s.syncDataForTableAndSyncTypeByHttpHook(tableConfig, dataSyncLog)
	case int(enum.SyncType_PAGING_HTTP_HOOK):
		s.syncDataForTableAndSyncTypeByPagingHttpHook(tableConfig, dataSyncLog)
	case int(enum.SyncType_ITERATIVE_HTTP_HOOK):
		s.syncDataForTableAndSyncTypeByIterativeHttpHook(tableConfig, dataSyncLog)
	case int(enum.SyncType_CSV):
		s.syncDataForTableAndSyncTypeByCsv(tableConfig, dataSyncLog)
	}
}

func (s *DataSync) syncDataForTableAndSyncTypeByBigDataDsp(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {
	dspOption := &DspOption{}
	err := json.Unmarshal([]byte(dataSyncLog.SyncParams), dspOption)
	if err != nil {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to unmarshal sync params, data: %s, err: %v", dataSyncLog.SyncParams, err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	// 循环读取csv文件的生成状态，直到生成成功
	err = s.waitForCsvFileReady(dataSyncLog, dspOption)
	if err != nil {
		return
	}

	// 读取csv文件内容
	resp, err := postDataForDsp(csvDownloadClient, "/api/open/ds/dataFile/download", dataSyncLog, dspOption)
	if err != nil {
		log.Println("请求失败:", err)
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while download csv file, err: %v", err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			log.Println("close response body for download csv file")
			resp.Body.Close() // 确保关闭响应体
		}
	}()

	// 读取文件流
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Println("请求失败:", resp.StatusCode, string(body))
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while read csv file stream, err: %v", err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	// 转化csv
	body, err = s.refineCsvContent(tableConfig, body)
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

	// 落地文件，便于调试
	if len(body) > 0 {
		os.WriteFile(fmt.Sprintf("/tmp/%s.csv", dataSyncLog.TableName), body, 0644)
	}

	dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Complete)
	dataSyncLog.CompleteSyncTime = timeutil.ConvertTimeToMilliseconds(time.Now())
	dataSyncLog.FailLog = ""
	admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
	log.Printf("======> data sync log for table %v is complete", dataSyncLog.TableName)
}

func (s *DataSync) waitForCsvFileReady(dataSyncLog *schema.DataSyncLog, dspOption *DspOption) (err error) {
	// 循环读取csv文件的生成状态，直到生成成功
	for {
		var resp *http.Response
		resp, err = postDataForDsp(csvQryClient, "/api/open/ds/dataFile/status", dataSyncLog, dspOption)
		if err != nil {
			log.Println("请求失败:", err)
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while query csv file status, err: %v", err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}
		defer func() {
			if resp != nil && resp.Body != nil {
				log.Println("close response body for waitForCsvFileReady")
				resp.Body.Close() // 确保关闭响应体
			}
		}()

		// 读取响应
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 3 && body[0] == 0xEF && body[1] == 0xBB && body[2] == 0xBF {
			// 确认是BOM，已通过Read操作跳过，文件指针已后移3字节
			fmt.Println("检测到BOM并已跳过")
			body = body[3:]
		}
		if resp.StatusCode != http.StatusOK {
			log.Println("请求失败:", resp.StatusCode, string(body))
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while read http response, status code: %v, body: %s", resp.StatusCode, body)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return errors.New(dataSyncLog.FailLog)
		}

		csvQryResp := &CsvQryResp{}
		err = json.Unmarshal(body, csvQryResp)
		if err != nil {
			log.Println("请求失败:", err)
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while unmarshal csv file query response, body: %s, err: %v", body, err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return errors.New(dataSyncLog.FailLog)
		}

		if csvQryResp.Data == "SUCCESS" {
			log.Printf("csv file for %v is ready\n", dataSyncLog.TableName)
			return nil
		} else {
			log.Printf("csv file for %v is not ready, will try to check again after 5 seconds\n", dataSyncLog.TableName)
			time.Sleep(5 * time.Second)
		}
	}
}

func postDataForDsp(client *http.Client, apiPath string, dataSyncLog *schema.DataSyncLog, dspOption *DspOption) (resp *http.Response, err error) {
	sign, timestamp := dspOption.GenerateSecret()

	// 查询csv文件的就绪状态
	option := map[string]interface{}{
		"source":    dspOption.AppKey,
		"timeStamp": strconv.Itoa(int(timestamp)),
		"password":  sign,
		"queryId":   dspOption.QueryId,
	}

	jsonData, _ := json.Marshal(option)

	url := fmt.Sprintf("%s/%s", dspOption.BaseUrl, apiPath)

	log.Printf("postDataForDsp, url: %s, postData: %s\n", url, jsonData)

	resp, err = client.Post(
		url,
		"application/json",
		bytes.NewBuffer(jsonData),
	)

	if err == nil && (resp == nil || resp.Body == nil) {
		return resp, fmt.Errorf("http resp is empty, resp == nil:%v, resp.Body == nil:%v", resp == nil, resp.Body == nil)
	}

	return
}

type HttpSyncParam struct {
	Source          string                 `json:"source"`
	Url             string                 `json:"url"`
	Headers         map[string]string      `json:"headers"`
	Method          string                 `json:"method"`
	BodyJson        map[string]interface{} `json:"bodyJson"`
	ResultType      string                 `json:"resultType"`
	UsePage         bool                   `json:"usePage"`
	PageArgType     int                    `json:"pageArgType"` // 分页参数的类型，0-在body请求体；1-在query参数列表
	PageNumField    string                 `json:"pageNumField"`
	PageSizeField   string                 `json:"pageSizeField"`
	PagingFrom      int                    `json:"pagingFrom"`
	PagingSize      int                    `json:"pagingSize"`
	TotalFieldPath  string                 `json:"totalFieldPath"`
	RecordFieldPath string                 `json:"recordFieldPath"`
	CsvHeader       string                 `json:"csvHeader"` // json时才需要，按需提取字段
	UseIteration    bool                   `json:"useIteration"`
	ArgsList        [][]interface{}        `json:"argsList"`
}

func (s *DataSync) syncDataForTableAndSyncTypeByHttpHook(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {

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

	// 准备请求体
	var bodyReader io.Reader
	if len(syncParam.BodyJson) > 0 && (syncParam.Method == "POST" || syncParam.Method == "PUT" || syncParam.Method == "PATCH") {
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

	switch syncParam.ResultType {
	case "csv":
		// 读取文件流
		body, _ := io.ReadAll(resp.Body)
		if len(body) > 3 && body[0] == 0xEF && body[1] == 0xBB && body[2] == 0xBF {
			// 确认是BOM，已通过Read操作跳过，文件指针已后移3字节
			fmt.Println("检测到BOM并已跳过")
			body = body[3:]
		}

		// 转化csv
		body, err = s.refineCsvContent(tableConfig, body)
		if err != nil {
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while refine csv, err: %v", err)
			admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
			domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
			return
		}

		// 落地文件，便于调试
		if len(body) > 0 {
			os.WriteFile(fmt.Sprintf("/tmp/%s.csv", dataSyncLog.TableName), body, 0644)
		}

		if resp.StatusCode != http.StatusOK {
			log.Println("请求失败:", resp.StatusCode, string(body))
			dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
			dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while read csv file stream, err: %v", err)
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
	case "json":

	}
}

func (s *DataSync) syncDataForTableAndSyncTypeByCsv(tableConfig *dbutil.TableConfig, dataSyncLog *schema.DataSyncLog) {
	log.Printf("=======> syncDataForTableAndSyncTypeByCsv: %s\n", tableConfig.TableAlias)
	// 读取文件流
	body := []byte(dataSyncLog.SyncParams)
	if len(body) > 3 && body[0] == 0xEF && body[1] == 0xBB && body[2] == 0xBF {
		// 确认是BOM，已通过Read操作跳过，文件指针已后移3字节
		fmt.Println("检测到BOM并已跳过")
		body = body[3:]
	}

	// 转化csv
	body, err := s.refineCsvContent(tableConfig, body)
	if err != nil {
		dataSyncLog.SyncPhase = int(enum.DataSyncLogPhase_Fail)
		dataSyncLog.FailLog = fmt.Sprintf("fail to sync data while refine csv, err: %v", err)
		admin_store.UpdateDataSyncLogById(s.dataSyncCfg.GetCentralDB(), dataSyncLog)
		domain_error.ProcessSevereError(false, 0, nil, err, dataSyncLog.FailLog)
		return
	}

	// 落地文件，便于调试
	if len(body) > 0 {
		os.WriteFile(fmt.Sprintf("/tmp/%s.csv", dataSyncLog.TableName), body, 0644)
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

func(s*DataSync)refineCsvContent(tableConfig*dbutil.TableConfig, body[]byte)([]byte, error){
	if s.dataSyncAdapter == nil {
		return body, nil
	}
	return s.dataSyncAdapter.RefineCsvContent(tableConfig, body)
}