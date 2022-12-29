package errhandling

import (
	"errors"

	handlederr "github.com/the-zucc/errhandling/handled-err"
)

func ErrHandling() {
	if err := recover(); err != nil {
		if he, ok := err.(handlederr.Error); ok {
			panic(errors.New(he.PrintableError()))
		}
		panic(err)
	}
}

func CheckErr[T any](val T, err error) func(*T, *error) T {
	return func(t *T, errAddr *error) T {
		if err != nil {
			*t = val
			*errAddr = err
			panic(nil)
		}
		return val
	}
}

/*
WithCause returns a function that takes a format as parameter. That
function, when called, will:
  - check for a non-nil error
  - if non-nil, return the value and an error (with the provided message)
    decorated with the underlying error as cause
  - if the error is nil, return the value and the error
*/
func WithCause[T any](val T, underlyingErr error) func(errMsg string) (v T, e error) {
	// if error is nil
	if underlyingErr == nil {
		return func(errMsg string) (T, error) {
			return val, nil
		}
	}
	return func(errMsg string) (T, error) {
		return val, handlederr.New(
			errMsg,
			handlederr.New(
				underlyingErr.Error(),
			),
		)
	}
}

/*
ReturnErr needs to be paired with a deferred call to ErrHandling().

When called, it will return any non-nil error returned by the function
call passed in the parameters.
*/
func ReturnErr[T any](val T, underlyingErr error) func(errAddr *error) T {
	if underlyingErr != nil {
		return func(errAddr *error) T {
			*errAddr = underlyingErr
			panic(nil)
		}
	}
	return func(errAddr *error) T {
		return val
	}
}
