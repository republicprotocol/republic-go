package hyper_test

import (
	"context"
	"strconv"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Preparations", func() {

	blocks := NewSharedBlocks(1, 1)
	validator, _ := NewTestValidator(blocks, uint8(4))

	Context("when processing prepares", func() {

		It("should return errors on shutdown", func() {
			ctx, cancel := context.WithCancel(context.Background())
			prepareChIn := make(chan Prepare)

			_, _, errCh := ProcessPreparation(ctx, prepareChIn, validator)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()

				for err := range errCh {
					Ω(err).Should(HaveOccurred())
					Ω(err).Should(Equal(context.Canceled))
				}
			}()

			cancel()
			wg.Wait()
		})

		It("should return commit after processing a threshold number of prepares", func() {
			ctx, cancel := context.WithCancel(context.Background())
			prepareChIn := make(chan Prepare)
			commitCh, _, errCh := ProcessPreparation(ctx, prepareChIn, validator)

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()

				for {
					select {
					case err := <-errCh:
						Ω(err).Should(HaveOccurred())
						Ω(err).Should(Equal(context.Canceled))
						return
					case commit, ok := <-commitCh:
						if !ok {
							return
						}
						Ω(commit.Rank).Should(Equal(Rank(1)))
						Ω(commit.Height).Should(Equal(uint64(1)))
						Ω(commit.Block).Should(Equal(Block{
							Tuples{},
							Signature("Proposer"),
						}))
						Ω(commit.ThresholdSignature).Should(Equal(ThresholdSignature("Threshold_BLS")))
						cancel()
					}
				}
			}()

			for i := uint8(0); i < validator.Threshold(); i++ {
				prepare := Prepare{
					Signature("Signature of " + strconv.Itoa(int(i))),
					Block{
						Tuples{},
						Signature("Proposer"),
					},
					Rank(1),
					uint64(1),
				}
				prepareChIn <- prepare
			}
			wg.Wait()
		})
	})
})
