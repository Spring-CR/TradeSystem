package app_store

// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.

import (
	"database/sql"
	"github.com/linchunquan/sqlgen/db"
	"rhino-core/schema"
)

const CreateTradeOrderStmt = `
CREATE TABLE IF NOT EXISTS trade_orders (
 f_id                         BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code                VARCHAR(32)
,f_business_code              VARCHAR(32)
,f_cl_group_ord_id            VARCHAR(188)
,f_cl_ord_id                  VARCHAR(188)
,f_account                    VARCHAR(64)
,f_handl_inst                 VARCHAR(2)
,f_app_ord_id                 VARCHAR(188)
,f_ord_id                     VARCHAR(188)
,f_parent_cl_ord_id           VARCHAR(188)
,f_is_direct_ord              BOOLEAN
,f_is_alg_ord                 BOOLEAN
,f_is_sub_alg_ord             BOOLEAN
,f_is_instr_ord               BOOLEAN
,f_is_sub_instr_ord           BOOLEAN
,f_is_cross_date_ord          BOOLEAN
,f_is_sub_cross_date_ord      BOOLEAN
,f_min_qty                    DOUBLE
,f_security_exchange          VARCHAR(8)
,f_security_exchange_region   VARCHAR(4)
,f_symbol                     VARCHAR(64)
,f_symbol_sfx                 VARCHAR(8)
,f_security_id                VARCHAR(64)
,f_id_source                  VARCHAR(2)
,f_security_type              VARCHAR(2)
,f_side                       VARCHAR(2)
,f_transact_time              BIGINT
,f_trade_date                 BIGINT
,f_order_qty                  DOUBLE
,f_cash_order_qty             DOUBLE
,f_ord_type                   VARCHAR(2)
,f_price                      DOUBLE
,f_currency                   VARCHAR(4)
,f_effective_time             BIGINT
,f_expire_date                BIGINT
,f_expire_time                BIGINT
,f_open_close                 VARCHAR(2)
,f_contract_multiplier        DOUBLE
,f_ord_changed_count          INTEGER
,f_ord_cancel_count           INTEGER
,f_extend_attr                MEDIUMTEXT
,f_alg_params                 MEDIUMTEXT
,f_alg_name                   VARCHAR(32)
,f_ord_creator                VARCHAR(64)
,f_ord_create_time            BIGINT
,f_ord_draft_update_user      VARCHAR(64)
,f_ord_draft_update_time      BIGINT
,f_ord_draft_del_flag         INTEGER
,f_ord_draft_del_user         VARCHAR(64)
,f_ord_draft_del_time         BIGINT
,f_ord_exec_user_scope        VARCHAR(512)
,f_ord_exec_user              VARCHAR(64)
,f_ord_status_update_time     BIGINT
,f_ord_status                 VARCHAR(2)
,f_review_flag                VARCHAR(2)
,f_reviewer_scope             VARCHAR(512)
,f_reviewer                   VARCHAR(512)
,f_approve_status             INTEGER
,f_review_time                BIGINT
,f_order_submit_fail_reason   VARCHAR(512)
,f_push_in_queue_before_trade BOOLEAN
,f_latest_action_type         VARCHAR(2)
,f_channel_code               VARCHAR(32)
,f_db_insert_time             BIGINT
,f_msg_seq                    BIGINT
,f_worker_affinity            INTEGER
,f_quota_validate_time        BIGINT
);
`

const InsertTradeOrderStmt = `
INSERT INTO trade_orders (
 f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeOrderStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
`

const SelectTradeOrderRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
LIMIT ? OFFSET ?
`

const SelectTradeOrderCountStmt = `
SELECT count(1)
FROM trade_orders 
`

const SelectTradeOrderByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
WHERE f_id=?
`

const UpdateTradeOrderByIdStmt = `
UPDATE trade_orders SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_cl_group_ord_id=?
,f_cl_ord_id=?
,f_account=?
,f_handl_inst=?
,f_app_ord_id=?
,f_ord_id=?
,f_parent_cl_ord_id=?
,f_is_direct_ord=?
,f_is_alg_ord=?
,f_is_sub_alg_ord=?
,f_is_instr_ord=?
,f_is_sub_instr_ord=?
,f_is_cross_date_ord=?
,f_is_sub_cross_date_ord=?
,f_min_qty=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_transact_time=?
,f_trade_date=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_date=?
,f_expire_time=?
,f_open_close=?
,f_contract_multiplier=?
,f_ord_changed_count=?
,f_ord_cancel_count=?
,f_extend_attr=?
,f_alg_params=?
,f_alg_name=?
,f_ord_creator=?
,f_ord_create_time=?
,f_ord_draft_update_user=?
,f_ord_draft_update_time=?
,f_ord_draft_del_flag=?
,f_ord_draft_del_user=?
,f_ord_draft_del_time=?
,f_ord_exec_user_scope=?
,f_ord_exec_user=?
,f_ord_status_update_time=?
,f_ord_status=?
,f_review_flag=?
,f_reviewer_scope=?
,f_reviewer=?
,f_approve_status=?
,f_review_time=?
,f_order_submit_fail_reason=?
,f_push_in_queue_before_trade=?
,f_latest_action_type=?
,f_channel_code=?
,f_db_insert_time=?
,f_msg_seq=?
,f_worker_affinity=?
,f_quota_validate_time=? 
WHERE f_id=?
`

const DeleteTradeOrderByIdStmt = `
DELETE FROM trade_orders 
WHERE f_id=?
`

const CreateIdxToClordidStmt = `
CREATE INDEX idx_to_clordid ON trade_orders (f_cl_ord_id);
`

const SelectTradeOrderByClOrdIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
WHERE f_cl_ord_id=?
`

const SelectTradeOrderCountByClOrdIdStmt = `
SELECT count(1)
FROM trade_orders 
WHERE f_cl_ord_id=?
`

const SelectTradeOrderRangeByClOrdIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
WHERE f_cl_ord_id=?
LIMIT ? OFFSET ?
`

const DeleteTradeOrderByClOrdIdStmt = `
DELETE FROM trade_orders 
WHERE f_cl_ord_id=?
`

const CreateUqAppordidStmt = `
CREATE UNIQUE INDEX uq_appordid ON trade_orders (f_app_ord_id);
`

const SelectTradeOrderByAppOrdIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_cl_ord_id
,f_account
,f_handl_inst
,f_app_ord_id
,f_ord_id
,f_parent_cl_ord_id
,f_is_direct_ord
,f_is_alg_ord
,f_is_sub_alg_ord
,f_is_instr_ord
,f_is_sub_instr_ord
,f_is_cross_date_ord
,f_is_sub_cross_date_ord
,f_min_qty
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_transact_time
,f_trade_date
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_date
,f_expire_time
,f_open_close
,f_contract_multiplier
,f_ord_changed_count
,f_ord_cancel_count
,f_extend_attr
,f_alg_params
,f_alg_name
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_flag
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_ord_status
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
,f_order_submit_fail_reason
,f_push_in_queue_before_trade
,f_latest_action_type
,f_channel_code
,f_db_insert_time
,f_msg_seq
,f_worker_affinity
,f_quota_validate_time
FROM trade_orders 
WHERE f_app_ord_id=?
`

const SelectTradeOrderCountByAppOrdIdStmt = `
SELECT count(1)
FROM trade_orders 
WHERE f_app_ord_id=?
`

const UpdateTradeOrderByAppOrdIdStmt = `
UPDATE trade_orders SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_cl_group_ord_id=?
,f_cl_ord_id=?
,f_account=?
,f_handl_inst=?
,f_app_ord_id=?
,f_ord_id=?
,f_parent_cl_ord_id=?
,f_is_direct_ord=?
,f_is_alg_ord=?
,f_is_sub_alg_ord=?
,f_is_instr_ord=?
,f_is_sub_instr_ord=?
,f_is_cross_date_ord=?
,f_is_sub_cross_date_ord=?
,f_min_qty=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_transact_time=?
,f_trade_date=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_date=?
,f_expire_time=?
,f_open_close=?
,f_contract_multiplier=?
,f_ord_changed_count=?
,f_ord_cancel_count=?
,f_extend_attr=?
,f_alg_params=?
,f_alg_name=?
,f_ord_creator=?
,f_ord_create_time=?
,f_ord_draft_update_user=?
,f_ord_draft_update_time=?
,f_ord_draft_del_flag=?
,f_ord_draft_del_user=?
,f_ord_draft_del_time=?
,f_ord_exec_user_scope=?
,f_ord_exec_user=?
,f_ord_status_update_time=?
,f_ord_status=?
,f_review_flag=?
,f_reviewer_scope=?
,f_reviewer=?
,f_approve_status=?
,f_review_time=?
,f_order_submit_fail_reason=?
,f_push_in_queue_before_trade=?
,f_latest_action_type=?
,f_channel_code=?
,f_db_insert_time=?
,f_msg_seq=?
,f_worker_affinity=?
,f_quota_validate_time=? 
WHERE f_app_ord_id=?
`

const DeleteTradeOrderByAppOrdIdStmt = `
DELETE FROM trade_orders 
WHERE f_app_ord_id=?
`

func scanTradeOrder(row *sql.Row) (*schema.TradeOrder, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullBool
	var v11 sql.NullBool
	var v12 sql.NullBool
	var v13 sql.NullBool
	var v14 sql.NullBool
	var v15 sql.NullBool
	var v16 sql.NullBool
	var v17 sql.NullFloat64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullInt64
	var v27 sql.NullInt64
	var v28 sql.NullFloat64
	var v29 sql.NullFloat64
	var v30 sql.NullString
	var v31 sql.NullFloat64
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullInt64
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullFloat64
	var v38 sql.NullInt64
	var v39 sql.NullInt64
	var v40 sql.NullString
	var v41 sql.NullString
	var v42 sql.NullString
	var v43 sql.NullString
	var v44 sql.NullInt64
	var v45 sql.NullString
	var v46 sql.NullInt64
	var v47 sql.NullInt64
	var v48 sql.NullString
	var v49 sql.NullInt64
	var v50 sql.NullString
	var v51 sql.NullString
	var v52 sql.NullInt64
	var v53 sql.NullString
	var v54 sql.NullString
	var v55 sql.NullString
	var v56 sql.NullString
	var v57 sql.NullInt64
	var v58 sql.NullInt64
	var v59 sql.NullString
	var v60 sql.NullBool
	var v61 sql.NullString
	var v62 sql.NullString
	var v63 sql.NullInt64
	var v64 sql.NullInt64
	var v65 sql.NullInt64
	var v66 sql.NullInt64

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
		&v64,
		&v65,
		&v66,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeOrder{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.SystemCode = v1.String
	} else {
		v.SystemCode = ""
	}

	if v2.Valid {
		v.BusinessCode = v2.String
	} else {
		v.BusinessCode = ""
	}

	if v3.Valid {
		v.ClGroupOrdID = v3.String
	} else {
		v.ClGroupOrdID = ""
	}

	if v4.Valid {
		v.ClOrdID = v4.String
	} else {
		v.ClOrdID = ""
	}

	if v5.Valid {
		v.Account = v5.String
	} else {
		v.Account = ""
	}

	if v6.Valid {
		v.HandlInst = v6.String
	} else {
		v.HandlInst = ""
	}

	if v7.Valid {
		v.AppOrdID = v7.String
	} else {
		v.AppOrdID = ""
	}

	if v8.Valid {
		v.OrdID = v8.String
	} else {
		v.OrdID = ""
	}

	if v9.Valid {
		v.ParentClOrdID = v9.String
	} else {
		v.ParentClOrdID = ""
	}

	if v10.Valid {
		v.IsDirectOrd = v10.Bool
	} else {
		v.IsDirectOrd = false
	}

	if v11.Valid {
		v.IsAlgOrd = v11.Bool
	} else {
		v.IsAlgOrd = false
	}

	if v12.Valid {
		v.IsSubAlgOrd = v12.Bool
	} else {
		v.IsSubAlgOrd = false
	}

	if v13.Valid {
		v.IsInstrOrd = v13.Bool
	} else {
		v.IsInstrOrd = false
	}

	if v14.Valid {
		v.IsSubInstrOrd = v14.Bool
	} else {
		v.IsSubInstrOrd = false
	}

	if v15.Valid {
		v.IsCrossDateOrd = v15.Bool
	} else {
		v.IsCrossDateOrd = false
	}

	if v16.Valid {
		v.IsSubCrossDateOrd = v16.Bool
	} else {
		v.IsSubCrossDateOrd = false
	}

	if v17.Valid {
		v.MinQty = v17.Float64
	} else {
		v.MinQty = 0
	}

	if v18.Valid {
		v.SecurityExchange = v18.String
	} else {
		v.SecurityExchange = ""
	}

	if v19.Valid {
		v.SecurityExchangeRegion = v19.String
	} else {
		v.SecurityExchangeRegion = ""
	}

	if v20.Valid {
		v.Symbol = v20.String
	} else {
		v.Symbol = ""
	}

	if v21.Valid {
		v.SymbolSfx = v21.String
	} else {
		v.SymbolSfx = ""
	}

	if v22.Valid {
		v.SecurityID = v22.String
	} else {
		v.SecurityID = ""
	}

	if v23.Valid {
		v.IDSource = v23.String
	} else {
		v.IDSource = ""
	}

	if v24.Valid {
		v.SecurityType = v24.String
	} else {
		v.SecurityType = ""
	}

	if v25.Valid {
		v.Side = v25.String
	} else {
		v.Side = ""
	}

	if v26.Valid {
		v.TransactTime = v26.Int64
	} else {
		v.TransactTime = 0
	}

	if v27.Valid {
		v.TradeDate = v27.Int64
	} else {
		v.TradeDate = 0
	}

	if v28.Valid {
		v.OrderQty = v28.Float64
	} else {
		v.OrderQty = 0
	}

	if v29.Valid {
		v.CashOrderQty = v29.Float64
	} else {
		v.CashOrderQty = 0
	}

	if v30.Valid {
		v.OrdType = v30.String
	} else {
		v.OrdType = ""
	}

	if v31.Valid {
		v.Price = v31.Float64
	} else {
		v.Price = 0
	}

	if v32.Valid {
		v.Currency = v32.String
	} else {
		v.Currency = ""
	}

	if v33.Valid {
		v.EffectiveTime = v33.Int64
	} else {
		v.EffectiveTime = 0
	}

	if v34.Valid {
		v.ExpireDate = v34.Int64
	} else {
		v.ExpireDate = 0
	}

	if v35.Valid {
		v.ExpireTime = v35.Int64
	} else {
		v.ExpireTime = 0
	}

	if v36.Valid {
		v.OpenClose = v36.String
	} else {
		v.OpenClose = ""
	}

	if v37.Valid {
		v.ContractMultiplier = v37.Float64
	} else {
		v.ContractMultiplier = 0
	}

	if v38.Valid {
		v.OrdChangedCount = int(v38.Int64)
	} else {
		v.OrdChangedCount = 0
	}

	if v39.Valid {
		v.OrdCancelCount = int(v39.Int64)
	} else {
		v.OrdCancelCount = 0
	}

	if v40.Valid {
		v.ExtendAttr = v40.String
	} else {
		v.ExtendAttr = ""
	}

	if v41.Valid {
		v.AlgParams = v41.String
	} else {
		v.AlgParams = ""
	}

	if v42.Valid {
		v.AlgName = v42.String
	} else {
		v.AlgName = ""
	}

	if v43.Valid {
		v.OrdCreator = v43.String
	} else {
		v.OrdCreator = ""
	}

	if v44.Valid {
		v.OrdCreateTime = v44.Int64
	} else {
		v.OrdCreateTime = 0
	}

	if v45.Valid {
		v.OrdDraftUpdateUser = v45.String
	} else {
		v.OrdDraftUpdateUser = ""
	}

	if v46.Valid {
		v.OrdDraftUpdateTime = v46.Int64
	} else {
		v.OrdDraftUpdateTime = 0
	}

	if v47.Valid {
		v.OrdDraftDelFlag = int(v47.Int64)
	} else {
		v.OrdDraftDelFlag = 0
	}

	if v48.Valid {
		v.OrdDraftDelUser = v48.String
	} else {
		v.OrdDraftDelUser = ""
	}

	if v49.Valid {
		v.OrdDraftDelTime = v49.Int64
	} else {
		v.OrdDraftDelTime = 0
	}

	if v50.Valid {
		v.OrdExecUserScope = v50.String
	} else {
		v.OrdExecUserScope = ""
	}

	if v51.Valid {
		v.OrdExecUser = v51.String
	} else {
		v.OrdExecUser = ""
	}

	if v52.Valid {
		v.OrdStatusUpdateTime = v52.Int64
	} else {
		v.OrdStatusUpdateTime = 0
	}

	if v53.Valid {
		v.OrdStatus = v53.String
	} else {
		v.OrdStatus = ""
	}

	if v54.Valid {
		v.ReviewFlag = v54.String
	} else {
		v.ReviewFlag = ""
	}

	if v55.Valid {
		v.ReviewerScope = v55.String
	} else {
		v.ReviewerScope = ""
	}

	if v56.Valid {
		v.Reviewer = v56.String
	} else {
		v.Reviewer = ""
	}

	if v57.Valid {
		v.ApproveStatus = int(v57.Int64)
	} else {
		v.ApproveStatus = 0
	}

	if v58.Valid {
		v.ReviewTime = v58.Int64
	} else {
		v.ReviewTime = 0
	}

	if v59.Valid {
		v.OrderSubmitFailReason = v59.String
	} else {
		v.OrderSubmitFailReason = ""
	}

	if v60.Valid {
		v.PushInQueueBeforeTrade = v60.Bool
	} else {
		v.PushInQueueBeforeTrade = false
	}

	if v61.Valid {
		v.LatestActionType = v61.String
	} else {
		v.LatestActionType = ""
	}

	if v62.Valid {
		v.ChannelCode = v62.String
	} else {
		v.ChannelCode = ""
	}

	if v63.Valid {
		v.DBInsertTime = v63.Int64
	} else {
		v.DBInsertTime = 0
	}

	if v64.Valid {
		v.MsgSeq = v64.Int64
	} else {
		v.MsgSeq = 0
	}

	if v65.Valid {
		v.WorkerAffinity = int(v65.Int64)
	} else {
		v.WorkerAffinity = 0
	}

	if v66.Valid {
		v.QuotaValidateTime = v66.Int64
	} else {
		v.QuotaValidateTime = 0
	}

	return v, nil
}

func scanTradeOrders(rows *sql.Rows) ([]*schema.TradeOrder, error) {
	var err error
	var vv []*schema.TradeOrder

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullBool
	var v11 sql.NullBool
	var v12 sql.NullBool
	var v13 sql.NullBool
	var v14 sql.NullBool
	var v15 sql.NullBool
	var v16 sql.NullBool
	var v17 sql.NullFloat64
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullInt64
	var v27 sql.NullInt64
	var v28 sql.NullFloat64
	var v29 sql.NullFloat64
	var v30 sql.NullString
	var v31 sql.NullFloat64
	var v32 sql.NullString
	var v33 sql.NullInt64
	var v34 sql.NullInt64
	var v35 sql.NullInt64
	var v36 sql.NullString
	var v37 sql.NullFloat64
	var v38 sql.NullInt64
	var v39 sql.NullInt64
	var v40 sql.NullString
	var v41 sql.NullString
	var v42 sql.NullString
	var v43 sql.NullString
	var v44 sql.NullInt64
	var v45 sql.NullString
	var v46 sql.NullInt64
	var v47 sql.NullInt64
	var v48 sql.NullString
	var v49 sql.NullInt64
	var v50 sql.NullString
	var v51 sql.NullString
	var v52 sql.NullInt64
	var v53 sql.NullString
	var v54 sql.NullString
	var v55 sql.NullString
	var v56 sql.NullString
	var v57 sql.NullInt64
	var v58 sql.NullInt64
	var v59 sql.NullString
	var v60 sql.NullBool
	var v61 sql.NullString
	var v62 sql.NullString
	var v63 sql.NullInt64
	var v64 sql.NullInt64
	var v65 sql.NullInt64
	var v66 sql.NullInt64

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
			&v64,
			&v65,
			&v66,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeOrder{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.SystemCode = v1.String
		} else {
			v.SystemCode = ""
		}

		if v2.Valid {
			v.BusinessCode = v2.String
		} else {
			v.BusinessCode = ""
		}

		if v3.Valid {
			v.ClGroupOrdID = v3.String
		} else {
			v.ClGroupOrdID = ""
		}

		if v4.Valid {
			v.ClOrdID = v4.String
		} else {
			v.ClOrdID = ""
		}

		if v5.Valid {
			v.Account = v5.String
		} else {
			v.Account = ""
		}

		if v6.Valid {
			v.HandlInst = v6.String
		} else {
			v.HandlInst = ""
		}

		if v7.Valid {
			v.AppOrdID = v7.String
		} else {
			v.AppOrdID = ""
		}

		if v8.Valid {
			v.OrdID = v8.String
		} else {
			v.OrdID = ""
		}

		if v9.Valid {
			v.ParentClOrdID = v9.String
		} else {
			v.ParentClOrdID = ""
		}

		if v10.Valid {
			v.IsDirectOrd = v10.Bool
		} else {
			v.IsDirectOrd = false
		}

		if v11.Valid {
			v.IsAlgOrd = v11.Bool
		} else {
			v.IsAlgOrd = false
		}

		if v12.Valid {
			v.IsSubAlgOrd = v12.Bool
		} else {
			v.IsSubAlgOrd = false
		}

		if v13.Valid {
			v.IsInstrOrd = v13.Bool
		} else {
			v.IsInstrOrd = false
		}

		if v14.Valid {
			v.IsSubInstrOrd = v14.Bool
		} else {
			v.IsSubInstrOrd = false
		}

		if v15.Valid {
			v.IsCrossDateOrd = v15.Bool
		} else {
			v.IsCrossDateOrd = false
		}

		if v16.Valid {
			v.IsSubCrossDateOrd = v16.Bool
		} else {
			v.IsSubCrossDateOrd = false
		}

		if v17.Valid {
			v.MinQty = v17.Float64
		} else {
			v.MinQty = 0
		}

		if v18.Valid {
			v.SecurityExchange = v18.String
		} else {
			v.SecurityExchange = ""
		}

		if v19.Valid {
			v.SecurityExchangeRegion = v19.String
		} else {
			v.SecurityExchangeRegion = ""
		}

		if v20.Valid {
			v.Symbol = v20.String
		} else {
			v.Symbol = ""
		}

		if v21.Valid {
			v.SymbolSfx = v21.String
		} else {
			v.SymbolSfx = ""
		}

		if v22.Valid {
			v.SecurityID = v22.String
		} else {
			v.SecurityID = ""
		}

		if v23.Valid {
			v.IDSource = v23.String
		} else {
			v.IDSource = ""
		}

		if v24.Valid {
			v.SecurityType = v24.String
		} else {
			v.SecurityType = ""
		}

		if v25.Valid {
			v.Side = v25.String
		} else {
			v.Side = ""
		}

		if v26.Valid {
			v.TransactTime = v26.Int64
		} else {
			v.TransactTime = 0
		}

		if v27.Valid {
			v.TradeDate = v27.Int64
		} else {
			v.TradeDate = 0
		}

		if v28.Valid {
			v.OrderQty = v28.Float64
		} else {
			v.OrderQty = 0
		}

		if v29.Valid {
			v.CashOrderQty = v29.Float64
		} else {
			v.CashOrderQty = 0
		}

		if v30.Valid {
			v.OrdType = v30.String
		} else {
			v.OrdType = ""
		}

		if v31.Valid {
			v.Price = v31.Float64
		} else {
			v.Price = 0
		}

		if v32.Valid {
			v.Currency = v32.String
		} else {
			v.Currency = ""
		}

		if v33.Valid {
			v.EffectiveTime = v33.Int64
		} else {
			v.EffectiveTime = 0
		}

		if v34.Valid {
			v.ExpireDate = v34.Int64
		} else {
			v.ExpireDate = 0
		}

		if v35.Valid {
			v.ExpireTime = v35.Int64
		} else {
			v.ExpireTime = 0
		}

		if v36.Valid {
			v.OpenClose = v36.String
		} else {
			v.OpenClose = ""
		}

		if v37.Valid {
			v.ContractMultiplier = v37.Float64
		} else {
			v.ContractMultiplier = 0
		}

		if v38.Valid {
			v.OrdChangedCount = int(v38.Int64)
		} else {
			v.OrdChangedCount = 0
		}

		if v39.Valid {
			v.OrdCancelCount = int(v39.Int64)
		} else {
			v.OrdCancelCount = 0
		}

		if v40.Valid {
			v.ExtendAttr = v40.String
		} else {
			v.ExtendAttr = ""
		}

		if v41.Valid {
			v.AlgParams = v41.String
		} else {
			v.AlgParams = ""
		}

		if v42.Valid {
			v.AlgName = v42.String
		} else {
			v.AlgName = ""
		}

		if v43.Valid {
			v.OrdCreator = v43.String
		} else {
			v.OrdCreator = ""
		}

		if v44.Valid {
			v.OrdCreateTime = v44.Int64
		} else {
			v.OrdCreateTime = 0
		}

		if v45.Valid {
			v.OrdDraftUpdateUser = v45.String
		} else {
			v.OrdDraftUpdateUser = ""
		}

		if v46.Valid {
			v.OrdDraftUpdateTime = v46.Int64
		} else {
			v.OrdDraftUpdateTime = 0
		}

		if v47.Valid {
			v.OrdDraftDelFlag = int(v47.Int64)
		} else {
			v.OrdDraftDelFlag = 0
		}

		if v48.Valid {
			v.OrdDraftDelUser = v48.String
		} else {
			v.OrdDraftDelUser = ""
		}

		if v49.Valid {
			v.OrdDraftDelTime = v49.Int64
		} else {
			v.OrdDraftDelTime = 0
		}

		if v50.Valid {
			v.OrdExecUserScope = v50.String
		} else {
			v.OrdExecUserScope = ""
		}

		if v51.Valid {
			v.OrdExecUser = v51.String
		} else {
			v.OrdExecUser = ""
		}

		if v52.Valid {
			v.OrdStatusUpdateTime = v52.Int64
		} else {
			v.OrdStatusUpdateTime = 0
		}

		if v53.Valid {
			v.OrdStatus = v53.String
		} else {
			v.OrdStatus = ""
		}

		if v54.Valid {
			v.ReviewFlag = v54.String
		} else {
			v.ReviewFlag = ""
		}

		if v55.Valid {
			v.ReviewerScope = v55.String
		} else {
			v.ReviewerScope = ""
		}

		if v56.Valid {
			v.Reviewer = v56.String
		} else {
			v.Reviewer = ""
		}

		if v57.Valid {
			v.ApproveStatus = int(v57.Int64)
		} else {
			v.ApproveStatus = 0
		}

		if v58.Valid {
			v.ReviewTime = v58.Int64
		} else {
			v.ReviewTime = 0
		}

		if v59.Valid {
			v.OrderSubmitFailReason = v59.String
		} else {
			v.OrderSubmitFailReason = ""
		}

		if v60.Valid {
			v.PushInQueueBeforeTrade = v60.Bool
		} else {
			v.PushInQueueBeforeTrade = false
		}

		if v61.Valid {
			v.LatestActionType = v61.String
		} else {
			v.LatestActionType = ""
		}

		if v62.Valid {
			v.ChannelCode = v62.String
		} else {
			v.ChannelCode = ""
		}

		if v63.Valid {
			v.DBInsertTime = v63.Int64
		} else {
			v.DBInsertTime = 0
		}

		if v64.Valid {
			v.MsgSeq = v64.Int64
		} else {
			v.MsgSeq = 0
		}

		if v65.Valid {
			v.WorkerAffinity = int(v65.Int64)
		} else {
			v.WorkerAffinity = 0
		}

		if v66.Valid {
			v.QuotaValidateTime = v66.Int64
		} else {
			v.QuotaValidateTime = 0
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeOrder(v *schema.TradeOrder) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 string
	var v10 bool
	var v11 bool
	var v12 bool
	var v13 bool
	var v14 bool
	var v15 bool
	var v16 bool
	var v17 float64
	var v18 string
	var v19 string
	var v20 string
	var v21 string
	var v22 string
	var v23 string
	var v24 string
	var v25 string
	var v26 int64
	var v27 int64
	var v28 float64
	var v29 float64
	var v30 string
	var v31 float64
	var v32 string
	var v33 int64
	var v34 int64
	var v35 int64
	var v36 string
	var v37 float64
	var v38 int
	var v39 int
	var v40 string
	var v41 string
	var v42 string
	var v43 string
	var v44 int64
	var v45 string
	var v46 int64
	var v47 int
	var v48 string
	var v49 int64
	var v50 string
	var v51 string
	var v52 int64
	var v53 string
	var v54 string
	var v55 string
	var v56 string
	var v57 int
	var v58 int64
	var v59 string
	var v60 bool
	var v61 string
	var v62 string
	var v63 int64
	var v64 int64
	var v65 int
	var v66 int64

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ClGroupOrdID
	v4 = v.ClOrdID
	v5 = v.Account
	v6 = v.HandlInst
	v7 = v.AppOrdID
	v8 = v.OrdID
	v9 = v.ParentClOrdID
	v10 = v.IsDirectOrd
	v11 = v.IsAlgOrd
	v12 = v.IsSubAlgOrd
	v13 = v.IsInstrOrd
	v14 = v.IsSubInstrOrd
	v15 = v.IsCrossDateOrd
	v16 = v.IsSubCrossDateOrd
	v17 = v.MinQty
	v18 = v.SecurityExchange
	v19 = v.SecurityExchangeRegion
	v20 = v.Symbol
	v21 = v.SymbolSfx
	v22 = v.SecurityID
	v23 = v.IDSource
	v24 = v.SecurityType
	v25 = v.Side
	v26 = v.TransactTime
	v27 = v.TradeDate
	v28 = v.OrderQty
	v29 = v.CashOrderQty
	v30 = v.OrdType
	v31 = v.Price
	v32 = v.Currency
	v33 = v.EffectiveTime
	v34 = v.ExpireDate
	v35 = v.ExpireTime
	v36 = v.OpenClose
	v37 = v.ContractMultiplier
	v38 = v.OrdChangedCount
	v39 = v.OrdCancelCount
	v40 = v.ExtendAttr
	v41 = v.AlgParams
	v42 = v.AlgName
	v43 = v.OrdCreator
	v44 = v.OrdCreateTime
	v45 = v.OrdDraftUpdateUser
	v46 = v.OrdDraftUpdateTime
	v47 = v.OrdDraftDelFlag
	v48 = v.OrdDraftDelUser
	v49 = v.OrdDraftDelTime
	v50 = v.OrdExecUserScope
	v51 = v.OrdExecUser
	v52 = v.OrdStatusUpdateTime
	v53 = v.OrdStatus
	v54 = v.ReviewFlag
	v55 = v.ReviewerScope
	v56 = v.Reviewer
	v57 = v.ApproveStatus
	v58 = v.ReviewTime
	v59 = v.OrderSubmitFailReason
	v60 = v.PushInQueueBeforeTrade
	v61 = v.LatestActionType
	v62 = v.ChannelCode
	v63 = v.DBInsertTime
	v64 = v.MsgSeq
	v65 = v.WorkerAffinity
	v66 = v.QuotaValidateTime

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
		v64,
		v65,
		v66,
	}
}

func genericSelectTradeOrder(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeOrder, error) {
	row := db.QueryRow(query, args...)
	return scanTradeOrder(row)
}

func genericSelectTradeOrders(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeOrder, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeOrders(rows)
}

func InsertTradeOrder(db db.SimpleDB, v *schema.TradeOrder) error {

	res, err := db.Exec(InsertTradeOrderStmt, sliceTradeOrder(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeOrderById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeOrderByIdStmt, args...)
	return err
}

func DeleteTradeOrderByClOrdId(db db.SimpleDB, clOrdID string) error {
	args := []interface{}{clOrdID}
	_, err := db.Exec(DeleteTradeOrderByClOrdIdStmt, args...)
	return err
}

func DeleteTradeOrderByAppOrdId(db db.SimpleDB, appOrdID string) error {
	args := []interface{}{appOrdID}
	_, err := db.Exec(DeleteTradeOrderByAppOrdIdStmt, args...)
	return err
}

func UpdateTradeOrderById(db db.SimpleDB, v *schema.TradeOrder) error {
	args := sliceTradeOrder(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeOrderByIdStmt, args...)
	return err
}

func UpdateTradeOrderByAppOrdId(db db.SimpleDB, v *schema.TradeOrder) error {
	args := sliceTradeOrder(v)
	args = append(args, v.AppOrdID)
	_, err := db.Exec(UpdateTradeOrderByAppOrdIdStmt, args...)
	return err
}

func GetTradeOrderById(db db.SimpleDB, iD int64) (*schema.TradeOrder, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeOrder(db, SelectTradeOrderByIdStmt, args...)
	return v, err
}

func GetTradeOrderByAppOrdId(db db.SimpleDB, appOrdID string) (*schema.TradeOrder, error) {
	args := []interface{}{appOrdID}
	v, err := genericSelectTradeOrder(db, SelectTradeOrderByAppOrdIdStmt, args...)
	return v, err
}

func FindAllTradeOrders(db db.SimpleDB) ([]*schema.TradeOrder, error) {
	args := []interface{}{}
	v, err := genericSelectTradeOrders(db, SelectTradeOrderStmt, args...)
	return v, err
}

func FindAllTradeOrdersInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeOrder, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeOrders(db, SelectTradeOrderRangeStmt, args...)
	return v, err
}

func FindTradeOrdersByClOrdId(db db.SimpleDB, clOrdID string) ([]*schema.TradeOrder, error) {
	args := []interface{}{clOrdID}
	v, err := genericSelectTradeOrders(db, SelectTradeOrderByClOrdIdStmt, args...)
	return v, err
}

func FindTradeOrdersByClOrdIdInRange(db db.SimpleDB, clOrdID string, limit int64, offset int64) ([]*schema.TradeOrder, error) {
	args := []interface{}{clOrdID, limit, offset}
	v, err := genericSelectTradeOrders(db, SelectTradeOrderRangeByClOrdIdStmt, args...)
	return v, err
}

func CountTradeOrder(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeOrderCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeOrderByClOrdId(db db.SimpleDB, clOrdID string) (int, error) {
	var count int
	args := []interface{}{clOrdID}
	row := db.QueryRow(SelectTradeOrderCountByClOrdIdStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTradeOrderByAppOrdId(db db.SimpleDB, appOrdID string) (int, error) {
	var count int
	args := []interface{}{appOrdID}
	row := db.QueryRow(SelectTradeOrderCountByAppOrdIdStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateGroupTradeOrderStmt = `
