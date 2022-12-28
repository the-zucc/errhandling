package errhandling

import "fmt"

func Catch() {
	if err := recover(); err != nil {
		panic(err)
	}
}

/*
func RetErr[T any](err error) T {

}
*/

/*
WithCause returns a function that takes a format as parameter. That
function, when called, will:
  - check for a non-nil error
  - if non-nil, return an error decorated with the underlying error
  - if the error is nil,

describes the cause of the
*/
func WithCause[T any](val T, underlyingErr error) func(string) (T, error) {
	// if error is nil
	if underlyingErr == nil {
		return func(s string) (T, error) {
			return val, nil
		}
	}
	// else, error is not nil
	_, isRootError := underlyingErr.(error)
	_, isHandledError := underlyingErr.(error)

	// if error has not yet been handled
	if !(isRootError || isHandledError) {
		return func(errDesc string) (T, error) {
			var re = fmt.Errorf("%s -> root cause: %s, %s", errDesc, underlyingErr, "%s")
			return val, re // todo check this
		}
	} else if isRootError {
		return func(errDesc string) (T, error) {
			return val, nil // todo recheck this as well
		}
	}
	// else if error has been handled
	return func(errDesc string) (val T, returnedErr error) {
		return val, fmt.Errorf("%s -> caused by: %s", errDesc, underlyingErr)
	}
}
