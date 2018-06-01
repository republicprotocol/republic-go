package dispatch_test

import (
	"sync"
	"sync/atomic"
	"time"

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
			broadcaster := NewBroadcaster()
			broadcaster.Close()

			done := make(chan struct{})
			lis := broadcaster.Listen(done)
			select {
			case _, ok := <-lis:
				Expect(ok).Should(BeFalse())
			}
		})

		It("should not block when shutting down under heavy usage", func() {
			broadcaster := NewBroadcaster()

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				CoForAll(int(MaxListeners), func(i int) {
					done := make(chan struct{})
					for range broadcaster.Listen(done) {
					}
				})
			}()
			time.Sleep(time.Second)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer broadcaster.Close()

				CoForAll(int(MaxListeners), func(i int) {
					done := make(chan struct{})
					ch := make(chan interface{})
					defer close(ch)
					go broadcaster.Broadcast(done, ch)
					ch <- struct{}{}
				})
			}()

			wg.Wait()
		})

	})

	Context("when broadcasting", func() {

		It("should restrict the maximum number of listeners", func() {
			broadcaster := NewBroadcaster()
			listeners := make([]<-chan interface{}, 2*MaxListeners)
			closed := int32(0)
			CoForAll(int(2*MaxListeners), func(i int) {
				done := make(chan struct{})
				lis := broadcaster.Listen(done)
				listeners[i] = lis
			})
			CoForAll(int(2*MaxListeners), func(i int) {
				select {
				case _, ok := <-listeners[i]:
					if !ok {
						atomic.AddInt32(&closed, 1)
					}
				default:
				}
			})
			Expect(closed).Should(Equal(MaxListeners))
		})

		It("should send message from one broadcast to many listeners", func(done Done) {
			defer close(done)
			broadcaster := NewBroadcaster()

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				n := int64(0)
				CoForAll(int(MaxListeners), func(i int) {
					done := make(chan struct{})
					for range broadcaster.Listen(done) {
						atomic.AddInt64(&n, 1)
					}
				})
				Expect(n).Should(Equal(int64(MaxListeners)))
			}()
			time.Sleep(time.Second)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer broadcaster.Close()

				done := make(chan struct{})
				ch := make(chan interface{})
				defer close(ch)
				go broadcaster.Broadcast(done, ch)
				ch <- struct{}{}
				time.Sleep(time.Second)
			}()

			wg.Wait()
		}, 30 /* 30 second timeout */)

		It("should send messages from many broadcasts to many listeners", func(done Done) {
			defer close(done)
			broadcaster := NewBroadcaster()

			var wg sync.WaitGroup
			wg.Add(2)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()

				n := int64(0)
				CoForAll(int(MaxListeners), func(i int) {
					done := make(chan struct{})
					for range broadcaster.Listen(done) {
						atomic.AddInt64(&n, 1)
					}
				})
				Expect(n).Should(Equal(int64(MaxListeners * MaxListeners)))
			}()
			time.Sleep(time.Second)

			go func() {
				defer GinkgoRecover()
				defer wg.Done()
				defer broadcaster.Close()

				CoForAll(int(MaxListeners), func(i int) {
					done := make(chan struct{})
					ch := make(chan interface{})
					defer close(ch)
					go broadcaster.Broadcast(done, ch)
					ch <- struct{}{}
				})
				time.Sleep(time.Second * 2)
			}()

			wg.Wait()
		}, 60 /* 1 minute timeout */)

	})

})
