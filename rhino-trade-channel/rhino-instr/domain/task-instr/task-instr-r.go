package domain_task_instr

import (
	"rhino-common/context"
	"rhino-common/domain_error"
	"rhino-common/utils/dbutil"
	"rhino-instr/schema"
	"rhino-instr/store"
)

func FindTaskInstrs(fieldConditions []*dbutil.FieldCondition, limit, offset int) (result []*schema.TaskInstrView, total int, de *domain_error.Error) {
	var err error
	result, err = store.FindTaskInstrs(context.DB, fieldConditions, limit, offset)
	if err != nil && dbutil.IsDbRecordEmptyError(err){
		err = nil
	}

	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_FIND_INSTR_RECORD_ERR_CODE, err)
		return
	}

	total, err = store.FindTaskInstrsCount(context.DB, fieldConditions)
	if err != nil {
		de = domain_error.Build(domain_error.CANNOT_FIND_INSTR_RECORD_ERR_CODE, err)
		return
	}

	return
}