package transduce

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/internal/errors"
)

const (
	NilStreamYieldError = errors.Const("yield must not be nil")
	NilStreamError      = errors.Const("stream must not be nil")
)

type StreamYield[T any] func(context.Context, T) error

type Stream[T any] func(ctx context.Context, yield StreamYield[T]) error
