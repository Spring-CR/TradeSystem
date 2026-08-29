package schema

// 资产单元
type AssetUnit struct {
	// 表主键，自增
	ID int64 `json:"-"`
	// 账号id
	AccountNo string `sql:"unique: uq_au, index: i_au_an, size: 16" json:"account_no"`
	// 账号名称
	AccountName string `sql:"size: 256" json:"account_name"`
	// 组合编号
	CombiNo string `sql:"index: uq_au, size: 16" json:"combi_no"`
	// 组合名称
	CombiName string `sql:"size: 256" json:"combi_name"`
}
