package client

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/frankli0324/go-requests/internal/request"
	"github.com/frankli0324/go-requests/internal/response"
)

type Handler = func(*RequestCtx) error
type Middleware = func(next Handler) Handler

type Client struct {
	Client http.Client

	// I miss partial classes in C#
	middlewares    []Middleware
	chainedHandler Handler
	mwLock         sync.RWMutex
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
	cli.chainedHandler = cli.request
	return cli, cli.Configure(opts...)
}

func (c *Client) CtxDo(
	ctx context.Context, req *request.Request,
) (func(), *response.Response, error) {
	rctx := getRequestCtx(ctx)
	rctx.Request = req
	c.mwLock.RLock()
	call := c.chainedHandler
	c.mwLock.RUnlock()
	return rctx.Done, rctx.Response, call(rctx)
}

func (c *Client) request(rc *RequestCtx) error {
	r, err := rc.Request.Build(rc.Context)
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
