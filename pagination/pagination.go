package pagination

import (
	"net/http"
	"strconv"
	"strings"
)

type Metadata struct {
	TotalCount int
	Page       int
	Limit      int
	PrevPage   *int
	NextPage   *int
}

type Page[T any] struct {
	Data       []T
	Pagination Metadata
}

func ParseHeaders(headers http.Header) Metadata {
	return Metadata{
		TotalCount: parseHeaderInt(headers, "X-Total-Count"),
		Page:       parseHeaderInt(headers, "X-Page"),
		Limit:      parseHeaderInt(headers, "X-Limit"),
		PrevPage:   parseOptionalHeaderInt(headers, "X-Prev-Page"),
		NextPage:   parseOptionalHeaderInt(headers, "X-Next-Page"),
	}
}

func parseHeaderInt(headers http.Header, key string) int {
	value := strings.TrimSpace(headers.Get(key))
	if value == "" {
		return 0
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}

func parseOptionalHeaderInt(headers http.Header, key string) *int {
	value := strings.TrimSpace(headers.Get(key))
	if value == "" || strings.EqualFold(value, "null") {
		return nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &parsed
}
