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

var ERROR_IN_CATCH = errstack.New("Catch() and CatchVal() must be called with a non-nil pointer")

/*
Catch() and Catch_() perform the cleanup operation after function
execution. If errors were Thrown, it ensures they are returned up
the call stack.

The arguments passed to both of these functions should be the
address of the returned value and/or returned error.

In the case of a function that returns a value and an error, a
deferred call to Catch() should appear as the function's first
statement.

Example:

	func SomeFunc() (s string, e error) {
		defer Catch(&s, &e)
		func(){
			Return("hello world!", nil) // this returns the values
		}()
		return "", nil
	}

	var str, _ := SomeFunc()
*/
func Catch[T any](valAddr *T, errAddr *error) {
	if errAddr == nil {
		panic(ERROR_IN_CATCH)
	}
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
Catch() and Catch_() perform the cleanup operation after function
execution. If errors were Thrown, it ensures they are returned up
the call stack.

The arguments passed to both of these functions should be the
address of the returned value and/or returned error.

In the case of a function that only returns an error, a deferred
call to Catch_() should appear as the function's first statement.

example:

	func SomeFunc() (e error) {
		defer Catch_(&e)
		func(){
			Throw(errstack.New("some error occurred")) // this returns the values
		}()
		return "", nil
	}

	var _ := SomeFunc()
*/
func Catch_(errAddr *error) {
	if errAddr == nil {
		panic(ERROR_IN_CATCH)
	}
	if panicInfo := recover(); panicInfo != nil {
		if err, ok := panicInfo.(errstack.Error); ok {
			*errAddr = errors.New(err.PrintableError())
		}
		panic(panicInfo)
	}
}

/*
WithCauseVal() returns a function that takes a format as parameter. That
function, when called, will:

  - check for a non-nil error
  - if non-nil, return the value and an error (with the provided message)
    decorated with the underlying error as cause
  - if the error is nil, return the value and the error

example:

	func SomeFunction() (string, error) // this returns an error
	var _, err = WithCauseVal(SomeFunction())("some error occurred")
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
WithCause_() returns a function that takes a format as parameter. That
function, should be called as soon as it is returned, with the error
message that should be generated will:

  - check for a non-nil error
  - if non-nil, return the value and an error (with the provided message)
    decorated with the underlying error as cause
  - if the error is nil, return the value and the error

example:

	func SomeFunction() (string, error) // this returns an error
	var _, err = With_cause(SomeFunction())("some error occurred")
	// the above decorates the underlying error with the error that resulted
	// from it.
*/

func WithCause_(err error) func(errMsg string) (e error) {
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
Throw() needs to be paired with a deferred call to Catch().

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
func Throw[T any](val T, err error) T {
	if err != nil {
		panic(valErr[T]{
			val: val,
			err: err,
		})
	}
	return val
}

func Throw_(err error) {
	if err != nil {
		panic(_err{err: err})
	}
}

/*
Return() and Return_() throw an error up the call stack. They effectively
panic on a wrapped error or wrapped value-error pair, which get intercepted
by Catch() and Catch_(). This can be used in anonymous functions to return
to the higher-level named function.

Return_() Example:

	func SomeFunction() (e error) {
		defer Catch_(&e)
		func(){
			Return_(errors.New("oops!"))
		}()
		return nil
	}

	var _ = SomeFunction() // this returns an error with "oops!" as message.
*/
func Return_(err error) {
	panic(_err{err: err})
}

/*
Return() and Return_() throw an error up the call stack. They effectively
panic on a wrapped error or wrapped value-error pair, which get intercepted
by Catch() and Catch_(). This can be used in anonymous functions to return
to the higher-level named function.

Return() Example:

	func SomeFunction() (s string, e error) {
		defer Catch_(&s, &e)
		func(){
			Return("Hello world!", nil)
		}()
		return nil
	}

	var str, _ = SomeFunction() // this returns "Hello world!" and a nil error
*/
func Return[T any](val T, err error) {
	if _, ok := err.(errstack.Error); ok {
		panic(valErr[T]{
			val: val,
			err: err,
		})
	}
	panic(valErr[T]{
		val: val,
		err: err,
	})
}

// TODO check those two
/*
Must() and Must_() will panic on the provided error if not nil.
This is useful for critical operations during application execution,
and statements which's failure would prevent the application from
running at all.

Must() Example:

	func someCriticalFunction() (string, error)

	func main() {
		// TODO add deferred call to finalization function here
		str := Must(SomeCriticalFunction()) // this will panic on error
	}
*/
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

/*
Must() and Must_() will panic on the provided error if not nil.
This is useful for critical operations during application execution,
and statements which's failure would prevent the application from
running at all.

Must_() Example:

	func someCriticalFunction() (error)

	func main() {
		// TODO add deferred call to finalization function here
		Must_(SomeCriticalFunction()) // this will panic on error
	}
*/
func Must_(err error) {
	if err != nil {
		panic(err)
	}
}

/*
OnErr() and OnErr_() will run the provided function on the returned
error if it is not nil.

OnErr() Example:

	func someFunction() (string, error)

	func main() {
		str, err := OnErr(someFunction())(func(err error){
			fmt.Println("error - %s", err)
		})
	}
*/
func OnErr[T any](val T, err error) func(f func(error)) (T, error) {
	return func(f func(error)) (T, error) {
		if err != nil {
			f(err)
			return val, err
		}
		return val, err
	}
}

/*
OnErr() and OnErr_() will run the provided function on the returned
error if it is not nil.

OnErr_() Example:

	func someFunction() (error)

	func main() {
		OnErr_(someFunction())(func(err error){
			fmt.Println("error - %s", err)
		})
	}
*/
func OnErr_(err error) func(f func(error)) {
	return func(f func(error)) {
		if err != nil {
			f(err)
		}
	}
}

/*
OnSuccess() runs the provided function on the result of the function
call if the error is nil. If the error is not nil, the provided function
is not run.

OnSuccess() Example:

	func someFunction() (string, error)

	func main() {
		str := OnSuccess(someFunction())(func(s string) {
			fmt.Println(str) // just some code that runs on success
		})
	}
*/
func OnSuccess[T any](val T, err error) func(f func(T)) (T, error) {
	return func(f func(val T)) (T, error) {
		if err != nil {
			return val, err
		}
		return val, err
	}
}

/*
OnSuccess() runs the provided function if the error is nil. If the
error is not nil, the provided function is not run.

OnSuccess_() Example:

	func someFunction() (error)

	func main() {
		 := OnSuccess(someFunction())(func(){
			fmt.Println(str) // just some code that runs on success
		})
	}
*/
func OnSuccess_(err error) func(f func()) {
	return func(f func()) {
		if err != nil {
			return
		}
		f()
	}
}
