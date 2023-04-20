package pkg

import gorest "github.com/FTChinese/go-rest"

type AsyncResult[T any] struct {
	Value T
	Err   error
}

type PagedList[T any] struct {
	Total int64 `json:"total" db:"row_count"`
	gorest.Pagination
	Data []T `json:"data"`
}
