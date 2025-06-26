package window

import (
	"context"
)

type Cursor[T any] interface {
	Store(ctx context.Context, value T) error
	Load(ctx context.Context) (T, error)
}
