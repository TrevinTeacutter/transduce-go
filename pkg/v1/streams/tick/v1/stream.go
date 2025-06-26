package tick

import (
	"context"
	"time"

	transduce "github.com/TrevinTeacutter/transduce-go/pkg/v1"
)

type stream struct {
	interval time.Duration
}

func Stream(interval time.Duration) transduce.Stream[Tick] {
	s := &stream{
		interval: interval,
	}

	return s.Stream
}

func (s *stream) Stream(ctx context.Context, yield transduce.StreamYield[Tick]) error {
	if yield == nil {
		return transduce.NilStreamYieldError
	}

	ticker := time.NewTicker(s.interval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// We could return ctx.Err(), but most consumers typically don't care about
			// bubbling up cancellation/deadline errors because they are the ones that
			// cancelled them in the first place.
			return nil
		case tick, ok := <-ticker.C:
			if !ok {
				// Realistically this should never be reached, but it doesn't hurt to add a
				// safeguard in case there is aberrant behavior that is introduced at some point
				return nil
			}

			err := yield(ctx, Tick{
				// We could use the current time, but for most use cases, this value is likely
				// to be ignored as we are just using this to trigger some logic on a frequent
				// interval
				timestamp: tick,
			})
			if err != nil {
				// Would prefer to log this err and move on, but avoiding a logger dependency
				// for simplicity seems better for such an example.
				return err
			}
		}
	}
}
