package request

import (
	"context"
	"net/http"
	"net/url"
)

type TRedirect string

var KRedirect TRedirect = "redirect"

type Request struct {
	Method  string
	URL     string
	Headers http.Header
	Content interface{}

	NoRedirect bool
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

	if r.NoRedirect {
		ctx = context.WithValue(ctx, KRedirect, false)
	}
	req = req.WithContext(ctx) // copy occurred here

	r.buildContent(req)
	return req, nil
}
