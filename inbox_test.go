package dark_test
//
//import (
//	"math/big"
//	"math/rand"
//	"sync"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	"github.com/republicprotocol/go-dark-network"
//	"github.com/republicprotocol/go-do"
//	"github.com/republicprotocol/go-order-compute"
//)
//
//const Number_Of_New_Result = 1000
//
//func newResult(i int) *compute.Final {
//	return &compute.Final{
//		ID:          []byte{uint8(i)},
//		BuyOrderID:  []byte{},
//		SellOrderID: []byte{},
//		FstCode:     big.NewInt(0),
//		SndCode:     big.NewInt(0),
//		Price:       big.NewInt(0),
//		MaxVolume:   big.NewInt(0),
//		MinVolume:   big.NewInt(0),
//	}
//}
//
//var _ = Describe("Inbox", func() {
//	var inbox *dark.Inbox
//
//	BeforeEach(func() {
//		inbox = dark.NewInbox()
//	})
//
//	Context("add/retrieve result to/from the box", func() {
//		It("should be able to sequentially add and retrieve result ", func() {
//			//  Add results
//			newResults := make([]*compute.Final, Number_Of_New_Result)
//			for index := range newResults {
//				newResults[index] = newResult(index)
//			}
//			do.CoForAll(newResults, func(i int) {
//				inbox.AddNewResult(newResults[i])
//			})
//
//			// Get all new results
//			allNewResults := inbox.GetAllNewResults()
//			Ω(len(allNewResults)).Should(Equal(Number_Of_New_Result))
//
//			// Get all results
//			wg := new(sync.WaitGroup)
//			wg.Add(Number_Of_New_Result)
//
//			for i := 0; i < Number_Of_New_Result; i++ {
//				go func() {
//					defer GinkgoRecover()
//
//					results := inbox.GetAllResults()
//					Ω(len(results)).Should(Equal(Number_Of_New_Result))
//					wg.Done()
//				}()
//			}
//			wg.Wait()
//		})
//	})
//
//	Context("simulate random functions concurrently", func() {
//		It("should handle concurrent calls properly", func() {
//			results := make([]*compute.Final, Number_Of_New_Result)
//			for index := range results {
//				results[index] = newResult(index)
//			}
//			wg := new(sync.WaitGroup)
//			wg.Add(Number_Of_New_Result)
//
//			operations := []string{"add", "all"}
//			for i := 0; i < Number_Of_New_Result; i++ {
//				go func(i int) {
//					defer GinkgoRecover()
//
//					operation := operations[rand.Intn(2)]
//					switch operation {
//					case "add":
//						result := newResult(i)
//						inbox.AddNewResult(result)
//					case "all":
//						_ = inbox.GetAllResults()
//					}
//					wg.Done()
//				}(i)
//			}
//			wg.Wait()
//		})
//	})
//
//})
