package xing_test

import (
	"log"
	"math/big"
	"math/rand"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-order-compute"
	"github.com/republicprotocol/go-xing"
)

const (
	Number_Of_New_Result = 1000
)

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
			results := make([]*compute.Result, Number_Of_New_Result)
			for i := range results {
				results[i] = newResult(i)
			}
			do.CoForAll(results, func(i int) {
				inbox.AddNewResult(results[i])
			})

			// Get all results
			wg := new(sync.WaitGroup)
			wg.Add(Number_Of_New_Result)
			for i := 0; i < Number_Of_New_Result; i++ {
				go func() {
					defer GinkgoRecover()

					results := inbox.GetAllResults()
					wg.Done()
					Ω(len(results)).Should(Equal(1000))
				}()
			}
			wg.Wait()

			// Get new results
			wg = new(sync.WaitGroup)
			wg.Add(Number_Of_New_Result)
			for i := 0; i < Number_Of_New_Result; i++ {
				go func() {
					defer GinkgoRecover()

					_ = inbox.GetNewResult()
					wg.Done()
				}()
			}
			wg.Wait()
		})
	})

	FContext("simulate random functions concurrently", func() {
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
			operations := []string{"add", "get" , "all"}
			for i:= 0; i <Number_Of_New_Result; i ++{
				go func() {
					operation := operations[rand.Intn(3)]
					switch operation {
					case "add":
						result:= newResult(i)
						inbox.AddNewResult(result)
					case "get":
						_ = inbox.GetNewResult()
					case "all":
						_ = inbox.GetAllResults()
					}
					wg.Done()
				}()
			}
			wg.Wait()
			log.Println(len(inbox.GetAllResults()))
		})
	})

})
