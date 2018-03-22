package stackint_test

// import (
// 	// "math/big"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/republicprotocol/republic-go/stackint"
// )

// var _ = Describe("Int1024 montgomery", func() {

// 	prime := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")

// 	Context("when multiplying", func() {
// 		It("should return the right result for 1024 bit numbers", func() {

// 			testcases := []TestCase{
// 				TestCase{inputsStr: []string{"3", "2"}, expectedStr: "6"},
// 				TestCase{inputsStr: []string{"2", "3"}, expectedStr: "6"},
// 				TestCase{inputsStr: []string{"8", "1"}, expectedStr: "8"},
// 				TestCase{inputsStr: []string{"7", "0"}, expectedStr: "0"},
// 				TestCase{inputsStr: []string{"101", "102"}, expectedStr: "10302"},
// 				TestCase{inputsStr: []string{"179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137215", "100"}, expectedStr: "10400"},
// 				TestCase{inputsStr: []string{"179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137107", "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137105"}, expectedStr: "24"},
// 				TestCase{inputsStr: []string{"179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137105", "179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137107"}, expectedStr: "24"},
// 			}

// 			for _, test := range testcases {
// 				left := FromString(test.inputsStr[0])
// 				right := FromString(test.inputsStr[1])
// 				expected := left.MulModulo(&right, &prime)

// 				// expected := FromString(test.expectedStr)

// 				leftM := PrimeM.ToMont(&left)
// 				rightM := PrimeM.ToMont(&right)
// 				actM := leftM.MontMul(&rightM)
// 				act := actM.ToInt1024()

// 				Ω(act).Should(Equal(expected))

// 				// for i := 0; i < 1000; i++ {
// 				// 	expected = left.MulModulo(&right, &prime)
// 				// 	actM = leftM.MontMul(&rightM)
// 				// }

// 				oneM := PrimeM.ToMont(&one)
// 				expected = left.AddModulo(&one, &prime)
// 				actM = leftM.MontAdd(&oneM)
// 				act = actM.ToInt1024()

// 				Ω(act).Should(Equal(expected))
// 			}
// 		})
// 	})

// 	Context("when finding multiplicative inverse", func() {
// 		It("should return the right result for 1024 bit numbers", func() {
// 			twoM := PrimeM.ToMont(&two)

// 			invM := twoM.MontInv()
// 			inv := invM.ToInt1024()
// 			expInv := two.ModInverse(&prime)
// 			Ω(inv).Should(Equal(expInv))

// 		})
// 	})

// 	Context("when adding numbers", func() {
// 		It("should return the right result for 1024 bit numbers", func() {
// 			left := FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137107")
// 			leftM := PrimeM.ToMont(&left)
// 			oneM := PrimeM.ToMont(&one)
// 			expected := left.AddModulo(&one, &prime)
// 			actM := leftM.MontAdd(&oneM)
// 			act := actM.ToInt1024()

// 			Ω(act).Should(Equal(expected))
// 		})
// 	})

// 	Context("when subtracting numbers", func() {
// 		It("should return the right result for 1024 bit numbers", func() {
// 			left := FromString("0")
// 			leftM := PrimeM.ToMont(&left)
// 			oneM := PrimeM.ToMont(&one)
// 			expected := left.SubModulo(&one, &prime)
// 			// fmt.Println(leftM.String())
// 			actM := leftM.MontSub(&oneM)
// 			// fmt.Println(actM.String())
// 			act := actM.ToInt1024()

// 			Ω(act).Should(Equal(expected))
// 		})
// 	})
// })
