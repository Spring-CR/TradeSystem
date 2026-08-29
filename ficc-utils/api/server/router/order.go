package router

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"ficc-utils/api/api_const"
	"ficc-utils/api/options"
	"ficc-utils/common/domain_error"
	"ficc-utils/common/server/middleware"
	"ficc-utils/common/utils/orderutils"
	"ficc-utils/common/utils/request"
	"ficc-utils/common/utils/timeutil"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type CsvColumn struct {
	header    string        // CSV表头名称
	field     string        // order中的字段名
	converter any			// 字段值转换函数, 支持 func(any)string / func(string)string / func(map[string]any)string
	condition func(map[string]any) bool // 列导出的条件，为空时判断field是否存在
}

func setCurrentOrdersQueryRouter(e *gin.Engine, tradeOrdersServiceUrl string) {

	r := e.Group(api_const.RouteCurrentOrders)
	{
		r.GET("", func(c *gin.Context) {
			handleOrdersQuery(c, tradeOrdersServiceUrl)
		})
	}
}

func setHistoryOrdersQueryRouter(e *gin.Engine, hisTradeOrdersServiceUrl string) {

	r := e.Group(api_const.RouteHistoryOrders)
	{
		r.GET("", func(c *gin.Context) {
			handleOrdersQuery(c, hisTradeOrdersServiceUrl)
		})
	}
}

// 测试环境
// 当前订单：http://olts-dev.gf.com.cn/unitrade/test/titans/ficc/order_report/api/v1/report/current/trade_orders
func setCurrentOrdersExportRouter(e *gin.Engine, tradeOrdersServiceUrl string) {

	r := e.Group(api_const.RouteExportCurrentTradeOrders)
	{
		r.POST("", func(c *gin.Context) {
			handleOrdersExport(c, tradeOrdersServiceUrl, false)
		})
	}
}

// 测试环境
// 历史订单：http://olts-dev.gf.com.cn/unitrade/test/titans/ficc/order_report/api/v1/report/history/trade_orders
func setHistoryOrdersExportRouter(e *gin.Engine, hisTradeOrdersServiceUrl string) {

	r := e.Group(api_const.RouteExportHisTradeOrders)
	{
		r.POST("", func(c *gin.Context) {
			handleOrdersExport(c, hisTradeOrdersServiceUrl, true)
		})
	}
}

// 测试环境
// 历史订单：http://olts-dev.gf.com.cn/unitrade/test/titans/ficc/order_report/api/v1/report/history/trade_orders
func handleOrdersExport(c *gin.Context, tradeOrdersServiceUrl string, isHis bool) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_EXPORT_ORDERS_INFO_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}

	result, err := ForwardQueryTradeOrdersService[map[string]any](tradeOrdersServiceUrl, body)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, err, "")
		middleware.ProcessDomainError(de, c)
		return
	}

	if result.Code != 0 {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, nil, result.Message)
		middleware.ProcessDomainError(de, c)
		return
	}

	err = exportCSV(c, result, isHis);
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_EXPORT_ORDERS_INFO_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}
}

