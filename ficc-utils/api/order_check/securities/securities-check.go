package securities

import (
	"ficc-utils/common/utils/data_qry"
	"ficc-utils/common/utils/wechat"
	"log"
)

const (
	CheckSymbolsT0ResultHead = "【T0券表记录缺失提示】\n"
)

func Task_CheckSymbolsT0(dataQryUrl, WebhookUrl string) error {
	passed, checkInfo, err := checkSymbolsT0(dataQryUrl)
	if err != nil {
		log.Printf("Run Task Task_CheckSymbolsT0 error: %v", err)
		return err
	}

	if !passed {
		msg := CheckSymbolsT0ResultHead + checkInfo
		log.Printf("Run Task Task_CheckSymbolsT0 not passed:\n%s", msg)
		err = wechat.SendToWeChat(WebhookUrl, msg)
		if err != nil {
			log.Printf("Run Task Task_CheckSymbolsT0 SendToWeChat error: %v", err)
			return err
		}
		return nil
	}
	log.Println("Run Task Task_CheckSymbolsT0 check passed!")
	return nil
}

func checkSymbolsT0(dataQryUrl string) (passed bool, checkInfo string, err error) {
	symbolsT0, err := data_qry.GetSymbolsT0(dataQryUrl)
	if err != nil {
		return false, "", err
	}

	if len(symbolsT0) == 0 {
		checkInfo += "T0券表记录数 = 0"
		return false, checkInfo, nil
	}

	log.Printf("checkSymbolsT0 T0券表记录数: %d", len(symbolsT0))
	return true, "", nil
}