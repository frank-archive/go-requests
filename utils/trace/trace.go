package trace

import (
	"context"
	"crypto/tls"
	"net/http/httptrace"
	"time"
)

type tTraceCtx string

var kDisableTrace tTraceCtx = "disable-trace"
var kWithTraceCtx tTraceCtx = "with-trace-ctx"

// WithoutTrace disables tracing on a single request
func WithoutTrace(ctx context.Context) context.Context {
	return context.WithValue(ctx, kDisableTrace, true)
}

func WithTraceContext(
	ctx context.Context, hook func(ctx *TraceContext, event string), watch watchType,
) context.Context {
	tc := &TraceContext{}
	tc.Init()
	return context.WithValue(ctx, kWithTraceCtx, tc)
}

type watchType int

const (
	WatchStart = watchType(iota)
	WatchDone
	WatchAll
)

type TraceContext struct {
	start time.Time
	trace *httptrace.ClientTrace

	DNS     TimeSpan
	DNSName string // DNSName is populated after "dns_start" event
	DNSErr  error  // DNSErr is populated after "dns_done" event

	TLSHandshake TimeSpan
	TLSConnState *tls.ConnectionState // TLSConnState is populated after "tls_handshake_done" event
	TLSErr       error                // TLSErr is populated after "tls_handshake_done" event

	Connect        TimeSpan
	ConnectAddress string
	ConnectErr     error

	Total     time.Duration // from start to PutIdleConn
	FirstByte time.Duration // from start to GotFirstResponseByte

	when    watchType
	OnEvent func(ctx *TraceContext, event string)
}

func (tc *TraceContext) hook(e string, is_start, is_done bool) {
	if tc.when == WatchAll ||
		(is_start && tc.when == WatchStart) ||
		(is_done && tc.when == WatchDone) {
		tc.OnEvent(tc, e)
	}
}

func (tc *TraceContext) Stop() {
	tc.Total = time.Since(tc.start)
}

func (tc *TraceContext) Reset() {
	tc.start = time.Now()
}

func (tc *TraceContext) Init() {
	tc.trace = &httptrace.ClientTrace{
		DNSStart: func(di httptrace.DNSStartInfo) {
			tc.DNS.Start = time.Now()
			tc.DNSName = di.Host
			tc.hook("dns_start", true, false)
		},
		DNSDone: func(di httptrace.DNSDoneInfo) {
			tc.DNS.Done = time.Now()
			tc.DNSErr = di.Err
			tc.hook("dns_done", false, true)
		},
		TLSHandshakeStart: func() {
			tc.TLSHandshake.Start = time.Now()
			tc.hook("tls_handshake_start", true, false)
		},
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			tc.TLSHandshake.Done = time.Now()
			tc.TLSConnState = &cs
			tc.TLSErr = err
			tc.hook("tls_handshake_done", false, true)
		},
		ConnectStart: func(network, addr string) {
			tc.Connect.Start = time.Now()
			tc.ConnectAddress = addr
			tc.hook("connect_start", true, false)
		},
		ConnectDone: func(_, _ string, err error) {
			tc.Connect.Done = time.Now()
			tc.ConnectErr = err
			tc.hook("connect_done", false, true)
		},
		GotFirstResponseByte: func() {
			tc.FirstByte = time.Since(tc.start)
			tc.hook("first_byte", false, false)
		},
	}
}
