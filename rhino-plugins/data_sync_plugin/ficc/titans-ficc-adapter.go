package ficc

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/dbutil"
	"rhino-common/utils/timeutil"
	"rhino-common/utils/tradedate"
	"rhino-core/domain_cfg"
	"rhino-plugins/data_sync_plugin/utils"
	"time"
)

type TitansFiccDataSyncAdapter struct {
	dataSyncConfig *domain_cfg.DataSyncConfig
}

func NewTitansFiccDataSyncAdapter(dataSyncConfig *domain_cfg.DataSyncConfig) (adapter *TitansFiccDataSyncAdapter, de *domain_error.Error) {
	log.Printf("construct TitansFiccDataSyncAdapter...")
	adapter = &TitansFiccDataSyncAdapter{dataSyncConfig: dataSyncConfig}
	return
}

func (a *TitansFiccDataSyncAdapter) RefineCsvContent(tableConfig *dbutil.TableConfig, rawCsv []byte) (newCsv []byte, err error) {

	if tableConfig.TableAlias != "PositionBase" {
		return rawCsv, nil
	}

	// T日
	var currTradeDate string
	// T+1日
	var nextTradeDate string

	currTradeDate, nextTradeDate, err = a.getCurrentAndNextTradeDates()
	if err != nil {
		return rawCsv, err
	}

	log.Printf("compute trade dates for PositionBase, currTradeDate:%s, nextTradeDate:%s\n", currTradeDate, nextTradeDate)

	// 解析csv数据内容
	contractPositionTableCinfig, err1 := a.getContractPositionTableConfig()
	if err1 != nil {
		return rawCsv, err1
	}
	lines, err1 := utils.ConvertCsvToMap(rawCsv, contractPositionTableCinfig)
	if err1 != nil {
		return rawCsv, err1
	}

	parValMap, err1 := GetParValueMap(a.dataSyncConfig.GetAppDB())
	if err1 != nil {
		return rawCsv, err1
	}

	log.Printf("total lines of contractPositionTable:%d\n", len(lines))

	// 按 keyCtptyId、keyPlanId、keyInstrumentId、longShort 分组
	lineAggregateMap := make(map[string][]map[string]interface{})
	// 聚合
	for _, line := range lines {
		key, ok := getKeyOfLine(line)
		if !ok {
			continue
		}
		lineAggregateMap[key] = append(lineAggregateMap[key], line)
	}
	log.Printf("lineAggregateMap size=%d\n", len(lineAggregateMap))

	// 计算聚会记录
	var basePositions []map[string]interface{}
	for _, contractPositions := range lineAggregateMap {
		t0BasePosition, t1BasePosition, err := computeBasePositionRecord(currTradeDate, nextTradeDate, contractPositions, parValMap)
		if err != nil {
			return rawCsv, err
		}
		if t0BasePosition != nil {
			basePositions = append(basePositions, t0BasePosition)
		}
		if t1BasePosition != nil {
			basePositions = append(basePositions, t1BasePosition)
		}
	}

	return utils.ConvertMapToCsv(basePositions, tableConfig)
}

func (a *TitansFiccDataSyncAdapter) getContractPositionTableConfig() (tableConfig *dbutil.TableConfig, err error) {
	for _, v := range a.dataSyncConfig.GetTableConfigs() {
		if v.TableAlias == "ContractPosition" {
			return v, nil
		}
	}
	return nil, fmt.Errorf("cannot find out table config for ContractPosition")
}

func computeBasePositionRecord(currTradeDate string, nextTradeDate string, contractPositions []map[string]interface{}, parValMap map[string]*SecRecord) (t0BasePosition map[string]interface{}, t1BasePosition map[string]interface{}, err error) {
	var t0Qty float64
	var t1Qty float64
	var allQty float64
	var allNotional float64
	var allDirtyCost float64
	var allCleanCost float64
	var allInitCost float64
	var referData map[string]interface{}
	for _, contractPosition := range contractPositions {

		if referData == nil {
			referData = contractPosition
		}

		v, ok, _ := attrutil.GetAttrValue(contractPosition, "quantity", enum.AttrValueType_FLOAT)
		if !ok {
			continue
		}
		qty := v.(float64)

		allQty += qty

		v, _, _ = attrutil.GetAttrValue(contractPosition, "dynamicNotional", enum.AttrValueType_FLOAT)
		allNotional += v.(float64)

		v, _, _ = attrutil.GetAttrValue(contractPosition, "openBondGrossPrice", enum.AttrValueType_FLOAT)
		allDirtyCost += v.(float64) * qty

		v, _, _ = attrutil.GetAttrValue(contractPosition, "openBondNetPrice", enum.AttrValueType_FLOAT)
		allCleanCost += v.(float64) * qty

		v, _, _ = attrutil.GetAttrValue(contractPosition, "initPrice", enum.AttrValueType_FLOAT)
		allInitCost += v.(float64) * qty

		v, _, _ = attrutil.GetAttrValue(contractPosition, "interestSettlementDate", enum.AttrValueType_STRING)
		interestSettlementDate := v.(string)
		if len(interestSettlementDate) != 10 {
			continue
		}

		if interestSettlementDate <= currTradeDate {
			t0Qty += qty
		}

		if interestSettlementDate <= nextTradeDate {
			t1Qty += qty
		}
	}

	if t0Qty > 0 {
		t0BasePosition = initBasePosition(referData, parValMap, allQty, allNotional, allDirtyCost, allCleanCost, allInitCost)
		t0BasePosition["quantity"] = t0Qty
		t0BasePosition["tradeDateFlag"] = "T0"
	}

	if t1Qty > 0 {
		t1BasePosition = initBasePosition(referData, parValMap, allQty, allNotional, allDirtyCost, allCleanCost, allInitCost)
		t1BasePosition["quantity"] = t1Qty
		t1BasePosition["tradeDateFlag"] = "T1"
	}

	return
}

