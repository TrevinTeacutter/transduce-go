package transduce

import (
	"context"
	"errors"
	"iter"

	internalErrors "github.com/TrevinTeacutter/transduce-go/internal/errors"
)

const (
	NilSeqError = internalErrors.Const("seq must not be nil")
)

func Compose[A, B, C, R any](left Transducer[A, B, R], right Transducer[B, C, R]) (Transducer[A, C, R], error) {
	if left == nil || right == nil {
		return nil, NilTransducerErr
	}

	return func(reducer Reducer[C, R]) Reducer[A, R] {
		return left(right(reducer))
	}, nil
}

func Completing[A, R any](step Step[A, R]) (Reducer[A, R], error) {
	if step == nil {
		return nil, NilStepError
	}

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
	), nil
}

func Reduce[A, R any](ctx context.Context, value A, reducer Reducer[A, R], accumulator R) (R, error) {
	if reducer == nil {
		return accumulator, NilReducerError
	}

	result := accumulator

	result, err := reducer.Step(ctx, value, result)
	if err != nil && !errors.Is(err, Terminator) {
		return accumulator, err
	}

	return reducer.Result(ctx, result)
}

func ReduceSeq[A, R any](ctx context.Context, sequence iter.Seq[A], reducer Reducer[A, R], accumulator R) (R, error) {
	if sequence == nil {
		return accumulator, NilSeqError
	}

	if reducer == nil {
		return accumulator, NilReducerError
	}

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
	if stream == nil {
		return accumulator, NilStreamError
	}

	if reducer == nil {
		return accumulator, NilReducerError
	}

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
	var zero R

	if transducer == nil {
		return zero, NilStreamError
	}

	if reducer == nil {
		return zero, NilReducerError
	}

	complete := transducer(reducer)

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return Reduce[A, R](ctx, value, complete, accumulator)
}

func TransduceSeq[A, B, R any](ctx context.Context, sequence iter.Seq[A], transducer Transducer[A, B, R], reducer Reducer[B, R]) (R, error) {
	var zero R

	if sequence == nil {
		return zero, NilSeqError
	}

	if transducer == nil {
		return zero, NilStreamError
	}

	if reducer == nil {
		return zero, NilReducerError
	}

	complete := transducer(reducer)

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return ReduceSeq[A, R](ctx, sequence, complete, accumulator)
}

func TransduceStream[A, B, R any](ctx context.Context, stream Stream[A], transducer Transducer[A, B, R], reducer Reducer[B, R]) (R, error) {
	var zero R

	if stream == nil {
		return zero, NilStreamError
	}

	if transducer == nil {
		return zero, NilStreamError
	}

	if reducer == nil {
		return zero, NilReducerError
	}

	complete := transducer(reducer)

	accumulator, err := complete.Initial(ctx)
	if err != nil {
		return zero, err
	}

	return ReduceStream[A, R](ctx, stream, complete, accumulator)
}
