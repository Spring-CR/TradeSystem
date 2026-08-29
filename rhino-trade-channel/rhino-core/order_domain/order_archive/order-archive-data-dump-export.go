package order_archive

import (
	"database/sql"
	"rhino-common/domain_error"
	"rhino-core/schema"
)

func (a *OrderArchiver) DumpTradeOrders(db *sql.DB, tableName string, records []*schema.TradeOrder) (de *domain_error.Error) {
	return a.dumpTradeOrders(db, tableName, records, nil)
}

func (a *OrderArchiver) DumpTradeActionLatestResps(db *sql.DB, tableName string, records []*schema.TradeActionLatestResp) (de *domain_error.Error) {
	return a.dumpTradeActionLatestResps(db, tableName, records, nil)
}

func (a *OrderArchiver) DumpTradeActionResps(db *sql.DB, tableName string, records []*schema.TradeActionResp) (de *domain_error.Error) {
	return a.dumpTradeActionResps(db, tableName, records, nil)
}