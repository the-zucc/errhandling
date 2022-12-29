package errhandling_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	handlederr "github.com/the-zucc/errhandling/handled-err"
)

func TestErrHandling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "errhandling tests")
}

func testFunctionNoErr() (string, error) {
	return "", nil
}

func testFunctionErr() (string, error) {
	return "", errors.New("some error occurred")
}

type someErrorType error

var _ = Describe("errhandling tests", func() {
	It("the base assumptions needed for this framework to work", func() {
		Expect(func(cause ...error) bool {
			return cause == nil
		}()).To(BeTrue())
	})
	It("HandledError() should work properly", func() {
		// a root error
		e1 := handlederr.NewError("some root error")
		he1, ok := e1.(handlederr.Error)
		Expect(ok).To(BeTrue())

		// an error with a cause
		e2 := handlederr.NewError("some caused error", he1)
		_, ok = e2.(handlederr.Error)
		Expect(ok).To(BeTrue())
	})
	It("Error.Error() should work properly", func() {
		someError := "some error"
		someOtherError := "some other error"
		e1 := handlederr.NewError(someError)
		e2 := handlederr.NewError(someOtherError, e1)
		Expect(e1.Error()).To(Equal(someError))
		Expect(e2.Error()).To(Equal(someError + " -> " + someOtherError))
	})
	It("Error.PrintableError() should work properly", func() {
		e1 := handlederr.NewError("some root error")
		e2 := handlederr.NewError("some caused error", e1)
		he2, ok := e2.(handlederr.Error)
		Expect(ok).To(BeTrue())
		pe := he2.PrintableError()
		Expect(pe).To(Equal(
			"error:\n\tsome caused error\n\nRoot cause:\n\tsome root error\n\nFull error trace:\n\tsome caused error\n\tcaused by: some root error",
		))
	})
})
