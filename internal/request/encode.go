package request

import (
	"io"
	"net/http"
)

// Encode encodes the body stream with encoder, and appends the name to
// the "Content-Encoding" header. The encoder should return the name
// of the encoding and the encoded stream.
// Returning an empty string means don't append the name to the "Content-Encoding"
// header, but the returned stream will still be used.
// Note that:
// - The encoder may be called multiple times.
// - Encoder will not be called if Content is nil
func (r *Request) Encode(encoder func(io.ReadCloser) (string, io.ReadCloser)) {
	r.Encoders = append(r.Encoders, encoder)
}

func (r *Request) buildEncoding(req *http.Request) {
	if len(r.Encoders) == 0 || req.Body == nil {
		return
	}
	contentEncodings := req.Header.Get("Content-Encoding")
	var name string
	for _, encoder := range r.Encoders {
		name, req.Body = encoder(req.Body)
		if name != "" {
			if contentEncodings != "" {
				contentEncodings += ", "
			}
			contentEncodings += name
		}
	}
	req.Header.Set("Content-Encoding", contentEncodings)

	oGetBody := req.GetBody
	if oGetBody == nil {
		return
	}
	req.GetBody = func() (io.ReadCloser, error) {
		stream, err := oGetBody()
		if err != nil {
			return stream, err
		}
		for _, encoder := range r.Encoders {
			_, stream = encoder(stream)
		}
		return stream, nil
	}
}
