package fix

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-trade-channel/adapter/config"
	"rhino-trade-channel/adapter/store/fix_store"
	"strings"

	"github.com/quickfixgo/quickfix"
)

func (c *GenericFIXChannel) Reset(force bool) (de *domain_error.Error) {

	log.Println("Reset GenericFIXChannel...")
	
	c.initiator.Stop(force)

	log.Println("original initiator stop!")

	var configFileContent []byte
	configFileContent, de = c.cfg.GetTradeChannelCfgAdapter().ToAppConfig()
	if de != nil {
		return
	}
	// 将ResetOnLogon 设置为Y
	configFileContent = bytes.ReplaceAll(configFileContent, []byte("ResetOnLogon=N"), []byte("ResetOnLogon=Y"))

	// step1: 生成fix客户端的配置文件（为了调试而输出文件，实际上启动fix客户端只需要直接读取配置内容即可）
	configDir := strings.TrimSpace(c.cfg.GetTradeChannel().ConfigDir)
	if configDir == "" {
		configDir = config.DefaultTradeChannelConfigDir
	}
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}
	configPath := filepath.Join(configDir, "tradeclient.cfg")
	err = os.WriteFile(configPath, configFileContent, 0644)
	if err != nil {
		de = domain_error.Build(domain_error.GENERIC_ERR_CODE, err)
		return
	}

	// step2: 启动fix客户端
	appSettings, err := quickfix.ParseSettings(bytes.NewReader(configFileContent))
	if err != nil {
		de = domain_error.Build(domain_error.FIX_PARSE_CLI_CFG_ERR_CODE, err)
		return
	}
	fileLogFactory, err := quickfix.NewFileLogFactory(appSettings)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_INIT_LOG_ERR_CODE, err)
		return
	}

	// step3: 创建storeFactory，重置msgSeqNum
	c.msgSeqGen.Reset()
	storeFactory := fix_store.NewAdvanceMemoryStoreFactory(c.msgSeqGen)

	// step4: 启动iniator
	c.initiator, err = quickfix.NewInitiator(c.tradeClient, storeFactory, appSettings, fileLogFactory)
	if err != nil {
		de = domain_error.Build(domain_error.FIX_INIT_CLI_ERR_CODE, err)
		return
	}
	err = c.initiator.Start()
	if err != nil {
		de = domain_error.Build(domain_error.FIX_START_CLI_ERR_CODE, err)
		return
	}

	log.Println("Reset GenericFIXChannel finish!")

	return
}