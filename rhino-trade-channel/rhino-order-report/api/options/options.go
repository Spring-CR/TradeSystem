package options

type PagingRecord struct {
	Total int `json:"total"`
	Data interface{} `json:"data"`
}