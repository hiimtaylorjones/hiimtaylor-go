package models

type Pagination struct {
	CurrentPage int
	TotalPages 	int
	TotalCount 	int
	PerPage 		int
	HasPrev			bool
	HasNext			bool
}

func NewPagination(page, perPage, totalCount int) Pagination {
	totalPages := totalCount / perPage
	if totalCount%perPage != 0 {
		totalPages++
	}

	if page < 1 {
		page = 1
	}

	return Pagination{
		CurrentPage: 	page,
		TotalPages:		totalPages,
		TotalCount:		totalCount,
		PerPage:			perPage,
		HasPrev:			page > 1,
		HasNext:			page < totalPages,	
	}
}