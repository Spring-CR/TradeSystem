package domain_error

func init() {
	registerErrMsg(map[string]string{
		INSTR_ERR_CODE:                               `指令处理相关错误`,
		CANNOT_LOCK_INSTR_MAIN_RECORD_ERR_CODE:       `无法锁住指令主表以插入新数据`,
		CANNOT_LOCK_INSTR_SECONDLY_RECORD_ERR_CODE:   `无法锁住指令副表以插入新数据`,
		CANNOT_INSERT_MAIN_INSTR_RECORD_ERR_CODE:     `插入任务指令主表记录异常`,
		CANNOT_INSERT_SECONDLY_INSTR_RECORD_ERR_CODE: `插入任务指令副表记录异常`,
		CANNOT_FIND_INSTR_RECORD_ERR_CODE:            `查询任务指令记录异常`,
		ILLEGAL_INSTR_CODE_ERR_CODE:                  `指令编号的格式不正确，请输入{日期}-{指令序号}-{指令修改序号}`,
		CANNOT_GET_INSTR_STOCK_ERR_CODE:              "无法获得指令证券明细记录, %v/%v/%v/%v",
		ONLY_NONE_EXE_STOCK_INSTR_CAN_BE_EXE:         "只有未曾委托的证券才能执行，证券%s不符合该条件",
		CANNOT_STATIS_INSTR_STOCK_ERR_CODE:           "无法统计指令证券明细, %v/%v/%v/%v",
		CANNOT_OVER_ENTRUST_ERR_CODE:                 "当前已委托数量已经满额，请勿超量委托",
		CANNOT_GET_TRADE_DESK_ORDER_ID_ERR_CODE:      "无法根据指令编号获取交易台母单编号，交易指令编号：%s",
	})
}
