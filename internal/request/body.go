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

func (r *Request) buildContent(req *http.Request) {
	setContentType := req.Header.Get("Content-Type")
	setContentType, req.ContentLength, req.Body = r.getBodyStream(setContentType)
	if setContentType != "" {
		r.Headers.Set("Content-Type", setContentType)
	}

	// https://stackoverflow.com/questions/17605915/what-is-the-correct-behavior-expected-of-an-http-post-302-redirect-to-get
	// see: net/http/client.go func redirectBehavior
	// also see: net/http/request.go func isReplayable
	if lastBody, isReadSeeker := req.Body.(io.ReadSeekCloser); isReadSeeker {
		req.GetBody = func() (io.ReadCloser, error) {
			lastBody.Seek(0, io.SeekStart)
			return lastBody, nil
		}
	}
	if r.isReplayable() {
		req.GetBody = func() (io.ReadCloser, error) {
			_, _, body := r.getBodyStream(setContentType)
			return body, nil
		}
	}
}

func (r *Request) isReplayable() bool {
	switch r.Content.(type) {
	case io.ReadSeeker:
		return true
	case io.Reader:
		return false
	}
	return true
}

func (r *Request) getBodyStream(
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
		body = NopSeekerCloser(strings.NewReader(c))
		contentLength = int64(len(c))
	case []byte: // guess later
		body = NopSeekerCloser(bytes.NewReader(c))
		contentLength = int64(len(c))
	case io.ReadCloser:
		body = c
	case io.ReadSeeker:
		body = NopSeekerCloser(c)
	case io.Reader:
		body = io.NopCloser(c)
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
