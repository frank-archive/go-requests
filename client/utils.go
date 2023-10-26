package client

import (
	"net/http"
)

func shouldGetHttpTransport(tr http.RoundTripper) (*http.Transport, bool) {
	if tr == nil {
		return nil, false
	}
	for _, ok := tr.(*http.Transport); !ok; _, ok = tr.(*http.Transport) {
		if utr, ok := tr.(interface{ Unwrap() http.RoundTripper }); ok {
			tr = utr.Unwrap()
		} else {
			return nil, false
		}
	}
	return tr.(*http.Transport), true
}
