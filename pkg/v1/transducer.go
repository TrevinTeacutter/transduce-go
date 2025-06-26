package transduce

import (
	"github.com/TrevinTeacutter/transduce-go/internal/errors"
)

const (
	NilTransducerErr = errors.Const("transducer must not be nil")
)

type Transducer[A, B, R any] func(reducer Reducer[B, R]) Reducer[A, R]
