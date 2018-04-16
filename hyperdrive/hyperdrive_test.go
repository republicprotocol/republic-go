package hyper_test

// import (
// 	"context"
// 	"log"
// 	"strconv"
// 	"sync"

// 	. "github.com/onsi/ginkgo"
// 	. "github.com/republicprotocol/republic-go/hyperdrive"
// 	"golang.org/x/crypto/sha3"
// )

// var _ = Describe("Hyperdrive", func() {

// 	threshold := uint8(4)
// 	commanderCount := uint8(6)

// 	Context("Hyperdrive", func() {

// 		It("Achieves consensus on a block over 240 commanders with 75% threshold", func() {
// 			ctx, cancel := context.WithCancel(context.Background())

// 			// Network
// 			ingress := make([]ChannelSet, commanderCount)
// 			egress := make([]ChannelSet, commanderCount)

// 			for i := uint8(0); i < commanderCount; i++ {
// 				ingress[i] = EmptyChannelSet(ctx, commanderCount)
// 				egress[i] = EmptyChannelSet(ctx, commanderCount)
// 			}

// 			for i := uint8(0); i < commanderCount; i++ {
// 				go egress[i].Split(ingress)
// 			}

// 			log.Println("Network initialized.... ")

// 			// Hyperdrive

// 			// Initialize replicas
// 			replicas := make([]Replica, commanderCount)
// 			for i := uint8(0); i < commanderCount; i++ {
// 				blocks := NewSharedBlocks(0, 0)
// 				validator, _ := NewTestValidator(blocks, threshold)
// 				replicas[i] = NewReplica(ctx, validator, ingress[i])
// 			}

// 			// Run replicas
// 			for i := uint8(0); i < commanderCount; i++ {
// 				go egress[i].Copy(replicas[i].Run())
// 			}

// 			log.Println("Starting the hyperdrive.... ")

// 			// Broadcast proposal to all the nodes
// 			proposal := Proposal{
// 				Signature("Proposal"),
// 				Block{
// 					Tuples{},
// 					Signature("Proposal"),
// 				},
// 				Rank(0),
// 				0,
// 			}

// 			for i := 0; i < len(replicas); i++ {
// 				ingress[i].Proposal <- proposal
// 			}

// 			log.Println("Broadcasted the proposals")

// 			// Wait for the blocks from all the nodes
// 			var wg sync.WaitGroup
// 			for i := uint8(0); i < commanderCount; i++ {
// 				wg.Add(1)
// 				go func(i uint8) {
// 					defer wg.Done()
// 					_ = <-egress[i].Block
// 					// log.Println("Block recieved on", i)
// 				}(i)
// 			}

// 			log.Println("Waiting for the blocks")
// 			wg.Wait()
// 			log.Println("Success!!!!!")
// 			cancel()

// 		})

// 		It("Achieves consensus 50 blocks over 240 commanders with 2/3 threshold", func() {
// 			ctx, cancel := context.WithCancel(context.Background())

// 			Blocks := 5

// 			// Network
// 			ingress := make([]ChannelSet, commanderCount)
// 			egress := make([]ChannelSet, commanderCount)

// 			for i := uint8(0); i < commanderCount; i++ {
// 				ingress[i] = EmptyChannelSet(ctx, commanderCount)
// 				egress[i] = EmptyChannelSet(ctx, commanderCount)
// 			}

// 			for i := uint8(0); i < commanderCount; i++ {
// 				go egress[i].Split(ingress)
// 			}

// 			log.Println("Network initialized.... ")

// 			// Hyperdrive

// 			// Initialize replicas
// 			replicas := make([]Replica, commanderCount)
// 			for i := uint8(0); i < commanderCount; i++ {
// 				blocks := NewSharedBlocks(0, 0)
// 				validator, _ := NewTestValidator(blocks, threshold)
// 				replicas[i] = NewReplica(ctx, validator, ingress[i])
// 			}

// 			// Run replicas
// 			for i := uint8(0); i < commanderCount; i++ {
// 				go egress[i].Copy(replicas[i].Run())
// 			}

// 			log.Println("Starting the hyperdrive.... ")

// 			// Broadcast proposal to all the nodes
// 			go func() {
// 				defer log.Println("Broadcasted the proposals")
// 				for i := 0; i < Blocks; i++ {
// 					block := Block{
// 						Tuples{
// 							Tuple{
// 								ID: sha3.Sum256([]byte(strconv.Itoa(i))),
// 							},
// 						},
// 						Signature("Proposal"),
// 					}
// 					for j := 0; j < len(replicas); j++ {
// 						ingress[j].Proposal <- Proposal{
// 							Signature("Proposal"),
// 							block,
// 							Rank(0),
// 							uint64(i),
// 						}
// 					}
// 				}
// 			}()

// 			// Wait for the blocks from all the nodes

// 			var wg sync.WaitGroup
// 			for i := uint8(0); i < commanderCount; i++ {
// 				wg.Add(1)
// 				go func(i uint8) {
// 					defer wg.Done()
// 					for j := 0; j < Blocks; j++ {
// 						_ = <-egress[i].Block
// 						log.Println("Received blocks no", j, "on", i)
// 					}
// 				}(i)
// 			}

// 			log.Println("Waiting for the blocks")
// 			wg.Wait()
// 			// time.Sleep(20 * time.Second)
// 			log.Println("Success!!!!!")
// 			cancel()

// 		})
// 	})
// })
