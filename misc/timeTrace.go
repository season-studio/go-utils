package misc

import (
	"fmt"
	"time"
)

type TimeTrace struct {
	startTime time.Time
}

func TraceTime() *TimeTrace {
	return &TimeTrace{
		startTime: time.Now(),
	}
}

func (t *TimeTrace) Trace() time.Duration {
	now := time.Now()
	return now.Sub(t.startTime)
}

func (t *TimeTrace) Format(d time.Duration) string {
	return fmt.Sprintf("%v:%02d:%02d.%03d", int(d.Hours()), int(d.Minutes()), int(d.Seconds()), d.Nanoseconds()/1e6)
}
