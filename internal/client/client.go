package client

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/frankli0324/go-requests/internal/request"
	"github.com/frankli0324/go-requests/internal/response"
)

type Handler = func(*RequestCtx) error
type Middleware = func(next Handler) Handler

type Client struct {
	Client      http.Client
	Middlewares []Middleware
}

func New(opts ...Option) (*Client, error) {
	cli := &Client{
		Client: http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   30, // the default value is too low
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if v := via[len(via)-1].Context().Value(request.KRedirect); v != false {
					return nil
				} else {
					return http.ErrUseLastResponse
				}
			},
		},
	}
	return cli, cli.Configure(opts...)
}

func (c *Client) Use(mw Middleware) {
	c.Middlewares = append(c.Middlewares, mw)
}

func (c *Client) CtxDo(
	ctx context.Context, req *request.Request,
) (func(), *response.Response, error) {
	call := func(rc *RequestCtx) error {
		r, err := req.Build(ctx)
		if err != nil {
			return err
		}
		resp, err := c.Client.Do(r)
		if err != nil {
			// on error, return a Response that won't be recovered
			return err
		}
		rc.Response = response.Wrap(resp)
		return nil
	}
	for _, mw := range c.Middlewares {
		call = mw(call)
	}
	rctx := getRequestCtx(ctx)
	rctx.Request = req
	if err := call(rctx); err != nil {
		return func() {}, rctx.Response, err
	}

	return rctx.Done, rctx.Response, nil
}
