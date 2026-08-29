package api_const

const (
	ApiVersionPrefix = "/api/v1"
)

const (
	RoutePing       = ApiVersionPrefix + "/ping"
	SubRouteQuery   = "/query"
	SubRouteDelete  = "/delete"
	SubRouteFindAll = "/all"

	ParamKey              = "key"
	ParamDate             = "date"
	ParamDailyInstrNo     = "dailyInstrNo"
	ParamIndexDailyModify = "indexDailyModify"
	ParamStockSerialNo    = "stockSerialNo"
)
