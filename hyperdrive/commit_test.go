package hyper_test

// import (
// 	"context"
// 	"sync"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/republicprotocol/republic-go/hyperdrive"
// )

// var _ = Describe("Commits", func() {

// 	blocks := NewSharedBlocks(1, 1)
// 	validator, _ := NewTestValidator(blocks, 4)

// 	Context("when processing commits", func() {

// 		It("should return errors on shutdown", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			commitChIn := make(chan Commit)

// 			_, _, _, errCh := ProcessCommit(ctx, commitChIn, validator)

// 			var wg sync.WaitGroup
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()

// 				for err := range errCh {
// 					Ω(err).Should(HaveOccurred())
// 					Ω(err).Should(Equal(context.Canceled))
// 				}
// 			}()

// 			cancel()
// 			wg.Wait()
// 		})

// 		It("should return commit after processing a threshold number of prepares", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			commitChIn := make(chan Commit)
// 			_, blockCh, _, errCh := ProcessCommit(ctx, commitChIn, validator)

// 			var wg sync.WaitGroup
// 			wg.Add(1)
// 			go func() {
// 				defer wg.Done()

// 				for {
// 					select {
// 					case err := <-errCh:
// 						Ω(err).Should(HaveOccurred())
// 						Ω(err).Should(Equal(context.Canceled))
// 						return
// 					case block, ok := <-blockCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(block).Should(Equal(Block{
// 							Tuples{},
// 							Signature("Proposer"),
// 						}))
// 						cancel()
// 					}
// 				}
// 			}()

// 			for i := uint8(0); i < validator.Threshold(); i++ {
// 				commit := Commit{
// 					ThresholdSignature: ThresholdSignature("Threshold_BLS"),
// 					Signature:          Signature("Proposer"),
// 					Block: Block{
// 						Tuples{},
// 						Signature("Proposer"),
// 					},
// 					Rank:   Rank(1),
// 					Height: Height(1),
// 				}
// 				commitChIn <- commit
// 			}
// 			wg.Wait()
// 		})
// 	})
// })
