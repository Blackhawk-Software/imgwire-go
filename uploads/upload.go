package uploads

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

func Put(
	ctx context.Context,
	client *http.Client,
	url string,
	upload *ResolvedUpload,
) error {
	if ctx == nil {
		ctx = context.Background()
	}
	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewReader(upload.Body),
	)
	if err != nil {
		return err
	}
	request.ContentLength = upload.ContentLength
	request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(upload.Body)), nil
	}
	if upload.MimeType != "" {
		request.Header.Set("Content-Type", upload.MimeType)
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("upload failed with status %d", response.StatusCode)
	}

	return nil
}