func handleOrdersQuery(c *gin.Context, tradeOrdersServiceUrl string) {
	input := map[string]any{}
	fieldConditions := []map[string]any{}

	account, de := request.GetQueryAsInt(c, api_const.ParamAccount, false)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	accountField := map[string]any{
		"field": "account",
		"field_type": 2, //整形
		"value_type": 0, //单值
		"value": account,
	}
	fieldConditions = append(fieldConditions, accountField)

	symbol, de := request.GetQueryAsString(c, api_const.ParamSymbol, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	if symbol != "" {
		symbolField := map[string]any{
		"field": "symbol",
		"field_type": 3, //字符串
		"value_type": 0, //单值
		"value": symbol,
		}
		fieldConditions = append(fieldConditions, symbolField)
	}

	side, de := request.GetQueryAsString(c, api_const.ParamSide, true)
	if middleware.ProcessDomainError(de, c) {
		return
	}
	if side != "" {
		sideField := map[string]any{
		"field": "side",
		"field_type": 3, //字符串
		"value_type": 0, //单值
		"value": side,
		}
		fieldConditions = append(fieldConditions, sideField)
	}

	sortType := 1 // 降序
	//查询历史订单
	if strings.Contains(tradeOrdersServiceUrl, "history") {
		sortType = 0 // 升序
		beginDate, de := request.GetQueryAsString(c, api_const.ParamBeginDate, false)
		if middleware.ProcessDomainError(de, c) {
			return
		}
		endDate, de := request.GetQueryAsString(c, api_const.ParamEndDate, false)
		if middleware.ProcessDomainError(de, c) {
			return
		}
		startTime, err := time.Parse("20060102", beginDate)
		if err != nil {
			de := domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, err, api_const.ParamBeginDate , beginDate)
			middleware.ProcessDomainError(de, c)
			return
		}
		endTime, err := time.Parse("20060102", endDate)
		if err != nil {
			de := domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, err, api_const.ParamEndDate , endDate)
			middleware.ProcessDomainError(de, c)
			return
		}
		if endTime.Sub(startTime) > time.Hour * 24 * 365 {
			de := domain_error.Build(domain_error.API_PARAM_PARSING_ERR_CODE, fmt.Errorf("单次查询跨度不超过1年"), api_const.ParamBeginDate + "," + api_const.ParamEndDate , beginDate + "," + endDate)
			middleware.ProcessDomainError(de, c)
			return
		}
		dateRangeField := map[string]any{
			"field": "f_ord_create_time",
			"field_type": 2, //整形
			"value_type": 1, //区间
			"value": []int64{startTime.In(timeutil.CnTimeLocation).UnixMilli(), endTime.In(timeutil.CnTimeLocation).Add(time.Hour*24).UnixMilli()},
		}
		fieldConditions = append(fieldConditions, dateRangeField)
	}
	input["field_conditions"] = fieldConditions
	input["sort_fields"] = []string{"f_ord_create_time"}
	input["sort_type"] = sortType
	input["limit"] = 1000000000
	input["offset"] = 0
	inputJson, err := json.Marshal(input)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}

	log.Printf("handleOrdersQuery ForwardQueryTradeOrdersService input:%s", inputJson)
	result, err := ForwardQueryTradeOrdersService[*options.Order](tradeOrdersServiceUrl, inputJson)
	if err != nil {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, err)
		middleware.ProcessDomainError(de, c)
		return
	}

	if result.Code != 0 {
		de := domain_error.Build(domain_error.CANNOT_QUERY_ORDERS_INFO_ERR_CODE, nil, result.Message)
		middleware.ProcessDomainError(de, c)
		return
	}

	output := []*options.OrderOut{}
	for _, order := range result.Data {
		if order.F_cum_qty == 0 {
			order.DirtyPriceWithFee = 0
		}
		output = append(output, order.ToOrderOut())
	}

	middleware.ResponseJson(c, output)
}

