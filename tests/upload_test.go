package tests

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	imgwire "github.com/imgwire/imgwire-go"
)

func TestUploadSupportsByteSlices(t *testing.T) {
	var uploadedBody string

	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v1/images/standard_upload":
			return jsonResponse(
				http.StatusOK,
				`{"upload_url":"https://uploads.example.com/upload","image":{"id":"img_123","cdn_url":"https://cdn.example.com/1","created_at":"2026-01-01T00:00:00Z","custom_metadata":{},"deleted_at":null,"environment_id":null,"exif_data":{},"extension":"jpg","hash_sha256":null,"height":1,"idempotency_key":null,"mime_type":"image/jpeg","original_filename":"one.jpg","processed_metadata_at":null,"purpose":null,"size_bytes":3,"status":"READY","updated_at":"2026-01-01T00:00:00Z","upload_token_id":null,"width":1}}`,
				nil,
			), nil
		case "/upload":
			body, _ := io.ReadAll(r.Body)
			uploadedBody = string(body)
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
			return nil, nil
		}
	})}

	client := imgwire.NewClient(
		"sk_test",
		imgwire.WithBaseURL("https://api.example.com"),
		imgwire.WithHTTPClient(httpClient),
	)

	image, err := client.Images.Upload(context.Background(), []byte("abc"), imgwire.UploadInput{
		FileName: "file.jpg",
		MimeType: "image/jpeg",
	})
	if err != nil {
		t.Fatalf("upload bytes: %v", err)
	}

	if image.Id != "img_123" {
		t.Fatalf("unexpected image id %q", image.Id)
	}
	if uploadedBody != "abc" {
		t.Fatalf("unexpected uploaded body %q", uploadedBody)
	}
	width := 150
	height := 150
	url, err := image.URL(imgwire.ImageURLOptions{
		Width:  &width,
		Height: &height,
	})
	if err != nil {
		t.Fatalf("build image url from uploaded image: %v", err)
	}
	if url != "https://cdn.example.com/1?height=150&width=150" {
		t.Fatalf("unexpected transformed url %q", url)
	}
}

func TestUploadSupportsFiles(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "upload-*.jpg")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString("payload"); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		t.Fatalf("seek temp file: %v", err)
	}

	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Path {
		case "/api/v1/images/standard_upload":
			return jsonResponse(
				http.StatusOK,
				`{"upload_url":"https://uploads.example.com/upload","image":{"id":"img_file","cdn_url":"https://cdn.example.com/1","created_at":"2026-01-01T00:00:00Z","custom_metadata":{},"deleted_at":null,"environment_id":null,"exif_data":{},"extension":"jpg","hash_sha256":null,"height":1,"idempotency_key":null,"mime_type":"image/jpeg","original_filename":"one.jpg","processed_metadata_at":null,"purpose":null,"size_bytes":7,"status":"READY","updated_at":"2026-01-01T00:00:00Z","upload_token_id":null,"width":1}}`,
				nil,
			), nil
		case "/upload":
			body, _ := io.ReadAll(r.Body)
			if strings.TrimSpace(string(body)) != "payload" {
				t.Fatalf("unexpected uploaded file body %q", string(body))
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		default:
			t.Fatalf("unexpected path %q", r.URL.Path)
			return nil, nil
		}
	})}

	client := imgwire.NewClient(
		"sk_test",
		imgwire.WithBaseURL("https://api.example.com"),
		imgwire.WithHTTPClient(httpClient),
	)

	image, err := client.Images.Upload(context.Background(), file, imgwire.UploadInput{})
	if err != nil {
		t.Fatalf("upload file: %v", err)
	}
	if image.Id != "img_file" {
		t.Fatalf("unexpected image id %q", image.Id)
	}
}
