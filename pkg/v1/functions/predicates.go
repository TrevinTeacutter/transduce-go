package functions

import (
	"context"
	"errors"
)

func And[T any](children ...Predicate[T]) (Predicate[T], error) {
	if err := childrenCheck(children...); err != nil {
		return nil, err
	}

	return func(ctx context.Context, value T) (bool, error) {
		for _, child := range children {
			test, err := child(ctx, value)
			if err != nil {
				return false, err
			}

			if !test {
				return false, nil
			}
		}

		return true, nil
	}, nil
}

func Or[T any](children ...Predicate[T]) (Predicate[T], error) {
	if err := childrenCheck(children...); err != nil {
		return nil, err
	}

	return func(ctx context.Context, value T) (bool, error) {
		if len(children) <= 0 {
			// If no predicates are provided, we will do the same as the Filter behavior
			// and return true. However like Filter, this may be changed in a future state
			// as this could be unexpected behavior
			return true, nil
		}

		for _, child := range children {
			test, err := child(ctx, value)
			if err != nil {
				return false, err
			}

			if test {
				return true, nil
			}
		}

		return false, nil
	}, nil
}

func Negate[T any](child Predicate[T]) (Predicate[T], error) {
	if child == nil {
		return nil, NilPredicateError
	}

	return func(ctx context.Context, value T) (bool, error) {
		test, err := child(ctx, value)
		if err != nil {
			return false, err
		}

		return !test, nil
	}, nil
}

func childrenCheck[T any](children ...Predicate[T]) error {
	// Let's avoid unexpected behavior by guarding against no child predicates or nil
	// predicates. It would be nice to not have to guard against this.
	if len(children) <= 0 {
		return errors.New("children must have at least one element")
	}

	for _, child := range children {
		if child == nil {
			return NilPredicateError
		}
	}

	return nil
}
