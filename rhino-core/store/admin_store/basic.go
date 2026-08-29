package admin_store

// THIS FILE WAS AUTO-GENERATED. DO NOT MODIFY.

import (
	"database/sql"
	"github.com/linchunquan/sqlgen/db"
	"rhino-core/schema"
)

const CreateApplicationArchivingCfgItemStmt = `
CREATE TABLE IF NOT EXISTS application_archiving_cfg_items (
 f_id                          BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code                 VARCHAR(32)
,f_business_code               VARCHAR(32)
,f_task_name                   VARCHAR(32)
,f_match_channels              VARCHAR(512)
,f_data_archive_cn_begin_time  VARCHAR(64)
,f_data_archive_cn_latest_time VARCHAR(64)
,f_is_dst_sensitive            BOOLEAN
,f_is_last                     BOOLEAN
);
`

const InsertApplicationArchivingCfgItemStmt = `
INSERT INTO application_archiving_cfg_items (
 f_system_code
,f_business_code
,f_task_name
,f_match_channels
,f_data_archive_cn_begin_time
,f_data_archive_cn_latest_time
,f_is_dst_sensitive
,f_is_last
) VALUES (?,?,?,?,?,?,?,?)
`

const SelectApplicationArchivingCfgItemStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_task_name
,f_match_channels
,f_data_archive_cn_begin_time
,f_data_archive_cn_latest_time
,f_is_dst_sensitive
,f_is_last
FROM application_archiving_cfg_items 
`

const SelectApplicationArchivingCfgItemRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_task_name
,f_match_channels
,f_data_archive_cn_begin_time
,f_data_archive_cn_latest_time
,f_is_dst_sensitive
,f_is_last
FROM application_archiving_cfg_items 
LIMIT ? OFFSET ?
`

const SelectApplicationArchivingCfgItemCountStmt = `
SELECT count(1)
FROM application_archiving_cfg_items 
`

const SelectApplicationArchivingCfgItemByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_task_name
,f_match_channels
,f_data_archive_cn_begin_time
,f_data_archive_cn_latest_time
,f_is_dst_sensitive
,f_is_last
FROM application_archiving_cfg_items 
WHERE f_id=?
`

const UpdateApplicationArchivingCfgItemByIdStmt = `
UPDATE application_archiving_cfg_items SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_task_name=?
,f_match_channels=?
,f_data_archive_cn_begin_time=?
,f_data_archive_cn_latest_time=?
,f_is_dst_sensitive=?
,f_is_last=? 
WHERE f_id=?
`

const DeleteApplicationArchivingCfgItemByIdStmt = `
DELETE FROM application_archiving_cfg_items 
WHERE f_id=?
`

const CreatePkAaciStmt = `
CREATE UNIQUE INDEX pk_aaci ON application_archiving_cfg_items (f_system_code,f_business_code,f_task_name);
`

const SelectApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_task_name
,f_match_channels
,f_data_archive_cn_begin_time
,f_data_archive_cn_latest_time
,f_is_dst_sensitive
,f_is_last
FROM application_archiving_cfg_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_task_name=?
`

const SelectApplicationArchivingCfgItemCountBySystemCodeAndBusinessCodeAndTaskNameStmt = `
SELECT count(1)
FROM application_archiving_cfg_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_task_name=?
`

const UpdateApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt = `
UPDATE application_archiving_cfg_items SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_task_name=?
,f_match_channels=?
,f_data_archive_cn_begin_time=?
,f_data_archive_cn_latest_time=?
,f_is_dst_sensitive=?
,f_is_last=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_task_name=?
`

const DeleteApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt = `
DELETE FROM application_archiving_cfg_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_task_name=?
`

func scanApplicationArchivingCfgItem(row *sql.Row) (*schema.ApplicationArchivingCfgItem, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullBool
	var v8 sql.NullBool

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.ApplicationArchivingCfgItem{}

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
		v.TaskName = v3.String
	} else {
		v.TaskName = ""
	}

	if v4.Valid {
		v.MatchChannels = v4.String
	} else {
		v.MatchChannels = ""
	}

	if v5.Valid {
		v.DataArchiveCnBeginTime = v5.String
	} else {
		v.DataArchiveCnBeginTime = ""
	}

	if v6.Valid {
		v.DataArchiveCnLatestTime = v6.String
	} else {
		v.DataArchiveCnLatestTime = ""
	}

	if v7.Valid {
		v.IsDSTSensitive = v7.Bool
	} else {
		v.IsDSTSensitive = false
	}

	if v8.Valid {
		v.IsLast = v8.Bool
	} else {
		v.IsLast = false
	}

	return v, nil
}

func scanApplicationArchivingCfgItems(rows *sql.Rows) ([]*schema.ApplicationArchivingCfgItem, error) {
	var err error
	var vv []*schema.ApplicationArchivingCfgItem

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullBool
	var v8 sql.NullBool

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.ApplicationArchivingCfgItem{}

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
			v.TaskName = v3.String
		} else {
			v.TaskName = ""
		}

		if v4.Valid {
			v.MatchChannels = v4.String
		} else {
			v.MatchChannels = ""
		}

		if v5.Valid {
			v.DataArchiveCnBeginTime = v5.String
		} else {
			v.DataArchiveCnBeginTime = ""
		}

		if v6.Valid {
			v.DataArchiveCnLatestTime = v6.String
		} else {
			v.DataArchiveCnLatestTime = ""
		}

		if v7.Valid {
			v.IsDSTSensitive = v7.Bool
		} else {
			v.IsDSTSensitive = false
		}

		if v8.Valid {
			v.IsLast = v8.Bool
		} else {
			v.IsLast = false
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceApplicationArchivingCfgItem(v *schema.ApplicationArchivingCfgItem) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 bool
	var v8 bool

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.TaskName
	v4 = v.MatchChannels
	v5 = v.DataArchiveCnBeginTime
	v6 = v.DataArchiveCnLatestTime
	v7 = v.IsDSTSensitive
	v8 = v.IsLast

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
	}
}

func genericSelectApplicationArchivingCfgItem(db db.SimpleDB, query string, args ...interface{}) (*schema.ApplicationArchivingCfgItem, error) {
	row := db.QueryRow(query, args...)
	return scanApplicationArchivingCfgItem(row)
}

func genericSelectApplicationArchivingCfgItems(db db.SimpleDB, query string, args ...interface{}) ([]*schema.ApplicationArchivingCfgItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApplicationArchivingCfgItems(rows)
}

func InsertApplicationArchivingCfgItem(db db.SimpleDB, v *schema.ApplicationArchivingCfgItem) error {

	res, err := db.Exec(InsertApplicationArchivingCfgItemStmt, sliceApplicationArchivingCfgItem(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteApplicationArchivingCfgItemById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteApplicationArchivingCfgItemByIdStmt, args...)
	return err
}

func DeleteApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskName(db db.SimpleDB, systemCode string, businessCode string, taskName string) error {
	args := []interface{}{systemCode, businessCode, taskName}
	_, err := db.Exec(DeleteApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt, args...)
	return err
}

func UpdateApplicationArchivingCfgItemById(db db.SimpleDB, v *schema.ApplicationArchivingCfgItem) error {
	args := sliceApplicationArchivingCfgItem(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateApplicationArchivingCfgItemByIdStmt, args...)
	return err
}

func UpdateApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskName(db db.SimpleDB, v *schema.ApplicationArchivingCfgItem) error {
	args := sliceApplicationArchivingCfgItem(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.TaskName)
	_, err := db.Exec(UpdateApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt, args...)
	return err
}

func GetApplicationArchivingCfgItemById(db db.SimpleDB, iD int64) (*schema.ApplicationArchivingCfgItem, error) {
	args := []interface{}{iD}
	v, err := genericSelectApplicationArchivingCfgItem(db, SelectApplicationArchivingCfgItemByIdStmt, args...)
	return v, err
}

func GetApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskName(db db.SimpleDB, systemCode string, businessCode string, taskName string) (*schema.ApplicationArchivingCfgItem, error) {
	args := []interface{}{systemCode, businessCode, taskName}
	v, err := genericSelectApplicationArchivingCfgItem(db, SelectApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskNameStmt, args...)
	return v, err
}

func FindAllApplicationArchivingCfgItems(db db.SimpleDB) ([]*schema.ApplicationArchivingCfgItem, error) {
	args := []interface{}{}
	v, err := genericSelectApplicationArchivingCfgItems(db, SelectApplicationArchivingCfgItemStmt, args...)
	return v, err
}

func FindAllApplicationArchivingCfgItemsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.ApplicationArchivingCfgItem, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectApplicationArchivingCfgItems(db, SelectApplicationArchivingCfgItemRangeStmt, args...)
	return v, err
}

