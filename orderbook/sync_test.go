package orderbook_test

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/order"
)

var ErrCannotFindOrder = errors.New("cannot find order")

var _ = Describe("Syncer", func() {

	Context("when syncing with the ledger", func() {

		It("should sync with renledger by updating memory and returning a ChangeSet", func() {
			numberOfOrderPairs := 40
			renLimit := 10

			renLedger := newMockRenLedger()

			// Initial sync with an empty Ren Ledger should return an empty changeset
			syncer := NewSyncer(&renLedger, renLimit)
			changeSet, err := syncer.Sync()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(changeSet)).Should(Equal(0))

			// Opening some buy and sell orders on the Ren Ledger
			err = renLedger.openBuyAndSellOrders(numberOfOrderPairs)
			Expect(err).ShouldNot(HaveOccurred())

			// Sync and expect 2*renLimit number of change sets to be created
			changeSet, err = syncer.Sync()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(changeSet)).Should(Equal(renLimit * 2))

			// Randomly cancel a quarter of the orders
			for i := 0; i < numberOfOrderPairs/4; i++ {
				id := generateRandomOrderID(2 * numberOfOrderPairs)
				err = renLedger.CancelOrder([65]byte{}, [32]byte{byte(id)})
			}
			// Sync and expect canceled orders to be purged
			changeSet, err = syncer.Sync()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(changeSet) < renLimit*4).Should(BeTrue())

			// Randomly confirm half of the total orders
			for i := 0; i < numberOfOrderPairs/2; i++ {
				id := generateRandomOrderID(2 * numberOfOrderPairs)
				matchID := [32]byte{byte(id + 1)}
				err = renLedger.ConfirmOrder([32]byte{byte(id)}, matchID)
			}

			// Sync and expect confirmed orders to have purged
			changeSet, err = syncer.Sync()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(changeSet) < renLimit*8).Should(BeTrue())
		})
	})
})

type mockOrder struct {
	status   order.Status
	parity   order.Parity
	priority uint64
}

type mockRenLedger struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu *sync.Mutex
	orders   map[order.ID]mockOrder
}

func newMockRenLedger() mockRenLedger {
	return mockRenLedger{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu: new(sync.Mutex),
		orders:   map[order.ID]mockOrder{},
	}
}

func (renLedger *mockRenLedger) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; !ok {
		renLedger.orders[orderID] = mockOrder{
			status:   order.Open,
			parity:   order.ParityBuy,
			priority: binary.LittleEndian.Uint64(orderID[:]),
		}
		renLedger.buyOrders = append(renLedger.buyOrders, orderID)
		return nil
	}
	return errors.New("cannot open order that is already open")
}

func (renLedger *mockRenLedger) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	renLedger.ordersMu.Lock()
	renLedger.sellOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.sellOrdersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; !ok {
		renLedger.orders[orderID] = mockOrder{
			status:   order.Open,
			parity:   order.ParitySell,
			priority: binary.LittleEndian.Uint64(orderID[:]),
		}
		renLedger.sellOrders = append(renLedger.sellOrders, orderID)
		return nil
	}
	return errors.New("cannot open order that is already open")
}

func (renLedger *mockRenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	return renLedger.setOrderStatus(orderID, order.Canceled)
}

func (renLedger *mockRenLedger) ConfirmOrder(id order.ID, match order.ID) error {
	if err := renLedger.setOrderStatus(id, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	renLedger.setOrderStatus(match, order.Confirmed)
	return nil
}

func (renLedger *mockRenLedger) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (renLedger *mockRenLedger) Status(orderID order.ID) (order.Status, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if ord, ok := renLedger.orders[orderID]; ok {
		return ord.status, nil
	}
	return order.Nil, ErrCannotFindOrder
}

func (renLedger *mockRenLedger) Priority(orderID order.ID) (uint64, error) {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if ord, ok := renLedger.orders[orderID]; ok {
		return ord.priority, nil
	}
	return uint64(0), ErrCannotFindOrder
}

func (renLedger *mockRenLedger) Trader(orderID order.ID) (string, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) OrderMatch(orderID order.ID) (order.ID, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) Depth(orderID order.ID) (uint, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) BuyOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	orders := []order.ID{}
	end := offset + limit
	if end > len(renLedger.buyOrders) {
		end = len(renLedger.buyOrders)
	}
	for i := offset; i < end; i++ {
		orderID := renLedger.buyOrders[i]
		if buyOrder, ok := renLedger.orders[orderID]; ok {
			if buyOrder.parity == order.ParityBuy && buyOrder.status == order.Open {
				orders = append(orders, orderID)
			}
		}

	}
	return orders, nil
}

func (renLedger *mockRenLedger) SellOrders(offset, limit int) ([]order.ID, error) {
	renLedger.ordersMu.Lock()
	renLedger.buyOrdersMu.Lock()
	defer renLedger.ordersMu.Unlock()
	defer renLedger.buyOrdersMu.Unlock()

	orders := []order.ID{}
	end := offset + limit
	if end > len(renLedger.sellOrders) {
		end = len(renLedger.sellOrders)
	}
	for i := offset; i < end; i++ {
		orderID := renLedger.sellOrders[i]
		if sellOrder, ok := renLedger.orders[orderID]; ok {
			if sellOrder.parity == order.ParitySell && sellOrder.status == order.Open {
				orders = append(orders, orderID)
			}
		}
	}
	return orders, nil
}

func (renLedger *mockRenLedger) setOrderStatus(orderID order.ID, status order.Status) error {
	renLedger.ordersMu.Lock()
	defer renLedger.ordersMu.Unlock()

	if _, ok := renLedger.orders[orderID]; ok {
		ord := renLedger.orders[orderID]
		ord.status = status
		renLedger.orders[orderID] = ord
		return nil
	}
	return ErrCannotFindOrder
}

func (renLedger *mockRenLedger) openBuyAndSellOrders(n int) error {
	for i := 0; i < 2*n; i += 2 {
		if err := renLedger.OpenBuyOrder([65]byte{}, [32]byte{byte(i)}); err != nil {
			return err
		}
		if err := renLedger.OpenSellOrder([65]byte{}, [32]byte{byte(i + 1)}); err != nil {
			return err
		}
	}
	return nil
}

func generateRandomOrderID(numberOfOrders int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Intn(numberOfOrders)
}
