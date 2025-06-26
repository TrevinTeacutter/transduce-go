package functions

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/internal/errors"
	"github.com/TrevinTeacutter/transduce-go/pkg/v1"
)

const (
	NilFunctionError       = errors.Const("function must not be nil")
	NilPredicateError      = errors.Const("predicate must not be nil")
	NilWithFunctionError   = errors.Const("with function must not be nil")
	NilBinaryFunctionError = errors.Const("binary function must not be nil")
)

// Function is what it sounds like, it takes in a value and returns either another value (which could be of the same type)
// or an error.
type Function[A, B any] func(ctx context.Context, value A) (B, error)

// Predicate is a special function that returns a boolean, this is used particularly by Filter-esque transducers.
type Predicate[T any] Function[T, bool]

// WithFunction is a bit more overhead than your typical function, but allows you to do cleanup after all downstream
// processing has taken place. A good example of a use case is times when you need to consume a file chunk by chunk,
// but also don't want to read all of it into memory like you would with a normal Function. This allows you to do that.
type WithFunction[A, B, R any] func(ctx context.Context, value A, reducer transduce.Reducer[B, R], accumulator R) (R, error)

// BinaryFunction is like a Function but has two arguments to process rather than one, useful
// for cases like handling maps or slices where you need the index as part of the processing.
type BinaryFunction[A, B, C any] func(ctx context.Context, left A, right B) (C, error)

// note: technically there might be a case for a TrinaryFunction and the like, but would prefer to avoid that
// explosion and keep the surface area of the API pretty small.
