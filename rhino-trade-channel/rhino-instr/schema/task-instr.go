package schema

// 投资任务指令
type TaskInstr struct {
	// 表主键，自增
	ID int64

	// 业务日期，记录创建的日期(8位，类型万徳)
	Date int `sql:"index: pk_ti;i_ti_date"`
	// 指令的ID，每天从1开始，全系统当日唯一
	DailyInstrNo int64 `sql:"unique: pk_ti"`
	// 指令新下达时，该字段的值为1，当指令有修改时，该字段的值依次递增
	IndexDailyModify int64 `sql:"index: pk_ti"`

	// 批量指令的ID，每天从1开始，全系统当日唯一
	BatchSerialNo int64 
	// 前指令修改次序
	IndexLastModify int64
	// 账号id
	AccountNo string `sql:"size: 16"`
	// Combination，组合编号
	CombiNo string `sql:"size: 16"`
	// 指令类型：1:个股 2:组合 3:个股批量 4:组合批量
	InstrType string `sql:"size: 1"`
	// 指令开始日期，8位，yyyyMMdd
	BeginDate int
	// 指令结束日期，8位，yyyyMMdd
	EndDate int
	// 指令开始时间，6位，HHmmss
	BeginTime int
	// 指令结束时间，6位，HHmmss
	EndTime int
	// 指令下达日期
	DirectDate int
	// 指令下达的时间，6位，HHmmss
	DirectTime int
	// 指令下达人
	DirectOperator string `sql:"index: i_ti_direct_operator, size: 32"`
	// 指令修改日期，8位
	ModifyDate int
	// 指令修改时间，6位，HHmmss
	ModifyTime int
	// 指令修改人
	ModifyOperator string `sql:"size: 32"`
	// 指令修改原因
	ModifyReason string`sql:"size: 128"`
	// 分发日期（8位）
	DispenseDate int
	// 分发时间，6位，HHmmss
	DispenseTime int
	// 分发人
	DispenseOperator string `sql:"size: 32"`
	// 分发拒绝原因
	DispenseRefuseReason string`sql:"size: 128"`
	// 撤销日期，8位
	CancelDate int
	// 撤销时间，6位，HHmmss
	CancelTime int
	// 撤销人
	CancelOperator string `sql:"size: 32"`
	// 撤销原因
	CancelReason string`sql:"size: 128"`
	// 执行人
	Operator string `sql:"index: i_ti_operator, size: 32"`
	// 指令状态：1、有效 2、已修改 3、已撤销 4、已暂停 5、审批拒绝 6、分发拒绝 7、指令录入 8、分仓失败 9、草稿指令 a、临时下达 b、临时修改 c、撤消失败
	InstrStatus string`sql:"size: 1"`
	// 分发状态：1：未分发 2：已分发 3：分发拒绝
	DispenseStatus string`sql:"size: 1"`
	// 委托完成状态：1：未执行 2：部分执行 3：完成
	EntrustExecuteStatus string`sql:"size: 1"`
	// 成交完成状态：1：未执行 2：部分执行 3：完成
	// 因交易台不区分委托和成交状态，这里直接复用交易台的状态，并增加
	DealExecuteStatus string`sql:"size: 1"`
	// 记录创建的微秒时间戳
	CreateTime int64
	// 业务类型：1：交易所网上 2：银行间二级市场 3：QDII 4：股指期货 5：LOF指令 A：开放式基金 B：存款类 C：交易所网下申购 D：银行间债券一级市场
    // E：交易所大宗交易 F：交易所债券打包 G：公开市场招投标 H：项目投标业务 I：上交易固定收益平台 J：上交易固定收益平台做市
	BusinessType string `sql:"size: 2"`
	// 锁定标志
	LockFlag int
	// 指定交易员
	LimitOperator string `sql:"size: 32"`
	// 公司号
	OrgId int64
	// 部门号
	DeptId int64
	// IP地址
	IpAddress string `sql:"size: 16"`
	// Mac地址
	Mac string `sql:"size: 20"`
	// 磁盘序列号
	VolserialNo string `sql:"size: 10"`
}

