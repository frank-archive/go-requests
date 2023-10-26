package request

import (
	"net/http"

	"github.com/frankli0324/go-requests/internal/request"
)

type FormBody = request.FormBody
type MultiPartFormBody = request.MultiPartFormBody

func Head(url string, headers http.Header) *request.Request {
	return &request.Request{Method: "HEAD", URL: url, Headers: headers}
}

func Get(url string, headers http.Header) *request.Request {
	return &request.Request{Method: "GET", URL: url, Headers: headers}
}

func Post(url string, headers http.Header, content interface{}) *request.Request {
	return &request.Request{Method: "POST", URL: url, Headers: headers, Content: content}
}

func Put(url string, headers http.Header, content interface{}) *request.Request {
	return &request.Request{Method: "PUT", URL: url, Headers: headers, Content: content}
}

func Delete(url string, headers http.Header, content interface{}) *request.Request {
	return &request.Request{Method: "DELETE", URL: url, Headers: headers, Content: content}
}

func Patch(url string, headers http.Header, content interface{}) *request.Request {
	return &request.Request{Method: "PATCH", URL: url, Headers: headers, Content: content}
}
