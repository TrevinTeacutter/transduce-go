package transduce

import (
	"context"

	"github.com/TrevinTeacutter/transduce-go/internal/errors"
)

const (
	NilStepError = errors.Const("step must not be nil")
)

type Step[A, R any] func(ctx context.Context, value A, accumulator R) (R, error)
