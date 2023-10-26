package response

import "net/http"

type Response struct {
	*http.Response
	done bool
}

func Wrap(resp *http.Response) *Response {
	ret := getResponse()
	ret.done = false
	ret.Response = resp
	return ret
}

func (r *Response) Done() {
	if r.done {
		return
	}
	r.Response.Body.Close()
	putResponse(r)
	r.done = true
}
