package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/republicprotocol/republic-go/stackint"
)

var mulModBigFn = func(inputs ...Int1024) Int1024 { return inputs[0].MulModuloBig(&inputs[1], &inputs[2]) }

var _ = Describe("Int1024 Double (2048 bits)", func() {
	Context("arithmetic", func() {
		It("should work", func() {
			RunAllCases(mulModBigFn, []TestCase{
				TestCase{inputsStr: []string{"3", "2", "3"}, expectedStr: "0"},
			})
		})
	})
})
