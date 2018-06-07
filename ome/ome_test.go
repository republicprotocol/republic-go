package ome_test

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/identity"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Ome", func() {
	Context("ome should manage everything about order matching ", func() {
		var confirmer Confirmer
		var ranker Ranker

		var book orderbook.Orderbook
		var ledger cal.RenLedger
		var accounts cal.DarkpoolAccounts
		var smpcer smpc.Smpcer
		var storer Storer
		var matcher Matcher
		var settler Settler

		BeforeEach(func() {
			// Generate mock instance for all the parts we need
			book = newMockOrderbook()
			ledger = newMockRenLedger()
			smpcer = newMockSmpcer()
			storer = newMockStorer()
			accounts = newMockAccounts()

			confirmer = NewConfirmer(storer, ledger, 2*time.Second, 0)
			ranker = NewRanker(1, 0, storer)
			matcher = NewMatcher(storer, smpcer)
			settler = NewSettler(storer, smpcer, accounts)
		})

		It("should be able to sync with the order book ", func() {
			done := make(chan struct{})
			go func() {
				defer close(done)
				time.Sleep(10 * time.Second)
			}()
			ome := NewOme(ranker, matcher, confirmer, settler, storer, book, smpcer)
			errs := ome.Run(done)
			for err := range errs {
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

		It("should be able to listen for epoch change event", func() {
			done := make(chan struct{})
			ome := NewOme(ranker, matcher, confirmer, settler, storer, book, smpcer)
			go func() {
				defer close(done)

				time.Sleep(3 * time.Second)
				epoch := cal.Epoch{}
				ome.OnChangeEpoch(epoch)
				time.Sleep(3 * time.Second)
			}()
			errs := ome.Run(done)
			for err := range errs {
				Ω(err).ShouldNot(HaveOccurred())
			}
		})

	})
})

type mockOrderbook struct {
	numberOfOrderMatches int64
	hasSynced            bool
	rsaKey               crypto.RsaKey
	orderFragments       map[order.ID]order.Fragment
	orders               map[order.ID]order.Order
}

func (book *mockOrderbook) OrderFragment(orderID order.ID) (order.Fragment, error) {
	return book.orderFragments[orderID], nil
}

func (book *mockOrderbook) Order(orderID order.ID) (order.Order, error) {
	return book.orders[orderID], nil
}

func (book *mockOrderbook) ConfirmOrderMatch(buy order.ID, sell order.ID) error {
	atomic.AddInt64(&book.numberOfOrderMatches, 1)
	return nil
}

func (book *mockOrderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	var err error
	book.orderFragments[orderFragment.OrderID], err = orderFragment.Decrypt(*book.rsaKey.PrivateKey)
	return err
}

func (book *mockOrderbook) Sync() (orderbook.ChangeSet, error) {
	if !book.hasSynced {
		changes := make(orderbook.ChangeSet, 5)
		i := 0
		for _, orderFragment := range book.orderFragments {
			changes[i] = orderbook.Change{
				OrderID:       orderFragment.OrderID,
				OrderParity:   orderFragment.OrderParity,
				OrderPriority: uint64(i),
				OrderStatus:   order.Open,
			}
			i++
		}
		book.hasSynced = true
		return changes, nil
	}
	return orderbook.ChangeSet{}, nil
}

func (book *mockOrderbook) AddOrder(ord order.Order) {
	if _, ok := book.orders[ord.ID]; !ok {
		book.orders[ord.ID] = ord
	}
}

func newMockOrderbook() *mockOrderbook {
	// Generate new RSA key
	rsaKey, err := crypto.RandomRsaKey()
	Expect(err).ShouldNot(HaveOccurred())

	return &mockOrderbook{
		numberOfOrderMatches: int64(0),
		hasSynced:            false,
		rsaKey:               rsaKey,
		orderFragments:       make(map[order.ID]order.Fragment),
		orders:               make(map[order.ID]order.Order),
	}
}

// mockStorer is a mock implementation of the storer it stores the order and
// orderFragments in the cache.
type mockStorer struct {
	mu             *sync.Mutex
	orders         map[[32]byte]order.Order
	orderFragments map[[32]byte]order.Fragment
}

func newMockStorer() *mockStorer {
	return &mockStorer{
		mu:             new(sync.Mutex),
		orders:         map[[32]byte]order.Order{},
		orderFragments: map[[32]byte]order.Fragment{},
	}
}

func (storer mockStorer) InsertOrderFragment(fragment order.Fragment) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orderFragments[fragment.OrderID] = fragment
	return nil
}

func (storer mockStorer) InsertOrder(order order.Order) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	storer.orders[order.ID] = order
	return nil
}

func (storer mockStorer) OrderFragment(id order.ID) (order.Fragment, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	if _, ok := storer.orderFragments[id]; !ok {
		return order.Fragment{}, errors.New("no such fragment")
	}
	return storer.orderFragments[id], nil
}

func (storer mockStorer) Order(id order.ID) (order.Order, error) {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	return storer.orders[id], nil
}

func (storer mockStorer) RemoveOrderFragment(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orderFragments, id)
	return nil
}

func (storer mockStorer) RemoveOrder(id order.ID) error {
	storer.mu.Lock()
	defer storer.mu.Unlock()

	delete(storer.orders, id)
	return nil
}

func (store *mockStorer) InsertComputation(computations Computation) error {
	return nil
}

func (store *mockStorer) Computation(id ComputationID) (Computation, error) {
	return Computation{}, ErrComputationNotFound
}

func (store *mockStorer) Computations() (Computations, error) {
	return Computations{}, nil
}

type mockSmpcer struct {
	result uint64
}

func (smpcer *mockSmpcer) Connect(networkID smpc.NetworkID, nodes identity.Addresses) {
}

func (smpcer *mockSmpcer) Disconnect(networkID smpc.NetworkID) {
}

func (smpcer *mockSmpcer) Join(networkID smpc.NetworkID, join smpc.Join, callback smpc.Callback) error {
	return nil
}

func newMockSmpcer() *mockSmpcer {
	return &mockSmpcer{
		result: uint64(0),
	}
}

func newOrder(isBuy bool) order.Order {
	price := uint64(rand.Intn(2000))
	volume := uint64(rand.Intn(2000))
	nonce := int64(rand.Intn(1000000000))
	parity := order.ParityBuy
	if !isBuy {
		parity = order.ParitySell
	}
	return order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce)
}

func computeID(computation Computation) [32]byte {
	id := [32]byte{}
	copy(id[:], crypto.Keccak256(computation.Buy[:], computation.Sell[:]))
	return id
}

func randomOrderFragment() (order.Fragment, order.Fragment, error) {
	buyOrder := newRandomOrder()
	sellOrder := order.NewOrder(buyOrder.Type, 1-buyOrder.Parity, buyOrder.Expiry, buyOrder.Tokens, buyOrder.Price, buyOrder.Volume, buyOrder.MinimumVolume, buyOrder.Nonce)

	buyShares, err := buyOrder.Split(10, 7)
	if err != nil {
		return order.Fragment{}, order.Fragment{}, err
	}
	sellShares, err := sellOrder.Split(10, 7)
	if err != nil {
		return order.Fragment{}, order.Fragment{}, err
	}

	return buyShares[0], sellShares[0], nil
}

type mockAccounts struct {
}

func newMockAccounts() cal.DarkpoolAccounts {
	return &mockAccounts{}
}

func (accounts *mockAccounts) Settle(buy order.Order, sell order.Order) error {
	return nil
}

func (accounts *mockAccounts) Balance(trader string, token order.Token) (float64, error) {
	return 0, nil
}
