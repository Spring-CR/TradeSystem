package store

// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.

import (
	"database/sql"
	"github.com/linchunquan/sqlgen/db"
	"rhino-instr/schema"
)

const CreateTaskInstrStmt = `
CREATE TABLE IF NOT EXISTS task_instrs (
 f_id                     BIGINT PRIMARY KEY AUTO_INCREMENT
,f_date                   INTEGER
,f_daily_instr_no         BIGINT
,f_index_daily_modify     BIGINT
,f_batch_serial_no        BIGINT
,f_index_last_modify      BIGINT
,f_account_no             VARCHAR(16)
,f_combi_no               VARCHAR(16)
,f_instr_type             VARCHAR(1)
,f_begin_date             INTEGER
,f_end_date               INTEGER
,f_begin_time             INTEGER
,f_end_time               INTEGER
,f_direct_date            INTEGER
,f_direct_time            INTEGER
,f_direct_operator        VARCHAR(32)
,f_modify_date            INTEGER
,f_modify_time            INTEGER
,f_modify_operator        VARCHAR(32)
,f_modify_reason          VARCHAR(128)
,f_dispense_date          INTEGER
,f_dispense_time          INTEGER
,f_dispense_operator      VARCHAR(32)
,f_dispense_refuse_reason VARCHAR(128)
,f_cancel_date            INTEGER
,f_cancel_time            INTEGER
,f_cancel_operator        VARCHAR(32)
,f_cancel_reason          VARCHAR(128)
,f_operator               VARCHAR(32)
,f_instr_status           VARCHAR(1)
,f_dispense_status        VARCHAR(1)
,f_entrust_execute_status VARCHAR(1)
,f_deal_execute_status    VARCHAR(1)
,f_create_time            BIGINT
,f_business_type          VARCHAR(2)
,f_lock_flag              INTEGER
,f_limit_operator         VARCHAR(32)
,f_org_id                 BIGINT
,f_dept_id                BIGINT
,f_ip_address             VARCHAR(16)
,f_mac                    VARCHAR(20)
,f_volserial_no           VARCHAR(10)
);
`

const InsertTaskInstrStmt = `
INSERT INTO task_instrs (
 f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTaskInstrStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
`

const SelectTaskInstrRangeStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
LIMIT ? OFFSET ?
`

const SelectTaskInstrCountStmt = `
SELECT count(1)
FROM task_instrs 
`

const SelectTaskInstrByIdStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_id=?
`

const UpdateTaskInstrByIdStmt = `
UPDATE task_instrs SET 
 f_id=?
,f_date=?
,f_daily_instr_no=?
,f_index_daily_modify=?
,f_batch_serial_no=?
,f_index_last_modify=?
,f_account_no=?
,f_combi_no=?
,f_instr_type=?
,f_begin_date=?
,f_end_date=?
,f_begin_time=?
,f_end_time=?
,f_direct_date=?
,f_direct_time=?
,f_direct_operator=?
,f_modify_date=?
,f_modify_time=?
,f_modify_operator=?
,f_modify_reason=?
,f_dispense_date=?
,f_dispense_time=?
,f_dispense_operator=?
,f_dispense_refuse_reason=?
,f_cancel_date=?
,f_cancel_time=?
,f_cancel_operator=?
,f_cancel_reason=?
,f_operator=?
,f_instr_status=?
,f_dispense_status=?
,f_entrust_execute_status=?
,f_deal_execute_status=?
,f_create_time=?
,f_business_type=?
,f_lock_flag=?
,f_limit_operator=?
,f_org_id=?
,f_dept_id=?
,f_ip_address=?
,f_mac=?
,f_volserial_no=? 
WHERE f_id=?
`

const DeleteTaskInstrByIdStmt = `
DELETE FROM task_instrs 
WHERE f_id=?
`

const CreatePkTiStmt = `
CREATE UNIQUE INDEX pk_ti ON task_instrs (f_date,f_daily_instr_no,f_index_daily_modify);
`

const SelectTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const SelectTaskInstrCountByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT count(1)
FROM task_instrs 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const UpdateTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
UPDATE task_instrs SET 
 f_id=?
,f_date=?
,f_daily_instr_no=?
,f_index_daily_modify=?
,f_batch_serial_no=?
,f_index_last_modify=?
,f_account_no=?
,f_combi_no=?
,f_instr_type=?
,f_begin_date=?
,f_end_date=?
,f_begin_time=?
,f_end_time=?
,f_direct_date=?
,f_direct_time=?
,f_direct_operator=?
,f_modify_date=?
,f_modify_time=?
,f_modify_operator=?
,f_modify_reason=?
,f_dispense_date=?
,f_dispense_time=?
,f_dispense_operator=?
,f_dispense_refuse_reason=?
,f_cancel_date=?
,f_cancel_time=?
,f_cancel_operator=?
,f_cancel_reason=?
,f_operator=?
,f_instr_status=?
,f_dispense_status=?
,f_entrust_execute_status=?
,f_deal_execute_status=?
,f_create_time=?
,f_business_type=?
,f_lock_flag=?
,f_limit_operator=?
,f_org_id=?
,f_dept_id=?
,f_ip_address=?
,f_mac=?
,f_volserial_no=? 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const DeleteTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
DELETE FROM task_instrs 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const CreateITiDateStmt = `
CREATE INDEX i_ti_date ON task_instrs (f_date);
`

const SelectTaskInstrByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_date=?
`

const SelectTaskInstrCountByDateStmt = `
SELECT count(1)
FROM task_instrs 
WHERE f_date=?
`

const SelectTaskInstrRangeByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_date=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrByDateStmt = `
DELETE FROM task_instrs 
WHERE f_date=?
`

const CreateITiDirectOperatorStmt = `
CREATE INDEX i_ti_direct_operator ON task_instrs (f_direct_operator);
`

const SelectTaskInstrByDirectOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_direct_operator=?
`

const SelectTaskInstrCountByDirectOperatorStmt = `
SELECT count(1)
FROM task_instrs 
WHERE f_direct_operator=?
`

const SelectTaskInstrRangeByDirectOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_direct_operator=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrByDirectOperatorStmt = `
DELETE FROM task_instrs 
WHERE f_direct_operator=?
`

const CreateITiOperatorStmt = `
CREATE INDEX i_ti_operator ON task_instrs (f_operator);
`

const SelectTaskInstrByOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_operator=?
`

const SelectTaskInstrCountByOperatorStmt = `
SELECT count(1)
FROM task_instrs 
WHERE f_operator=?
`

const SelectTaskInstrRangeByOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
FROM task_instrs 
WHERE f_operator=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrByOperatorStmt = `
DELETE FROM task_instrs 
WHERE f_operator=?
`

func scanTaskInstr(row *sql.Row) (*schema.TaskInstr, error) {
	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
	var v13 sql.NullInt64
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullInt64
	var v17 sql.NullInt64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullInt64
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullInt64
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullString
	var v29 sql.NullString
	var v30 sql.NullString
	var v31 sql.NullString
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullString
	var v40 sql.NullString
	var v41 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
		&v7,
		&v8,
		&v9,
		&v10,
		&v11,
		&v12,
		&v13,
		&v14,
		&v15,
		&v16,
		&v17,
		&v18,
		&v19,
		&v20,
		&v21,
		&v22,
		&v23,
		&v24,
		&v25,
		&v26,
		&v27,
		&v28,
		&v29,
		&v30,
		&v31,
		&v32,
		&v33,
		&v34,
		&v35,
		&v36,
		&v37,
		&v38,
		&v39,
		&v40,
		&v41,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TaskInstr{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.Date = int(v1.Int64)
	} else {
		v.Date = 0
	}

	if v2.Valid {
		v.DailyInstrNo = v2.Int64
	} else {
		v.DailyInstrNo = 0
	}

	if v3.Valid {
		v.IndexDailyModify = v3.Int64
	} else {
		v.IndexDailyModify = 0
	}

	if v4.Valid {
		v.BatchSerialNo = v4.Int64
	} else {
		v.BatchSerialNo = 0
	}

	if v5.Valid {
		v.IndexLastModify = v5.Int64
	} else {
		v.IndexLastModify = 0
	}

	if v6.Valid {
		v.AccountNo = v6.String
	} else {
		v.AccountNo = ""
	}

	if v7.Valid {
		v.CombiNo = v7.String
	} else {
		v.CombiNo = ""
	}

	if v8.Valid {
		v.InstrType = v8.String
	} else {
		v.InstrType = ""
	}

	if v9.Valid {
		v.BeginDate = int(v9.Int64)
	} else {
		v.BeginDate = 0
	}

	if v10.Valid {
		v.EndDate = int(v10.Int64)
	} else {
		v.EndDate = 0
	}

	if v11.Valid {
		v.BeginTime = int(v11.Int64)
	} else {
		v.BeginTime = 0
	}

	if v12.Valid {
		v.EndTime = int(v12.Int64)
	} else {
		v.EndTime = 0
	}

	if v13.Valid {
		v.DirectDate = int(v13.Int64)
	} else {
		v.DirectDate = 0
	}

	if v14.Valid {
		v.DirectTime = int(v14.Int64)
	} else {
		v.DirectTime = 0
	}

	if v15.Valid {
		v.DirectOperator = v15.String
	} else {
		v.DirectOperator = ""
	}

	if v16.Valid {
		v.ModifyDate = int(v16.Int64)
	} else {
		v.ModifyDate = 0
	}

	if v17.Valid {
		v.ModifyTime = int(v17.Int64)
	} else {
		v.ModifyTime = 0
	}

	if v18.Valid {
		v.ModifyOperator = v18.String
	} else {
		v.ModifyOperator = ""
	}

	if v19.Valid {
		v.ModifyReason = v19.String
	} else {
		v.ModifyReason = ""
	}

	if v20.Valid {
		v.DispenseDate = int(v20.Int64)
	} else {
		v.DispenseDate = 0
	}

	if v21.Valid {
		v.DispenseTime = int(v21.Int64)
	} else {
		v.DispenseTime = 0
	}

	if v22.Valid {
		v.DispenseOperator = v22.String
	} else {
		v.DispenseOperator = ""
	}

	if v23.Valid {
		v.DispenseRefuseReason = v23.String
	} else {
		v.DispenseRefuseReason = ""
	}

	if v24.Valid {
		v.CancelDate = int(v24.Int64)
	} else {
		v.CancelDate = 0
	}

	if v25.Valid {
		v.CancelTime = int(v25.Int64)
	} else {
		v.CancelTime = 0
	}

	if v26.Valid {
		v.CancelOperator = v26.String
	} else {
		v.CancelOperator = ""
	}

	if v27.Valid {
		v.CancelReason = v27.String
	} else {
		v.CancelReason = ""
	}

	if v28.Valid {
		v.Operator = v28.String
	} else {
		v.Operator = ""
	}

	if v29.Valid {
		v.InstrStatus = v29.String
	} else {
		v.InstrStatus = ""
	}

	if v30.Valid {
		v.DispenseStatus = v30.String
	} else {
		v.DispenseStatus = ""
	}

	if v31.Valid {
		v.EntrustExecuteStatus = v31.String
	} else {
		v.EntrustExecuteStatus = ""
	}

	if v32.Valid {
		v.DealExecuteStatus = v32.String
	} else {
		v.DealExecuteStatus = ""
	}

	if v33.Valid {
		v.CreateTime = v33.Int64
	} else {
		v.CreateTime = 0
	}

	if v34.Valid {
		v.BusinessType = v34.String
	} else {
		v.BusinessType = ""
	}

	if v35.Valid {
		v.LockFlag = int(v35.Int64)
	} else {
		v.LockFlag = 0
	}

	if v36.Valid {
		v.LimitOperator = v36.String
	} else {
		v.LimitOperator = ""
	}

	if v37.Valid {
		v.OrgId = v37.Int64
	} else {
		v.OrgId = 0
	}

	if v38.Valid {
		v.DeptId = v38.Int64
	} else {
		v.DeptId = 0
	}

	if v39.Valid {
		v.IpAddress = v39.String
	} else {
		v.IpAddress = ""
	}

	if v40.Valid {
		v.Mac = v40.String
	} else {
		v.Mac = ""
	}

	if v41.Valid {
		v.VolserialNo = v41.String
	} else {
		v.VolserialNo = ""
	}

	return v, nil
}

func scanTaskInstrs(rows *sql.Rows) ([]*schema.TaskInstr, error) {
	var err error
	var vv []*schema.TaskInstr

	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
	var v13 sql.NullInt64
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullInt64
	var v17 sql.NullInt64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullInt64
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullInt64
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullString
	var v29 sql.NullString
	var v30 sql.NullString
	var v31 sql.NullString
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullString
	var v40 sql.NullString
	var v41 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
			&v10,
			&v11,
			&v12,
			&v13,
			&v14,
			&v15,
			&v16,
			&v17,
			&v18,
			&v19,
			&v20,
			&v21,
			&v22,
			&v23,
			&v24,
			&v25,
			&v26,
			&v27,
			&v28,
			&v29,
			&v30,
			&v31,
			&v32,
			&v33,
			&v34,
			&v35,
			&v36,
			&v37,
			&v38,
			&v39,
			&v40,
			&v41,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TaskInstr{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.Date = int(v1.Int64)
		} else {
			v.Date = 0
		}

		if v2.Valid {
			v.DailyInstrNo = v2.Int64
		} else {
			v.DailyInstrNo = 0
		}

		if v3.Valid {
			v.IndexDailyModify = v3.Int64
		} else {
			v.IndexDailyModify = 0
		}

		if v4.Valid {
			v.BatchSerialNo = v4.Int64
		} else {
			v.BatchSerialNo = 0
		}

		if v5.Valid {
			v.IndexLastModify = v5.Int64
		} else {
			v.IndexLastModify = 0
		}

		if v6.Valid {
			v.AccountNo = v6.String
		} else {
			v.AccountNo = ""
		}

		if v7.Valid {
			v.CombiNo = v7.String
		} else {
			v.CombiNo = ""
		}

		if v8.Valid {
			v.InstrType = v8.String
		} else {
			v.InstrType = ""
		}

		if v9.Valid {
			v.BeginDate = int(v9.Int64)
		} else {
			v.BeginDate = 0
		}

		if v10.Valid {
			v.EndDate = int(v10.Int64)
		} else {
			v.EndDate = 0
		}

		if v11.Valid {
			v.BeginTime = int(v11.Int64)
		} else {
			v.BeginTime = 0
		}

		if v12.Valid {
			v.EndTime = int(v12.Int64)
		} else {
			v.EndTime = 0
		}

		if v13.Valid {
			v.DirectDate = int(v13.Int64)
		} else {
			v.DirectDate = 0
		}

		if v14.Valid {
			v.DirectTime = int(v14.Int64)
		} else {
			v.DirectTime = 0
		}

		if v15.Valid {
			v.DirectOperator = v15.String
		} else {
			v.DirectOperator = ""
		}

		if v16.Valid {
			v.ModifyDate = int(v16.Int64)
		} else {
			v.ModifyDate = 0
		}

		if v17.Valid {
			v.ModifyTime = int(v17.Int64)
		} else {
			v.ModifyTime = 0
		}

		if v18.Valid {
			v.ModifyOperator = v18.String
		} else {
			v.ModifyOperator = ""
		}

		if v19.Valid {
			v.ModifyReason = v19.String
		} else {
			v.ModifyReason = ""
		}

		if v20.Valid {
			v.DispenseDate = int(v20.Int64)
		} else {
			v.DispenseDate = 0
		}

		if v21.Valid {
			v.DispenseTime = int(v21.Int64)
		} else {
			v.DispenseTime = 0
		}

		if v22.Valid {
			v.DispenseOperator = v22.String
		} else {
			v.DispenseOperator = ""
		}

		if v23.Valid {
			v.DispenseRefuseReason = v23.String
		} else {
			v.DispenseRefuseReason = ""
		}

		if v24.Valid {
			v.CancelDate = int(v24.Int64)
		} else {
			v.CancelDate = 0
		}

		if v25.Valid {
			v.CancelTime = int(v25.Int64)
		} else {
			v.CancelTime = 0
		}

		if v26.Valid {
			v.CancelOperator = v26.String
		} else {
			v.CancelOperator = ""
		}

		if v27.Valid {
			v.CancelReason = v27.String
		} else {
			v.CancelReason = ""
		}

		if v28.Valid {
			v.Operator = v28.String
		} else {
			v.Operator = ""
		}

		if v29.Valid {
			v.InstrStatus = v29.String
		} else {
			v.InstrStatus = ""
		}

		if v30.Valid {
			v.DispenseStatus = v30.String
		} else {
			v.DispenseStatus = ""
		}

		if v31.Valid {
			v.EntrustExecuteStatus = v31.String
		} else {
			v.EntrustExecuteStatus = ""
		}

		if v32.Valid {
			v.DealExecuteStatus = v32.String
		} else {
			v.DealExecuteStatus = ""
		}

		if v33.Valid {
			v.CreateTime = v33.Int64
		} else {
			v.CreateTime = 0
		}

		if v34.Valid {
			v.BusinessType = v34.String
		} else {
			v.BusinessType = ""
		}

		if v35.Valid {
			v.LockFlag = int(v35.Int64)
		} else {
			v.LockFlag = 0
		}

		if v36.Valid {
			v.LimitOperator = v36.String
		} else {
			v.LimitOperator = ""
		}

		if v37.Valid {
			v.OrgId = v37.Int64
		} else {
			v.OrgId = 0
		}

		if v38.Valid {
			v.DeptId = v38.Int64
		} else {
			v.DeptId = 0
		}

		if v39.Valid {
			v.IpAddress = v39.String
		} else {
			v.IpAddress = ""
		}

		if v40.Valid {
			v.Mac = v40.String
		} else {
			v.Mac = ""
		}

		if v41.Valid {
			v.VolserialNo = v41.String
		} else {
			v.VolserialNo = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTaskInstr(v *schema.TaskInstr) []interface{} {
	var v0 int64
	var v1 int
	var v2 int64
	var v3 int64
	var v4 int64
	var v5 int64
	var v6 string
	var v7 string
	var v8 string
	var v9 int
	var v10 int
	var v11 int
	var v12 int
	var v13 int
	var v14 int
	var v15 string
	var v16 int
	var v17 int
	var v18 string
	var v19 string
	var v20 int
	var v21 int
	var v22 string
	var v23 string
	var v24 int
	var v25 int
	var v26 string
	var v27 string
	var v28 string
	var v29 string
	var v30 string
	var v31 string
	var v32 string
	var v33 int64
	var v34 string
	var v35 int
	var v36 string
	var v37 int64
	var v38 int64
	var v39 string
	var v40 string
	var v41 string

	v0 = v.ID
	v1 = v.Date
	v2 = v.DailyInstrNo
	v3 = v.IndexDailyModify
	v4 = v.BatchSerialNo
	v5 = v.IndexLastModify
	v6 = v.AccountNo
	v7 = v.CombiNo
	v8 = v.InstrType
	v9 = v.BeginDate
	v10 = v.EndDate
	v11 = v.BeginTime
	v12 = v.EndTime
	v13 = v.DirectDate
	v14 = v.DirectTime
	v15 = v.DirectOperator
	v16 = v.ModifyDate
	v17 = v.ModifyTime
	v18 = v.ModifyOperator
	v19 = v.ModifyReason
	v20 = v.DispenseDate
	v21 = v.DispenseTime
	v22 = v.DispenseOperator
	v23 = v.DispenseRefuseReason
	v24 = v.CancelDate
	v25 = v.CancelTime
	v26 = v.CancelOperator
	v27 = v.CancelReason
	v28 = v.Operator
	v29 = v.InstrStatus
	v30 = v.DispenseStatus
	v31 = v.EntrustExecuteStatus
	v32 = v.DealExecuteStatus
	v33 = v.CreateTime
	v34 = v.BusinessType
	v35 = v.LockFlag
	v36 = v.LimitOperator
	v37 = v.OrgId
	v38 = v.DeptId
	v39 = v.IpAddress
	v40 = v.Mac
	v41 = v.VolserialNo

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
		v9,
		v10,
		v11,
		v12,
		v13,
		v14,
		v15,
		v16,
		v17,
		v18,
		v19,
		v20,
		v21,
		v22,
		v23,
		v24,
		v25,
		v26,
		v27,
		v28,
		v29,
		v30,
		v31,
		v32,
		v33,
		v34,
		v35,
		v36,
		v37,
		v38,
		v39,
		v40,
		v41,
	}
}

func genericSelectTaskInstr(db db.SimpleDB, query string, args ...interface{}) (*schema.TaskInstr, error) {
	row := db.QueryRow(query, args...)
	return scanTaskInstr(row)
}

func genericSelectTaskInstrs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TaskInstr, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTaskInstrs(rows)
}

func InsertTaskInstr(db db.SimpleDB, v *schema.TaskInstr) error {

	res, err := db.Exec(InsertTaskInstrStmt, sliceTaskInstr(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTaskInstrById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTaskInstrByIdStmt, args...)
	return err
}

func DeleteTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) error {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	_, err := db.Exec(DeleteTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return err
}

func DeleteTaskInstrByDate(db db.SimpleDB, date int) error {
	args := []interface{}{date}
	_, err := db.Exec(DeleteTaskInstrByDateStmt, args...)
	return err
}

func DeleteTaskInstrByDirectOperator(db db.SimpleDB, directOperator string) error {
	args := []interface{}{directOperator}
	_, err := db.Exec(DeleteTaskInstrByDirectOperatorStmt, args...)
	return err
}

func DeleteTaskInstrByOperator(db db.SimpleDB, operator string) error {
	args := []interface{}{operator}
	_, err := db.Exec(DeleteTaskInstrByOperatorStmt, args...)
	return err
}

func UpdateTaskInstrById(db db.SimpleDB, v *schema.TaskInstr) error {
	args := sliceTaskInstr(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTaskInstrByIdStmt, args...)
	return err
}

func UpdateTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, v *schema.TaskInstr) error {
	args := sliceTaskInstr(v)
	args = append(args, v.Date, v.DailyInstrNo, v.IndexDailyModify)
	_, err := db.Exec(UpdateTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return err
}

func GetTaskInstrById(db db.SimpleDB, iD int64) (*schema.TaskInstr, error) {
	args := []interface{}{iD}
	v, err := genericSelectTaskInstr(db, SelectTaskInstrByIdStmt, args...)
	return v, err
}

func GetTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) (*schema.TaskInstr, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	v, err := genericSelectTaskInstr(db, SelectTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return v, err
}

func FindAllTaskInstrs(db db.SimpleDB) ([]*schema.TaskInstr, error) {
	args := []interface{}{}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrStmt, args...)
	return v, err
}

func FindAllTaskInstrsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TaskInstr, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrRangeStmt, args...)
	return v, err
}

func FindTaskInstrsByDate(db db.SimpleDB, date int) ([]*schema.TaskInstr, error) {
	args := []interface{}{date}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrByDateStmt, args...)
	return v, err
}

func FindTaskInstrsByDateInRange(db db.SimpleDB, date int, limit int64, offset int64) ([]*schema.TaskInstr, error) {
	args := []interface{}{date, limit, offset}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrRangeByDateStmt, args...)
	return v, err
}

func FindTaskInstrsByDirectOperator(db db.SimpleDB, directOperator string) ([]*schema.TaskInstr, error) {
	args := []interface{}{directOperator}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrByDirectOperatorStmt, args...)
	return v, err
}

func FindTaskInstrsByDirectOperatorInRange(db db.SimpleDB, directOperator string, limit int64, offset int64) ([]*schema.TaskInstr, error) {
	args := []interface{}{directOperator, limit, offset}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrRangeByDirectOperatorStmt, args...)
	return v, err
}

func FindTaskInstrsByOperator(db db.SimpleDB, operator string) ([]*schema.TaskInstr, error) {
	args := []interface{}{operator}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrByOperatorStmt, args...)
	return v, err
}

func FindTaskInstrsByOperatorInRange(db db.SimpleDB, operator string, limit int64, offset int64) ([]*schema.TaskInstr, error) {
	args := []interface{}{operator, limit, offset}
	v, err := genericSelectTaskInstrs(db, SelectTaskInstrRangeByOperatorStmt, args...)
	return v, err
}

func CountTaskInstr(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTaskInstrCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) (int, error) {
	var count int
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	row := db.QueryRow(SelectTaskInstrCountByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrByDate(db db.SimpleDB, date int) (int, error) {
	var count int
	args := []interface{}{date}
	row := db.QueryRow(SelectTaskInstrCountByDateStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrByDirectOperator(db db.SimpleDB, directOperator string) (int, error) {
	var count int
	args := []interface{}{directOperator}
	row := db.QueryRow(SelectTaskInstrCountByDirectOperatorStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrByOperator(db db.SimpleDB, operator string) (int, error) {
	var count int
	args := []interface{}{operator}
	row := db.QueryRow(SelectTaskInstrCountByOperatorStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTaskInstrStockStmt = `
