package dispatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Dispatch Package", func() {

	FContext("Wait", func() {

		It("should wait for a single signal channels", func() {
			sigCh := make(chan struct{})
			signal := true

			go func() {
				signal = false
				close(sigCh)
			}()

			Wait(sigCh)

			Ω(signal).Should(BeFalse())
		})

		It("should wait for multiple signal channels", func() {
			sigCh1 := make(chan struct{})
			sigCh2 := make(chan struct{})
			sigCh3 := make(chan struct{})
			sigCh4 := make(chan struct{})

			signal1 := true
			signal2 := true
			signal3 := true
			signal4 := true

			go func() {
				signal1 = false
				close(sigCh1)

				signal2 = false
				close(sigCh2)

				signal3 = false
				close(sigCh3)

				signal4 = false
				close(sigCh4)
			}()

			Wait(sigCh1, sigCh2, sigCh3, sigCh4)

			Ω(signal1).Should(BeFalse())
			Ω(signal2).Should(BeFalse())
			Ω(signal3).Should(BeFalse())
			Ω(signal4).Should(BeFalse())
		})

	})

	FContext("Close", func() {

		It("should close a single channel", func() {
			ch := make(chan int)
			Close(ch)
			_, ok := <-ch
			Ω(ok).Should(BeFalse())
		})

		It("should close a single channel", func() {
			ch1 := make(chan int)
			ch2 := make(chan int)
			ch3 := make(chan int)
			ch4 := make(chan int)

			Close(ch1, ch2, ch3, ch4)

			_, ok1 := <-ch1
			_, ok2 := <-ch2
			_, ok3 := <-ch3
			_, ok4 := <-ch4

			Ω(ok1).Should(BeFalse())
			Ω(ok2).Should(BeFalse())
			Ω(ok3).Should(BeFalse())
			Ω(ok4).Should(BeFalse())
		})

	})

	FContext("Split", func() {

		It("should split the channel into an array of channels", func() {
			inCh := make(chan int)
			outChs := make([]chan int, 100)
			for i := 0; i < 100; i++ {
				outChs[i] = make(chan int)
			}

			go Split(inCh, outChs)

			go func() {
				defer close(inCh)
				inCh <- 1729
			}()

			for _, ch := range outChs {
				i := <-ch
				Ω(i).Should(Equal(1729))
			}
		})

		It("should split the channel into multiple channels", func() {
			inCh := make(chan int)
			outCh1 := make(chan int)
			outCh2 := make(chan int)
			outCh3 := make(chan int)

			go Split(inCh, outCh1, outCh2, outCh3)

			go func() {
				defer close(inCh)
				inCh <- 1729
			}()

			o1 := <-outCh1
			Ω(o1).Should(Equal(1729))

			o2 := <-outCh2
			Ω(o2).Should(Equal(1729))

			o3 := <-outCh3
			Ω(o3).Should(Equal(1729))

		})

		It("should panic when the output channel is of a different type", func() {
			inCh := make(chan int)
			outCh := make(chan struct{})

			go func() {
				defer close(inCh)
				inCh <- 1729
			}()

			Ω(func() { Split(inCh, outCh) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := 10
			out := false

			Ω(func() { Split(in, out) }).Should(Panic())
		})
	})

	FContext("Merge", func() {

		It("should merge an array of channels into a channel", func() {

			outCh := make(chan int)
			inChs := make([]chan int, 100)

			for i := 0; i < 100; i++ {
				inChs[i] = make(chan int)
				go func(i int) {
					defer close(inChs[i])
					inChs[i] <- i
				}(i)
			}

			go Merge(outCh, inChs)

			for i := 0; i < 100; i++ {
				_ = <-outCh

			}
		})

		It("should merge multiple channels into a channel", func() {
			outCh := make(chan int)
			inCh1 := make(chan int)
			inCh2 := make(chan int)
			inCh3 := make(chan int)

			go Merge(outCh, inCh1, inCh2, inCh3)

			go func() {
				defer close(inCh1)
				defer close(inCh2)
				defer close(inCh3)
				inCh1 <- 1
				inCh2 <- 2
				inCh3 <- 3
			}()

			_ = <-outCh
			_ = <-outCh
			_ = <-outCh

			close(outCh)
		})

		It("should panic when the output channels are of different types", func() {
			inCh := make(chan int)
			outCh := make(chan struct{})

			go func() {
				defer close(inCh)
				inCh <- 1729
			}()

			Ω(func() { Merge(outCh, inCh) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := 10
			out := false

			Ω(func() { Merge(in, out) }).Should(Panic())
		})
	})
})
