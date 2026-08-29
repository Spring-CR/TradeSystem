package domain_trade_instr_test

import (
	"database/sql"
	"log"
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	domain_asset_unit "rhino-instr/domain/asset-unit"
	"rhino-instr/domain/status"
	domain_task_instr "rhino-instr/domain/task-instr"
	trade_channel "rhino-instr/domain/trade-channel"
	domain_trade_instr "rhino-instr/domain/trade-instr"
	"rhino-instr/schema"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbUrl = "root:guangfa4cool@tcp(10.51.136.72:56322)/olts_tradedesk?charset=utf8"
)

var (
	tradeChannel *trade_channel.KafkaTradeChannel
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Second * 500)
	context.DB = db

	var de *domain_error.Error
	tradeChannel, de = trade_channel.NewKafkaTradeChannel("rhino-instr-req-dev", "rhino-instr-resp-dev", []string{"10.128.13.230:30490", "10.128.13.233:30490", "10.128.13.232:30490"})
	processDomainError(de)

	status.StatusObserve(tradeChannel)

	context.RetryIntervalSeconds = 0
	context.RetryIntervalTimes = 1
}

func TestE2ESingleTrade(t *testing.T) {
	reportCode := "IF2406"
	queryAssetUnit(t)
	date, dailyInstrNo, indexDailyModify := issueInstruction(reportCode)
	stockSerialNo := queryInstruction(date, dailyInstrNo, indexDailyModify, reportCode, t)
	executeInstruction(date, dailyInstrNo, indexDailyModify, stockSerialNo)
	queryInstructionDealStatus()
}

// 下达指令
func issueInstruction(reportCode string) (date int, dailyInstrNo int64, indexDailyModify int64) {
	dateStr := time.Now().Format("20060102")
	dateNum, err := strconv.Atoi(dateStr)
	processError(err)

	var de *domain_error.Error
	dailyInstrNo, indexDailyModify, _, date, de = domain_task_instr.CreateInsertTaskInstr(
		"4012", "400012001", schema.TaskInstrInstrTypeStock,
		//"3999", "400005|4005_000", schema.TaskInstrInstrTypeStock,
		dateNum, dateNum, 90000, 150000,
		"user1",
		"",
		"user2",
		[]*schema.TaskStock{
			{
				MarketNo:         schema.TaskInstrMarketNoZJ,
				ReportCode:       reportCode,
				EntrustDirection: schema.TaskInstrEntrustDirectionBuy,
				OpenClose:        schema.TaskInstrOpenCloseOpen,
				//Amount:           500,
				//Amount:           3,
				Amount:           90,
				Balance:          0, // 不需要
				Price:            5000,
				ContractSize:     300,
				InvestType:       "a", // 投机、套保、套利，由于实际没有使用，可以不用传
			},
		},
	)

	processDomainError(de)

	log.Println("下达指令成功！")

	return
}

// 查询指令
func queryInstruction(date int, dailyInstrNo int64, indexDailyModify int64, reportCode string, t *testing.T) (stockSerialNo int64) {

	result, total, de := domain_task_instr.FindTaskInstrs([]*dbutil.FieldCondition{
		{Field: "date", ValueType: 0, Value: date},
		{Field: "daily_instr_no", ValueType: 0, Value: dailyInstrNo},
		{Field: "index_daily_modify", ValueType: 0, Value: indexDailyModify},
		{Field: "report_code", ValueType: 0, Value: reportCode},
	}, 1, 0)
	processDomainError(de)

	expectedTotal := 1
	if total != expectedTotal {
		t.Errorf("total = %d; want %d", total, expectedTotal)
	}

	stockSerialNo = result[0].StockSerialNo

	log.Println("查询指令成功！")

	return
}

// 查询资产单元
func queryAssetUnit(t *testing.T) {
	assetUnits, de := domain_asset_unit.FindAllAssetUnits()
	processDomainError(de)
	found := false
	account_no := "4012"
	combi_no := "400012001"
	for _, assetUnit := range assetUnits {
		if assetUnit.AccountNo == account_no && assetUnit.CombiNo == combi_no {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("cannot find asset unit with account_no:%s, combi_no:%s", account_no, combi_no)
	}

	log.Println("查询资产单元成功！")
}

// 执行指令
func executeInstruction(date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) {
	//de := domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "user2", "4012", tradeChannel)
	//de := domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 3520, 100, "user2", "3999", tradeChannel)
	//de := domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 1, "user2", "3999", tradeChannel)
	_, de := domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 30, "user2", context.DefaultATPUser, tradeChannel)
	processDomainError(de)

	log.Println("#1 执行指令成功！")

	//de = domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 1, "user2", "3999", tradeChannel)
	_, de = domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 30, "user2", context.DefaultATPUser, tradeChannel)
	processDomainError(de)

	log.Println("#2 执行指令成功！")

	//de = domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 1, "user2", "3999", tradeChannel)
	_, de = domain_trade_instr.ExecuteSingleTradeInstr(date, dailyInstrNo, indexDailyModify, stockSerialNo, "2", 5000, 30, "user2", context.DefaultATPUser, tradeChannel)
	processDomainError(de)

	log.Println("#3 执行指令成功！")
}

// 查询指令证券的成交状态
func queryInstructionDealStatus() {
	time.Sleep(20 * time.Minute)
}

func processError(err error) {
	if err != nil {
		panic(err)
	}
}

func processDomainError(err *domain_error.Error) {
	if err != nil {
		panic(err.ErrorString())
	}
}
