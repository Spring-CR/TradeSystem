package api_const

const (
	ApiVersionPrefix = "/api/v1"
)

const (
	RoutePing  = ApiVersionPrefix + "/ping"
	RouteLogin = ApiVersionPrefix + "/login"
)

const (
	ParamAccount = "account" // 产品账号ID，即交易对手ID，整形
	ParamPlanCode = "planCode" // 业务方案编号，文本
	ParamSymbol = "symbol" //标的代码, 文本
	ParamLongShort = "longShort" //多空方向, 文本
	ParamSide      = "side"      //交易方向, 文本
	ParamPageNum = "pageNum" //分页，页码, 整形
	ParamPageSize = "pageSize" //分页，每页长度，整形
	ParamTradeDate = "tradeDate" //交易日, 形如2020-01-01 文本
	ParamBeginDate = "beginDate" //查询开始时间, 文本，形如20200101
	ParamEndDate   = "endDate"   //查询结束时间, 文本，形如20200101
	ParamYtm = "ytm" // 债券到期收益率, 形如0.05，浮点数
)