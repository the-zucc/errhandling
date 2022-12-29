package handlederr

import "fmt"

type StackedError interface {
	PrintableError() string
}

/*
This struct is used to store all information regarding an error.
It decorates the error and reports it properly (with the nested
causes and such) to the developer.
*/
type Error struct {
	msg       string // the error message
	RootCause *error // the root cause
	Cause     *error // the underlying cause of the error
}

/*
Returns the error message of this error (this comes straight

from the error interface in Go). It returns a string of the
following format:

<some error> -> <some other error> -> some other error
*/
func (e Error) Error() string {
	if e.Cause == nil {
		return e.msg
	}
	return fmt.Sprintf("%s -> %s", *(e.Cause), e.msg)
}

/*
Returns the full printable error message, with the root cause.

Sample format:

	error:
		[<some resulting error>]

	Root cause:
		[<some root error>]

	Full error trace:
		<some resulting error>
		caused by: <some error>
		caused by: <some error>
		caused by: <some root error>

Example usage:

	func Example() handlederror.Error{
		val, err := someFunctionCall()
		if err != nil { // good old if err != nil
			return handlederror.New("oops, something went wrong.", err)
		}
	}

	var errMsg := Example().PrintableError() // this prints
*/
func (e Error) PrintableError() string {
	if se, ok := (*e.RootCause).(Error); ok {
		return fmt.Sprintf(
			"error:\n\t%s\n\nRoot cause:\n\t%s\n\nFull error trace:\n%s",
			e.msg,
			se.msg,
			e.errorTrace(false),
		)
	}
	return fmt.Sprintf(
		"error:\n\t%s\n\nRoot cause:\n\t%s\n\nFull error trace:\n%s",
		e.msg,
		*(e.RootCause),
		e.errorTrace(false),
	)
}

/*
This returns the error trace as a printable string
*/
func (e Error) errorTrace(isCause bool) string {
	// if he is a cause
	if isCause {
		// if he has a cause
		if e.Cause != nil {
			// is the cause a stackedError?
			if cause_, ok := (*e.Cause).(Error); ok {
				// if it is, include its stack trace in the returned message
				return fmt.Sprintf("\tcaused by: %s\n%s", e.msg, cause_.errorTrace(true))
			}
			// if not a stackedError, return the normal message, decorated.
			return fmt.Sprintf("\tcaused by: %s\n\tcaused by: %s", e.msg, *e.Cause)
		}
		return fmt.Sprintf("\tcaused by: %s", e.msg)
	}

	// if we are here, then he is not a cause.
	// if he *has* a cause
	if e.Cause != nil {
		// is the cause a stackedError?
		if cause_, ok := (*e.Cause).(Error); ok {
			// if it is, include its stack trace in the returned message
			return fmt.Sprintf("\t%s\n%s", e.msg, cause_.errorTrace(true))
		}
		// if not a stackedError, return the normal message, decorated.
		return fmt.Sprintf("\t%s\n\tcaused by: %s", e.msg, *(e.Cause))
	}
	return fmt.Sprintf("\t%s", e.msg)
}

// this instanciates a stackedError
func New(msg string, cause ...error) error {
	returnedErr := new(error) // instantiate an error pointer
	if len(cause) == 0 {      // if no cause was provided
		*returnedErr = Error{ // set the error pointer's pointed value to a stackedError
			msg:       msg,
			RootCause: returnedErr,
			Cause:     nil, // this error has no cause, it's a root cause
		}
		return *returnedErr // return the struct
	}

	// if we are here, a cause was provided
	// if the cause is a stackedError
	if hc, isCauseStacked := (cause[0]).(Error); isCauseStacked {
		*returnedErr = Error{
			msg:       msg,
			RootCause: hc.RootCause,
			Cause:     &cause[0],
		}
		return *returnedErr
	}
	// if we are here, the cause is an outside error.
	*returnedErr = Error{
		msg:       msg,
		RootCause: returnedErr,
		Cause:     &cause[0],
	}
	return *returnedErr
}
