package hyper_test

import (
	"context"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Channel Set", func() {

	Context("Channel Set", func() {

		It("should be able to create an empty channel set", func() {
			emptyChanSet := EmptyChannelSet(100)
			Ω(emptyChanSet).Should(Not(BeNil()))
		})

		It("should be able to create, write, read and close a channel set from channels", func() {
			chanSet := EmptyChannelSet(100)
			chanSet.Proposal <- Proposal{
				Height: 1,
			}
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				for {
					select {
					case proposal, ok := <-chanSet.Proposal:
						if !ok {
							return
						}
						Ω(uint64(1)).Should(Equal(proposal.Height))
						wg.Done()
					}
				}
			}()
			time.Sleep(10 * time.Microsecond)
			chanSet.Close()
			wg.Wait()
		})

		It("should be able to broadcast a channel set to multiple channel sets", func() {
			chanSet := EmptyChannelSet(100)
			chanSet.Proposal <- Proposal{
				Height: 1,
			}
			var wg sync.WaitGroup
			chanSets := make([]ChannelSet, 100)

			for i := 0; i < 100; i++ {
				chanSets[i] = EmptyChannelSet(100)
			}
			ctx, cancel := context.WithCancel(context.Background())
			go chanSet.Split(ctx, chanSets)

			wg.Add(100)
			for i := 0; i < 100; i++ {
				go func(i int) {
					for {
						select {
						case proposal, ok := <-chanSets[i].Proposal:
							if !ok {
								return
							}
							Ω(uint64(1)).Should(Equal(proposal.Height))
							wg.Done()
							return
						}
					}
				}(i)
			}

			wg.Wait()
			cancel()
			chanSet.Close()
		})

		It("should be able to pipe a channel set into a different channel set", func() {
			chanSet := EmptyChannelSet(100)
			chanSet.Proposal <- Proposal{
				Height: 1,
			}
			chanSet2 := EmptyChannelSet(100)
			ctx, cancel := context.WithCancel(context.Background())
			go chanSet2.Copy(ctx, chanSet)
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				for {
					select {
					case proposal, ok := <-chanSet2.Proposal:
						if !ok {
							return
						}
						Ω(uint64(1)).Should(Equal(proposal.Height))
						wg.Done()
						return
					}
				}
			}()
			wg.Wait()
			cancel()
			chanSet.Close()
			chanSet2.Close()
		})
	})
})
