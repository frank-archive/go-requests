package requests

import (
	"context"
	"errors"

	"github.com/frankli0324/go-requests/internal/client"
	"github.com/frankli0324/go-requests/internal/profiles"
	"github.com/frankli0324/go-requests/internal/request"
	"github.com/frankli0324/go-requests/internal/response"
)

type Client = client.Client
type Request = request.Request
type Response = response.Response

var defaultClient = Client{}

func Configure(baseProfile string, opts ...client.Option) error {
	if p, ok := profiles.Get(baseProfile); ok {
		if err := defaultClient.Configure(p...); err != nil {
			return err
		}
	} else {
		return errors.New("invalid profile: " + baseProfile)
	}
	return defaultClient.Configure(opts...)
}

func CtxDo(ctx context.Context, req *Request) (func(), *Response, error) {
	return defaultClient.CtxDo(ctx, req)
}

func RegisterProfile(name string, p []client.Option) {
	profiles.Register(name, p)
}
