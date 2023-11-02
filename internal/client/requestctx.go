package client

import (
	"context"
	"sync"

	"github.com/frankli0324/go-requests/internal/request"
	"github.com/frankli0324/go-requests/internal/response"
)

type RequestCtx struct {
	context.Context
	Error    error
	Request  *request.Request
	Response *response.Response
}

// Done recycles the resources used by request context,
// after calling this function, the request context
// MUST NOT be used again.
func (rc *RequestCtx) Done() {
	rc.Request.Done()
	rc.Response.Done()
	putRequestCtx(rc)
}

var p = sync.Pool{New: func() interface{} {
	return &RequestCtx{}
}}

func getRequestCtx(ctx context.Context) *RequestCtx {
	c := p.Get().(*RequestCtx)
	c.Context = ctx
	return c
}

func putRequestCtx(ctx *RequestCtx) {
	ctx.Context = nil
	ctx.Request = nil
	ctx.Response = nil
	ctx.Error = nil
	p.Put(ctx)
}
