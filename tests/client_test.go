package tests

import (
	"context"
	"net/http"
	"testing"

	imgwire "github.com/imgwire/imgwire-go"
)

func TestClientSetsDefaultHeaders(t *testing.T) {
	var authorization string
	var environmentID string
	var userAgent string

	httpClient := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		authorization = r.Header.Get("Authorization")
		environmentID = r.Header.Get("X-Environment-Id")
		userAgent = r.Header.Get("User-Agent")
		return jsonResponse(
			http.StatusOK,
			`{"id":"cd_123","hostname":"images.example.com","environment_id":"env_123","status":"PENDING","certificate_status":"PENDING","cname_record":"images.example.com","cname_value":"cname.imgwire.dev","dcv_cname_record":"_acme.images.example.com","dcv_cname_value":"dcv.imgwire.dev","created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z","last_verified_at":null}`,
			nil,
		), nil
	})}

	client := imgwire.NewClient(
		"sk_test",
		imgwire.WithBaseURL("https://api.example.com"),
		imgwire.WithEnvironmentID("env_123"),
		imgwire.WithHTTPClient(httpClient),
	)

	_, err := client.CustomDomain.Retrieve(context.Background())
	if err != nil {
		t.Fatalf("retrieve custom domain: %v", err)
	}

	if authorization != "Bearer sk_test" {
		t.Fatalf("unexpected authorization header %q", authorization)
	}
	if environmentID != "env_123" {
		t.Fatalf("unexpected environment id header %q", environmentID)
	}
	if userAgent == "" {
		t.Fatalf("expected user agent header to be set")
	}
}
