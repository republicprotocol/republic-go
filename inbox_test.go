package xing_test

import (
	"math/big"
	"math/rand"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-xing"
)

const Number_Of_New_Result = 1000

func newResult(i int) *compute.Result {
	return &compute.Result{
		ID:          []byte{uint8(i)},
		BuyOrderID:  []byte{},
		SellOrderID: []byte{},
		FstCode:     big.NewInt(0),
		SndCode:     big.NewInt(0),
		Price:       big.NewInt(0),
		MaxVolume:   big.NewInt(0),
		MinVolume:   big.NewInt(0),
	}
}

var _ = Describe("Inbox", func() {
	var inbox *xing.Inbox

	BeforeEach(func() {
		inbox = xing.NewInbox()
	})

	Context("add/retrieve result to/from the box", func() {
		It("should be able to sequentially add and retrieve result ", func() {
			//  Add results
			newResults := make([]*compute.Result, Number_Of_New_Result)
			for i := range newResults {
				newResults[i] = newResult(i)
			}
			do.CoForAll(newResults, func(i int) {
				inbox.AddNewResult(newResults[i])
			})

			// Get all results
			wg := new(sync.WaitGroup)
			wg.Add(Number_Of_New_Result)
			for i := 0; i < Number_Of_New_Result; i++ {
				go func() {
					defer GinkgoRecover()

					results := inbox.GetAllResults()
					Ω(len(results)).Should(Equal(1000))
					wg.Done()
				}()
			}
			wg.Wait()

			results := inbox.GetAllNewResults()
			Ω(len(results)).Should(BeNumerically(">", 0))
		})
	})

	Context("simulate random functions concurrently", func() {
		It("should handle concurrent calls properly", func() {
			results := make([]*compute.Result, Number_Of_New_Result)
			for i := range results {
				results[i] = newResult(i)
			}
			do.CoForAll(results, func(i int) {
				inbox.AddNewResult(results[i])
			})

			wg := new(sync.WaitGroup)
			wg.Add(Number_Of_New_Result)

			operations := []string{"add", "all"}
			for i:= 0; i <Number_Of_New_Result; i ++{
				go func(i int) {
					defer GinkgoRecover()

					operation := operations[rand.Intn(2)]
					switch operation {
					case "add":
						result:= newResult(i)
						inbox.AddNewResult(result)
					case "all":
						_ = inbox.GetAllResults()
					}
					wg.Done()
				}(i)
			}
			go func() {
				defer GinkgoRecover()
				results := inbox.GetAllNewResults()
				Ω(len(results)).Should(BeNumerically(">", 0))
			}()
			wg.Wait()
		})
	})

})
