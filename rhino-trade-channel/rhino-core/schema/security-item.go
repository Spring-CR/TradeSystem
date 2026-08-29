package schema

type SecurityItem struct {
	ID                       int64
	SecurityLibCode          string `sql:"unique: pk_si, size: 32"`
	SecurityZhName           string
	SecurityEnName           string
	Symbol                   string `sql:"index: pk_si, size: 64"`
	SymbolSfx                string `sql:"size: 8"`
	SecurityExchangeSymbol   string `sql:"index: pk_si, size: 64"`
	SecurityISIN             string `sql:"index: pk_si, size: 64"`
	SecurityRIC              string `sql:"index: pk_si, size: 64"`
	SecurityExchange         string `sql:"size: 8"`
	SecurityExchangeRegion   string `sql:"size: 4"`
	ContractMultiplier       float64
	Currency                 string `sql:"size: 4"`
	LotSize                  int
	IssueDate                string `sql:"size: 16"`
	ContractMonth            string `sql:"size: 12"`
	ExpireDate               string `sql:"size: 16"`
	SecurityType             string `sql:"size: 16"`
	UnderlyingSecurityCode   string `sql:"size: 64"`
	UnderlyingSecurityZhName string `sql:"size: 128"`
	UnderlyingSecurityEnName string `sql:"size: 64"`
	PutOrCall                string `sql:"size: 2"`
}
