package lcl

import "fmt"

type Result[T comparable] struct {
	Value T
	Error error
}

func NewResult[T comparable](v T, err error) *Result[T] {
	return &Result[T]{Value: v, Error: err}
}

func NewValueResult[T comparable](v T) *Result[T] {
	return NewResult(v, nil)
}

func NewErrorResult(err error) *Result[any] {
	return NewResult(any(nil), err)
}

func (r *Result[T]) Unwrap() (T, error) {
	return r.Value, r.Error
}

func (r *Result[T]) Bind(f func(v T) *Result[T]) *Result[T] {
	if r.Error != nil {
		return r
	}
	return f(r.Value)
}

func (r *Result[T]) MustPass(msg string, v ...any) {
	if r.Error != nil {
		panic(fmt.Sprintf(msg+": "+r.Error.Error(), v...))
	}
}

func (r *Result[T]) MustGet(msg string, v ...any) T {
	r.MustPass(msg, v...)
	return r.Value
}

func (r *Result[T]) MustPresent(msg string, v ...any) {
	if IsEmpty(r.Value) {
		panic(fmt.Sprintf(msg, v...))
	}
}
