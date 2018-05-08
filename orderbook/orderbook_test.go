package orderbook_test

import (
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

const (
	MaxConnections = 128
	NumberOfOrders = 100
	)

var _ = Describe("order book", func() {
	Context("order status change event", func() {
		var book orderbook.Orderbook
		var orders [NumberOfOrders]order.ID
		var err error

		BeforeEach(func() {
			book, err  = orderbook.NewOrderbook()
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			book.Close()

			err:= ClearDB()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should be able to open an order", func() {
			for i:=0 ; i < NumberOfOrders; i ++ {
				ord := NewOrder(order.ParityBuy, uint(i))
				orders[i] = ord.ID
				fragments, err  := ord.Split(1,1, &smpc.Prime)
				Expect(err).ShouldNot(HaveOccurred())

				err = book.Open( *fragments[0])
				Expect(err).ShouldNot(HaveOccurred())
			}

			for _, id := range orders{
				status, err := book.Status(id)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(status).Should(Equal(order.Open))
			}
		})
	})
})

// ClearDB will delete the db files from previous tests.
func ClearDB() error {
	cmd := exec.Command("/bin/sh", "-c", "rm -rf $HOME/.darknode/orderbook")
	return cmd.Run()
}

func NewOrder(parity order.Parity, i uint) *order.Order{
	return order.NewOrder(order.TypeLimit, parity, time.Now(), order.CurrencyCodeREN, order.CurrencyCodeBTC, stackint.FromUint(i),stackint.FromUint(i),stackint.FromUint(i),stackint.FromUint(i) )
}

