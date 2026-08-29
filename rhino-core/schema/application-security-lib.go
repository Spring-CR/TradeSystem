package schema

type ApplicationSecurityLib struct {
	ID              int64
	SystemCode      string `sql:"unique: uq_asl, size: 32"`
	BusinessCode    string `sql:"index: uq_asl, size: 32"`
	SecurityLibCode string `sql:"index: uq_asl, size: 32"`
}
