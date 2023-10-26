package request

func New(method, url string) (func(), *Request) {
	req := getRequest()
	req.Method = method
	req.URL = url
	return func() { putRequest(req) }, req
}
