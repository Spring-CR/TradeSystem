package ficc_fut

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/attrutil"
	"rhino-common/utils/dbutil"
	"rhino-core/domain_cfg"
	"rhino-plugins/data_sync_plugin/utils"
	"sort"
	"strconv"
	"strings"

	ficc_fut_posi "rhino-plugins/order_position_plugin/ficc_fut"
)

type FiccFutDataSyncAdapter struct {
	dataSyncConfig *domain_cfg.DataSyncConfig
}

func NewFiccFutDataSyncAdapter(dataSyncConfig *domain_cfg.DataSyncConfig) (adapter *FiccFutDataSyncAdapter, de *domain_error.Error) {
	log.Printf("construct FiccFutDataSyncAdapter...")
	adapter = &FiccFutDataSyncAdapter{dataSyncConfig: dataSyncConfig}
	return
}

func (a *FiccFutDataSyncAdapter) RefineCsvContent(tableConfig *dbutil.TableConfig, rawCsv []byte) (newCsv []byte, err error) {

	if !strings.HasPrefix(tableConfig.TableAlias, "PositionBase") {
		return rawCsv, nil
	}

	defer func() {
		os.WriteFile("/tmp/"+tableConfig.TableName+".csv", newCsv, 0644)
	}()

	marginRatioMap, err1 := a.getMarginRatioRecordMap(a.dataSyncConfig.GetAppDB())
	if err1 != nil {
		return rawCsv, err1
	}

	exchangeRateRecordMap, err1 := a.getExchangeRateRecordMap(a.dataSyncConfig.GetAppDB())
	if err1 != nil {
		return rawCsv, err1
	}

	// 解析csv数据内容
	contractPositionTableCinfig, err1 := a.getContractPositionTableConfig()
	if err1 != nil {
		return rawCsv, err1
	}

	lines, err1 := utils.ConvertCsvToMap(rawCsv, contractPositionTableCinfig)
	if err1 != nil {
		return rawCsv, err1
	}

	log.Printf("RefineCsvContent for:%s, contractPosition table line count:%d\n", tableConfig.TableName, lines)

	m := make(map[string]*ficc_fut_posi.PositionRecord)
	for _, line := range lines {

		key, ok := getKeyOfLine(line)
		if !ok {
			continue
		}

		positionRecord, ok1 := m[key]
		if ok1 {
			log.Printf("aggregate for key:%v\n", key)
			a.aggregate(tableConfig, marginRatioMap, positionRecord, line)
		} else {
			log.Printf("newPositionRecordFromLine for key:%v\n", key)
			positionRecord = a.newPositionRecordFromLine(tableConfig, marginRatioMap, line)
			m[key] = positionRecord
		}
	}

	if len(m) == 0 {
		return
	}

	var positions []*ficc_fut_posi.PositionRecord
	for k, v1 := range m {
		pearKey := k
		if strings.HasSuffix(k, "LONG") {
			pearKey = strings.ReplaceAll(pearKey, "LONG", "SHORT")
		} else {
			pearKey = strings.ReplaceAll(pearKey, "SHORT", "LONG")
		}
		v2, ok := m[pearKey]
		if ok {
			v3 := a.diff(v1, v2)
			if v3 != nil {
				positions = append(positions, v3)
			}
			delete(m, pearKey)
		} else {
			positions = append(positions, v1)
		}
	}

	positions, err = a.updateByTradeActionResps(tableConfig, positions)
	if err != nil {
		return nil, err
	}

	sort.Slice(positions, func(i, j int) bool {
		if positions[i].Account < positions[j].Account {
			return true
		}
		return positions[i].Symbol2 < positions[j].Symbol2
	})

	// 填充值
	for i, p := range positions {
		a.fill(p, exchangeRateRecordMap)
		js, _ := json.Marshal(p)
		log.Printf("#%d Position: %s\n", i, js)
	}

	return a.generateCSV(positions)
}

