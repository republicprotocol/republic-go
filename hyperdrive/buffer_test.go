package hyper_test

// import (
// 	"log"
// 	"math/rand"
// 	"sync"
// 	"time"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/republicprotocol/republic-go/hyperdrive"
// )

// var _ = Describe("Buffer", func() {

// 	Context("Proposals", func() {

// 		It("should only return proposals of current static height", func() {
// 			chanSetIn := EmptyChannelSet()
// 			sb := NewSharedBlocks(1, 1)
// 			chanSetOut := ProcessBuffer(chanSetIn, &sb)

// 			go func() {
// 				for {
// 					select {
// 					case proposal, ok := <-chanSetOut.Proposal:
// 						if !ok {
// 							return
// 						}
// 						Ω(proposal.Height).Should(Equal(sb.Height))
// 					}
// 				}
// 			}()

// 			for i := 0; i < 100; i++ {
// 				h := rand.Intn(4) + 1
// 				chanSetIn.Proposal <- Proposal{
// 					Height: Height(h),
// 					Rank:   Rank(1),
// 					Block:  Block{},
// 				}
// 			}

// 			chanSetIn.Close()

// 		})

// 		It("should only return proposals of current dynamic height which changes every second", func() {
// 			chanSetIn := EmptyChannelSet()
// 			counter := map[Height]uint64{}

// 			sb := NewSharedBlocks(1, 1)

// 			chanSetOut := ProcessBuffer(chanSetIn, &sb)
// 			var wg sync.WaitGroup
// 			go func() {
// 				defer GinkgoRecover()
// 				defer log.Println("Channel closed")
// 				for {
// 					select {
// 					case proposal, ok := <-chanSetOut.Prepare:
// 						if !ok {
// 							return
// 						}
// 						counter[proposal.Height]++
// 						Ω(proposal.Height).Should(Equal(sb.Height))
// 						log.Println("Counter", counter[1]+counter[2]+counter[3]+counter[4]+counter[5])
// 					}
// 				}
// 			}()

// 			for i := 0; i < 1000; i++ {
// 				h := rand.Intn(5) + 1
// 				chanSetIn.Proposal <- Proposal{
// 					Height: Height(h),
// 					Rank:   Rank(1),
// 					Block:  Block{},
// 				}
// 			}
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()
// 				defer GinkgoRecover()
// 				defer time.Sleep(1 * time.Second)
// 				for i := 1; i < 5; i++ {
// 					time.Sleep(1 * time.Second)
// 					sb.IncrementHeight()
// 				}
// 			}()

// 			chanSetIn.Close()
// 			wg.Wait()

// 			Ω(uint64(1000)).Should(Equal(counter[1] + counter[2] + counter[3] + counter[4] + counter[5]))
// 		})

// 	})
// })