func initBasePosition(referData map[string]interface{}, parValueMap map[string]*SecRecord, allQty float64, allNotional float64, allDirtyCost float64, allCleanCost float64, allInitCost float64) map[string]interface{} {
	v := map[string]interface{}{}
	v["keyCtptyId"] = referData["keyCtptyId"]
	v["ctptyShortName"] = referData["ctptyShortName"]
	v["keyPlanId"] = referData["keyPlanId"]
	v["planCode"] = referData["planCode"]
	v["ultraContractId"] = referData["ultraContractId"]
	v["ultraContractCode"] = referData["ultraContractCode"]
	v["keyInstrumentId"] = referData["keyInstrumentId"]
	v["windCode"] = referData["windCode"]
	v["insShtDesc"] = referData["insShtDesc"]
	v["currency"] = referData["currency"]
	v["exchange"] = referData["exchange"]
	v["longShort"] = referData["longShort"]

	windCode, _ := v["windCode"].(string)
	secRecord, ok := parValueMap[windCode]
	if !ok {
		errStr := fmt.Sprintf("no parValue found for %v, use 100", v["windCode"])
		domain_error.ProcessSevereError(false, 0, nil, errors.New(errStr), errStr)
		secRecord = &SecRecord{
			ParValue: 100,
			SecurityType: "BOND",
		}
	}
	v["parValue"] = secRecord.ParValue
	v["securityType"] = secRecord.SecurityType
	v["baseCashQty"] = allQty * float64(secRecord.ParValue)
	v["baseNotional"] = allNotional
	v["baseDirtyCost"] = allDirtyCost
	v["baseCleanCost"] = allCleanCost
	v["baseInitCost"] = allInitCost

	return v
}

func getKeyOfLine(line map[string]interface{}) (key string, ok bool) {
	//按keyCtptyId、keyPlanId、keyInstrumentId、longShort 分组
	_keyCtptyId, ok := line["keyCtptyId"]
	if !ok {
		return "", false
	}
	keyCtptyId, ok := _keyCtptyId.(int64)
	if !ok {
		return "", false
	}

	_keyPlanId, ok := line["keyPlanId"]
	if !ok {
		return "", false
	}
	keyPlanId, ok := _keyPlanId.(int64)
	if !ok {
		return "", false
	}

	_keyInstrumentId, ok := line["keyInstrumentId"]
	if !ok {
		return "", false
	}
	keyInstrumentId, ok := _keyInstrumentId.(int64)
	if !ok {
		return "", false
	}

	_longShort, ok := line["longShort"]
	if !ok {
		return "", false
	}
	longShort, ok := _longShort.(string)
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%v-%v-%v-%v", keyCtptyId, keyPlanId, keyInstrumentId, longShort), true
}

func (a *TitansFiccDataSyncAdapter) getCurrentAndNextTradeDates() (currTradeDate string, nextTradeDate string, err error) {

	mrkCloseTime, mrkCloseTimeZone := a.dataSyncConfig.GetMrkCloseTime()
	location := timeutil.GetTimeZone(mrkCloseTimeZone)
	dateStr := time.Now().In(location).Format(time.DateOnly)
	// 自然日
	datetimeStr := dateStr + " " + mrkCloseTime

	closeTime, _ := time.ParseInLocation(time.DateTime, datetimeStr, location)
	currTime := time.Now().In(location)
	if currTime.After(closeTime) {
		dateStr = currTime.Add(24 * time.Hour).Format(time.DateOnly)
	}

	appId, appSecret := a.dataSyncConfig.GetSecretInfo()
	currTradeDate, err = tradedate.GetTradeDay(a.dataSyncConfig.GetTradeDateServiceUrl(), dateStr, 0, "NIB", appId, appSecret)
	if err != nil {
		return
	}
	if currTradeDate != datetimeStr[:len(currTradeDate)] {
		log.Printf("currTradeDate=%s, datetimeStr=%s\n", currTradeDate, datetimeStr)
		currTradeDate, err = tradedate.GetTradeDay(a.dataSyncConfig.GetTradeDateServiceUrl(), dateStr, 1, "NIB", appId, appSecret)
		if err != nil {
			return
		}
	}

	nextTradeDate, err = tradedate.GetTradeDay(a.dataSyncConfig.GetTradeDateServiceUrl(), currTradeDate, 1, "NIB", appId, appSecret)
	if err != nil {
		return
	}

	return
}

type SecRecord struct {
	ParValue     int
	SecurityType string
}

// GetParValueMap 从securities表读取WIND_CODE和PAR_VALUE字段，构造map[string]int
func GetParValueMap(db *sql.DB) (result map[string]*SecRecord, err error) {
	// 初始化map
	result = make(map[string]*SecRecord)

	// 执行SQL查询
	rows, err := db.Query("SELECT WIND_CODE, PAR_VALUE, INS_FAMILY FROM securities")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer rows.Close() // 确保查询结果集被关闭[6](@ref)

	// 遍历每一行数据
	for rows.Next() {
		var windCode string
		var parValue int
		var securityType string

		// 将当前行的数据扫描到变量中[9](@ref)
		err := rows.Scan(&windCode, &parValue, &securityType)
		if err != nil {
			return nil, fmt.Errorf("数据扫描失败: %v", err)
		}

		// 将数据存入map
		result[windCode] = &SecRecord{ParValue: parValue, SecurityType: securityType}
	}

	// 检查遍历过程中是否发生错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历行时发生错误: %v", err)
	}

	return result, nil
}
