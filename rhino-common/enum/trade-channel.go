package enum

// 作者：林春泉
// 参考：http://wiki.gf.com.cn/pages/viewpage.action?pageId=259950569

// 交易通道的协议类型
type ChannelProtocolType string

const (
	ChannelProtocolType_O32     ChannelProtocolType = "O32"
	ChannelProtocolType_O45     ChannelProtocolType = "O45"
	ChannelProtocolType_UF20    ChannelProtocolType = "UF20"
	ChannelProtocolType_ATP     ChannelProtocolType = "UF20"
	ChannelProtocolType_ATP_FUT ChannelProtocolType = "ATP_FUT"
	ChannelProtocolType_ATP_OPT ChannelProtocolType = "ATP_OPT"
	ChannelProtocolType_CTP     ChannelProtocolType = "CTP"
	ChannelProtocolType_ESUNNY  ChannelProtocolType = "ESUNNY"
	ChannelProtocolType_FIX42   ChannelProtocolType = "FIX42"
	ChannelProtocolType_FIX44   ChannelProtocolType = "FIX44"
	ChannelProtocolType_STARS   ChannelProtocolType = "stars"
)

// 交易通道的状态
type ChannelStatus int

const (
	ChannelStatus_OFFLINE    ChannelStatus = 0
	ChannelStatus_SUBHEALTHY ChannelStatus = 1
	ChannelStatus_HEALTHY    ChannelStatus = 2
)
