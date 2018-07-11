package orderbook_test

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Orderbook", func() {

	var (
		numberOfOrders = 20
		done           chan struct{}
	)

	BeforeEach(func() {
		done = make(chan struct{})
	})

	AfterEach(func() {
		close(done)
	})

	Context("when opening new orders", func() {

		It("should not return an error and must add fragment to storer", func() {
			// Generate new RSA key
			rsaKey, err := crypto.RandomRsaKey()
			Ω(err).ShouldNot(HaveOccurred())

			// Create mock syncer and storer
			// syncer := testutils.NewSyncer(numberOfOrders)
			storer, err := leveldb.NewStore("./data.out", 72*time.Hour)
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), testutils.NewMockContractBinder(), time.Hour, 100)

			orderbook.Sync(done)
			orderbook.OnChangeEpoch(registry.Epoch{})

			// Create encryptedOrderFragments
			encryptedOrderFragments := make([]order.EncryptedFragment, numberOfOrders)
			for i := 0; i < numberOfOrders; i++ {
				ord := testutils.RandomOrder()
				fragments, err := ord.Split(5, 4)
				encryptedOrderFragments[i], err = fragments[0].Encrypt(rsaKey.PublicKey)
				Ω(err).ShouldNot(HaveOccurred())
			}

			ctx, cancelCtx := context.WithCancel(context.Background())
			defer cancelCtx()

			// Open all encrypted order fragments
			for i := 0; i < numberOfOrders; i++ {
				err = orderbook.OpenOrder(ctx, encryptedOrderFragments[i])
				Ω(err).ShouldNot(HaveOccurred())
			}

			time.Sleep(time.Second)
			iter, err := storer.OrderbookOrderFragmentStore().OrderFragments(registry.Epoch{})
			Expect(err).ShouldNot(HaveOccurred())
			defer iter.Release()
			collection, err := iter.Collect()
			Expect(err).ShouldNot(HaveOccurred())
			Ω(len(collection)).Should(Equal(numberOfOrders))
		})

		It("should be able to sync with the ledger by the syncer", func() {
			// Generate new RSA key
			rsaKey, err := crypto.RandomRsaKey()
			Ω(err).ShouldNot(HaveOccurred())

			// Create mock syncer and storer
			// syncer := testutils.NewSyncer(numberOfOrders)
			storer, err := leveldb.NewStore("./data.out", 72*time.Hour)
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, storer.OrderbookPointerStore(), storer.OrderbookOrderStore(), storer.OrderbookOrderFragmentStore(), testutils.NewMockContractBinder(), time.Hour, 100)

			// Ω(syncer.HasSynced()).Should(BeFalse())
			doneChan := make(<-chan struct{})
			changeset, _ := orderbook.Sync(doneChan)
			Ω(len(changeset)).Should(BeZero())
			// Ω(syncer.HasSynced()).Should(BeTrue())
		})
	})
})

// ErrOpenOpenedOrder is returned when trying to open an opened order.
var ErrOpenOpenedOrder = errors.New("cannot open order that is already open")

// Constant value of trader address for testing
const (
	GenesisBuyer  = "0x90e6572eF66a11690b09dd594a18f36Cf76055C8"
	GenesisSeller = "0x8DF05f77e8aa74D3D8b5342e6007319A470a64ce"
)

// orderbookBinder is a mock implementation of the orderbook.ContractBinder interface.
type orderbookBinder struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu    *sync.Mutex
	orderStatus map[order.ID]order.Status
}

// newOrderbookBinder returns a mock orderbookBinder.
func newOrderbookBinder() *orderbookBinder {
	return &orderbookBinder{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu:    new(sync.Mutex),
		orderStatus: map[order.ID]order.Status{},
	}
}

// BuyOrders returns a limit of buy orders starting from the offset.
func (binder *orderbookBinder) BuyOrders(offset, limit int) ([]order.ID, error) {
	binder.ordersMu.Lock()
	binder.buyOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.buyOrdersMu.Unlock()

	if offset > len(binder.buyOrders) {
		return []order.ID{}, errors.New("index out of range")
	}
	end := offset + limit
	if end > len(binder.buyOrders) {
		end = len(binder.buyOrders)
	}
	return binder.buyOrders[offset:end], nil
}

