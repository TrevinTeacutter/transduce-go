package functions

func And[T any](children ...Predicate[T]) Predicate[T] {
	return func(value T) bool {
		if len(children) <= 0 {
			return false
		}

		for _, child := range children {
			if !child(value) {
				return false
			}
		}

		return true
	}
}

func Or[T any](children ...Predicate[T]) Predicate[T] {
	return func(value T) bool {
		for _, child := range children {
			if child(value) {
				return true
			}
		}

		return false
	}
}

func Negate[T any](child Predicate[T]) Predicate[T] {
	return func(value T) bool {
		if child == nil {
			return false
		}

		return !child(value)
	}
}
