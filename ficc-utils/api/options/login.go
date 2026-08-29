package options

type LoginOpt struct {
	PhoneNum string `json:"phoneNum"`
	Password string `json:"password"`
}

type LoginResult struct {
	ApiToken      string `json:"apiToken"`
	EffectiveDate string `json:"effectiveDate"`
	Accounts      []int  `json:"accounts"`
}
