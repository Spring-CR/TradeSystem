package order_purge

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"rhino-common/domain_error"
	"rhino-common/enum"
	"rhino-common/utils/kafka"
	"rhino-core/schema"
	"rhino-core/store/admin_store"
)

// 第4步：分析内存订单数据，拆分待清除和继续保留的数据
func (a *OrderPurger) resetKafka(purgingLog *schema.DataPurgingLog) (err error) {

	if purgingLog.CompletePhase >= int(enum.DataPurgingLogPhase_ResetKafka) {
		return
	}

	systemCode, businessCode := a.applicationCfg.GetSystemAndBusinessCodes()

	topics := []string{
		fmt.Sprintf("%s-%s-req", systemCode, businessCode),
		fmt.Sprintf("%s-%s-resp", systemCode, businessCode),
		fmt.Sprintf("%s-%s-channel-req", systemCode, businessCode),
		fmt.Sprintf("%s-%s-order-status", systemCode, businessCode),
		fmt.Sprintf("%s-%s-order-status-ack", systemCode, businessCode),
	}

	for i := 0; i < 3; i++ {
		err = kafka.PurgeMessageForTopics(a.applicationCfg.GetKafkaBrokers(), topics)
		if err == nil {
			break
		}
	}

	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to purge kafka message")
		return
	}

	purgingLog.CompletePhase = int(enum.DataPurgingLogPhase_ResetKafka)
	err = admin_store.UpdateDataPurgingLogById(a.applicationCfg.GetCentralDB(), purgingLog)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "error occurs in OrderPurger::resetKafka")
	}

	// truncate缓存文件
	workDir := a.applicationCfg.GetWorkingDir()
	// 清空状态同步文件
	err = os.Truncate(filepath.Join(workDir, "order_cace", "order_status_msg.log"), 0)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "fail to truncate order_status_msg.log")
	} else {
		log.Println("success truncate order_status_msg.log")
	}
	// 清空成交回报文件
	os.Truncate(filepath.Join(workDir, "stream_api", "trade_resp.log"), 0)
	if err != nil {
		domain_error.ProcessSevereError(false, 0, nil, err, "trade_resp.log")
	} else {
		log.Println("success truncate trade_resp.log")
	}

	return
}
