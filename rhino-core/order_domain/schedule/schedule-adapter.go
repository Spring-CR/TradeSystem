package schedule

import (
	"rhino-core/order_domain/order_position_manager"
	"rhino-core/schema"
)

type ScheduleAdapter interface {
	ExecCustomizedPurgingTask(tradeOrdersToKeep []*schema.TradeOrder, tradeActionLatestRespsToKeep []*schema.TradeActionLatestResp, tradeActionRespsToKeep []*schema.TradeActionResp, tradeOrdersToArchive []*schema.TradeOrder, tradeActionLatestRespsToArchive []*schema.TradeActionLatestResp, purgingLog *schema.DataPurgingLog, positionManager *order_position_manager.PositionManager) (err error)
}
