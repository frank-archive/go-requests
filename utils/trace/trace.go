package trace

import (
	"context"
	"net/http/httptrace"
	"sync"
)

var pool = sync.Pool{
	New: func() interface{} {
		return &TraceContext{}
	},
}

// WithTraceContext should be called before sending the request on the request context
// the returning context should be set on the request.
// keep hold of the returning [TraceContext] and call [Recycle] on it after all tracing
// job is done (e.g. finished sending the metrics)
func WithTraceContext(
	ctx context.Context, hook func(ctx *TraceContext, event string), watch watchType,
) (context.Context, *TraceContext) {
	tc := pool.Get().(*TraceContext)
	tc.OnEvent = hook
	tc.when = watch
	tc.Init()
	return httptrace.WithClientTrace(ctx, tc.trace), tc
}

func Recycle(tc *TraceContext) {
	pool.Put(tc)
}

type watchType int

const (
	WatchStart = watchType(iota)
	WatchDone
	WatchAll
)
