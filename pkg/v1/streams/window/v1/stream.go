package window

import (
	"context"
	"time"

	"github.com/TrevinTeacutter/goback/pkg/v1"
	"github.com/TrevinTeacutter/goback/pkg/v1/backoffs/noop/v1"
	transduce "github.com/TrevinTeacutter/transduce-go/pkg/v1"
	"github.com/TrevinTeacutter/transduce-go/pkg/v1/streams/tick/v1"
)

type Option func(*stream)

type stream struct {
	// dependencies
	backoff backoff.Backoff
	cursor  Cursor[time.Time]

	// configuration
	interval time.Duration
	minimum  time.Duration
	maximum  time.Duration
}

func Stream(options ...Option) (transduce.Stream[Window], error) {
	s := &stream{}

	for _, option := range options {
		option(s)
	}

	if s.backoff == nil {
		s.backoff = &noop.Backoff{}
	}

	if s.cursor == nil {
		s.cursor = &Memory[time.Time]{}
	}

	return s.Stream, nil
}

func (s *stream) Stream(ctx context.Context, yield transduce.StreamYield[Window]) error {
	if yield == nil {
		return transduce.NilStreamYieldError
	}

	var timer <-chan time.Time

	switch {
	case s.interval <= 0:
		// by default interval is zero, so this means no ticks will ever happen
		// this will deadlock if nothing else is running in the application
		temp := make(chan time.Time)

		defer close(temp)

		timer = temp
	default:
		ticker := time.NewTicker(s.interval)

		defer ticker.Stop()

		timer = ticker.C
	}

	for {
		select {
		case <-ctx.Done():
			// We could return ctx.Err(), but most consumers typically don't care about
			// bubbling up cancellation/deadline internalErrors because they are the ones that
			// cancelled them in the first place.
			return nil
		case timestamp, ok := <-timer:
			if !ok {
				// Realistically this should never be reached, but it doesn't hurt to add a
				// safeguard in case there is aberrant behavior that is introduced at some point
				return nil
			}

			if err := s.handleTick(ctx, timestamp, yield); err != nil {
				return err
			}
		}
	}
}

func (s *stream) handleTick(ctx context.Context, timestamp time.Time, yield transduce.StreamYield[Window]) error {
	cursor, err := s.cursor.Load(ctx)
	if err != nil {
		// swallow this error, once I figure out a logger interface we won't anymore
		// we don't bubble this up because these errors could be transient
		return nil
	}

	window, ok := s.calculateWindow(timestamp, cursor)
	if !ok {
		return nil
	}

	if err = yield(ctx, window); err != nil {
		s.handleBackoff(ctx)

		// swallow this error, once I figure out a logger interface we won't anymore
		// we don't bubble this up because these errors could be transient
		return nil
	}

	// easier to just reset on every success than try to only reset on the first success
	// after a failure
	s.backoff.Reset()

	// only try to store the cursor if we were successful
	if err = s.cursor.Store(ctx, window.End); err != nil {
		// swallow this error, once I figure out a logger interface we won't anymore
		// we don't bubble this up because these errors could be transient
		return nil
	}

	return nil
}

func (s *stream) calculateWindow(timestamp time.Time, cursor time.Time) (Window, bool) {
	window := Window{
		Tick: tick.Tick{
			Timestamp: timestamp,
		},
		Start: cursor,
		End:   time.Now(),
	}

	switch {
	case cursor.IsZero():
		// by default, we want to be greedy and grab the maximum window size
		window.Start = window.End.Add(-s.maximum)
	default:
		difference := window.Difference()

		// to allow for an interval that is more frequent than the smallest window we want,
		// we just return early and check again next interval
		if difference < s.minimum {
			return Window{}, false
		}

		if difference > s.maximum {
			window.End = window.Start.Add(s.maximum)
		}
	}

	return window, true
}

func (s *stream) handleBackoff(ctx context.Context) {
	duration, err := s.backoff.NextAttempt()
	if err != nil {
		// swallow this error, realistically you shouldn't use a backoff that could
		// return an error, but we would want to log this regardless for that case....
		// once I figure out a logger interface to use
		return
	}

	// with 1.23, this is the better option over time.Sleep since we can be more
	// responsive with context closures
	select {
	case <-ctx.Done():
	case <-time.After(duration):
	}
}
