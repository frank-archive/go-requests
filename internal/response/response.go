package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"sync/atomic"
)

type Response struct {
	*http.Response

	done       atomic.Bool
	bodyClosed atomic.Bool
}

func (r *Response) Body() ([]byte, error) {
	if r.Response.Body == nil {
		return nil, errors.New("response doesn't have valid body")
	}
	if !r.bodyClosed.CompareAndSwap(false, true) {
		return nil, errors.New("body already consumed")
	}
	defer r.Response.Body.Close()
	return io.ReadAll(r.Response.Body)
}

func (r *Response) String() (string, error) {
	b, err := r.Body()
	return string(b), err
}

func (r *Response) JSON(i interface{}) error {
	if r.Response.Body == nil {
		return errors.New("response doesn't have valid body")
	}
	if !r.bodyClosed.CompareAndSwap(false, true) {
		return errors.New("body already consumed")
	}
	dec := json.NewDecoder(r.Response.Body)
	err := dec.Decode(i)
	r.Response.Body.Close()
	return err
}

func Wrap(resp *http.Response) *Response {
	if resp == nil {
		return nil
	}
	ret := getResponse()
	ret.done.Store(false)
	ret.bodyClosed.Store(false)
	ret.Response = resp
	return ret
}

func (r *Response) Done() {
	if !r.done.CompareAndSwap(false, true) {
		return
	}
	if r.bodyClosed.CompareAndSwap(false, true) {
		r.Response.Body.Close()
	}
	putResponse(r)
}
