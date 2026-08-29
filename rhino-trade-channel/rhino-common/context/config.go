package context

import (
	"database/sql"
	"rhino-common/context/constant"
)

var (
	Lang                 constant.LangType
	DB                   *sql.DB
	DB_GFQG              *sql.DB
	DefaultATPUser       = "4006"
	RetryIntervalTimes   = 10
	RetryIntervalSeconds = 3
)
