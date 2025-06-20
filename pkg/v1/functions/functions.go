package functions

import (
	"context"
	"iter"

	"github.com/TrevinTeacutter/transduce-go/v1"
)

func Into[A, B any, R []B](ctx context.Context, value A, transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	step := func(ctx context.Context, value B, accumulator R) (R, error) {
		return append(accumulator, value), nil
	}

	return transduce.Reduce(ctx, value, transducer(transduce.Completing(step)), initial)
}

func IntoSeq[A, B any, R []B](ctx context.Context, iterator iter.Seq[A], transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	step := func(ctx context.Context, value B, accumulator R) (R, error) {
		return append(accumulator, value), nil
	}

	return transduce.ReduceSeq(ctx, iterator, transducer(transduce.Completing(step)), initial)
}

func IntoStream[A, B any, R []B](ctx context.Context, stream transduce.Stream[A], transducer transduce.Transducer[A, B, R], initial R) (R, error) {
	step := func(ctx context.Context, value B, accumulator R) (R, error) {
		return append(accumulator, value), nil
	}

	return transduce.ReduceStream(ctx, stream, transducer(transduce.Completing(step)), initial)
}

func Map[A, B, R any](function Function[A, B]) transduce.Transducer[A, B, R] {
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
	}
}

func Filter[A, R any](predicate Predicate[A]) transduce.Transducer[A, A, R] {
	return func(reducer transduce.Reducer[A, R]) transduce.Reducer[A, R] {
		return transduce.NewReducer(
			reducer.Initial,
			reducer.Result,
			func(ctx context.Context, value A, accumulator R) (R, error) {
				if predicate(value) {
					return reducer.Step(ctx, value, accumulator)
				}

				return accumulator, nil
			},
		)
	}
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

func MapSplitSeq[A any, B iter.Seq[C], C, R any](function Function[A, B]) transduce.Transducer[A, C, R] {
	return transduce.Compose(Map[A, B, R](function), SplitSeq[B, C, R]())
}

func MapSplitStream[A any, B transduce.Stream[C], C, R any](function Function[A, B]) transduce.Transducer[A, C, R] {
	return transduce.Compose(Map[A, B, R](function), SplitStream[B, C, R]())
}
