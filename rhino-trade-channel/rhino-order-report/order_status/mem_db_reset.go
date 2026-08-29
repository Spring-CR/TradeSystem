package order_status

import (
	"log"
	"rhino-common/domain_error"
	"rhino-core/schema"
)

func (s *OrderStatusReplica) afterReset(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, tradeActionRespsToArchive []*schema.TradeActionResp, purgingLog *schema.DataPurgingLog) {

	log.Printf("start to reset OrderStatusReplica, orderCache:%v\n", s.orderCache)

	if s.dbWrite != nil {
		err := s.dbWrite.Close()
		if err != nil {
			domain_error.ProcessSevereError(false, 0, nil, err, "关闭内存数据库dbWrite连接失败: ")
		}
	}

	// 当前s.dbRead 指向了 s.dbWrite
	// if s.dbRead != nil {
	// 	err := s.dbRead.Close()
	// 	if err != nil {
	// 		domain_error.ProcessSevereError(false, 0, nil, err, "关闭内存数据库dbRead连接失败: ")
	// 	}
	// }

	log.Printf("start to initMemDB in reset function, orderCache=%v\n", s.orderCache)
	s.initMemDb()
	s.memInitByOrderTopology()

	// 清空然后重载持仓信息
	if s.positionManager != nil {
		s.positionManager.PurgePositionByPurgingTask(tradeOrdersToKeep, tradeActionLatestRespsToKeep, tradeActionRespsToKeep, tradeOrdersToArchive, tradeActionLatestRespsToArchive, tradeActionRespsToArchive, purgingLog)
		// 重新加载持仓
		s.positionManager.ReloadPositionRecordsForPurgingTask(purgingLog)
	}
}