// SellOrders returns a limit sell orders starting from the offset.
func (binder *orderbookBinder) SellOrders(offset, limit int) ([]order.ID, error) {
	binder.ordersMu.Lock()
	binder.buyOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.buyOrdersMu.Unlock()

	if offset > len(binder.sellOrders) {
		return []order.ID{}, errors.New("index out of range")
	}
	end := offset + limit
	if end > len(binder.sellOrders) {
		end = len(binder.sellOrders)
	}
	return binder.sellOrders[offset:end], nil
}

// Status returns the status of the order by the order ID.
func (binder *orderbookBinder) Status(orderID order.ID) (order.Status, error) {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	if status, ok := binder.orderStatus[orderID]; ok {
		return status, nil
	}
	return order.Nil, ErrOrderNotFound
}

// Priority returns the priority of the order by the order ID.
func (binder *orderbookBinder) Priority(orderID order.ID) (uint64, error) {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	for _, id := range binder.buyOrders {
		if orderID.Equal(id) {
			return uint64(order.ParityBuy), nil
		}
	}
	for _, id := range binder.sellOrders {
		if orderID.Equal(id) {
			return uint64(order.ParitySell), nil
		}
	}

	return uint64(0), ErrOrderNotFound
}

// Trader returns the trader of the order by the order ID.
func (binder *orderbookBinder) Trader(orderID order.ID) (string, error) {
	return GenesisBuyer, nil
}

// BlockNumber returns the block number when the order being last modified.
func (binder *orderbookBinder) BlockNumber(orderID order.ID) (*big.Int, error) {
	return big.NewInt(100), nil
}

// Depth returns the depth of an order.
func (binder *orderbookBinder) Depth(orderID order.ID) (uint, error) {
	return 10, nil
}

// Depth returns the depth of an order.
func (binder *orderbookBinder) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(10), nil
}

// OpenBuyOrder in the mock orderbookBinder.
func (binder *orderbookBinder) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.buyOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.buyOrdersMu.Unlock()

	if _, ok := binder.orderStatus[orderID]; !ok {
		binder.buyOrders = append(binder.buyOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return errors.New("cannot open order that is already open")
}

// OpenSellOrder in the mock orderbookBinder.
func (binder *orderbookBinder) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.sellOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.sellOrdersMu.Unlock()

	if _, ok := binder.orderStatus[orderID]; !ok {
		binder.sellOrders = append(binder.sellOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return ErrOpenOpenedOrder
}

// CancelOrder in the mock orderbookBinder.
func (binder *orderbookBinder) CancelOrder(signature [65]byte, orderID order.ID) error {
	return binder.setOrderStatus(orderID, order.Canceled)
}

// ConfirmOrder confirm a order pair is a match.
func (binder *orderbookBinder) ConfirmOrder(id order.ID, match order.ID) error {
	if err := binder.setOrderStatus(id, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	if err := binder.setOrderStatus(match, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	return nil
}

func (binder *orderbookBinder) Orders(offset, limit int) ([]order.ID, []order.Status, []string, error) {
	return []order.ID{}, []order.Status{}, []string{}, nil
}

func (binder *orderbookBinder) setOrderStatus(orderID order.ID, status order.Status) error {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	switch status {
	case order.Open:
		binder.orderStatus[orderID] = order.Open
	case order.Confirmed:
		if binder.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		binder.orderStatus[orderID] = order.Confirmed
	case order.Canceled:
		if binder.orderStatus[orderID] != order.Open {
			return errors.New("order not open")
		}
		binder.orderStatus[orderID] = order.Canceled
	}
	return nil
}

// newEpoch returns a new epoch with only one pod and one darknode.
func newEpoch(i int, node identity.Address) registry.Epoch {
	return registry.Epoch{
		Hash: testutils.Random32Bytes(),
		Pods: []registry.Pod{
			{
				Position:  0,
				Hash:      testutils.Random32Bytes(),
				Darknodes: []identity.Address{node},
			},
		},
		Darknodes:     []identity.Address{node},
		BlockNumber:   big.NewInt(int64(i)),
		BlockInterval: big.NewInt(int64(2)),
	}
}
