package pagination

import (
	"net/http"
	"strconv"
	"strings"
)

// url?page_size=12
func parseIntQueryParam(r *http.Request, key string, defaultValue int) (int, error) {
	qs := r.URL.Query()

	if queryParam := qs.Get(key); queryParam != "" {
		val, err := strconv.Atoi(queryParam)

		if err != nil {
			return defaultValue, err
		}

		return val, nil
	}

	return defaultValue, nil
}

func parseStrQueryParam(r *http.Request, key string, defaultValue string) (string, error) {
	qs := r.URL.Query()

	if queryParam := qs.Get(key); queryParam != "" {
		return queryParam, nil
	}

	return defaultValue, nil
}

func parseCSVQueryParam(r *http.Request, key string) ([]string, error) {
	qs := r.URL.Query()

	queryParam := qs.Get(key)
	if queryParam == "" {
		return []string{}, nil
	}

	tags := strings.Split(queryParam, ",")
	return tags, nil
}
