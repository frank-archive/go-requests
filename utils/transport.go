package utils

import "net/http"

type WrappedRoundTripper interface {
	http.RoundTripper
	Unwrap() http.RoundTripper
}

func GetHttpTransport(tr http.RoundTripper) (*http.Transport, bool) {
	if tr == nil {
		return nil, false
	}
	for _, ok := tr.(*http.Transport); !ok; _, ok = tr.(*http.Transport) {
		if utr, ok := tr.(WrappedRoundTripper); ok {
			tr = utr.Unwrap()
		} else {
			return nil, false
		}
	}
	return tr.(*http.Transport), true
}
