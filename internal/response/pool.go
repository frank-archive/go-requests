package response

import "sync"

var p = sync.Pool{New: func() interface{} { return &Response{} }}

func getResponse() *Response {
	return p.Get().(*Response)
}

func putResponse(req *Response) {
	p.Put(req)
}
