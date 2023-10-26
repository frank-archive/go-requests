package client

import (
	"context"
	"net/http"

	"github.com/frankli0324/go-requests/internal/request"
	"github.com/frankli0324/go-requests/internal/response"
)

type Client struct {
	http.Client
}

func (c *Client) CtxDo(
	ctx context.Context, req *request.Request,
) (func(), *response.Response, error) {
	r, err := req.Build(ctx)
	if err != nil {
		return func() {}, nil, err
	}
	resp, err := c.Client.Do(r)
	if err != nil {
		// on error, return a Response that won't be recovered
		return func() {}, response.Wrap(resp), err
	}
	ret := response.Wrap(resp)
	return ret.Done, ret, nil
}
