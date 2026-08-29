package ficc_fut

import (
	"log"
	"rhino-common/domain_error"
	"rhino-common/utils/kafka"
	"rhino-core/domain_cfg"
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
	"strings"
)

type FiccFutScheduleAdapter struct {
	applicationCfg *domain_cfg.ApplicationCfg
}

func NewFiccFutScheduleAdapter(applicationCfg *domain_cfg.ApplicationCfg) (adapter *FiccFutScheduleAdapter, de *domain_error.Error) {
	log.Printf("construct FiccFutScheduleAdapter...")
	adapter = &FiccFutScheduleAdapter{applicationCfg: applicationCfg}
	return
}

func (a *FiccFutScheduleAdapter) ExecCustomizedPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog, positionManager *order_position_manager.PositionManager) (err error) {
	log.Printf("ExecCustomizedPurgingTask...")

	if purgingLog.TaskName == "dms" {
		// 清理境内交易台kafka
		tradeChannels := a.applicationCfg.GetTradeChannels()

		for _, tradeChannel := range tradeChannels {

			if tradeChannel.TradeChannel.ChannelCode != "olts-fut" {
				continue
			}

			// 清理交易通道的kafka
			tradeChannelDetailsList := a.applicationCfg.GetTradeChannels()
			for _, tradeChannelDetails := range tradeChannelDetailsList {
				chCfgMap := tradeChannelDetails.GetTradeChannelCfgItemMap()
				brokerCfgItem := chCfgMap["KafkaBroker"]
				requestTopicCfgItem := chCfgMap["RequestTopic"]
				responseTopicCfgItem := chCfgMap["ResponseTopic"]
				if brokerCfgItem != nil && requestTopicCfgItem != nil && responseTopicCfgItem != nil {
					brokers := getBrokers(brokerCfgItem.ConfigItemValue)
					err = kafka.PurgeMessageForTopics(brokers, []string{strings.TrimSpace(requestTopicCfgItem.ConfigItemValue), strings.TrimSpace(responseTopicCfgItem.ConfigItemValue)})
					if err != nil {
						return
					}
				}
				log.Printf("success purging star trading topic for channel %s\n", tradeChannelDetails.TradeChannel.ChannelCode)
			}
		}
	}

	positionManager.PurgePositionByPurgingTask(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, purgingLog)

	// 重新加载持仓
	positionManager.ReloadPositionRecordsForPurgingTask(purgingLog)

	// // Todo，发起一次数据上场，并且确认状态是正常的
	// // 清理持仓
	// if positionCalculator != nil {
	// 	positionCalculator.Reset()
	// }
	// // 执行数据上场的数据加载
	// if a.applicationCfg.GetAutoSyncRepo() != nil {
	// 	a.applicationCfg.GetAutoSyncRepo().Reset()
	// }

	return
}

func getBrokers(brokerStr string) (brokers []string) {
	strs := strings.Split(brokerStr, ",")
	for _, str := range strs {
		str = strings.TrimSpace(str)
		if str != "" {
			brokers = append(brokers, str)
		}
	}
	return
}
