package dbutil

import (
	"log"
	"strings"
)

type StructuralQuery struct {
	SelectFields    []string          `json:"select_fields"`
	Table           string            `json:"table"`
	FieldConditions []*FieldCondition `json:"field_conditions"`
	GroupByFields   []string          `json:"group_by_fields"`
	SortFields      []string          `json:"sort_fields"`
	SortType        int               `json:"sort_type"` // 0 - 升序，1 - 降序
	Limit           int               `json:"limit"`
	Offset          int               `json:"offset"`
}

func (q *StructuralQuery) GenerateQueryStatement() (query string, args []interface{}, err error) {
	query = "SELECT "
	if len(q.SelectFields) > 0 {
		query += strings.Join(q.SelectFields, ",")
	} else {
		query += "*"
	}
	query += " FROM " + q.Table
	if len(q.FieldConditions) > 0 {
		var whereClause string
		whereClause, args = PrepareQueryStatmentRecursion("AND", q.FieldConditions)
		query += " WHERE " + whereClause
	}
	if len(q.GroupByFields) > 0 {
		query += " GROUP BY " + strings.Join(q.GroupByFields, ",")
	}
	if len(q.SortFields) > 0 {
		query += " ORDER BY " + strings.Join(q.SortFields, ",")
	}
	if q.SortType == 1 {
		query += " DESC"
	}
	query += " LIMIT ? OFFSET ?"
	args = append(args, q.Limit, q.Offset)

	log.Printf("query sql:%s\n", query)
	return
}

func (q *StructuralQuery) GenerateQueryCountStatement(queryStatement string, queryStatementArgs []interface{}) (query string, args []interface{}) {
	/*begin := strings.Index(queryStatement, "FROM")
	end := strings.Index(queryStatement, "LIMIT")
	queryStatement = queryStatement[begin:end]

	idx := strings.Index(queryStatement, "ORDER BY")
	if idx > 0 {
		queryStatement = queryStatement[:idx]
	}

	groupBy := "GROUP BY"
	idx = strings.Index(queryStatement, groupBy)
	if idx > 0 {
		distinct := queryStatement[idx+len(groupBy):]
		queryStatement = queryStatement[:idx]
		query = "SELECT COUNT(DISTINCT "+distinct+") as total_records " + queryStatement
	} else {
		query = "SELECT COUNT(1) as total_records " + queryStatement
	}
	args = queryStatementArgs[:len(queryStatementArgs)-2]
	log.Printf("count query:%s\nargs len:%d\n", query, len(args))
	return*/
	begin := strings.Index(queryStatement, "FROM")
	end := strings.Index(queryStatement, "LIMIT")
	queryStatement = queryStatement[begin:end]

	idx := strings.Index(queryStatement, "ORDER BY")
	if idx > 0 {
		queryStatement = queryStatement[:idx]
	}

	groupBy := "GROUP BY"
	idx = strings.Index(queryStatement, groupBy)
	if idx > 0 {
		query = "SELECT COUNT(1) as total_records FROM ( SELECT 1 " + queryStatement + ")subquery"
	} else {
		query = "SELECT COUNT(1) as total_records " + queryStatement
	}
	args = queryStatementArgs[:len(queryStatementArgs)-2]
	log.Printf("count query:%s,\n args:%v,\n queryStatementArgs:%v\n", query, args, queryStatementArgs)
	return
}

func PrepareQueryStatmentRecursion(operator string, fieldConditions []*FieldCondition) (whereClause string, args []interface{}) {

	if len(fieldConditions) > 0 {
		whereClause += "("
	}

	for i, fieldCondition := range fieldConditions {
		if i > 0 {
			whereClause += ") " + operator + " ("
		}

		if len(fieldCondition.SubFieldConditions) == 0 {
			switch fieldCondition.ValueType {
			case FieldConditionValueTypeItem:
				if fieldCondition.IsNot {
					whereClause += fieldPrefix + fieldCondition.Field + " != ?"
				} else {
					if fieldCondition.FieldType == FieldTypeString {
						whereClause += fieldPrefix + fieldCondition.Field + " like ?"
					} else {
						whereClause += fieldPrefix + fieldCondition.Field + "=?"
					}
				}
				args = append(args, fieldCondition.Value)
			case FieldConditionValueTypeRange:
				values := fieldCondition.getValues(nil)
				whereClause += fieldPrefix + fieldCondition.Field + ">=?"
				if len(values) > 1 {
					whereClause += " AND " + fieldPrefix + fieldCondition.Field + "<?"
				}
				args = append(args, values...)
			case FieldConditionValueTypeSet:
				values := fieldCondition.getValues(nil)
				if fieldCondition.IsNot {
					whereClause += fieldPrefix + fieldCondition.Field + " NOT IN ("
				} else {
					whereClause += fieldPrefix + fieldCondition.Field + " IN ("
				}
				for i := range values {
					if i > 0 {
						whereClause += ","
					}
					whereClause += "?"
				}
				whereClause += ")"
				args = append(args, values...)
			}
		} else {
			whereClause2, args2 := PrepareQueryStatmentRecursion(fieldCondition.Operator, fieldCondition.SubFieldConditions)
			whereClause += whereClause2
			args = append(args, args2...)
		}
	}

	if len(fieldConditions) > 0 {
		whereClause += ")"
	}

	return whereClause, args
}
