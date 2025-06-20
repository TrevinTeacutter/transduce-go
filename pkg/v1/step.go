package transduce

import "context"

type Step[A, R any] func(ctx context.Context, value A, accumulator R) (R, error)
