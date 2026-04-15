# `imgwire-go`

[![Go Reference](https://pkg.go.dev/badge/github.com/Blackhawk-Software/imgwire-go.svg)](https://pkg.go.dev/github.com/Blackhawk-Software/imgwire-go)
[![CI](https://github.com/Blackhawk-Software/imgwire-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Blackhawk-Software/imgwire-go/actions/workflows/ci.yml)

`imgwire-go` is the server-side Go SDK for the imgwire API.

Use it in backend services, workers, and jobs to authenticate with a Server API Key, upload files from Go readers, file handles, or byte slices, manage server-side resources, and call the imgwire API without hand-writing request plumbing.

Image values returned from the handwritten resource layer also expose a URL builder so you can generate imgwire transformation URLs directly from API results.

## Installation

```bash
go get github.com/Blackhawk-Software/imgwire-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"os"

	imgwire "github.com/Blackhawk-Software/imgwire-go"
)

func main() {
	client := imgwire.NewClient("sk_...")

	file, err := os.Open("hero.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	image, err := client.Images.Upload(context.Background(), file, imgwire.UploadInput{
		MimeType: "image/jpeg",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(image.Id)
	width := 300
	url, err := image.URL(imgwire.ImageURLOptions{Width: &width})
	if err != nil {
		panic(err)
	}

	fmt.Println(url)
}
```

## Client Setup

Create a client with your server key:

```go
client := imgwire.NewClient("sk_...")
```

Optional configuration:

```go
client := imgwire.NewClient(
	"sk_...",
	imgwire.WithBaseURL("https://api.imgwire.dev"),
	imgwire.WithEnvironmentID("env_123"),
	imgwire.WithTimeout(10*time.Second),
	imgwire.WithMaxRetries(2),
	imgwire.WithBackoff(500*time.Millisecond),
)
```

## Resources

The current handwritten SDK surface exposes these grouped resources:

- `client.Images`
- `client.CustomDomain`
- `client.CorsOrigins`
- `client.Metrics`

### `client.Images`

Image operations and upload workflows.

Image-returning methods return an extended image type with `URL(...)` so your backend can generate transformation URLs without rebuilding imgwire’s query rules itself.

Supported methods:

- `List(ctx, page, limit)`
- `ListPages(ctx, page, limit)`
- `ListAll(ctx, page, limit)`
- `Retrieve(ctx, imageID)`
- `Create(ctx, input, uploadToken)`
- `Upload(ctx, file, options...)`
- `CreateUploadToken(ctx)`
- `CreateBulkDownloadJob(ctx, input)`
- `RetrieveBulkDownloadJob(ctx, imageDownloadJobID)`
- `BulkDelete(ctx, input)`
- `Delete(ctx, imageID)`

List images:

```go
page, err := client.Images.List(context.Background(), 1, 25)
if err != nil {
	panic(err)
}

fmt.Println(page.Data)
fmt.Println(page.Pagination.TotalCount)
```

Iterate page-by-page:

```go
pages := client.Images.ListPages(context.Background(), 1, 100)
for pages.Next() {
	page := pages.Page()
	fmt.Println(page.Pagination.Page, len(page.Data))
}
if err := pages.Err(); err != nil {
	panic(err)
}
```

Iterate every image record:

```go
items := client.Images.ListAll(context.Background(), 1, 100)
for items.Next() {
	image := items.Item()
	fmt.Println(image.Id)
}
if err := items.Err(); err != nil {
	panic(err)
}
```

Retrieve an image by id:

```go
image, err := client.Images.Retrieve(context.Background(), "img_123")
if err != nil {
	panic(err)
}

fmt.Println(image.Id)

width := 300
height := 300
url, err := image.URL(imgwire.ImageURLOptions{
	Width:  &width,
	Height: &height,
})
if err != nil {
	panic(err)
}

fmt.Println(url)
```

Create a standard upload intent directly:

```go
input := imgwire.StandardUploadCreateSchema{}
input.SetFileName("hero.png")
input.SetContentLength(1024)

upload, err := client.Images.Create(context.Background(), input, "")
if err != nil {
	panic(err)
}

fmt.Println(upload.UploadUrl)
```

Upload from a file handle:

```go
file, err := os.Open("hero.jpg")
if err != nil {
	panic(err)
}
defer file.Close()

image, err := client.Images.Upload(context.Background(), file, imgwire.UploadInput{
	MimeType: "image/jpeg",
})
if err != nil {
	panic(err)
}

preset := imgwire.PresetThumbnail
thumbnailURL, err := image.URL(imgwire.ImageURLOptions{
	Preset: &preset,
})
if err != nil {
	panic(err)
}

fmt.Println(thumbnailURL)
```

Upload from a byte slice:

```go
image, err := client.Images.Upload(context.Background(), imageBytes, imgwire.UploadInput{
	FileName: "hero.png",
	MimeType: "image/png",
})
if err != nil {
	panic(err)
}
```

### Image URL Transformations

Image-returning endpoints return image values with a `URL(...)` helper:

```go
image, err := client.Images.Retrieve(context.Background(), "img_123")
if err != nil {
	panic(err)
}

preset := imgwire.PresetThumbnail
background := "#ffffff"
width := 300
height := 300
rotate := 90

thumbnailURL, err := image.URL(imgwire.ImageURLOptions{
	Preset:     &preset,
	Background: &background,
	Width:      &width,
	Height:     &height,
	Rotate:     &rotate,
})
if err != nil {
	panic(err)
}

fmt.Println(thumbnailURL)
```

Supported transformation options:

| Option field        | Output rule           | Description                                                                      |
| ------------------- | --------------------- | -------------------------------------------------------------------------------- |
| `Preset`            | path suffix           | Applies a named preset such as `thumbnail`, `small`, `medium`, or `large`.       |
| `Width`             | `width`               | Sets output width.                                                               |
| `Height`            | `height`              | Sets output height.                                                              |
| `MinWidth`          | `min-width`           | Sets minimum width constraint.                                                   |
| `MinHeight`         | `min-height`          | Sets minimum height constraint.                                                  |
| `ResizingType`      | `resizing_type`       | Controls resize strategy such as `fit`, `fill`, `fill-down`, `force`, or `auto`. |
| `Zoom`              | `zoom`                | Applies zoom scaling.                                                            |
| `DPR`               | `dpr`                 | Sets device pixel ratio scaling.                                                 |
| `Crop`              | `crop`                | Applies crop dimensions and optional gravity.                                    |
| `Gravity`           | `gravity`             | Sets crop/focus gravity.                                                         |
| `Padding`           | `padding`             | Adds padding using 1 to 4 numeric values.                                        |
| `Extend`            | `extend`              | Enables extension behavior with optional gravity.                                |
| `ExtendAspectRatio` | `extend_aspect_ratio` | Extends to preserve aspect ratio with optional gravity.                          |
| `Enlarge`           | `enlarge`             | Allows enlarging beyond the original size when `true`.                           |
| `Background`        | `background`          | Sets background color using hex or `r:g:b`.                                      |
| `Rotate`            | `rotate`              | Rotates by `0`, `90`, `180`, `270`, or `360`.                                    |
| `Flip`              | `flip`                | Flips horizontally/vertically using a `true:false`-style value.                  |
| `Blur`              | `blur`                | Applies blur.                                                                    |
| `Sharpen`           | `sharpen`             | Applies sharpening.                                                              |
| `Pixelate`          | `pixelate`            | Applies pixelation.                                                              |
| `Format`            | `format`              | Changes output format such as `jpg`, `png`, `avif`, `gif`, or `webp`.            |
| `Quality`           | `quality`             | Sets output quality from `0` to `100`.                                           |
| `StripMetadata`     | `strip_metadata`      | Strips metadata when `true`.                                                     |
| `StripColorProfile` | `strip_color_profile` | Strips embedded color profiles when `true`.                                      |
| `KeepCopyright`     | `keep_copyright`      | Preserves copyright metadata when `true`.                                        |

Examples:

```go
background := "#ffffff"
width := 150
height := 150
url, _ := image.URL(imgwire.ImageURLOptions{
	Background: &background,
	Width:      &width,
	Height:     &height,
})
```

```go
stripMetadata := true
url, _ := image.URL(imgwire.ImageURLOptions{
	StripMetadata: &stripMetadata,
})
```

```go
format := imgwire.FormatWEBP
quality := 80
url, _ := image.URL(imgwire.ImageURLOptions{
	Format:  &format,
	Quality: &quality,
})
```

Create an upload token:

```go
uploadToken, err := client.Images.CreateUploadToken(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(uploadToken.Token)
```

Create and inspect a bulk download job:

```go
job, err := client.Images.CreateBulkDownloadJob(
	context.Background(),
	imgwire.ImageDownloadJobCreateSchema{
		ImageIds: []string{"img_123", "img_456"},
	},
)
if err != nil {
	panic(err)
}

refreshed, err := client.Images.RetrieveBulkDownloadJob(context.Background(), job.Id)
if err != nil {
	panic(err)
}

fmt.Println(refreshed.Id)
```

Delete multiple images:

```go
_, err := client.Images.BulkDelete(
	context.Background(),
	imgwire.BulkDeleteImagesSchema{
		ImageIds: []string{"img_123", "img_456"},
	},
)
if err != nil {
	panic(err)
}
```

### `client.CustomDomain`

Custom domain management for your imgwire environment.

Supported methods:

- `Create(ctx, input)`
- `Retrieve(ctx)`
- `TestConnection(ctx)`
- `Delete(ctx)`

Example:

```go
created, err := client.CustomDomain.Create(
	context.Background(),
	imgwire.CustomDomainCreateSchema{
		Hostname: "images.example.com",
	},
)
if err != nil {
	panic(err)
}

current, err := client.CustomDomain.Retrieve(context.Background())
if err != nil {
	panic(err)
}

verification, err := client.CustomDomain.TestConnection(context.Background())
if err != nil {
	panic(err)
}

fmt.Println(created.Id, current.Hostname, verification.Status)
```

### `client.CorsOrigins`

CORS origin management for server-controlled environments.

Supported methods:

- `List(ctx, page, limit)`
- `ListPages(ctx, page, limit)`
- `ListAll(ctx, page, limit)`
- `Create(ctx, input)`
- `Retrieve(ctx, corsOriginID)`
- `Update(ctx, corsOriginID, input)`
- `Delete(ctx, corsOriginID)`

Example:

```go
created, err := client.CorsOrigins.Create(
	context.Background(),
	imgwire.CorsOriginCreateSchema{
		Pattern: "app.example.com",
	},
)
if err != nil {
	panic(err)
}

origins, err := client.CorsOrigins.List(context.Background(), 1, 50)
if err != nil {
	panic(err)
}

updated, err := client.CorsOrigins.Update(
	context.Background(),
	created.Id,
	imgwire.CorsOriginUpdateSchema{
		Pattern: "dashboard.example.com",
	},
)
if err != nil {
	panic(err)
}

fmt.Println(len(origins.Data), updated.Pattern)
```

### `client.Metrics`

Server-side metrics endpoints for dashboards, reporting, and internal tooling.

Supported methods:

- `GetDatasets(ctx, query)`
- `GetStats(ctx, query)`

Example:

```go
dateStart := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
dateEnd := time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC)
interval := imgwire.MetricsDatasetInterval("DAILY")

datasets, err := client.Metrics.GetDatasets(context.Background(), imgwire.MetricsQuery{
	DateStart: &dateStart,
	DateEnd:   &dateEnd,
	Interval:  &interval,
	TZ:        "America/Chicago",
})
if err != nil {
	panic(err)
}

stats, err := client.Metrics.GetStats(context.Background(), imgwire.MetricsQuery{
	DateStart: &dateStart,
	DateEnd:   &dateEnd,
	Interval:  &interval,
	TZ:        "America/Chicago",
})
if err != nil {
	panic(err)
}

fmt.Println(datasets, stats)
```

## Response Shape Notes

- List endpoints exposed through handwritten wrappers return `Page[T]` values with `Data` and parsed pagination metadata.
- `ListPages()` yields paginated result objects across pages through an iterator with `Next()`, `Page()`, and `Err()`.
- `ListAll()` yields individual items across every page through an iterator with `Next()`, `Item()`, and `Err()`.
- Image-returning methods return handwritten image values with `URL(...)` for transformation URL generation.
- Upload helpers return the created image record after the presigned upload completes.

## Development

For local development from this repository:

```bash
make install
```

Regenerate checked-in artifacts:

```bash
make generate
```

Verify generated artifacts are current:

```bash
make verify-generated
```

Run tests and build the module:

```bash
make test
make build
```

Run formatting:

```bash
make format
```

Run the full local CI flow:

```bash
make ci
```

## Repository Notes

- `generated/` is disposable OpenAPI Generator output and should not be edited manually.
- `openapi/raw.openapi.json` is the checked-in raw backend contract snapshot.
- `openapi/sdk.openapi.json` is the SDK-shaped contract emitted by `@imgwire/codegen-core`.
- Handwritten Go code lives outside `generated/`.
- Yarn Classic is used for codegen tooling, and Go modules are used for runtime dependency management.
