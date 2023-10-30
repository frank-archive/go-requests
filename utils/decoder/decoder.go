package decoder

import (
	"io"
	"net/http"

	"github.com/frankli0324/go-requests/utils"
)

type Decoder func(io.Reader) io.Reader

func GetOrWrap(tr http.RoundTripper) (*RoundTripper, bool) {
	if tr == nil {
		return nil, false
	}
	innerTr := tr
	for _, ok := innerTr.(*RoundTripper); !ok; _, ok = innerTr.(*RoundTripper) {
		if utr, ok := innerTr.(utils.WrappedRoundTripper); ok {
			innerTr = utr.Unwrap()
		} else {
			return &RoundTripper{tr, map[string]Decoder{}}, false
		}
	}
	return innerTr.(*RoundTripper), true // found decoder in inner roundtrippers
}

type readCloserStruct struct {
	io.Reader
	io.Closer
}
