package imgwirehttp

import (
	"bytes"
	"io"
	nethttp "net/http"
	"time"
)

type RetryTransport struct {
	Base       nethttp.RoundTripper
	MaxRetries int
	Backoff    time.Duration
}

func (t *RetryTransport) RoundTrip(request *nethttp.Request) (*nethttp.Response, error) {
	base := t.Base
	if base == nil {
		base = nethttp.DefaultTransport
	}

	backoff := t.Backoff
	if backoff <= 0 {
		backoff = 500 * time.Millisecond
	}

	var lastResponse *nethttp.Response
	var lastError error

	for attempt := 0; attempt <= t.MaxRetries; attempt++ {
		cloned, err := cloneRequest(request)
		if err != nil {
			return nil, err
		}

		response, err := base.RoundTrip(cloned)
		if !shouldRetry(response, err, attempt, t.MaxRetries) {
			return response, err
		}

		if response != nil && response.Body != nil {
			_ = response.Body.Close()
		}
		lastResponse = response
		lastError = err
		time.Sleep(backoff * time.Duration(attempt+1))
	}

	return lastResponse, lastError
}

func cloneRequest(request *nethttp.Request) (*nethttp.Request, error) {
	cloned := request.Clone(request.Context())
	if request.Body == nil {
		return cloned, nil
	}

	if request.GetBody != nil {
		body, err := request.GetBody()
		if err != nil {
			return nil, err
		}
		cloned.Body = body
		return cloned, nil
	}

	if request.ContentLength == 0 {
		cloned.Body = nil
		return cloned, nil
	}

	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, err
	}
	_ = request.Body.Close()
	request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyBytes)), nil
	}
	cloned.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	return cloned, nil
}

func shouldRetry(response *nethttp.Response, err error, attempt int, maxRetries int) bool {
	if attempt >= maxRetries {
		return false
	}
	if err != nil {
		return true
	}
	if response == nil {
		return false
	}
	return response.StatusCode == nethttp.StatusTooManyRequests || response.StatusCode >= 500
}
