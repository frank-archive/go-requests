package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/frankli0324/go-requests/internal/request/form"
)

type FormBody map[string]interface{}
type MultiPartFormBody map[string]io.Reader

func (r *Request) buildContent(
	defaultContentType string,
) (contentType string, contentLength int64, body io.ReadCloser) {
	if r.Content == nil {
		return "", 0, nil
	}
	switch c := r.Content.(type) {
	case FormBody:
		reader, writer := io.Pipe()
		encoder := form.NewEncoder(writer)
		go func() { writer.CloseWithError(encoder.Encode(c)) }()
		return "application/x-www-form-urlencoded", 0, reader
	case string: // guess later
		body = io.NopCloser(strings.NewReader(c))
		contentLength = int64(len(c))
	case []byte: // guess later
		body = io.NopCloser(bytes.NewReader(c))
		contentLength = int64(len(c))
	case interface{}:
		reader, writer := io.Pipe()
		encoder := json.NewEncoder(writer)
		// errors are returned when Client.Do is called
		go func() { writer.CloseWithError(encoder.Encode(c)) }()
		return "application/json", 0, reader
	}
	if defaultContentType != "" {
		return defaultContentType, 0, body
	}
	// unknown content type yet
	buf := make([]byte, 512)
	contentType = "text/plain"
	if v, ok := body.(io.ReadSeeker); ok {
		if l, err := v.Read(buf); err == nil {
			contentType = http.DetectContentType(buf[:l])
		}
		v.Seek(0, io.SeekStart)
	}
	return
}
