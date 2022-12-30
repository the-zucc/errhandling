package errhandling

import (
	"errors"

	errstack "github.com/the-zucc/errhandling/err-stack"
)

/*
this holds the value and error for returning val-err pairs
up the call stack
*/
type valErr[T any] struct {
	val T
	err error
}

/*
_err wraps errors that are thrown by this library, in order to
identify them when they are caught by the cleanup operations.
*/
type _err struct {
	err error
}

/*
CatchVal() performs the cleanup operation after function execution.
If errors were Thrown, it ensures they are returned up the call stack.
A deferred call to this function should appear in all functions that
are written using the errhandling library and that return a value and
an error.

Example:

	func SomeFunc() (s string, e error) {
		defer CatchVal(&s, &e)
		func(){
			Return("hello world!", nil) // this returns the values
		}()
		return "", nil
	}

	var str, _ := SomeFunc()
*/
func CatchVal[T any](valAddr *T, errAddr *error) {
	if panicInfo := recover(); panicInfo != nil {
		// in the case of a Return[T any](T, error) we need this type check
		if ve, ok := panicInfo.(valErr[T]); ok {
			*valAddr = ve.val
			*errAddr = ve.err
			return
		}
		// in the case of a Throw(error) we need this type check
		if err_, ok := panicInfo.(_err); ok {
			*errAddr = err_.err
			return
		}
		// if we panicked on a stacked error we need to print it out
		if err, ok := panicInfo.(errstack.Error); ok {
			panic(errors.New(err.PrintableError()))
		}
		// otherwise any other panic will panic
		panic(panicInfo)
	}
}

/*
Catch() performs the cleanup operation after function execution.
Similar to CatchVal(), if errors were Thrown, it ensures they are
returned up the call stack. A deferred call to this function should
appear in all functions that are written using the errhandling
library and that return a value and an error.

example:

	func SomeFunc() (e error) {
		defer Catch(&e)
		func(){
			Return("Hello world!", nil) // this returns the values
		}()
		return "", nil
	}

	var _ := SomeFunc()
*/
func Catch(errAddr *error) {
	if panicInfo := recover(); panicInfo != nil {
		if err, ok := panicInfo.(errstack.Error); ok {
			*errAddr = errors.New(err.PrintableError())
		}
		panic(panicInfo)
	}
}

/*
WithCause() returns a function that takes a format as parameter. That
function, when called, will:

  - check for a non-nil error
  - if non-nil, return the value and an error (with the provided message)
    decorated with the underlying error as cause
  - if the error is nil, return the value and the error

example:

	func SomeFunction() (string, error) // this returns an error
	var _, err = WithCause(SomeFunction())("some error occurred")
	// the above decorates the underlying error with the error that resulted
	// from it.
*/
func WithCause[T any](val T, err error) func(errMsg string) (v T, e error) {
	// if error is nil
	if err == nil {
		return func(errMsg string) (T, error) {
			return val, nil
		}
	}
	if se, ok := err.(errstack.Error); ok {
		return func(errMsg string) (T, error) {
			return val, errstack.New(
				errMsg,
				se,
			)
		}
	}
	return func(errMsg string) (T, error) {
		return val, errstack.New(
			errMsg,
			errstack.New(
				err.Error(),
			),
		)
	}
}

/*
With_cause() returns a function that takes a format as parameter. That
function, when called, will:

  - check for a non-nil error
  - if non-nil, return the value and an error (with the provided message)
    decorated with the underlying error as cause
  - if the error is nil, return the value and the error

example:

	func SomeFunction() (string, error) // this returns an error
	var _, err = WithCause(SomeFunction())("some error occurred")
	// the above decorates the underlying error with the error that resulted
	// from it.
*/
func With_cause(err error) func(errMsg string) (e error) {
	return func(errMsg string) error {
		if err == nil {
			return nil
		}
		if se, ok := err.(errstack.Error); ok {
			return errstack.New(errMsg, se)
		}
		return errstack.New(
			errMsg,
			errstack.New(
				err.Error(),
			),
		)
	}
}

/*
ReturnErr needs to be paired with a deferred call to ErrHandling().

When called, it will pass up the call stack any non-nil error
returned by the function call passed in the parameters.

example:

	func SomeOtherFunction() (s string, e error) { // this returns an error
		return "", errors.New("some error occurred")
	}

	func SomeFunction() (e error) {
		defer ErrHandling()
		someStringVar := ReturnErr(SomeOtherFunction())(&e) //
		return nil
	}
*/
func ThrowOnErr[T any](val T, err error) T {
	val, err = WithCause(val, err)("")
	if err != nil {
		panic(valErr[T]{
			val: val,
			err: err,
		})
	}
	return val
}

/*
This throws an error up the call stack. It effectively panics on
a wrapped error, which gets intercepted by Catch() and CatchVal().
This can be used in anonymous functions to return to the higher-
level named function.

Example:

	func SomeFunction() (e error) {
		defer Catch(&e)
		func(){
			Throw(errors.New("oops!"))
		}()
		return nil
	}

	var _ = SomeFunction() // this returns an error with "oops!" as message.
*/
func Throw(err error) {
	panic(_err{err: err})
}

func Return[T any](val T, err error) {
	if _, ok := err.(errstack.Error); ok {
		panic(valErr[T]{
			val: val,
			err: err,
		})
	}
	panic(valErr[T]{
		val: val,
		err: errstack.New(err.Error()),
	})
}
