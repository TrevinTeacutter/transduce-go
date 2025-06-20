package functions

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/v1"
)

type Predicate[T any] func(value T) bool

type Function[A, B any] func(ctx context.Context, value A) (B, error)

type WithFunction[A, B, R any] func(ctx context.Context, value A, reducer transduce.Reducer[B, R], accumulator R) (R, error)

type BinaryFunction[A, B, C any] func(ctx context.Context, left A, right B) (C, error)
