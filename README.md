# `imgwire-go`

`imgwire-go` is the server-side Go SDK for the imgwire API.

It authenticates with Server API Keys, exposes the current server API surface, and layers Go-specific helpers for uploads, pagination, retries, and configuration on top of an OpenAPI-generated base client.

## Installation

```bash
go get github.com/imgwire/imgwire-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"

	imgwire "github.com/imgwire/imgwire-go"
)

func main() {
	client := imgwire.NewClient("sk_...")

	page, err := client.Images.List(context.Background(), 1, 25)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(page.Data))
}
```

## Client Setup

```go
client := imgwire.NewClient(
	"sk_...",
	imgwire.WithBaseURL("https://api.imgwire.dev"),
	imgwire.WithEnvironmentID("env_123"),
	imgwire.WithTimeout(10*time.Second),
	imgwire.WithMaxRetries(2),
)
```

## Resources

The current handwritten SDK surface exposes:

- `client.Images`
- `client.CustomDomain`
- `client.CorsOrigins`
- `client.Metrics`

### Uploads

Upload from an `io.Reader`, file handle, or byte slice:

```go
file, _ := os.Open("file.jpg")
defer file.Close()

image, err := client.Images.Upload(context.Background(), file, imgwire.UploadInput{
	MimeType: "image/jpeg",
})
```

```go
image, err := client.Images.Upload(context.Background(), []byte("payload"), imgwire.UploadInput{
	FileName: "file.jpg",
	MimeType: "image/jpeg",
})
```

### Pagination

List methods return data plus parsed pagination metadata:

```go
page, err := client.Images.List(context.Background(), 1, 25)
```

Iterate over pages:

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

Iterate over every item:

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

## Generation

Install tooling and regenerate:

```bash
yarn install --frozen-lockfile
yarn generate
```

Validation:

```bash
yarn verify-generated
go test ./...
go build ./...
```

Generated files in `generated/` are disposable. Handwritten code belongs outside that directory.
