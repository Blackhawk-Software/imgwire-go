package tests

import (
	"bytes"
	"io"
	"net/http"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return fn(request)
}

func jsonResponse(statusCode int, body string, headers map[string]string) *http.Response {
	header := make(http.Header, len(headers)+1)
	header.Set("Content-Type", "application/json")
	for key, value := range headers {
		header.Set(key, value)
	}

	return &http.Response{
		StatusCode: statusCode,
		Header:     header,
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
