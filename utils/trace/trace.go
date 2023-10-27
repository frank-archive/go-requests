package trace

import (
	"context"
	"net/http"
	"net/http/httptrace"
)

type tTraceCtx string

var kDisableTrace tTraceCtx = "disable-trace"

// WithoutTrace disables tracing on a single request
func WithoutTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, kDisableTrace, true)
}

type TracingTransport struct {
	http.RoundTripper
	Trace *httptrace.ClientTrace
}

func (tr TracingTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if r.Context().Value(kDisableTrace) == true {
		return tr.RoundTripper.RoundTrip(r)
	}
	return tr.RoundTripper.RoundTrip(
		r.WithContext(httptrace.WithClientTrace(r.Context(), tr.Trace)),
	)
}

func (tr TracingTransport) Unwrap() http.RoundTripper {
	return tr.RoundTripper
}
