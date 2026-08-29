package ficc

import (
	"database/sql"
)

type CyptGroupInfo struct {
	GroupID   int64
	GroupName string
}

// GetCyptGroupMap 返回一个 map，key 为交易对手 ID，value 为所属组的信息。
// 若交易对手无归属组或组记录不存在，则不包含在 map 中。
func GetCyptGroupMap(db *sql.DB) (map[int64]*CyptGroupInfo, error) {
	query := `
		SELECT t.KEY_CTPTY_ID, g.KEY_CTPTY_ID, g.CTPTY_SHORT_NAME
		FROM Counterparties t
		JOIN Counterparties g ON t.CTPTY_GROUP_ID = g.KEY_CTPTY_ID
		WHERE (t.GROUP_FLAG IS NULL OR t.GROUP_FLAG != 'Y')
		  AND t.CTPTY_GROUP_ID IS NOT NULL
		  AND g.GROUP_FLAG = 'Y'
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]*CyptGroupInfo)
	for rows.Next() {
		var counterpartyID int64
		var groupID int64
		var groupName string
		if err := rows.Scan(&counterpartyID, &groupID, &groupName); err != nil {
			return nil, err
		}
		result[counterpartyID] = &CyptGroupInfo{
			GroupID:   groupID,
			GroupName: groupName,
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

type PlanInfo struct {
	CounterpartyID    int64
	PlanID            int64
	PlanCode          string
	UltraContractID   int64
	UltraContractCode string
	BusinessType      string
}

func GetPlanInfoMap(db *sql.DB) (map[int64]*PlanInfo, error) {
	query := `
	select KEY_CTPTY_ID as CounterpartyID, 
	KEY_PLAN_ID as PlanID, PLAN_CODE as PlanCode, 
	ULTRA_CONTRACT_ID as UltraContractID, ULTRA_CONTRACT_CODE as UltraContractCode, 
	BUSINESS_TYPE as BusinessType from business_plans where KEY_CTPTY_ID is not null and PLAN_CODE  is not null	
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[int64]*PlanInfo)

	for rows.Next() {
		var (
			CounterpartyID    sql.NullInt64
			PlanID            sql.NullInt64
			PlanCode          sql.NullString
			UltraContractID   sql.NullInt64
			UltraContractCode sql.NullString
			BusinessType      sql.NullString
		)
		if err := rows.Scan(&CounterpartyID, &PlanID, &PlanCode, &UltraContractID, &UltraContractCode, &BusinessType); err != nil {
			return nil, err
		}
		result[CounterpartyID.Int64] = &PlanInfo{
			CounterpartyID:    CounterpartyID.Int64,
			PlanID:            PlanID.Int64,
			PlanCode:          PlanCode.String,
			UltraContractID:   UltraContractID.Int64,
			UltraContractCode: UltraContractCode.String,
			BusinessType:      BusinessType.String,
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
