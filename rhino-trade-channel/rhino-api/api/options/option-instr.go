package options

import "rhino-instr/schema"

type IssueInstructionOption struct {
	AccountNo      string              `json:"account_no"`
	CombiNo        string              `json:"combi_no"`
	InstrType      string              `json:"instr_type"`
	BeginDate      int                 `json:"begin_date"`
	EndDate        int                 `json:"end_date"`
	BeginTime      int                 `json:"begin_time"`
	EndTime        int                 `json:"end_time"`
	DirectOperator string              `json:"direct_operator"`
	BusinessType   string              `json:"business_type"`
	LimitOperator  string              `json:"limit_operator"`
	Stocks         []*schema.TaskStock `json:"stocks"`
	CreateTime     int64               `json:"createTime"` // 可缺省
}

type IssueInstructionRet struct {
	DailyInstrNo     int64 `json:"daily_instr_no"`
	IndexDailyModify int64 `json:"index_daily_modify"`
	BatchSerialNo    int64 `json:"batch_serial_no"`
	DateNum          int   `json:"date_num"`
}

type ExecuteSingleTradeInstrOption struct {
	Date             int     `json:"date"`
	DailyInstrNo     int64   `json:"daily_instr_no"`
	IndexDailyModify int64   `json:"index_daily_modify"`
	StockSerialNo    int64   `json:"stock_serial_no"`
	InstrOperator    string  `json:"instr_operator"`
	OrdType          string  `json:"ord_type"`
	Price            float64 `json:"price"`
	Amount           float64 `json:"amount"`
}
