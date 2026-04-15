package client

import (
	nethttp "net/http"

	generated "github.com/Blackhawk-Software/imgwire-go/generated"
	imgwirehttp "github.com/Blackhawk-Software/imgwire-go/http"
	"github.com/Blackhawk-Software/imgwire-go/resources"
)

type Client struct {
	apiClient  *generated.APIClient
	httpClient *nethttp.Client
	Options    Options

	CorsOrigins  *resources.CorsOriginsResource
	CustomDomain *resources.CustomDomainResource
	Images       *resources.ImagesResource
	Metrics      *resources.MetricsResource
}

func New(apiKey string, opts ...Option) *Client {
	options := defaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	options = options.normalized()

	httpClient := options.HTTPClient
	if httpClient == nil {
		httpClient = imgwirehttp.NewClient(imgwirehttp.Options{
			Timeout:    options.Timeout,
			MaxRetries: options.MaxRetries,
			Backoff:    options.Backoff,
		})
	}

	cfg := generated.NewConfiguration()
	cfg.Servers = generated.ServerConfigurations{
		{
			URL:         options.BaseURL,
			Description: "imgwire API",
		},
	}
	cfg.HTTPClient = httpClient
	cfg.UserAgent = options.UserAgent
	cfg.AddDefaultHeader("Authorization", "Bearer "+apiKey)
	if options.EnvironmentID != "" {
		cfg.AddDefaultHeader("X-Environment-Id", options.EnvironmentID)
	}

	apiClient := generated.NewAPIClient(cfg)

	return &Client{
		apiClient:    apiClient,
		httpClient:   httpClient,
		Options:      options,
		CorsOrigins:  resources.NewCorsOriginsResource(apiClient),
		CustomDomain: resources.NewCustomDomainResource(apiClient),
		Images:       resources.NewImagesResource(apiClient, httpClient),
		Metrics:      resources.NewMetricsResource(apiClient),
	}
}

func (c *Client) APIClient() *generated.APIClient {
	return c.apiClient
}

func (c *Client) HTTPClient() *nethttp.Client {
	return c.httpClient
}
