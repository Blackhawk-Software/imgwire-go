package resources

import (
	"context"

	generated "github.com/imgwire/imgwire-go/generated"
	"github.com/imgwire/imgwire-go/pagination"
)

type CorsOriginsResource struct {
	api generated.CorsOriginsAPI
}

func NewCorsOriginsResource(apiClient *generated.APIClient) *CorsOriginsResource {
	return &CorsOriginsResource{api: apiClient.CorsOriginsAPI}
}

func (r *CorsOriginsResource) List(
	ctx context.Context,
	page int,
	limit int,
) (pagination.Page[generated.CorsOriginSchema], error) {
	request := r.api.CorsOriginsList(ctx)
	if page > 0 {
		request = request.Page(int32(page))
	}
	if limit > 0 {
		request = request.Limit(int32(limit))
	}

	data, response, err := request.Execute()
	if err != nil {
		return pagination.Page[generated.CorsOriginSchema]{}, err
	}

	return pagination.Page[generated.CorsOriginSchema]{
		Data:       data,
		Pagination: pagination.ParseHeaders(response.Header),
	}, nil
}

func (r *CorsOriginsResource) ListPages(
	ctx context.Context,
	page int,
	limit int,
) *pagination.PageIterator[generated.CorsOriginSchema] {
	return pagination.NewPageIterator(ctx, page, limit, r.List)
}

func (r *CorsOriginsResource) ListAll(
	ctx context.Context,
	page int,
	limit int,
) *pagination.ItemIterator[generated.CorsOriginSchema] {
	return pagination.NewItemIterator(r.ListPages(ctx, page, limit))
}

func (r *CorsOriginsResource) Create(
	ctx context.Context,
	input generated.CorsOriginCreateSchema,
) (*generated.CorsOriginSchema, error) {
	value, _, err := r.api.CorsOriginsCreate(ctx).
		CorsOriginCreateSchema(input).
		Execute()
	return value, err
}

func (r *CorsOriginsResource) Retrieve(
	ctx context.Context,
	corsOriginID string,
) (*generated.CorsOriginSchema, error) {
	value, _, err := r.api.CorsOriginsRetrieve(ctx, corsOriginID).Execute()
	return value, err
}

func (r *CorsOriginsResource) Update(
	ctx context.Context,
	corsOriginID string,
	input generated.CorsOriginUpdateSchema,
) (*generated.CorsOriginSchema, error) {
	value, _, err := r.api.CorsOriginsUpdate(ctx, corsOriginID).
		CorsOriginUpdateSchema(input).
		Execute()
	return value, err
}

func (r *CorsOriginsResource) Delete(
	ctx context.Context,
	corsOriginID string,
) (map[string]string, error) {
	value, _, err := r.api.CorsOriginsDelete(ctx, corsOriginID).Execute()
	return value, err
}