func CountApplicationArchivingCfgItem(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectApplicationArchivingCfgItemCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountApplicationArchivingCfgItemBySystemCodeAndBusinessCodeAndTaskName(db db.SimpleDB, systemCode string, businessCode string, taskName string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, taskName}
	row := db.QueryRow(SelectApplicationArchivingCfgItemCountBySystemCodeAndBusinessCodeAndTaskNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeAreaStmt = `
CREATE TABLE IF NOT EXISTS trade_areas (
 f_id               BIGINT PRIMARY KEY AUTO_INCREMENT
,f_area_code        VARCHAR(64)
,f_area_zh_name     VARCHAR(128)
,f_area_en_name     VARCHAR(128)
,f_trade_begin_time VARCHAR(64)
,f_trade_end_time   VARCHAR(64)
);
`

const InsertTradeAreaStmt = `
INSERT INTO trade_areas (
 f_area_code
,f_area_zh_name
,f_area_en_name
,f_trade_begin_time
,f_trade_end_time
) VALUES (?,?,?,?,?)
`

const SelectTradeAreaStmt = `
SELECT 
 f_id
,f_area_code
,f_area_zh_name
,f_area_en_name
,f_trade_begin_time
,f_trade_end_time
FROM trade_areas 
`

const SelectTradeAreaRangeStmt = `
SELECT 
 f_id
,f_area_code
,f_area_zh_name
,f_area_en_name
,f_trade_begin_time
,f_trade_end_time
FROM trade_areas 
LIMIT ? OFFSET ?
`

const SelectTradeAreaCountStmt = `
SELECT count(1)
FROM trade_areas 
`

const SelectTradeAreaByIdStmt = `
SELECT 
 f_id
,f_area_code
,f_area_zh_name
,f_area_en_name
,f_trade_begin_time
,f_trade_end_time
FROM trade_areas 
WHERE f_id=?
`

const UpdateTradeAreaByIdStmt = `
UPDATE trade_areas SET 
 f_id=?
,f_area_code=?
,f_area_zh_name=?
,f_area_en_name=?
,f_trade_begin_time=?
,f_trade_end_time=? 
WHERE f_id=?
`

const DeleteTradeAreaByIdStmt = `
DELETE FROM trade_areas 
WHERE f_id=?
`

const CreatePkAreaStmt = `
CREATE UNIQUE INDEX pk_area ON trade_areas (f_area_code);
`

const SelectTradeAreaByAreaCodeStmt = `
SELECT 
 f_id
,f_area_code
,f_area_zh_name
,f_area_en_name
,f_trade_begin_time
,f_trade_end_time
FROM trade_areas 
WHERE f_area_code=?
`

const SelectTradeAreaCountByAreaCodeStmt = `
SELECT count(1)
FROM trade_areas 
WHERE f_area_code=?
`

const UpdateTradeAreaByAreaCodeStmt = `
UPDATE trade_areas SET 
 f_id=?
,f_area_code=?
,f_area_zh_name=?
,f_area_en_name=?
,f_trade_begin_time=?
,f_trade_end_time=? 
WHERE f_area_code=?
`

const DeleteTradeAreaByAreaCodeStmt = `
DELETE FROM trade_areas 
WHERE f_area_code=?
`

func scanTradeArea(row *sql.Row) (*schema.TradeArea, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

	v := &schema.TradeArea{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.AreaCode = v1.String
	} else {
		v.AreaCode = ""
	}

	if v2.Valid {
		v.AreaZhName = v2.String
	} else {
		v.AreaZhName = ""
	}

	if v3.Valid {
		v.AreaEnName = v3.String
	} else {
		v.AreaEnName = ""
	}

	if v4.Valid {
		v.TradeBeginTime = v4.String
	} else {
		v.TradeBeginTime = ""
	}

	if v5.Valid {
		v.TradeEndTime = v5.String
	} else {
		v.TradeEndTime = ""
	}

	return v, nil
}

func scanTradeAreas(rows *sql.Rows) ([]*schema.TradeArea, error) {
	var err error
	var vv []*schema.TradeArea

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

		v := &schema.TradeArea{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.AreaCode = v1.String
		} else {
			v.AreaCode = ""
		}

		if v2.Valid {
			v.AreaZhName = v2.String
		} else {
			v.AreaZhName = ""
		}

		if v3.Valid {
			v.AreaEnName = v3.String
		} else {
			v.AreaEnName = ""
		}

		if v4.Valid {
			v.TradeBeginTime = v4.String
		} else {
			v.TradeBeginTime = ""
		}

		if v5.Valid {
			v.TradeEndTime = v5.String
		} else {
			v.TradeEndTime = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeArea(v *schema.TradeArea) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string

	v0 = v.ID
	v1 = v.AreaCode
	v2 = v.AreaZhName
	v3 = v.AreaEnName
	v4 = v.TradeBeginTime
	v5 = v.TradeEndTime

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
	}
}

func genericSelectTradeArea(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeArea, error) {
	row := db.QueryRow(query, args...)
	return scanTradeArea(row)
}

func genericSelectTradeAreas(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeArea, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeAreas(rows)
}

func InsertTradeArea(db db.SimpleDB, v *schema.TradeArea) error {

	res, err := db.Exec(InsertTradeAreaStmt, sliceTradeArea(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeAreaById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeAreaByIdStmt, args...)
	return err
}

func DeleteTradeAreaByAreaCode(db db.SimpleDB, areaCode string) error {
	args := []interface{}{areaCode}
	_, err := db.Exec(DeleteTradeAreaByAreaCodeStmt, args...)
	return err
}

func UpdateTradeAreaById(db db.SimpleDB, v *schema.TradeArea) error {
	args := sliceTradeArea(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeAreaByIdStmt, args...)
	return err
}

func UpdateTradeAreaByAreaCode(db db.SimpleDB, v *schema.TradeArea) error {
	args := sliceTradeArea(v)
	args = append(args, v.AreaCode)
	_, err := db.Exec(UpdateTradeAreaByAreaCodeStmt, args...)
	return err
}

func GetTradeAreaById(db db.SimpleDB, iD int64) (*schema.TradeArea, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeArea(db, SelectTradeAreaByIdStmt, args...)
	return v, err
}

func GetTradeAreaByAreaCode(db db.SimpleDB, areaCode string) (*schema.TradeArea, error) {
	args := []interface{}{areaCode}
	v, err := genericSelectTradeArea(db, SelectTradeAreaByAreaCodeStmt, args...)
	return v, err
}

func FindAllTradeAreas(db db.SimpleDB) ([]*schema.TradeArea, error) {
	args := []interface{}{}
	v, err := genericSelectTradeAreas(db, SelectTradeAreaStmt, args...)
	return v, err
}

func FindAllTradeAreasInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeArea, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeAreas(db, SelectTradeAreaRangeStmt, args...)
	return v, err
}

func CountTradeArea(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeAreaCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAreaByAreaCode(db db.SimpleDB, areaCode string) (int, error) {
	var count int
	args := []interface{}{areaCode}
	row := db.QueryRow(SelectTradeAreaCountByAreaCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateExtendAttrItemStmt = `
CREATE TABLE IF NOT EXISTS extend_attr_items (
 f_id                    BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code           VARCHAR(32)
,f_business_code         VARCHAR(32)
,f_required              BOOLEAN
,f_attr_name             VARCHAR(32)
,f_attr_zh_name          VARCHAR(128)
,f_attr_value_type       INTEGER
,f_attr_value_len        INTEGER
,f_attr_min_value        DOUBLE
,f_attr_max_value        DOUBLE
,f_attr_value_range_type INTEGER
,f_attr_value_regex      VARCHAR(512)
,f_enum_range            MEDIUMTEXT
,f_index                 BOOLEAN
,f_unique                BOOLEAN
);
`

const InsertExtendAttrItemStmt = `
INSERT INTO extend_attr_items (
 f_system_code
,f_business_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
,f_index
,f_unique
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectExtendAttrItemStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
,f_index
,f_unique
FROM extend_attr_items 
`

const SelectExtendAttrItemRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
,f_index
,f_unique
FROM extend_attr_items 
LIMIT ? OFFSET ?
`

const SelectExtendAttrItemCountStmt = `
SELECT count(1)
FROM extend_attr_items 
`

const SelectExtendAttrItemByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
,f_index
,f_unique
FROM extend_attr_items 
WHERE f_id=?
`

const UpdateExtendAttrItemByIdStmt = `
UPDATE extend_attr_items SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_required=?
,f_attr_name=?
,f_attr_zh_name=?
,f_attr_value_type=?
,f_attr_value_len=?
,f_attr_min_value=?
,f_attr_max_value=?
,f_attr_value_range_type=?
,f_attr_value_regex=?
,f_enum_range=?
,f_index=?
,f_unique=? 
WHERE f_id=?
`

const DeleteExtendAttrItemByIdStmt = `
DELETE FROM extend_attr_items 
WHERE f_id=?
`

const CreateUqEatStmt = `
CREATE UNIQUE INDEX uq_eat ON extend_attr_items (f_system_code,f_business_code,f_attr_name);
`

const SelectExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
,f_index
,f_unique
FROM extend_attr_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_attr_name=?
`

const SelectExtendAttrItemCountBySystemCodeAndBusinessCodeAndAttrNameStmt = `
SELECT count(1)
FROM extend_attr_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_attr_name=?
`

const UpdateExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt = `
UPDATE extend_attr_items SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_required=?
,f_attr_name=?
,f_attr_zh_name=?
,f_attr_value_type=?
,f_attr_value_len=?
,f_attr_min_value=?
,f_attr_max_value=?
,f_attr_value_range_type=?
,f_attr_value_regex=?
,f_enum_range=?
,f_index=?
,f_unique=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_attr_name=?
`

const DeleteExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt = `
DELETE FROM extend_attr_items 
WHERE f_system_code=?
AND f_business_code=?
AND f_attr_name=?
`

func scanExtendAttrItem(row *sql.Row) (*schema.ExtendAttrItem, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullBool
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullBool
	var v14 sql.NullBool

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.ExtendAttrItem{}

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
		v.Required = v3.Bool
	} else {
		v.Required = false
	}

	if v4.Valid {
		v.AttrName = v4.String
	} else {
		v.AttrName = ""
	}

	if v5.Valid {
		v.AttrZhName = v5.String
	} else {
		v.AttrZhName = ""
	}

	if v6.Valid {
		v.AttrValueType = int(v6.Int64)
	} else {
		v.AttrValueType = 0
	}

	if v7.Valid {
		v.AttrValueLen = int(v7.Int64)
	} else {
		v.AttrValueLen = 0
	}

	if v8.Valid {
		v.AttrMinValue = v8.Float64
	} else {
		v.AttrMinValue = 0
	}

	if v9.Valid {
		v.AttrMaxValue = v9.Float64
	} else {
		v.AttrMaxValue = 0
	}

	if v10.Valid {
		v.AttrValueRangeType = int(v10.Int64)
	} else {
		v.AttrValueRangeType = 0
	}

	if v11.Valid {
		v.AttrValueRegex = v11.String
	} else {
		v.AttrValueRegex = ""
	}

	if v12.Valid {
		v.EnumRange = v12.String
	} else {
		v.EnumRange = ""
	}

	if v13.Valid {
		v.Index = v13.Bool
	} else {
		v.Index = false
	}

	if v14.Valid {
		v.Unique = v14.Bool
	} else {
		v.Unique = false
	}

	return v, nil
}

func scanExtendAttrItems(rows *sql.Rows) ([]*schema.ExtendAttrItem, error) {
	var err error
	var vv []*schema.ExtendAttrItem

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullBool
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString
	var v13 sql.NullBool
	var v14 sql.NullBool

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.ExtendAttrItem{}

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
			v.Required = v3.Bool
		} else {
			v.Required = false
		}

		if v4.Valid {
			v.AttrName = v4.String
		} else {
			v.AttrName = ""
		}

		if v5.Valid {
			v.AttrZhName = v5.String
		} else {
			v.AttrZhName = ""
		}

		if v6.Valid {
			v.AttrValueType = int(v6.Int64)
		} else {
			v.AttrValueType = 0
		}

		if v7.Valid {
			v.AttrValueLen = int(v7.Int64)
		} else {
			v.AttrValueLen = 0
		}

		if v8.Valid {
			v.AttrMinValue = v8.Float64
		} else {
			v.AttrMinValue = 0
		}

		if v9.Valid {
			v.AttrMaxValue = v9.Float64
		} else {
			v.AttrMaxValue = 0
		}

		if v10.Valid {
			v.AttrValueRangeType = int(v10.Int64)
		} else {
			v.AttrValueRangeType = 0
		}

		if v11.Valid {
			v.AttrValueRegex = v11.String
		} else {
			v.AttrValueRegex = ""
		}

		if v12.Valid {
			v.EnumRange = v12.String
		} else {
			v.EnumRange = ""
		}

		if v13.Valid {
			v.Index = v13.Bool
		} else {
			v.Index = false
		}

		if v14.Valid {
			v.Unique = v14.Bool
		} else {
			v.Unique = false
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceExtendAttrItem(v *schema.ExtendAttrItem) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 bool
	var v4 string
	var v5 string
	var v6 int
	var v7 int
	var v8 float64
	var v9 float64
	var v10 int
	var v11 string
	var v12 string
	var v13 bool
	var v14 bool

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.Required
	v4 = v.AttrName
	v5 = v.AttrZhName
	v6 = v.AttrValueType
	v7 = v.AttrValueLen
	v8 = v.AttrMinValue
	v9 = v.AttrMaxValue
	v10 = v.AttrValueRangeType
	v11 = v.AttrValueRegex
	v12 = v.EnumRange
	v13 = v.Index
	v14 = v.Unique

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
	}
}

func genericSelectExtendAttrItem(db db.SimpleDB, query string, args ...interface{}) (*schema.ExtendAttrItem, error) {
	row := db.QueryRow(query, args...)
	return scanExtendAttrItem(row)
}

func genericSelectExtendAttrItems(db db.SimpleDB, query string, args ...interface{}) ([]*schema.ExtendAttrItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanExtendAttrItems(rows)
}

func InsertExtendAttrItem(db db.SimpleDB, v *schema.ExtendAttrItem) error {

	res, err := db.Exec(InsertExtendAttrItemStmt, sliceExtendAttrItem(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteExtendAttrItemById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteExtendAttrItemByIdStmt, args...)
	return err
}

func DeleteExtendAttrItemBySystemCodeAndBusinessCodeAndAttrName(db db.SimpleDB, systemCode string, businessCode string, attrName string) error {
	args := []interface{}{systemCode, businessCode, attrName}
	_, err := db.Exec(DeleteExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt, args...)
	return err
}

func UpdateExtendAttrItemById(db db.SimpleDB, v *schema.ExtendAttrItem) error {
	args := sliceExtendAttrItem(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateExtendAttrItemByIdStmt, args...)
	return err
}

func UpdateExtendAttrItemBySystemCodeAndBusinessCodeAndAttrName(db db.SimpleDB, v *schema.ExtendAttrItem) error {
	args := sliceExtendAttrItem(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.AttrName)
	_, err := db.Exec(UpdateExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt, args...)
	return err
}

func GetExtendAttrItemById(db db.SimpleDB, iD int64) (*schema.ExtendAttrItem, error) {
	args := []interface{}{iD}
	v, err := genericSelectExtendAttrItem(db, SelectExtendAttrItemByIdStmt, args...)
	return v, err
}

func GetExtendAttrItemBySystemCodeAndBusinessCodeAndAttrName(db db.SimpleDB, systemCode string, businessCode string, attrName string) (*schema.ExtendAttrItem, error) {
	args := []interface{}{systemCode, businessCode, attrName}
	v, err := genericSelectExtendAttrItem(db, SelectExtendAttrItemBySystemCodeAndBusinessCodeAndAttrNameStmt, args...)
	return v, err
}

func FindAllExtendAttrItems(db db.SimpleDB) ([]*schema.ExtendAttrItem, error) {
	args := []interface{}{}
	v, err := genericSelectExtendAttrItems(db, SelectExtendAttrItemStmt, args...)
	return v, err
}

func FindAllExtendAttrItemsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.ExtendAttrItem, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectExtendAttrItems(db, SelectExtendAttrItemRangeStmt, args...)
	return v, err
}

func CountExtendAttrItem(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectExtendAttrItemCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountExtendAttrItemBySystemCodeAndBusinessCodeAndAttrName(db db.SimpleDB, systemCode string, businessCode string, attrName string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, attrName}
	row := db.QueryRow(SelectExtendAttrItemCountBySystemCodeAndBusinessCodeAndAttrNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateSubOrderProviderStmt = `
CREATE TABLE IF NOT EXISTS sub_order_providers (
 f_id               BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code      VARCHAR(32)
,f_business_code    VARCHAR(32)
,f_provider_code    VARCHAR(32)
,f_provider_zh_name VARCHAR(128)
,f_provider_en_name VARCHAR(32)
,f_description      VARCHAR(512)
,f_invoke_url       VARCHAR(512)
,f_api_token        VARCHAR(256)
);
`

const InsertSubOrderProviderStmt = `
INSERT INTO sub_order_providers (
 f_system_code
,f_business_code
,f_provider_code
,f_provider_zh_name
,f_provider_en_name
,f_description
,f_invoke_url
,f_api_token
) VALUES (?,?,?,?,?,?,?,?)
`

const SelectSubOrderProviderStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_provider_code
,f_provider_zh_name
,f_provider_en_name
,f_description
,f_invoke_url
,f_api_token
FROM sub_order_providers 
`

const SelectSubOrderProviderRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_provider_code
,f_provider_zh_name
,f_provider_en_name
,f_description
,f_invoke_url
,f_api_token
FROM sub_order_providers 
LIMIT ? OFFSET ?
`

const SelectSubOrderProviderCountStmt = `
SELECT count(1)
FROM sub_order_providers 
`

const SelectSubOrderProviderByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_provider_code
,f_provider_zh_name
,f_provider_en_name
,f_description
,f_invoke_url
,f_api_token
FROM sub_order_providers 
WHERE f_id=?
`

const UpdateSubOrderProviderByIdStmt = `
UPDATE sub_order_providers SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_provider_code=?
,f_provider_zh_name=?
,f_provider_en_name=?
,f_description=?
,f_invoke_url=?
,f_api_token=? 
WHERE f_id=?
`

const DeleteSubOrderProviderByIdStmt = `
DELETE FROM sub_order_providers 
WHERE f_id=?
`

const CreateUqSopStmt = `
CREATE UNIQUE INDEX uq_sop ON sub_order_providers (f_system_code,f_business_code,f_provider_code);
`

const SelectSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_provider_code
,f_provider_zh_name
,f_provider_en_name
,f_description
,f_invoke_url
,f_api_token
FROM sub_order_providers 
WHERE f_system_code=?
AND f_business_code=?
AND f_provider_code=?
`

const SelectSubOrderProviderCountBySystemCodeAndBusinessCodeAndProviderCodeStmt = `
SELECT count(1)
FROM sub_order_providers 
WHERE f_system_code=?
AND f_business_code=?
AND f_provider_code=?
`

const UpdateSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt = `
UPDATE sub_order_providers SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_provider_code=?
,f_provider_zh_name=?
,f_provider_en_name=?
,f_description=?
,f_invoke_url=?
,f_api_token=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_provider_code=?
`

const DeleteSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt = `
DELETE FROM sub_order_providers 
WHERE f_system_code=?
AND f_business_code=?
AND f_provider_code=?
`

func scanSubOrderProvider(row *sql.Row) (*schema.SubOrderProvider, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.SubOrderProvider{}

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
		v.ProviderCode = v3.String
	} else {
		v.ProviderCode = ""
	}

	if v4.Valid {
		v.ProviderZhName = v4.String
	} else {
		v.ProviderZhName = ""
	}

	if v5.Valid {
		v.ProviderEnName = v5.String
	} else {
		v.ProviderEnName = ""
	}

	if v6.Valid {
		v.Description = v6.String
	} else {
		v.Description = ""
	}

	if v7.Valid {
		v.InvokeUrl = v7.String
	} else {
		v.InvokeUrl = ""
	}

	if v8.Valid {
		v.ApiToken = v8.String
	} else {
		v.ApiToken = ""
	}

	return v, nil
}

func scanSubOrderProviders(rows *sql.Rows) ([]*schema.SubOrderProvider, error) {
	var err error
	var vv []*schema.SubOrderProvider

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.SubOrderProvider{}

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
			v.ProviderCode = v3.String
		} else {
			v.ProviderCode = ""
		}

		if v4.Valid {
			v.ProviderZhName = v4.String
		} else {
			v.ProviderZhName = ""
		}

		if v5.Valid {
			v.ProviderEnName = v5.String
		} else {
			v.ProviderEnName = ""
		}

		if v6.Valid {
			v.Description = v6.String
		} else {
			v.Description = ""
		}

		if v7.Valid {
			v.InvokeUrl = v7.String
		} else {
			v.InvokeUrl = ""
		}

		if v8.Valid {
			v.ApiToken = v8.String
		} else {
			v.ApiToken = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceSubOrderProvider(v *schema.SubOrderProvider) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ProviderCode
	v4 = v.ProviderZhName
	v5 = v.ProviderEnName
	v6 = v.Description
	v7 = v.InvokeUrl
	v8 = v.ApiToken

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
	}
}

func genericSelectSubOrderProvider(db db.SimpleDB, query string, args ...interface{}) (*schema.SubOrderProvider, error) {
	row := db.QueryRow(query, args...)
	return scanSubOrderProvider(row)
}

func genericSelectSubOrderProviders(db db.SimpleDB, query string, args ...interface{}) ([]*schema.SubOrderProvider, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSubOrderProviders(rows)
}

func InsertSubOrderProvider(db db.SimpleDB, v *schema.SubOrderProvider) error {

	res, err := db.Exec(InsertSubOrderProviderStmt, sliceSubOrderProvider(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteSubOrderProviderById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteSubOrderProviderByIdStmt, args...)
	return err
}

func DeleteSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCode(db db.SimpleDB, systemCode string, businessCode string, providerCode string) error {
	args := []interface{}{systemCode, businessCode, providerCode}
	_, err := db.Exec(DeleteSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt, args...)
	return err
}

func UpdateSubOrderProviderById(db db.SimpleDB, v *schema.SubOrderProvider) error {
	args := sliceSubOrderProvider(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateSubOrderProviderByIdStmt, args...)
	return err
}

func UpdateSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCode(db db.SimpleDB, v *schema.SubOrderProvider) error {
	args := sliceSubOrderProvider(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.ProviderCode)
	_, err := db.Exec(UpdateSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt, args...)
	return err
}

func GetSubOrderProviderById(db db.SimpleDB, iD int64) (*schema.SubOrderProvider, error) {
	args := []interface{}{iD}
	v, err := genericSelectSubOrderProvider(db, SelectSubOrderProviderByIdStmt, args...)
	return v, err
}

func GetSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCode(db db.SimpleDB, systemCode string, businessCode string, providerCode string) (*schema.SubOrderProvider, error) {
	args := []interface{}{systemCode, businessCode, providerCode}
	v, err := genericSelectSubOrderProvider(db, SelectSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCodeStmt, args...)
	return v, err
}

func FindAllSubOrderProviders(db db.SimpleDB) ([]*schema.SubOrderProvider, error) {
	args := []interface{}{}
	v, err := genericSelectSubOrderProviders(db, SelectSubOrderProviderStmt, args...)
	return v, err
}

func FindAllSubOrderProvidersInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.SubOrderProvider, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectSubOrderProviders(db, SelectSubOrderProviderRangeStmt, args...)
	return v, err
}

func CountSubOrderProvider(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectSubOrderProviderCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountSubOrderProviderBySystemCodeAndBusinessCodeAndProviderCode(db db.SimpleDB, systemCode string, businessCode string, providerCode string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, providerCode}
	row := db.QueryRow(SelectSubOrderProviderCountBySystemCodeAndBusinessCodeAndProviderCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateSecurityLibStmt = `
CREATE TABLE IF NOT EXISTS security_libs (
 f_id                           BIGINT PRIMARY KEY AUTO_INCREMENT
,f_security_lib_code            VARCHAR(32)
,f_security_type                VARCHAR(16)
,f_security_lib_zh_name         VARCHAR(128)
,f_security_lib_en_name         VARCHAR(32)
,f_preferred_security_id_source VARCHAR(2)
,f_data_source                  VARCHAR(128)
,f_description                  VARCHAR(512)
,f_last_sync_datetime           VARCHAR(32)
);
`

const InsertSecurityLibStmt = `
INSERT INTO security_libs (
 f_security_lib_code
,f_security_type
,f_security_lib_zh_name
,f_security_lib_en_name
,f_preferred_security_id_source
,f_data_source
,f_description
,f_last_sync_datetime
) VALUES (?,?,?,?,?,?,?,?)
`

const SelectSecurityLibStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_type
,f_security_lib_zh_name
,f_security_lib_en_name
,f_preferred_security_id_source
,f_data_source
,f_description
,f_last_sync_datetime
FROM security_libs 
`

const SelectSecurityLibRangeStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_type
,f_security_lib_zh_name
,f_security_lib_en_name
,f_preferred_security_id_source
,f_data_source
,f_description
,f_last_sync_datetime
FROM security_libs 
LIMIT ? OFFSET ?
`

const SelectSecurityLibCountStmt = `
SELECT count(1)
FROM security_libs 
`

const SelectSecurityLibByIdStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_type
,f_security_lib_zh_name
,f_security_lib_en_name
,f_preferred_security_id_source
,f_data_source
,f_description
,f_last_sync_datetime
FROM security_libs 
WHERE f_id=?
`

const UpdateSecurityLibByIdStmt = `
UPDATE security_libs SET 
 f_id=?
,f_security_lib_code=?
,f_security_type=?
,f_security_lib_zh_name=?
,f_security_lib_en_name=?
,f_preferred_security_id_source=?
,f_data_source=?
,f_description=?
,f_last_sync_datetime=? 
WHERE f_id=?
`

const DeleteSecurityLibByIdStmt = `
DELETE FROM security_libs 
WHERE f_id=?
`

const CreatePkSlStmt = `
CREATE UNIQUE INDEX pk_sl ON security_libs (f_security_lib_code);
`

const SelectSecurityLibBySecurityLibCodeStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_type
,f_security_lib_zh_name
,f_security_lib_en_name
,f_preferred_security_id_source
,f_data_source
,f_description
,f_last_sync_datetime
FROM security_libs 
WHERE f_security_lib_code=?
`

const SelectSecurityLibCountBySecurityLibCodeStmt = `
SELECT count(1)
FROM security_libs 
WHERE f_security_lib_code=?
`

const UpdateSecurityLibBySecurityLibCodeStmt = `
UPDATE security_libs SET 
 f_id=?
,f_security_lib_code=?
,f_security_type=?
,f_security_lib_zh_name=?
,f_security_lib_en_name=?
,f_preferred_security_id_source=?
,f_data_source=?
,f_description=?
,f_last_sync_datetime=? 
WHERE f_security_lib_code=?
`

const DeleteSecurityLibBySecurityLibCodeStmt = `
DELETE FROM security_libs 
WHERE f_security_lib_code=?
`

func scanSecurityLib(row *sql.Row) (*schema.SecurityLib, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.SecurityLib{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.SecurityLibCode = v1.String
	} else {
		v.SecurityLibCode = ""
	}

	if v2.Valid {
		v.SecurityType = v2.String
	} else {
		v.SecurityType = ""
	}

	if v3.Valid {
		v.SecurityLibZhName = v3.String
	} else {
		v.SecurityLibZhName = ""
	}

	if v4.Valid {
		v.SecurityLibEnName = v4.String
	} else {
		v.SecurityLibEnName = ""
	}

	if v5.Valid {
		v.PreferredSecurityIDSource = v5.String
	} else {
		v.PreferredSecurityIDSource = ""
	}

	if v6.Valid {
		v.DataSource = v6.String
	} else {
		v.DataSource = ""
	}

	if v7.Valid {
		v.Description = v7.String
	} else {
		v.Description = ""
	}

	if v8.Valid {
		v.LastSyncDatetime = v8.String
	} else {
		v.LastSyncDatetime = ""
	}

	return v, nil
}

func scanSecurityLibs(rows *sql.Rows) ([]*schema.SecurityLib, error) {
	var err error
	var vv []*schema.SecurityLib

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.SecurityLib{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.SecurityLibCode = v1.String
		} else {
			v.SecurityLibCode = ""
		}

		if v2.Valid {
			v.SecurityType = v2.String
		} else {
			v.SecurityType = ""
		}

		if v3.Valid {
			v.SecurityLibZhName = v3.String
		} else {
			v.SecurityLibZhName = ""
		}

		if v4.Valid {
			v.SecurityLibEnName = v4.String
		} else {
			v.SecurityLibEnName = ""
		}

		if v5.Valid {
			v.PreferredSecurityIDSource = v5.String
		} else {
			v.PreferredSecurityIDSource = ""
		}

		if v6.Valid {
			v.DataSource = v6.String
		} else {
			v.DataSource = ""
		}

		if v7.Valid {
			v.Description = v7.String
		} else {
			v.Description = ""
		}

		if v8.Valid {
			v.LastSyncDatetime = v8.String
		} else {
			v.LastSyncDatetime = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceSecurityLib(v *schema.SecurityLib) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string

	v0 = v.ID
	v1 = v.SecurityLibCode
	v2 = v.SecurityType
	v3 = v.SecurityLibZhName
	v4 = v.SecurityLibEnName
	v5 = v.PreferredSecurityIDSource
	v6 = v.DataSource
	v7 = v.Description
	v8 = v.LastSyncDatetime

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
	}
}

func genericSelectSecurityLib(db db.SimpleDB, query string, args ...interface{}) (*schema.SecurityLib, error) {
	row := db.QueryRow(query, args...)
	return scanSecurityLib(row)
}

func genericSelectSecurityLibs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.SecurityLib, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSecurityLibs(rows)
}

func InsertSecurityLib(db db.SimpleDB, v *schema.SecurityLib) error {

	res, err := db.Exec(InsertSecurityLibStmt, sliceSecurityLib(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteSecurityLibById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteSecurityLibByIdStmt, args...)
	return err
}

func DeleteSecurityLibBySecurityLibCode(db db.SimpleDB, securityLibCode string) error {
	args := []interface{}{securityLibCode}
	_, err := db.Exec(DeleteSecurityLibBySecurityLibCodeStmt, args...)
	return err
}

func UpdateSecurityLibById(db db.SimpleDB, v *schema.SecurityLib) error {
	args := sliceSecurityLib(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateSecurityLibByIdStmt, args...)
	return err
}

func UpdateSecurityLibBySecurityLibCode(db db.SimpleDB, v *schema.SecurityLib) error {
	args := sliceSecurityLib(v)
	args = append(args, v.SecurityLibCode)
	_, err := db.Exec(UpdateSecurityLibBySecurityLibCodeStmt, args...)
	return err
}

func GetSecurityLibById(db db.SimpleDB, iD int64) (*schema.SecurityLib, error) {
	args := []interface{}{iD}
	v, err := genericSelectSecurityLib(db, SelectSecurityLibByIdStmt, args...)
	return v, err
}

func GetSecurityLibBySecurityLibCode(db db.SimpleDB, securityLibCode string) (*schema.SecurityLib, error) {
	args := []interface{}{securityLibCode}
	v, err := genericSelectSecurityLib(db, SelectSecurityLibBySecurityLibCodeStmt, args...)
	return v, err
}

func FindAllSecurityLibs(db db.SimpleDB) ([]*schema.SecurityLib, error) {
	args := []interface{}{}
	v, err := genericSelectSecurityLibs(db, SelectSecurityLibStmt, args...)
	return v, err
}

func FindAllSecurityLibsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.SecurityLib, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectSecurityLibs(db, SelectSecurityLibRangeStmt, args...)
	return v, err
}

func CountSecurityLib(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectSecurityLibCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountSecurityLibBySecurityLibCode(db db.SimpleDB, securityLibCode string) (int, error) {
	var count int
	args := []interface{}{securityLibCode}
	row := db.QueryRow(SelectSecurityLibCountBySecurityLibCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateSecurityItemStmt = `
CREATE TABLE IF NOT EXISTS security_items (
 f_id                          BIGINT PRIMARY KEY AUTO_INCREMENT
,f_security_lib_code           VARCHAR(32)
,f_security_zh_name            VARCHAR(512)
,f_security_en_name            VARCHAR(512)
,f_symbol                      VARCHAR(64)
,f_symbol_sfx                  VARCHAR(8)
,f_security_exchange_symbol    VARCHAR(64)
,f_security_isin               VARCHAR(64)
,f_security_ric                VARCHAR(64)
,f_security_exchange           VARCHAR(8)
,f_security_exchange_region    VARCHAR(4)
,f_contract_multiplier         DOUBLE
,f_currency                    VARCHAR(4)
,f_lot_size                    INTEGER
,f_issue_date                  VARCHAR(16)
,f_contract_month              VARCHAR(12)
,f_expire_date                 VARCHAR(16)
,f_security_type               VARCHAR(16)
,f_underlying_security_code    VARCHAR(64)
,f_underlying_security_zh_name VARCHAR(128)
,f_underlying_security_en_name VARCHAR(64)
,f_put_or_call                 VARCHAR(2)
);
`

const InsertSecurityItemStmt = `
INSERT INTO security_items (
 f_security_lib_code
,f_security_zh_name
,f_security_en_name
,f_symbol
,f_symbol_sfx
,f_security_exchange_symbol
,f_security_isin
,f_security_ric
,f_security_exchange
,f_security_exchange_region
,f_contract_multiplier
,f_currency
,f_lot_size
,f_issue_date
,f_contract_month
,f_expire_date
,f_security_type
,f_underlying_security_code
,f_underlying_security_zh_name
,f_underlying_security_en_name
,f_put_or_call
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectSecurityItemStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_zh_name
,f_security_en_name
,f_symbol
,f_symbol_sfx
,f_security_exchange_symbol
,f_security_isin
,f_security_ric
,f_security_exchange
,f_security_exchange_region
,f_contract_multiplier
,f_currency
,f_lot_size
,f_issue_date
,f_contract_month
,f_expire_date
,f_security_type
,f_underlying_security_code
,f_underlying_security_zh_name
,f_underlying_security_en_name
,f_put_or_call
FROM security_items 
`

const SelectSecurityItemRangeStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_zh_name
,f_security_en_name
,f_symbol
,f_symbol_sfx
,f_security_exchange_symbol
,f_security_isin
,f_security_ric
,f_security_exchange
,f_security_exchange_region
,f_contract_multiplier
,f_currency
,f_lot_size
,f_issue_date
,f_contract_month
,f_expire_date
,f_security_type
,f_underlying_security_code
,f_underlying_security_zh_name
,f_underlying_security_en_name
,f_put_or_call
FROM security_items 
LIMIT ? OFFSET ?
`

const SelectSecurityItemCountStmt = `
SELECT count(1)
FROM security_items 
`

const SelectSecurityItemByIdStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_zh_name
,f_security_en_name
,f_symbol
,f_symbol_sfx
,f_security_exchange_symbol
,f_security_isin
,f_security_ric
,f_security_exchange
,f_security_exchange_region
,f_contract_multiplier
,f_currency
,f_lot_size
,f_issue_date
,f_contract_month
,f_expire_date
,f_security_type
,f_underlying_security_code
,f_underlying_security_zh_name
,f_underlying_security_en_name
,f_put_or_call
FROM security_items 
WHERE f_id=?
`

const UpdateSecurityItemByIdStmt = `
UPDATE security_items SET 
 f_id=?
,f_security_lib_code=?
,f_security_zh_name=?
,f_security_en_name=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_exchange_symbol=?
,f_security_isin=?
,f_security_ric=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_contract_multiplier=?
,f_currency=?
,f_lot_size=?
,f_issue_date=?
,f_contract_month=?
,f_expire_date=?
,f_security_type=?
,f_underlying_security_code=?
,f_underlying_security_zh_name=?
,f_underlying_security_en_name=?
,f_put_or_call=? 
WHERE f_id=?
`

const DeleteSecurityItemByIdStmt = `
DELETE FROM security_items 
WHERE f_id=?
`

const CreatePkSiStmt = `
CREATE UNIQUE INDEX pk_si ON security_items (f_security_lib_code,f_symbol,f_security_exchange_symbol,f_security_isin,f_security_ric);
`

const SelectSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt = `
SELECT 
 f_id
,f_security_lib_code
,f_security_zh_name
,f_security_en_name
,f_symbol
,f_symbol_sfx
,f_security_exchange_symbol
,f_security_isin
,f_security_ric
,f_security_exchange
,f_security_exchange_region
,f_contract_multiplier
,f_currency
,f_lot_size
,f_issue_date
,f_contract_month
,f_expire_date
,f_security_type
,f_underlying_security_code
,f_underlying_security_zh_name
,f_underlying_security_en_name
,f_put_or_call
FROM security_items 
WHERE f_security_lib_code=?
AND f_symbol=?
AND f_security_exchange_symbol=?
AND f_security_isin=?
AND f_security_ric=?
`

const SelectSecurityItemCountBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt = `
SELECT count(1)
FROM security_items 
WHERE f_security_lib_code=?
AND f_symbol=?
AND f_security_exchange_symbol=?
AND f_security_isin=?
AND f_security_ric=?
`

const UpdateSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt = `
UPDATE security_items SET 
 f_id=?
,f_security_lib_code=?
,f_security_zh_name=?
,f_security_en_name=?
,f_symbol=?
,f_symbol_sfx=?
,f_security_exchange_symbol=?
,f_security_isin=?
,f_security_ric=?
,f_security_exchange=?
,f_security_exchange_region=?
,f_contract_multiplier=?
,f_currency=?
,f_lot_size=?
,f_issue_date=?
,f_contract_month=?
,f_expire_date=?
,f_security_type=?
,f_underlying_security_code=?
,f_underlying_security_zh_name=?
,f_underlying_security_en_name=?
,f_put_or_call=? 
WHERE f_security_lib_code=?
AND f_symbol=?
AND f_security_exchange_symbol=?
AND f_security_isin=?
AND f_security_ric=?
`

const DeleteSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt = `
DELETE FROM security_items 
WHERE f_security_lib_code=?
AND f_symbol=?
AND f_security_exchange_symbol=?
AND f_security_isin=?
AND f_security_ric=?
`

func scanSecurityItem(row *sql.Row) (*schema.SecurityItem, error) {
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
	var v11 sql.NullFloat64
	var v12 sql.NullString
	var v13 sql.NullInt64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.SecurityItem{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.SecurityLibCode = v1.String
	} else {
		v.SecurityLibCode = ""
	}

	if v2.Valid {
		v.SecurityZhName = v2.String
	} else {
		v.SecurityZhName = ""
	}

	if v3.Valid {
		v.SecurityEnName = v3.String
	} else {
		v.SecurityEnName = ""
	}

	if v4.Valid {
		v.Symbol = v4.String
	} else {
		v.Symbol = ""
	}

	if v5.Valid {
		v.SymbolSfx = v5.String
	} else {
		v.SymbolSfx = ""
	}

	if v6.Valid {
		v.SecurityExchangeSymbol = v6.String
	} else {
		v.SecurityExchangeSymbol = ""
	}

	if v7.Valid {
		v.SecurityISIN = v7.String
	} else {
		v.SecurityISIN = ""
	}

	if v8.Valid {
		v.SecurityRIC = v8.String
	} else {
		v.SecurityRIC = ""
	}

	if v9.Valid {
		v.SecurityExchange = v9.String
	} else {
		v.SecurityExchange = ""
	}

	if v10.Valid {
		v.SecurityExchangeRegion = v10.String
	} else {
		v.SecurityExchangeRegion = ""
	}

	if v11.Valid {
		v.ContractMultiplier = v11.Float64
	} else {
		v.ContractMultiplier = 0
	}

	if v12.Valid {
		v.Currency = v12.String
	} else {
		v.Currency = ""
	}

	if v13.Valid {
		v.LotSize = int(v13.Int64)
	} else {
		v.LotSize = 0
	}

	if v14.Valid {
		v.IssueDate = v14.String
	} else {
		v.IssueDate = ""
	}

	if v15.Valid {
		v.ContractMonth = v15.String
	} else {
		v.ContractMonth = ""
	}

	if v16.Valid {
		v.ExpireDate = v16.String
	} else {
		v.ExpireDate = ""
	}

	if v17.Valid {
		v.SecurityType = v17.String
	} else {
		v.SecurityType = ""
	}

	if v18.Valid {
		v.UnderlyingSecurityCode = v18.String
	} else {
		v.UnderlyingSecurityCode = ""
	}

	if v19.Valid {
		v.UnderlyingSecurityZhName = v19.String
	} else {
		v.UnderlyingSecurityZhName = ""
	}

	if v20.Valid {
		v.UnderlyingSecurityEnName = v20.String
	} else {
		v.UnderlyingSecurityEnName = ""
	}

	if v21.Valid {
		v.PutOrCall = v21.String
	} else {
		v.PutOrCall = ""
	}

	return v, nil
}

func scanSecurityItems(rows *sql.Rows) ([]*schema.SecurityItem, error) {
	var err error
	var vv []*schema.SecurityItem

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
	var v11 sql.NullFloat64
	var v12 sql.NullString
	var v13 sql.NullInt64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullString
	var v18 sql.NullString
	var v19 sql.NullString
	var v20 sql.NullString
	var v21 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.SecurityItem{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.SecurityLibCode = v1.String
		} else {
			v.SecurityLibCode = ""
		}

		if v2.Valid {
			v.SecurityZhName = v2.String
		} else {
			v.SecurityZhName = ""
		}

		if v3.Valid {
			v.SecurityEnName = v3.String
		} else {
			v.SecurityEnName = ""
		}

		if v4.Valid {
			v.Symbol = v4.String
		} else {
			v.Symbol = ""
		}

		if v5.Valid {
			v.SymbolSfx = v5.String
		} else {
			v.SymbolSfx = ""
		}

		if v6.Valid {
			v.SecurityExchangeSymbol = v6.String
		} else {
			v.SecurityExchangeSymbol = ""
		}

		if v7.Valid {
			v.SecurityISIN = v7.String
		} else {
			v.SecurityISIN = ""
		}

		if v8.Valid {
			v.SecurityRIC = v8.String
		} else {
			v.SecurityRIC = ""
		}

		if v9.Valid {
			v.SecurityExchange = v9.String
		} else {
			v.SecurityExchange = ""
		}

		if v10.Valid {
			v.SecurityExchangeRegion = v10.String
		} else {
			v.SecurityExchangeRegion = ""
		}

		if v11.Valid {
			v.ContractMultiplier = v11.Float64
		} else {
			v.ContractMultiplier = 0
		}

		if v12.Valid {
			v.Currency = v12.String
		} else {
			v.Currency = ""
		}

		if v13.Valid {
			v.LotSize = int(v13.Int64)
		} else {
			v.LotSize = 0
		}

		if v14.Valid {
			v.IssueDate = v14.String
		} else {
			v.IssueDate = ""
		}

		if v15.Valid {
			v.ContractMonth = v15.String
		} else {
			v.ContractMonth = ""
		}

		if v16.Valid {
			v.ExpireDate = v16.String
		} else {
			v.ExpireDate = ""
		}

		if v17.Valid {
			v.SecurityType = v17.String
		} else {
			v.SecurityType = ""
		}

		if v18.Valid {
			v.UnderlyingSecurityCode = v18.String
		} else {
			v.UnderlyingSecurityCode = ""
		}

		if v19.Valid {
			v.UnderlyingSecurityZhName = v19.String
		} else {
			v.UnderlyingSecurityZhName = ""
		}

		if v20.Valid {
			v.UnderlyingSecurityEnName = v20.String
		} else {
			v.UnderlyingSecurityEnName = ""
		}

		if v21.Valid {
			v.PutOrCall = v21.String
		} else {
			v.PutOrCall = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceSecurityItem(v *schema.SecurityItem) []interface{} {
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
	var v11 float64
	var v12 string
	var v13 int
	var v14 string
	var v15 string
	var v16 string
	var v17 string
	var v18 string
	var v19 string
	var v20 string
	var v21 string

	v0 = v.ID
	v1 = v.SecurityLibCode
	v2 = v.SecurityZhName
	v3 = v.SecurityEnName
	v4 = v.Symbol
	v5 = v.SymbolSfx
	v6 = v.SecurityExchangeSymbol
	v7 = v.SecurityISIN
	v8 = v.SecurityRIC
	v9 = v.SecurityExchange
	v10 = v.SecurityExchangeRegion
	v11 = v.ContractMultiplier
	v12 = v.Currency
	v13 = v.LotSize
	v14 = v.IssueDate
	v15 = v.ContractMonth
	v16 = v.ExpireDate
	v17 = v.SecurityType
	v18 = v.UnderlyingSecurityCode
	v19 = v.UnderlyingSecurityZhName
	v20 = v.UnderlyingSecurityEnName
	v21 = v.PutOrCall

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
	}
}

func genericSelectSecurityItem(db db.SimpleDB, query string, args ...interface{}) (*schema.SecurityItem, error) {
	row := db.QueryRow(query, args...)
	return scanSecurityItem(row)
}

func genericSelectSecurityItems(db db.SimpleDB, query string, args ...interface{}) ([]*schema.SecurityItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSecurityItems(rows)
}

func InsertSecurityItem(db db.SimpleDB, v *schema.SecurityItem) error {

	res, err := db.Exec(InsertSecurityItemStmt, sliceSecurityItem(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteSecurityItemById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteSecurityItemByIdStmt, args...)
	return err
}

func DeleteSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRic(db db.SimpleDB, securityLibCode string, symbol string, securityExchangeSymbol string, securityISIN string, securityRIC string) error {
	args := []interface{}{securityLibCode, symbol, securityExchangeSymbol, securityISIN, securityRIC}
	_, err := db.Exec(DeleteSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt, args...)
	return err
}

func UpdateSecurityItemById(db db.SimpleDB, v *schema.SecurityItem) error {
	args := sliceSecurityItem(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateSecurityItemByIdStmt, args...)
	return err
}

func UpdateSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRic(db db.SimpleDB, v *schema.SecurityItem) error {
	args := sliceSecurityItem(v)
	args = append(args, v.SecurityLibCode, v.Symbol, v.SecurityExchangeSymbol, v.SecurityISIN, v.SecurityRIC)
	_, err := db.Exec(UpdateSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt, args...)
	return err
}

func GetSecurityItemById(db db.SimpleDB, iD int64) (*schema.SecurityItem, error) {
	args := []interface{}{iD}
	v, err := genericSelectSecurityItem(db, SelectSecurityItemByIdStmt, args...)
	return v, err
}

func GetSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRic(db db.SimpleDB, securityLibCode string, symbol string, securityExchangeSymbol string, securityISIN string, securityRIC string) (*schema.SecurityItem, error) {
	args := []interface{}{securityLibCode, symbol, securityExchangeSymbol, securityISIN, securityRIC}
	v, err := genericSelectSecurityItem(db, SelectSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt, args...)
	return v, err
}

func FindAllSecurityItems(db db.SimpleDB) ([]*schema.SecurityItem, error) {
	args := []interface{}{}
	v, err := genericSelectSecurityItems(db, SelectSecurityItemStmt, args...)
	return v, err
}

func FindAllSecurityItemsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.SecurityItem, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectSecurityItems(db, SelectSecurityItemRangeStmt, args...)
	return v, err
}

func CountSecurityItem(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectSecurityItemCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountSecurityItemBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRic(db db.SimpleDB, securityLibCode string, symbol string, securityExchangeSymbol string, securityISIN string, securityRIC string) (int, error) {
	var count int
	args := []interface{}{securityLibCode, symbol, securityExchangeSymbol, securityISIN, securityRIC}
	row := db.QueryRow(SelectSecurityItemCountBySecurityLibCodeAndSymbolAndSecurityExchangeSymbolAndSecurityIsinAndSecurityRicStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateApplicationSecurityLibStmt = `
CREATE TABLE IF NOT EXISTS application_security_libs (
 f_id                BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code       VARCHAR(32)
,f_business_code     VARCHAR(32)
,f_security_lib_code VARCHAR(32)
);
`

const InsertApplicationSecurityLibStmt = `
INSERT INTO application_security_libs (
 f_system_code
,f_business_code
,f_security_lib_code
) VALUES (?,?,?)
`

const SelectApplicationSecurityLibStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_security_lib_code
FROM application_security_libs 
`

const SelectApplicationSecurityLibRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_security_lib_code
FROM application_security_libs 
LIMIT ? OFFSET ?
`

const SelectApplicationSecurityLibCountStmt = `
SELECT count(1)
FROM application_security_libs 
`

const SelectApplicationSecurityLibByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_security_lib_code
FROM application_security_libs 
WHERE f_id=?
`

const UpdateApplicationSecurityLibByIdStmt = `
UPDATE application_security_libs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_security_lib_code=? 
WHERE f_id=?
`

const DeleteApplicationSecurityLibByIdStmt = `
DELETE FROM application_security_libs 
WHERE f_id=?
`

const CreateUqAslStmt = `
CREATE UNIQUE INDEX uq_asl ON application_security_libs (f_system_code,f_business_code,f_security_lib_code);
`

const SelectApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_security_lib_code
FROM application_security_libs 
WHERE f_system_code=?
AND f_business_code=?
AND f_security_lib_code=?
`

const SelectApplicationSecurityLibCountBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt = `
SELECT count(1)
FROM application_security_libs 
WHERE f_system_code=?
AND f_business_code=?
AND f_security_lib_code=?
`

const UpdateApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt = `
UPDATE application_security_libs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_security_lib_code=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_security_lib_code=?
`

const DeleteApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt = `
DELETE FROM application_security_libs 
WHERE f_system_code=?
AND f_business_code=?
AND f_security_lib_code=?
`

func scanApplicationSecurityLib(row *sql.Row) (*schema.ApplicationSecurityLib, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.ApplicationSecurityLib{}

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
		v.SecurityLibCode = v3.String
	} else {
		v.SecurityLibCode = ""
	}

	return v, nil
}

func scanApplicationSecurityLibs(rows *sql.Rows) ([]*schema.ApplicationSecurityLib, error) {
	var err error
	var vv []*schema.ApplicationSecurityLib

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.ApplicationSecurityLib{}

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
			v.SecurityLibCode = v3.String
		} else {
			v.SecurityLibCode = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceApplicationSecurityLib(v *schema.ApplicationSecurityLib) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.SecurityLibCode

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
	}
}

func genericSelectApplicationSecurityLib(db db.SimpleDB, query string, args ...interface{}) (*schema.ApplicationSecurityLib, error) {
	row := db.QueryRow(query, args...)
	return scanApplicationSecurityLib(row)
}

func genericSelectApplicationSecurityLibs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.ApplicationSecurityLib, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApplicationSecurityLibs(rows)
}

func InsertApplicationSecurityLib(db db.SimpleDB, v *schema.ApplicationSecurityLib) error {

	res, err := db.Exec(InsertApplicationSecurityLibStmt, sliceApplicationSecurityLib(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteApplicationSecurityLibById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteApplicationSecurityLibByIdStmt, args...)
	return err
}

func DeleteApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCode(db db.SimpleDB, systemCode string, businessCode string, securityLibCode string) error {
	args := []interface{}{systemCode, businessCode, securityLibCode}
	_, err := db.Exec(DeleteApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt, args...)
	return err
}

func UpdateApplicationSecurityLibById(db db.SimpleDB, v *schema.ApplicationSecurityLib) error {
	args := sliceApplicationSecurityLib(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateApplicationSecurityLibByIdStmt, args...)
	return err
}

func UpdateApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCode(db db.SimpleDB, v *schema.ApplicationSecurityLib) error {
	args := sliceApplicationSecurityLib(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.SecurityLibCode)
	_, err := db.Exec(UpdateApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt, args...)
	return err
}

func GetApplicationSecurityLibById(db db.SimpleDB, iD int64) (*schema.ApplicationSecurityLib, error) {
	args := []interface{}{iD}
	v, err := genericSelectApplicationSecurityLib(db, SelectApplicationSecurityLibByIdStmt, args...)
	return v, err
}

func GetApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCode(db db.SimpleDB, systemCode string, businessCode string, securityLibCode string) (*schema.ApplicationSecurityLib, error) {
	args := []interface{}{systemCode, businessCode, securityLibCode}
	v, err := genericSelectApplicationSecurityLib(db, SelectApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt, args...)
	return v, err
}

func FindAllApplicationSecurityLibs(db db.SimpleDB) ([]*schema.ApplicationSecurityLib, error) {
	args := []interface{}{}
	v, err := genericSelectApplicationSecurityLibs(db, SelectApplicationSecurityLibStmt, args...)
	return v, err
}

func FindAllApplicationSecurityLibsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.ApplicationSecurityLib, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectApplicationSecurityLibs(db, SelectApplicationSecurityLibRangeStmt, args...)
	return v, err
}

func CountApplicationSecurityLib(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectApplicationSecurityLibCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountApplicationSecurityLibBySystemCodeAndBusinessCodeAndSecurityLibCode(db db.SimpleDB, systemCode string, businessCode string, securityLibCode string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, securityLibCode}
	row := db.QueryRow(SelectApplicationSecurityLibCountBySystemCodeAndBusinessCodeAndSecurityLibCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeChannelStmt = `
CREATE TABLE IF NOT EXISTS trade_channels (
 f_id                    BIGINT PRIMARY KEY AUTO_INCREMENT
,f_channel_code          VARCHAR(32)
,f_channel_zh_name       VARCHAR(128)
,f_channel_en_name       VARCHAR(32)
,f_channel_protocol_type VARCHAR(32)
,f_channel_adapter_name  VARCHAR(64)
,f_description           VARCHAR(512)
,f_addresses             VARCHAR(512)
,f_exchange              VARCHAR(8)
,f_time_zone             VARCHAR(512)
,f_display_time_zone     VARCHAR(8)
,f_begin_time            VARCHAR(8)
,f_end_time              VARCHAR(8)
,f_data_num_adj          INTEGER
,f_active_real_address   VARCHAR(64)
,f_real_address          VARCHAR(64)
,f_export_address        VARCHAR(64)
,f_export_http_port      INTEGER
,f_export_ws_port        INTEGER
,f_api_token             VARCHAR(256)
,f_status                INTEGER
,f_offline_reason        VARCHAR(512)
,f_config_dir            VARCHAR(512)
,f_adapter_path          VARCHAR(128)
);
`

const InsertTradeChannelStmt = `
INSERT INTO trade_channels (
 f_channel_code
,f_channel_zh_name
,f_channel_en_name
,f_channel_protocol_type
,f_channel_adapter_name
,f_description
,f_addresses
,f_exchange
,f_time_zone
,f_display_time_zone
,f_begin_time
,f_end_time
,f_data_num_adj
,f_active_real_address
,f_real_address
,f_export_address
,f_export_http_port
,f_export_ws_port
,f_api_token
,f_status
,f_offline_reason
,f_config_dir
,f_adapter_path
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeChannelStmt = `
SELECT 
 f_id
,f_channel_code
,f_channel_zh_name
,f_channel_en_name
,f_channel_protocol_type
,f_channel_adapter_name
,f_description
,f_addresses
,f_exchange
,f_time_zone
,f_display_time_zone
,f_begin_time
,f_end_time
,f_data_num_adj
,f_active_real_address
,f_real_address
,f_export_address
,f_export_http_port
,f_export_ws_port
,f_api_token
,f_status
,f_offline_reason
,f_config_dir
,f_adapter_path
FROM trade_channels 
`

const SelectTradeChannelRangeStmt = `
SELECT 
 f_id
,f_channel_code
,f_channel_zh_name
,f_channel_en_name
,f_channel_protocol_type
,f_channel_adapter_name
,f_description
,f_addresses
,f_exchange
,f_time_zone
,f_display_time_zone
,f_begin_time
,f_end_time
,f_data_num_adj
,f_active_real_address
,f_real_address
,f_export_address
,f_export_http_port
,f_export_ws_port
,f_api_token
,f_status
,f_offline_reason
,f_config_dir
,f_adapter_path
FROM trade_channels 
LIMIT ? OFFSET ?
`

const SelectTradeChannelCountStmt = `
SELECT count(1)
FROM trade_channels 
`

const SelectTradeChannelByIdStmt = `
SELECT 
 f_id
,f_channel_code
,f_channel_zh_name
,f_channel_en_name
,f_channel_protocol_type
,f_channel_adapter_name
,f_description
,f_addresses
,f_exchange
,f_time_zone
,f_display_time_zone
,f_begin_time
,f_end_time
,f_data_num_adj
,f_active_real_address
,f_real_address
,f_export_address
,f_export_http_port
,f_export_ws_port
,f_api_token
,f_status
,f_offline_reason
,f_config_dir
,f_adapter_path
FROM trade_channels 
WHERE f_id=?
`

const UpdateTradeChannelByIdStmt = `
UPDATE trade_channels SET 
 f_id=?
,f_channel_code=?
,f_channel_zh_name=?
,f_channel_en_name=?
,f_channel_protocol_type=?
,f_channel_adapter_name=?
,f_description=?
,f_addresses=?
,f_exchange=?
,f_time_zone=?
,f_display_time_zone=?
,f_begin_time=?
,f_end_time=?
,f_data_num_adj=?
,f_active_real_address=?
,f_real_address=?
,f_export_address=?
,f_export_http_port=?
,f_export_ws_port=?
,f_api_token=?
,f_status=?
,f_offline_reason=?
,f_config_dir=?
,f_adapter_path=? 
WHERE f_id=?
`

const DeleteTradeChannelByIdStmt = `
DELETE FROM trade_channels 
WHERE f_id=?
`

const CreatePkTcStmt = `
CREATE UNIQUE INDEX pk_tc ON trade_channels (f_channel_code);
`

const SelectTradeChannelByChannelCodeStmt = `
SELECT 
 f_id
,f_channel_code
,f_channel_zh_name
,f_channel_en_name
,f_channel_protocol_type
,f_channel_adapter_name
,f_description
,f_addresses
,f_exchange
,f_time_zone
,f_display_time_zone
,f_begin_time
,f_end_time
,f_data_num_adj
,f_active_real_address
,f_real_address
,f_export_address
,f_export_http_port
,f_export_ws_port
,f_api_token
,f_status
,f_offline_reason
,f_config_dir
,f_adapter_path
FROM trade_channels 
WHERE f_channel_code=?
`

const SelectTradeChannelCountByChannelCodeStmt = `
SELECT count(1)
FROM trade_channels 
WHERE f_channel_code=?
`

const UpdateTradeChannelByChannelCodeStmt = `
UPDATE trade_channels SET 
 f_id=?
,f_channel_code=?
,f_channel_zh_name=?
,f_channel_en_name=?
,f_channel_protocol_type=?
,f_channel_adapter_name=?
,f_description=?
,f_addresses=?
,f_exchange=?
,f_time_zone=?
,f_display_time_zone=?
,f_begin_time=?
,f_end_time=?
,f_data_num_adj=?
,f_active_real_address=?
,f_real_address=?
,f_export_address=?
,f_export_http_port=?
,f_export_ws_port=?
,f_api_token=?
,f_status=?
,f_offline_reason=?
,f_config_dir=?
,f_adapter_path=? 
WHERE f_channel_code=?
`

const DeleteTradeChannelByChannelCodeStmt = `
DELETE FROM trade_channels 
WHERE f_channel_code=?
`

func scanTradeChannel(row *sql.Row) (*schema.TradeChannel, error) {
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
	var v13 sql.NullInt64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullInt64
	var v18 sql.NullInt64
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeChannel{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ChannelCode = v1.String
	} else {
		v.ChannelCode = ""
	}

	if v2.Valid {
		v.ChannelZhName = v2.String
	} else {
		v.ChannelZhName = ""
	}

	if v3.Valid {
		v.ChannelEnName = v3.String
	} else {
		v.ChannelEnName = ""
	}

	if v4.Valid {
		v.ChannelProtocolType = v4.String
	} else {
		v.ChannelProtocolType = ""
	}

	if v5.Valid {
		v.ChannelAdapterName = v5.String
	} else {
		v.ChannelAdapterName = ""
	}

	if v6.Valid {
		v.Description = v6.String
	} else {
		v.Description = ""
	}

	if v7.Valid {
		v.Addresses = v7.String
	} else {
		v.Addresses = ""
	}

	if v8.Valid {
		v.Exchange = v8.String
	} else {
		v.Exchange = ""
	}

	if v9.Valid {
		v.TimeZone = v9.String
	} else {
		v.TimeZone = ""
	}

	if v10.Valid {
		v.DisplayTimeZone = v10.String
	} else {
		v.DisplayTimeZone = ""
	}

	if v11.Valid {
		v.BeginTime = v11.String
	} else {
		v.BeginTime = ""
	}

	if v12.Valid {
		v.EndTime = v12.String
	} else {
		v.EndTime = ""
	}

	if v13.Valid {
		v.DataNumAdj = int(v13.Int64)
	} else {
		v.DataNumAdj = 0
	}

	if v14.Valid {
		v.ActiveRealAddress = v14.String
	} else {
		v.ActiveRealAddress = ""
	}

	if v15.Valid {
		v.RealAddress = v15.String
	} else {
		v.RealAddress = ""
	}

	if v16.Valid {
		v.ExportAddress = v16.String
	} else {
		v.ExportAddress = ""
	}

	if v17.Valid {
		v.ExportHttpPort = int(v17.Int64)
	} else {
		v.ExportHttpPort = 0
	}

	if v18.Valid {
		v.ExportWSPort = int(v18.Int64)
	} else {
		v.ExportWSPort = 0
	}

	if v19.Valid {
		v.ApiToken = v19.String
	} else {
		v.ApiToken = ""
	}

	if v20.Valid {
		v.Status = int(v20.Int64)
	} else {
		v.Status = 0
	}

	if v21.Valid {
		v.OfflineReason = v21.String
	} else {
		v.OfflineReason = ""
	}

	if v22.Valid {
		v.ConfigDir = v22.String
	} else {
		v.ConfigDir = ""
	}

	if v23.Valid {
		v.AdapterPath = v23.String
	} else {
		v.AdapterPath = ""
	}

	return v, nil
}

func scanTradeChannels(rows *sql.Rows) ([]*schema.TradeChannel, error) {
	var err error
	var vv []*schema.TradeChannel

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
	var v13 sql.NullInt64
	var v14 sql.NullString
	var v15 sql.NullString
	var v16 sql.NullString
	var v17 sql.NullInt64
	var v18 sql.NullInt64
	var v19 sql.NullString
	var v20 sql.NullInt64
	var v21 sql.NullString
	var v22 sql.NullString
	var v23 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeChannel{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ChannelCode = v1.String
		} else {
			v.ChannelCode = ""
		}

		if v2.Valid {
			v.ChannelZhName = v2.String
		} else {
			v.ChannelZhName = ""
		}

		if v3.Valid {
			v.ChannelEnName = v3.String
		} else {
			v.ChannelEnName = ""
		}

		if v4.Valid {
			v.ChannelProtocolType = v4.String
		} else {
			v.ChannelProtocolType = ""
		}

		if v5.Valid {
			v.ChannelAdapterName = v5.String
		} else {
			v.ChannelAdapterName = ""
		}

		if v6.Valid {
			v.Description = v6.String
		} else {
			v.Description = ""
		}

		if v7.Valid {
			v.Addresses = v7.String
		} else {
			v.Addresses = ""
		}

		if v8.Valid {
			v.Exchange = v8.String
		} else {
			v.Exchange = ""
		}

		if v9.Valid {
			v.TimeZone = v9.String
		} else {
			v.TimeZone = ""
		}

		if v10.Valid {
			v.DisplayTimeZone = v10.String
		} else {
			v.DisplayTimeZone = ""
		}

		if v11.Valid {
			v.BeginTime = v11.String
		} else {
			v.BeginTime = ""
		}

		if v12.Valid {
			v.EndTime = v12.String
		} else {
			v.EndTime = ""
		}

		if v13.Valid {
			v.DataNumAdj = int(v13.Int64)
		} else {
			v.DataNumAdj = 0
		}

		if v14.Valid {
			v.ActiveRealAddress = v14.String
		} else {
			v.ActiveRealAddress = ""
		}

		if v15.Valid {
			v.RealAddress = v15.String
		} else {
			v.RealAddress = ""
		}

		if v16.Valid {
			v.ExportAddress = v16.String
		} else {
			v.ExportAddress = ""
		}

		if v17.Valid {
			v.ExportHttpPort = int(v17.Int64)
		} else {
			v.ExportHttpPort = 0
		}

		if v18.Valid {
			v.ExportWSPort = int(v18.Int64)
		} else {
			v.ExportWSPort = 0
		}

		if v19.Valid {
			v.ApiToken = v19.String
		} else {
			v.ApiToken = ""
		}

		if v20.Valid {
			v.Status = int(v20.Int64)
		} else {
			v.Status = 0
		}

		if v21.Valid {
			v.OfflineReason = v21.String
		} else {
			v.OfflineReason = ""
		}

		if v22.Valid {
			v.ConfigDir = v22.String
		} else {
			v.ConfigDir = ""
		}

		if v23.Valid {
			v.AdapterPath = v23.String
		} else {
			v.AdapterPath = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeChannel(v *schema.TradeChannel) []interface{} {
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
	var v13 int
	var v14 string
	var v15 string
	var v16 string
	var v17 int
	var v18 int
	var v19 string
	var v20 int
	var v21 string
	var v22 string
	var v23 string

	v0 = v.ID
	v1 = v.ChannelCode
	v2 = v.ChannelZhName
	v3 = v.ChannelEnName
	v4 = v.ChannelProtocolType
	v5 = v.ChannelAdapterName
	v6 = v.Description
	v7 = v.Addresses
	v8 = v.Exchange
	v9 = v.TimeZone
	v10 = v.DisplayTimeZone
	v11 = v.BeginTime
	v12 = v.EndTime
	v13 = v.DataNumAdj
	v14 = v.ActiveRealAddress
	v15 = v.RealAddress
	v16 = v.ExportAddress
	v17 = v.ExportHttpPort
	v18 = v.ExportWSPort
	v19 = v.ApiToken
	v20 = v.Status
	v21 = v.OfflineReason
	v22 = v.ConfigDir
	v23 = v.AdapterPath

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
	}
}

func genericSelectTradeChannel(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeChannel, error) {
	row := db.QueryRow(query, args...)
	return scanTradeChannel(row)
}

func genericSelectTradeChannels(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeChannel, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeChannels(rows)
}

func InsertTradeChannel(db db.SimpleDB, v *schema.TradeChannel) error {

	res, err := db.Exec(InsertTradeChannelStmt, sliceTradeChannel(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeChannelById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeChannelByIdStmt, args...)
	return err
}

func DeleteTradeChannelByChannelCode(db db.SimpleDB, channelCode string) error {
	args := []interface{}{channelCode}
	_, err := db.Exec(DeleteTradeChannelByChannelCodeStmt, args...)
	return err
}

func UpdateTradeChannelById(db db.SimpleDB, v *schema.TradeChannel) error {
	args := sliceTradeChannel(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeChannelByIdStmt, args...)
	return err
}

func UpdateTradeChannelByChannelCode(db db.SimpleDB, v *schema.TradeChannel) error {
	args := sliceTradeChannel(v)
	args = append(args, v.ChannelCode)
	_, err := db.Exec(UpdateTradeChannelByChannelCodeStmt, args...)
	return err
}

func GetTradeChannelById(db db.SimpleDB, iD int64) (*schema.TradeChannel, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeChannel(db, SelectTradeChannelByIdStmt, args...)
	return v, err
}

func GetTradeChannelByChannelCode(db db.SimpleDB, channelCode string) (*schema.TradeChannel, error) {
	args := []interface{}{channelCode}
	v, err := genericSelectTradeChannel(db, SelectTradeChannelByChannelCodeStmt, args...)
	return v, err
}

func FindAllTradeChannels(db db.SimpleDB) ([]*schema.TradeChannel, error) {
	args := []interface{}{}
	v, err := genericSelectTradeChannels(db, SelectTradeChannelStmt, args...)
	return v, err
}

func FindAllTradeChannelsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeChannel, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeChannels(db, SelectTradeChannelRangeStmt, args...)
	return v, err
}

func CountTradeChannel(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeChannelCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeChannelByChannelCode(db db.SimpleDB, channelCode string) (int, error) {
	var count int
	args := []interface{}{channelCode}
	row := db.QueryRow(SelectTradeChannelCountByChannelCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeChannelCfgItemStmt = `
CREATE TABLE IF NOT EXISTS trade_channel_cfg_items (
 f_id                        BIGINT PRIMARY KEY AUTO_INCREMENT
,f_channel_code              VARCHAR(32)
,f_config_item_name          VARCHAR(32)
,f_config_item_value         VARCHAR(512)
,f_config_item_default_value VARCHAR(512)
,f_description               VARCHAR(512)
,f_required                  INTEGER
);
`

const InsertTradeChannelCfgItemStmt = `
INSERT INTO trade_channel_cfg_items (
 f_channel_code
,f_config_item_name
,f_config_item_value
,f_config_item_default_value
,f_description
,f_required
) VALUES (?,?,?,?,?,?)
`

const SelectTradeChannelCfgItemStmt = `
SELECT 
 f_id
,f_channel_code
,f_config_item_name
,f_config_item_value
,f_config_item_default_value
,f_description
,f_required
FROM trade_channel_cfg_items 
`

const SelectTradeChannelCfgItemRangeStmt = `
SELECT 
 f_id
,f_channel_code
,f_config_item_name
,f_config_item_value
,f_config_item_default_value
,f_description
,f_required
FROM trade_channel_cfg_items 
LIMIT ? OFFSET ?
`

const SelectTradeChannelCfgItemCountStmt = `
SELECT count(1)
FROM trade_channel_cfg_items 
`

const SelectTradeChannelCfgItemByIdStmt = `
SELECT 
 f_id
,f_channel_code
,f_config_item_name
,f_config_item_value
,f_config_item_default_value
,f_description
,f_required
FROM trade_channel_cfg_items 
WHERE f_id=?
`

const UpdateTradeChannelCfgItemByIdStmt = `
UPDATE trade_channel_cfg_items SET 
 f_id=?
,f_channel_code=?
,f_config_item_name=?
,f_config_item_value=?
,f_config_item_default_value=?
,f_description=?
,f_required=? 
WHERE f_id=?
`

const DeleteTradeChannelCfgItemByIdStmt = `
DELETE FROM trade_channel_cfg_items 
WHERE f_id=?
`

const CreatePkTcciStmt = `
CREATE UNIQUE INDEX pk_tcci ON trade_channel_cfg_items (f_channel_code,f_config_item_name);
`

const SelectTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt = `
SELECT 
 f_id
,f_channel_code
,f_config_item_name
,f_config_item_value
,f_config_item_default_value
,f_description
,f_required
FROM trade_channel_cfg_items 
WHERE f_channel_code=?
AND f_config_item_name=?
`

const SelectTradeChannelCfgItemCountByChannelCodeAndConfigItemNameStmt = `
SELECT count(1)
FROM trade_channel_cfg_items 
WHERE f_channel_code=?
AND f_config_item_name=?
`

const UpdateTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt = `
UPDATE trade_channel_cfg_items SET 
 f_id=?
,f_channel_code=?
,f_config_item_name=?
,f_config_item_value=?
,f_config_item_default_value=?
,f_description=?
,f_required=? 
WHERE f_channel_code=?
AND f_config_item_name=?
`

const DeleteTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt = `
DELETE FROM trade_channel_cfg_items 
WHERE f_channel_code=?
AND f_config_item_name=?
`

func scanTradeChannelCfgItem(row *sql.Row) (*schema.TradeChannelCfgItem, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeChannelCfgItem{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ChannelCode = v1.String
	} else {
		v.ChannelCode = ""
	}

	if v2.Valid {
		v.ConfigItemName = v2.String
	} else {
		v.ConfigItemName = ""
	}

	if v3.Valid {
		v.ConfigItemValue = v3.String
	} else {
		v.ConfigItemValue = ""
	}

	if v4.Valid {
		v.ConfigItemDefaultValue = v4.String
	} else {
		v.ConfigItemDefaultValue = ""
	}

	if v5.Valid {
		v.Description = v5.String
	} else {
		v.Description = ""
	}

	if v6.Valid {
		v.Required = int(v6.Int64)
	} else {
		v.Required = 0
	}

	return v, nil
}

func scanTradeChannelCfgItems(rows *sql.Rows) ([]*schema.TradeChannelCfgItem, error) {
	var err error
	var vv []*schema.TradeChannelCfgItem

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeChannelCfgItem{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ChannelCode = v1.String
		} else {
			v.ChannelCode = ""
		}

		if v2.Valid {
			v.ConfigItemName = v2.String
		} else {
			v.ConfigItemName = ""
		}

		if v3.Valid {
			v.ConfigItemValue = v3.String
		} else {
			v.ConfigItemValue = ""
		}

		if v4.Valid {
			v.ConfigItemDefaultValue = v4.String
		} else {
			v.ConfigItemDefaultValue = ""
		}

		if v5.Valid {
			v.Description = v5.String
		} else {
			v.Description = ""
		}

		if v6.Valid {
			v.Required = int(v6.Int64)
		} else {
			v.Required = 0
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeChannelCfgItem(v *schema.TradeChannelCfgItem) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 int

	v0 = v.ID
	v1 = v.ChannelCode
	v2 = v.ConfigItemName
	v3 = v.ConfigItemValue
	v4 = v.ConfigItemDefaultValue
	v5 = v.Description
	v6 = v.Required

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
	}
}

func genericSelectTradeChannelCfgItem(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeChannelCfgItem, error) {
	row := db.QueryRow(query, args...)
	return scanTradeChannelCfgItem(row)
}

func genericSelectTradeChannelCfgItems(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeChannelCfgItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeChannelCfgItems(rows)
}

func InsertTradeChannelCfgItem(db db.SimpleDB, v *schema.TradeChannelCfgItem) error {

	res, err := db.Exec(InsertTradeChannelCfgItemStmt, sliceTradeChannelCfgItem(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeChannelCfgItemById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeChannelCfgItemByIdStmt, args...)
	return err
}

func DeleteTradeChannelCfgItemByChannelCodeAndConfigItemName(db db.SimpleDB, channelCode string, configItemName string) error {
	args := []interface{}{channelCode, configItemName}
	_, err := db.Exec(DeleteTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt, args...)
	return err
}

func UpdateTradeChannelCfgItemById(db db.SimpleDB, v *schema.TradeChannelCfgItem) error {
	args := sliceTradeChannelCfgItem(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeChannelCfgItemByIdStmt, args...)
	return err
}

func UpdateTradeChannelCfgItemByChannelCodeAndConfigItemName(db db.SimpleDB, v *schema.TradeChannelCfgItem) error {
	args := sliceTradeChannelCfgItem(v)
	args = append(args, v.ChannelCode, v.ConfigItemName)
	_, err := db.Exec(UpdateTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt, args...)
	return err
}

func GetTradeChannelCfgItemById(db db.SimpleDB, iD int64) (*schema.TradeChannelCfgItem, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeChannelCfgItem(db, SelectTradeChannelCfgItemByIdStmt, args...)
	return v, err
}

func GetTradeChannelCfgItemByChannelCodeAndConfigItemName(db db.SimpleDB, channelCode string, configItemName string) (*schema.TradeChannelCfgItem, error) {
	args := []interface{}{channelCode, configItemName}
	v, err := genericSelectTradeChannelCfgItem(db, SelectTradeChannelCfgItemByChannelCodeAndConfigItemNameStmt, args...)
	return v, err
}

func FindAllTradeChannelCfgItems(db db.SimpleDB) ([]*schema.TradeChannelCfgItem, error) {
	args := []interface{}{}
	v, err := genericSelectTradeChannelCfgItems(db, SelectTradeChannelCfgItemStmt, args...)
	return v, err
}

func FindAllTradeChannelCfgItemsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeChannelCfgItem, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeChannelCfgItems(db, SelectTradeChannelCfgItemRangeStmt, args...)
	return v, err
}

func CountTradeChannelCfgItem(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeChannelCfgItemCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeChannelCfgItemByChannelCodeAndConfigItemName(db db.SimpleDB, channelCode string, configItemName string) (int, error) {
	var count int
	args := []interface{}{channelCode, configItemName}
	row := db.QueryRow(SelectTradeChannelCfgItemCountByChannelCodeAndConfigItemNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateApplicationTradeChannelStmt = `
CREATE TABLE IF NOT EXISTS application_trade_channels (
 f_id                    BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code           VARCHAR(32)
,f_business_code         VARCHAR(32)
,f_channel_code          VARCHAR(32)
,f_default_trade_account VARCHAR(128)
,f_activated             BOOLEAN
);
`

const InsertApplicationTradeChannelStmt = `
INSERT INTO application_trade_channels (
 f_system_code
,f_business_code
,f_channel_code
,f_default_trade_account
,f_activated
) VALUES (?,?,?,?,?)
`

const SelectApplicationTradeChannelStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_default_trade_account
,f_activated
FROM application_trade_channels 
`

const SelectApplicationTradeChannelRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_default_trade_account
,f_activated
FROM application_trade_channels 
LIMIT ? OFFSET ?
`

const SelectApplicationTradeChannelCountStmt = `
SELECT count(1)
FROM application_trade_channels 
`

const SelectApplicationTradeChannelByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_default_trade_account
,f_activated
FROM application_trade_channels 
WHERE f_id=?
`

const UpdateApplicationTradeChannelByIdStmt = `
UPDATE application_trade_channels SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_channel_code=?
,f_default_trade_account=?
,f_activated=? 
WHERE f_id=?
`

const DeleteApplicationTradeChannelByIdStmt = `
DELETE FROM application_trade_channels 
WHERE f_id=?
`

const CreatePkAtcStmt = `
CREATE UNIQUE INDEX pk_atc ON application_trade_channels (f_system_code,f_business_code,f_channel_code);
`

const SelectApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_default_trade_account
,f_activated
FROM application_trade_channels 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
`

const SelectApplicationTradeChannelCountBySystemCodeAndBusinessCodeAndChannelCodeStmt = `
SELECT count(1)
FROM application_trade_channels 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
`

const UpdateApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt = `
UPDATE application_trade_channels SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_channel_code=?
,f_default_trade_account=?
,f_activated=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
`

const DeleteApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt = `
DELETE FROM application_trade_channels 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
`

func scanApplicationTradeChannel(row *sql.Row) (*schema.ApplicationTradeChannel, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullBool

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

	v := &schema.ApplicationTradeChannel{}

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
		v.ChannelCode = v3.String
	} else {
		v.ChannelCode = ""
	}

	if v4.Valid {
		v.DefaultTradeAccount = v4.String
	} else {
		v.DefaultTradeAccount = ""
	}

	if v5.Valid {
		v.Activated = v5.Bool
	} else {
		v.Activated = false
	}

	return v, nil
}

func scanApplicationTradeChannels(rows *sql.Rows) ([]*schema.ApplicationTradeChannel, error) {
	var err error
	var vv []*schema.ApplicationTradeChannel

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullBool

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

		v := &schema.ApplicationTradeChannel{}

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
			v.ChannelCode = v3.String
		} else {
			v.ChannelCode = ""
		}

		if v4.Valid {
			v.DefaultTradeAccount = v4.String
		} else {
			v.DefaultTradeAccount = ""
		}

		if v5.Valid {
			v.Activated = v5.Bool
		} else {
			v.Activated = false
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceApplicationTradeChannel(v *schema.ApplicationTradeChannel) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 bool

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ChannelCode
	v4 = v.DefaultTradeAccount
	v5 = v.Activated

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
	}
}

func genericSelectApplicationTradeChannel(db db.SimpleDB, query string, args ...interface{}) (*schema.ApplicationTradeChannel, error) {
	row := db.QueryRow(query, args...)
	return scanApplicationTradeChannel(row)
}

func genericSelectApplicationTradeChannels(db db.SimpleDB, query string, args ...interface{}) ([]*schema.ApplicationTradeChannel, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApplicationTradeChannels(rows)
}

func InsertApplicationTradeChannel(db db.SimpleDB, v *schema.ApplicationTradeChannel) error {

	res, err := db.Exec(InsertApplicationTradeChannelStmt, sliceApplicationTradeChannel(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteApplicationTradeChannelById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteApplicationTradeChannelByIdStmt, args...)
	return err
}

func DeleteApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string) error {
	args := []interface{}{systemCode, businessCode, channelCode}
	_, err := db.Exec(DeleteApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt, args...)
	return err
}

func UpdateApplicationTradeChannelById(db db.SimpleDB, v *schema.ApplicationTradeChannel) error {
	args := sliceApplicationTradeChannel(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateApplicationTradeChannelByIdStmt, args...)
	return err
}

func UpdateApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCode(db db.SimpleDB, v *schema.ApplicationTradeChannel) error {
	args := sliceApplicationTradeChannel(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.ChannelCode)
	_, err := db.Exec(UpdateApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt, args...)
	return err
}

func GetApplicationTradeChannelById(db db.SimpleDB, iD int64) (*schema.ApplicationTradeChannel, error) {
	args := []interface{}{iD}
	v, err := genericSelectApplicationTradeChannel(db, SelectApplicationTradeChannelByIdStmt, args...)
	return v, err
}

func GetApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string) (*schema.ApplicationTradeChannel, error) {
	args := []interface{}{systemCode, businessCode, channelCode}
	v, err := genericSelectApplicationTradeChannel(db, SelectApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCodeStmt, args...)
	return v, err
}

func FindAllApplicationTradeChannels(db db.SimpleDB) ([]*schema.ApplicationTradeChannel, error) {
	args := []interface{}{}
	v, err := genericSelectApplicationTradeChannels(db, SelectApplicationTradeChannelStmt, args...)
	return v, err
}

func FindAllApplicationTradeChannelsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.ApplicationTradeChannel, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectApplicationTradeChannels(db, SelectApplicationTradeChannelRangeStmt, args...)
	return v, err
}

func CountApplicationTradeChannel(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectApplicationTradeChannelCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountApplicationTradeChannelBySystemCodeAndBusinessCodeAndChannelCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, channelCode}
	row := db.QueryRow(SelectApplicationTradeChannelCountBySystemCodeAndBusinessCodeAndChannelCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeAlgorithmStmt = `
CREATE TABLE IF NOT EXISTS trade_algorithms (
 f_id                BIGINT PRIMARY KEY AUTO_INCREMENT
,f_channel_code      VARCHAR(32)
,f_algorithm_code    VARCHAR(32)
,f_algorithm_zh_name VARCHAR(128)
,f_algorithm_en_name VARCHAR(64)
,f_description       VARCHAR(512)
);
`

const InsertTradeAlgorithmStmt = `
INSERT INTO trade_algorithms (
 f_channel_code
,f_algorithm_code
,f_algorithm_zh_name
,f_algorithm_en_name
,f_description
) VALUES (?,?,?,?,?)
`

const SelectTradeAlgorithmStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_algorithm_zh_name
,f_algorithm_en_name
,f_description
FROM trade_algorithms 
`

const SelectTradeAlgorithmRangeStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_algorithm_zh_name
,f_algorithm_en_name
,f_description
FROM trade_algorithms 
LIMIT ? OFFSET ?
`

const SelectTradeAlgorithmCountStmt = `
SELECT count(1)
FROM trade_algorithms 
`

const SelectTradeAlgorithmByIdStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_algorithm_zh_name
,f_algorithm_en_name
,f_description
FROM trade_algorithms 
WHERE f_id=?
`

const UpdateTradeAlgorithmByIdStmt = `
UPDATE trade_algorithms SET 
 f_id=?
,f_channel_code=?
,f_algorithm_code=?
,f_algorithm_zh_name=?
,f_algorithm_en_name=?
,f_description=? 
WHERE f_id=?
`

const DeleteTradeAlgorithmByIdStmt = `
DELETE FROM trade_algorithms 
WHERE f_id=?
`

const CreatePkTaStmt = `
CREATE UNIQUE INDEX pk_ta ON trade_algorithms (f_channel_code,f_algorithm_code);
`

const SelectTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_algorithm_zh_name
,f_algorithm_en_name
,f_description
FROM trade_algorithms 
WHERE f_channel_code=?
AND f_algorithm_code=?
`

const SelectTradeAlgorithmCountByChannelCodeAndAlgorithmCodeStmt = `
SELECT count(1)
FROM trade_algorithms 
WHERE f_channel_code=?
AND f_algorithm_code=?
`

const UpdateTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt = `
UPDATE trade_algorithms SET 
 f_id=?
,f_channel_code=?
,f_algorithm_code=?
,f_algorithm_zh_name=?
,f_algorithm_en_name=?
,f_description=? 
WHERE f_channel_code=?
AND f_algorithm_code=?
`

const DeleteTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt = `
DELETE FROM trade_algorithms 
WHERE f_channel_code=?
AND f_algorithm_code=?
`

func scanTradeAlgorithm(row *sql.Row) (*schema.TradeAlgorithm, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

	v := &schema.TradeAlgorithm{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ChannelCode = v1.String
	} else {
		v.ChannelCode = ""
	}

	if v2.Valid {
		v.AlgorithmCode = v2.String
	} else {
		v.AlgorithmCode = ""
	}

	if v3.Valid {
		v.AlgorithmZhName = v3.String
	} else {
		v.AlgorithmZhName = ""
	}

	if v4.Valid {
		v.AlgorithmEnName = v4.String
	} else {
		v.AlgorithmEnName = ""
	}

	if v5.Valid {
		v.Description = v5.String
	} else {
		v.Description = ""
	}

	return v, nil
}

func scanTradeAlgorithms(rows *sql.Rows) ([]*schema.TradeAlgorithm, error) {
	var err error
	var vv []*schema.TradeAlgorithm

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

		v := &schema.TradeAlgorithm{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ChannelCode = v1.String
		} else {
			v.ChannelCode = ""
		}

		if v2.Valid {
			v.AlgorithmCode = v2.String
		} else {
			v.AlgorithmCode = ""
		}

		if v3.Valid {
			v.AlgorithmZhName = v3.String
		} else {
			v.AlgorithmZhName = ""
		}

		if v4.Valid {
			v.AlgorithmEnName = v4.String
		} else {
			v.AlgorithmEnName = ""
		}

		if v5.Valid {
			v.Description = v5.String
		} else {
			v.Description = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeAlgorithm(v *schema.TradeAlgorithm) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string

	v0 = v.ID
	v1 = v.ChannelCode
	v2 = v.AlgorithmCode
	v3 = v.AlgorithmZhName
	v4 = v.AlgorithmEnName
	v5 = v.Description

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
	}
}

func genericSelectTradeAlgorithm(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeAlgorithm, error) {
	row := db.QueryRow(query, args...)
	return scanTradeAlgorithm(row)
}

func genericSelectTradeAlgorithms(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeAlgorithm, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeAlgorithms(rows)
}

func InsertTradeAlgorithm(db db.SimpleDB, v *schema.TradeAlgorithm) error {

	res, err := db.Exec(InsertTradeAlgorithmStmt, sliceTradeAlgorithm(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeAlgorithmById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeAlgorithmByIdStmt, args...)
	return err
}

func DeleteTradeAlgorithmByChannelCodeAndAlgorithmCode(db db.SimpleDB, channelCode string, algorithmCode string) error {
	args := []interface{}{channelCode, algorithmCode}
	_, err := db.Exec(DeleteTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt, args...)
	return err
}

func UpdateTradeAlgorithmById(db db.SimpleDB, v *schema.TradeAlgorithm) error {
	args := sliceTradeAlgorithm(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeAlgorithmByIdStmt, args...)
	return err
}

func UpdateTradeAlgorithmByChannelCodeAndAlgorithmCode(db db.SimpleDB, v *schema.TradeAlgorithm) error {
	args := sliceTradeAlgorithm(v)
	args = append(args, v.ChannelCode, v.AlgorithmCode)
	_, err := db.Exec(UpdateTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt, args...)
	return err
}

func GetTradeAlgorithmById(db db.SimpleDB, iD int64) (*schema.TradeAlgorithm, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeAlgorithm(db, SelectTradeAlgorithmByIdStmt, args...)
	return v, err
}

func GetTradeAlgorithmByChannelCodeAndAlgorithmCode(db db.SimpleDB, channelCode string, algorithmCode string) (*schema.TradeAlgorithm, error) {
	args := []interface{}{channelCode, algorithmCode}
	v, err := genericSelectTradeAlgorithm(db, SelectTradeAlgorithmByChannelCodeAndAlgorithmCodeStmt, args...)
	return v, err
}

func FindAllTradeAlgorithms(db db.SimpleDB) ([]*schema.TradeAlgorithm, error) {
	args := []interface{}{}
	v, err := genericSelectTradeAlgorithms(db, SelectTradeAlgorithmStmt, args...)
	return v, err
}

func FindAllTradeAlgorithmsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeAlgorithm, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeAlgorithms(db, SelectTradeAlgorithmRangeStmt, args...)
	return v, err
}

func CountTradeAlgorithm(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeAlgorithmCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAlgorithmByChannelCodeAndAlgorithmCode(db db.SimpleDB, channelCode string, algorithmCode string) (int, error) {
	var count int
	args := []interface{}{channelCode, algorithmCode}
	row := db.QueryRow(SelectTradeAlgorithmCountByChannelCodeAndAlgorithmCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeAlgorithmAttrItemStmt = `
CREATE TABLE IF NOT EXISTS trade_algorithm_attr_items (
 f_id                    BIGINT PRIMARY KEY AUTO_INCREMENT
,f_channel_code          VARCHAR(32)
,f_algorithm_code        VARCHAR(32)
,f_required              BOOLEAN
,f_attr_name             VARCHAR(32)
,f_attr_zh_name          VARCHAR(128)
,f_attr_value_type       INTEGER
,f_attr_value_len        INTEGER
,f_attr_min_value        DOUBLE
,f_attr_max_value        DOUBLE
,f_attr_value_range_type INTEGER
,f_attr_value_regex      VARCHAR(512)
,f_enum_range            MEDIUMTEXT
);
`

const InsertTradeAlgorithmAttrItemStmt = `
INSERT INTO trade_algorithm_attr_items (
 f_channel_code
,f_algorithm_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectTradeAlgorithmAttrItemStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
FROM trade_algorithm_attr_items 
`

const SelectTradeAlgorithmAttrItemRangeStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
FROM trade_algorithm_attr_items 
LIMIT ? OFFSET ?
`

const SelectTradeAlgorithmAttrItemCountStmt = `
SELECT count(1)
FROM trade_algorithm_attr_items 
`

const SelectTradeAlgorithmAttrItemByIdStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
FROM trade_algorithm_attr_items 
WHERE f_id=?
`

const UpdateTradeAlgorithmAttrItemByIdStmt = `
UPDATE trade_algorithm_attr_items SET 
 f_id=?
,f_channel_code=?
,f_algorithm_code=?
,f_required=?
,f_attr_name=?
,f_attr_zh_name=?
,f_attr_value_type=?
,f_attr_value_len=?
,f_attr_min_value=?
,f_attr_max_value=?
,f_attr_value_range_type=?
,f_attr_value_regex=?
,f_enum_range=? 
WHERE f_id=?
`

const DeleteTradeAlgorithmAttrItemByIdStmt = `
DELETE FROM trade_algorithm_attr_items 
WHERE f_id=?
`

const CreatePkTaaiStmt = `
CREATE UNIQUE INDEX pk_taai ON trade_algorithm_attr_items (f_channel_code,f_algorithm_code,f_attr_name);
`

const SelectTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt = `
SELECT 
 f_id
,f_channel_code
,f_algorithm_code
,f_required
,f_attr_name
,f_attr_zh_name
,f_attr_value_type
,f_attr_value_len
,f_attr_min_value
,f_attr_max_value
,f_attr_value_range_type
,f_attr_value_regex
,f_enum_range
FROM trade_algorithm_attr_items 
WHERE f_channel_code=?
AND f_algorithm_code=?
AND f_attr_name=?
`

const SelectTradeAlgorithmAttrItemCountByChannelCodeAndAlgorithmCodeAndAttrNameStmt = `
SELECT count(1)
FROM trade_algorithm_attr_items 
WHERE f_channel_code=?
AND f_algorithm_code=?
AND f_attr_name=?
`

const UpdateTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt = `
UPDATE trade_algorithm_attr_items SET 
 f_id=?
,f_channel_code=?
,f_algorithm_code=?
,f_required=?
,f_attr_name=?
,f_attr_zh_name=?
,f_attr_value_type=?
,f_attr_value_len=?
,f_attr_min_value=?
,f_attr_max_value=?
,f_attr_value_range_type=?
,f_attr_value_regex=?
,f_enum_range=? 
WHERE f_channel_code=?
AND f_algorithm_code=?
AND f_attr_name=?
`

const DeleteTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt = `
DELETE FROM trade_algorithm_attr_items 
WHERE f_channel_code=?
AND f_algorithm_code=?
AND f_attr_name=?
`

func scanTradeAlgorithmAttrItem(row *sql.Row) (*schema.TradeAlgorithmAttrItem, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullBool
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeAlgorithmAttrItem{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ChannelCode = v1.String
	} else {
		v.ChannelCode = ""
	}

	if v2.Valid {
		v.AlgorithmCode = v2.String
	} else {
		v.AlgorithmCode = ""
	}

	if v3.Valid {
		v.Required = v3.Bool
	} else {
		v.Required = false
	}

	if v4.Valid {
		v.AttrName = v4.String
	} else {
		v.AttrName = ""
	}

	if v5.Valid {
		v.AttrZhName = v5.String
	} else {
		v.AttrZhName = ""
	}

	if v6.Valid {
		v.AttrValueType = int(v6.Int64)
	} else {
		v.AttrValueType = 0
	}

	if v7.Valid {
		v.AttrValueLen = int(v7.Int64)
	} else {
		v.AttrValueLen = 0
	}

	if v8.Valid {
		v.AttrMinValue = v8.Float64
	} else {
		v.AttrMinValue = 0
	}

	if v9.Valid {
		v.AttrMaxValue = v9.Float64
	} else {
		v.AttrMaxValue = 0
	}

	if v10.Valid {
		v.AttrValueRangeType = int(v10.Int64)
	} else {
		v.AttrValueRangeType = 0
	}

	if v11.Valid {
		v.AttrValueRegex = v11.String
	} else {
		v.AttrValueRegex = ""
	}

	if v12.Valid {
		v.EnumRange = v12.String
	} else {
		v.EnumRange = ""
	}

	return v, nil
}

func scanTradeAlgorithmAttrItems(rows *sql.Rows) ([]*schema.TradeAlgorithmAttrItem, error) {
	var err error
	var vv []*schema.TradeAlgorithmAttrItem

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullBool
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullFloat64
	var v9 sql.NullFloat64
	var v10 sql.NullInt64
	var v11 sql.NullString
	var v12 sql.NullString

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeAlgorithmAttrItem{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ChannelCode = v1.String
		} else {
			v.ChannelCode = ""
		}

		if v2.Valid {
			v.AlgorithmCode = v2.String
		} else {
			v.AlgorithmCode = ""
		}

		if v3.Valid {
			v.Required = v3.Bool
		} else {
			v.Required = false
		}

		if v4.Valid {
			v.AttrName = v4.String
		} else {
			v.AttrName = ""
		}

		if v5.Valid {
			v.AttrZhName = v5.String
		} else {
			v.AttrZhName = ""
		}

		if v6.Valid {
			v.AttrValueType = int(v6.Int64)
		} else {
			v.AttrValueType = 0
		}

		if v7.Valid {
			v.AttrValueLen = int(v7.Int64)
		} else {
			v.AttrValueLen = 0
		}

		if v8.Valid {
			v.AttrMinValue = v8.Float64
		} else {
			v.AttrMinValue = 0
		}

		if v9.Valid {
			v.AttrMaxValue = v9.Float64
		} else {
			v.AttrMaxValue = 0
		}

		if v10.Valid {
			v.AttrValueRangeType = int(v10.Int64)
		} else {
			v.AttrValueRangeType = 0
		}

		if v11.Valid {
			v.AttrValueRegex = v11.String
		} else {
			v.AttrValueRegex = ""
		}

		if v12.Valid {
			v.EnumRange = v12.String
		} else {
			v.EnumRange = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeAlgorithmAttrItem(v *schema.TradeAlgorithmAttrItem) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 bool
	var v4 string
	var v5 string
	var v6 int
	var v7 int
	var v8 float64
	var v9 float64
	var v10 int
	var v11 string
	var v12 string

	v0 = v.ID
	v1 = v.ChannelCode
	v2 = v.AlgorithmCode
	v3 = v.Required
	v4 = v.AttrName
	v5 = v.AttrZhName
	v6 = v.AttrValueType
	v7 = v.AttrValueLen
	v8 = v.AttrMinValue
	v9 = v.AttrMaxValue
	v10 = v.AttrValueRangeType
	v11 = v.AttrValueRegex
	v12 = v.EnumRange

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
	}
}

func genericSelectTradeAlgorithmAttrItem(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeAlgorithmAttrItem, error) {
	row := db.QueryRow(query, args...)
	return scanTradeAlgorithmAttrItem(row)
}

func genericSelectTradeAlgorithmAttrItems(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeAlgorithmAttrItem, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeAlgorithmAttrItems(rows)
}

func InsertTradeAlgorithmAttrItem(db db.SimpleDB, v *schema.TradeAlgorithmAttrItem) error {

	res, err := db.Exec(InsertTradeAlgorithmAttrItemStmt, sliceTradeAlgorithmAttrItem(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeAlgorithmAttrItemById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeAlgorithmAttrItemByIdStmt, args...)
	return err
}

func DeleteTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrName(db db.SimpleDB, channelCode string, algorithmCode string, attrName string) error {
	args := []interface{}{channelCode, algorithmCode, attrName}
	_, err := db.Exec(DeleteTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt, args...)
	return err
}

func UpdateTradeAlgorithmAttrItemById(db db.SimpleDB, v *schema.TradeAlgorithmAttrItem) error {
	args := sliceTradeAlgorithmAttrItem(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeAlgorithmAttrItemByIdStmt, args...)
	return err
}

func UpdateTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrName(db db.SimpleDB, v *schema.TradeAlgorithmAttrItem) error {
	args := sliceTradeAlgorithmAttrItem(v)
	args = append(args, v.ChannelCode, v.AlgorithmCode, v.AttrName)
	_, err := db.Exec(UpdateTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt, args...)
	return err
}

func GetTradeAlgorithmAttrItemById(db db.SimpleDB, iD int64) (*schema.TradeAlgorithmAttrItem, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeAlgorithmAttrItem(db, SelectTradeAlgorithmAttrItemByIdStmt, args...)
	return v, err
}

func GetTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrName(db db.SimpleDB, channelCode string, algorithmCode string, attrName string) (*schema.TradeAlgorithmAttrItem, error) {
	args := []interface{}{channelCode, algorithmCode, attrName}
	v, err := genericSelectTradeAlgorithmAttrItem(db, SelectTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrNameStmt, args...)
	return v, err
}

func FindAllTradeAlgorithmAttrItems(db db.SimpleDB) ([]*schema.TradeAlgorithmAttrItem, error) {
	args := []interface{}{}
	v, err := genericSelectTradeAlgorithmAttrItems(db, SelectTradeAlgorithmAttrItemStmt, args...)
	return v, err
}

func FindAllTradeAlgorithmAttrItemsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeAlgorithmAttrItem, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeAlgorithmAttrItems(db, SelectTradeAlgorithmAttrItemRangeStmt, args...)
	return v, err
}

func CountTradeAlgorithmAttrItem(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeAlgorithmAttrItemCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAlgorithmAttrItemByChannelCodeAndAlgorithmCodeAndAttrName(db db.SimpleDB, channelCode string, algorithmCode string, attrName string) (int, error) {
	var count int
	args := []interface{}{channelCode, algorithmCode, attrName}
	row := db.QueryRow(SelectTradeAlgorithmAttrItemCountByChannelCodeAndAlgorithmCodeAndAttrNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeAccountGroupStmt = `
CREATE TABLE IF NOT EXISTS trade_account_groups (
 f_id                BIGINT PRIMARY KEY AUTO_INCREMENT
,f_parent_group_code VARCHAR(128)
,f_group_code        VARCHAR(128)
,f_channel_code      VARCHAR(32)
,f_group_zh_name     VARCHAR(128)
,f_group_en_name     VARCHAR(32)
,f_description       VARCHAR(512)
);
`

const InsertTradeAccountGroupStmt = `
INSERT INTO trade_account_groups (
 f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
) VALUES (?,?,?,?,?,?)
`

const SelectTradeAccountGroupStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
`

const SelectTradeAccountGroupRangeStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
LIMIT ? OFFSET ?
`

const SelectTradeAccountGroupCountStmt = `
SELECT count(1)
FROM trade_account_groups 
`

const SelectTradeAccountGroupByIdStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
WHERE f_id=?
`

const UpdateTradeAccountGroupByIdStmt = `
UPDATE trade_account_groups SET 
 f_id=?
,f_parent_group_code=?
,f_group_code=?
,f_channel_code=?
,f_group_zh_name=?
,f_group_en_name=?
,f_description=? 
WHERE f_id=?
`

const DeleteTradeAccountGroupByIdStmt = `
DELETE FROM trade_account_groups 
WHERE f_id=?
`

const CreateIdxTagPgcStmt = `
CREATE INDEX idx_tag_pgc ON trade_account_groups (f_parent_group_code);
`

const SelectTradeAccountGroupByParentGroupCodeStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
WHERE f_parent_group_code=?
`

const SelectTradeAccountGroupCountByParentGroupCodeStmt = `
SELECT count(1)
FROM trade_account_groups 
WHERE f_parent_group_code=?
`

const SelectTradeAccountGroupRangeByParentGroupCodeStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
WHERE f_parent_group_code=?
LIMIT ? OFFSET ?
`

const DeleteTradeAccountGroupByParentGroupCodeStmt = `
DELETE FROM trade_account_groups 
WHERE f_parent_group_code=?
`

const CreatePkTagStmt = `
CREATE UNIQUE INDEX pk_tag ON trade_account_groups (f_group_code,f_channel_code);
`

const SelectTradeAccountGroupByGroupCodeAndChannelCodeStmt = `
SELECT 
 f_id
,f_parent_group_code
,f_group_code
,f_channel_code
,f_group_zh_name
,f_group_en_name
,f_description
FROM trade_account_groups 
WHERE f_group_code=?
AND f_channel_code=?
`

const SelectTradeAccountGroupCountByGroupCodeAndChannelCodeStmt = `
SELECT count(1)
FROM trade_account_groups 
WHERE f_group_code=?
AND f_channel_code=?
`

const UpdateTradeAccountGroupByGroupCodeAndChannelCodeStmt = `
UPDATE trade_account_groups SET 
 f_id=?
,f_parent_group_code=?
,f_group_code=?
,f_channel_code=?
,f_group_zh_name=?
,f_group_en_name=?
,f_description=? 
WHERE f_group_code=?
AND f_channel_code=?
`

const DeleteTradeAccountGroupByGroupCodeAndChannelCodeStmt = `
DELETE FROM trade_account_groups 
WHERE f_group_code=?
AND f_channel_code=?
`

func scanTradeAccountGroup(row *sql.Row) (*schema.TradeAccountGroup, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeAccountGroup{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.ParentGroupCode = v1.String
	} else {
		v.ParentGroupCode = ""
	}

	if v2.Valid {
		v.GroupCode = v2.String
	} else {
		v.GroupCode = ""
	}

	if v3.Valid {
		v.ChannelCode = v3.String
	} else {
		v.ChannelCode = ""
	}

	if v4.Valid {
		v.GroupZhName = v4.String
	} else {
		v.GroupZhName = ""
	}

	if v5.Valid {
		v.GroupEnName = v5.String
	} else {
		v.GroupEnName = ""
	}

	if v6.Valid {
		v.Description = v6.String
	} else {
		v.Description = ""
	}

	return v, nil
}

func scanTradeAccountGroups(rows *sql.Rows) ([]*schema.TradeAccountGroup, error) {
	var err error
	var vv []*schema.TradeAccountGroup

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeAccountGroup{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.ParentGroupCode = v1.String
		} else {
			v.ParentGroupCode = ""
		}

		if v2.Valid {
			v.GroupCode = v2.String
		} else {
			v.GroupCode = ""
		}

		if v3.Valid {
			v.ChannelCode = v3.String
		} else {
			v.ChannelCode = ""
		}

		if v4.Valid {
			v.GroupZhName = v4.String
		} else {
			v.GroupZhName = ""
		}

		if v5.Valid {
			v.GroupEnName = v5.String
		} else {
			v.GroupEnName = ""
		}

		if v6.Valid {
			v.Description = v6.String
		} else {
			v.Description = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeAccountGroup(v *schema.TradeAccountGroup) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string

	v0 = v.ID
	v1 = v.ParentGroupCode
	v2 = v.GroupCode
	v3 = v.ChannelCode
	v4 = v.GroupZhName
	v5 = v.GroupEnName
	v6 = v.Description

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
	}
}

func genericSelectTradeAccountGroup(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeAccountGroup, error) {
	row := db.QueryRow(query, args...)
	return scanTradeAccountGroup(row)
}

func genericSelectTradeAccountGroups(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeAccountGroup, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeAccountGroups(rows)
}

func InsertTradeAccountGroup(db db.SimpleDB, v *schema.TradeAccountGroup) error {

	res, err := db.Exec(InsertTradeAccountGroupStmt, sliceTradeAccountGroup(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeAccountGroupById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeAccountGroupByIdStmt, args...)
	return err
}

func DeleteTradeAccountGroupByParentGroupCode(db db.SimpleDB, parentGroupCode string) error {
	args := []interface{}{parentGroupCode}
	_, err := db.Exec(DeleteTradeAccountGroupByParentGroupCodeStmt, args...)
	return err
}

func DeleteTradeAccountGroupByGroupCodeAndChannelCode(db db.SimpleDB, groupCode string, channelCode string) error {
	args := []interface{}{groupCode, channelCode}
	_, err := db.Exec(DeleteTradeAccountGroupByGroupCodeAndChannelCodeStmt, args...)
	return err
}

func UpdateTradeAccountGroupById(db db.SimpleDB, v *schema.TradeAccountGroup) error {
	args := sliceTradeAccountGroup(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeAccountGroupByIdStmt, args...)
	return err
}

func UpdateTradeAccountGroupByGroupCodeAndChannelCode(db db.SimpleDB, v *schema.TradeAccountGroup) error {
	args := sliceTradeAccountGroup(v)
	args = append(args, v.GroupCode, v.ChannelCode)
	_, err := db.Exec(UpdateTradeAccountGroupByGroupCodeAndChannelCodeStmt, args...)
	return err
}

func GetTradeAccountGroupById(db db.SimpleDB, iD int64) (*schema.TradeAccountGroup, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeAccountGroup(db, SelectTradeAccountGroupByIdStmt, args...)
	return v, err
}

func GetTradeAccountGroupByGroupCodeAndChannelCode(db db.SimpleDB, groupCode string, channelCode string) (*schema.TradeAccountGroup, error) {
	args := []interface{}{groupCode, channelCode}
	v, err := genericSelectTradeAccountGroup(db, SelectTradeAccountGroupByGroupCodeAndChannelCodeStmt, args...)
	return v, err
}

func FindAllTradeAccountGroups(db db.SimpleDB) ([]*schema.TradeAccountGroup, error) {
	args := []interface{}{}
	v, err := genericSelectTradeAccountGroups(db, SelectTradeAccountGroupStmt, args...)
	return v, err
}

func FindAllTradeAccountGroupsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeAccountGroup, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeAccountGroups(db, SelectTradeAccountGroupRangeStmt, args...)
	return v, err
}

func FindTradeAccountGroupsByParentGroupCode(db db.SimpleDB, parentGroupCode string) ([]*schema.TradeAccountGroup, error) {
	args := []interface{}{parentGroupCode}
	v, err := genericSelectTradeAccountGroups(db, SelectTradeAccountGroupByParentGroupCodeStmt, args...)
	return v, err
}

func FindTradeAccountGroupsByParentGroupCodeInRange(db db.SimpleDB, parentGroupCode string, limit int64, offset int64) ([]*schema.TradeAccountGroup, error) {
	args := []interface{}{parentGroupCode, limit, offset}
	v, err := genericSelectTradeAccountGroups(db, SelectTradeAccountGroupRangeByParentGroupCodeStmt, args...)
	return v, err
}

func CountTradeAccountGroup(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeAccountGroupCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAccountGroupByParentGroupCode(db db.SimpleDB, parentGroupCode string) (int, error) {
	var count int
	args := []interface{}{parentGroupCode}
	row := db.QueryRow(SelectTradeAccountGroupCountByParentGroupCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAccountGroupByGroupCodeAndChannelCode(db db.SimpleDB, groupCode string, channelCode string) (int, error) {
	var count int
	args := []interface{}{groupCode, channelCode}
	row := db.QueryRow(SelectTradeAccountGroupCountByGroupCodeAndChannelCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateTradeAccountStmt = `
CREATE TABLE IF NOT EXISTS trade_accounts (
 f_id              BIGINT PRIMARY KEY AUTO_INCREMENT
,f_account_code    VARCHAR(128)
,f_account_zh_name VARCHAR(512)
,f_account_en_name VARCHAR(512)
,f_channel_code    VARCHAR(32)
,f_group_code      VARCHAR(128)
,f_description     VARCHAR(512)
);
`

const InsertTradeAccountStmt = `
INSERT INTO trade_accounts (
 f_account_code
,f_account_zh_name
,f_account_en_name
,f_channel_code
,f_group_code
,f_description
) VALUES (?,?,?,?,?,?)
`

const SelectTradeAccountStmt = `
SELECT 
 f_id
,f_account_code
,f_account_zh_name
,f_account_en_name
,f_channel_code
,f_group_code
,f_description
FROM trade_accounts 
`

const SelectTradeAccountRangeStmt = `
SELECT 
 f_id
,f_account_code
,f_account_zh_name
,f_account_en_name
,f_channel_code
,f_group_code
,f_description
FROM trade_accounts 
LIMIT ? OFFSET ?
`

const SelectTradeAccountCountStmt = `
SELECT count(1)
FROM trade_accounts 
`

const SelectTradeAccountByIdStmt = `
SELECT 
 f_id
,f_account_code
,f_account_zh_name
,f_account_en_name
,f_channel_code
,f_group_code
,f_description
FROM trade_accounts 
WHERE f_id=?
`

const UpdateTradeAccountByIdStmt = `
UPDATE trade_accounts SET 
 f_id=?
,f_account_code=?
,f_account_zh_name=?
,f_account_en_name=?
,f_channel_code=?
,f_group_code=?
,f_description=? 
WHERE f_id=?
`

const DeleteTradeAccountByIdStmt = `
DELETE FROM trade_accounts 
WHERE f_id=?
`

const CreatePkTacctStmt = `
CREATE UNIQUE INDEX pk_tacct ON trade_accounts (f_account_code,f_channel_code,f_group_code);
`

const SelectTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt = `
SELECT 
 f_id
,f_account_code
,f_account_zh_name
,f_account_en_name
,f_channel_code
,f_group_code
,f_description
FROM trade_accounts 
WHERE f_account_code=?
AND f_channel_code=?
AND f_group_code=?
`

const SelectTradeAccountCountByAccountCodeAndChannelCodeAndGroupCodeStmt = `
SELECT count(1)
FROM trade_accounts 
WHERE f_account_code=?
AND f_channel_code=?
AND f_group_code=?
`

const UpdateTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt = `
UPDATE trade_accounts SET 
 f_id=?
,f_account_code=?
,f_account_zh_name=?
,f_account_en_name=?
,f_channel_code=?
,f_group_code=?
,f_description=? 
WHERE f_account_code=?
AND f_channel_code=?
AND f_group_code=?
`

const DeleteTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt = `
DELETE FROM trade_accounts 
WHERE f_account_code=?
AND f_channel_code=?
AND f_group_code=?
`

func scanTradeAccount(row *sql.Row) (*schema.TradeAccount, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString

	err := row.Scan(
		&v0,
		&v1,
		&v2,
		&v3,
		&v4,
		&v5,
		&v6,
	)
	if err != nil {
		return nil, err
	}

	v := &schema.TradeAccount{}

	if v0.Valid {
		v.ID = v0.Int64
	} else {
		v.ID = 0
	}

	if v1.Valid {
		v.AccountCode = v1.String
	} else {
		v.AccountCode = ""
	}

	if v2.Valid {
		v.AccountZhName = v2.String
	} else {
		v.AccountZhName = ""
	}

	if v3.Valid {
		v.AccountEnName = v3.String
	} else {
		v.AccountEnName = ""
	}

	if v4.Valid {
		v.ChannelCode = v4.String
	} else {
		v.ChannelCode = ""
	}

	if v5.Valid {
		v.GroupCode = v5.String
	} else {
		v.GroupCode = ""
	}

	if v6.Valid {
		v.Description = v6.String
	} else {
		v.Description = ""
	}

	return v, nil
}

func scanTradeAccounts(rows *sql.Rows) ([]*schema.TradeAccount, error) {
	var err error
	var vv []*schema.TradeAccount

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString

	for rows.Next() {
		err = rows.Scan(
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
			&v6,
		)
		if err != nil {
			return vv, err
		}

		v := &schema.TradeAccount{}

		if v0.Valid {
			v.ID = v0.Int64
		} else {
			v.ID = 0
		}

		if v1.Valid {
			v.AccountCode = v1.String
		} else {
			v.AccountCode = ""
		}

		if v2.Valid {
			v.AccountZhName = v2.String
		} else {
			v.AccountZhName = ""
		}

		if v3.Valid {
			v.AccountEnName = v3.String
		} else {
			v.AccountEnName = ""
		}

		if v4.Valid {
			v.ChannelCode = v4.String
		} else {
			v.ChannelCode = ""
		}

		if v5.Valid {
			v.GroupCode = v5.String
		} else {
			v.GroupCode = ""
		}

		if v6.Valid {
			v.Description = v6.String
		} else {
			v.Description = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceTradeAccount(v *schema.TradeAccount) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string

	v0 = v.ID
	v1 = v.AccountCode
	v2 = v.AccountZhName
	v3 = v.AccountEnName
	v4 = v.ChannelCode
	v5 = v.GroupCode
	v6 = v.Description

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
		v6,
	}
}

func genericSelectTradeAccount(db db.SimpleDB, query string, args ...interface{}) (*schema.TradeAccount, error) {
	row := db.QueryRow(query, args...)
	return scanTradeAccount(row)
}

func genericSelectTradeAccounts(db db.SimpleDB, query string, args ...interface{}) ([]*schema.TradeAccount, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTradeAccounts(rows)
}

func InsertTradeAccount(db db.SimpleDB, v *schema.TradeAccount) error {

	res, err := db.Exec(InsertTradeAccountStmt, sliceTradeAccount(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteTradeAccountById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteTradeAccountByIdStmt, args...)
	return err
}

func DeleteTradeAccountByAccountCodeAndChannelCodeAndGroupCode(db db.SimpleDB, accountCode string, channelCode string, groupCode string) error {
	args := []interface{}{accountCode, channelCode, groupCode}
	_, err := db.Exec(DeleteTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt, args...)
	return err
}

func UpdateTradeAccountById(db db.SimpleDB, v *schema.TradeAccount) error {
	args := sliceTradeAccount(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateTradeAccountByIdStmt, args...)
	return err
}

func UpdateTradeAccountByAccountCodeAndChannelCodeAndGroupCode(db db.SimpleDB, v *schema.TradeAccount) error {
	args := sliceTradeAccount(v)
	args = append(args, v.AccountCode, v.ChannelCode, v.GroupCode)
	_, err := db.Exec(UpdateTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt, args...)
	return err
}

func GetTradeAccountById(db db.SimpleDB, iD int64) (*schema.TradeAccount, error) {
	args := []interface{}{iD}
	v, err := genericSelectTradeAccount(db, SelectTradeAccountByIdStmt, args...)
	return v, err
}

func GetTradeAccountByAccountCodeAndChannelCodeAndGroupCode(db db.SimpleDB, accountCode string, channelCode string, groupCode string) (*schema.TradeAccount, error) {
	args := []interface{}{accountCode, channelCode, groupCode}
	v, err := genericSelectTradeAccount(db, SelectTradeAccountByAccountCodeAndChannelCodeAndGroupCodeStmt, args...)
	return v, err
}

func FindAllTradeAccounts(db db.SimpleDB) ([]*schema.TradeAccount, error) {
	args := []interface{}{}
	v, err := genericSelectTradeAccounts(db, SelectTradeAccountStmt, args...)
	return v, err
}

func FindAllTradeAccountsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.TradeAccount, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectTradeAccounts(db, SelectTradeAccountRangeStmt, args...)
	return v, err
}

func CountTradeAccount(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectTradeAccountCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountTradeAccountByAccountCodeAndChannelCodeAndGroupCode(db db.SimpleDB, accountCode string, channelCode string, groupCode string) (int, error) {
	var count int
	args := []interface{}{accountCode, channelCode, groupCode}
	row := db.QueryRow(SelectTradeAccountCountByAccountCodeAndChannelCodeAndGroupCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateApplicationTradeAccountStmt = `
CREATE TABLE IF NOT EXISTS application_trade_accounts (
 f_id            BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code   VARCHAR(32)
,f_business_code VARCHAR(32)
,f_channel_code  VARCHAR(32)
,f_account_code  VARCHAR(128)
,f_user_id       VARCHAR(256)
);
`

const InsertApplicationTradeAccountStmt = `
INSERT INTO application_trade_accounts (
 f_system_code
,f_business_code
,f_channel_code
,f_account_code
,f_user_id
) VALUES (?,?,?,?,?)
`

const SelectApplicationTradeAccountStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_account_code
,f_user_id
FROM application_trade_accounts 
`

const SelectApplicationTradeAccountRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_account_code
,f_user_id
FROM application_trade_accounts 
LIMIT ? OFFSET ?
`

const SelectApplicationTradeAccountCountStmt = `
SELECT count(1)
FROM application_trade_accounts 
`

const SelectApplicationTradeAccountByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_account_code
,f_user_id
FROM application_trade_accounts 
WHERE f_id=?
`

const UpdateApplicationTradeAccountByIdStmt = `
UPDATE application_trade_accounts SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_channel_code=?
,f_account_code=?
,f_user_id=? 
WHERE f_id=?
`

const DeleteApplicationTradeAccountByIdStmt = `
DELETE FROM application_trade_accounts 
WHERE f_id=?
`

const CreatePkAtaStmt = `
CREATE UNIQUE INDEX pk_ata ON application_trade_accounts (f_system_code,f_business_code,f_channel_code,f_account_code);
`

const SelectApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_channel_code
,f_account_code
,f_user_id
FROM application_trade_accounts 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
AND f_account_code=?
`

const SelectApplicationTradeAccountCountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt = `
SELECT count(1)
FROM application_trade_accounts 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
AND f_account_code=?
`

const UpdateApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt = `
UPDATE application_trade_accounts SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_channel_code=?
,f_account_code=?
,f_user_id=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
AND f_account_code=?
`

const DeleteApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt = `
DELETE FROM application_trade_accounts 
WHERE f_system_code=?
AND f_business_code=?
AND f_channel_code=?
AND f_account_code=?
`

func scanApplicationTradeAccount(row *sql.Row) (*schema.ApplicationTradeAccount, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

	v := &schema.ApplicationTradeAccount{}

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
		v.ChannelCode = v3.String
	} else {
		v.ChannelCode = ""
	}

	if v4.Valid {
		v.AccountCode = v4.String
	} else {
		v.AccountCode = ""
	}

	if v5.Valid {
		v.UserId = v5.String
	} else {
		v.UserId = ""
	}

	return v, nil
}

func scanApplicationTradeAccounts(rows *sql.Rows) ([]*schema.ApplicationTradeAccount, error) {
	var err error
	var vv []*schema.ApplicationTradeAccount

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
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

		v := &schema.ApplicationTradeAccount{}

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
			v.ChannelCode = v3.String
		} else {
			v.ChannelCode = ""
		}

		if v4.Valid {
			v.AccountCode = v4.String
		} else {
			v.AccountCode = ""
		}

		if v5.Valid {
			v.UserId = v5.String
		} else {
			v.UserId = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceApplicationTradeAccount(v *schema.ApplicationTradeAccount) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ChannelCode
	v4 = v.AccountCode
	v5 = v.UserId

	return []interface{}{
		v0,
		v1,
		v2,
		v3,
		v4,
		v5,
	}
}

func genericSelectApplicationTradeAccount(db db.SimpleDB, query string, args ...interface{}) (*schema.ApplicationTradeAccount, error) {
	row := db.QueryRow(query, args...)
	return scanApplicationTradeAccount(row)
}

func genericSelectApplicationTradeAccounts(db db.SimpleDB, query string, args ...interface{}) ([]*schema.ApplicationTradeAccount, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanApplicationTradeAccounts(rows)
}

func InsertApplicationTradeAccount(db db.SimpleDB, v *schema.ApplicationTradeAccount) error {

	res, err := db.Exec(InsertApplicationTradeAccountStmt, sliceApplicationTradeAccount(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteApplicationTradeAccountById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteApplicationTradeAccountByIdStmt, args...)
	return err
}

func DeleteApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string, accountCode string) error {
	args := []interface{}{systemCode, businessCode, channelCode, accountCode}
	_, err := db.Exec(DeleteApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt, args...)
	return err
}

func UpdateApplicationTradeAccountById(db db.SimpleDB, v *schema.ApplicationTradeAccount) error {
	args := sliceApplicationTradeAccount(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateApplicationTradeAccountByIdStmt, args...)
	return err
}

func UpdateApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCode(db db.SimpleDB, v *schema.ApplicationTradeAccount) error {
	args := sliceApplicationTradeAccount(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.ChannelCode, v.AccountCode)
	_, err := db.Exec(UpdateApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt, args...)
	return err
}

func GetApplicationTradeAccountById(db db.SimpleDB, iD int64) (*schema.ApplicationTradeAccount, error) {
	args := []interface{}{iD}
	v, err := genericSelectApplicationTradeAccount(db, SelectApplicationTradeAccountByIdStmt, args...)
	return v, err
}

func GetApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string, accountCode string) (*schema.ApplicationTradeAccount, error) {
	args := []interface{}{systemCode, businessCode, channelCode, accountCode}
	v, err := genericSelectApplicationTradeAccount(db, SelectApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt, args...)
	return v, err
}

func FindAllApplicationTradeAccounts(db db.SimpleDB) ([]*schema.ApplicationTradeAccount, error) {
	args := []interface{}{}
	v, err := genericSelectApplicationTradeAccounts(db, SelectApplicationTradeAccountStmt, args...)
	return v, err
}

func FindAllApplicationTradeAccountsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.ApplicationTradeAccount, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectApplicationTradeAccounts(db, SelectApplicationTradeAccountRangeStmt, args...)
	return v, err
}

func CountApplicationTradeAccount(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectApplicationTradeAccountCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountApplicationTradeAccountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCode(db db.SimpleDB, systemCode string, businessCode string, channelCode string, accountCode string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, channelCode, accountCode}
	row := db.QueryRow(SelectApplicationTradeAccountCountBySystemCodeAndBusinessCodeAndChannelCodeAndAccountCodeStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateDataArchivingLogStmt = `
CREATE TABLE IF NOT EXISTS data_archiving_logs (
 f_id                     BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code            VARCHAR(32)
,f_business_code          VARCHAR(32)
,f_archiving_date         VARCHAR(8)
,f_task_name              VARCHAR(32)
,f_first_archiving_time   BIGINT
,f_current_archiving_time BIGINT
,f_exec_count             INTEGER
,f_complete               BOOLEAN
,f_complete_phase         INTEGER
);
`

const InsertDataArchivingLogStmt = `
INSERT INTO data_archiving_logs (
 f_system_code
,f_business_code
,f_archiving_date
,f_task_name
,f_first_archiving_time
,f_current_archiving_time
,f_exec_count
,f_complete
,f_complete_phase
) VALUES (?,?,?,?,?,?,?,?,?)
`

const SelectDataArchivingLogStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_archiving_date
,f_task_name
,f_first_archiving_time
,f_current_archiving_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_archiving_logs 
`

const SelectDataArchivingLogRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_archiving_date
,f_task_name
,f_first_archiving_time
,f_current_archiving_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_archiving_logs 
LIMIT ? OFFSET ?
`

const SelectDataArchivingLogCountStmt = `
SELECT count(1)
FROM data_archiving_logs 
`

const SelectDataArchivingLogByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_archiving_date
,f_task_name
,f_first_archiving_time
,f_current_archiving_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_archiving_logs 
WHERE f_id=?
`

const UpdateDataArchivingLogByIdStmt = `
UPDATE data_archiving_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_archiving_date=?
,f_task_name=?
,f_first_archiving_time=?
,f_current_archiving_time=?
,f_exec_count=?
,f_complete=?
,f_complete_phase=? 
WHERE f_id=?
`

const DeleteDataArchivingLogByIdStmt = `
DELETE FROM data_archiving_logs 
WHERE f_id=?
`

const CreatePkDalStmt = `
CREATE UNIQUE INDEX pk_dal ON data_archiving_logs (f_system_code,f_business_code,f_archiving_date,f_task_name);
`

const SelectDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_archiving_date
,f_task_name
,f_first_archiving_time
,f_current_archiving_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_archiving_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_archiving_date=?
AND f_task_name=?
`

const SelectDataArchivingLogCountBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt = `
SELECT count(1)
FROM data_archiving_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_archiving_date=?
AND f_task_name=?
`

const UpdateDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt = `
UPDATE data_archiving_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_archiving_date=?
,f_task_name=?
,f_first_archiving_time=?
,f_current_archiving_time=?
,f_exec_count=?
,f_complete=?
,f_complete_phase=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_archiving_date=?
AND f_task_name=?
`

const DeleteDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt = `
DELETE FROM data_archiving_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_archiving_date=?
AND f_task_name=?
`

func scanDataArchivingLog(row *sql.Row) (*schema.DataArchivingLog, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullBool
	var v9 sql.NullInt64

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
	)
	if err != nil {
		return nil, err
	}

	v := &schema.DataArchivingLog{}

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
		v.ArchivingDate = v3.String
	} else {
		v.ArchivingDate = ""
	}

	if v4.Valid {
		v.TaskName = v4.String
	} else {
		v.TaskName = ""
	}

	if v5.Valid {
		v.FirstArchivingTime = v5.Int64
	} else {
		v.FirstArchivingTime = 0
	}

	if v6.Valid {
		v.CurrentArchivingTime = v6.Int64
	} else {
		v.CurrentArchivingTime = 0
	}

	if v7.Valid {
		v.ExecCount = int(v7.Int64)
	} else {
		v.ExecCount = 0
	}

	if v8.Valid {
		v.Complete = v8.Bool
	} else {
		v.Complete = false
	}

	if v9.Valid {
		v.CompletePhase = int(v9.Int64)
	} else {
		v.CompletePhase = 0
	}

	return v, nil
}

func scanDataArchivingLogs(rows *sql.Rows) ([]*schema.DataArchivingLog, error) {
	var err error
	var vv []*schema.DataArchivingLog

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullInt64
	var v7 sql.NullInt64
	var v8 sql.NullBool
	var v9 sql.NullInt64

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
		)
		if err != nil {
			return vv, err
		}

		v := &schema.DataArchivingLog{}

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
			v.ArchivingDate = v3.String
		} else {
			v.ArchivingDate = ""
		}

		if v4.Valid {
			v.TaskName = v4.String
		} else {
			v.TaskName = ""
		}

		if v5.Valid {
			v.FirstArchivingTime = v5.Int64
		} else {
			v.FirstArchivingTime = 0
		}

		if v6.Valid {
			v.CurrentArchivingTime = v6.Int64
		} else {
			v.CurrentArchivingTime = 0
		}

		if v7.Valid {
			v.ExecCount = int(v7.Int64)
		} else {
			v.ExecCount = 0
		}

		if v8.Valid {
			v.Complete = v8.Bool
		} else {
			v.Complete = false
		}

		if v9.Valid {
			v.CompletePhase = int(v9.Int64)
		} else {
			v.CompletePhase = 0
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceDataArchivingLog(v *schema.DataArchivingLog) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 int64
	var v6 int64
	var v7 int
	var v8 bool
	var v9 int

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.ArchivingDate
	v4 = v.TaskName
	v5 = v.FirstArchivingTime
	v6 = v.CurrentArchivingTime
	v7 = v.ExecCount
	v8 = v.Complete
	v9 = v.CompletePhase

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
	}
}

func genericSelectDataArchivingLog(db db.SimpleDB, query string, args ...interface{}) (*schema.DataArchivingLog, error) {
	row := db.QueryRow(query, args...)
	return scanDataArchivingLog(row)
}

func genericSelectDataArchivingLogs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.DataArchivingLog, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDataArchivingLogs(rows)
}

func InsertDataArchivingLog(db db.SimpleDB, v *schema.DataArchivingLog) error {

	res, err := db.Exec(InsertDataArchivingLogStmt, sliceDataArchivingLog(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteDataArchivingLogById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteDataArchivingLogByIdStmt, args...)
	return err
}

func DeleteDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, archivingDate string, taskName string) error {
	args := []interface{}{systemCode, businessCode, archivingDate, taskName}
	_, err := db.Exec(DeleteDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt, args...)
	return err
}

func UpdateDataArchivingLogById(db db.SimpleDB, v *schema.DataArchivingLog) error {
	args := sliceDataArchivingLog(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateDataArchivingLogByIdStmt, args...)
	return err
}

func UpdateDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(db db.SimpleDB, v *schema.DataArchivingLog) error {
	args := sliceDataArchivingLog(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.ArchivingDate, v.TaskName)
	_, err := db.Exec(UpdateDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt, args...)
	return err
}

func GetDataArchivingLogById(db db.SimpleDB, iD int64) (*schema.DataArchivingLog, error) {
	args := []interface{}{iD}
	v, err := genericSelectDataArchivingLog(db, SelectDataArchivingLogByIdStmt, args...)
	return v, err
}

func GetDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, archivingDate string, taskName string) (*schema.DataArchivingLog, error) {
	args := []interface{}{systemCode, businessCode, archivingDate, taskName}
	v, err := genericSelectDataArchivingLog(db, SelectDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt, args...)
	return v, err
}

func FindAllDataArchivingLogs(db db.SimpleDB) ([]*schema.DataArchivingLog, error) {
	args := []interface{}{}
	v, err := genericSelectDataArchivingLogs(db, SelectDataArchivingLogStmt, args...)
	return v, err
}

func FindAllDataArchivingLogsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.DataArchivingLog, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectDataArchivingLogs(db, SelectDataArchivingLogRangeStmt, args...)
	return v, err
}

func CountDataArchivingLog(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectDataArchivingLogCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountDataArchivingLogBySystemCodeAndBusinessCodeAndArchivingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, archivingDate string, taskName string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, archivingDate, taskName}
	row := db.QueryRow(SelectDataArchivingLogCountBySystemCodeAndBusinessCodeAndArchivingDateAndTaskNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateDataPurgingLogStmt = `
CREATE TABLE IF NOT EXISTS data_purging_logs (
 f_id                               BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code                      VARCHAR(32)
,f_business_code                    VARCHAR(32)
,f_purging_date                     VARCHAR(8)
,f_task_name                        VARCHAR(32)
,f_group_trade_order_purging        LONGTEXT
,f_trade_order_purging              LONGTEXT
,f_trade_action_latest_resp_purging LONGTEXT
,f_trade_action_resp_purge          LONGTEXT
,f_first_purging_time               BIGINT
,f_current_purging_time             BIGINT
,f_exec_count                       INTEGER
,f_complete                         BOOLEAN
,f_complete_phase                   INTEGER
);
`

const InsertDataPurgingLogStmt = `
INSERT INTO data_purging_logs (
 f_system_code
,f_business_code
,f_purging_date
,f_task_name
,f_group_trade_order_purging
,f_trade_order_purging
,f_trade_action_latest_resp_purging
,f_trade_action_resp_purge
,f_first_purging_time
,f_current_purging_time
,f_exec_count
,f_complete
,f_complete_phase
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectDataPurgingLogStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_purging_date
,f_task_name
,f_group_trade_order_purging
,f_trade_order_purging
,f_trade_action_latest_resp_purging
,f_trade_action_resp_purge
,f_first_purging_time
,f_current_purging_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_purging_logs 
`

const SelectDataPurgingLogRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_purging_date
,f_task_name
,f_group_trade_order_purging
,f_trade_order_purging
,f_trade_action_latest_resp_purging
,f_trade_action_resp_purge
,f_first_purging_time
,f_current_purging_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_purging_logs 
LIMIT ? OFFSET ?
`

const SelectDataPurgingLogCountStmt = `
SELECT count(1)
FROM data_purging_logs 
`

const SelectDataPurgingLogByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_purging_date
,f_task_name
,f_group_trade_order_purging
,f_trade_order_purging
,f_trade_action_latest_resp_purging
,f_trade_action_resp_purge
,f_first_purging_time
,f_current_purging_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_purging_logs 
WHERE f_id=?
`

const UpdateDataPurgingLogByIdStmt = `
UPDATE data_purging_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_purging_date=?
,f_task_name=?
,f_group_trade_order_purging=?
,f_trade_order_purging=?
,f_trade_action_latest_resp_purging=?
,f_trade_action_resp_purge=?
,f_first_purging_time=?
,f_current_purging_time=?
,f_exec_count=?
,f_complete=?
,f_complete_phase=? 
WHERE f_id=?
`

const DeleteDataPurgingLogByIdStmt = `
DELETE FROM data_purging_logs 
WHERE f_id=?
`

const CreatePkDclStmt = `
CREATE UNIQUE INDEX pk_dcl ON data_purging_logs (f_system_code,f_business_code,f_purging_date,f_task_name);
`

const SelectDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_purging_date
,f_task_name
,f_group_trade_order_purging
,f_trade_order_purging
,f_trade_action_latest_resp_purging
,f_trade_action_resp_purge
,f_first_purging_time
,f_current_purging_time
,f_exec_count
,f_complete
,f_complete_phase
FROM data_purging_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_purging_date=?
AND f_task_name=?
`

const SelectDataPurgingLogCountBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt = `
SELECT count(1)
FROM data_purging_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_purging_date=?
AND f_task_name=?
`

const UpdateDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt = `
UPDATE data_purging_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_purging_date=?
,f_task_name=?
,f_group_trade_order_purging=?
,f_trade_order_purging=?
,f_trade_action_latest_resp_purging=?
,f_trade_action_resp_purge=?
,f_first_purging_time=?
,f_current_purging_time=?
,f_exec_count=?
,f_complete=?
,f_complete_phase=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_purging_date=?
AND f_task_name=?
`

const DeleteDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt = `
DELETE FROM data_purging_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_purging_date=?
AND f_task_name=?
`

func scanDataPurgingLog(row *sql.Row) (*schema.DataPurgingLog, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullBool
	var v13 sql.NullInt64

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

	v := &schema.DataPurgingLog{}

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
		v.PurgingDate = v3.String
	} else {
		v.PurgingDate = ""
	}

	if v4.Valid {
		v.TaskName = v4.String
	} else {
		v.TaskName = ""
	}

	if v5.Valid {
		v.GroupTradeOrderPurging = v5.String
	} else {
		v.GroupTradeOrderPurging = ""
	}

	if v6.Valid {
		v.TradeOrderPurging = v6.String
	} else {
		v.TradeOrderPurging = ""
	}

	if v7.Valid {
		v.TradeActionLatestRespPurging = v7.String
	} else {
		v.TradeActionLatestRespPurging = ""
	}

	if v8.Valid {
		v.TradeActionRespPurge = v8.String
	} else {
		v.TradeActionRespPurge = ""
	}

	if v9.Valid {
		v.FirstPurgingTime = v9.Int64
	} else {
		v.FirstPurgingTime = 0
	}

	if v10.Valid {
		v.CurrentPurgingTime = v10.Int64
	} else {
		v.CurrentPurgingTime = 0
	}

	if v11.Valid {
		v.ExecCount = int(v11.Int64)
	} else {
		v.ExecCount = 0
	}

	if v12.Valid {
		v.Complete = v12.Bool
	} else {
		v.Complete = false
	}

	if v13.Valid {
		v.CompletePhase = int(v13.Int64)
	} else {
		v.CompletePhase = 0
	}

	return v, nil
}

func scanDataPurgingLogs(rows *sql.Rows) ([]*schema.DataPurgingLog, error) {
	var err error
	var vv []*schema.DataPurgingLog

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullString
	var v6 sql.NullString
	var v7 sql.NullString
	var v8 sql.NullString
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullBool
	var v13 sql.NullInt64

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

		v := &schema.DataPurgingLog{}

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
			v.PurgingDate = v3.String
		} else {
			v.PurgingDate = ""
		}

		if v4.Valid {
			v.TaskName = v4.String
		} else {
			v.TaskName = ""
		}

		if v5.Valid {
			v.GroupTradeOrderPurging = v5.String
		} else {
			v.GroupTradeOrderPurging = ""
		}

		if v6.Valid {
			v.TradeOrderPurging = v6.String
		} else {
			v.TradeOrderPurging = ""
		}

		if v7.Valid {
			v.TradeActionLatestRespPurging = v7.String
		} else {
			v.TradeActionLatestRespPurging = ""
		}

		if v8.Valid {
			v.TradeActionRespPurge = v8.String
		} else {
			v.TradeActionRespPurge = ""
		}

		if v9.Valid {
			v.FirstPurgingTime = v9.Int64
		} else {
			v.FirstPurgingTime = 0
		}

		if v10.Valid {
			v.CurrentPurgingTime = v10.Int64
		} else {
			v.CurrentPurgingTime = 0
		}

		if v11.Valid {
			v.ExecCount = int(v11.Int64)
		} else {
			v.ExecCount = 0
		}

		if v12.Valid {
			v.Complete = v12.Bool
		} else {
			v.Complete = false
		}

		if v13.Valid {
			v.CompletePhase = int(v13.Int64)
		} else {
			v.CompletePhase = 0
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceDataPurgingLog(v *schema.DataPurgingLog) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 string
	var v6 string
	var v7 string
	var v8 string
	var v9 int64
	var v10 int64
	var v11 int
	var v12 bool
	var v13 int

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.PurgingDate
	v4 = v.TaskName
	v5 = v.GroupTradeOrderPurging
	v6 = v.TradeOrderPurging
	v7 = v.TradeActionLatestRespPurging
	v8 = v.TradeActionRespPurge
	v9 = v.FirstPurgingTime
	v10 = v.CurrentPurgingTime
	v11 = v.ExecCount
	v12 = v.Complete
	v13 = v.CompletePhase

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

func genericSelectDataPurgingLog(db db.SimpleDB, query string, args ...interface{}) (*schema.DataPurgingLog, error) {
	row := db.QueryRow(query, args...)
	return scanDataPurgingLog(row)
}

func genericSelectDataPurgingLogs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.DataPurgingLog, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDataPurgingLogs(rows)
}

func InsertDataPurgingLog(db db.SimpleDB, v *schema.DataPurgingLog) error {

	res, err := db.Exec(InsertDataPurgingLogStmt, sliceDataPurgingLog(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteDataPurgingLogById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteDataPurgingLogByIdStmt, args...)
	return err
}

func DeleteDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, purgingDate string, taskName string) error {
	args := []interface{}{systemCode, businessCode, purgingDate, taskName}
	_, err := db.Exec(DeleteDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt, args...)
	return err
}

func UpdateDataPurgingLogById(db db.SimpleDB, v *schema.DataPurgingLog) error {
	args := sliceDataPurgingLog(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateDataPurgingLogByIdStmt, args...)
	return err
}

func UpdateDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(db db.SimpleDB, v *schema.DataPurgingLog) error {
	args := sliceDataPurgingLog(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.PurgingDate, v.TaskName)
	_, err := db.Exec(UpdateDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt, args...)
	return err
}

func GetDataPurgingLogById(db db.SimpleDB, iD int64) (*schema.DataPurgingLog, error) {
	args := []interface{}{iD}
	v, err := genericSelectDataPurgingLog(db, SelectDataPurgingLogByIdStmt, args...)
	return v, err
}

func GetDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, purgingDate string, taskName string) (*schema.DataPurgingLog, error) {
	args := []interface{}{systemCode, businessCode, purgingDate, taskName}
	v, err := genericSelectDataPurgingLog(db, SelectDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt, args...)
	return v, err
}

func FindAllDataPurgingLogs(db db.SimpleDB) ([]*schema.DataPurgingLog, error) {
	args := []interface{}{}
	v, err := genericSelectDataPurgingLogs(db, SelectDataPurgingLogStmt, args...)
	return v, err
}

func FindAllDataPurgingLogsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.DataPurgingLog, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectDataPurgingLogs(db, SelectDataPurgingLogRangeStmt, args...)
	return v, err
}

func CountDataPurgingLog(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectDataPurgingLogCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountDataPurgingLogBySystemCodeAndBusinessCodeAndPurgingDateAndTaskName(db db.SimpleDB, systemCode string, businessCode string, purgingDate string, taskName string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, purgingDate, taskName}
	row := db.QueryRow(SelectDataPurgingLogCountBySystemCodeAndBusinessCodeAndPurgingDateAndTaskNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

const CreateDataSyncLogStmt = `
CREATE TABLE IF NOT EXISTS data_sync_logs (
 f_id                 BIGINT PRIMARY KEY AUTO_INCREMENT
,f_system_code        VARCHAR(32)
,f_business_code      VARCHAR(32)
,f_sync_date          VARCHAR(8)
,f_table_name         VARCHAR(32)
,f_sync_type          INTEGER
,f_sync_params        MEDIUMTEXT
,f_report_time        BIGINT
,f_first_sync_time    BIGINT
,f_current_sync_time  BIGINT
,f_complete_sync_time BIGINT
,f_exec_count         INTEGER
,f_sync_phase         INTEGER
,f_fail_log           MEDIUMTEXT
);
`

const InsertDataSyncLogStmt = `
INSERT INTO data_sync_logs (
 f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)
`

const SelectDataSyncLogStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
`

const SelectDataSyncLogRangeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
LIMIT ? OFFSET ?
`

const SelectDataSyncLogCountStmt = `
SELECT count(1)
FROM data_sync_logs 
`

const SelectDataSyncLogByIdStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_id=?
`

const UpdateDataSyncLogByIdStmt = `
UPDATE data_sync_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_sync_date=?
,f_table_name=?
,f_sync_type=?
,f_sync_params=?
,f_report_time=?
,f_first_sync_time=?
,f_current_sync_time=?
,f_complete_sync_time=?
,f_exec_count=?
,f_sync_phase=?
,f_fail_log=? 
WHERE f_id=?
`

const DeleteDataSyncLogByIdStmt = `
DELETE FROM data_sync_logs 
WHERE f_id=?
`

const CreateDslSbdStmt = `
CREATE INDEX dsl_sbd ON data_sync_logs (f_system_code,f_business_code,f_sync_date,f_table_name);
`

const SelectDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
`

const SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt = `
SELECT count(1)
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
`

const SelectDataSyncLogRangeBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
LIMIT ? OFFSET ?
`

const DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt = `
DELETE FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
`

const CreateDlsSbtsStmt = `
CREATE INDEX dls_sbts ON data_sync_logs (f_system_code,f_business_code,f_table_name,f_sync_phase);
`

const SelectDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_table_name=?
AND f_sync_phase=?
`

const SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt = `
SELECT count(1)
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_table_name=?
AND f_sync_phase=?
`

const SelectDataSyncLogRangeBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_table_name=?
AND f_sync_phase=?
LIMIT ? OFFSET ?
`

const DeleteDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt = `
DELETE FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_table_name=?
AND f_sync_phase=?
`

const CreatePkDslStmt = `
CREATE UNIQUE INDEX pk_dsl ON data_sync_logs (f_system_code,f_business_code,f_sync_date,f_table_name,f_report_time);
`

const SelectDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt = `
SELECT 
 f_id
,f_system_code
,f_business_code
,f_sync_date
,f_table_name
,f_sync_type
,f_sync_params
,f_report_time
,f_first_sync_time
,f_current_sync_time
,f_complete_sync_time
,f_exec_count
,f_sync_phase
,f_fail_log
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
AND f_report_time=?
`

const SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt = `
SELECT count(1)
FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
AND f_report_time=?
`

const UpdateDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt = `
UPDATE data_sync_logs SET 
 f_id=?
,f_system_code=?
,f_business_code=?
,f_sync_date=?
,f_table_name=?
,f_sync_type=?
,f_sync_params=?
,f_report_time=?
,f_first_sync_time=?
,f_current_sync_time=?
,f_complete_sync_time=?
,f_exec_count=?
,f_sync_phase=?
,f_fail_log=? 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
AND f_report_time=?
`

const DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt = `
DELETE FROM data_sync_logs 
WHERE f_system_code=?
AND f_business_code=?
AND f_sync_date=?
AND f_table_name=?
AND f_report_time=?
`

func scanDataSyncLog(row *sql.Row) (*schema.DataSyncLog, error) {
	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullInt64
	var v8 sql.NullInt64
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
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

	v := &schema.DataSyncLog{}

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
		v.SyncDate = v3.String
	} else {
		v.SyncDate = ""
	}

	if v4.Valid {
		v.TableName = v4.String
	} else {
		v.TableName = ""
	}

	if v5.Valid {
		v.SyncType = int(v5.Int64)
	} else {
		v.SyncType = 0
	}

	if v6.Valid {
		v.SyncParams = v6.String
	} else {
		v.SyncParams = ""
	}

	if v7.Valid {
		v.ReportTime = v7.Int64
	} else {
		v.ReportTime = 0
	}

	if v8.Valid {
		v.FirstSyncTime = v8.Int64
	} else {
		v.FirstSyncTime = 0
	}

	if v9.Valid {
		v.CurrentSyncTime = v9.Int64
	} else {
		v.CurrentSyncTime = 0
	}

	if v10.Valid {
		v.CompleteSyncTime = v10.Int64
	} else {
		v.CompleteSyncTime = 0
	}

	if v11.Valid {
		v.ExecCount = int(v11.Int64)
	} else {
		v.ExecCount = 0
	}

	if v12.Valid {
		v.SyncPhase = int(v12.Int64)
	} else {
		v.SyncPhase = 0
	}

	if v13.Valid {
		v.FailLog = v13.String
	} else {
		v.FailLog = ""
	}

	return v, nil
}

func scanDataSyncLogs(rows *sql.Rows) ([]*schema.DataSyncLog, error) {
	var err error
	var vv []*schema.DataSyncLog

	var v0 sql.NullInt64
	var v1 sql.NullString
	var v2 sql.NullString
	var v3 sql.NullString
	var v4 sql.NullString
	var v5 sql.NullInt64
	var v6 sql.NullString
	var v7 sql.NullInt64
	var v8 sql.NullInt64
	var v9 sql.NullInt64
	var v10 sql.NullInt64
	var v11 sql.NullInt64
	var v12 sql.NullInt64
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

		v := &schema.DataSyncLog{}

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
			v.SyncDate = v3.String
		} else {
			v.SyncDate = ""
		}

		if v4.Valid {
			v.TableName = v4.String
		} else {
			v.TableName = ""
		}

		if v5.Valid {
			v.SyncType = int(v5.Int64)
		} else {
			v.SyncType = 0
		}

		if v6.Valid {
			v.SyncParams = v6.String
		} else {
			v.SyncParams = ""
		}

		if v7.Valid {
			v.ReportTime = v7.Int64
		} else {
			v.ReportTime = 0
		}

		if v8.Valid {
			v.FirstSyncTime = v8.Int64
		} else {
			v.FirstSyncTime = 0
		}

		if v9.Valid {
			v.CurrentSyncTime = v9.Int64
		} else {
			v.CurrentSyncTime = 0
		}

		if v10.Valid {
			v.CompleteSyncTime = v10.Int64
		} else {
			v.CompleteSyncTime = 0
		}

		if v11.Valid {
			v.ExecCount = int(v11.Int64)
		} else {
			v.ExecCount = 0
		}

		if v12.Valid {
			v.SyncPhase = int(v12.Int64)
		} else {
			v.SyncPhase = 0
		}

		if v13.Valid {
			v.FailLog = v13.String
		} else {
			v.FailLog = ""
		}

		vv = append(vv, v)
	}
	return vv, rows.Err()
}

func sliceDataSyncLog(v *schema.DataSyncLog) []interface{} {
	var v0 int64
	var v1 string
	var v2 string
	var v3 string
	var v4 string
	var v5 int
	var v6 string
	var v7 int64
	var v8 int64
	var v9 int64
	var v10 int64
	var v11 int
	var v12 int
	var v13 string

	v0 = v.ID
	v1 = v.SystemCode
	v2 = v.BusinessCode
	v3 = v.SyncDate
	v4 = v.TableName
	v5 = v.SyncType
	v6 = v.SyncParams
	v7 = v.ReportTime
	v8 = v.FirstSyncTime
	v9 = v.CurrentSyncTime
	v10 = v.CompleteSyncTime
	v11 = v.ExecCount
	v12 = v.SyncPhase
	v13 = v.FailLog

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

func genericSelectDataSyncLog(db db.SimpleDB, query string, args ...interface{}) (*schema.DataSyncLog, error) {
	row := db.QueryRow(query, args...)
	return scanDataSyncLog(row)
}

func genericSelectDataSyncLogs(db db.SimpleDB, query string, args ...interface{}) ([]*schema.DataSyncLog, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDataSyncLogs(rows)
}

func InsertDataSyncLog(db db.SimpleDB, v *schema.DataSyncLog) error {

	res, err := db.Exec(InsertDataSyncLogStmt, sliceDataSyncLog(v)[1:]...)
	if err != nil {
		return err
	}

	v.ID, err = res.LastInsertId()
	return err
}

func DeleteDataSyncLogById(db db.SimpleDB, iD int64) error {
	args := []interface{}{iD}
	_, err := db.Exec(DeleteDataSyncLogByIdStmt, args...)
	return err
}

func DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableName(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string) error {
	args := []interface{}{systemCode, businessCode, syncDate, tableName}
	_, err := db.Exec(DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt, args...)
	return err
}

func DeleteDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhase(db db.SimpleDB, systemCode string, businessCode string, tableName string, syncPhase int) error {
	args := []interface{}{systemCode, businessCode, tableName, syncPhase}
	_, err := db.Exec(DeleteDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt, args...)
	return err
}

func DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTime(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string, reportTime int64) error {
	args := []interface{}{systemCode, businessCode, syncDate, tableName, reportTime}
	_, err := db.Exec(DeleteDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt, args...)
	return err
}

func UpdateDataSyncLogById(db db.SimpleDB, v *schema.DataSyncLog) error {
	args := sliceDataSyncLog(v)
	args = append(args, v.ID)
	_, err := db.Exec(UpdateDataSyncLogByIdStmt, args...)
	return err
}

func UpdateDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTime(db db.SimpleDB, v *schema.DataSyncLog) error {
	args := sliceDataSyncLog(v)
	args = append(args, v.SystemCode, v.BusinessCode, v.SyncDate, v.TableName, v.ReportTime)
	_, err := db.Exec(UpdateDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt, args...)
	return err
}

func GetDataSyncLogById(db db.SimpleDB, iD int64) (*schema.DataSyncLog, error) {
	args := []interface{}{iD}
	v, err := genericSelectDataSyncLog(db, SelectDataSyncLogByIdStmt, args...)
	return v, err
}

func GetDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTime(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string, reportTime int64) (*schema.DataSyncLog, error) {
	args := []interface{}{systemCode, businessCode, syncDate, tableName, reportTime}
	v, err := genericSelectDataSyncLog(db, SelectDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt, args...)
	return v, err
}

func FindAllDataSyncLogs(db db.SimpleDB) ([]*schema.DataSyncLog, error) {
	args := []interface{}{}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogStmt, args...)
	return v, err
}

func FindAllDataSyncLogsInRange(db db.SimpleDB, limit int64, offset int64) ([]*schema.DataSyncLog, error) {
	args := []interface{}{limit, offset}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogRangeStmt, args...)
	return v, err
}

func FindDataSyncLogsBySystemCodeAndBusinessCodeAndSyncDateAndTableName(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string) ([]*schema.DataSyncLog, error) {
	args := []interface{}{systemCode, businessCode, syncDate, tableName}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt, args...)
	return v, err
}

func FindDataSyncLogsBySystemCodeAndBusinessCodeAndSyncDateAndTableNameInRange(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string, limit int64, offset int64) ([]*schema.DataSyncLog, error) {
	args := []interface{}{systemCode, businessCode, syncDate, tableName, limit, offset}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogRangeBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt, args...)
	return v, err
}

func FindDataSyncLogsBySystemCodeAndBusinessCodeAndTableNameAndSyncPhase(db db.SimpleDB, systemCode string, businessCode string, tableName string, syncPhase int) ([]*schema.DataSyncLog, error) {
	args := []interface{}{systemCode, businessCode, tableName, syncPhase}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt, args...)
	return v, err
}

func FindDataSyncLogsBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseInRange(db db.SimpleDB, systemCode string, businessCode string, tableName string, syncPhase int, limit int64, offset int64) ([]*schema.DataSyncLog, error) {
	args := []interface{}{systemCode, businessCode, tableName, syncPhase, limit, offset}
	v, err := genericSelectDataSyncLogs(db, SelectDataSyncLogRangeBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt, args...)
	return v, err
}

func CountDataSyncLog(db db.SimpleDB) (int, error) {
	var count int
	row := db.QueryRow(SelectDataSyncLogCountStmt)
	err := row.Scan(&count)
	return count, err
}

func CountDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableName(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, syncDate, tableName}
	row := db.QueryRow(SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndSyncDateAndTableNameStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountDataSyncLogBySystemCodeAndBusinessCodeAndTableNameAndSyncPhase(db db.SimpleDB, systemCode string, businessCode string, tableName string, syncPhase int) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, tableName, syncPhase}
	row := db.QueryRow(SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndTableNameAndSyncPhaseStmt, args...)
	err := row.Scan(&count)
	return count, err
}

func CountDataSyncLogBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTime(db db.SimpleDB, systemCode string, businessCode string, syncDate string, tableName string, reportTime int64) (int, error) {
	var count int
	args := []interface{}{systemCode, businessCode, syncDate, tableName, reportTime}
	row := db.QueryRow(SelectDataSyncLogCountBySystemCodeAndBusinessCodeAndSyncDateAndTableNameAndReportTimeStmt, args...)
	err := row.Scan(&count)
	return count, err
}
