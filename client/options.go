package client

import (
	nethttp "net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL    = "https://api.imgwire.dev"
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 2
	defaultBackoff    = 500 * time.Millisecond
	defaultUserAgent  = "imgwire-go/0.1.0"
)

type Options struct {
	BaseURL       string
	EnvironmentID string
	Timeout       time.Duration
	MaxRetries    int
	Backoff       time.Duration
	UserAgent     string
	HTTPClient    *nethttp.Client
}

type Option func(*Options)

func defaultOptions() Options {
	return Options{
		BaseURL:    defaultBaseURL,
		Timeout:    defaultTimeout,
		MaxRetries: defaultMaxRetries,
		Backoff:    defaultBackoff,
		UserAgent:  defaultUserAgent,
	}
}

func (o Options) normalized() Options {
	if o.BaseURL == "" {
		o.BaseURL = defaultBaseURL
	}
	o.BaseURL = strings.TrimRight(o.BaseURL, "/")
	if o.Timeout <= 0 {
		o.Timeout = defaultTimeout
	}
	if o.MaxRetries < 0 {
		o.MaxRetries = 0
	}
	if o.Backoff <= 0 {
		o.Backoff = defaultBackoff
	}
	if o.UserAgent == "" {
		o.UserAgent = defaultUserAgent
	}
	return o
}

func WithBaseURL(baseURL string) Option {
	return func(options *Options) {
		options.BaseURL = baseURL
	}
}

func WithEnvironmentID(environmentID string) Option {
	return func(options *Options) {
		options.EnvironmentID = environmentID
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

func WithMaxRetries(maxRetries int) Option {
	return func(options *Options) {
		options.MaxRetries = maxRetries
	}
}

func WithBackoff(backoff time.Duration) Option {
	return func(options *Options) {
		options.Backoff = backoff
	}
}

func WithUserAgent(userAgent string) Option {
	return func(options *Options) {
		options.UserAgent = userAgent
	}
}

func WithHTTPClient(httpClient *nethttp.Client) Option {
	return func(options *Options) {
		options.HTTPClient = httpClient
	}
}
