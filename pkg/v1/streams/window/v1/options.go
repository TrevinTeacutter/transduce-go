package window

import (
	"time"

	backoff "github.com/TrevinTeacutter/goback/pkg/v1"
)

func WithBackoff(value backoff.Backoff) Option {
	return func(s *stream) {
		s.backoff = value
	}
}

func WithCursor(value Cursor[time.Time]) Option {
	return func(s *stream) {
		s.cursor = value
	}
}

func WithIntervalAndConstantWindow(interval time.Duration, window time.Duration) Option {
	return func(s *stream) {
		s.interval = interval
		s.minimum = window
		s.maximum = window
	}
}

func WithIntervalAndAdaptiveWindow(interval time.Duration, minimum time.Duration, maximum time.Duration) Option {
	return func(s *stream) {
		s.interval = interval
		s.minimum = minimum
		s.maximum = maximum
	}
}
