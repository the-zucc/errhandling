package handlederror

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
func (he Error) Error() string {
	return fmt.Sprintf("%s -> %s", *(he.Cause), he.msg)
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
*/
func (he Error) PrintableError() string {
	return fmt.Sprintf(
		"error:\n\t%s\n\nRoot cause:\n\t%s\n\nFull error trace:\n%s",
		he.msg,
		*(he.RootCause),
		he.errorTrace(false),
	)
}

/*
This returns the error trace as a printable string
*/
func (he Error) errorTrace(isCause bool) string {
	// if he is a cause
	if isCause {
		// if he has a cause
		if he.Cause != nil {
			// is the cause a handledError?
			if cause_, ok := (*he.Cause).(Error); ok {
				// if it is, include its stack trace in the returned message
				return fmt.Sprintf("\tcaused by: %s\n%s", he.msg, cause_.errorTrace(true))
			}
			// if not a handledError, return the normal message, decorated.
			return fmt.Sprintf("\tcaused by: %s\n\tcaused by: %s", he.msg, *he.Cause)
		}
	}

	// if we are here, then he is not a cause.
	// if he *has* a cause
	if he.Cause != nil {
		// is the cause a handledError?
		if cause_, ok := (*he.Cause).(Error); ok {
			// if it is, include its stack trace in the returned message
			return fmt.Sprintf("\t%s\n\tcaused by: %s", he.msg, cause_.errorTrace(true))
		}
		// if not a handledError, return the normal message, decorated.
		return fmt.Sprintf("\t%s\n\tcaused by: %s", he.msg, *(he.Cause))
	}
	return fmt.Sprintf("caused by: %s", he.msg)
}

// this instanciates a handledError
func HandledError(msg string, cause ...error) error {
	returnedErr := new(error) // instantiate an error pointer
	if len(cause) == 0 {      // if no cause was provided
		*returnedErr = Error{ // set the error pointer's pointed value to a handledError
			msg:       msg,
			RootCause: returnedErr,
			Cause:     nil, // this error has no cause, it's a root cause
		}
		return *returnedErr // return the struct
	}

	// if we are here, a cause was provided
	// if the cause is a handledError
	if hc, isCauseHandled := (cause[0]).(Error); isCauseHandled {
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
