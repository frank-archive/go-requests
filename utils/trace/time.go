package trace

import "time"

type TimeSpan struct {
	Start time.Time
	Done  time.Time
}

func (t TimeSpan) Duration() time.Duration {
	return t.Done.Sub(t.Start)
}

func (t TimeSpan) Milliseconds() int64 {
	return t.Done.Sub(t.Start).Milliseconds()
}

func NewTimeSpan(start time.Time) TimeSpan {
	return TimeSpan{Start: start}
}
