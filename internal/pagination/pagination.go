package pagination

import (
	"net/http"
	"strconv"
)

const (
	DefaultPageSize    int = 10
	DefaultCurrentPage int = 1
)

type Pagination struct {
	CurrentPage  int `json:"current_page" validate:"min=1"`
	PageSize     int `json:"page_size" validate:"min=1,max=20"`
	FirstPage    int
	LastPage     int
	TotalRecords int
}

func (p *Pagination) Paginate() {
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}

	if p.PageSize <= 0 {
		p.PageSize = DefaultPageSize
	}

	// empty return
	if p.TotalRecords == 0 {
		p.CurrentPage = 0
		p.FirstPage = 0
		p.LastPage = 0
		return
	}

	p.FirstPage = 1
	p.LastPage = (p.TotalRecords + p.PageSize - 1) / p.PageSize

	if p.CurrentPage > p.LastPage {
		p.CurrentPage = p.LastPage
	}
}

func (p *Pagination) Offset() int {
	if p.CurrentPage <= 1 {
		return 0
	}

	return (p.CurrentPage - 1) * p.PageSize
}

func (p *Pagination) WriteHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Current-Page", strconv.Itoa(p.CurrentPage))
	w.Header().Set("X-Page-Size", strconv.Itoa(p.PageSize))
	w.Header().Set("X-First-Page", strconv.Itoa(p.FirstPage))
	w.Header().Set("X-Last-Page", strconv.Itoa(p.LastPage))
	w.Header().Set("X-Total-Records", strconv.Itoa(p.TotalRecords))
}

func ParsePaginationParams(r *http.Request) (Pagination, error) {
	page, err := parseIntQueryParam(r, "page", DefaultCurrentPage)
	if err != nil {
		return Pagination{}, err
	}

	pageSize, err := parseIntQueryParam(r, "page_size", DefaultPageSize)
	if err != nil {
		return Pagination{}, err
	}

	return Pagination{
		CurrentPage: page,
		PageSize:    pageSize,
	}, nil
}
