package transduce

import "context"

type StreamYield[T any] func(context.Context, T) error

type Stream[T any] func(ctx context.Context, yield StreamYield[T]) error
