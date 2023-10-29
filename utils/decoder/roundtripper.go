package decoder

import (
	"io"
	"net/http"
	"strings"
)

type RoundTripper struct {
	http.RoundTripper
	Decoders map[string]Decoder
}

func (rt *RoundTripper) Unwrap() http.RoundTripper {
	return rt.RoundTripper
}

func (rt *RoundTripper) acceptEncoding() (ret string) {
	if len(rt.Decoders) == 0 {
		return ""
	}
	for name := range rt.Decoders {
		ret += ", " + name
	}
	return ret[2:]
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if hdr := rt.acceptEncoding(); hdr != "" {
		req.Header.Set("Accept-Encoding", hdr)
	}
	resp, err := rt.RoundTripper.RoundTrip(req)
	if err != nil || resp.Body == nil {
		return resp, err
	}
	ce := resp.Header.Get("Content-Encoding")
	if ce == "" {
		return resp, err
	}

	var closer io.Closer = resp.Body
	var reader io.Reader = resp.Body
	for {
		v, rem, _ := strings.Cut(ce, ",")
		if dec, ok := rt.Decoders[strings.TrimSpace(v)]; ok {
			ce = rem
			reader = dec(reader)
			// content length is no longer predictable
			resp.ContentLength = -1
			resp.Header.Del("Content-Length")
		} else {
			break
		}
	}
	if ce := strings.TrimSpace(ce); ce != "" {
		resp.Header.Set("Content-Encoding", ce)
	} else {
		resp.Uncompressed = true
		resp.Header.Del("Content-Encoding")
	}
	resp.Body = readCloserStruct{reader, closer}

	return resp, err
}
