package errhandling_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/the-zucc/errhandling"
	errstack "github.com/the-zucc/errhandling/err-stack"
)

const SAMPLE_STRING = "Hello world!"
const ROOT_ERROR = "some error occurred"

func TestErrHandling(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "errhandling2 tests")
}

var _ = Describe("errhandling2 tests", func() {
	It("CatchVal() should work properly for Return()", func() {
		str, err := func() (s string, e error) {
			defer CatchVal(&s, &e)
			func() {
				Return("some string", errors.New("oopsie"))
			}()
			return "", nil
		}()
		Expect(str).To(Equal("some string"))
		Expect(err.Error()).To(Equal("oopsie"))
	})
	It("CatchVal() should work properly for Throw()", func() {
		str, err := func() (s string, e error) {
			defer CatchVal(&s, &e)
			func() {
				Throw(errors.New("oopsie"))
			}()
			return "", nil
		}()
		Expect(str).To(Equal(""))
		Expect(err.Error()).To(Equal("oopsie"))
	})
	It("CatchVal() should work properly for a panic on a errstack.Error", func() {
		var e error
		func() {
			defer func() {
				if err := recover(); err != nil {
					if err, ok := err.(error); ok {
						e = err
					}
				}
			}()
			func() (s string, e error) {
				defer CatchVal(&s, &e)
				func() {
					panic(errstack.New("oops !", errors.New(ROOT_ERROR)))
				}()
				return "", nil
			}()
		}()
		Expect(e).NotTo(BeNil())
	})
	It("Catch() should return an errstack.Error for Return()", func() {
		str, err := func() (s string, e error) {
			defer CatchVal(&s, &e)
			func() { Return(SAMPLE_STRING, errors.New("oops !")) }()
			return "", nil
		}()

		Expect(str).To(Equal(SAMPLE_STRING))
		Expect(err).NotTo(BeNil())
		_, ok := err.(errstack.Error)
		Expect(ok).To(BeTrue())
	})
})
