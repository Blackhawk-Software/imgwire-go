package resources

import (
	"context"
	"time"

	generated "github.com/Blackhawk-Software/imgwire-go/generated"
)

type MetricsResource struct {
	api generated.MetricsAPI
}

func NewMetricsResource(apiClient *generated.APIClient) *MetricsResource {
	return &MetricsResource{api: apiClient.MetricsAPI}
}

type MetricsQuery struct {
	DateStart *time.Time
	DateEnd   *time.Time
	Interval  *generated.MetricsDatasetInterval
	TZ        string
}

func (r *MetricsResource) GetDatasets(
	ctx context.Context,
	query MetricsQuery,
) (*generated.MetricsDatasetsSchema, error) {
	request := r.api.MetricsGetDatasets(ctx)
	if query.DateStart != nil {
		request = request.DateStart(*query.DateStart)
	}
	if query.DateEnd != nil {
		request = request.DateEnd(*query.DateEnd)
	}
	if query.Interval != nil {
		request = request.Interval(*query.Interval)
	}
	if query.TZ != "" {
		request = request.Tz(query.TZ)
	}
	value, _, err := request.Execute()
	return value, err
}

func (r *MetricsResource) GetStats(
	ctx context.Context,
	query MetricsQuery,
) (*generated.MetricsStatsSchema, error) {
	request := r.api.MetricsGetStats(ctx)
	if query.DateStart != nil {
		request = request.DateStart(*query.DateStart)
	}
	if query.DateEnd != nil {
		request = request.DateEnd(*query.DateEnd)
	}
	if query.Interval != nil {
		request = request.Interval(*query.Interval)
	}
	if query.TZ != "" {
		request = request.Tz(query.TZ)
	}
	value, _, err := request.Execute()
	return value, err
}
