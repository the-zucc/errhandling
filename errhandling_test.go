package errhandling_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	handlederror "github.com/the-zucc/errhandling/handled-error"
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
	It("the HandledError() function should work properly", func() {
		// a root error
		e1 := handlederror.HandledError("some root error")
		he1, ok := e1.(handlederror.Error)
		Expect(ok).To(BeTrue())
		pe := he1.PrintableError()

		Expect(ok).To(BeTrue())
		// an error with a cause
		e2 := handlederror.HandledError("some caused error", he1)
		he2, ok := e2.(handlederror.Error)
		Expect(ok).To(BeTrue())
		pe = he2.PrintableError()
		Expect(pe).To(Equal(
			"error:\n\tsome caused error\n\nRoot cause:\n\tsome root error\n\nFull error trace:\n\tsome caused error\n\tcaused by:\nsome root error",
		))
	})
})
