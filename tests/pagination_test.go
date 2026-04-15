package tests

import (
	"context"
	"net/http"
	"testing"

	imgwire "github.com/imgwire/imgwire-go"
	"github.com/imgwire/imgwire-go/pagination"
)

func TestImagesListParsesPaginationHeaders(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		return jsonResponse(
			http.StatusOK,
			`[{"id":"img_1","cdn_url":"https://cdn.example.com/1","created_at":"2026-01-01T00:00:00Z","custom_metadata":{},"deleted_at":null,"environment_id":null,"exif_data":{},"extension":"jpg","hash_sha256":null,"height":1,"idempotency_key":null,"mime_type":"image/jpeg","original_filename":"one.jpg","processed_metadata_at":null,"purpose":null,"size_bytes":1,"status":"READY","updated_at":"2026-01-01T00:00:00Z","upload_token_id":null,"width":1}]`,
			map[string]string{
				"X-Total-Count": "3",
				"X-Page":        "1",
				"X-Limit":       "2",
				"X-Next-Page":   "2",
			},
		), nil
	})}

	client := imgwire.NewClient(
		"sk_test",
		imgwire.WithBaseURL("https://api.example.com"),
		imgwire.WithHTTPClient(httpClient),
	)

	page, err := client.Images.List(context.Background(), 1, 2)
	if err != nil {
		t.Fatalf("list images: %v", err)
	}

	if page.Pagination.TotalCount != 3 {
		t.Fatalf("unexpected total count %d", page.Pagination.TotalCount)
	}
	if page.Pagination.NextPage == nil || *page.Pagination.NextPage != 2 {
		t.Fatalf("unexpected next page %#v", page.Pagination.NextPage)
	}
	if len(page.Data) != 1 {
		t.Fatalf("unexpected page size %d", len(page.Data))
	}
}

func TestImagesListAllIteratesAcrossPages(t *testing.T) {
	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			return jsonResponse(
				http.StatusOK,
				`[{"id":"img_1","cdn_url":"https://cdn.example.com/1","created_at":"2026-01-01T00:00:00Z","custom_metadata":{},"deleted_at":null,"environment_id":null,"exif_data":{},"extension":"jpg","hash_sha256":null,"height":1,"idempotency_key":null,"mime_type":"image/jpeg","original_filename":"one.jpg","processed_metadata_at":null,"purpose":null,"size_bytes":1,"status":"READY","updated_at":"2026-01-01T00:00:00Z","upload_token_id":null,"width":1}]`,
				map[string]string{
					"X-Total-Count": "2",
					"X-Page":        "1",
					"X-Limit":       "1",
					"X-Next-Page":   "2",
				},
			), nil
		case "2":
			return jsonResponse(
				http.StatusOK,
				`[{"id":"img_2","cdn_url":"https://cdn.example.com/2","created_at":"2026-01-01T00:00:00Z","custom_metadata":{},"deleted_at":null,"environment_id":null,"exif_data":{},"extension":"jpg","hash_sha256":null,"height":1,"idempotency_key":null,"mime_type":"image/jpeg","original_filename":"two.jpg","processed_metadata_at":null,"purpose":null,"size_bytes":1,"status":"READY","updated_at":"2026-01-01T00:00:00Z","upload_token_id":null,"width":1}]`,
				map[string]string{
					"X-Total-Count": "2",
					"X-Page":        "2",
					"X-Limit":       "1",
				},
			), nil
		default:
			t.Fatalf("unexpected page query %q", r.URL.Query().Get("page"))
			return nil, nil
		}
	})}

	client := imgwire.NewClient(
		"sk_test",
		imgwire.WithBaseURL("https://api.example.com"),
		imgwire.WithHTTPClient(httpClient),
	)
	iterator := client.Images.ListAll(context.Background(), 1, 1)

	var ids []string
	for iterator.Next() {
		ids = append(ids, iterator.Item().Id)
	}
	if err := iterator.Err(); err != nil {
		t.Fatalf("iterate images: %v", err)
	}
	if len(ids) != 2 || ids[0] != "img_1" || ids[1] != "img_2" {
		t.Fatalf("unexpected ids %#v", ids)
	}
}

func TestPaginationParserTreatsNullLikeMissing(t *testing.T) {
	headers := http.Header{
		"X-Total-Count": []string{"10"},
		"X-Page":        []string{"1"},
		"X-Limit":       []string{"25"},
		"X-Prev-Page":   []string{"null"},
		"X-Next-Page":   []string{""},
	}

	parsed := pagination.ParseHeaders(headers)

	if parsed.TotalCount != 10 {
		t.Fatalf("unexpected total count %d", parsed.TotalCount)
	}
	if parsed.Page != 1 {
		t.Fatalf("unexpected page %d", parsed.Page)
	}
	if parsed.Limit != 25 {
		t.Fatalf("unexpected limit %d", parsed.Limit)
	}
	if parsed.PrevPage != nil {
		t.Fatalf("expected nil prev page, got %#v", parsed.PrevPage)
	}
	if parsed.NextPage != nil {
		t.Fatalf("expected nil next page, got %#v", parsed.NextPage)
	}
}

func TestPaginationParserTreatsNoneLikeMissing(t *testing.T) {
	headers := http.Header{
		"X-Total-Count": []string{"none"},
		"X-Page":        []string{"None"},
		"X-Limit":       []string{"  none  "},
		"X-Prev-Page":   []string{"NONE"},
		"X-Next-Page":   []string{" none "},
	}

	parsed := pagination.ParseHeaders(headers)

	if parsed.TotalCount != 0 {
		t.Fatalf("expected zero total count, got %d", parsed.TotalCount)
	}
	if parsed.Page != 0 {
		t.Fatalf("expected zero page, got %d", parsed.Page)
	}
	if parsed.Limit != 0 {
		t.Fatalf("expected zero limit, got %d", parsed.Limit)
	}
	if parsed.PrevPage != nil {
		t.Fatalf("expected nil prev page, got %#v", parsed.PrevPage)
	}
	if parsed.NextPage != nil {
		t.Fatalf("expected nil next page, got %#v", parsed.NextPage)
	}
}
