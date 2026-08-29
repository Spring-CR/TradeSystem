package options

type GenericQueryResult[T any] struct {
	Total   int              `json:"total"`
	Data    []T              `json:"data"`
	Code    int              `json:"Code"`
	Message string           `json:"Message"`
}
