package options

type QueryResult struct {
	Total      int         `json:"total"`
	Data       interface{} `json:"data"`
	DisplayLen int         `json:"displayLen"`
}
