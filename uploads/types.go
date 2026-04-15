package uploads

import "io"

type ResolvedUpload struct {
	Body          []byte
	ContentLength int64
	FileName      string
	MimeType      string
}

type CreateInput struct {
	FileName       string
	MimeType       string
	CustomMetadata map[string]any
	HashSHA256     string
	IdempotencyKey string
	Purpose        string
}

type Input struct {
	File          any
	FileName      string
	MimeType      string
	ContentLength int64
}

type ReaderAtSeeker interface {
	io.Reader
	io.Seeker
}
