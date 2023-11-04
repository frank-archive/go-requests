package requests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/frankli0324/go-requests/request"
)

func ExampleCtxDo() {
	done, resp, err := CtxDo(context.Background(), request.Get("https://www.google.com", nil))
	if err != nil {
		panic(err)
	}
	defer done()
	fmt.Println(resp.Header.Get("Date"))
	c, err := resp.Body()
	fmt.Println(string(c))
}

func ExampleMiddleware() {
	cli, _ := NewClient()
	cli.Use(func(next Handler) Handler {
		return func(rc *RequestCtx) error {
			if rc.Request.Headers == nil {
				rc.Request.Headers = make(http.Header)
			}
			rc.Request.Headers.Add("handler-ok", "1")
			err := next(rc)
			if err == nil {
				fmt.Println(rc.Response.Header.Get("Date"))
			}
			return err
		}
	})
	done, resp, err := cli.CtxDo(context.Background(), request.Get("https://www.google.com", nil))
	if err != nil {
		panic(err)
	}
	defer done()

	c, err := resp.Body()
	fmt.Println(string(c))
}
