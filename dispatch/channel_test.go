package dispatch_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Channels", func() {

	Context("when merging a channel of channels", func() {

		It("should merge a single input channel", func() {
			done := make(chan struct{})
			ins := make(chan (chan int))
			out := make(chan int)

			go Merge(done, ins, out)
			go func() {
				defer close(ins)
				in := make(chan int)
				ins <- in
				for i := 0; i < 10; i++ {
					in <- i
				}
			}()
			for i := 0; i < 10; i++ {
				Expect(<-out).Should(Equal(i))
			}

			close(done)
			time.Sleep(time.Millisecond)
		})

		It("should merge multiple input channels", func() {
			done := make(chan struct{})
			ins := make(chan (chan int))
			out := make(chan int)

			go Merge(done, ins, out)
			go func() {
				defer close(ins)
				for n := 0; n < 10; n++ {
					in := make(chan int)
					ins <- in
					for i := 0; i < 10; i++ {
						in <- i
					}
				}
			}()
			js := map[int]int{}
			for i := 0; i < 10*10; i++ {
				j := <-out
				Expect(j).Should(BeNumerically(">=", 0))
				Expect(j).Should(BeNumerically("<", 10))
				js[j]++
			}
			for i := 0; i < 10; i++ {
				Expect(js[i]).Should(Equal(10))
			}

			close(done)
			time.Sleep(time.Millisecond)
		})

		It("should panic for invalid input types", func() {
			Expect(func() {
				Merge(make(chan struct{}), make(chan (chan float32)), make(chan int))
			}).Should(Panic())

			Expect(func() {
				Merge(make(chan struct{}), make(chan int), make(chan int))
			}).Should(Panic())

			Expect(func() {
				Merge(make(chan struct{}), make(chan float32), make(chan int))
			}).Should(Panic())

			Expect(func() {
				Merge(make(chan struct{}), make(chan (chan float32)), 0)
			}).Should(Panic())

			Expect(func() {
				Merge(make(chan struct{}), 0, make(chan int))
			}).Should(Panic())

			Expect(func() {
				Merge(make(chan struct{}), 0, 0)
			}).Should(Panic())
		})
	})

	Context("when forwarding a channel", func() {

		It("should forward an input channel", func() {
			done := make(chan struct{})
			in := make(chan int)
			out := make(chan int)

			go Forward(done, in, out)
			go func() {
				defer close(in)
				for i := 0; i < 10; i++ {
					in <- i
				}
			}()
			for i := 0; i < 10; i++ {
				Expect(<-out).Should(Equal(i))
			}

			close(done)
			time.Sleep(time.Millisecond)
		})

		It("should panic for invalid input types", func() {
			Expect(func() {
				Forward(make(chan struct{}), make(chan (chan float32)), make(chan int))
			}).Should(Panic())

			Expect(func() {
				Forward(make(chan struct{}), make(chan float32), make(chan int))
			}).Should(Panic())

			Expect(func() {
				Forward(make(chan struct{}), make(chan (chan float32)), 0)
			}).Should(Panic())

			Expect(func() {
				Forward(make(chan struct{}), 0, make(chan int))
			}).Should(Panic())

			Expect(func() {
				Forward(make(chan struct{}), 0, 0)
			}).Should(Panic())
		})
	})
})
