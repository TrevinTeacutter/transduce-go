package transduce

import (
	"context"
	"errors"
	"iter"
)

func Compose[A, B, C, R any](left Transducer[A, B, R], right Transducer[B, C, R]) Transducer[A, C, R] {
	return func(reducer Reducer[C, R]) Reducer[A, R] {
		return left(right(reducer))
	}
}

func Completing[A, R any](step Step[A, R]) Reducer[A, R] {
	return NewReducer[A, R](
		func(ctx context.Context) (R, error) {
			var zero R

			return zero, nil
		},
		func(ctx context.Context, accumulator R) (R, error) {
			return accumulator, nil
		},
		func(ctx context.Context, value A, accumulator R) (R, error) {
			return step(ctx, value, accumulator)
		},
	)
}

func Reduce[A, R any](ctx context.Context, value A, reducer Reducer[A, R], accumulator R) (R, error) {
	result := accumulator

	result, err := reducer.Step(ctx, value, result)
	if err != nil && !errors.Is(err, Terminator) {
		return accumulator, err
	}

	return reducer.Result(ctx, result)
}

func ReduceSeq[A, R any](ctx context.Context, sequence iter.Seq[A], reducer Reducer[A, R], accumulator R) (R, error) {
	result := accumulator

	var err error

	for value := range sequence {
		result, err = reducer.Step(ctx, value, result)
		if err != nil && !errors.Is(err, Terminator) {
			return accumulator, err
		}
	}

	return reducer.Result(ctx, result)
}

func ReduceStream[A, R any](ctx context.Context, stream Stream[A], reducer Reducer[A, R], accumulator R) (R, error) {
	result := accumulator

	err := stream(ctx, func(ctx context.Context, value A) error {
		var err error

		result, err = reducer.Step(ctx, value, result)

		return err
	})
	if err != nil && !errors.Is(err, Terminator) {
		return accumulator, err
	}

	return reducer.Result(ctx, result)
}

func Transduce[A, B, R any](ctx context.Context, value A, transducer Transducer[A, B, R], reducer Reducer[B, R]) (R, error) {
	complete := transducer(reducer)

	var zero R

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return Reduce[A, R](ctx, value, complete, accumulator)
}

func TransduceSeq[A, B, R any](ctx context.Context, sequence iter.Seq[A], transducer Transducer[A, B, R], reducer Reducer[B, R]) (R, error) {
	complete := transducer(reducer)

	var zero R

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return ReduceSeq[A, R](ctx, sequence, complete, accumulator)
}

func TransduceStream[A, B, R any](ctx context.Context, stream Stream[A], transducer Transducer[A, B, R], reducer Reducer[B, R]) (R, error) {
	complete := transducer(reducer)

	var zero R

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return ReduceStream[A, R](ctx, stream, complete, accumulator)
}
