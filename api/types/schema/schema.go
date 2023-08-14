package schema

// pagination
type PaginationParam struct {
	PageNo   int `json:"pageNo"`
	PageSize int `json:"pageSize"`
}

type PaginationResult struct {
	Total    int64 `json:"total"`
	PageNo   int   `json:"pageNo"`
	PageSize int   `json:"pageSize"`
}

func (a PaginationParam) GetPageNo() int {
	return a.PageNo
}

func (a PaginationParam) GetPageSize() int {
	pageSize := a.PageSize
	if a.PageSize <= 0 {
		pageSize = 10
	}
	return pageSize
}
