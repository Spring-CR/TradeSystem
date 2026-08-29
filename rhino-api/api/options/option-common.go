package options

import "rhino-common/utils/dbutil"

type CommonQueryOption struct {
	FieldConditions []*dbutil.FieldCondition `json:"field_conditions"`
	Limit           int                      `json:"limit"`
	Offset          int                      `json:"offset"`
}

type PagingRecord struct {
	Total int `json:"total"`
	Data interface{} `json:"data"`
}