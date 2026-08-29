package schema

type DataArchivingLog struct {
	ID                   int64
	SystemCode           string `sql:"unique: pk_dal, size: 32"`
	BusinessCode         string `sql:"index: pk_dal, size: 32"`
	ArchivingDate        string `sql:"index: pk_dal, size: 8"`
	TaskName             string `sql:"index: pk_dal, size: 32"`
	FirstArchivingTime   int64
	CurrentArchivingTime int64
	ExecCount            int
	Complete             bool
	CompletePhase        int
	//FailLog              string `sql:"type: MEDIUMTEXT"`
}
