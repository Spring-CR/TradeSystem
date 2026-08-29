package dbutil

import "github.com/linchunquan/sqlgen/db"

// 泛型分页查询模板函数
func PaginateQueryAll[T any](
	db db.SimpleDB,
	pageSize int64,
	queryFn func(db db.SimpleDB, limit, offset int64) ([]*T, error),
) ([]*T, error) {

	var results []*T
	var offset int64

	for {
		chunk, err := queryFn(db, pageSize, offset)
		if err != nil && IsDbRecordEmptyError(err) {
			err = nil
		}
		if err != nil {
			return nil, err
		}

		results = append(results, chunk...)

		// 判断是否最后一页
		if int64(len(chunk)) < pageSize {
			break
		}

		offset += pageSize
	}

	return results, nil
}
