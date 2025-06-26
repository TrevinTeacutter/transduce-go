package window

import (
	"context"
	"sync/atomic"
)

type Memory[T any] struct {
	value atomic.Pointer[T]
}

func (m *Memory[T]) Store(_ context.Context, value T) error {
	m.value.Store(&value)

	return nil
}

func (m *Memory[T]) Load(_ context.Context) (value T, err error) {
	actual := m.value.Load()

	if actual == nil {
		return *new(T), nil
	}

	return *actual, nil
}
