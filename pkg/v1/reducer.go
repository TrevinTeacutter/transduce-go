package transduce

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/internal/errors"
)

const (
	NilReducerError = errors.Const("reducer must not be nil")
)

var _ Reducer[any, any] = (*reducer[any, any])(nil)

type Reducer[A, R any] interface {
	Initial(ctx context.Context) (R, error)
	Result(ctx context.Context, accumulator R) (R, error)
	Step(ctx context.Context, value A, accumulator R) (R, error)
}

type reducer[A, R any] struct {
	initial func(context.Context) (R, error)
	result  func(context.Context, R) (R, error)
	step    Step[A, R]
}

func NewReducer[A, R any](
	initial func(context.Context) (R, error),
	result func(context.Context, R) (R, error),
	step Step[A, R],
) Reducer[A, R] {
	return &reducer[A, R]{
		initial: initial,
		result:  result,
		step:    step,
	}
}

func (r *reducer[A, R]) Initial(ctx context.Context) (R, error) {
	return r.initial(ctx)
}

func (r *reducer[A, R]) Result(ctx context.Context, accumulator R) (R, error) {
	return r.result(ctx, accumulator)
}

func (r *reducer[A, R]) Step(ctx context.Context, value A, accumulator R) (R, error) {
	return r.step(ctx, value, accumulator)
}
