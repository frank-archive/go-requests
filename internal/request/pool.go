package request

import (
	"net/http"
	"sync"
)

var p = sync.Pool{New: func() interface{} {
	return &http.Request{}
}}

func getRequest() *http.Request {
	return p.Get().(*http.Request)
}

func putRequest(req *http.Request) {
	p.Put(req)
}