// 调用订单查询服务（tradeOrdersServiceUrl），并将结果返回
func ForwardQueryTradeOrdersService[T options.TradeOrder](tradeOrdersServiceUrl string, payload []byte) (result *options.GenericQueryResult[T], err error) {
	log.Printf("ForwardQueryTradeOrdersService %s input:%s\n", tradeOrdersServiceUrl, payload)
	resp, err := http.Post(tradeOrdersServiceUrl, "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	log.Printf("tradeOrdersServiceUrl:%s", tradeOrdersServiceUrl)

	body, err := io.ReadAll(resp.Body)
	log.Printf("tradeOrdersServiceUrl body:%s", body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}

	if r, ok := any(result).(*options.GenericQueryResult[*options.Order]); ok {
		for _, order := range r.Data {
			order.OrdStatusText = orderutils.ConvertTrsOrderStatus(order.F_ord_status, order.F_approve_status, order.F_cum_qty)
		}
	}

	bs, _ := json.MarshalIndent(result, "", "    ")
	log.Printf("tradeOrdersServiceUrl result:%s", bs)
	return
}

//func invokeQueryTradeOrdersService(tradeOrdersServiceUrl string, isHis bool, accounts []int, symbol string, openClose string, orderStatusList []string, startDate, endDate, pageNum, pageSize int) (result *options.ExportTradeOrdersResult, err error){
//
//	input := make(map[string]any)
//	input["select_fields"] = []string{"clOrdID","account","counterparty","planCode","ultraContractCode","symbol","symbolName","side","openClose","handlInst", "settlType","ytm","price","dirtyPrice","quantity","f_cum_qty","currency","remark","f_ord_create_time","f_ord_status","f_approve_status","f_ord_rej_reason","locked", "f_trade_date"}
//
//	fieldConditions := []map[string]any{}
//	if len(accounts) > 0 {
//		accountField := map[string]any{
//			"field": "account",
//			"field_type": 2, //整形
//			"value_type": 2, //集合
//			"value": accounts,
//		}
//		fieldConditions = append(fieldConditions, accountField)
//	}
//	if len(symbol) >0 {
//		symbolField := map[string]any{
//			"field": "symbol",
//			"field_type": 3, //字符串
//			"value_type": 0, //单值
//			"value": symbol,
//		}
//		fieldConditions = append(fieldConditions, symbolField)
//	}
//	if openClose != "" {
//		openCloseField := map[string]any{
//			"field": "openClose",
//			"field_type": 3, //字符串
//			"value_type": 0, //单值
//			"value": openClose,
//		}
//		fieldConditions = append(fieldConditions, openCloseField)
//	}
//
//	if len(orderStatusList) > 0 {
//		orderStatusField := map[string]any{
//			"field": "f_ord_status",
//			"field_type": 3, //字符串
//			"value_type": 2, //集合
//			"value": orderStatusList,
//		}
//		fieldConditions = append(fieldConditions, orderStatusField)
//	}
//
//	if isHis && startDate > 0 && endDate > 0 {
//		dateField := map[string]any{
//			"field": "f_trade_date",
//			"field_type": 1, //时间
//			"value_type": 1, //范围
//			"value": []int{startDate, endDate},
//
//		}
//		fieldConditions = append(fieldConditions, dateField)
//	}
//
//	input["field_conditions"] = fieldConditions
//	input["sort_fields"] = []string{"f_ord_create_time"}
//	input["sort_type"] = 1 //降序
//	input["limit"] = pageSize
//	input["offset"] = (pageNum-1) * pageSize
//
//	inputJson, err := json.Marshal(input)
//	if err != nil {
//		return
//	}
//	log.Printf("tradeOrdersServiceUrl input:%s", inputJson)
//	resp, err := http.Post(tradeOrdersServiceUrl, "application/json", bytes.NewReader(inputJson))
//	if err != nil {
//		return
//	}
//	defer resp.Body.Close()
//	log.Printf("tradeOrdersServiceUrl:%s", tradeOrdersServiceUrl)
//
//	body, err := io.ReadAll(resp.Body)
//	log.Printf("tradeOrdersServiceUrl body:%s", body)
//
//	if err != nil {
//		return
//	}
//
//	err = json.Unmarshal(body, &result)
//	if err != nil {
//		return
//	}
//	log.Printf("tradeOrdersServiceUrl result:%v", result)
//	return
//}

var exportColumns = []CsvColumn{
	{header: "日期", field: "f_ord_create_time", converter: formatDate},
	{header: "订单号", field: "clOrdID", converter: formatString},
	{header: "账户名称", field: "counterparty", converter: formatString},
	{header: "大合约编号", field: "ultraContractCode", converter: formatString},
	{header: "标的代码", field: "symbol", converter: formatString},
	{header: "标的名称", field: "symbolName", converter: formatString},
	{header: "交易方向", field: "side", converter: orderutils.ConvertSide},
	//{header: "交易方向", field: "openClose", converter: convertOpenClose},
	{header: "交易效率", field: "handlInst", converter: orderutils.ConvertHandlInst},
	{header: "清算速度", field: "settlType", converter: formatString},
	{header: "意向到期收益(%)", field: "ytm", converter: formatFloat64(-1) },
	{header: "意向净价", field: "price", converter: formatFloat64(-1) },
	{header: "意向全价", field: "dirtyPrice", converter: formatFloat64(-1) },
	{header: "意向全价(含费)", field: "dirtyPriceWithFee", converter: formatFloat64(-1) },
	{header: "券面总额(万元)", field: "quantity", converter: formatFloat64(-1, 10000) }, // 除以10000
	{header: "成交进度", field:  "", converter: formatExecPercent, condition: keysInOrder("quantity", "f_cum_qty")},
	{header: "成交面额(万元)", field:  "f_cum_qty", converter: formatFloat64(-1, 10000) }, // 除以10000
	{header: "成交净价(不含费)", field:  "avgPrice", converter: formatExecPrice("avgPrice"), condition: keysInOrder("f_cum_qty", "avgPrice")},
	{header: "成交全价(不含费)", field:  "avgDirtyPrice",converter: formatExecPrice("avgDirtyPrice"), condition: keysInOrder("f_cum_qty", "avgDirtyPrice")},
	{header: "成交全价(含费)", field: "avgDirtyPriceWithFee", converter: formatExecPrice("avgDirtyPriceWithFee"), condition: keysInOrder("f_cum_qty", "avgDirtyPriceWithFee")},
	{header: "佣金费率", field: "commissionRate", converter: formatSpreadAndFee("commissionRate"), condition: keysInOrder("f_cum_qty", "commissionRate")},
	{header: "利差", field: "spread", converter: formatSpreadAndFee("spread"), condition: keysInOrder("f_cum_qty", "spread")},
	{header: "币种", field: "currency", converter: formatString},
	{header: "订单来源", field: "ordSource", converter: orderutils.ConvertOrdSource},
	{header: "交易申请时间", field: "f_ord_create_time", converter: formatDateTime },
	{header: "状态", field:  "f_ord_status", converter: formatOrderStatus, condition: keysInOrder("f_ord_status", "f_cum_qty", "f_approve_status")},
}


// "日期",         "账户名称",      "大合约编号",         "标的代码", "标的名称",   "交易方向",   "交易效率",  "清算速度",  "意向到期收益(%)", "意向净价", "意向全价",    "券面总额(万元)", "成交进度", "成交面额(万元)", "成交净价(不含费)", "成交全价(不含费)", "币种",     "交易申请时间",        "状态"
//"f_trade_date", "counterparty", "ultraContractCode",  "symbol",  "symbolName", "side",   , "handlInst", "settlType", "ytm",            "price",   "dirtyPrice", "quantity",                  "f_cum_qty",                                         "currency", "f_ord_create_time", "f_ord_status"
func exportCSV(c *gin.Context, data *options.GenericQueryResult[map[string]any], isHis bool) error {
	c.Header("Content-Disposition", "attachment;filename=export_orders.csv")
	c.Header("Content-Type", "text/csv;charset=utf-8")

	byteBuf := bytes.NewBuffer(nil)

	writer := csv.NewWriter(byteBuf)

	header := make([]string, 0, len(exportColumns))
	records := make([][]string, 0, len(data.Data))
	headerSeted := false

	//第一列显示日期
	exportColumns[0].condition = func(m map[string]any) bool { return true }
	for _, order := range data.Data {
		record := make([]string, 0, len(exportColumns))
		for _, col := range exportColumns {
			if col.condition == nil {
				if col.field != "" && !keysIn(order, col.field) {
					continue
				}
			} else if !col.condition(order) {
				continue
			}
			if !headerSeted {
				header = append(header, col.header)
			}

			val := getColVal(order, col)
			record = append(record, val)
		}
		records = append(records, record)
		headerSeted = true
	}

	if len(header) > 0 {
		if err := writer.Write(header); err != nil {
			return fmt.Errorf("write csv header failed: %w", err)
		}
	}
	if len(records) > 0 {
		if err := writer.WriteAll(records); err != nil {
			return fmt.Errorf("write csv records failed: %w", err)
		}
	}

	writer.Flush()

	// GBK 编码器
    encoder := simplifiedchinese.GBK.NewEncoder()
	gbkBytes, _, err := transform.Bytes(encoder, byteBuf.Bytes())
    if err != nil {
        return err
    }
	
	c.Writer.Write(gbkBytes)

	return nil
}

func keysIn(m map[string]any, keys ...string) bool {
	if len(keys) == 0 {
		return false
	}
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func keysInOrder(keys ...string) func(m map[string]any) bool {
	return func(m map[string]any) bool {
		return keysIn(m, keys...)
	}
}

func formatString(v any) string {
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

func formatFloat64(prec int, divisor ...float64) func(v any) string {
	return func(v any) string {
		val, ok := v.(float64)
		if !ok {
			return ""
		}
		if len(divisor) > 0 && divisor[0] > 1e-9 {
			val /= divisor[0]
		}
		return strconv.FormatFloat(val, 'f', prec, 64)
	}
}

func formatDateTime(v any) string {
	timeStamp, ok := v.(float64)
	if !ok {
		return ""
	}
	return time.Unix(int64(timeStamp/1000), 0).Format(time.DateTime)
}


func formatDate(v any) string {
	timeStamp, ok := v.(float64)
	if !ok {
		return ""
	}
	return time.Unix(int64(timeStamp/1000), 0).Format(time.DateOnly)
}

func formatExecPercent(order map[string]any) string {
	quantity := getFloat64(order, "quantity")
	fCumQty := getFloat64(order, "f_cum_qty")
	if quantity <= 1e-9 { // 避免除零
		return "0.00%"
	}
	rate := fCumQty / quantity * 100
	return fmt.Sprintf("%.2f%%", rate)
}

func formatExecPrice(key string) func(order map[string]any) string {
	return func (order map[string]any) string {
		cumQty := getFloat64(order, "f_cum_qty")
		if cumQty > 1e-9 {
			price := getFloat64(order, key)
			return strconv.FormatFloat(price, 'f', -1, 64)
		}
		return ""
	}
}

func formatOrderStatus(order map[string]any) string {
	ordStatus := formatString(order["f_ord_status"])
	approveStatus := int(getFloat64(order, "f_approve_status"))
	fCumQty := int64(getFloat64(order, "f_cum_qty"))
	return orderutils.ConvertTrsOrderStatus(ordStatus, approveStatus, fCumQty)
}

func getFloat64(m map[string]any, key string) float64 {
	val, ok := m[key]
	if !ok {
		return 0
	}
	f, ok := val.(float64)
	if !ok {
		return 0
	}
	return f
}

func formatSpreadAndFee(key string) func(order map[string]any) string {
	return func (order map[string]any) string {
		cumQty := getFloat64(order, "f_cum_qty")
		if cumQty > 1e-9 {
			val, ok := order[key].(float64)
			if ok {
				return fmt.Sprintf("%.4f%%", val*100)
			}
		}
		return ""
	}
}

func getColVal(order map[string]any, col CsvColumn) string {
	val := ""
	if col.field == "" || keysIn(order, col.field) {
		switch conv := col.converter.(type) {
		case func(any) string: val = conv(order[col.field])
		case func(string) string: val = conv(formatString(order[col.field]))
		case func(map[string]any) string: val = conv(order)
		}
	}
	return val
}