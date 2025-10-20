package domain

type ListFilter struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
    Query  string `json:"query"`
}

func NewListFilter(limit int, offset int, query string) ListFilter {
	return ListFilter{
		Limit:  int64(limit),
		Offset: int64(offset),
        Query:  query,
	}
}

type Pagination struct {
	HasPrevious bool
	HasNext     bool
	NextOffset  int64
	PrevOffset  int64
	Limit       int64
	TotalCount  int64
}

func NewPagination(filter ListFilter, totalCount int64) Pagination {
	nextOffset := filter.Offset + filter.Limit
	prevOffset := filter.Offset - filter.Limit
	hasMore := nextOffset < totalCount

	return Pagination{
		HasPrevious: filter.Offset > 0,
		HasNext:     hasMore,
		NextOffset:  nextOffset,
		PrevOffset:  prevOffset,
		Limit:       filter.Limit,
		TotalCount:  totalCount,
	}
}
