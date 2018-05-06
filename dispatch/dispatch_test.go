package dispatch_test

import (
	"reflect"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Dispatch Package", func() {

	Context("Wait", func() {

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
				sigCh1 <- struct{}{}
				signal1 = false
				close(sigCh1)

				sigCh2 <- struct{}{}
				signal2 = false
				close(sigCh2)

				sigCh3 <- struct{}{}
				signal3 = false
				close(sigCh3)

				sigCh4 <- struct{}{}
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

	Context("Close", func() {

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

	Context("Split", func() {

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

	Context("Merge", func() {

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

		It("should panic for invalid arguments", func() {
			in := 10
			out := false

			Ω(func() { Merge(in, out) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := make(chan int)
			out := [5]int{1, 2, 3, 4, 5}

			Ω(func() { Merge(in, out) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := make(chan int)
			out := 10

			Ω(func() { Merge(in, out) }).Should(Panic())
		})

	})

	Context("Send", func() {
		It("should send a message to a channel", func() {
			in := make(chan int)
			msg := 1

			go Send(in, reflect.ValueOf(msg))

			Ω(<-in).Should(Equal(msg))
		})

		It("should send a message to an array of channel", func() {
			in := make([]chan int, 10)
			msg := 1

			for i := 0; i < 10; i++ {
				in[i] = make(chan int)
			}

			go Send(in, reflect.ValueOf(msg))

			for i := 0; i < 10; i++ {
				Ω(<-in[i]).Should(Equal(msg))
			}
		})

		It("should panic if the message and channel have different types", func() {
			in := make(chan int)
			msg := 1

			go func() {
				time.Sleep(time.Second)
				close(in)
			}()

			Ω(func() { Send(in, reflect.ValueOf(msg)) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := 2
			msg := 1

			Ω(func() { Send(in, reflect.ValueOf(msg)) }).Should(Panic())
		})

		It("should panic for invalid arguments", func() {
			in := []int{1, 2, 3, 4, 5}
			msg := 1

			Ω(func() { Send(in, reflect.ValueOf(msg)) }).Should(Panic())
		})

	})

	Context("Pipe", func() {
		It("should be able to pipe from one channel to another", func() {
			doneCh := make(chan struct{})
			inCh := make(chan int)
			outCh := make(chan int)

			var wg sync.WaitGroup

			wg.Add(1)
			go func() {
				defer wg.Done()
				Pipe(doneCh, inCh, outCh)
			}()

			for i := 0; i < 10; i++ {
				inCh <- i
				Ω(<-outCh).Should(Equal(i))
			}

			close(doneCh)
			wg.Wait()

			close(inCh)
			close(outCh)
		})

		It("should panic for an invalid producer", func() {
			doneCh := make(chan struct{})
			outCh := make(chan int)
			Ω(func() { Pipe(doneCh, 10, outCh) }).Should(Panic())
		})

		It("should panic for an invalid consumer", func() {
			doneCh := make(chan struct{})
			inCh := make(chan int)
			Ω(func() { Pipe(doneCh, inCh, 10) }).Should(Panic())
		})
	})
})
