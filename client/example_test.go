package client

import (
	"compress/flate"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/frankli0324/go-requests"
	"github.com/frankli0324/go-requests/request"
)

func ExampleWithDecoder() {
	cli, err := requests.NewClient(
		WithDecoder("gzip", func(r io.Reader) io.Reader {
			gr, err := gzip.NewReader(r)
			if err != nil {
				// server is sending content-encoding: gzip
				// but not returning valid gzipped data
				return r
			}
			return gr
		}),
		// client.WithDecoder("br", func(r io.Reader) io.Reader { return brotli.NewReader(r) }),
		WithDecoder("deflate", func(r io.Reader) io.Reader { return flate.NewReader(r) }),
	)
	done, resp, err := cli.CtxDo(
		context.Background(), request.Get("https://httpbin.org/headers", nil),
	)
	if err != nil {
		panic(err)
	}
	defer done()
	fmt.Println(resp.Header.Get("Date"))
	c, err := resp.Body()
	fmt.Println(string(c))
}

func TestAAA(t *testing.T) {
	ExampleWithDecoder()
}
