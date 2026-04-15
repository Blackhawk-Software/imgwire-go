package imgwirehttp

import (
	nethttp "net/http"
	"time"
)

type Options struct {
	Timeout    time.Duration
	MaxRetries int
	Backoff    time.Duration
	Transport  nethttp.RoundTripper
}

func NewClient(options Options) *nethttp.Client {
	transport := options.Transport
	if transport == nil {
		transport = nethttp.DefaultTransport
	}

	return &nethttp.Client{
		Timeout: options.Timeout,
		Transport: &RetryTransport{
			Base:       transport,
			MaxRetries: options.MaxRetries,
			Backoff:    options.Backoff,
		},
	}
}
