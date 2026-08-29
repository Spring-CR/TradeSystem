package admin_store

import (
	"rhino-core/schema"

	"github.com/linchunquan/sqlgen/db"
)

var (
	SelectDataSyncLogByIdForUpdateStmt = SelectDataSyncLogByIdStmt + ` for update`
)

func SelectDataSyncLogByIdForUpdate(db db.SimpleDB, id int64) (*schema.DataSyncLog, error) {
	args := []interface{}{id}
	v, err := genericSelectDataSyncLog(db, SelectDataSyncLogByIdForUpdateStmt, args...)
	return v, err
}
