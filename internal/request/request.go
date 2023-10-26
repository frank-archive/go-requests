package request

import (
	"context"
	"io"
	"net/http"
)

type Request struct {
	Method  string
	URL     string
	Proxy   string
	Headers http.Header
	Content interface{}
}

func (r *Request) Build(ctx context.Context) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, nil)
	if err != nil {
		return nil, err
	}
	req.Header = r.Headers
	var setContentType string
	setContentType, req.Body = r.buildContent(req.Header.Get("Content-Type"))
	r.Headers.Set("Content-Type", setContentType)
	req.GetBody = func() (io.ReadCloser, error) {
		_, body := r.buildContent(setContentType)
		return body, nil
	}
	return req, nil
}
