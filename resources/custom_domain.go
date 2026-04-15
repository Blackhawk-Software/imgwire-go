package resources

import (
	"context"

	generated "github.com/imgwire/imgwire-go/generated"
)

type CustomDomainResource struct {
	api generated.CustomDomainAPI
}

func NewCustomDomainResource(apiClient *generated.APIClient) *CustomDomainResource {
	return &CustomDomainResource{api: apiClient.CustomDomainAPI}
}

func (r *CustomDomainResource) Create(
	ctx context.Context,
	input generated.CustomDomainCreateSchema,
) (*generated.CustomDomainSchema, error) {
	value, _, err := r.api.CustomDomainCreate(ctx).
		CustomDomainCreateSchema(input).
		Execute()
	return value, err
}

func (r *CustomDomainResource) Retrieve(
	ctx context.Context,
) (*generated.CustomDomainSchema, error) {
	value, _, err := r.api.CustomDomainRetrieve(ctx).Execute()
	return value, err
}

func (r *CustomDomainResource) TestConnection(
	ctx context.Context,
) (*generated.CustomDomainSchema, error) {
	value, _, err := r.api.CustomDomainTestConnection(ctx).Execute()
	return value, err
}

func (r *CustomDomainResource) Delete(
	ctx context.Context,
) (map[string]string, error) {
	value, _, err := r.api.CustomDomainDelete(ctx).Execute()
	return value, err
}
