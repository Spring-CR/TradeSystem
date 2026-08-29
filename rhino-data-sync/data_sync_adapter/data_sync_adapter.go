package data_sync_adapter

import "rhino-common/utils/dbutil"

type DataSyncAdapter interface {
	RefineCsvContent(tableConfig *dbutil.TableConfig, rawCsv []byte) (newCsv []byte, err error)
}
