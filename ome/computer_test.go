package ome_test

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("Computer", func() {
	Context("when given computations", func() {
		var confirmer Confirmer
		var renLedger cal.RenLedger

		BeforeEach(func() {
			renLedger = newMockRenLedger()
			confirmer = newMockConfirmer()
		})

		It("should successfully complete all computations", func(d Done) {
			defer close(d)

			numberOfComputations := 20

			smpcer := newMockSmpcer(PassAll)
			err := smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())

			accounts := newMockAccounts()

			storer := NewMockStorer()
			computer := NewComputer(storer, smpcer, confirmer, renLedger, accounts)

			done := make(chan struct{})
			go func() {
				defer close(done)

				time.Sleep(3 * time.Second)
			}()
			computationsCh := make(chan ComputationEpoch)
			defer close(computationsCh)

			errs := computer.Compute(done, computationsCh)

			// Generate computations for newly created buy and sell orders
			computations := make(Computations, numberOfComputations)
			for i := 0; i < numberOfComputations; i++ {
				buyOrderFragment, sellOrderFragment, err := randomOrderFragment()
				Expect(err).ShouldNot(HaveOccurred())
				Expect(storer.InsertOrderFragment(buyOrderFragment)).ShouldNot(HaveOccurred())
				Expect(storer.InsertOrderFragment(sellOrderFragment)).ShouldNot(HaveOccurred())

				computations[i] = Computation{
					Buy:      buyOrderFragment.OrderID,
					Sell:     sellOrderFragment.OrderID,
					Priority: Priority(i),
				}

				Expect(renLedger.OpenBuyOrder([65]byte{}, computations[i].Buy)).ShouldNot(HaveOccurred())
				Expect(renLedger.OpenBuyOrder([65]byte{}, computations[i].Sell)).ShouldNot(HaveOccurred())

				computationsCh <- ComputationEpoch{
					Computation: computations[i],
					ID:          computeID(computations[i]),
					Epoch:       [32]byte{1},
				}
				log.Print(i)
			}

			err, ok := <-errs
			Expect(err).To(BeNil())
			Expect(ok).Should(BeFalse())

		}, 300 /* 5 minute timeout */)
	})
})

// mockStorer is a mock implementation of the storer it stores the order and
// orderFragments in the cache.
type mockStorer struct {
	mu             *sync.Mutex
	orders         map[[32]byte]order.Order
	orderFragments map[[32]byte]order.Fragment
}

func NewMockStorer() *mockStorer {
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

func (store *mockStorer) Computation(id [32]byte) (Computation, error) {
	return Computation{}, ErrComputationNotFound
}

type mockSmpcer struct {
	result       uint64
	instructions chan smpc.Inst
	results      chan smpc.Result
}

func (smpcer *mockSmpcer) Start() error {
	smpcer.instructions = make(chan smpc.Inst)
	smpcer.results = make(chan smpc.Result)

	go func() {
		for {
			select {
			case inst, ok := <-smpcer.instructions:
				if !ok {
					return
				}
				if inst.InstJ != nil {
					smpcer.results <- smpc.Result{
						InstID:    inst.InstID,
						NetworkID: inst.NetworkID,
						ResultJ: &smpc.ResultJ{
							Value: smpcer.result,
						},
					}
				}
			}
		}
	}()
	return nil
}

func (smpcer *mockSmpcer) Shutdown() error {
	close(smpcer.instructions)
	close(smpcer.results)
	return nil
}

func (smpcer *mockSmpcer) Instructions() chan<- smpc.Inst {
	return smpcer.instructions
}

func (smpcer *mockSmpcer) Results() <-chan smpc.Result {
	return smpcer.results
}

type ResultStatus uint8

// ResultStatus values.
const (
	FailAll     ResultStatus = 0
	PassAll     ResultStatus = 1
	PassPartial ResultStatus = 2
)

func newMockSmpcer(resultStatus ResultStatus) *mockSmpcer {
	result := uint64(0)

	switch resultStatus {
	case FailAll:
		result = uint64(shamir.Prime - 1)
	case PassAll:
		result = uint64(0)
	case PassPartial:
		result = uint64(1)
	}

	return &mockSmpcer{
		result: result,
	}
}

type mockConfirmer struct {
}

func newMockConfirmer() Confirmer {
	return &mockConfirmer{}
}

func (confirmer *mockConfirmer) ConfirmOrderMatches(done <-chan struct{}, orderMatches <-chan Computation) (<-chan Computation, <-chan error) {
	confirmedMatches := make(chan Computation)
	errs := make(chan error)

	go func() {
		for {
			select {
			case <-done:
				return
			case computation, ok := <-orderMatches:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case confirmedMatches <- computation:
				}
			}
		}
	}()

	return confirmedMatches, errs
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