CREATE TABLE IF NOT EXISTS group_trade_orders (
 f_id                     BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code            VARCHAR(32)
,f_business_code          VARCHAR(32)
,f_cl_group_ord_id        VARCHAR(128)
,f_security_type          VARCHAR(2)
,f_sub_order_derive_type  VARCHAR(2)
,f_transact_time          BIGINT
,f_ord_status             VARCHAR(2)
,f_ord_fill_status        VARCHAR(2)
,f_ord_creator            VARCHAR(64)
,f_ord_create_time        BIGINT
,f_ord_draft_update_user  VARCHAR(64)
,f_ord_draft_update_time  VARCHAR(512)
,f_ord_draft_del_user     VARCHAR(64)
,f_ord_draft_del_time     BIGINT
,f_ord_exec_user_scope    VARCHAR(512)
,f_ord_exec_user          VARCHAR(64)
,f_ord_status_update_time BIGINT
,f_review_flag            INTEGER
,f_reviewer_scope         VARCHAR(512)
,f_reviewer               VARCHAR(512)
,f_approve_status         VARCHAR(512)
,f_review_time            BIGINT
);
`

const InsertGroupTradeOrderStmt = `
INSERT INTO group_trade_orders (
 f_system_code
,f_business_code
,f_cl_group_ord_id
,f_security_type
,f_sub_order_derive_type
,f_transact_time
,f_ord_status
,f_ord_fill_status
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectGroupTradeOrderStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_security_type
,f_sub_order_derive_type
,f_transact_time
,f_ord_status
,f_ord_fill_status
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
FROM group_trade_orders 
`

const SelectGroupTradeOrderRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_security_type
,f_sub_order_derive_type
,f_transact_time
,f_ord_status
,f_ord_fill_status
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
FROM group_trade_orders 
LIMIT ? OFFSET ?
`

const SelectGroupTradeOrderCountStmt = `
SELECT count(1)
FROM group_trade_orders 
`

const SelectGroupTradeOrderByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_security_type
,f_sub_order_derive_type
,f_transact_time
,f_ord_status
,f_ord_fill_status
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
FROM group_trade_orders 
WHERE f_id=?
`

const UpdateGroupTradeOrderByIdStmt = `
UPDATE group_trade_orders SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_cl_group_ord_id=?
,f_security_type=?
,f_sub_order_derive_type=?
,f_transact_time=?
,f_ord_status=?
,f_ord_fill_status=?
,f_ord_creator=?
,f_ord_create_time=?
,f_ord_draft_update_user=?
,f_ord_draft_update_time=?
,f_ord_draft_del_user=?
,f_ord_draft_del_time=?
,f_ord_exec_user_scope=?
,f_ord_exec_user=?
,f_ord_status_update_time=?
,f_review_flag=?
,f_reviewer_scope=?
,f_reviewer=?
,f_approve_status=?
,f_review_time=? 
WHERE f_id=?
`

const DeleteGroupTradeOrderByIdStmt = `
DELETE FROM group_trade_orders 
WHERE f_id=?
`

const CreatePkGtoClidStmt = `
CREATE UNIQUE INDEX pk_gto_clid ON group_trade_orders (f_cl_group_ord_id);
`

const SelectGroupTradeOrderByClGroupOrdIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_cl_group_ord_id
,f_security_type
,f_sub_order_derive_type
,f_transact_time
,f_ord_status
,f_ord_fill_status
,f_ord_creator
,f_ord_create_time
,f_ord_draft_update_user
,f_ord_draft_update_time
,f_ord_draft_del_user
,f_ord_draft_del_time
,f_ord_exec_user_scope
,f_ord_exec_user
,f_ord_status_update_time
,f_review_flag
,f_reviewer_scope
,f_reviewer
,f_approve_status
,f_review_time
FROM group_trade_orders 
WHERE f_cl_group_ord_id=?
`

const SelectGroupTradeOrderCountByClGroupOrdIdStmt = `
SELECT count(1)
FROM group_trade_orders 
WHERE f_cl_group_ord_id=?
`

const UpdateGroupTradeOrderByClGroupOrdIdStmt = `
UPDATE group_trade_orders SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_cl_group_ord_id=?
,f_security_type=?
,f_sub_order_derive_type=?
,f_transact_time=?
,f_ord_status=?
,f_ord_fill_status=?
,f_ord_creator=?
,f_ord_create_time=?
,f_ord_draft_update_user=?
,f_ord_draft_update_time=?
,f_ord_draft_del_user=?
,f_ord_draft_del_time=?
,f_ord_exec_user_scope=?
,f_ord_exec_user=?
,f_ord_status_update_time=?
,f_review_flag=?
,f_reviewer_scope=?
,f_reviewer=?
,f_approve_status=?
,f_review_time=? 
WHERE f_cl_group_ord_id=?
`

const DeleteGroupTradeOrderByClGroupOrdIdStmt = `
DELETE FROM group_trade_orders 
WHERE f_cl_group_ord_id=?
`

func scanGroupTradeOrder(row *sql.Row) (*schema.GroupTradeOrder, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullInt64
	var v18 sql.NullInt64
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullInt64

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.GroupTradeOrder{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.SystemCode = v1.String
	} else {
		v.SystemCode = ""
	}

	if v2.Valid {
		v.BusinessCode = v2.String
	} else {
		v.BusinessCode = ""
	}

	if v3.Valid {
		v.ClGroupOrdID = v3.String
	} else {
		v.ClGroupOrdID = ""
	}

	if v4.Valid {
		v.SecurityType = v4.String
	} else {
		v.SecurityType = ""
	}

	if v5.Valid {
		v.SubOrderDeriveType = v5.String
	} else {
		v.SubOrderDeriveType = ""
	}

	if v6.Valid {
		v.TransactTime = v6.Int64
	} else {
		v.TransactTime = 0
	}

	if v7.Valid {
		v.OrdStatus = v7.String
	} else {
		v.OrdStatus = ""
	}

	if v8.Valid {
		v.OrdFillStatus = v8.String
	} else {
		v.OrdFillStatus = ""
	}

	if v9.Valid {
		v.OrdCreator = v9.String
	} else {
		v.OrdCreator = ""
	}

	if v10.Valid {
		v.OrdCreateTime = v10.Int64
	} else {
		v.OrdCreateTime = 0
	}

	if v11.Valid {
		v.OrdDraftUpdateUser = v11.String
	} else {
		v.OrdDraftUpdateUser = ""
	}

	if v12.Valid {
		v.OrdDraftUpdateTime = v12.String
	} else {
		v.OrdDraftUpdateTime = ""
	}

	if v13.Valid {
		v.OrdDraftDelUser = v13.String
	} else {
		v.OrdDraftDelUser = ""
	}

	if v14.Valid {
		v.OrdDraftDelTime = v14.Int64
	} else {
		v.OrdDraftDelTime = 0
	}

	if v15.Valid {
		v.OrdExecUserScope = v15.String
	} else {
		v.OrdExecUserScope = ""
	}

	if v16.Valid {
		v.OrdExecUser = v16.String
	} else {
		v.OrdExecUser = ""
	}

	if v17.Valid {
		v.OrdStatusUpdateTime = v17.Int64
	} else {
		v.OrdStatusUpdateTime = 0
	}

	if v18.Valid {
		v.ReviewFlag = int(v18.Int64)
	} else {
		v.ReviewFlag = 0
	}

	if v19.Valid {
		v.ReviewerScope = v19.String
	} else {
		v.ReviewerScope = ""
	}

	if v20.Valid {
		v.Reviewer = v20.String
	} else {
		v.Reviewer = ""
	}

	if v21.Valid {
		v.ApproveStatus = v21.String
	} else {
		v.ApproveStatus = ""
	}

	if v22.Valid {
		v.ReviewTime = v22.Int64
	} else {
		v.ReviewTime = 0
	}

	return v, nil
}

func scanGroupTradeOrders(rows *sql.Rows) ([]*schema.GroupTradeOrder, error) {
	var err error
	var vv []*schema.GroupTradeOrder

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullInt64
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullInt64
	var v18 sql.NullInt64
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullInt64

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.GroupTradeOrder{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.SystemCode = v1.String
		} else {
			v.SystemCode = ""
		}

		if v2.Valid {
			v.BusinessCode = v2.String
		} else {
			v.BusinessCode = ""
		}

		if v3.Valid {
			v.ClGroupOrdID = v3.String
		} else {
			v.ClGroupOrdID = ""
		}

		if v4.Valid {
			v.SecurityType = v4.String
		} else {
			v.SecurityType = ""
		}

		if v5.Valid {
			v.SubOrderDeriveType = v5.String
		} else {
			v.SubOrderDeriveType = ""
		}

		if v6.Valid {
			v.TransactTime = v6.Int64
		} else {
			v.TransactTime = 0
		}

		if v7.Valid {
			v.OrdStatus = v7.String
		} else {
			v.OrdStatus = ""
		}

		if v8.Valid {
			v.OrdFillStatus = v8.String
		} else {
			v.OrdFillStatus = ""
		}

		if v9.Valid {
			v.OrdCreator = v9.String
		} else {
			v.OrdCreator = ""
		}

		if v10.Valid {
			v.OrdCreateTime = v10.Int64
		} else {
			v.OrdCreateTime = 0
		}

		if v11.Valid {
			v.OrdDraftUpdateUser = v11.String
		} else {
			v.OrdDraftUpdateUser = ""
		}

		if v12.Valid {
			v.OrdDraftUpdateTime = v12.String
		} else {
			v.OrdDraftUpdateTime = ""
		}

		if v13.Valid {
			v.OrdDraftDelUser = v13.String
		} else {
			v.OrdDraftDelUser = ""
		}

		if v14.Valid {
			v.OrdDraftDelTime = v14.Int64
		} else {
			v.OrdDraftDelTime = 0
		}

		if v15.Valid {
			v.OrdExecUserScope = v15.String
		} else {
			v.OrdExecUserScope = ""
		}

		if v16.Valid {
			v.OrdExecUser = v16.String
		} else {
			v.OrdExecUser = ""
		}

		if v17.Valid {
			v.OrdStatusUpdateTime = v17.Int64
		} else {
			v.OrdStatusUpdateTime = 0
		}

		if v18.Valid {
			v.ReviewFlag = int(v18.Int64)
		} else {
			v.ReviewFlag = 0
		}

		if v19.Valid {
			v.ReviewerScope = v19.String
		} else {
			v.ReviewerScope = ""
		}

		if v20.Valid {
			v.Reviewer = v20.String
		} else {
			v.Reviewer = ""
		}

		if v21.Valid {
			v.ApproveStatus = v21.String
		} else {
			v.ApproveStatus = ""
		}

		if v22.Valid {
			v.ReviewTime = v22.Int64
		} else {
			v.ReviewTime = 0
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceGroupTradeOrder(v *schema.GroupTradeOrder) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 int64
	var v7 string
	var v8 string
	var v9 string
	var v10 int64
	var v11 string
	var v12 string
	var v13 string
	var v14 int64
	var v15 string
	var v16 string
	var v17 int64
	var v18 int
	var v19 string
	var v20 string
	var v21 string
	var v22 int64

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ClGroupOrdID
	v4 = v.SecurityType
	v5 = v.SubOrderDeriveType
	v6 = v.TransactTime
	v7 = v.OrdStatus
	v8 = v.OrdFillStatus
	v9 = v.OrdCreator
	v10 = v.OrdCreateTime
	v11 = v.OrdDraftUpdateUser
	v12 = v.OrdDraftUpdateTime
	v13 = v.OrdDraftDelUser
	v14 = v.OrdDraftDelTime
	v15 = v.OrdExecUserScope
	v16 = v.OrdExecUser
	v17 = v.OrdStatusUpdateTime
	v18 = v.ReviewFlag
	v19 = v.ReviewerScope
	v20 = v.Reviewer
	v21 = v.ApproveStatus
	v22 = v.ReviewTime

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
	}
}

func genericSelectGroupTradeOrder(db db.SimpleDB, query string, args ...interface{}) (*schema.GroupTradeOrder, error) {
	row := db.QueryRow(query, args...)
	return scanGroupTradeOrder(row)
}

func genericSelectGroupTradeOrders(db db.SimpleDB, query string, args ...interface{}) ([]*schema.GroupTradeOrder, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGroupTradeOrders(rows)
}

func InsertGroupTradeOrder(db db.SimpleDB, v *schema.GroupTradeOrder) error {

	res, err := db.Exec(InsertGroupTradeOrderStmt, sliceGroupTradeOrder(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteGroupTradeOrderById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteGroupTradeOrderByIdStmt, args...)
	return err
}

func DeleteGroupTradeOrderByClGroupOrdId(db db.SimpleDB, clGroupOrdID string) error {
	args := []interface{}{clGroupOrdID}
	_, err := db.Exec(DeleteGroupTradeOrderByClGroupOrdIdStmt, args...)
	return err
}

func UpdateGroupTradeOrderById(db db.SimpleDB, v *schema.GroupTradeOrder) error {
	args := sliceGroupTradeOrder(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateGroupTradeOrderByIdStmt, args...)
	return err
}

func UpdateGroupTradeOrderByClGroupOrdId(db db.SimpleDB, v *schema.GroupTradeOrder) error {
	args := sliceGroupTradeOrder(v)
	args = append(args, v.ClGroupOrdID)
	_, err := db.Exec(UpdateGroupTradeOrderByClGroupOrdIdStmt, args...)
	return err
}

func GetGroupTradeOrderById(db db.SimpleDB, iD int64) (*schema.GroupTradeOrder, error) {
	args := []interface{}{iD}
	v, err := genericSelectGroupTradeOrder(db, SelectGroupTradeOrderByIdStmt, args...)
	return v, err
}

func GetGroupTradeOrderByClGroupOrdId(db db.SimpleDB, clGroupOrdID string) (*schema.GroupTradeOrder, error) {
	args := []interface{}{clGroupOrdID}
	v, err := genericSelectGroupTradeOrder(db, SelectGroupTradeOrderByClGroupOrdIdStmt, args...)
	return v, err
}

func FindAllGroupTradeOrders(db db.SimpleDB) ([]*schema.GroupTradeOrder, error) {
	args := []interface{}{}
	v, err := genericSelectGroupTradeOrders(db, SelectGroupTradeOrderStmt, args...)
	return v, err
}

func FindAllGroupTradeOrdersInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.GroupTradeOrder, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectGroupTradeOrders(db, SelectGroupTradeOrderRangeStmt, args...)
	return v, err
}

func CountGroupTradeOrder(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectGroupTradeOrderCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountGroupTradeOrderByClGroupOrdId(db db.SimpleDB, clGroupOrdID string) (int, error) {
	var count int
	args := []interface{}{clGroupOrdID}
	row := db.QueryRow(SelectGroupTradeOrderCountByClGroupOrdIdStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeActionLatestRespStmt = `
CREATE TABLE IF NOT EXISTS trade_action_latest_resps (
 f_id                       BIGINT PRIMARY KEY AUTO_INCREMENT
,f_action_user              VARCHAR(64)
,f_action_msg_time          BIGINT
,f_action_time              BIGINT
,f_action_type              VARCHAR(2)
,f_action_key               VARCHAR(188)
,f_ord_draft_before_update  MEDIUMTEXT
,f_app_ord_id               VARCHAR(188)
,f_order_id                 VARCHAR(188)
,f_root_cl_ord_id           VARCHAR(188)
,f_cl_ord_id                VARCHAR(188)
,f_orig_cl_ord_id           VARCHAR(188)
,f_exec_id                  VARCHAR(188)
,f_exec_type                VARCHAR(2)
,f_ord_status               VARCHAR(2)
,f_ord_rej_reason           VARCHAR(512)
,f_cxl_rej_response_to      VARCHAR(512)
,f_exec_restatement_reason  VARCHAR(512)
,f_account                  VARCHAR(64)
,f_security_exchange        VARCHAR(8)
,f_security_exchange_region VARCHAR(4)
,f_symbol                   VARCHAR(64)
,f_symbol_sfx               VARCHAR(8)
,f_security_id              VARCHAR(64)
,f_id_source                VARCHAR(2)
,f_security_type            VARCHAR(2)
,f_side                     VARCHAR(2)
,f_open_close               VARCHAR(2)
,f_order_qty                DOUBLE
,f_cash_order_qty           DOUBLE
,f_ord_type                 VARCHAR(2)
,f_price                    DOUBLE
,f_currency                 VARCHAR(4)
,f_effective_time           VARCHAR(64)
,f_expire_time              VARCHAR(64)
,f_last_shares              BIGINT
,f_last_px                  DOUBLE
,f_leaves_qty               BIGINT
,f_cum_qty                  BIGINT
,f_avg_px                   DOUBLE
,f_transact_time            BIGINT
,f_msg_time                 BIGINT
,f_msg_seq                  BIGINT
,f_stream_input_msg_seq     BIGINT
,f_channel_code             VARCHAR(32)
);
`

const InsertTradeActionLatestRespStmt = `
INSERT INTO trade_action_latest_resps (
 f_action_user
,f_action_msg_time
,f_action_time
,f_action_type
,f_action_key
,f_ord_draft_before_update
,f_app_ord_id
,f_order_id
,f_root_cl_ord_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_msg_time
,f_msg_seq
,f_stream_input_msg_seq
,f_channel_code
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeActionLatestRespStmt = `
SELECT 
 f_id
,f_action_user
,f_action_msg_time
,f_action_time
,f_action_type
,f_action_key
,f_ord_draft_before_update
,f_app_ord_id
,f_order_id
,f_root_cl_ord_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_msg_time
,f_msg_seq
,f_stream_input_msg_seq
,f_channel_code
FROM trade_action_latest_resps 
`

const SelectTradeActionLatestRespRangeStmt = `
SELECT 
 f_id
,f_action_user
,f_action_msg_time
,f_action_time
,f_action_type
,f_action_key
,f_ord_draft_before_update
,f_app_ord_id
,f_order_id
,f_root_cl_ord_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_msg_time
,f_msg_seq
,f_stream_input_msg_seq
,f_channel_code
FROM trade_action_latest_resps 
LIMIT ? OFFSET ?
`

const SelectTradeActionLatestRespCountStmt = `
SELECT count(1)
FROM trade_action_latest_resps 
`

const SelectTradeActionLatestRespByIdStmt = `
SELECT 
 f_id
,f_action_user
,f_action_msg_time
,f_action_time
,f_action_type
,f_action_key
,f_ord_draft_before_update
,f_app_ord_id
,f_order_id
,f_root_cl_ord_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_msg_time
,f_msg_seq
,f_stream_input_msg_seq
,f_channel_code
FROM trade_action_latest_resps 
WHERE f_id=?
`

const UpdateTradeActionLatestRespByIdStmt = `
UPDATE trade_action_latest_resps SET 
 f_id=?
,f_action_user=?
,f_action_msg_time=?
,f_action_time=?
,f_action_type=?
,f_action_key=?
,f_ord_draft_before_update=?
,f_app_ord_id=?
,f_order_id=?
,f_root_cl_ord_id=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=?
,f_exec_id=?
,f_exec_type=?
,f_ord_status=?
,f_ord_rej_reason=?
,f_cxl_rej_response_to=?
,f_exec_restatement_reason=?
,f_account=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_open_close=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_time=?
,f_last_shares=?
,f_last_px=?
,f_leaves_qty=?
,f_cum_qty=?
,f_avg_px=?
,f_transact_time=?
,f_msg_time=?
,f_msg_seq=?
,f_stream_input_msg_seq=?
,f_channel_code=? 
WHERE f_id=?
`

const DeleteTradeActionLatestRespByIdStmt = `
DELETE FROM trade_action_latest_resps 
WHERE f_id=?
`

const CreateUqTalrActkeyStmt = `
CREATE UNIQUE INDEX uq_talr_actkey ON trade_action_latest_resps (f_action_key);
`

const SelectTradeActionLatestRespByActionKeyStmt = `
SELECT 
 f_id
,f_action_user
,f_action_msg_time
,f_action_time
,f_action_type
,f_action_key
,f_ord_draft_before_update
,f_app_ord_id
,f_order_id
,f_root_cl_ord_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_security_exchange
,f_security_exchange_region
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_msg_time
,f_msg_seq
,f_stream_input_msg_seq
,f_channel_code
FROM trade_action_latest_resps 
WHERE f_action_key=?
`

const SelectTradeActionLatestRespCountByActionKeyStmt = `
SELECT count(1)
FROM trade_action_latest_resps 
WHERE f_action_key=?
`

const UpdateTradeActionLatestRespByActionKeyStmt = `
UPDATE trade_action_latest_resps SET 
 f_id=?
,f_action_user=?
,f_action_msg_time=?
,f_action_time=?
,f_action_type=?
,f_action_key=?
,f_ord_draft_before_update=?
,f_app_ord_id=?
,f_order_id=?
,f_root_cl_ord_id=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=?
,f_exec_id=?
,f_exec_type=?
,f_ord_status=?
,f_ord_rej_reason=?
,f_cxl_rej_response_to=?
,f_exec_restatement_reason=?
,f_account=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_open_close=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_time=?
,f_last_shares=?
,f_last_px=?
,f_leaves_qty=?
,f_cum_qty=?
,f_avg_px=?
,f_transact_time=?
,f_msg_time=?
,f_msg_seq=?
,f_stream_input_msg_seq=?
,f_channel_code=? 
WHERE f_action_key=?
`

const DeleteTradeActionLatestRespByActionKeyStmt = `
DELETE FROM trade_action_latest_resps 
WHERE f_action_key=?
`

func scanTradeActionLatestResp(row *sql.Row) (*schema.TradeActionLatestResp, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullFloat64
	var v29 sql.NullFloat64
	var v30 sql.NullString
	var v31 sql.NullFloat64
	var v32 sql.NullString
	var v33 sql.NullString
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullFloat64
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullFloat64
	var v40 sql.NullInt64
	var v41 sql.NullInt64
	var v42 sql.NullInt64
	var v43 sql.NullInt64
	var v44 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeActionLatestResp{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ActionUser = v1.String
	} else {
		v.ActionUser = ""
	}

	if v2.Valid {
		v.ActionMsgTime = v2.Int64
	} else {
		v.ActionMsgTime = 0
	}

	if v3.Valid {
		v.ActionTime = v3.Int64
	} else {
		v.ActionTime = 0
	}

	if v4.Valid {
		v.ActionType = v4.String
	} else {
		v.ActionType = ""
	}

	if v5.Valid {
		v.ActionKey = v5.String
	} else {
		v.ActionKey = ""
	}

	if v6.Valid {
		v.OrdDraftBeforeUpdate = v6.String
	} else {
		v.OrdDraftBeforeUpdate = ""
	}

	if v7.Valid {
		v.AppOrdID = v7.String
	} else {
		v.AppOrdID = ""
	}

	if v8.Valid {
		v.OrderID = v8.String
	} else {
		v.OrderID = ""
	}

	if v9.Valid {
		v.RootClOrdID = v9.String
	} else {
		v.RootClOrdID = ""
	}

	if v10.Valid {
		v.ClOrdID = v10.String
	} else {
		v.ClOrdID = ""
	}

	if v11.Valid {
		v.OrigClOrdID = v11.String
	} else {
		v.OrigClOrdID = ""
	}

	if v12.Valid {
		v.ExecID = v12.String
	} else {
		v.ExecID = ""
	}

	if v13.Valid {
		v.ExecType = v13.String
	} else {
		v.ExecType = ""
	}

	if v14.Valid {
		v.OrdStatus = v14.String
	} else {
		v.OrdStatus = ""
	}

	if v15.Valid {
		v.OrdRejReason = v15.String
	} else {
		v.OrdRejReason = ""
	}

	if v16.Valid {
		v.CxlRejResponseTo = v16.String
	} else {
		v.CxlRejResponseTo = ""
	}

	if v17.Valid {
		v.ExecRestatementReason = v17.String
	} else {
		v.ExecRestatementReason = ""
	}

	if v18.Valid {
		v.Account = v18.String
	} else {
		v.Account = ""
	}

	if v19.Valid {
		v.SecurityExchange = v19.String
	} else {
		v.SecurityExchange = ""
	}

	if v20.Valid {
		v.SecurityExchangeRegion = v20.String
	} else {
		v.SecurityExchangeRegion = ""
	}

	if v21.Valid {
		v.Symbol = v21.String
	} else {
		v.Symbol = ""
	}

	if v22.Valid {
		v.SymbolSfx = v22.String
	} else {
		v.SymbolSfx = ""
	}

	if v23.Valid {
		v.SecurityID = v23.String
	} else {
		v.SecurityID = ""
	}

	if v24.Valid {
		v.IDSource = v24.String
	} else {
		v.IDSource = ""
	}

	if v25.Valid {
		v.SecurityType = v25.String
	} else {
		v.SecurityType = ""
	}

	if v26.Valid {
		v.Side = v26.String
	} else {
		v.Side = ""
	}

	if v27.Valid {
		v.OpenClose = v27.String
	} else {
		v.OpenClose = ""
	}

	if v28.Valid {
		v.OrderQty = v28.Float64
	} else {
		v.OrderQty = 0
	}

	if v29.Valid {
		v.CashOrderQty = v29.Float64
	} else {
		v.CashOrderQty = 0
	}

	if v30.Valid {
		v.OrdType = v30.String
	} else {
		v.OrdType = ""
	}

	if v31.Valid {
		v.Price = v31.Float64
	} else {
		v.Price = 0
	}

	if v32.Valid {
		v.Currency = v32.String
	} else {
		v.Currency = ""
	}

	if v33.Valid {
		v.EffectiveTime = v33.String
	} else {
		v.EffectiveTime = ""
	}

	if v34.Valid {
		v.ExpireTime = v34.String
	} else {
		v.ExpireTime = ""
	}

	if v35.Valid {
		v.LastShares = v35.Int64
	} else {
		v.LastShares = 0
	}

	if v36.Valid {
		v.LastPx = v36.Float64
	} else {
		v.LastPx = 0
	}

	if v37.Valid {
		v.LeavesQty = v37.Int64
	} else {
		v.LeavesQty = 0
	}

	if v38.Valid {
		v.CumQty = v38.Int64
	} else {
		v.CumQty = 0
	}

	if v39.Valid {
		v.AvgPx = v39.Float64
	} else {
		v.AvgPx = 0
	}

	if v40.Valid {
		v.TransactTime = v40.Int64
	} else {
		v.TransactTime = 0
	}

	if v41.Valid {
		v.MsgTime = v41.Int64
	} else {
		v.MsgTime = 0
	}

	if v42.Valid {
		v.MsgSeq = v42.Int64
	} else {
		v.MsgSeq = 0
	}

	if v43.Valid {
		v.StreamInputMsgSeq = v43.Int64
	} else {
		v.StreamInputMsgSeq = 0
	}

	if v44.Valid {
		v.ChannelCode = v44.String
	} else {
		v.ChannelCode = ""
	}

	return v, nil
}

func scanTradeActionLatestResps(rows *sql.Rows) ([]*schema.TradeActionLatestResp, error) {
	var err error
	var vv []*schema.TradeActionLatestResp

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullInt64
	var v3 sql.NullInt64
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullString
	var v27 sql.NullString
	var v28 sql.NullFloat64
	var v29 sql.NullFloat64
	var v30 sql.NullString
	var v31 sql.NullFloat64
	var v32 sql.NullString
	var v33 sql.NullString
	var v34 sql.NullString
	var v35 sql.NullInt64
	var v36 sql.NullFloat64
	var v37 sql.NullInt64
	var v38 sql.NullInt64
	var v39 sql.NullFloat64
	var v40 sql.NullInt64
	var v41 sql.NullInt64
	var v42 sql.NullInt64
	var v43 sql.NullInt64
	var v44 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeActionLatestResp{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ActionUser = v1.String
		} else {
			v.ActionUser = ""
		}

		if v2.Valid {
			v.ActionMsgTime = v2.Int64
		} else {
			v.ActionMsgTime = 0
		}

		if v3.Valid {
			v.ActionTime = v3.Int64
		} else {
			v.ActionTime = 0
		}

		if v4.Valid {
			v.ActionType = v4.String
		} else {
			v.ActionType = ""
		}

		if v5.Valid {
			v.ActionKey = v5.String
		} else {
			v.ActionKey = ""
		}

		if v6.Valid {
			v.OrdDraftBeforeUpdate = v6.String
		} else {
			v.OrdDraftBeforeUpdate = ""
		}

		if v7.Valid {
			v.AppOrdID = v7.String
		} else {
			v.AppOrdID = ""
		}

		if v8.Valid {
			v.OrderID = v8.String
		} else {
			v.OrderID = ""
		}

		if v9.Valid {
			v.RootClOrdID = v9.String
		} else {
			v.RootClOrdID = ""
		}

		if v10.Valid {
			v.ClOrdID = v10.String
		} else {
			v.ClOrdID = ""
		}

		if v11.Valid {
			v.OrigClOrdID = v11.String
		} else {
			v.OrigClOrdID = ""
		}

		if v12.Valid {
			v.ExecID = v12.String
		} else {
			v.ExecID = ""
		}

		if v13.Valid {
			v.ExecType = v13.String
		} else {
			v.ExecType = ""
		}

		if v14.Valid {
			v.OrdStatus = v14.String
		} else {
			v.OrdStatus = ""
		}

		if v15.Valid {
			v.OrdRejReason = v15.String
		} else {
			v.OrdRejReason = ""
		}

		if v16.Valid {
			v.CxlRejResponseTo = v16.String
		} else {
			v.CxlRejResponseTo = ""
		}

		if v17.Valid {
			v.ExecRestatementReason = v17.String
		} else {
			v.ExecRestatementReason = ""
		}

		if v18.Valid {
			v.Account = v18.String
		} else {
			v.Account = ""
		}

		if v19.Valid {
			v.SecurityExchange = v19.String
		} else {
			v.SecurityExchange = ""
		}

		if v20.Valid {
			v.SecurityExchangeRegion = v20.String
		} else {
			v.SecurityExchangeRegion = ""
		}

		if v21.Valid {
			v.Symbol = v21.String
		} else {
			v.Symbol = ""
		}

		if v22.Valid {
			v.SymbolSfx = v22.String
		} else {
			v.SymbolSfx = ""
		}

		if v23.Valid {
			v.SecurityID = v23.String
		} else {
			v.SecurityID = ""
		}

		if v24.Valid {
			v.IDSource = v24.String
		} else {
			v.IDSource = ""
		}

		if v25.Valid {
			v.SecurityType = v25.String
		} else {
			v.SecurityType = ""
		}

		if v26.Valid {
			v.Side = v26.String
		} else {
			v.Side = ""
		}

		if v27.Valid {
			v.OpenClose = v27.String
		} else {
			v.OpenClose = ""
		}

		if v28.Valid {
			v.OrderQty = v28.Float64
		} else {
			v.OrderQty = 0
		}

		if v29.Valid {
			v.CashOrderQty = v29.Float64
		} else {
			v.CashOrderQty = 0
		}

		if v30.Valid {
			v.OrdType = v30.String
		} else {
			v.OrdType = ""
		}

		if v31.Valid {
			v.Price = v31.Float64
		} else {
			v.Price = 0
		}

		if v32.Valid {
			v.Currency = v32.String
		} else {
			v.Currency = ""
		}

		if v33.Valid {
			v.EffectiveTime = v33.String
		} else {
			v.EffectiveTime = ""
		}

		if v34.Valid {
			v.ExpireTime = v34.String
		} else {
			v.ExpireTime = ""
		}

		if v35.Valid {
			v.LastShares = v35.Int64
		} else {
			v.LastShares = 0
		}

		if v36.Valid {
			v.LastPx = v36.Float64
		} else {
			v.LastPx = 0
		}

		if v37.Valid {
			v.LeavesQty = v37.Int64
		} else {
			v.LeavesQty = 0
		}

		if v38.Valid {
			v.CumQty = v38.Int64
		} else {
			v.CumQty = 0
		}

		if v39.Valid {
			v.AvgPx = v39.Float64
		} else {
			v.AvgPx = 0
		}

		if v40.Valid {
			v.TransactTime = v40.Int64
		} else {
			v.TransactTime = 0
		}

		if v41.Valid {
			v.MsgTime = v41.Int64
		} else {
			v.MsgTime = 0
		}

		if v42.Valid {
			v.MsgSeq = v42.Int64
		} else {
			v.MsgSeq = 0
		}

		if v43.Valid {
			v.StreamInputMsgSeq = v43.Int64
		} else {
			v.StreamInputMsgSeq = 0
		}

		if v44.Valid {
			v.ChannelCode = v44.String
		} else {
			v.ChannelCode = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeActionLatestResp(v *schema.TradeActionLatestResp) []interface{} {
	var v0 int64
	var v1 string
	var v2 int64
	var v3 int64
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 string
	var v10 string
	var v11 string
	var v12 string
	var v13 string
	var v14 string
	var v15 string
	var v16 string
	var v17 string
	var v18 string
	var v19 string
	var v20 string
	var v21 string
	var v22 string
	var v23 string
	var v24 string
	var v25 string
	var v26 string
	var v27 string
	var v28 float64
	var v29 float64
	var v30 string
	var v31 float64
	var v32 string
	var v33 string
	var v34 string
	var v35 int64
	var v36 float64
	var v37 int64
	var v38 int64
	var v39 float64
	var v40 int64
	var v41 int64
	var v42 int64
	var v43 int64
	var v44 string

	v0 = v.ID
	v1 = v.ActionUser
	v2 = v.ActionMsgTime
	v3 = v.ActionTime
	v4 = v.ActionType
	v5 = v.ActionKey
	v6 = v.OrdDraftBeforeUpdate
	v7 = v.AppOrdID
	v8 = v.OrderID
	v9 = v.RootClOrdID
	v10 = v.ClOrdID
	v11 = v.OrigClOrdID
	v12 = v.ExecID
	v13 = v.ExecType
	v14 = v.OrdStatus
	v15 = v.OrdRejReason
	v16 = v.CxlRejResponseTo
	v17 = v.ExecRestatementReason
	v18 = v.Account
	v19 = v.SecurityExchange
	v20 = v.SecurityExchangeRegion
	v21 = v.Symbol
	v22 = v.SymbolSfx
	v23 = v.SecurityID
	v24 = v.IDSource
	v25 = v.SecurityType
	v26 = v.Side
	v27 = v.OpenClose
	v28 = v.OrderQty
	v29 = v.CashOrderQty
	v30 = v.OrdType
	v31 = v.Price
	v32 = v.Currency
	v33 = v.EffectiveTime
	v34 = v.ExpireTime
	v35 = v.LastShares
	v36 = v.LastPx
	v37 = v.LeavesQty
	v38 = v.CumQty
	v39 = v.AvgPx
	v40 = v.TransactTime
	v41 = v.MsgTime
	v42 = v.MsgSeq
	v43 = v.StreamInputMsgSeq
	v44 = v.ChannelCode

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
	}
}

func genericSelectTradeActionLatestResp(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeActionLatestResp, error) {
	row := db.QueryRow(query, args...)
	return scanTradeActionLatestResp(row)
}

func genericSelectTradeActionLatestResps(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeActionLatestResp, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeActionLatestResps(rows)
}

func InsertTradeActionLatestResp(db db.SimpleDB, v *schema.TradeActionLatestResp) error {

	res, err := db.Exec(InsertTradeActionLatestRespStmt, sliceTradeActionLatestResp(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeActionLatestRespById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeActionLatestRespByIdStmt, args...)
	return err
}

func DeleteTradeActionLatestRespByActionKey(db db.SimpleDB, actionKey string) error {
	args := []interface{}{actionKey}
	_, err := db.Exec(DeleteTradeActionLatestRespByActionKeyStmt, args...)
	return err
}

func UpdateTradeActionLatestRespById(db db.SimpleDB, v *schema.TradeActionLatestResp) error {
	args := sliceTradeActionLatestResp(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeActionLatestRespByIdStmt, args...)
	return err
}

func UpdateTradeActionLatestRespByActionKey(db db.SimpleDB, v *schema.TradeActionLatestResp) error {
	args := sliceTradeActionLatestResp(v)
	args = append(args, v.ActionKey)
	_, err := db.Exec(UpdateTradeActionLatestRespByActionKeyStmt, args...)
	return err
}

func GetTradeActionLatestRespById(db db.SimpleDB, iD int64) (*schema.TradeActionLatestResp, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeActionLatestResp(db, SelectTradeActionLatestRespByIdStmt, args...)
	return v, err
}

func GetTradeActionLatestRespByActionKey(db db.SimpleDB, actionKey string) (*schema.TradeActionLatestResp, error) {
	args := []interface{}{actionKey}
	v, err := genericSelectTradeActionLatestResp(db, SelectTradeActionLatestRespByActionKeyStmt, args...)
	return v, err
}

func FindAllTradeActionLatestResps(db db.SimpleDB) ([]*schema.TradeActionLatestResp, error) {
	args := []interface{}{}
	v, err := genericSelectTradeActionLatestResps(db, SelectTradeActionLatestRespStmt, args...)
	return v, err
}

func FindAllTradeActionLatestRespsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeActionLatestResp, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeActionLatestResps(db, SelectTradeActionLatestRespRangeStmt, args...)
	return v, err
}

func CountTradeActionLatestResp(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeActionLatestRespCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeActionLatestRespByActionKey(db db.SimpleDB, actionKey string) (int, error) {
	var count int
	args := []interface{}{actionKey}
	row := db.QueryRow(SelectTradeActionLatestRespCountByActionKeyStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeActionRespStmt = `
CREATE TABLE IF NOT EXISTS trade_action_resps (
 f_id                      BIGINT PRIMARY KEY AUTO_INCREMENT
,f_order_id                VARCHAR(188)
,f_cl_ord_id               VARCHAR(188)
,f_orig_cl_ord_id          VARCHAR(188)
,f_exec_id                 VARCHAR(188)
,f_exec_ref_id             VARCHAR(188)
,f_exec_trans_type         VARCHAR(2)
,f_exec_type               VARCHAR(2)
,f_ord_status              VARCHAR(2)
,f_ord_rej_reason          VARCHAR(512)
,f_cxl_rej_response_to     VARCHAR(512)
,f_exec_restatement_reason VARCHAR(512)
,f_account                 VARCHAR(64)
,f_symbol                  VARCHAR(64)
,f_symbol_sfx              VARCHAR(8)
,f_security_id             VARCHAR(64)
,f_id_source               VARCHAR(2)
,f_security_type           VARCHAR(2)
,f_side                    VARCHAR(2)
,f_open_close              VARCHAR(2)
,f_order_qty               DOUBLE
,f_cash_order_qty          DOUBLE
,f_ord_type                VARCHAR(2)
,f_price                   DOUBLE
,f_currency                VARCHAR(4)
,f_effective_time          VARCHAR(64)
,f_expire_time             VARCHAR(64)
,f_last_shares             BIGINT
,f_last_px                 DOUBLE
,f_leaves_qty              BIGINT
,f_cum_qty                 BIGINT
,f_avg_px                  DOUBLE
,f_transact_time           BIGINT
,f_exchange_trade_date     VARCHAR(16)
,f_msg_time                BIGINT
,f_db_insert_time          BIGINT
,f_msg_seq                 BIGINT
,f_channel_code            VARCHAR(32)
,f_extend_attr             MEDIUMTEXT
);
`

const InsertTradeActionRespStmt = `
INSERT INTO trade_action_resps (
 f_order_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_ref_id
,f_exec_trans_type
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_exchange_trade_date
,f_msg_time
,f_db_insert_time
,f_msg_seq
,f_channel_code
,f_extend_attr
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeActionRespStmt = `
SELECT 
 f_id
,f_order_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_ref_id
,f_exec_trans_type
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_exchange_trade_date
,f_msg_time
,f_db_insert_time
,f_msg_seq
,f_channel_code
,f_extend_attr
FROM trade_action_resps 
`

const SelectTradeActionRespRangeStmt = `
SELECT 
 f_id
,f_order_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_ref_id
,f_exec_trans_type
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_exchange_trade_date
,f_msg_time
,f_db_insert_time
,f_msg_seq
,f_channel_code
,f_extend_attr
FROM trade_action_resps 
LIMIT ? OFFSET ?
`

const SelectTradeActionRespCountStmt = `
SELECT count(1)
FROM trade_action_resps 
`

const SelectTradeActionRespByIdStmt = `
SELECT 
 f_id
,f_order_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_ref_id
,f_exec_trans_type
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_exchange_trade_date
,f_msg_time
,f_db_insert_time
,f_msg_seq
,f_channel_code
,f_extend_attr
FROM trade_action_resps 
WHERE f_id=?
`

const UpdateTradeActionRespByIdStmt = `
UPDATE trade_action_resps SET 
 f_id=?
,f_order_id=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=?
,f_exec_id=?
,f_exec_ref_id=?
,f_exec_trans_type=?
,f_exec_type=?
,f_ord_status=?
,f_ord_rej_reason=?
,f_cxl_rej_response_to=?
,f_exec_restatement_reason=?
,f_account=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_open_close=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_time=?
,f_last_shares=?
,f_last_px=?
,f_leaves_qty=?
,f_cum_qty=?
,f_avg_px=?
,f_transact_time=?
,f_exchange_trade_date=?
,f_msg_time=?
,f_db_insert_time=?
,f_msg_seq=?
,f_channel_code=?
,f_extend_attr=? 
WHERE f_id=?
`

const DeleteTradeActionRespByIdStmt = `
DELETE FROM trade_action_resps 
WHERE f_id=?
`

const CreatePkTarStmt = `
CREATE UNIQUE INDEX pk_tar ON trade_action_resps (f_cl_ord_id,f_exec_id,f_channel_code);
`

const SelectTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt = `
SELECT 
 f_id
,f_order_id
,f_cl_ord_id
,f_orig_cl_ord_id
,f_exec_id
,f_exec_ref_id
,f_exec_trans_type
,f_exec_type
,f_ord_status
,f_ord_rej_reason
,f_cxl_rej_response_to
,f_exec_restatement_reason
,f_account
,f_symbol
,f_symbol_sfx
,f_security_id
,f_id_source
,f_security_type
,f_side
,f_open_close
,f_order_qty
,f_cash_order_qty
,f_ord_type
,f_price
,f_currency
,f_effective_time
,f_expire_time
,f_last_shares
,f_last_px
,f_leaves_qty
,f_cum_qty
,f_avg_px
,f_transact_time
,f_exchange_trade_date
,f_msg_time
,f_db_insert_time
,f_msg_seq
,f_channel_code
,f_extend_attr
FROM trade_action_resps 
WHERE f_cl_ord_id=?
AND f_exec_id=?
AND f_channel_code=?
`

const SelectTradeActionRespCountByClOrdIdAndExecIdAndChannelCodeStmt = `
SELECT count(1)
FROM trade_action_resps 
WHERE f_cl_ord_id=?
AND f_exec_id=?
AND f_channel_code=?
`

const UpdateTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt = `
UPDATE trade_action_resps SET 
 f_id=?
,f_order_id=?
,f_cl_ord_id=?
,f_orig_cl_ord_id=?
,f_exec_id=?
,f_exec_ref_id=?
,f_exec_trans_type=?
,f_exec_type=?
,f_ord_status=?
,f_ord_rej_reason=?
,f_cxl_rej_response_to=?
,f_exec_restatement_reason=?
,f_account=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_id=?
,f_id_source=?
,f_security_type=?
,f_side=?
,f_open_close=?
,f_order_qty=?
,f_cash_order_qty=?
,f_ord_type=?
,f_price=?
,f_currency=?
,f_effective_time=?
,f_expire_time=?
,f_last_shares=?
,f_last_px=?
,f_leaves_qty=?
,f_cum_qty=?
,f_avg_px=?
,f_transact_time=?
,f_exchange_trade_date=?
,f_msg_time=?
,f_db_insert_time=?
,f_msg_seq=?
,f_channel_code=?
,f_extend_attr=? 
WHERE f_cl_ord_id=?
AND f_exec_id=?
AND f_channel_code=?
`

const DeleteTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt = `
DELETE FROM trade_action_resps 
WHERE f_cl_ord_id=?
AND f_exec_id=?
AND f_channel_code=?
`

func scanTradeActionResp(row *sql.Row) (*schema.TradeActionResp, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullFloat64
	var v21 sql.NullFloat64
	var v22 sql.NullString
	var v23 sql.NullFloat64
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullString
	var v27 sql.NullInt64
	var v28 sql.NullFloat64
	var v29 sql.NullInt64
	var v30 sql.NullInt64
	var v31 sql.NullFloat64
	var v32 sql.NullInt64
	var v33 sql.NullString
	var v34 sql.NullInt64
	var v35 sql.NullInt64
	var v36 sql.NullInt64
	var v37 sql.NullString
	var v38 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeActionResp{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.OrderID = v1.String
	} else {
		v.OrderID = ""
	}

	if v2.Valid {
		v.ClOrdID = v2.String
	} else {
		v.ClOrdID = ""
	}

	if v3.Valid {
		v.OrigClOrdID = v3.String
	} else {
		v.OrigClOrdID = ""
	}

	if v4.Valid {
		v.ExecID = v4.String
	} else {
		v.ExecID = ""
	}

	if v5.Valid {
		v.ExecRefID = v5.String
	} else {
		v.ExecRefID = ""
	}

	if v6.Valid {
		v.ExecTransType = v6.String
	} else {
		v.ExecTransType = ""
	}

	if v7.Valid {
		v.ExecType = v7.String
	} else {
		v.ExecType = ""
	}

	if v8.Valid {
		v.OrdStatus = v8.String
	} else {
		v.OrdStatus = ""
	}

	if v9.Valid {
		v.OrdRejReason = v9.String
	} else {
		v.OrdRejReason = ""
	}

	if v10.Valid {
		v.CxlRejResponseTo = v10.String
	} else {
		v.CxlRejResponseTo = ""
	}

	if v11.Valid {
		v.ExecRestatementReason = v11.String
	} else {
		v.ExecRestatementReason = ""
	}

	if v12.Valid {
		v.Account = v12.String
	} else {
		v.Account = ""
	}

	if v13.Valid {
		v.Symbol = v13.String
	} else {
		v.Symbol = ""
	}

	if v14.Valid {
		v.SymbolSfx = v14.String
	} else {
		v.SymbolSfx = ""
	}

	if v15.Valid {
		v.SecurityID = v15.String
	} else {
		v.SecurityID = ""
	}

	if v16.Valid {
		v.IDSource = v16.String
	} else {
		v.IDSource = ""
	}

	if v17.Valid {
		v.SecurityType = v17.String
	} else {
		v.SecurityType = ""
	}

	if v18.Valid {
		v.Side = v18.String
	} else {
		v.Side = ""
	}

	if v19.Valid {
		v.OpenClose = v19.String
	} else {
		v.OpenClose = ""
	}

	if v20.Valid {
		v.OrderQty = v20.Float64
	} else {
		v.OrderQty = 0
	}

	if v21.Valid {
		v.CashOrderQty = v21.Float64
	} else {
		v.CashOrderQty = 0
	}

	if v22.Valid {
		v.OrdType = v22.String
	} else {
		v.OrdType = ""
	}

	if v23.Valid {
		v.Price = v23.Float64
	} else {
		v.Price = 0
	}

	if v24.Valid {
		v.Currency = v24.String
	} else {
		v.Currency = ""
	}

	if v25.Valid {
		v.EffectiveTime = v25.String
	} else {
		v.EffectiveTime = ""
	}

	if v26.Valid {
		v.ExpireTime = v26.String
	} else {
		v.ExpireTime = ""
	}

	if v27.Valid {
		v.LastShares = v27.Int64
	} else {
		v.LastShares = 0
	}

	if v28.Valid {
		v.LastPx = v28.Float64
	} else {
		v.LastPx = 0
	}

	if v29.Valid {
		v.LeavesQty = v29.Int64
	} else {
		v.LeavesQty = 0
	}

	if v30.Valid {
		v.CumQty = v30.Int64
	} else {
		v.CumQty = 0
	}

	if v31.Valid {
		v.AvgPx = v31.Float64
	} else {
		v.AvgPx = 0
	}

	if v32.Valid {
		v.TransactTime = v32.Int64
	} else {
		v.TransactTime = 0
	}

	if v33.Valid {
		v.ExchangeTradeDate = v33.String
	} else {
		v.ExchangeTradeDate = ""
	}

	if v34.Valid {
		v.MsgTime = v34.Int64
	} else {
		v.MsgTime = 0
	}

	if v35.Valid {
		v.DBInsertTime = v35.Int64
	} else {
		v.DBInsertTime = 0
	}

	if v36.Valid {
		v.MsgSeq = v36.Int64
	} else {
		v.MsgSeq = 0
	}

	if v37.Valid {
		v.ChannelCode = v37.String
	} else {
		v.ChannelCode = ""
	}

	if v38.Valid {
		v.ExtendAttr = v38.String
	} else {
		v.ExtendAttr = ""
	}

	return v, nil
}

func scanTradeActionResps(rows *sql.Rows) ([]*schema.TradeActionResp, error) {
	var err error
	var vv []*schema.TradeActionResp

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullString
	var v10 sql.NullString
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullString
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullFloat64
	var v21 sql.NullFloat64
	var v22 sql.NullString
	var v23 sql.NullFloat64
	var v24 sql.NullString
	var v25 sql.NullString
	var v26 sql.NullString
	var v27 sql.NullInt64
	var v28 sql.NullFloat64
	var v29 sql.NullInt64
	var v30 sql.NullInt64
	var v31 sql.NullFloat64
	var v32 sql.NullInt64
	var v33 sql.NullString
	var v34 sql.NullInt64
	var v35 sql.NullInt64
	var v36 sql.NullInt64
	var v37 sql.NullString
	var v38 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeActionResp{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.OrderID = v1.String
		} else {
			v.OrderID = ""
		}

		if v2.Valid {
			v.ClOrdID = v2.String
		} else {
			v.ClOrdID = ""
		}

		if v3.Valid {
			v.OrigClOrdID = v3.String
		} else {
			v.OrigClOrdID = ""
		}

		if v4.Valid {
			v.ExecID = v4.String
		} else {
			v.ExecID = ""
		}

		if v5.Valid {
			v.ExecRefID = v5.String
		} else {
			v.ExecRefID = ""
		}

		if v6.Valid {
			v.ExecTransType = v6.String
		} else {
			v.ExecTransType = ""
		}

		if v7.Valid {
			v.ExecType = v7.String
		} else {
			v.ExecType = ""
		}

		if v8.Valid {
			v.OrdStatus = v8.String
		} else {
			v.OrdStatus = ""
		}

		if v9.Valid {
			v.OrdRejReason = v9.String
		} else {
			v.OrdRejReason = ""
		}

		if v10.Valid {
			v.CxlRejResponseTo = v10.String
		} else {
			v.CxlRejResponseTo = ""
		}

		if v11.Valid {
			v.ExecRestatementReason = v11.String
		} else {
			v.ExecRestatementReason = ""
		}

		if v12.Valid {
			v.Account = v12.String
		} else {
			v.Account = ""
		}

		if v13.Valid {
			v.Symbol = v13.String
		} else {
			v.Symbol = ""
		}

		if v14.Valid {
			v.SymbolSfx = v14.String
		} else {
			v.SymbolSfx = ""
		}

		if v15.Valid {
			v.SecurityID = v15.String
		} else {
			v.SecurityID = ""
		}

		if v16.Valid {
			v.IDSource = v16.String
		} else {
			v.IDSource = ""
		}

		if v17.Valid {
			v.SecurityType = v17.String
		} else {
			v.SecurityType = ""
		}

		if v18.Valid {
			v.Side = v18.String
		} else {
			v.Side = ""
		}

		if v19.Valid {
			v.OpenClose = v19.String
		} else {
			v.OpenClose = ""
		}

		if v20.Valid {
			v.OrderQty = v20.Float64
		} else {
			v.OrderQty = 0
		}

		if v21.Valid {
			v.CashOrderQty = v21.Float64
		} else {
			v.CashOrderQty = 0
		}

		if v22.Valid {
			v.OrdType = v22.String
		} else {
			v.OrdType = ""
		}

		if v23.Valid {
			v.Price = v23.Float64
		} else {
			v.Price = 0
		}

		if v24.Valid {
			v.Currency = v24.String
		} else {
			v.Currency = ""
		}

		if v25.Valid {
			v.EffectiveTime = v25.String
		} else {
			v.EffectiveTime = ""
		}

		if v26.Valid {
			v.ExpireTime = v26.String
		} else {
			v.ExpireTime = ""
		}

		if v27.Valid {
			v.LastShares = v27.Int64
		} else {
			v.LastShares = 0
		}

		if v28.Valid {
			v.LastPx = v28.Float64
		} else {
			v.LastPx = 0
		}

		if v29.Valid {
			v.LeavesQty = v29.Int64
		} else {
			v.LeavesQty = 0
		}

		if v30.Valid {
			v.CumQty = v30.Int64
		} else {
			v.CumQty = 0
		}

		if v31.Valid {
			v.AvgPx = v31.Float64
		} else {
			v.AvgPx = 0
		}

		if v32.Valid {
			v.TransactTime = v32.Int64
		} else {
			v.TransactTime = 0
		}

		if v33.Valid {
			v.ExchangeTradeDate = v33.String
		} else {
			v.ExchangeTradeDate = ""
		}

		if v34.Valid {
			v.MsgTime = v34.Int64
		} else {
			v.MsgTime = 0
		}

		if v35.Valid {
			v.DBInsertTime = v35.Int64
		} else {
			v.DBInsertTime = 0
		}

		if v36.Valid {
			v.MsgSeq = v36.Int64
		} else {
			v.MsgSeq = 0
		}

		if v37.Valid {
			v.ChannelCode = v37.String
		} else {
			v.ChannelCode = ""
		}

		if v38.Valid {
			v.ExtendAttr = v38.String
		} else {
			v.ExtendAttr = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeActionResp(v *schema.TradeActionResp) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 string
	var v10 string
	var v11 string
	var v12 string
	var v13 string
	var v14 string
	var v15 string
	var v16 string
	var v17 string
	var v18 string
	var v19 string
	var v20 float64
	var v21 float64
	var v22 string
	var v23 float64
	var v24 string
	var v25 string
	var v26 string
	var v27 int64
	var v28 float64
	var v29 int64
	var v30 int64
	var v31 float64
	var v32 int64
	var v33 string
	var v34 int64
	var v35 int64
	var v36 int64
	var v37 string
	var v38 string

	v0 = v.ID
	v1 = v.OrderID
	v2 = v.ClOrdID
	v3 = v.OrigClOrdID
	v4 = v.ExecID
	v5 = v.ExecRefID
	v6 = v.ExecTransType
	v7 = v.ExecType
	v8 = v.OrdStatus
	v9 = v.OrdRejReason
	v10 = v.CxlRejResponseTo
	v11 = v.ExecRestatementReason
	v12 = v.Account
	v13 = v.Symbol
	v14 = v.SymbolSfx
	v15 = v.SecurityID
	v16 = v.IDSource
	v17 = v.SecurityType
	v18 = v.Side
	v19 = v.OpenClose
	v20 = v.OrderQty
	v21 = v.CashOrderQty
	v22 = v.OrdType
	v23 = v.Price
	v24 = v.Currency
	v25 = v.EffectiveTime
	v26 = v.ExpireTime
	v27 = v.LastShares
	v28 = v.LastPx
	v29 = v.LeavesQty
	v30 = v.CumQty
	v31 = v.AvgPx
	v32 = v.TransactTime
	v33 = v.ExchangeTradeDate
	v34 = v.MsgTime
	v35 = v.DBInsertTime
	v36 = v.MsgSeq
	v37 = v.ChannelCode
	v38 = v.ExtendAttr

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
	}
}

func genericSelectTradeActionResp(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeActionResp, error) {
	row := db.QueryRow(query, args...)
	return scanTradeActionResp(row)
}

func genericSelectTradeActionResps(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeActionResp, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeActionResps(rows)
}

func InsertTradeActionResp(db db.SimpleDB, v *schema.TradeActionResp) error {

	res, err := db.Exec(InsertTradeActionRespStmt, sliceTradeActionResp(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeActionRespById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeActionRespByIdStmt, args...)
	return err
}

func DeleteTradeActionRespByClOrdIdAndExecIdAndChannelCode(db db.SimpleDB, clOrdID string, execID string, channelCode string) error {
	args := []interface{}{clOrdID, execID, channelCode}
	_, err := db.Exec(DeleteTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt, args...)
	return err
}

func UpdateTradeActionRespById(db db.SimpleDB, v *schema.TradeActionResp) error {
	args := sliceTradeActionResp(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeActionRespByIdStmt, args...)
	return err
}

func UpdateTradeActionRespByClOrdIdAndExecIdAndChannelCode(db db.SimpleDB, v *schema.TradeActionResp) error {
	args := sliceTradeActionResp(v)
	args = append(args, v.ClOrdID, v.ExecID, v.ChannelCode)
	_, err := db.Exec(UpdateTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt, args...)
	return err
}

func GetTradeActionRespById(db db.SimpleDB, iD int64) (*schema.TradeActionResp, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeActionResp(db, SelectTradeActionRespByIdStmt, args...)
	return v, err
}

func GetTradeActionRespByClOrdIdAndExecIdAndChannelCode(db db.SimpleDB, clOrdID string, execID string, channelCode string) (*schema.TradeActionResp, error) {
	args := []interface{}{clOrdID, execID, channelCode}
	v, err := genericSelectTradeActionResp(db, SelectTradeActionRespByClOrdIdAndExecIdAndChannelCodeStmt, args...)
	return v, err
}

func FindAllTradeActionResps(db db.SimpleDB) ([]*schema.TradeActionResp, error) {
	args := []interface{}{}
	v, err := genericSelectTradeActionResps(db, SelectTradeActionRespStmt, args...)
	return v, err
}

func FindAllTradeActionRespsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeActionResp, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeActionResps(db, SelectTradeActionRespRangeStmt, args...)
	return v, err
}

func CountTradeActionResp(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeActionRespCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeActionRespByClOrdIdAndExecIdAndChannelCode(db db.SimpleDB, clOrdID string, execID string, channelCode string) (int, error) {
	var count int
	args := []interface{}{clOrdID, execID, channelCode}
	row := db.QueryRow(SelectTradeActionRespCountByClOrdIdAndExecIdAndChannelCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateUtilFixMessageStmt = `
CREATE TABLE IF NOT EXISTS util_fix_messages (
 f_id           BIGINT PRIMARY KEY AUTO_INCREMENT
,f_msg_side     INTEGER
,f_msg_type     VARCHAR(2)
,f_msg_time     BIGINT
,f_data         MEDIUMBLOB
,f_channel_code VARCHAR(32)
);
`

const InsertUtilFixMessageStmt = `
INSERT INTO util_fix_messages (
 f_msg_side
,f_msg_type
,f_msg_time
,f_data
,f_channel_code
) VALUES (?,?,?,?,?)
`

const SelectUtilFixMessageStmt = `
SELECT 
 f_id
,f_msg_side
,f_msg_type
,f_msg_time
,f_data
,f_channel_code
FROM util_fix_messages 
`

const SelectUtilFixMessageRangeStmt = `
SELECT 
 f_id
,f_msg_side
,f_msg_type
,f_msg_time
,f_data
,f_channel_code
FROM util_fix_messages 
LIMIT ? OFFSET ?
`

const SelectUtilFixMessageCountStmt = `
SELECT count(1)
FROM util_fix_messages 
`

const SelectUtilFixMessageByIdStmt = `
SELECT 
 f_id
,f_msg_side
,f_msg_type
,f_msg_time
,f_data
,f_channel_code
FROM util_fix_messages 
WHERE f_id=?
`

const UpdateUtilFixMessageByIdStmt = `
UPDATE util_fix_messages SET 
 f_id=?
,f_msg_side=?
,f_msg_type=?
,f_msg_time=?
,f_data=?
,f_channel_code=? 
WHERE f_id=?
`

const DeleteUtilFixMessageByIdStmt = `
DELETE FROM util_fix_messages 
WHERE f_id=?
`

func scanUtilFixMessage(row *sql.Row) (*schema.UtilFixMessage, error) {
	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullString
	var v3 sql.NullInt64
	var v4 db.NullBytes
	var v5 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.UtilFixMessage{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.MsgSide = int(v1.Int64)
	} else {
		v.MsgSide = 0
	}

	if v2.Valid {
		v.MsgType = v2.String
	} else {
		v.MsgType = ""
	}

	if v3.Valid {
		v.MsgTime = v3.Int64
	} else {
		v.MsgTime = 0
	}

	if v4.Valid {
		v.Data = v4.Bytes
	} else {
		v.Data = nil
	}

	if v5.Valid {
		v.ChannelCode = v5.String
	} else {
		v.ChannelCode = ""
	}

	return v, nil
}

func scanUtilFixMessages(rows *sql.Rows) ([]*schema.UtilFixMessage, error) {
	var err error
	var vv []*schema.UtilFixMessage

	var v0 sql.NullInt64
	var v1 sql.NullInt64
	var v2 sql.NullString
	var v3 sql.NullInt64
	var v4 db.NullBytes
	var v5 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.UtilFixMessage{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.MsgSide = int(v1.Int64)
		} else {
			v.MsgSide = 0
		}

		if v2.Valid {
			v.MsgType = v2.String
		} else {
			v.MsgType = ""
		}

		if v3.Valid {
			v.MsgTime = v3.Int64
		} else {
			v.MsgTime = 0
		}

		if v4.Valid {
			v.Data = v4.Bytes
		} else {
			v.Data = nil
		}

		if v5.Valid {
			v.ChannelCode = v5.String
		} else {
			v.ChannelCode = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceUtilFixMessage(v *schema.UtilFixMessage) []interface{} {
	var v0 int64
	var v1 int
	var v2 string
	var v3 int64
	var v4 []byte
	var v5 string

	v0 = v.ID
	v1 = v.MsgSide
	v2 = v.MsgType
	v3 = v.MsgTime
	v4 = v.Data
	v5 = v.ChannelCode

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
	}
}

func genericSelectUtilFixMessage(db db.SimpleDB, query string, args ...interface{}) (*schema.UtilFixMessage, error) {
	row := db.QueryRow(query, args...)
	return scanUtilFixMessage(row)
}

func genericSelectUtilFixMessages(db db.SimpleDB, query string, args ...interface{}) ([]*schema.UtilFixMessage, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanUtilFixMessages(rows)
}

func InsertUtilFixMessage(db db.SimpleDB, v *schema.UtilFixMessage) error {

	res, err := db.Exec(InsertUtilFixMessageStmt, sliceUtilFixMessage(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteUtilFixMessageById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteUtilFixMessageByIdStmt, args...)
	return err
}

func UpdateUtilFixMessageById(db db.SimpleDB, v *schema.UtilFixMessage) error {
	args := sliceUtilFixMessage(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateUtilFixMessageByIdStmt, args...)
	return err
}

func GetUtilFixMessageById(db db.SimpleDB, iD int64) (*schema.UtilFixMessage, error) {
	args := []interface{}{iD}
	v, err := genericSelectUtilFixMessage(db, SelectUtilFixMessageByIdStmt, args...)
	return v, err
}

func FindAllUtilFixMessages(db db.SimpleDB) ([]*schema.UtilFixMessage, error) {
	args := []interface{}{}
	v, err := genericSelectUtilFixMessages(db, SelectUtilFixMessageStmt, args...)
	return v, err
}

func FindAllUtilFixMessagesInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.UtilFixMessage, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectUtilFixMessages(db, SelectUtilFixMessageRangeStmt, args...)
	return v, err
}

func CountUtilFixMessage(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectUtilFixMessageCountStmt)
	err := row.Scan(&count)
	return count, err
}
