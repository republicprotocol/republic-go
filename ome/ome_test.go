package ome_test

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/registry"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/testutils"
)

const (
	Depth        = uint(0)
	PollInterval = time.Second
)

var _ = Describe("Ome", func() {
	var (
		done     chan struct{}
		addr     identity.Address
		err      error
		epoch    registry.Epoch
		storer   ComputationStorer
		book     orderbook.Orderbook
		smpcer   smpc.Smpcer
		contract ContractBinder

		// Ome components
		computationsGenerator ComputationGenerator
		matcher               Matcher
		confirmer             Confirmer
		settler               Settler
	)

	Context("ome should manage everything about order matching ", func() {

		BeforeEach(func() {
			done = make(chan struct{})
			addr, epoch, err = testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())

			computationsGenerator = NewComputationGenerator()
			rsaKey, err := crypto.RandomRsaKey()
			Ω(err).ShouldNot(HaveOccurred())
			book = testutils.NewRandOrderbook(rsaKey)
			Ω(err).ShouldNot(HaveOccurred())
			smpcer = testutils.NewAlwaysMatchSmpc()
			contract = newOmeBinder()

			Ω(err).ShouldNot(HaveOccurred())
			matcher = NewMatcher(storer, smpcer)
			confirmer = NewConfirmer(storer, contract, PollInterval, Depth)
			settler = NewSettler(storer, smpcer, contract)

			store, err := leveldb.NewStore("./data.out", 72*time.Hour)
			Ω(err).ShouldNot(HaveOccurred())
			storer = store.SomerComputationStore()
		})

		AfterEach(func() {
			close(done)
			os.RemoveAll("./data.out")
		})

		It("should be able to sync with the order book ", func() {
			ome := NewOme(addr, computationsGenerator, matcher, confirmer, settler, storer, book, smpcer, epoch)
			errs := ome.Run(done)
			go func() {
				defer GinkgoRecover()

				for err := range errs {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()
		})

		It("should be able to listen for epoch change event", func() {
			ome := NewOme(addr, computationsGenerator, matcher, confirmer, settler, storer, book, smpcer, epoch)
			errs := ome.Run(done)

			go func() {
				defer GinkgoRecover()

				for err := range errs {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()

			_, epoch, err := testutils.RandomEpoch(0)
			Ω(err).ShouldNot(HaveOccurred())

			ome.OnChangeEpoch(epoch)
		})
	})
})

// ErrOpenOpenedOrder is returned when trying to open an opened order.
var ErrOpenOpenedOrder = errors.New("cannot open order that is already open")

// omeBinder is a mock implementation of ome.ContractBinder.
type omeBinder struct {
	buyOrdersMu *sync.Mutex
	buyOrders   []order.ID

	sellOrdersMu *sync.Mutex
	sellOrders   []order.ID

	ordersMu    *sync.Mutex
	orders      map[order.ID]int
	orderStatus map[order.ID]order.Status

	mu    *sync.Mutex
	comps int
	buys  map[order.ID]struct{}
	sells map[order.ID]struct{}
}

// newOmeBinder returns a mock omeBinder.
func newOmeBinder() *omeBinder {
	return &omeBinder{
		buyOrdersMu: new(sync.Mutex),
		buyOrders:   []order.ID{},

		sellOrdersMu: new(sync.Mutex),
		sellOrders:   []order.ID{},

		ordersMu:    new(sync.Mutex),
		orders:      map[order.ID]int{},
		orderStatus: map[order.ID]order.Status{},

		mu:    new(sync.Mutex),
		comps: 0,
		buys:  map[order.ID]struct{}{},
		sells: map[order.ID]struct{}{},
	}
}

// ConfirmOrder confirm a order pair is a match.
func (binder *omeBinder) ConfirmOrder(id order.ID, match order.ID) error {
	if err := binder.setOrderStatus(id, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	if err := binder.setOrderStatus(match, order.Confirmed); err != nil {
		return fmt.Errorf("cannot confirm order that is not open: %v", err)
	}
	return nil
}

// Status returns the status of the order by the order ID.
func (binder *omeBinder) Status(orderID order.ID) (order.Status, error) {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	if status, ok := binder.orderStatus[orderID]; ok {
		return status, nil
	}
	return order.Nil, ErrOrderNotFound
}

// Trader returns the matched order of the order by the order ID.
func (binder *omeBinder) OrderMatch(orderID order.ID) (order.ID, error) {
	binder.ordersMu.Lock()
	defer binder.ordersMu.Unlock()

	if binder.orderStatus[orderID] != order.Confirmed {
		return order.ID{}, errors.New("order is not open ")
	}
	for i, id := range binder.buyOrders {
		if orderID.Equal(id) {
			return binder.sellOrders[i], nil
		}
	}
	for i, id := range binder.sellOrders {
		if orderID.Equal(id) {
			return binder.buyOrders[i], nil
		}
	}
	return order.ID{}, ErrOrderNotFound
}

// Depth returns the block depth since the order been confirmed.
func (binder *omeBinder) Depth(orderID order.ID) (uint, error) {
	return 100, nil
}

// OpenBuyOrder in the mock omeBinder.
func (binder *omeBinder) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.buyOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.buyOrdersMu.Unlock()

	if _, ok := binder.orders[orderID]; !ok {
		binder.orders[orderID] = len(binder.buyOrders)
		binder.buyOrders = append(binder.buyOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return errors.New("cannot open order that is already open")
}

// OpenSellOrder in the mock omeBinder.
func (binder *omeBinder) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	binder.ordersMu.Lock()
	binder.sellOrdersMu.Lock()
	defer binder.ordersMu.Unlock()
	defer binder.sellOrdersMu.Unlock()

	if _, ok := binder.orders[orderID]; !ok {
		binder.orders[orderID] = len(binder.sellOrders)
		binder.sellOrders = append(binder.sellOrders, orderID)
		binder.orderStatus[orderID] = order.Open
		return nil
	}

	return ErrOpenOpenedOrder
}

// Settle implements the Settle function of the mock omeBinder
func (binder *omeBinder) Settle(buy order.Order, sell order.Order) error {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	binder.comps++
	binder.buys[buy.ID] = struct{}{}
	binder.sells[buy.ID] = struct{}{}

	return nil
}

func (binder *omeBinder) SettleCounts() int {
	binder.mu.Lock()
	defer binder.mu.Unlock()

	return binder.comps
}

func (binder *omeBinder) setOrderStatus(orderID order.ID, status order.Status) error {
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
