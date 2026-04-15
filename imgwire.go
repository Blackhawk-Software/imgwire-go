package imgwire

import (
	"github.com/imgwire/imgwire-go/client"
	generated "github.com/imgwire/imgwire-go/generated"
	"github.com/imgwire/imgwire-go/resources"
	"github.com/imgwire/imgwire-go/uploads"
)

type Client = client.Client
type Option = client.Option
type Options = client.Options

type MetricsQuery = resources.MetricsQuery
type UploadInput = uploads.CreateInput

type BulkDeleteImagesSchema = generated.BulkDeleteImagesSchema
type CorsOriginCreateSchema = generated.CorsOriginCreateSchema
type CorsOriginSchema = generated.CorsOriginSchema
type CorsOriginUpdateSchema = generated.CorsOriginUpdateSchema
type CustomDomainCreateSchema = generated.CustomDomainCreateSchema
type CustomDomainSchema = generated.CustomDomainSchema
type ImageDownloadJobCreateSchema = generated.ImageDownloadJobCreateSchema
type ImageDownloadJobSchema = generated.ImageDownloadJobSchema
type ImageSchema = generated.ImageSchema
type MetricsDatasetInterval = generated.MetricsDatasetInterval
type MetricsDatasetsSchema = generated.MetricsDatasetsSchema
type MetricsStatsSchema = generated.MetricsStatsSchema
type StandardUploadCreateSchema = generated.StandardUploadCreateSchema
type StandardUploadResponseSchema = generated.StandardUploadResponseSchema
type UploadTokenCreateResponseSchema = generated.UploadTokenCreateResponseSchema

var (
	WithBackoff       = client.WithBackoff
	WithBaseURL       = client.WithBaseURL
	WithEnvironmentID = client.WithEnvironmentID
	WithHTTPClient    = client.WithHTTPClient
	WithMaxRetries    = client.WithMaxRetries
	WithTimeout       = client.WithTimeout
	WithUserAgent     = client.WithUserAgent
)

func NewClient(apiKey string, opts ...Option) *Client {
	return client.New(apiKey, opts...)
}