func (a *FiccFutDataSyncAdapter) generateCSV(records []*ficc_fut_posi.PositionRecord) ([]byte, error) {
	buf := &bytes.Buffer{}
	writer := csv.NewWriter(buf)
	defer writer.Flush()

	// 获取结构体类型
	var typ reflect.Type
	if len(records) > 0 && records[0] != nil {
		typ = reflect.TypeOf(*records[0])
	} else {
		typ = reflect.TypeOf(ficc_fut_posi.PositionRecord{})
	}

	// 写入标题行
	numFields := typ.NumField()
	header := make([]string, numFields)
	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		header[i] = field.Name
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// 写入每行数据
	for _, rec := range records {
		row := make([]string, numFields)
		if rec == nil {
			// 空行全为空字符串
			continue
		} else {
			val := reflect.ValueOf(*rec)
			for i := 0; i < numFields; i++ {
				fieldVal := val.Field(i)
				row[i] = valueToString(fieldVal)
			}
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func valueToString(v reflect.Value) string {
	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)
	case reflect.Float32, reflect.Float64:
		// 使用fmt格式化，避免科学计数法，保留足够精度
		return strconv.FormatFloat(v.Float(), 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

func (a *FiccFutDataSyncAdapter) fill(position *ficc_fut_posi.PositionRecord, exchangeRateRecordMap map[string]*ExchangeRateRecord) {

	exchangeRateRecord, ok := exchangeRateRecordMap[position.Symbol2]
	if !ok {
		position.ExchangeRateCNY = 1.0
	} else {
		position.ExchangeRateCNY = exchangeRateRecord.ExchangeRateCNY
	}

	position.NetPosition = position.InitNetPosition
	position.Key = fmt.Sprintf("%v-%v", position.Account, position.Symbol2)

	if position.NetPosition > 0 {

		position.LongAvailablePosition = position.NetPosition
		position.ShortAvailablePosition = 0
		position.ShortPriceCost = 0
		position.ShortPriceWithFeeCost = 0
		position.ShortPriceCNYWithFeeCost = 0
		position.LongAvgPrice = position.LongPriceCost / position.NetPosition / position.ContractMultiplier
		position.LongAvgPriceWithFee = position.LongPriceWithFeeCost / position.NetPosition / position.ContractMultiplier
		position.ShortAvgPrice = 0
		position.ShortAvgPriceWithFee = 0
		position.BuyOrderLeftQty = 0
		position.SellOrderLeftQty = 0
		position.BuyOrderLeftCost = 0
		position.SellOrderLeftCost = 0

	} else if position.NetPosition < 0 {

		position.LongAvailablePosition = 0
		position.ShortAvailablePosition = -position.NetPosition
		position.LongPriceCost = 0
		position.LongPriceWithFeeCost = 0
		position.LongPriceCNYWithFeeCost = 0
		position.LongAvgPrice = 0
		position.LongAvgPriceWithFee = 0
		position.ShortAvgPrice = position.ShortPriceCost / -position.NetPosition / position.ContractMultiplier
		position.ShortAvgPriceWithFee = position.ShortPriceWithFeeCost / -position.NetPosition / position.ContractMultiplier
		position.BuyOrderLeftQty = 0
		position.SellOrderLeftQty = 0
		position.BuyOrderLeftCost = 0
		position.SellOrderLeftCost = 0

	} else {

		position.LongAvailablePosition = 0
		position.ShortAvailablePosition = 0
		position.LongPriceCost = 0
		position.LongPriceWithFeeCost = 0
		position.LongPriceCNYWithFeeCost = 0
		position.ShortPriceCost = 0
		position.ShortPriceWithFeeCost = 0
		position.ShortPriceCNYWithFeeCost = 0
		position.LongAvgPrice = 0
		position.LongAvgPriceWithFee = 0
		position.ShortAvgPrice = 0
		position.ShortAvgPriceWithFee = 0
		position.BuyOrderLeftQty = 0
		position.SellOrderLeftQty = 0
		position.BuyOrderLeftCost = 0
		position.SellOrderLeftCost = 0
	}

	position.InitLongPriceCost = position.LongPriceCost
	position.InitLongPriceWithFeeCost = position.LongPriceWithFeeCost
	position.InitShortPriceCost = position.ShortPriceCost
	position.InitShortPriceWithFeeCost = position.ShortPriceWithFeeCost
}

func (a *FiccFutDataSyncAdapter) diff(p1, p2 *ficc_fut_posi.PositionRecord) *ficc_fut_posi.PositionRecord {

	if p1.InitNetPosition+p2.InitNetPosition == 0 {
		return nil
	}

	if p1.ContractBaseDate > p2.ContractBaseDate {
		p2.ContractBaseDate = p1.ContractBaseDate
	}

	avgPrice := (p1.ShortPriceCost + p1.LongPriceCost + p2.ShortPriceCost + p2.LongPriceCost) / (math.Abs(p1.InitNetPosition) + math.Abs(p2.InitNetPosition)) / p2.ContractMultiplier
	avgPriceWithFee := (p1.ShortPriceWithFeeCost + p1.LongPriceWithFeeCost + p2.ShortPriceWithFeeCost + p2.LongPriceWithFeeCost) / (math.Abs(p1.InitNetPosition) + math.Abs(p2.InitNetPosition)) / p2.ContractMultiplier
	avgPriceCNYWithFee := (p1.ShortPriceCNYWithFeeCost + p1.LongPriceCNYWithFeeCost + p2.ShortPriceCNYWithFeeCost + p2.LongPriceCNYWithFeeCost) / (math.Abs(p1.InitNetPosition) + math.Abs(p2.InitNetPosition)) / p2.ContractMultiplier

	p2.InitNetPosition += p1.InitNetPosition

	if p2.InitNetPosition > 0 {

		p2.LongPriceCost = avgPrice * p2.InitNetPosition * p2.ContractMultiplier
		p2.ShortPriceCost = 0

		p2.LongPriceWithFeeCost = avgPriceWithFee * p2.InitNetPosition * p2.ContractMultiplier
		p2.ShortPriceWithFeeCost = 0

		p2.LongPriceCNYWithFeeCost = avgPriceCNYWithFee * p2.InitNetPosition * p2.ContractMultiplier
		p2.ShortPriceCNYWithFeeCost = 0

	} else {

		p2.LongPriceCost = 0
		p2.ShortPriceCost = avgPrice * -p2.InitNetPosition * p2.ContractMultiplier

		p2.LongPriceWithFeeCost = 0
		p2.ShortPriceWithFeeCost = avgPriceWithFee * -p2.InitNetPosition * p2.ContractMultiplier

		p2.LongPriceCNYWithFeeCost = 0
		p2.ShortPriceCNYWithFeeCost = avgPriceCNYWithFee * -p2.InitNetPosition * p2.ContractMultiplier
	}

	return p2
}

func (a *FiccFutDataSyncAdapter) aggregate(tableConfig *dbutil.TableConfig, marginRatioMap map[string]*MarginRatioRecord, base *ficc_fut_posi.PositionRecord, line map[string]interface{}) {

	next := a.newPositionRecordFromLine(tableConfig, marginRatioMap, line)

	if next.ContractBaseDate > base.ContractBaseDate {
		base.ContractBaseDate = next.ContractBaseDate
	}

	base.InitNetPosition += next.InitNetPosition
	base.LongPriceCost += next.LongPriceCost
	base.ShortPriceCost += next.ShortPriceCost
	base.LongPriceWithFeeCost += next.LongPriceWithFeeCost
	base.ShortPriceWithFeeCost += next.ShortPriceWithFeeCost
	base.LongPriceCNYWithFeeCost += next.LongPriceCNYWithFeeCost
	base.ShortPriceCNYWithFeeCost += next.ShortPriceCNYWithFeeCost
}

func (a *FiccFutDataSyncAdapter) newPositionRecordFromLine(tableConfig *dbutil.TableConfig, marginRatioMap map[string]*MarginRatioRecord, line map[string]interface{}) *ficc_fut_posi.PositionRecord {

	positionRecord := &ficc_fut_posi.PositionRecord{}

	account, _, _ := attrutil.GetAttrValue(line, "COUNTERPARTY_ID", enum.AttrValueType_INT)
	positionRecord.Account = account.(int)
	positionRecord.CounterpartyID = positionRecord.Account

	counterparty, _, _ := attrutil.GetAttrValue(line, "COUNTERPARTY", enum.AttrValueType_STRING)
	positionRecord.Counterparty = counterparty.(string)

	symbol2, _, _ := attrutil.GetAttrValue(line, "SYMBOL2", enum.AttrValueType_STRING)
	positionRecord.Symbol2 = symbol2.(string)

	symbolName, _, _ := attrutil.GetAttrValue(line, "SECURITY_NAME", enum.AttrValueType_STRING)
	positionRecord.SymbolName = symbolName.(string)

	currency, _, _ := attrutil.GetAttrValue(line, "CURRENCY_UNDERLYING", enum.AttrValueType_STRING)
	positionRecord.Currency = currency.(string)

	planCode, _, _ := attrutil.GetAttrValue(line, "PLAN_CODE", enum.AttrValueType_STRING)
	positionRecord.PlanCode = planCode.(string)

	ultraContractCode, _, _ := attrutil.GetAttrValue(line, "ULTRA_CONTRACT_CODE", enum.AttrValueType_STRING)
	positionRecord.UltraContractCode = ultraContractCode.(string)

	securityExchange, _, _ := attrutil.GetAttrValue(line, "SECURITY_EXCHANGE", enum.AttrValueType_STRING)
	positionRecord.SecurityExchange = securityExchange.(string)

	// 取了具体的期货类型
	securityType, _, _ := attrutil.GetAttrValue(line, "FUT_TYPE", enum.AttrValueType_STRING)
	positionRecord.SecurityType = securityType.(string)

	productCode, _, _ := attrutil.GetAttrValue(line, "FUT_CODE", enum.AttrValueType_STRING)
	positionRecord.ProductCode = productCode.(string)

	contractMultiplier, _, _ := attrutil.GetAttrValue(line, "CONTRACT_MULTIPLIER", enum.AttrValueType_FLOAT)
	positionRecord.ContractMultiplier = contractMultiplier.(float64)

	marginRatioRecord := marginRatioMap[fmt.Sprintf("%v-%v-%v", account, planCode, productCode)]
	if marginRatioRecord != nil {
		positionRecord.LongMarginRatio = marginRatioRecord.MarginRatio
		positionRecord.ShortMarginRatio = marginRatioRecord.MarginRatio
	}

	if strings.HasSuffix(tableConfig.TableAlias, "Dms") {
		positionRecord.ExchangeArea = "dms"
	}

	if strings.HasSuffix(tableConfig.TableAlias, "Ovs") {
		positionRecord.ExchangeArea = "ovs"
	}

	// 以下属性需要动态调整

	contractBaseDate, _, _ := attrutil.GetAttrValue(line, "TITANS_POS_DATE", enum.AttrValueType_STRING)
	positionRecord.ContractBaseDate = contractBaseDate.(string)
	positionRecord.ContractBaseDate = strings.ReplaceAll(positionRecord.ContractBaseDate, "-", "")
	positionRecord.ContractBaseDate = positionRecord.ContractBaseDate[:8]

	initNetPosition, _, _ := attrutil.GetAttrValue(line, "AMOUNT", enum.AttrValueType_FLOAT)
	positionRecord.InitNetPosition = initNetPosition.(float64)

	//cost, _, _ := attrutil.GetAttrValue(line, "BASE_DYNAMIC_NOTIONAL", enum.AttrValueType_FLOAT)
	//改用标的币种的名义本金
	cost, _, _ := attrutil.GetAttrValue(line, "DYNAMIC_NOTIONAL", enum.AttrValueType_FLOAT)
	costCNY, _, _ := attrutil.GetAttrValue(line, "BASE_DYNAMIC_NOTIONAL", enum.AttrValueType_FLOAT)

	logShort, _, _ := attrutil.GetAttrValue(line, "LONG_SHORT", enum.AttrValueType_STRING)

	if logShort.(string) != "LONG" {

		positionRecord.ShortPriceCost = cost.(float64)
		positionRecord.ShortPriceWithFeeCost = positionRecord.ShortPriceCost
		positionRecord.ShortPriceCNYWithFeeCost = costCNY.(float64)
		positionRecord.InitNetPosition = -positionRecord.InitNetPosition

	} else {

		positionRecord.LongPriceCost = cost.(float64)
		positionRecord.LongPriceWithFeeCost = positionRecord.LongPriceCost
		positionRecord.LongPriceCNYWithFeeCost = costCNY.(float64)
	}

	js, _ := json.Marshal(positionRecord)
	log.Printf("======>newPositionRecordFromLine:%s\n", js)

	return positionRecord
}

func getKeyOfLine(line map[string]interface{}) (key string, ok bool) {

	defer func() {
		if !ok {
			domain_error.ProcessSevereError(false, 0, nil, errors.New("fail to parse key"), fmt.Sprintf("key not found for line: %v\n", line))
		}
	}()

	_keyCtptyId, ok := line["COUNTERPARTY_ID"]
	if !ok {
		return "", false
	}
	keyCtptyId, ok := _keyCtptyId.(int64)
	if !ok {
		return "", false
	}

	_symbol2, ok := line["SYMBOL2"]
	if !ok {
		return "", false
	}
	symbol2, ok := _symbol2.(string)
	if !ok {
		return "", false
	}

	_longShort, ok := line["LONG_SHORT"]
	if !ok {
		return "", false
	}
	longShort, ok := _longShort.(string)
	if !ok {
		return "", false
	}

	return fmt.Sprintf("%v-%v-%v", keyCtptyId, symbol2, longShort), true
}

func (a *FiccFutDataSyncAdapter) getContractPositionTableConfig() (tableConfig *dbutil.TableConfig, err error) {
	for _, v := range a.dataSyncConfig.GetTableConfigs() {
		if v.TableAlias == "ContractPositionOvs" {
			return v, nil
		}
	}
	return nil, fmt.Errorf("cannot find out table config for ContractPosition")
}

type MarginRatioRecord struct {
	Account     int
	PlanCode    string
	ProductCode string
	MarginRatio float64
}

// GetMarginRatioMap 从securities表读取WIND_CODE和PAR_VALUE字段，构造map[string]int
func (a *FiccFutDataSyncAdapter) getMarginRatioRecordMap(db *sql.DB) (result map[string]*MarginRatioRecord, err error) {
	// 初始化map
	result = make(map[string]*MarginRatioRecord)

	// 执行SQL查询
	rows, err := db.Query("SELECT COUNTERPARTY_ID, PLAN_CODE, FUT_CODE, MARGIN_RATIO FROM MarginThreshold WHERE COUNTERPARTY_ID IS NOT NULL AND PLAN_CODE IS NOT NULL AND FUT_CODE IS NOT NULL AND MARGIN_RATIO IS NOT NULL")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer rows.Close() // 确保查询结果集被关闭[6](@ref)

	// 遍历每一行数据
	for rows.Next() {
		var account int
		var planCode string
		var productCode string
		var marginRatio float64

		// 将当前行的数据扫描到变量中[9](@ref)
		err := rows.Scan(&account, &planCode, &productCode, &marginRatio)
		if err != nil {
			return nil, fmt.Errorf("数据扫描失败: %v", err)
		}

		// 将数据存入map
		result[fmt.Sprintf("%v-%v-%v", account, planCode, productCode)] = &MarginRatioRecord{Account: account, PlanCode: planCode, ProductCode: productCode, MarginRatio: marginRatio}
	}

	// 检查遍历过程中是否发生错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历行时发生错误: %v", err)
	}

	return result, nil
}

type ExchangeRateRecord struct {
	Symbol2         string
	ExchangeRateCNY float64
}

// GetMarginRatioMap 从securities表读取WIND_CODE和PAR_VALUE字段，构造map[string]int
func (a *FiccFutDataSyncAdapter) getExchangeRateRecordMap(db *sql.DB) (result map[string]*ExchangeRateRecord, err error) {
	// 初始化map
	result = make(map[string]*ExchangeRateRecord)

	// 执行SQL查询
	rows, err := db.Query("SELECT DISTINCTROW  SYMBOL2, EXCHANGE_RATE_CNY from Security WHERE SYMBOL2 IS NOT NULL AND EXCHANGE_RATE_CNY IS NOT NULL")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	defer rows.Close() // 确保查询结果集被关闭[6](@ref)

	// 遍历每一行数据
	for rows.Next() {
		var symbol2 string
		var exchangeRateCNY float64

		err := rows.Scan(&symbol2, &exchangeRateCNY)
		if err != nil {
			return nil, fmt.Errorf("数据扫描失败: %v", err)
		}

		// 将数据存入map
		result[symbol2] = &ExchangeRateRecord{Symbol2: symbol2, ExchangeRateCNY: exchangeRateCNY}
	}

	// 检查遍历过程中是否发生错误
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历行时发生错误: %v", err)
	}

	return result, nil
}
