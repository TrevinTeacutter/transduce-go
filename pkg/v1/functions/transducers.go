package functions

import (
	"context"
	"iter"

	"github.com/TrevinTeacutter/transduce-go/pkg/v1"
)

func into[A, B any, R []B](transducer transduce.Transducer[A, B, R]) (transduce.Reducer[B, R], error) {
	if transducer == nil {
		return nil, transduce.NilTransducerErr
	}

	step := func(ctx context.Context, value B, accumulator R) (R, error) {
		return append(accumulator, value), nil
	}

	completed, err := transduce.Completing(step)
	if err != nil {
		return nil, err
	}

	return completed, nil
}

func Into[A, B any, R []B](ctx context.Context, value A, transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	completed, err := into(transducer)
	if err != nil {
		return initial, err
	}

	return transduce.Reduce(ctx, value, transducer(completed), initial)
}

func IntoSeq[A, B any, R []B](ctx context.Context, sequence iter.Seq[A], transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	if sequence == nil {
		return initial, transduce.NilSeqError
	}

	completed, err := into(transducer)
	if err != nil {
		return initial, err
	}

	return transduce.ReduceSeq(ctx, sequence, transducer(completed), initial)
}

func IntoStream[A, B any, R []B](ctx context.Context, stream transduce.Stream[A], transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	if stream == nil {
		return initial, transduce.NilStreamError
	}

	completed, err := into(transducer)
	if err != nil {
		return initial, err
	}

	return transduce.ReduceStream(ctx, stream, transducer(completed), initial)
}

func Map[A, B, R any](function Function[A, B]) (transduce.Transducer[A, B, R], error) {
	if function == nil {
		return nil, NilFunctionError
	}

	return func(reducer transduce.Reducer[B, R]) transduce.Reducer[A, R] {
		return transduce.NewReducer(
			reducer.Initial,
			reducer.Result,
			func(ctx context.Context, value A, accumulator R) (R, error) {
				result, err := function(ctx, value)
				if err != nil {
					return accumulator, err
				}

				return reducer.Step(ctx, result, accumulator)
			},
		)
	}, nil
}

func Filter[A, R any](predicate Predicate[A]) (transduce.Transducer[A, A, R], error) {
	if predicate == nil {
		return nil, NilPredicateError
	}

	return func(reducer transduce.Reducer[A, R]) transduce.Reducer[A, R] {
		return transduce.NewReducer(
			reducer.Initial,
			reducer.Result,
			func(ctx context.Context, value A, accumulator R) (R, error) {

				test, err := predicate(ctx, value)
				if err != nil {
					return accumulator, err
				}

				if test {
					return reducer.Step(ctx, value, accumulator)
				}

				return accumulator, nil
			},
		)
	}, nil
}

func SplitSeq[A iter.Seq[B], B, R any]() transduce.Transducer[A, B, R] {
	return func(reducer transduce.Reducer[B, R]) transduce.Reducer[A, R] {
		return transduce.NewReducer(
			reducer.Initial,
			reducer.Result,
			func(ctx context.Context, value A, accumulator R) (R, error) {
				return transduce.ReduceSeq[B, R](ctx, iter.Seq[B](value), reducer, accumulator)
			},
		)
	}
}

func SplitStream[A transduce.Stream[B], B, R any]() transduce.Transducer[A, B, R] {
	return func(reducer transduce.Reducer[B, R]) transduce.Reducer[A, R] {
		return transduce.NewReducer(
			reducer.Initial,
			reducer.Result,
			func(ctx context.Context, value A, accumulator R) (R, error) {
				return transduce.ReduceStream[B, R](ctx, transduce.Stream[B](value), reducer, accumulator)
			},
		)
	}
}

func MapSplitSeq[A any, B iter.Seq[C], C, R any](function Function[A, B]) (transduce.Transducer[A, C, R], error) {
	mapper, err := Map[A, B, R](function)
	if err != nil {
		return nil, err
	}

	return transduce.Compose(mapper, SplitSeq[B, C, R]())
}

func MapSplitStream[A any, B transduce.Stream[C], C, R any](function Function[A, B]) (transduce.Transducer[A, C, R], error) {
	mapper, err := Map[A, B, R](function)
	if err != nil {
		return nil, err
	}

	return transduce.Compose(mapper, SplitStream[B, C, R]())
}
