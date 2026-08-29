package dbutil

import (
	"database/sql"
	"fmt"
	"rhino-common/domain_error"
	"rhino-common/utils/emailalert"
	"strings"
	"time"
)

/*
"github.com/acsellers/inflections"

*/

func IsDbRecordEmptyError(err error) bool {
	return sql.ErrNoRows == err
}

func RollbackTx(tx *sql.Tx) *domain_error.Error {
	err := tx.Rollback()
	for i := 0; i < 5 && err != nil; i++ {
		time.Sleep(3 * time.Second)
		err = tx.Rollback()
	}
	if err != nil {
		emailalert.Send("rhino - RollbackTx with err", fmt.Sprintf("%+v", err))
		return domain_error.Build(domain_error.DATABASE_ROLLBACK_TRANS_ERR_CODE, err)
	}

	return nil
}

func CommitTx(tx *sql.Tx) *domain_error.Error {
	err := tx.Commit()
	if err != nil {
		de := RollbackTx(tx)
		if de != nil {
			return domain_error.Build(domain_error.DATABASE_ROLLBACK_TRANS_COMMIT_FAIL_ERR_CODE, de.Err)
		}
		return domain_error.Build(domain_error.DATABASE_COMMIT_TRANS_ERR_CODE, err, err)
	}
	return nil
}

func BeginTx(db *sql.DB) (tx *sql.Tx, de *domain_error.Error) {
	var err error
	tx, err = db.Begin()
	if err != nil {
		return nil, domain_error.Build(domain_error.DATABASE_OPEN_TRANS_ERR_CODE, err)
	}
	return tx, de
}

func IsMysqlDuplicateEntryError(err error) bool {
	if err != nil && strings.Contains(err.Error(), "Duplicate entry"){
		return true
	}
	return false
}