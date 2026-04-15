package resources

import (
	"context"
	nethttp "net/http"

	generated "github.com/Blackhawk-Software/imgwire-go/generated"
	"github.com/Blackhawk-Software/imgwire-go/images"
	"github.com/Blackhawk-Software/imgwire-go/pagination"
	"github.com/Blackhawk-Software/imgwire-go/uploads"
)

type ImagesResource struct {
	api        generated.ImagesAPI
	httpClient *nethttp.Client
}

func NewImagesResource(
	apiClient *generated.APIClient,
	httpClient *nethttp.Client,
) *ImagesResource {
	return &ImagesResource{
		api:        apiClient.ImagesAPI,
		httpClient: httpClient,
	}
}

func (r *ImagesResource) List(
	ctx context.Context,
	page int,
	limit int,
) (pagination.Page[images.ImgwireImage], error) {
	request := r.api.ImagesList(ctx)
	if page > 0 {
		request = request.Page(int32(page))
	}
	if limit > 0 {
		request = request.Limit(int32(limit))
	}

	data, response, err := request.Execute()
	if err != nil {
		return pagination.Page[images.ImgwireImage]{}, err
	}

	extended := make([]images.ImgwireImage, 0, len(data))
	for _, image := range data {
		extended = append(extended, images.ExtendImage(image))
	}

	return pagination.Page[images.ImgwireImage]{
		Data:       extended,
		Pagination: pagination.ParseHeaders(response.Header),
	}, nil
}

func (r *ImagesResource) ListPages(
	ctx context.Context,
	page int,
	limit int,
) *pagination.PageIterator[images.ImgwireImage] {
	return pagination.NewPageIterator(ctx, page, limit, r.List)
}

func (r *ImagesResource) ListAll(
	ctx context.Context,
	page int,
	limit int,
) *pagination.ItemIterator[images.ImgwireImage] {
	return pagination.NewItemIterator(r.ListPages(ctx, page, limit))
}

func (r *ImagesResource) Retrieve(
	ctx context.Context,
	imageID string,
) (*images.ImgwireImage, error) {
	value, _, err := r.api.ImagesRetrieve(ctx, imageID).Execute()
	if err != nil {
		return nil, err
	}
	return images.ExtendImagePtr(value), nil
}

func (r *ImagesResource) Create(
	ctx context.Context,
	input generated.StandardUploadCreateSchema,
	uploadToken string,
) (*images.StandardUploadResponse, error) {
	request := r.api.ImagesCreate(ctx).
		StandardUploadCreateSchema(input)
	if uploadToken != "" {
		request = request.UploadToken(uploadToken)
	}
	value, _, err := request.Execute()
	if err != nil {
		return nil, err
	}
	return images.ExtendStandardUploadResponse(value), nil
}

func (r *ImagesResource) CreateUploadToken(
	ctx context.Context,
) (*generated.UploadTokenCreateResponseSchema, error) {
	value, _, err := r.api.ImagesCreateUploadToken(ctx).Execute()
	return value, err
}

func (r *ImagesResource) CreateBulkDownloadJob(
	ctx context.Context,
	input generated.ImageDownloadJobCreateSchema,
) (*generated.ImageDownloadJobSchema, error) {
	value, _, err := r.api.ImagesCreateBulkDownloadJob(ctx).
		ImageDownloadJobCreateSchema(input).
		Execute()
	return value, err
}

func (r *ImagesResource) RetrieveBulkDownloadJob(
	ctx context.Context,
	imageDownloadJobID string,
) (*generated.ImageDownloadJobSchema, error) {
	value, _, err := r.api.ImagesRetrieveBulkDownloadJob(ctx, imageDownloadJobID).Execute()
	return value, err
}

func (r *ImagesResource) BulkDelete(
	ctx context.Context,
	input generated.BulkDeleteImagesSchema,
) (map[string]string, error) {
	value, _, err := r.api.ImagesBulkDelete(ctx).
		BulkDeleteImagesSchema(input).
		Execute()
	return value, err
}

func (r *ImagesResource) Delete(
	ctx context.Context,
	imageID string,
) (map[string]string, error) {
	value, _, err := r.api.ImagesDelete(ctx, imageID).Execute()
	return value, err
}

func (r *ImagesResource) Upload(
	ctx context.Context,
	file any,
	inputs ...uploads.CreateInput,
) (*images.ImgwireImage, error) {
	var input uploads.CreateInput
	if len(inputs) > 0 {
		input = inputs[0]
	}

	resolved, err := uploads.Resolve(uploads.Input{
		File:     file,
		FileName: input.FileName,
		MimeType: input.MimeType,
	})
	if err != nil {
		return nil, err
	}

	createInput := generated.NewStandardUploadCreateSchema(resolved.FileName)
	createInput.SetContentLength(int32(resolved.ContentLength))
	if input.MimeType != "" {
		mimeType := generated.SupportedMimeType(input.MimeType)
		createInput.SetMimeType(mimeType)
	}
	if input.HashSHA256 != "" {
		createInput.SetHashSha256(input.HashSHA256)
	}
	if input.IdempotencyKey != "" {
		createInput.SetIdempotencyKey(input.IdempotencyKey)
	}
	if input.Purpose != "" {
		createInput.SetPurpose(input.Purpose)
	}
	if len(input.CustomMetadata) > 0 {
		customMetadata, err := toCustomMetadata(input.CustomMetadata)
		if err != nil {
			return nil, err
		}
		createInput.SetCustomMetadata(customMetadata)
	}

	created, err := r.Create(ctx, *createInput, "")
	if err != nil {
		return nil, err
	}

	err = uploads.Put(ctx, r.httpClient, created.UploadURL, resolved)
	if err != nil {
		return nil, err
	}

	image := created.Image
	return &image, nil
}
