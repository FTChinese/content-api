package pkg

type AsyncResult[T interface{}] struct {
	Value T
	Err   error
}
