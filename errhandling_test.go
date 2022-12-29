package errhandling_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/the-zucc/errhandling"
	handlederr "github.com/the-zucc/errhandling/handled-err"
)

func TestErrHandling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "errhandling tests")
}

func testFunctionNoErr() (string, error) {
	return "", nil
}

const ROOT_ERROR = "some error occurred"

func testFunctionErr() (string, error) {
	return "", errors.New(ROOT_ERROR)
}

type someErrorType error

var _ = Describe("errhandling tests", func() {
	It("the base assumptions needed for this framework to work", func() {
		Expect(func(cause ...error) bool {
			return cause == nil
		}()).To(BeTrue())
	})
	It("StackedError() should work properly", func() {
		// a root error
		e1 := handlederr.New("some root error")
		he1, ok := e1.(handlederr.Error)
		Expect(ok).To(BeTrue())

		// an error with a cause
		e2 := handlederr.New("some caused error", he1)
		_, ok = e2.(handlederr.Error)
		Expect(ok).To(BeTrue())
	})
	It("Error.Error() should work properly", func() {
		someError := "some error"
		someOtherError := "some other error"
		e1 := handlederr.New(someError)
		e2 := handlederr.New(someOtherError, e1)
		Expect(e1.Error()).To(Equal(someError))
		Expect(e2.Error()).To(Equal(someError + " -> " + someOtherError))
	})
	It("Error.PrintableError() should work properly", func() {
		e1 := handlederr.New("some root error")
		e2 := handlederr.New("some caused error", e1)
		he2, ok := e2.(handlederr.Error)
		Expect(ok).To(BeTrue())
		pe := he2.PrintableError()
		Expect(pe).To(Equal(
			"error:\n\tsome caused error\n\nRoot cause:\n\tsome root error\n\nFull error trace:\n\tsome caused error\n\tcaused by: some root error",
		))
	})

	It("WithCause should return a StackedError", func() {
		errMsg := "oopsies"
		err := func() (e error) {
			defer ErrHandling()
			ReturnErr(WithCause(testFunctionErr())(errMsg))(&e)
			return nil
		}()
		he, ok := err.(handlederr.Error)
		Expect(ok).To(BeTrue())
		pe := he.PrintableError()
		Expect(pe).To(Equal(
			"error:\n\t" + errMsg +
				"\n\nRoot cause:\n\t" +
				ROOT_ERROR +
				"\n\nFull error trace:\n\t" +
				errMsg +
				"\n\tcaused by: " + ROOT_ERROR,
		))
	})
})
