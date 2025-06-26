package window

import (
	"time"

	"github.com/TrevinTeacutter/transduce-go/pkg/v1/streams/tick/v1"
)

type Window struct {
	tick.Tick
	Start time.Time
	End   time.Time
}

func (w *Window) Difference() time.Duration {
	return w.End.Sub(w.Start)
}
