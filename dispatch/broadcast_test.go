package dispatch_test

import (
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Broadcaster", func() {

	Context("when shutting down", func() {

		It("should not block existing broadcasts after shutting down", func(done Done) {
			signal := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(done)
				Eventually(signal, 10).Should(BeClosed())
			}()

			var wg sync.WaitGroup
			broadcaster := NewBroadcaster()
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					done := make(chan struct{})
					ch := make(chan interface{})

					CoBegin(func() {
						defer close(done)
						broadcaster.Broadcast(done, ch)
					}, func() {
						for j := 0; j < 1000; j++ {
							select {
							case <-done:
							case ch <- j:
							}
						}
					})
				}()
			}
			broadcaster.Close()
			wg.Wait()

			close(signal)
		}, 10 /* 10 second timeout */)

		It("should not block new broadcasts after shutting down", func(done Done) {
			signal := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(done)
				Eventually(signal, 10).Should(BeClosed())
			}()

			var wg sync.WaitGroup
			broadcaster := NewBroadcaster()
			broadcaster.Close()
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()

					done := make(chan struct{})
					ch := make(chan interface{})

					CoBegin(func() {
						defer close(done)
						broadcaster.Broadcast(done, ch)
					}, func() {
						for j := 0; j < 1000; j++ {
							select {
							case <-done:
							case ch <- j:
							}
						}
					})
				}()
			}
			wg.Wait()

			close(signal)
		}, 10 /* 10 second timeout */)

		It("should not block existing listeners after shutting down", func() {

		})

		It("should not block new listeners after shutting down", func() {

		})

		It("should not block when shutting down under heavy usage", func() {

		})

	})

	Context("when broadcasting", func() {

		It("should send message from one broadcast to many listeners", func() {

		})

		It("should send messages from many broadcasts to one listener", func() {

		})

		It("should send messages from many broadcasts to many listeners", func() {

		})

	})

})