CREATE TABLE IF NOT EXISTS task_instr_stocks (
 f_id                           BIGINT PRIMARY KEY AUTO_INCREMENT
,f_date                         INTEGER
,f_daily_instr_no               BIGINT
,f_index_daily_modify           BIGINT
,f_stock_serial_no              BIGINT
,f_market_no                    VARCHAR(8)
,f_report_code                  VARCHAR(16)
,f_entrust_direction            VARCHAR(2)
,f_open_close                   VARCHAR(8)
,f_invest_type                  VARCHAR(2)
,f_amount                       DOUBLE
,f_balance                      DOUBLE
,f_contract_size                DOUBLE
,f_price                        DOUBLE
,f_stock_entrust_execute_status VARCHAR(1)
,f_stock_deal_execute_status    VARCHAR(1)
,f_total_deal_amount            DOUBLE
,f_total_deal_balance           DOUBLE
,f_cum_avg_price                DOUBLE
,f_total_entrust_amount         DOUBLE
,f_total_entrust_balance        DOUBLE
,f_deal_complete_date_time      BIGINT
,f_estimate_fee                 DOUBLE
,f_stock_instr_execution_time   BIGINT
,f_stock_instr_operator         VARCHAR(32)
);
`

const InsertTaskInstrStockStmt = `
INSERT INTO task_instr_stocks (
 f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTaskInstrStockStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
`

const SelectTaskInstrStockRangeStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
LIMIT ? OFFSET ?
`

const SelectTaskInstrStockCountStmt = `
SELECT count(1)
FROM task_instr_stocks 
`

const SelectTaskInstrStockByIdStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_id=?
`

const UpdateTaskInstrStockByIdStmt = `
UPDATE task_instr_stocks SET 
 f_id=?
,f_date=?
,f_daily_instr_no=?
,f_index_daily_modify=?
,f_stock_serial_no=?
,f_market_no=?
,f_report_code=?
,f_entrust_direction=?
,f_open_close=?
,f_invest_type=?
,f_amount=?
,f_balance=?
,f_contract_size=?
,f_price=?
,f_stock_entrust_execute_status=?
,f_stock_deal_execute_status=?
,f_total_deal_amount=?
,f_total_deal_balance=?
,f_cum_avg_price=?
,f_total_entrust_amount=?
,f_total_entrust_balance=?
,f_deal_complete_date_time=?
,f_estimate_fee=?
,f_stock_instr_execution_time=?
,f_stock_instr_operator=? 
WHERE f_id=?
`

const DeleteTaskInstrStockByIdStmt = `
DELETE FROM task_instr_stocks 
WHERE f_id=?
`

const CreateIFTisStmt = `
CREATE INDEX i_f_tis ON task_instr_stocks (f_date,f_daily_instr_no,f_index_daily_modify);
`

const SelectTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const SelectTaskInstrStockCountByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT count(1)
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const SelectTaskInstrStockRangeByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
DELETE FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

const CreateITisDateStmt = `
CREATE INDEX i_tis_date ON task_instr_stocks (f_date);
`

const SelectTaskInstrStockByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
`

const SelectTaskInstrStockCountByDateStmt = `
SELECT count(1)
FROM task_instr_stocks 
WHERE f_date=?
`

const SelectTaskInstrStockRangeByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrStockByDateStmt = `
DELETE FROM task_instr_stocks 
WHERE f_date=?
`

const CreatePkTisStmt = `
CREATE UNIQUE INDEX pk_tis ON task_instr_stocks (f_date,f_daily_instr_no,f_index_daily_modify,f_stock_serial_no);
`

const SelectTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const SelectTaskInstrStockCountByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
SELECT count(1)
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const UpdateTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
UPDATE task_instr_stocks SET 
 f_id=?
,f_date=?
,f_daily_instr_no=?
,f_index_daily_modify=?
,f_stock_serial_no=?
,f_market_no=?
,f_report_code=?
,f_entrust_direction=?
,f_open_close=?
,f_invest_type=?
,f_amount=?
,f_balance=?
,f_contract_size=?
,f_price=?
,f_stock_entrust_execute_status=?
,f_stock_deal_execute_status=?
,f_total_deal_amount=?
,f_total_deal_balance=?
,f_cum_avg_price=?
,f_total_entrust_amount=?
,f_total_entrust_balance=?
,f_deal_complete_date_time=?
,f_estimate_fee=?
,f_stock_instr_execution_time=?
,f_stock_instr_operator=? 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
DELETE FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const CreateFkTisToTiStmt = `
ALTER TABLE task_instr_stocks ADD FOREIGN KEY (f_date,f_daily_instr_no,f_index_daily_modify) REFERENCES task_instrs (f_date,f_daily_instr_no,f_index_daily_modify);
`

const SelectTaskInstrStockOfTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_stocks 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
`

func scanTaskInstrStock(row *sql.Row) (*schema.TaskInstrStock, error) {
	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullFloat64
	var v11 sql.NullFloat64
	var v12 sql.NullFloat64
	var v13 sql.NullFloat64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullFloat64
	var v17 sql.NullFloat64
	var v18 sql.NullFloat64
	var v19 sql.NullFloat64
	var v20 sql.NullFloat64
	var v21 sql.NullInt64
	var v22 sql.NullFloat64
	var v23 sql.NullInt64
	var v24 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
		&v7,
		&v8,
		&v9,
		&v10,
		&v11,
		&v12,
		&v13,
		&v14,
		&v15,
		&v16,
		&v17,
		&v18,
		&v19,
		&v20,
		&v21,
		&v22,
		&v23,
		&v24,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TaskInstrStock{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.Date = int(v1.Int64)
	} else {
		v.Date = 0
	}

	if v2.Valid {
		v.DailyInstrNo = v2.Int64
	} else {
		v.DailyInstrNo = 0
	}

	if v3.Valid {
		v.IndexDailyModify = v3.Int64
	} else {
		v.IndexDailyModify = 0
	}

	if v4.Valid {
		v.StockSerialNo = v4.Int64
	} else {
		v.StockSerialNo = 0
	}

	if v5.Valid {
		v.MarketNo = v5.String
	} else {
		v.MarketNo = ""
	}

	if v6.Valid {
		v.ReportCode = v6.String
	} else {
		v.ReportCode = ""
	}

	if v7.Valid {
		v.EntrustDirection = v7.String
	} else {
		v.EntrustDirection = ""
	}

	if v8.Valid {
		v.OpenClose = v8.String
	} else {
		v.OpenClose = ""
	}

	if v9.Valid {
		v.InvestType = v9.String
	} else {
		v.InvestType = ""
	}

	if v10.Valid {
		v.Amount = v10.Float64
	} else {
		v.Amount = 0
	}

	if v11.Valid {
		v.Balance = v11.Float64
	} else {
		v.Balance = 0
	}

	if v12.Valid {
		v.ContractSize = v12.Float64
	} else {
		v.ContractSize = 0
	}

	if v13.Valid {
		v.Price = v13.Float64
	} else {
		v.Price = 0
	}

	if v14.Valid {
		v.StockEntrustExecuteStatus = v14.String
	} else {
		v.StockEntrustExecuteStatus = ""
	}

	if v15.Valid {
		v.StockDealExecuteStatus = v15.String
	} else {
		v.StockDealExecuteStatus = ""
	}

	if v16.Valid {
		v.TotalDealAmount = v16.Float64
	} else {
		v.TotalDealAmount = 0
	}

	if v17.Valid {
		v.TotalDealBalance = v17.Float64
	} else {
		v.TotalDealBalance = 0
	}

	if v18.Valid {
		v.CumAvgPrice = v18.Float64
	} else {
		v.CumAvgPrice = 0
	}

	if v19.Valid {
		v.TotalEntrustAmount = v19.Float64
	} else {
		v.TotalEntrustAmount = 0
	}

	if v20.Valid {
		v.TotalEntrustBalance = v20.Float64
	} else {
		v.TotalEntrustBalance = 0
	}

	if v21.Valid {
		v.DealCompleteDateTime = v21.Int64
	} else {
		v.DealCompleteDateTime = 0
	}

	if v22.Valid {
		v.EstimateFee = v22.Float64
	} else {
		v.EstimateFee = 0
	}

	if v23.Valid {
		v.StockInstrExecutionTime = v23.Int64
	} else {
		v.StockInstrExecutionTime = 0
	}

	if v24.Valid {
		v.StockInstrOperator = v24.String
	} else {
		v.StockInstrOperator = ""
	}

	return v, nil
}

func scanTaskInstrStocks(rows *sql.Rows) ([]*schema.TaskInstrStock, error) {
	var err error
	var vv []*schema.TaskInstrStock

	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullFloat64
	var v11 sql.NullFloat64
	var v12 sql.NullFloat64
	var v13 sql.NullFloat64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullFloat64
	var v17 sql.NullFloat64
	var v18 sql.NullFloat64
	var v19 sql.NullFloat64
	var v20 sql.NullFloat64
	var v21 sql.NullInt64
	var v22 sql.NullFloat64
	var v23 sql.NullInt64
	var v24 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
			&v10,
			&v11,
			&v12,
			&v13,
			&v14,
			&v15,
			&v16,
			&v17,
			&v18,
			&v19,
			&v20,
			&v21,
			&v22,
			&v23,
			&v24,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TaskInstrStock{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.Date = int(v1.Int64)
		} else {
			v.Date = 0
		}

		if v2.Valid {
			v.DailyInstrNo = v2.Int64
		} else {
			v.DailyInstrNo = 0
		}

		if v3.Valid {
			v.IndexDailyModify = v3.Int64
		} else {
			v.IndexDailyModify = 0
		}

		if v4.Valid {
			v.StockSerialNo = v4.Int64
		} else {
			v.StockSerialNo = 0
		}

		if v5.Valid {
			v.MarketNo = v5.String
		} else {
			v.MarketNo = ""
		}

		if v6.Valid {
			v.ReportCode = v6.String
		} else {
			v.ReportCode = ""
		}

		if v7.Valid {
			v.EntrustDirection = v7.String
		} else {
			v.EntrustDirection = ""
		}

		if v8.Valid {
			v.OpenClose = v8.String
		} else {
			v.OpenClose = ""
		}

		if v9.Valid {
			v.InvestType = v9.String
		} else {
			v.InvestType = ""
		}

		if v10.Valid {
			v.Amount = v10.Float64
		} else {
			v.Amount = 0
		}

		if v11.Valid {
			v.Balance = v11.Float64
		} else {
			v.Balance = 0
		}

		if v12.Valid {
			v.ContractSize = v12.Float64
		} else {
			v.ContractSize = 0
		}

		if v13.Valid {
			v.Price = v13.Float64
		} else {
			v.Price = 0
		}

		if v14.Valid {
			v.StockEntrustExecuteStatus = v14.String
		} else {
			v.StockEntrustExecuteStatus = ""
		}

		if v15.Valid {
			v.StockDealExecuteStatus = v15.String
		} else {
			v.StockDealExecuteStatus = ""
		}

		if v16.Valid {
			v.TotalDealAmount = v16.Float64
		} else {
			v.TotalDealAmount = 0
		}

		if v17.Valid {
			v.TotalDealBalance = v17.Float64
		} else {
			v.TotalDealBalance = 0
		}

		if v18.Valid {
			v.CumAvgPrice = v18.Float64
		} else {
			v.CumAvgPrice = 0
		}

		if v19.Valid {
			v.TotalEntrustAmount = v19.Float64
		} else {
			v.TotalEntrustAmount = 0
		}

		if v20.Valid {
			v.TotalEntrustBalance = v20.Float64
		} else {
			v.TotalEntrustBalance = 0
		}

		if v21.Valid {
			v.DealCompleteDateTime = v21.Int64
		} else {
			v.DealCompleteDateTime = 0
		}

		if v22.Valid {
			v.EstimateFee = v22.Float64
		} else {
			v.EstimateFee = 0
		}

		if v23.Valid {
			v.StockInstrExecutionTime = v23.Int64
		} else {
			v.StockInstrExecutionTime = 0
		}

		if v24.Valid {
			v.StockInstrOperator = v24.String
		} else {
			v.StockInstrOperator = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTaskInstrStock(v *schema.TaskInstrStock) []interface{} {
	var v0 int64
	var v1 int
	var v2 int64
	var v3 int64
	var v4 int64
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 string
	var v10 float64
	var v11 float64
	var v12 float64
	var v13 float64
	var v14 string
	var v15 string
	var v16 float64
	var v17 float64
	var v18 float64
	var v19 float64
	var v20 float64
	var v21 int64
	var v22 float64
	var v23 int64
	var v24 string

	v0 = v.ID
	v1 = v.Date
	v2 = v.DailyInstrNo
	v3 = v.IndexDailyModify
	v4 = v.StockSerialNo
	v5 = v.MarketNo
	v6 = v.ReportCode
	v7 = v.EntrustDirection
	v8 = v.OpenClose
	v9 = v.InvestType
	v10 = v.Amount
	v11 = v.Balance
	v12 = v.ContractSize
	v13 = v.Price
	v14 = v.StockEntrustExecuteStatus
	v15 = v.StockDealExecuteStatus
	v16 = v.TotalDealAmount
	v17 = v.TotalDealBalance
	v18 = v.CumAvgPrice
	v19 = v.TotalEntrustAmount
	v20 = v.TotalEntrustBalance
	v21 = v.DealCompleteDateTime
	v22 = v.EstimateFee
	v23 = v.StockInstrExecutionTime
	v24 = v.StockInstrOperator

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
		v9,
		v10,
		v11,
		v12,
		v13,
		v14,
		v15,
		v16,
		v17,
		v18,
		v19,
		v20,
		v21,
		v22,
		v23,
		v24,
	}
}

func genericSelectTaskInstrStock(db db.SimpleDB, query string, args ...interface{}) (*schema.TaskInstrStock, error) {
	row := db.QueryRow(query, args...)
	return scanTaskInstrStock(row)
}

func genericSelectTaskInstrStocks(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TaskInstrStock, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTaskInstrStocks(rows)
}

func InsertTaskInstrStock(db db.SimpleDB, v *schema.TaskInstrStock) error {

	res, err := db.Exec(InsertTaskInstrStockStmt, sliceTaskInstrStock(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTaskInstrStockById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTaskInstrStockByIdStmt, args...)
	return err
}

func DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) error {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	_, err := db.Exec(DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return err
}

func DeleteTaskInstrStockByDate(db db.SimpleDB, date int) error {
	args := []interface{}{date}
	_, err := db.Exec(DeleteTaskInstrStockByDateStmt, args...)
	return err
}

func DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) error {
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	_, err := db.Exec(DeleteTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	return err
}

func UpdateTaskInstrStockById(db db.SimpleDB, v *schema.TaskInstrStock) error {
	args := sliceTaskInstrStock(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTaskInstrStockByIdStmt, args...)
	return err
}

func UpdateTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, v *schema.TaskInstrStock) error {
	args := sliceTaskInstrStock(v)
	args = append(args, v.Date, v.DailyInstrNo, v.IndexDailyModify, v.StockSerialNo)
	_, err := db.Exec(UpdateTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	return err
}

func GetTaskInstrStockById(db db.SimpleDB, iD int64) (*schema.TaskInstrStock, error) {
	args := []interface{}{iD}
	v, err := genericSelectTaskInstrStock(db, SelectTaskInstrStockByIdStmt, args...)
	return v, err
}

func GetTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) (*schema.TaskInstrStock, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	v, err := genericSelectTaskInstrStock(db, SelectTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	return v, err
}

func FindAllTaskInstrStocks(db db.SimpleDB) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockStmt, args...)
	return v, err
}

func FindAllTaskInstrStocksInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockRangeStmt, args...)
	return v, err
}

func FindTaskInstrStocksByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return v, err
}

func FindTaskInstrStocksByDateAndDailyInstrNoAndIndexDailyModifyInRange(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, limit int64, offset int64) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify, limit, offset}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockRangeByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return v, err
}

func FindTaskInstrStocksByDate(db db.SimpleDB, date int) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{date}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockByDateStmt, args...)
	return v, err
}

func FindTaskInstrStocksByDateInRange(db db.SimpleDB, date int, limit int64, offset int64) ([]*schema.TaskInstrStock, error) {
	args := []interface{}{date, limit, offset}
	v, err := genericSelectTaskInstrStocks(db, SelectTaskInstrStockRangeByDateStmt, args...)
	return v, err
}

func GetTaskInstrStockOfTaskInstrByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) (*schema.TaskInstrStock, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	v, err := genericSelectTaskInstrStock(db, SelectTaskInstrStockOfTaskInstrByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	return v, err
}

func CountTaskInstrStock(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTaskInstrStockCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModify(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64) (int, error) {
	var count int
	args := []interface{}{date, dailyInstrNo, indexDailyModify}
	row := db.QueryRow(SelectTaskInstrStockCountByDateAndDailyInstrNoAndIndexDailyModifyStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrStockByDate(db db.SimpleDB, date int) (int, error) {
	var count int
	args := []interface{}{date}
	row := db.QueryRow(SelectTaskInstrStockCountByDateStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrStockByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) (int, error) {
	var count int
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	row := db.QueryRow(SelectTaskInstrStockCountByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const SelectTaskInstrViewStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
`

const SelectTaskInstrViewRangeStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
LIMIT ? OFFSET ?
`

const SelectTaskInstrViewCountStmt = `
SELECT count(1)
FROM task_instr_views 
`

const SelectTaskInstrViewByIdStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_id=?
`

const SelectTaskInstrViewByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const SelectTaskInstrViewCountByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
SELECT count(1)
FROM task_instr_views 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const DeleteTaskInstrViewByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt = `
DELETE FROM task_instr_views 
WHERE f_date=?
AND f_daily_instr_no=?
AND f_index_daily_modify=?
AND f_stock_serial_no=?
`

const SelectTaskInstrViewByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_date=?
`

const SelectTaskInstrViewCountByDateStmt = `
SELECT count(1)
FROM task_instr_views 
WHERE f_date=?
`

const SelectTaskInstrViewRangeByDateStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_date=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrViewByDateStmt = `
DELETE FROM task_instr_views 
WHERE f_date=?
`

const SelectTaskInstrViewByDirectOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_direct_operator=?
`

const SelectTaskInstrViewCountByDirectOperatorStmt = `
SELECT count(1)
FROM task_instr_views 
WHERE f_direct_operator=?
`

const SelectTaskInstrViewRangeByDirectOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_direct_operator=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrViewByDirectOperatorStmt = `
DELETE FROM task_instr_views 
WHERE f_direct_operator=?
`

const SelectTaskInstrViewByOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_operator=?
`

const SelectTaskInstrViewCountByOperatorStmt = `
SELECT count(1)
FROM task_instr_views 
WHERE f_operator=?
`

const SelectTaskInstrViewRangeByOperatorStmt = `
SELECT 
 f_id
,f_date
,f_daily_instr_no
,f_index_daily_modify
,f_batch_serial_no
,f_index_last_modify
,f_account_no
,f_combi_no
,f_instr_type
,f_begin_date
,f_end_date
,f_begin_time
,f_end_time
,f_direct_date
,f_direct_time
,f_direct_operator
,f_modify_date
,f_modify_time
,f_modify_operator
,f_modify_reason
,f_dispense_date
,f_dispense_time
,f_dispense_operator
,f_dispense_refuse_reason
,f_cancel_date
,f_cancel_time
,f_cancel_operator
,f_cancel_reason
,f_operator
,f_instr_status
,f_dispense_status
,f_entrust_execute_status
,f_deal_execute_status
,f_create_time
,f_business_type
,f_lock_flag
,f_limit_operator
,f_org_id
,f_dept_id
,f_ip_address
,f_mac
,f_volserial_no
,f_stock_serial_no
,f_market_no
,f_report_code
,f_entrust_direction
,f_open_close
,f_invest_type
,f_amount
,f_entrustable_amount
,f_balance
,f_contract_size
,f_price
,f_stock_entrust_execute_status
,f_stock_deal_execute_status
,f_total_deal_amount
,f_total_deal_balance
,f_cum_avg_price
,f_total_entrust_amount
,f_total_entrust_balance
,f_deal_complete_date_time
,f_estimate_fee
,f_stock_instr_execution_time
,f_stock_instr_operator
FROM task_instr_views 
WHERE f_operator=?
LIMIT ? OFFSET ?
`

const DeleteTaskInstrViewByOperatorStmt = `
DELETE FROM task_instr_views 
WHERE f_operator=?
`

func scanTaskInstrView(row *sql.Row) (*schema.TaskInstrView, error) {
	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
	var v13 sql.NullInt64
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullInt64
	var v17 sql.NullInt64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullInt64
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullInt64
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullString
	var v29 sql.NullString
	var v30 sql.NullString
	var v31 sql.NullString
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullString
	var v40 sql.NullString
	var v41 sql.NullString
	var v42 sql.NullInt64
	var v43 sql.NullString
	var v44 sql.NullString
	var v45 sql.NullString
	var v46 sql.NullString
	var v47 sql.NullString
	var v48 sql.NullFloat64
	var v49 sql.NullFloat64
	var v50 sql.NullFloat64
	var v51 sql.NullFloat64
	var v52 sql.NullFloat64
	var v53 sql.NullString
	var v54 sql.NullString
	var v55 sql.NullFloat64
	var v56 sql.NullFloat64
	var v57 sql.NullFloat64
	var v58 sql.NullFloat64
	var v59 sql.NullFloat64
	var v60 sql.NullInt64
	var v61 sql.NullFloat64
	var v62 sql.NullInt64
	var v63 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
		&v7,
		&v8,
		&v9,
		&v10,
		&v11,
		&v12,
		&v13,
		&v14,
		&v15,
		&v16,
		&v17,
		&v18,
		&v19,
		&v20,
		&v21,
		&v22,
		&v23,
		&v24,
		&v25,
		&v26,
		&v27,
		&v28,
		&v29,
		&v30,
		&v31,
		&v32,
		&v33,
		&v34,
		&v35,
		&v36,
		&v37,
		&v38,
		&v39,
		&v40,
		&v41,
		&v42,
		&v43,
		&v44,
		&v45,
		&v46,
		&v47,
		&v48,
		&v49,
		&v50,
		&v51,
		&v52,
		&v53,
		&v54,
		&v55,
		&v56,
		&v57,
		&v58,
		&v59,
		&v60,
		&v61,
		&v62,
		&v63,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TaskInstrView{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.Date = int(v1.Int64)
	} else {
		v.Date = 0
	}

	if v2.Valid {
		v.DailyInstrNo = v2.Int64
	} else {
		v.DailyInstrNo = 0
	}

	if v3.Valid {
		v.IndexDailyModify = v3.Int64
	} else {
		v.IndexDailyModify = 0
	}

	if v4.Valid {
		v.BatchSerialNo = v4.Int64
	} else {
		v.BatchSerialNo = 0
	}

	if v5.Valid {
		v.IndexLastModify = v5.Int64
	} else {
		v.IndexLastModify = 0
	}

	if v6.Valid {
		v.AccountNo = v6.String
	} else {
		v.AccountNo = ""
	}

	if v7.Valid {
		v.CombiNo = v7.String
	} else {
		v.CombiNo = ""
	}

	if v8.Valid {
		v.InstrType = v8.String
	} else {
		v.InstrType = ""
	}

	if v9.Valid {
		v.BeginDate = int(v9.Int64)
	} else {
		v.BeginDate = 0
	}

	if v10.Valid {
		v.EndDate = int(v10.Int64)
	} else {
		v.EndDate = 0
	}

	if v11.Valid {
		v.BeginTime = int(v11.Int64)
	} else {
		v.BeginTime = 0
	}

	if v12.Valid {
		v.EndTime = int(v12.Int64)
	} else {
		v.EndTime = 0
	}

	if v13.Valid {
		v.DirectDate = int(v13.Int64)
	} else {
		v.DirectDate = 0
	}

	if v14.Valid {
		v.DirectTime = int(v14.Int64)
	} else {
		v.DirectTime = 0
	}

	if v15.Valid {
		v.DirectOperator = v15.String
	} else {
		v.DirectOperator = ""
	}

	if v16.Valid {
		v.ModifyDate = int(v16.Int64)
	} else {
		v.ModifyDate = 0
	}

	if v17.Valid {
		v.ModifyTime = int(v17.Int64)
	} else {
		v.ModifyTime = 0
	}

	if v18.Valid {
		v.ModifyOperator = v18.String
	} else {
		v.ModifyOperator = ""
	}

	if v19.Valid {
		v.ModifyReason = v19.String
	} else {
		v.ModifyReason = ""
	}

	if v20.Valid {
		v.DispenseDate = int(v20.Int64)
	} else {
		v.DispenseDate = 0
	}

	if v21.Valid {
		v.DispenseTime = int(v21.Int64)
	} else {
		v.DispenseTime = 0
	}

	if v22.Valid {
		v.DispenseOperator = v22.String
	} else {
		v.DispenseOperator = ""
	}

	if v23.Valid {
		v.DispenseRefuseReason = v23.String
	} else {
		v.DispenseRefuseReason = ""
	}

	if v24.Valid {
		v.CancelDate = int(v24.Int64)
	} else {
		v.CancelDate = 0
	}

	if v25.Valid {
		v.CancelTime = int(v25.Int64)
	} else {
		v.CancelTime = 0
	}

	if v26.Valid {
		v.CancelOperator = v26.String
	} else {
		v.CancelOperator = ""
	}

	if v27.Valid {
		v.CancelReason = v27.String
	} else {
		v.CancelReason = ""
	}

	if v28.Valid {
		v.Operator = v28.String
	} else {
		v.Operator = ""
	}

	if v29.Valid {
		v.InstrStatus = v29.String
	} else {
		v.InstrStatus = ""
	}

	if v30.Valid {
		v.DispenseStatus = v30.String
	} else {
		v.DispenseStatus = ""
	}

	if v31.Valid {
		v.EntrustExecuteStatus = v31.String
	} else {
		v.EntrustExecuteStatus = ""
	}

	if v32.Valid {
		v.DealExecuteStatus = v32.String
	} else {
		v.DealExecuteStatus = ""
	}

	if v33.Valid {
		v.CreateTime = v33.Int64
	} else {
		v.CreateTime = 0
	}

	if v34.Valid {
		v.BusinessType = v34.String
	} else {
		v.BusinessType = ""
	}

	if v35.Valid {
		v.LockFlag = int(v35.Int64)
	} else {
		v.LockFlag = 0
	}

	if v36.Valid {
		v.LimitOperator = v36.String
	} else {
		v.LimitOperator = ""
	}

	if v37.Valid {
		v.OrgId = v37.Int64
	} else {
		v.OrgId = 0
	}

	if v38.Valid {
		v.DeptId = v38.Int64
	} else {
		v.DeptId = 0
	}

	if v39.Valid {
		v.IpAddress = v39.String
	} else {
		v.IpAddress = ""
	}

	if v40.Valid {
		v.Mac = v40.String
	} else {
		v.Mac = ""
	}

	if v41.Valid {
		v.VolserialNo = v41.String
	} else {
		v.VolserialNo = ""
	}

	if v42.Valid {
		v.StockSerialNo = v42.Int64
	} else {
		v.StockSerialNo = 0
	}

	if v43.Valid {
		v.MarketNo = v43.String
	} else {
		v.MarketNo = ""
	}

	if v44.Valid {
		v.ReportCode = v44.String
	} else {
		v.ReportCode = ""
	}

	if v45.Valid {
		v.EntrustDirection = v45.String
	} else {
		v.EntrustDirection = ""
	}

	if v46.Valid {
		v.OpenClose = v46.String
	} else {
		v.OpenClose = ""
	}

	if v47.Valid {
		v.InvestType = v47.String
	} else {
		v.InvestType = ""
	}

	if v48.Valid {
		v.Amount = v48.Float64
	} else {
		v.Amount = 0
	}

	if v49.Valid {
		v.EntrustableAmount = v49.Float64
	} else {
		v.EntrustableAmount = 0
	}

	if v50.Valid {
		v.Balance = v50.Float64
	} else {
		v.Balance = 0
	}

	if v51.Valid {
		v.ContractSize = v51.Float64
	} else {
		v.ContractSize = 0
	}

	if v52.Valid {
		v.Price = v52.Float64
	} else {
		v.Price = 0
	}

	if v53.Valid {
		v.StockEntrustExecuteStatus = v53.String
	} else {
		v.StockEntrustExecuteStatus = ""
	}

	if v54.Valid {
		v.StockDealExecuteStatus = v54.String
	} else {
		v.StockDealExecuteStatus = ""
	}

	if v55.Valid {
		v.TotalDealAmount = v55.Float64
	} else {
		v.TotalDealAmount = 0
	}

	if v56.Valid {
		v.TotalDealBalance = v56.Float64
	} else {
		v.TotalDealBalance = 0
	}

	if v57.Valid {
		v.CumAvgPrice = v57.Float64
	} else {
		v.CumAvgPrice = 0
	}

	if v58.Valid {
		v.TotalEntrustAmount = v58.Float64
	} else {
		v.TotalEntrustAmount = 0
	}

	if v59.Valid {
		v.TotalEntrustBalance = v59.Float64
	} else {
		v.TotalEntrustBalance = 0
	}

	if v60.Valid {
		v.DealCompleteDateTime = v60.Int64
	} else {
		v.DealCompleteDateTime = 0
	}

	if v61.Valid {
		v.EstimateFee = v61.Float64
	} else {
		v.EstimateFee = 0
	}

	if v62.Valid {
		v.StockInstrExecutionTime = v62.Int64
	} else {
		v.StockInstrExecutionTime = 0
	}

	if v63.Valid {
		v.StockInstrOperator = v63.String
	} else {
		v.StockInstrOperator = ""
	}

	return v, nil
}

func scanTaskInstrViews(rows *sql.Rows) ([]*schema.TaskInstrView, error) {
	var err error
	var vv []*schema.TaskInstrView

	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullInt64
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
	var v13 sql.NullInt64
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullInt64
	var v17 sql.NullInt64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullInt64
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullInt64
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullString
	var v29 sql.NullString
	var v30 sql.NullString
	var v31 sql.NullString
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullString
	var v40 sql.NullString
	var v41 sql.NullString
	var v42 sql.NullInt64
	var v43 sql.NullString
	var v44 sql.NullString
	var v45 sql.NullString
	var v46 sql.NullString
	var v47 sql.NullString
	var v48 sql.NullFloat64
	var v49 sql.NullFloat64
	var v50 sql.NullFloat64
	var v51 sql.NullFloat64
	var v52 sql.NullFloat64
	var v53 sql.NullString
	var v54 sql.NullString
	var v55 sql.NullFloat64
	var v56 sql.NullFloat64
	var v57 sql.NullFloat64
	var v58 sql.NullFloat64
	var v59 sql.NullFloat64
	var v60 sql.NullInt64
	var v61 sql.NullFloat64
	var v62 sql.NullInt64
	var v63 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
			&v10,
			&v11,
			&v12,
			&v13,
			&v14,
			&v15,
			&v16,
			&v17,
			&v18,
			&v19,
			&v20,
			&v21,
			&v22,
			&v23,
			&v24,
			&v25,
			&v26,
			&v27,
			&v28,
			&v29,
			&v30,
			&v31,
			&v32,
			&v33,
			&v34,
			&v35,
			&v36,
			&v37,
			&v38,
			&v39,
			&v40,
			&v41,
			&v42,
			&v43,
			&v44,
			&v45,
			&v46,
			&v47,
			&v48,
			&v49,
			&v50,
			&v51,
			&v52,
			&v53,
			&v54,
			&v55,
			&v56,
			&v57,
			&v58,
			&v59,
			&v60,
			&v61,
			&v62,
			&v63,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TaskInstrView{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.Date = int(v1.Int64)
		} else {
			v.Date = 0
		}

		if v2.Valid {
			v.DailyInstrNo = v2.Int64
		} else {
			v.DailyInstrNo = 0
		}

		if v3.Valid {
			v.IndexDailyModify = v3.Int64
		} else {
			v.IndexDailyModify = 0
		}

		if v4.Valid {
			v.BatchSerialNo = v4.Int64
		} else {
			v.BatchSerialNo = 0
		}

		if v5.Valid {
			v.IndexLastModify = v5.Int64
		} else {
			v.IndexLastModify = 0
		}

		if v6.Valid {
			v.AccountNo = v6.String
		} else {
			v.AccountNo = ""
		}

		if v7.Valid {
			v.CombiNo = v7.String
		} else {
			v.CombiNo = ""
		}

		if v8.Valid {
			v.InstrType = v8.String
		} else {
			v.InstrType = ""
		}

		if v9.Valid {
			v.BeginDate = int(v9.Int64)
		} else {
			v.BeginDate = 0
		}

		if v10.Valid {
			v.EndDate = int(v10.Int64)
		} else {
			v.EndDate = 0
		}

		if v11.Valid {
			v.BeginTime = int(v11.Int64)
		} else {
			v.BeginTime = 0
		}

		if v12.Valid {
			v.EndTime = int(v12.Int64)
		} else {
			v.EndTime = 0
		}

		if v13.Valid {
			v.DirectDate = int(v13.Int64)
		} else {
			v.DirectDate = 0
		}

		if v14.Valid {
			v.DirectTime = int(v14.Int64)
		} else {
			v.DirectTime = 0
		}

		if v15.Valid {
			v.DirectOperator = v15.String
		} else {
			v.DirectOperator = ""
		}

		if v16.Valid {
			v.ModifyDate = int(v16.Int64)
		} else {
			v.ModifyDate = 0
		}

		if v17.Valid {
			v.ModifyTime = int(v17.Int64)
		} else {
			v.ModifyTime = 0
		}

		if v18.Valid {
			v.ModifyOperator = v18.String
		} else {
			v.ModifyOperator = ""
		}

		if v19.Valid {
			v.ModifyReason = v19.String
		} else {
			v.ModifyReason = ""
		}

		if v20.Valid {
			v.DispenseDate = int(v20.Int64)
		} else {
			v.DispenseDate = 0
		}

		if v21.Valid {
			v.DispenseTime = int(v21.Int64)
		} else {
			v.DispenseTime = 0
		}

		if v22.Valid {
			v.DispenseOperator = v22.String
		} else {
			v.DispenseOperator = ""
		}

		if v23.Valid {
			v.DispenseRefuseReason = v23.String
		} else {
			v.DispenseRefuseReason = ""
		}

		if v24.Valid {
			v.CancelDate = int(v24.Int64)
		} else {
			v.CancelDate = 0
		}

		if v25.Valid {
			v.CancelTime = int(v25.Int64)
		} else {
			v.CancelTime = 0
		}

		if v26.Valid {
			v.CancelOperator = v26.String
		} else {
			v.CancelOperator = ""
		}

		if v27.Valid {
			v.CancelReason = v27.String
		} else {
			v.CancelReason = ""
		}

		if v28.Valid {
			v.Operator = v28.String
		} else {
			v.Operator = ""
		}

		if v29.Valid {
			v.InstrStatus = v29.String
		} else {
			v.InstrStatus = ""
		}

		if v30.Valid {
			v.DispenseStatus = v30.String
		} else {
			v.DispenseStatus = ""
		}

		if v31.Valid {
			v.EntrustExecuteStatus = v31.String
		} else {
			v.EntrustExecuteStatus = ""
		}

		if v32.Valid {
			v.DealExecuteStatus = v32.String
		} else {
			v.DealExecuteStatus = ""
		}

		if v33.Valid {
			v.CreateTime = v33.Int64
		} else {
			v.CreateTime = 0
		}

		if v34.Valid {
			v.BusinessType = v34.String
		} else {
			v.BusinessType = ""
		}

		if v35.Valid {
			v.LockFlag = int(v35.Int64)
		} else {
			v.LockFlag = 0
		}

		if v36.Valid {
			v.LimitOperator = v36.String
		} else {
			v.LimitOperator = ""
		}

		if v37.Valid {
			v.OrgId = v37.Int64
		} else {
			v.OrgId = 0
		}

		if v38.Valid {
			v.DeptId = v38.Int64
		} else {
			v.DeptId = 0
		}

		if v39.Valid {
			v.IpAddress = v39.String
		} else {
			v.IpAddress = ""
		}

		if v40.Valid {
			v.Mac = v40.String
		} else {
			v.Mac = ""
		}

		if v41.Valid {
			v.VolserialNo = v41.String
		} else {
			v.VolserialNo = ""
		}

		if v42.Valid {
			v.StockSerialNo = v42.Int64
		} else {
			v.StockSerialNo = 0
		}

		if v43.Valid {
			v.MarketNo = v43.String
		} else {
			v.MarketNo = ""
		}

		if v44.Valid {
			v.ReportCode = v44.String
		} else {
			v.ReportCode = ""
		}

		if v45.Valid {
			v.EntrustDirection = v45.String
		} else {
			v.EntrustDirection = ""
		}

		if v46.Valid {
			v.OpenClose = v46.String
		} else {
			v.OpenClose = ""
		}

		if v47.Valid {
			v.InvestType = v47.String
		} else {
			v.InvestType = ""
		}

		if v48.Valid {
			v.Amount = v48.Float64
		} else {
			v.Amount = 0
		}

		if v49.Valid {
			v.EntrustableAmount = v49.Float64
		} else {
			v.EntrustableAmount = 0
		}

		if v50.Valid {
			v.Balance = v50.Float64
		} else {
			v.Balance = 0
		}

		if v51.Valid {
			v.ContractSize = v51.Float64
		} else {
			v.ContractSize = 0
		}

		if v52.Valid {
			v.Price = v52.Float64
		} else {
			v.Price = 0
		}

		if v53.Valid {
			v.StockEntrustExecuteStatus = v53.String
		} else {
			v.StockEntrustExecuteStatus = ""
		}

		if v54.Valid {
			v.StockDealExecuteStatus = v54.String
		} else {
			v.StockDealExecuteStatus = ""
		}

		if v55.Valid {
			v.TotalDealAmount = v55.Float64
		} else {
			v.TotalDealAmount = 0
		}

		if v56.Valid {
			v.TotalDealBalance = v56.Float64
		} else {
			v.TotalDealBalance = 0
		}

		if v57.Valid {
			v.CumAvgPrice = v57.Float64
		} else {
			v.CumAvgPrice = 0
		}

		if v58.Valid {
			v.TotalEntrustAmount = v58.Float64
		} else {
			v.TotalEntrustAmount = 0
		}

		if v59.Valid {
			v.TotalEntrustBalance = v59.Float64
		} else {
			v.TotalEntrustBalance = 0
		}

		if v60.Valid {
			v.DealCompleteDateTime = v60.Int64
		} else {
			v.DealCompleteDateTime = 0
		}

		if v61.Valid {
			v.EstimateFee = v61.Float64
		} else {
			v.EstimateFee = 0
		}

		if v62.Valid {
			v.StockInstrExecutionTime = v62.Int64
		} else {
			v.StockInstrExecutionTime = 0
		}

		if v63.Valid {
			v.StockInstrOperator = v63.String
		} else {
			v.StockInstrOperator = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTaskInstrView(v *schema.TaskInstrView) []interface{} {
	var v0 int64
	var v1 int
	var v2 int64
	var v3 int64
	var v4 int64
	var v5 int64
	var v6 string
	var v7 string
	var v8 string
	var v9 int
	var v10 int
	var v11 int
	var v12 int
	var v13 int
	var v14 int
	var v15 string
	var v16 int
	var v17 int
	var v18 string
	var v19 string
	var v20 int
	var v21 int
	var v22 string
	var v23 string
	var v24 int
	var v25 int
	var v26 string
	var v27 string
	var v28 string
	var v29 string
	var v30 string
	var v31 string
	var v32 string
	var v33 int64
	var v34 string
	var v35 int
	var v36 string
	var v37 int64
	var v38 int64
	var v39 string
	var v40 string
	var v41 string
	var v42 int64
	var v43 string
	var v44 string
	var v45 string
	var v46 string
	var v47 string
	var v48 float64
	var v49 float64
	var v50 float64
	var v51 float64
	var v52 float64
	var v53 string
	var v54 string
	var v55 float64
	var v56 float64
	var v57 float64
	var v58 float64
	var v59 float64
	var v60 int64
	var v61 float64
	var v62 int64
	var v63 string

	v0 = v.ID
	v1 = v.Date
	v2 = v.DailyInstrNo
	v3 = v.IndexDailyModify
	v4 = v.BatchSerialNo
	v5 = v.IndexLastModify
	v6 = v.AccountNo
	v7 = v.CombiNo
	v8 = v.InstrType
	v9 = v.BeginDate
	v10 = v.EndDate
	v11 = v.BeginTime
	v12 = v.EndTime
	v13 = v.DirectDate
	v14 = v.DirectTime
	v15 = v.DirectOperator
	v16 = v.ModifyDate
	v17 = v.ModifyTime
	v18 = v.ModifyOperator
	v19 = v.ModifyReason
	v20 = v.DispenseDate
	v21 = v.DispenseTime
	v22 = v.DispenseOperator
	v23 = v.DispenseRefuseReason
	v24 = v.CancelDate
	v25 = v.CancelTime
	v26 = v.CancelOperator
	v27 = v.CancelReason
	v28 = v.Operator
	v29 = v.InstrStatus
	v30 = v.DispenseStatus
	v31 = v.EntrustExecuteStatus
	v32 = v.DealExecuteStatus
	v33 = v.CreateTime
	v34 = v.BusinessType
	v35 = v.LockFlag
	v36 = v.LimitOperator
	v37 = v.OrgId
	v38 = v.DeptId
	v39 = v.IpAddress
	v40 = v.Mac
	v41 = v.VolserialNo
	v42 = v.StockSerialNo
	v43 = v.MarketNo
	v44 = v.ReportCode
	v45 = v.EntrustDirection
	v46 = v.OpenClose
	v47 = v.InvestType
	v48 = v.Amount
	v49 = v.EntrustableAmount
	v50 = v.Balance
	v51 = v.ContractSize
	v52 = v.Price
	v53 = v.StockEntrustExecuteStatus
	v54 = v.StockDealExecuteStatus
	v55 = v.TotalDealAmount
	v56 = v.TotalDealBalance
	v57 = v.CumAvgPrice
	v58 = v.TotalEntrustAmount
	v59 = v.TotalEntrustBalance
	v60 = v.DealCompleteDateTime
	v61 = v.EstimateFee
	v62 = v.StockInstrExecutionTime
	v63 = v.StockInstrOperator

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
		v9,
		v10,
		v11,
		v12,
		v13,
		v14,
		v15,
		v16,
		v17,
		v18,
		v19,
		v20,
		v21,
		v22,
		v23,
		v24,
		v25,
		v26,
		v27,
		v28,
		v29,
		v30,
		v31,
		v32,
		v33,
		v34,
		v35,
		v36,
		v37,
		v38,
		v39,
		v40,
		v41,
		v42,
		v43,
		v44,
		v45,
		v46,
		v47,
		v48,
		v49,
		v50,
		v51,
		v52,
		v53,
		v54,
		v55,
		v56,
		v57,
		v58,
		v59,
		v60,
		v61,
		v62,
		v63,
	}
}

func genericSelectTaskInstrView(db db.SimpleDB, query string, args ...interface{}) (*schema.TaskInstrView, error) {
	row := db.QueryRow(query, args...)
	return scanTaskInstrView(row)
}

func genericSelectTaskInstrViews(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TaskInstrView, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTaskInstrViews(rows)
}

func GetTaskInstrViewById(db db.SimpleDB, iD int64) (*schema.TaskInstrView, error) {
	args := []interface{}{iD}
	v, err := genericSelectTaskInstrView(db, SelectTaskInstrViewByIdStmt, args...)
	return v, err
}

func GetTaskInstrViewByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) (*schema.TaskInstrView, error) {
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	v, err := genericSelectTaskInstrView(db, SelectTaskInstrViewByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	return v, err
}

func FindAllTaskInstrViews(db db.SimpleDB) ([]*schema.TaskInstrView, error) {
	args := []interface{}{}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewStmt, args...)
	return v, err
}

func FindAllTaskInstrViewsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TaskInstrView, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewRangeStmt, args...)
	return v, err
}

func FindTaskInstrViewsByDate(db db.SimpleDB, date int) ([]*schema.TaskInstrView, error) {
	args := []interface{}{date}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewByDateStmt, args...)
	return v, err
}

func FindTaskInstrViewsByDateInRange(db db.SimpleDB, date int, limit int64, offset int64) ([]*schema.TaskInstrView, error) {
	args := []interface{}{date, limit, offset}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewRangeByDateStmt, args...)
	return v, err
}

func FindTaskInstrViewsByDirectOperator(db db.SimpleDB, directOperator string) ([]*schema.TaskInstrView, error) {
	args := []interface{}{directOperator}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewByDirectOperatorStmt, args...)
	return v, err
}

func FindTaskInstrViewsByDirectOperatorInRange(db db.SimpleDB, directOperator string, limit int64, offset int64) ([]*schema.TaskInstrView, error) {
	args := []interface{}{directOperator, limit, offset}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewRangeByDirectOperatorStmt, args...)
	return v, err
}

func FindTaskInstrViewsByOperator(db db.SimpleDB, operator string) ([]*schema.TaskInstrView, error) {
	args := []interface{}{operator}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewByOperatorStmt, args...)
	return v, err
}

func FindTaskInstrViewsByOperatorInRange(db db.SimpleDB, operator string, limit int64, offset int64) ([]*schema.TaskInstrView, error) {
	args := []interface{}{operator, limit, offset}
	v, err := genericSelectTaskInstrViews(db, SelectTaskInstrViewRangeByOperatorStmt, args...)
	return v, err
}

func CountTaskInstrView(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTaskInstrViewCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrViewByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNo(db db.SimpleDB, date int, dailyInstrNo int64, indexDailyModify int64, stockSerialNo int64) (int, error) {
	var count int
	args := []interface{}{date, dailyInstrNo, indexDailyModify, stockSerialNo}
	row := db.QueryRow(SelectTaskInstrViewCountByDateAndDailyInstrNoAndIndexDailyModifyAndStockSerialNoStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrViewByDate(db db.SimpleDB, date int) (int, error) {
	var count int
	args := []interface{}{date}
	row := db.QueryRow(SelectTaskInstrViewCountByDateStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrViewByDirectOperator(db db.SimpleDB, directOperator string) (int, error) {
	var count int
	args := []interface{}{directOperator}
	row := db.QueryRow(SelectTaskInstrViewCountByDirectOperatorStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTaskInstrViewByOperator(db db.SimpleDB, operator string) (int, error) {
	var count int
	args := []interface{}{operator}
	row := db.QueryRow(SelectTaskInstrViewCountByOperatorStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateAssetUnitStmt = `
CREATE TABLE IF NOT EXISTS asset_units (
 f_id           BIGINT PRIMARY KEY AUTO_INCREMENT
,f_account_no   VARCHAR(16)
,f_account_name VARCHAR(256)
,f_combi_no     VARCHAR(16)
,f_combi_name   VARCHAR(256)
);
`

const InsertAssetUnitStmt = `
INSERT INTO asset_units (
 f_account_no
,f_account_name
,f_combi_no
,f_combi_name
) VALUES (?,?,?,?)
`

const SelectAssetUnitStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
`

const SelectAssetUnitRangeStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
LIMIT ? OFFSET ?
`

const SelectAssetUnitCountStmt = `
SELECT count(1)
FROM asset_units 
`

const SelectAssetUnitByIdStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
WHERE f_id=?
`

const UpdateAssetUnitByIdStmt = `
UPDATE asset_units SET 
 f_id=?
,f_account_no=?
,f_account_name=?
,f_combi_no=?
,f_combi_name=? 
WHERE f_id=?
`

const DeleteAssetUnitByIdStmt = `
DELETE FROM asset_units 
WHERE f_id=?
`

const CreateIAuAnStmt = `
CREATE INDEX i_au_an ON asset_units (f_account_no);
`

const SelectAssetUnitByAccountNoStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
WHERE f_account_no=?
`

const SelectAssetUnitCountByAccountNoStmt = `
SELECT count(1)
FROM asset_units 
WHERE f_account_no=?
`

const SelectAssetUnitRangeByAccountNoStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
WHERE f_account_no=?
LIMIT ? OFFSET ?
`

const DeleteAssetUnitByAccountNoStmt = `
DELETE FROM asset_units 
WHERE f_account_no=?
`

const CreateUqAuStmt = `
CREATE UNIQUE INDEX uq_au ON asset_units (f_account_no,f_combi_no);
`

const SelectAssetUnitByAccountNoAndCombiNoStmt = `
SELECT 
 f_id
,f_account_no
,f_account_name
,f_combi_no
,f_combi_name
FROM asset_units 
WHERE f_account_no=?
AND f_combi_no=?
`

const SelectAssetUnitCountByAccountNoAndCombiNoStmt = `
SELECT count(1)
FROM asset_units 
WHERE f_account_no=?
AND f_combi_no=?
`

const UpdateAssetUnitByAccountNoAndCombiNoStmt = `
UPDATE asset_units SET 
 f_id=?
,f_account_no=?
,f_account_name=?
,f_combi_no=?
,f_combi_name=? 
WHERE f_account_no=?
AND f_combi_no=?
`

const DeleteAssetUnitByAccountNoAndCombiNoStmt = `
DELETE FROM asset_units 
WHERE f_account_no=?
AND f_combi_no=?
`

func scanAssetUnit(row *sql.Row) (*schema.AssetUnit, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.AssetUnit{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.AccountNo = v1.String
	} else {
		v.AccountNo = ""
	}

	if v2.Valid {
		v.AccountName = v2.String
	} else {
		v.AccountName = ""
	}

	if v3.Valid {
		v.CombiNo = v3.String
	} else {
		v.CombiNo = ""
	}

	if v4.Valid {
		v.CombiName = v4.String
	} else {
		v.CombiName = ""
	}

	return v, nil
}

func scanAssetUnits(rows *sql.Rows) ([]*schema.AssetUnit, error) {
	var err error
	var vv []*schema.AssetUnit

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.AssetUnit{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.AccountNo = v1.String
		} else {
			v.AccountNo = ""
		}

		if v2.Valid {
			v.AccountName = v2.String
		} else {
			v.AccountName = ""
		}

		if v3.Valid {
			v.CombiNo = v3.String
		} else {
			v.CombiNo = ""
		}

		if v4.Valid {
			v.CombiName = v4.String
		} else {
			v.CombiName = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceAssetUnit(v *schema.AssetUnit) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string

	v0 = v.ID
	v1 = v.AccountNo
	v2 = v.AccountName
	v3 = v.CombiNo
	v4 = v.CombiName

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
	}
}

func genericSelectAssetUnit(db db.SimpleDB, query string, args ...interface{}) (*schema.AssetUnit, error) {
	row := db.QueryRow(query, args...)
	return scanAssetUnit(row)
}

func genericSelectAssetUnits(db db.SimpleDB, query string, args ...interface{}) ([]*schema.AssetUnit, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanAssetUnits(rows)
}

func InsertAssetUnit(db db.SimpleDB, v *schema.AssetUnit) error {

	res, err := db.Exec(InsertAssetUnitStmt, sliceAssetUnit(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteAssetUnitById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteAssetUnitByIdStmt, args...)
	return err
}

func DeleteAssetUnitByAccountNo(db db.SimpleDB, accountNo string) error {
	args := []interface{}{accountNo}
	_, err := db.Exec(DeleteAssetUnitByAccountNoStmt, args...)
	return err
}

func DeleteAssetUnitByAccountNoAndCombiNo(db db.SimpleDB, accountNo string, combiNo string) error {
	args := []interface{}{accountNo, combiNo}
	_, err := db.Exec(DeleteAssetUnitByAccountNoAndCombiNoStmt, args...)
	return err
}

func UpdateAssetUnitById(db db.SimpleDB, v *schema.AssetUnit) error {
	args := sliceAssetUnit(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateAssetUnitByIdStmt, args...)
	return err
}

func UpdateAssetUnitByAccountNoAndCombiNo(db db.SimpleDB, v *schema.AssetUnit) error {
	args := sliceAssetUnit(v)
	args = append(args, v.AccountNo, v.CombiNo)
	_, err := db.Exec(UpdateAssetUnitByAccountNoAndCombiNoStmt, args...)
	return err
}

func GetAssetUnitById(db db.SimpleDB, iD int64) (*schema.AssetUnit, error) {
	args := []interface{}{iD}
	v, err := genericSelectAssetUnit(db, SelectAssetUnitByIdStmt, args...)
	return v, err
}

func GetAssetUnitByAccountNoAndCombiNo(db db.SimpleDB, accountNo string, combiNo string) (*schema.AssetUnit, error) {
	args := []interface{}{accountNo, combiNo}
	v, err := genericSelectAssetUnit(db, SelectAssetUnitByAccountNoAndCombiNoStmt, args...)
	return v, err
}

func FindAllAssetUnits(db db.SimpleDB) ([]*schema.AssetUnit, error) {
	args := []interface{}{}
	v, err := genericSelectAssetUnits(db, SelectAssetUnitStmt, args...)
	return v, err
}

func FindAllAssetUnitsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.AssetUnit, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectAssetUnits(db, SelectAssetUnitRangeStmt, args...)
	return v, err
}

func FindAssetUnitsByAccountNo(db db.SimpleDB, accountNo string) ([]*schema.AssetUnit, error) {
	args := []interface{}{accountNo}
	v, err := genericSelectAssetUnits(db, SelectAssetUnitByAccountNoStmt, args...)
	return v, err
}

func FindAssetUnitsByAccountNoInRange(db db.SimpleDB, accountNo string, limit int64, offset int64) ([]*schema.AssetUnit, error) {
	args := []interface{}{accountNo, limit, offset}
	v, err := genericSelectAssetUnits(db, SelectAssetUnitRangeByAccountNoStmt, args...)
	return v, err
}

func CountAssetUnit(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectAssetUnitCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountAssetUnitByAccountNo(db db.SimpleDB, accountNo string) (int, error) {
	var count int
	args := []interface{}{accountNo}
	row := db.QueryRow(SelectAssetUnitCountByAccountNoStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountAssetUnitByAccountNoAndCombiNo(db db.SimpleDB, accountNo string, combiNo string) (int, error) {
	var count int
	args := []interface{}{accountNo, combiNo}
	row := db.QueryRow(SelectAssetUnitCountByAccountNoAndCombiNoStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeInstrStmt = `
CREATE TABLE IF NOT EXISTS trade_instrs (
 f_id                       BIGINT PRIMARY KEY AUTO_INCREMENT
,f_msg_type                 VARCHAR(1)
,f_client_id                VARCHAR(16)
,f_parent_key               VARCHAR(64)
,f_secondary_cl_ord_id      VARCHAR(128)
,f_security_id              VARCHAR(16)
,f_symbol                   VARCHAR(16)
,f_side                     VARCHAR(1)
,f_transact_time            VARCHAR(17)
,f_order_qty                DOUBLE
,f_ord_type                 VARCHAR(1)
,f_time_in_force            VARCHAR(1)
,f_price                    DOUBLE
,f_target_strategy          VARCHAR(1)
,f_strategy_parameters_text MEDIUMTEXT
,f_market_code              VARCHAR(8)
,f_user_text                VARCHAR(512)
,f_open_close               VARCHAR(1)
,f_api_operator             VARCHAR(32)
,f_avg_px                   DOUBLE
,f_cum_amt                  DOUBLE
,f_cum_qty                  DOUBLE
,f_cum_total_fee            DOUBLE
,f_ord_status               VARCHAR(2)
,f_status_update_time       BIGINT
,f_status_update_text       MEDIUMTEXT
,f_status_kafka_offset      BIGINT
,f_message_time             BIGINT
,f_cl_ord_id                VARCHAR(128)
,f_orig_cl_ord_id           VARCHAR(128)
);
`

const InsertTradeInstrStmt = `
INSERT INTO trade_instrs (
 f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeInstrStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
`

const SelectTradeInstrRangeStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
LIMIT ? OFFSET ?
`

const SelectTradeInstrCountStmt = `
SELECT count(1)
FROM trade_instrs 
`

const SelectTradeInstrByIdStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
WHERE f_id=?
`

const UpdateTradeInstrByIdStmt = `
UPDATE trade_instrs SET 
 f_id=?
,f_msg_type=?
,f_client_id=?
,f_parent_key=?
,f_secondary_cl_ord_id=?
,f_security_id=?
,f_symbol=?
,f_side=?
,f_transact_time=?
,f_order_qty=?
,f_ord_type=?
,f_time_in_force=?
,f_price=?
,f_target_strategy=?
,f_strategy_parameters_text=?
,f_market_code=?
,f_user_text=?
,f_open_close=?
,f_api_operator=?
,f_avg_px=?
,f_cum_amt=?
,f_cum_qty=?
,f_cum_total_fee=?
,f_ord_status=?
,f_status_update_time=?
,f_status_update_text=?
,f_status_kafka_offset=?
,f_message_time=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=? 
WHERE f_id=?
`

const DeleteTradeInstrByIdStmt = `
DELETE FROM trade_instrs 
WHERE f_id=?
`

const CreateITiParentStmt = `
CREATE INDEX i_ti_parent ON trade_instrs (f_parent_key);
`

const SelectTradeInstrByParentKeyStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
WHERE f_parent_key=?
`

const SelectTradeInstrCountByParentKeyStmt = `
SELECT count(1)
FROM trade_instrs 
WHERE f_parent_key=?
`

const SelectTradeInstrRangeByParentKeyStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
WHERE f_parent_key=?
LIMIT ? OFFSET ?
`

const DeleteTradeInstrByParentKeyStmt = `
DELETE FROM trade_instrs 
WHERE f_parent_key=?
`

const CreatePkTradeInstrStmt = `
CREATE UNIQUE INDEX pk_trade_instr ON trade_instrs (f_secondary_cl_ord_id);
`

const SelectTradeInstrBySecondaryClOrdIdStmt = `
SELECT 
 f_id
,f_msg_type
,f_client_id
,f_parent_key
,f_secondary_cl_ord_id
,f_security_id
,f_symbol
,f_side
,f_transact_time
,f_order_qty
,f_ord_type
,f_time_in_force
,f_price
,f_target_strategy
,f_strategy_parameters_text
,f_market_code
,f_user_text
,f_open_close
,f_api_operator
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_ord_status
,f_status_update_time
,f_status_update_text
,f_status_kafka_offset
,f_message_time
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instrs 
WHERE f_secondary_cl_ord_id=?
`

const SelectTradeInstrCountBySecondaryClOrdIdStmt = `
SELECT count(1)
FROM trade_instrs 
WHERE f_secondary_cl_ord_id=?
`

const UpdateTradeInstrBySecondaryClOrdIdStmt = `
UPDATE trade_instrs SET 
 f_id=?
,f_msg_type=?
,f_client_id=?
,f_parent_key=?
,f_secondary_cl_ord_id=?
,f_security_id=?
,f_symbol=?
,f_side=?
,f_transact_time=?
,f_order_qty=?
,f_ord_type=?
,f_time_in_force=?
,f_price=?
,f_target_strategy=?
,f_strategy_parameters_text=?
,f_market_code=?
,f_user_text=?
,f_open_close=?
,f_api_operator=?
,f_avg_px=?
,f_cum_amt=?
,f_cum_qty=?
,f_cum_total_fee=?
,f_ord_status=?
,f_status_update_time=?
,f_status_update_text=?
,f_status_kafka_offset=?
,f_message_time=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=? 
WHERE f_secondary_cl_ord_id=?
`

const DeleteTradeInstrBySecondaryClOrdIdStmt = `
DELETE FROM trade_instrs 
WHERE f_secondary_cl_ord_id=?
`

func scanTradeInstr(row *sql.Row) (*schema.TradeInstr, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullFloat64
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullFloat64
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullFloat64
	var v20 sql.NullFloat64
	var v21 sql.NullFloat64
	var v22 sql.NullFloat64
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullString
	var v26 sql.NullInt64
	var v27 sql.NullInt64
	var v28 sql.NullString
	var v29 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
		&v7,
		&v8,
		&v9,
		&v10,
		&v11,
		&v12,
		&v13,
		&v14,
		&v15,
		&v16,
		&v17,
		&v18,
		&v19,
		&v20,
		&v21,
		&v22,
		&v23,
		&v24,
		&v25,
		&v26,
		&v27,
		&v28,
		&v29,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeInstr{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.MsgType = v1.String
	} else {
		v.MsgType = ""
	}

	if v2.Valid {
		v.ClientID = v2.String
	} else {
		v.ClientID = ""
	}

	if v3.Valid {
		v.ParentKey = v3.String
	} else {
		v.ParentKey = ""
	}

	if v4.Valid {
		v.SecondaryClOrdID = v4.String
	} else {
		v.SecondaryClOrdID = ""
	}

	if v5.Valid {
		v.SecurityID = v5.String
	} else {
		v.SecurityID = ""
	}

	if v6.Valid {
		v.Symbol = v6.String
	} else {
		v.Symbol = ""
	}

	if v7.Valid {
		v.Side = v7.String
	} else {
		v.Side = ""
	}

	if v8.Valid {
		v.TransactTime = v8.String
	} else {
		v.TransactTime = ""
	}

	if v9.Valid {
		v.OrderQty = v9.Float64
	} else {
		v.OrderQty = 0
	}

	if v10.Valid {
		v.OrdType = v10.String
	} else {
		v.OrdType = ""
	}

	if v11.Valid {
		v.TimeInForce = v11.String
	} else {
		v.TimeInForce = ""
	}

	if v12.Valid {
		v.Price = v12.Float64
	} else {
		v.Price = 0
	}

	if v13.Valid {
		v.TargetStrategy = v13.String
	} else {
		v.TargetStrategy = ""
	}

	if v14.Valid {
		v.StrategyParametersText = v14.String
	} else {
		v.StrategyParametersText = ""
	}

	if v15.Valid {
		v.MarketCode = v15.String
	} else {
		v.MarketCode = ""
	}

	if v16.Valid {
		v.UserText = v16.String
	} else {
		v.UserText = ""
	}

	if v17.Valid {
		v.OpenClose = v17.String
	} else {
		v.OpenClose = ""
	}

	if v18.Valid {
		v.ApiOperator = v18.String
	} else {
		v.ApiOperator = ""
	}

	if v19.Valid {
		v.AvgPx = v19.Float64
	} else {
		v.AvgPx = 0
	}

	if v20.Valid {
		v.CumAmt = v20.Float64
	} else {
		v.CumAmt = 0
	}

	if v21.Valid {
		v.CumQty = v21.Float64
	} else {
		v.CumQty = 0
	}

	if v22.Valid {
		v.CumTotalFee = v22.Float64
	} else {
		v.CumTotalFee = 0
	}

	if v23.Valid {
		v.OrdStatus = v23.String
	} else {
		v.OrdStatus = ""
	}

	if v24.Valid {
		v.StatusUpdateTime = v24.Int64
	} else {
		v.StatusUpdateTime = 0
	}

	if v25.Valid {
		v.StatusUpdateText = v25.String
	} else {
		v.StatusUpdateText = ""
	}

	if v26.Valid {
		v.StatusKafkaOffset = v26.Int64
	} else {
		v.StatusKafkaOffset = 0
	}

	if v27.Valid {
		v.MessageTime = v27.Int64
	} else {
		v.MessageTime = 0
	}

	if v28.Valid {
		v.ClOrdID = v28.String
	} else {
		v.ClOrdID = ""
	}

	if v29.Valid {
		v.OrigClOrdID = v29.String
	} else {
		v.OrigClOrdID = ""
	}

	return v, nil
}

func scanTradeInstrs(rows *sql.Rows) ([]*schema.TradeInstr, error) {
	var err error
	var vv []*schema.TradeInstr

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullFloat64
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullFloat64
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullFloat64
	var v20 sql.NullFloat64
	var v21 sql.NullFloat64
	var v22 sql.NullFloat64
	var v23 sql.NullString
	var v24 sql.NullInt64
	var v25 sql.NullString
	var v26 sql.NullInt64
	var v27 sql.NullInt64
	var v28 sql.NullString
	var v29 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
			&v10,
			&v11,
			&v12,
			&v13,
			&v14,
			&v15,
			&v16,
			&v17,
			&v18,
			&v19,
			&v20,
			&v21,
			&v22,
			&v23,
			&v24,
			&v25,
			&v26,
			&v27,
			&v28,
			&v29,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeInstr{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.MsgType = v1.String
		} else {
			v.MsgType = ""
		}

		if v2.Valid {
			v.ClientID = v2.String
		} else {
			v.ClientID = ""
		}

		if v3.Valid {
			v.ParentKey = v3.String
		} else {
			v.ParentKey = ""
		}

		if v4.Valid {
			v.SecondaryClOrdID = v4.String
		} else {
			v.SecondaryClOrdID = ""
		}

		if v5.Valid {
			v.SecurityID = v5.String
		} else {
			v.SecurityID = ""
		}

		if v6.Valid {
			v.Symbol = v6.String
		} else {
			v.Symbol = ""
		}

		if v7.Valid {
			v.Side = v7.String
		} else {
			v.Side = ""
		}

		if v8.Valid {
			v.TransactTime = v8.String
		} else {
			v.TransactTime = ""
		}

		if v9.Valid {
			v.OrderQty = v9.Float64
		} else {
			v.OrderQty = 0
		}

		if v10.Valid {
			v.OrdType = v10.String
		} else {
			v.OrdType = ""
		}

		if v11.Valid {
			v.TimeInForce = v11.String
		} else {
			v.TimeInForce = ""
		}

		if v12.Valid {
			v.Price = v12.Float64
		} else {
			v.Price = 0
		}

		if v13.Valid {
			v.TargetStrategy = v13.String
		} else {
			v.TargetStrategy = ""
		}

		if v14.Valid {
			v.StrategyParametersText = v14.String
		} else {
			v.StrategyParametersText = ""
		}

		if v15.Valid {
			v.MarketCode = v15.String
		} else {
			v.MarketCode = ""
		}

		if v16.Valid {
			v.UserText = v16.String
		} else {
			v.UserText = ""
		}

		if v17.Valid {
			v.OpenClose = v17.String
		} else {
			v.OpenClose = ""
		}

		if v18.Valid {
			v.ApiOperator = v18.String
		} else {
			v.ApiOperator = ""
		}

		if v19.Valid {
			v.AvgPx = v19.Float64
		} else {
			v.AvgPx = 0
		}

		if v20.Valid {
			v.CumAmt = v20.Float64
		} else {
			v.CumAmt = 0
		}

		if v21.Valid {
			v.CumQty = v21.Float64
		} else {
			v.CumQty = 0
		}

		if v22.Valid {
			v.CumTotalFee = v22.Float64
		} else {
			v.CumTotalFee = 0
		}

		if v23.Valid {
			v.OrdStatus = v23.String
		} else {
			v.OrdStatus = ""
		}

		if v24.Valid {
			v.StatusUpdateTime = v24.Int64
		} else {
			v.StatusUpdateTime = 0
		}

		if v25.Valid {
			v.StatusUpdateText = v25.String
		} else {
			v.StatusUpdateText = ""
		}

		if v26.Valid {
			v.StatusKafkaOffset = v26.Int64
		} else {
			v.StatusKafkaOffset = 0
		}

		if v27.Valid {
			v.MessageTime = v27.Int64
		} else {
			v.MessageTime = 0
		}

		if v28.Valid {
			v.ClOrdID = v28.String
		} else {
			v.ClOrdID = ""
		}

		if v29.Valid {
			v.OrigClOrdID = v29.String
		} else {
			v.OrigClOrdID = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeInstr(v *schema.TradeInstr) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 float64
	var v10 string
	var v11 string
	var v12 float64
	var v13 string
	var v14 string
	var v15 string
	var v16 string
	var v17 string
	var v18 string
	var v19 float64
	var v20 float64
	var v21 float64
	var v22 float64
	var v23 string
	var v24 int64
	var v25 string
	var v26 int64
	var v27 int64
	var v28 string
	var v29 string

	v0 = v.ID
	v1 = v.MsgType
	v2 = v.ClientID
	v3 = v.ParentKey
	v4 = v.SecondaryClOrdID
	v5 = v.SecurityID
	v6 = v.Symbol
	v7 = v.Side
	v8 = v.TransactTime
	v9 = v.OrderQty
	v10 = v.OrdType
	v11 = v.TimeInForce
	v12 = v.Price
	v13 = v.TargetStrategy
	v14 = v.StrategyParametersText
	v15 = v.MarketCode
	v16 = v.UserText
	v17 = v.OpenClose
	v18 = v.ApiOperator
	v19 = v.AvgPx
	v20 = v.CumAmt
	v21 = v.CumQty
	v22 = v.CumTotalFee
	v23 = v.OrdStatus
	v24 = v.StatusUpdateTime
	v25 = v.StatusUpdateText
	v26 = v.StatusKafkaOffset
	v27 = v.MessageTime
	v28 = v.ClOrdID
	v29 = v.OrigClOrdID

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
		v9,
		v10,
		v11,
		v12,
		v13,
		v14,
		v15,
		v16,
		v17,
		v18,
		v19,
		v20,
		v21,
		v22,
		v23,
		v24,
		v25,
		v26,
		v27,
		v28,
		v29,
	}
}

func genericSelectTradeInstr(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeInstr, error) {
	row := db.QueryRow(query, args...)
	return scanTradeInstr(row)
}

func genericSelectTradeInstrs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeInstr, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeInstrs(rows)
}

func InsertTradeInstr(db db.SimpleDB, v *schema.TradeInstr) error {

	res, err := db.Exec(InsertTradeInstrStmt, sliceTradeInstr(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeInstrById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeInstrByIdStmt, args...)
	return err
}

func DeleteTradeInstrByParentKey(db db.SimpleDB, parentKey string) error {
	args := []interface{}{parentKey}
	_, err := db.Exec(DeleteTradeInstrByParentKeyStmt, args...)
	return err
}

func DeleteTradeInstrBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) error {
	args := []interface{}{secondaryClOrdID}
	_, err := db.Exec(DeleteTradeInstrBySecondaryClOrdIdStmt, args...)
	return err
}

func UpdateTradeInstrById(db db.SimpleDB, v *schema.TradeInstr) error {
	args := sliceTradeInstr(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeInstrByIdStmt, args...)
	return err
}

func UpdateTradeInstrBySecondaryClOrdId(db db.SimpleDB, v *schema.TradeInstr) error {
	args := sliceTradeInstr(v)
	args = append(args, v.SecondaryClOrdID)
	_, err := db.Exec(UpdateTradeInstrBySecondaryClOrdIdStmt, args...)
	return err
}

func GetTradeInstrById(db db.SimpleDB, iD int64) (*schema.TradeInstr, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeInstr(db, SelectTradeInstrByIdStmt, args...)
	return v, err
}

func GetTradeInstrBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) (*schema.TradeInstr, error) {
	args := []interface{}{secondaryClOrdID}
	v, err := genericSelectTradeInstr(db, SelectTradeInstrBySecondaryClOrdIdStmt, args...)
	return v, err
}

func FindAllTradeInstrs(db db.SimpleDB) ([]*schema.TradeInstr, error) {
	args := []interface{}{}
	v, err := genericSelectTradeInstrs(db, SelectTradeInstrStmt, args...)
	return v, err
}

func FindAllTradeInstrsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeInstr, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeInstrs(db, SelectTradeInstrRangeStmt, args...)
	return v, err
}

func FindTradeInstrsByParentKey(db db.SimpleDB, parentKey string) ([]*schema.TradeInstr, error) {
	args := []interface{}{parentKey}
	v, err := genericSelectTradeInstrs(db, SelectTradeInstrByParentKeyStmt, args...)
	return v, err
}

func FindTradeInstrsByParentKeyInRange(db db.SimpleDB, parentKey string, limit int64, offset int64) ([]*schema.TradeInstr, error) {
	args := []interface{}{parentKey, limit, offset}
	v, err := genericSelectTradeInstrs(db, SelectTradeInstrRangeByParentKeyStmt, args...)
	return v, err
}

func CountTradeInstr(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeInstrCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeInstrByParentKey(db db.SimpleDB, parentKey string) (int, error) {
	var count int
	args := []interface{}{parentKey}
	row := db.QueryRow(SelectTradeInstrCountByParentKeyStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTradeInstrBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) (int, error) {
	var count int
	args := []interface{}{secondaryClOrdID}
	row := db.QueryRow(SelectTradeInstrCountBySecondaryClOrdIdStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeInstrRespStmt = `
CREATE TABLE IF NOT EXISTS trade_instr_resps (
 f_id                  BIGINT PRIMARY KEY AUTO_INCREMENT
,f_secondary_cl_ord_id VARCHAR(128)
,f_message_time        BIGINT
,f_msg_type            VARCHAR(1)
,f_transact_time       VARCHAR(17)
,f_status_kafka_offset BIGINT
,f_ord_status          VARCHAR(2)
,f_status_update_text  MEDIUMTEXT
,f_avg_px              DOUBLE
,f_cum_amt             DOUBLE
,f_cum_qty             DOUBLE
,f_cum_total_fee       DOUBLE
,f_cl_ord_id           VARCHAR(128)
,f_orig_cl_ord_id      VARCHAR(128)
);
`

const InsertTradeInstrRespStmt = `
INSERT INTO trade_instr_resps (
 f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeInstrRespStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
`

const SelectTradeInstrRespRangeStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
LIMIT ? OFFSET ?
`

const SelectTradeInstrRespCountStmt = `
SELECT count(1)
FROM trade_instr_resps 
`

const SelectTradeInstrRespByIdStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_id=?
`

const UpdateTradeInstrRespByIdStmt = `
UPDATE trade_instr_resps SET 
 f_id=?
,f_secondary_cl_ord_id=?
,f_message_time=?
,f_msg_type=?
,f_transact_time=?
,f_status_kafka_offset=?
,f_ord_status=?
,f_status_update_text=?
,f_avg_px=?
,f_cum_amt=?
,f_cum_qty=?
,f_cum_total_fee=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=? 
WHERE f_id=?
`

const DeleteTradeInstrRespByIdStmt = `
DELETE FROM trade_instr_resps 
WHERE f_id=?
`

const CreateITradeInstrRespStmt = `
CREATE INDEX i_trade_instr_resp ON trade_instr_resps (f_secondary_cl_ord_id);
`

const SelectTradeInstrRespBySecondaryClOrdIdStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
`

const SelectTradeInstrRespCountBySecondaryClOrdIdStmt = `
SELECT count(1)
FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
`

const SelectTradeInstrRespRangeBySecondaryClOrdIdStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
LIMIT ? OFFSET ?
`

const DeleteTradeInstrRespBySecondaryClOrdIdStmt = `
DELETE FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
`

const CreatePkTirStmt = `
CREATE UNIQUE INDEX pk_tir ON trade_instr_resps (f_secondary_cl_ord_id,f_status_kafka_offset);
`

const SelectTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
AND f_status_kafka_offset=?
`

const SelectTradeInstrRespCountBySecondaryClOrdIdAndStatusKafkaOffsetStmt = `
SELECT count(1)
FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
AND f_status_kafka_offset=?
`

const UpdateTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt = `
UPDATE trade_instr_resps SET 
 f_id=?
,f_secondary_cl_ord_id=?
,f_message_time=?
,f_msg_type=?
,f_transact_time=?
,f_status_kafka_offset=?
,f_ord_status=?
,f_status_update_text=?
,f_avg_px=?
,f_cum_amt=?
,f_cum_qty=?
,f_cum_total_fee=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=? 
WHERE f_secondary_cl_ord_id=?
AND f_status_kafka_offset=?
`

const DeleteTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt = `
DELETE FROM trade_instr_resps 
WHERE f_secondary_cl_ord_id=?
AND f_status_kafka_offset=?
`

const CreateITirOStmt = `
CREATE INDEX i_tir_o ON trade_instr_resps (f_status_kafka_offset);
`

const SelectTradeInstrRespByStatusKafkaOffsetStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_status_kafka_offset=?
`

const SelectTradeInstrRespCountByStatusKafkaOffsetStmt = `
SELECT count(1)
FROM trade_instr_resps 
WHERE f_status_kafka_offset=?
`

const SelectTradeInstrRespRangeByStatusKafkaOffsetStmt = `
SELECT 
 f_id
,f_secondary_cl_ord_id
,f_message_time
,f_msg_type
,f_transact_time
,f_status_kafka_offset
,f_ord_status
,f_status_update_text
,f_avg_px
,f_cum_amt
,f_cum_qty
,f_cum_total_fee
,f_cl_ord_id
,f_orig_cl_ord_id
FROM trade_instr_resps 
WHERE f_status_kafka_offset=?
LIMIT ? OFFSET ?
`

const DeleteTradeInstrRespByStatusKafkaOffsetStmt = `
DELETE FROM trade_instr_resps 
WHERE f_status_kafka_offset=?
`

func scanTradeInstrResp(row *sql.Row) (*schema.TradeInstrResp, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullInt64
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullFloat64
	var v11 sql.NullFloat64
	var v12 sql.NullString
	var v13 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
		&v7,
		&v8,
		&v9,
		&v10,
		&v11,
		&v12,
		&v13,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeInstrResp{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.SecondaryClOrdID = v1.String
	} else {
		v.SecondaryClOrdID = ""
	}

	if v2.Valid {
		v.MessageTime = v2.Int64
	} else {
		v.MessageTime = 0
	}

	if v3.Valid {
		v.MsgType = v3.String
	} else {
		v.MsgType = ""
	}

	if v4.Valid {
		v.TransactTime = v4.String
	} else {
		v.TransactTime = ""
	}

	if v5.Valid {
		v.StatusKafkaOffset = v5.Int64
	} else {
		v.StatusKafkaOffset = 0
	}

	if v6.Valid {
		v.OrdStatus = v6.String
	} else {
		v.OrdStatus = ""
	}

	if v7.Valid {
		v.StatusUpdateText = v7.String
	} else {
		v.StatusUpdateText = ""
	}

	if v8.Valid {
		v.AvgPx = v8.Float64
	} else {
		v.AvgPx = 0
	}

	if v9.Valid {
		v.CumAmt = v9.Float64
	} else {
		v.CumAmt = 0
	}

	if v10.Valid {
		v.CumQty = v10.Float64
	} else {
		v.CumQty = 0
	}

	if v11.Valid {
		v.CumTotalFee = v11.Float64
	} else {
		v.CumTotalFee = 0
	}

	if v12.Valid {
		v.ClOrdID = v12.String
	} else {
		v.ClOrdID = ""
	}

	if v13.Valid {
		v.OrigClOrdID = v13.String
	} else {
		v.OrigClOrdID = ""
	}

	return v, nil
}

func scanTradeInstrResps(rows *sql.Rows) ([]*schema.TradeInstrResp, error) {
	var err error
	var vv []*schema.TradeInstrResp

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullInt64
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullFloat64
	var v11 sql.NullFloat64
	var v12 sql.NullString
	var v13 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
			&v7,
			&v8,
			&v9,
			&v10,
			&v11,
			&v12,
			&v13,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeInstrResp{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.SecondaryClOrdID = v1.String
		} else {
			v.SecondaryClOrdID = ""
		}

		if v2.Valid {
			v.MessageTime = v2.Int64
		} else {
			v.MessageTime = 0
		}

		if v3.Valid {
			v.MsgType = v3.String
		} else {
			v.MsgType = ""
		}

		if v4.Valid {
			v.TransactTime = v4.String
		} else {
			v.TransactTime = ""
		}

		if v5.Valid {
			v.StatusKafkaOffset = v5.Int64
		} else {
			v.StatusKafkaOffset = 0
		}

		if v6.Valid {
			v.OrdStatus = v6.String
		} else {
			v.OrdStatus = ""
		}

		if v7.Valid {
			v.StatusUpdateText = v7.String
		} else {
			v.StatusUpdateText = ""
		}

		if v8.Valid {
			v.AvgPx = v8.Float64
		} else {
			v.AvgPx = 0
		}

		if v9.Valid {
			v.CumAmt = v9.Float64
		} else {
			v.CumAmt = 0
		}

		if v10.Valid {
			v.CumQty = v10.Float64
		} else {
			v.CumQty = 0
		}

		if v11.Valid {
			v.CumTotalFee = v11.Float64
		} else {
			v.CumTotalFee = 0
		}

		if v12.Valid {
			v.ClOrdID = v12.String
		} else {
			v.ClOrdID = ""
		}

		if v13.Valid {
			v.OrigClOrdID = v13.String
		} else {
			v.OrigClOrdID = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeInstrResp(v *schema.TradeInstrResp) []interface{} {
	var v0 int64
	var v1 string
	var v2 int64
	var v3 string
	var v4 string
	var v5 int64
	var v6 string
	var v7 string
	var v8 float64
	var v9 float64
	var v10 float64
	var v11 float64
	var v12 string
	var v13 string

	v0 = v.ID
	v1 = v.SecondaryClOrdID
	v2 = v.MessageTime
	v3 = v.MsgType
	v4 = v.TransactTime
	v5 = v.StatusKafkaOffset
	v6 = v.OrdStatus
	v7 = v.StatusUpdateText
	v8 = v.AvgPx
	v9 = v.CumAmt
	v10 = v.CumQty
	v11 = v.CumTotalFee
	v12 = v.ClOrdID
	v13 = v.OrigClOrdID

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
		v7,
		v8,
		v9,
		v10,
		v11,
		v12,
		v13,
	}
}

func genericSelectTradeInstrResp(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeInstrResp, error) {
	row := db.QueryRow(query, args...)
	return scanTradeInstrResp(row)
}

func genericSelectTradeInstrResps(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeInstrResp, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeInstrResps(rows)
}

func InsertTradeInstrResp(db db.SimpleDB, v *schema.TradeInstrResp) error {

	res, err := db.Exec(InsertTradeInstrRespStmt, sliceTradeInstrResp(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeInstrRespById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeInstrRespByIdStmt, args...)
	return err
}

func DeleteTradeInstrRespBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) error {
	args := []interface{}{secondaryClOrdID}
	_, err := db.Exec(DeleteTradeInstrRespBySecondaryClOrdIdStmt, args...)
	return err
}

func DeleteTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffset(db db.SimpleDB, secondaryClOrdID string, statusKafkaOffset int64) error {
	args := []interface{}{secondaryClOrdID, statusKafkaOffset}
	_, err := db.Exec(DeleteTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt, args...)
	return err
}

func DeleteTradeInstrRespByStatusKafkaOffset(db db.SimpleDB, statusKafkaOffset int64) error {
	args := []interface{}{statusKafkaOffset}
	_, err := db.Exec(DeleteTradeInstrRespByStatusKafkaOffsetStmt, args...)
	return err
}

func UpdateTradeInstrRespById(db db.SimpleDB, v *schema.TradeInstrResp) error {
	args := sliceTradeInstrResp(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeInstrRespByIdStmt, args...)
	return err
}

func UpdateTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffset(db db.SimpleDB, v *schema.TradeInstrResp) error {
	args := sliceTradeInstrResp(v)
	args = append(args, v.SecondaryClOrdID, v.StatusKafkaOffset)
	_, err := db.Exec(UpdateTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt, args...)
	return err
}

func GetTradeInstrRespById(db db.SimpleDB, iD int64) (*schema.TradeInstrResp, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeInstrResp(db, SelectTradeInstrRespByIdStmt, args...)
	return v, err
}

func GetTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffset(db db.SimpleDB, secondaryClOrdID string, statusKafkaOffset int64) (*schema.TradeInstrResp, error) {
	args := []interface{}{secondaryClOrdID, statusKafkaOffset}
	v, err := genericSelectTradeInstrResp(db, SelectTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffsetStmt, args...)
	return v, err
}

func FindAllTradeInstrResps(db db.SimpleDB) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespStmt, args...)
	return v, err
}

func FindAllTradeInstrRespsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespRangeStmt, args...)
	return v, err
}

func FindTradeInstrRespsBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{secondaryClOrdID}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespBySecondaryClOrdIdStmt, args...)
	return v, err
}

func FindTradeInstrRespsBySecondaryClOrdIdInRange(db db.SimpleDB, secondaryClOrdID string, limit int64, offset int64) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{secondaryClOrdID, limit, offset}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespRangeBySecondaryClOrdIdStmt, args...)
	return v, err
}

