package transduce

import "context"

type Stream[T any] func(ctx context.Context, yield func(context.Context, T) error) error
