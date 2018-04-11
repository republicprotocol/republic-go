package hyper_test

// import (
// 	"context"
// 	"sync"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// 	. "github.com/republicprotocol/republic-go/hyperdrive"
// )

// var _ = Describe("Processors", func() {

// 	signer, _ := NewTestSigner()
// 	proposer, _ := NewTestSigner()
// 	blocks := NewSharedBlocks(1, 1)
// 	validator, _ := NewTestValidator(blocks)

// 	Context("when processing proposals", func() {

// 		It("should return errors on shutdown", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			proposalChIn := make(chan Proposal)

// 			_, _, errCh := ProcessProposal(ctx, proposalChIn, signer, validator)

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

// 		It("should return prepares after processing a valid proposal", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			proposalChIn := make(chan Proposal)

// 			proposal := Proposal{
// 				Signature(proposer.ID().String()),
// 				Block{
// 					Tuples{},
// 					Signature(proposer.ID().String()),
// 				},
// 				Rank(1),
// 				Height(1),
// 			}
// 			prepCh, _, errCh := ProcessProposal(ctx, proposalChIn, signer, validator)

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
// 					case prepare, ok := <-prepCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(prepare.Rank).Should(Equal(proposal.Rank))
// 						Ω(prepare.Height).Should(Equal(proposal.Height))
// 						Ω(prepare.Block).Should(Equal(proposal.Block))
// 						Ω(prepare.Signature).Should(Equal(signer.Sign()))
// 						cancel()
// 					}
// 				}
// 			}()

// 			proposalChIn <- proposal
// 			wg.Wait()
// 		})

// 		It("should return faults after processing an invalid proposal", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			proposalChIn := make(chan Proposal)

// 			proposal := Proposal{
// 				Signature(proposer.ID().String()),
// 				Block{
// 					Tuples{},
// 					Signature(proposer.ID().String()),
// 				},
// 				Rank(2),
// 				Height(1),
// 			}
// 			_, faultCh, errCh := ProcessProposal(ctx, proposalChIn, signer, validator)

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
// 					case fault, ok := <-faultCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(fault.Rank).Should(Equal(proposal.Rank))
// 						Ω(fault.Height).Should(Equal(proposal.Height))
// 						Ω(fault.Signature).Should(Equal(signer.Sign()))
// 						cancel()
// 					}
// 				}
// 			}()

// 			proposalChIn <- proposal
// 			wg.Wait()
// 		})

// 		It("should not return a prepare after processing an invalid proposal", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			proposalChIn := make(chan Proposal)

// 			proposal := Proposal{
// 				Signature(proposer.ID().String()),
// 				Block{
// 					Tuples{},
// 					Signature(proposer.ID().String()),
// 				},
// 				Rank(2),
// 				Height(1),
// 			}
// 			prepCh, faultCh, errCh := ProcessProposal(ctx, proposalChIn, signer, validator)

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
// 					case prepare, ok := <-prepCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(prepare).Should(Not(HaveOccurred()))
// 						cancel()
// 					case fault, ok := <-faultCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(fault.Rank).Should(Equal(proposal.Rank))
// 						Ω(fault.Height).Should(Equal(proposal.Height))
// 						Ω(fault.Signature).Should(Equal(signer.Sign()))
// 						cancel()
// 					}
// 				}
// 			}()

// 			proposalChIn <- proposal
// 			wg.Wait()
// 		})

// 		It("should not return a fault after processing a valid proposal", func() {
// 			ctx, cancel := context.WithCancel(context.Background())
// 			proposalChIn := make(chan Proposal)

// 			proposal := Proposal{
// 				Signature(proposer.ID().String()),
// 				Block{
// 					Tuples{},
// 					Signature(proposer.ID().String()),
// 				},
// 				Rank(1),
// 				Height(1),
// 			}
// 			prepCh, faultCh, errCh := ProcessProposal(ctx, proposalChIn, validator)

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
// 					case prepare, ok := <-prepCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(prepare.Rank).Should(Equal(proposal.Rank))
// 						Ω(prepare.Height).Should(Equal(proposal.Height))
// 						Ω(prepare.Block).Should(Equal(proposal.Block))
// 						Ω(prepare.Signature).Should(Equal(signer.Sign()))
// 						cancel()
// 					case fault, ok := <-faultCh:
// 						if !ok {
// 							return
// 						}
// 						Ω(fault).Should(Not(HaveOccurred()))
// 						cancel()
// 					}
// 				}
// 			}()

// 			proposalChIn <- proposal
// 			wg.Wait()
// 		})
// 	})
// })