func FindTradeInstrRespsByStatusKafkaOffset(db db.SimpleDB, statusKafkaOffset int64) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{statusKafkaOffset}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespByStatusKafkaOffsetStmt, args...)
	return v, err
}

func FindTradeInstrRespsByStatusKafkaOffsetInRange(db db.SimpleDB, statusKafkaOffset int64, limit int64, offset int64) ([]*schema.TradeInstrResp, error) {
	args := []interface{}{statusKafkaOffset, limit, offset}
	v, err := genericSelectTradeInstrResps(db, SelectTradeInstrRespRangeByStatusKafkaOffsetStmt, args...)
	return v, err
}

func CountTradeInstrResp(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeInstrRespCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeInstrRespBySecondaryClOrdId(db db.SimpleDB, secondaryClOrdID string) (int, error) {
	var count int
	args := []interface{}{secondaryClOrdID}
	row := db.QueryRow(SelectTradeInstrRespCountBySecondaryClOrdIdStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTradeInstrRespBySecondaryClOrdIdAndStatusKafkaOffset(db db.SimpleDB, secondaryClOrdID string, statusKafkaOffset int64) (int, error) {
	var count int
	args := []interface{}{secondaryClOrdID, statusKafkaOffset}
	row := db.QueryRow(SelectTradeInstrRespCountBySecondaryClOrdIdAndStatusKafkaOffsetStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTradeInstrRespByStatusKafkaOffset(db db.SimpleDB, statusKafkaOffset int64) (int, error) {
	var count int
	args := []interface{}{statusKafkaOffset}
	row := db.QueryRow(SelectTradeInstrRespCountByStatusKafkaOffsetStmt, args...)
	err := row.Scan(&count)
	return count, err
}
