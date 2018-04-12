package stackint_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var mulModBigFn = func(inputs ...Int1024) Int1024 { return inputs[0].MulModuloBig(&inputs[1], &inputs[2]) }

var _ = Describe("Int1024 Double (2048 bits)", func() {
	Context("MulModuloBig", func() {
		It("should work", func() {
			mod, err := FromString("6277101735386680763835789423207666416102355444464034512896")
			Ω(err).Should(BeNil())
			mod1 := mod.Add(&one)
			RunAllCases(mulModBigFn, []TestCase{
				TestCase{inputsStr: []string{"3", "2", "3"}, expectedStr: "0"},
				TestCase{inputsInt: []Int1024{max, max, max}, expectedStr: "0"},
				TestCase{inputsInt: []Int1024{max, max, mod}, expectedStr: "1"},
				TestCase{inputsInt: []Int1024{max, max, mod1}, expectedStr: "340282366920938463500268095579187314689"},
				TestCase{inputsInt: []Int1024{max, max, two64}, expectedStr: "1"},
			})
		})

		It("should handle edge cases", func() {
			Ω(func() { max.MulModuloBig(&max, &three) }).Should(Panic())
			Ω(func() { max.MulModuloBig(&max, &zero) }).Should(Panic())
		})
	})
})
