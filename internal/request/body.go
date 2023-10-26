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

func (r *Request) buildContent(defaultContentType string) (string, io.ReadCloser) {
	var reader io.ReadCloser
	switch c := r.Content.(type) {
	case FormBody:
		reader, writer := io.Pipe()
		encoder := form.NewEncoder(writer)
		go func() { writer.CloseWithError(encoder.Encode(c)) }()
		return "application/x-www-form-urlencoded", reader
	case string: // guess later
		reader = io.NopCloser(strings.NewReader(c))
	case []byte: // guess later
		reader = io.NopCloser(bytes.NewReader(c))
	case interface{}:
		reader, writer := io.Pipe()
		encoder := json.NewEncoder(writer)
		// errors are returned when Client.Do is called
		go func() { writer.CloseWithError(encoder.Encode(c)) }()
		return "application/json", reader
	}
	if defaultContentType != "" {
		return defaultContentType, reader
	}
	// unknown content type yet
	buf := make([]byte, 512)
	if l, err := reader.Read(buf); err == nil {
		return http.DetectContentType(buf[:l]), reader
	}
	return "text/plain", reader
}
