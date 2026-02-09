package pagination

import (
	"net/http"
	"slices"
	"strings"
)

type SortOrder string

const (
	SortAsc                 = SortOrder("ASC")
	SortDesc                = SortOrder("DESC")
	DefaultSortField string = "-id"
)

type Filter struct {
	Search       string   `json:"search" validate:"max=100"` // includes title and content
	Tags         []string `json:"tags" validate:"max=5"`
	Sort         string   `json:"sort" validate:"oneof=id created_at -id -created_at"`
	SortSafeList []string
	// Since  string   `json:"since"`
}

func ParseFilterParams(r *http.Request) (Filter, error) {
	search, err := parseStrQueryParam(r, "search", "")
	if err != nil {
		return Filter{}, err
	}

	tags, err := parseCSVQueryParam(r, "tags")
	if err != nil {
		return Filter{}, err
	}

	sort, err := parseStrQueryParam(r, "sort", DefaultSortField)
	if err != nil {
		return Filter{}, err
	}

	return Filter{
		Search: search,
		Tags:   tags,
		Sort:   sort,
	}, nil
}

func (f Filter) SortColumn() string {
	if slices.Contains(f.SortSafeList, f.Sort) {
		return strings.TrimPrefix(f.Sort, "-")
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (p Filter) SortDirection() SortOrder {
	if strings.HasPrefix(p.Sort, "-") {
		return SortDesc
	}

	return SortAsc
}
