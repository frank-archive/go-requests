package request

import (
	"net/http"
	"sync"
)

var p = sync.Pool{New: func() interface{} {
	return &http.Request{
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}
}}

func getRequest() *http.Request {
	return p.Get().(*http.Request)
}

func PutRequest(req *http.Request) {
	p.Put(req)
}
