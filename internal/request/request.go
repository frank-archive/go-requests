package request

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Content interface{}
}

func (r *Request) Build(ctx context.Context) (*http.Request, error) {
	rURL, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}
	req := getRequest()
	req.Method = r.Method
	req.URL = rURL
	req.Header = r.Headers
	req.Host = rURL.Hostname()
	req = req.WithContext(ctx) // copy occurred here

	var setContentType string
	setContentType, req.ContentLength, req.Body = r.buildContent(req.Header.Get("Content-Type"))
	r.Headers.Set("Content-Type", setContentType)
	req.GetBody = func() (io.ReadCloser, error) {
		_, _, body := r.buildContent(setContentType)
		return body, nil
	}
	return req, nil
}
