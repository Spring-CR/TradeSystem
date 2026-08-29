package domain_cfg

import (
	"bytes"
	"rhino-common/domain_error"
	"rhino-core/schema"
)

/*
*
作为FIX Initiator 所需要的配置项
*/
var FixInitiatorCfgItemCheckers = []*CfgItemChecker{
	{
		ConfigItemName:         "BeginString",
		ConfigItemDefaultValue: "",
		Required:               1,
		NoOverWrite:            0,
		Description:            "此会话应使用的FIX版本",
	}, {
		ConfigItemName:         "SenderCompID",
		ConfigItemDefaultValue: "",
		Required:               1,
		NoOverWrite:            0,
		Description:            "发送方公司ID，用于标识发送方",
	}, {
		ConfigItemName:         "TargetCompID",
		ConfigItemDefaultValue: "",
		Required:               1,
		NoOverWrite:            0,
		Description:            "接收方公司ID，用于标识接收方",
	}, {
		ConfigItemName:         "ResetOnLogon",
		ConfigItemDefaultValue: "N",
		Required:               1,
		NoOverWrite:            1,
		Description:            "是否在登录后重置消息序号",
	}, {
		ConfigItemName:         "TimeStampPrecision",
		ConfigItemDefaultValue: "MILLIS",
		Required:               1,
		NoOverWrite:            1,
		Description:            "时间戳精度，强制为毫秒",
	}, {
		ConfigItemName:         "RejectInvalidMessage",
		ConfigItemDefaultValue: "N",
		Required:               1,
		NoOverWrite:            1,
		Description:            "系统在接收到不符合数据字典验证的消息时不会拒绝该消息（允许系统处理未完全符合数据字典的消息，特别是与不遵循严格规范的对方系统进行交互时；可以避免因小的格式错误或不一致而导致的消息丢失，从而保持交易的连续性）",
	}, {
		ConfigItemName:         "AllowUnknownMessageFields",
		ConfigItemDefaultValue: "Y",
		Required:               1,
		NoOverWrite:            1,
		Description:            "系统将允许消息中包含未在数据字典中定义的字段（标签小于5000）",
	}, {
		ConfigItemName:         "CheckUserDefinedFields",
		ConfigItemDefaultValue: "N",
		Required:               1,
		NoOverWrite:            1,
		Description:            "系统在接收消息时不会验证用户自定义字段（即字段标签大于或等于5000）的有效性",
	}, {
		ConfigItemName:         "CheckLatency",
		ConfigItemDefaultValue: "N",
		Required:               1,
		NoOverWrite:            1,
		Description:            "系统在接收消息时不会检查消息的延迟",
	}, /*{
		ConfigItemName:         "MaxLatency",
		ConfigItemDefaultValue: "240",
		Required:               1,
		NoOverWrite:            0,
		Description:            "允许的最大消息延迟，默认是240秒",
	},*/ {
		ConfigItemName:         "ReconnectInterval",
		ConfigItemDefaultValue: "5",
		Required:               0,
		NoOverWrite:            0,
		Description:            "定义重新连接尝试之间的时间间隔（秒），如果Initiator与Acceptor失去连接，系统将在此时间间隔后尝试重新连接",
	}, {
		ConfigItemName:         "LogoutTimeout",
		ConfigItemDefaultValue: "10",
		Required:               1,
		NoOverWrite:            1,
		Description:            "定义等待注销响应的时间（秒），在此时间内未收到响应将断开连接",
	}, {
		ConfigItemName:         "LogonTimeout",
		ConfigItemDefaultValue: "10",
		Required:               1,
		NoOverWrite:            1,
		Description:            "定义等待登录响应的时间（秒），在此时间内未收到响应将断开连接",
	}, {
		ConfigItemName:         "HeartBtInt",
		ConfigItemDefaultValue: "5",
		Required:               0,
		NoOverWrite:            0,
		Description:            "这是Initiator所需的配置，心跳机制用于检测连接是否仍然有效",
	}, {
		ConfigItemName:         "SocketConnectHost",
		ConfigItemDefaultValue: "<根据TradeChannel.Address来设置>",
		Required:               1,
		NoOverWrite:            1,
		SetByProgram:           1,
		Description:            "设置尝试连接的主机地址，可以指定多个主机以实现故障转移，使用形式如 SocketConnectHost1, SocketConnectHost2 等",
	}, {
		ConfigItemName:         "SocketConnectPort",
		ConfigItemDefaultValue: "<根据TradeChannel.Address来设置>",
		Required:               1,
		NoOverWrite:            1,
		SetByProgram:           1,
		Description:            "设置连接会话的套接字端口，与 SocketConnectHost 配合使用，确保故障转移时的端口匹配",
	}, {
		ConfigItemName:         "FileLogPath",
		ConfigItemDefaultValue: "/opt/rhino-trade-channel/log",
		Required:               0,
		NoOverWrite:            0,
		Description:            "quickfix的日志目录",
	},  {
		ConfigItemName:         "SocketUseSSL",
		ConfigItemDefaultValue: "N",
		Required:               0,
		NoOverWrite:            0,
		Description:            "是否启用SSL",
	},
}

type TradeChannelCfgAdapterForFixInitiator struct {
	tradeChannel         *schema.TradeChannel
	tradeChannelCfgItems []*schema.TradeChannelCfgItem
}

func newTradeChannelCfgAdapterForFixInitiator(tradeChannel *schema.TradeChannel, tradeChannelCfgItems []*schema.TradeChannelCfgItem) (cfg *TradeChannelCfgAdapterForFixInitiator) {
	cfg = &TradeChannelCfgAdapterForFixInitiator{tradeChannel: tradeChannel, tradeChannelCfgItems: tradeChannelCfgItems}
	return cfg
}

func (cfg *TradeChannelCfgAdapterForFixInitiator) ToAppConfig() (configFileContent []byte, de *domain_error.Error) {
	var itemValues map[string]string
	itemValues, de = checkCfgItems(cfg.tradeChannelCfgItems, FixInitiatorCfgItemCheckers)
	if de != nil {
		return
	}

	strBuf := bytes.NewBuffer([]byte("[DEFAULT]\n"))
	strBuf.Write([]byte("FileLogPath=" + itemValues["FileLogPath"] + "\n"))

	strBuf.Write([]byte("[SESSION]\n"))
	hostAndPortPaires := parseHostAndPortPaires(cfg.tradeChannel.Addresses)
	hostAndPortLen := len(hostAndPortPaires)
	for _, checker := range FixInitiatorCfgItemCheckers {
		configItemName := checker.ConfigItemName
		if configItemName == "FileLogPath" {
			continue
		}
		strBuf.WriteString(configItemName)
		strBuf.WriteString("=")
		if checker.SetByProgram > 0 {
			if configItemName == "SocketConnectHost" {
				for i, hostAndPort := range hostAndPortPaires {
					strBuf.WriteString(hostAndPort.Host)
					if i < hostAndPortLen-1 {
						strBuf.WriteString(",")
					}
				}
			}
			if configItemName == "SocketConnectPort" {
				for i, hostAndPort := range hostAndPortPaires {
					strBuf.WriteString(hostAndPort.Port)
					if i < hostAndPortLen-1 {
						strBuf.WriteString(",")
					}
				}

			}
		} else {
			strBuf.WriteString(itemValues[configItemName])
		}
		strBuf.WriteString("\n")
	}
	configFileContent = strBuf.Bytes()
	return
}
