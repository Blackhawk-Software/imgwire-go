package uploads

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Resolve(input Input) (*ResolvedUpload, error) {
	switch file := input.File.(type) {
	case []byte:
		fileName := input.FileName
		if fileName == "" {
			fileName = "upload.bin"
		}
		return &ResolvedUpload{
			Body:          file,
			ContentLength: int64(len(file)),
			FileName:      fileName,
			MimeType:      input.MimeType,
		}, nil
	case *os.File:
		body, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		fileName := input.FileName
		if fileName == "" {
			fileName = filepath.Base(file.Name())
		}
		return &ResolvedUpload{
			Body:          body,
			ContentLength: int64(len(body)),
			FileName:      fileName,
			MimeType:      input.MimeType,
		}, nil
	case io.Reader:
		body, err := io.ReadAll(file)
		if err != nil {
			return nil, err
		}
		fileName := input.FileName
		if fileName == "" {
			fileName = "upload.bin"
		}
		return &ResolvedUpload{
			Body:          body,
			ContentLength: int64(len(body)),
			FileName:      fileName,
			MimeType:      input.MimeType,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported upload input %T", input.File)
	}
}

func (u *ResolvedUpload) Reader() io.Reader {
	return bytes.NewReader(u.Body)
}
