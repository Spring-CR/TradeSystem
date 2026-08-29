package domain_error

func init() {
	registerErrMsg(map[string]string{
		ASSET_UNIT_ERR_CODE:                  `资产单元相关错误`,
		CANNOT_FIND_ALL_ASSET_UNITS_ERR_CODE: `无法获得资产单元记录`,
		CANNOT_CREATE_ASSET_UNIT_ERR_CODE:    `无法创建资产单元，账号：%s，账号名称：%s，组合编号：%s，组合名称：%s`,
		CANNOT_DELETE_ASSET_UNIT_ERR_CODE:    `无法删除资产单元，账号：%s，组合编号：%s`,
	})
}
