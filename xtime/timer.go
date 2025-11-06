package xtime

import (
	"time"
)

type Timer struct {
	start time.Time
}

func NewTimer() Timer {
	return Timer{
		start: time.Now(),
	}
}

func (t *Timer) Start() time.Time {
	return t.start
}

func (t *Timer) Stop() time.Duration {
	return time.Now().Sub(t.start)
}

func (t *Timer) Stops() string {
	return ManualDuration(t.Stop())
}
