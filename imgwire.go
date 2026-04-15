package imgwire

import (
	"github.com/imgwire/imgwire-go/client"
	generated "github.com/imgwire/imgwire-go/generated"
	"github.com/imgwire/imgwire-go/images"
	"github.com/imgwire/imgwire-go/resources"
	"github.com/imgwire/imgwire-go/uploads"
)

type Client = client.Client
type Option = client.Option
type Options = client.Options

type Image = images.ImgwireImage
type ImageURLOptions = images.URLOptions
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
type ImageSchema = images.ImgwireImage
type MetricsDatasetInterval = generated.MetricsDatasetInterval
type MetricsDatasetsSchema = generated.MetricsDatasetsSchema
type MetricsStatsSchema = generated.MetricsStatsSchema
type StandardUploadCreateSchema = generated.StandardUploadCreateSchema
type StandardUploadResponseSchema = images.StandardUploadResponse
type UploadTokenCreateResponseSchema = generated.UploadTokenCreateResponseSchema
type URLPreset = images.URLPreset
type GravityType = images.GravityType
type ResizingType = images.ResizingType
type OutputFormat = images.OutputFormat

var (
	FormatAVIF        = images.FormatAVIF
	FormatGIF         = images.FormatGIF
	FormatJPG         = images.FormatJPG
	FormatPNG         = images.FormatPNG
	FormatWEBP        = images.FormatWEBP
	GravityCenter     = images.GravityCenter
	GravityEast       = images.GravityEast
	GravityNorth      = images.GravityNorth
	GravityNorthEast  = images.GravityNorthEast
	GravityNorthWest  = images.GravityNorthWest
	GravitySouth      = images.GravitySouth
	GravitySouthEast  = images.GravitySouthEast
	GravitySouthWest  = images.GravitySouthWest
	GravityWest       = images.GravityWest
	PresetLarge       = images.PresetLarge
	PresetMedium      = images.PresetMedium
	PresetSmall       = images.PresetSmall
	PresetThumbnail   = images.PresetThumbnail
	ResizingAuto      = images.ResizingAuto
	ResizingFill      = images.ResizingFill
	ResizingFillDown  = images.ResizingFillDown
	ResizingFit       = images.ResizingFit
	ResizingForce     = images.ResizingForce
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
