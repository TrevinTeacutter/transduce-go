package tick

import "time"

func WithInterval(value time.Duration) Option {
	return func(s *stream) {
		s.interval = value
	}
}
