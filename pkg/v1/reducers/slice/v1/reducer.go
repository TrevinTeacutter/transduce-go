package slice

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/pkg/v1"
)

var _ transduce.Reducer[any, []any] = (*Reducer[any])(nil)

type Reducer[A any] struct{}

func (r *Reducer[A]) Initial(_ context.Context) ([]A, error) {
	return nil, nil
}

func (r *Reducer[A]) Result(_ context.Context, accumulator []A) ([]A, error) {
	return accumulator, nil
}

func (r *Reducer[A]) Step(_ context.Context, value A, accumulator []A) ([]A, error) {
	return append(accumulator, value), nil
}
