package transduce

type Transducer[A, B, R any] func(reducer Reducer[B, R]) Reducer[A, R]
