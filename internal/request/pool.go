package request

import "sync"

var p = sync.Pool{New: func() interface{} { return &Request{} }}

func getRequest() *Request {
	return p.Get().(*Request)
}

func putRequest(req *Request) {
	p.Put(req)
}
