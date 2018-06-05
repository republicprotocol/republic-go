package ome_test

import (
	"errors"
	"math/big"
	"math/rand"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Confirmer", func() {
	var confirmer Confirmer
	var renLedger cal.RenLedger
	var storer Storer

	BeforeEach(func() {
		depth, pollInterval := uint(0), time.Second
		renLedger = newMockRenLedger()
		storer = NewMockStorer()
		confirmer = NewConfirmer(depth, pollInterval, renLedger, storer)
	})

	It("should be able to confirm order on the ren ledger", func(d Done) {
		defer close(d)

		var numberOfComputationsToTest = 100

		done := make(chan struct{})
		orderMatches := make(chan Computation)
		orderIDs := map[[32]byte]struct{}{}
		computations := make([]Computation, numberOfComputationsToTest)
		for i := 0; i < numberOfComputationsToTest; i++ {
			computations[i] = randomComputaion()
			orderIDs[computations[i].Buy] = struct{}{}
			orderIDs[computations[i].Sell] = struct{}{}
		}

		// Open all the orders
		for i := 0; i < numberOfComputationsToTest; i++ {
			err := renLedger.OpenBuyOrder([65]byte{}, computations[i].Buy)
			Expect(err).ShouldNot(HaveOccurred())
			err = renLedger.OpenSellOrder([65]byte{}, computations[i].Sell)
			Expect(err).ShouldNot(HaveOccurred())
		}

		go func() {
			defer GinkgoRecover()
			defer close(done)

			for i := 0; i < numberOfComputationsToTest; i++ {
				orderMatches <- computations[i]
			}
			time.Sleep(5 * time.Second)
		}()

		confirmedMatches, errs := confirmer.ConfirmOrderMatches(done, orderMatches)

		go func() {
			defer GinkgoRecover()

			for err := range errs {
				Ω(err).ShouldNot(HaveOccurred())
			}
		}()

		for match := range confirmedMatches {
			_, ok := orderIDs[match.Buy]
			Ω(ok).Should(BeTrue())
			delete(orderIDs, match.Buy)

			_, ok = orderIDs[match.Sell]
			Ω(ok).Should(BeTrue())
			delete(orderIDs, match.Sell)
		}

		Ω(len(orderIDs)).Should(Equal(0))
	}, 100)
})

const (
	GenesisBuyer  = "0x90e6572eF66a11690b09dd594a18f36Cf76055C8"
	GenesisSeller = "0x8DF05f77e8aa74D3D8b5342e6007319A470a64ce"
)

type mockRenLedger struct {
	mu         *sync.Mutex
	buyOrders  map[[32]byte]struct{}
	sellOrders map[[32]byte]struct{}
	matches    map[[32]byte][32]byte
	states     map[[32]byte]order.Status
}

func newMockRenLedger() *mockRenLedger {
	return &mockRenLedger{
		mu:         new(sync.Mutex),
		buyOrders:  map[[32]byte]struct{}{},
		sellOrders: map[[32]byte]struct{}{},
		matches:    map[[32]byte][32]byte{},
		states:     map[[32]byte]order.Status{},
	}
}

func (renLedger *mockRenLedger) OpenBuyOrder(signature [65]byte, orderID order.ID) error {
	renLedger.mu.Lock()
	defer renLedger.mu.Unlock()

	renLedger.buyOrders[orderID] = struct{}{}
	renLedger.states[orderID] = order.Open

	return nil
}

func (renLedger *mockRenLedger) OpenSellOrder(signature [65]byte, orderID order.ID) error {
	renLedger.mu.Lock()
	defer renLedger.mu.Unlock()

	renLedger.sellOrders[orderID] = struct{}{}
	renLedger.states[orderID] = order.Open

	return nil
}

func (renLedger *mockRenLedger) CancelOrder(signature [65]byte, orderID order.ID) error {
	renLedger.mu.Lock()
	defer renLedger.mu.Unlock()

	renLedger.states[orderID] = order.Canceled

	return nil
}

func (renLedger *mockRenLedger) ConfirmOrder(id order.ID, match order.ID) error {
	renLedger.mu.Lock()
	defer renLedger.mu.Unlock()

	if _, ok := renLedger.matches[id]; renLedger.states[id] != order.Open || ok {
		return errors.New("you can only confirm open order ")
	}
	if _, ok := renLedger.matches[match]; renLedger.states[match] != order.Open || ok {
		return errors.New("you can only confirm open order ")
	}

	renLedger.states[id] = order.Confirmed
	renLedger.states[match] = order.Confirmed

	renLedger.matches[id] = match
	renLedger.matches[match] = id

	return nil
}

func (renLedger *mockRenLedger) Fee() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (renLedger *mockRenLedger) Status(id order.ID) (order.Status, error) {
	if _, ok := renLedger.matches[id]; ok {
		return order.Confirmed, nil
	}

	_, ok1 := renLedger.buyOrders[id]
	_, ok2 := renLedger.sellOrders[id]
	if ok1 || ok2 {
		return order.Open, nil
	}

	return order.Nil, nil
}

func (renLedger *mockRenLedger) Priority(id order.ID) (uint64, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) Depth(id order.ID) (uint, error) {
	// Return a big number which is deep enough.
	return 100, nil
}

func (renLedger *mockRenLedger) BuyOrders(offset, limit int) ([]order.ID, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) SellOrders(offset, limit int) ([]order.ID, error) {
	panic("unimplemented")
}

func (renLedger *mockRenLedger) OrderMatch(id order.ID) (order.ID, error) {
	renLedger.mu.Lock()
	defer renLedger.mu.Unlock()

	if _, ok := renLedger.matches[id]; !ok {
		return order.ID{}, errors.New("no match found the order.")
	}
	return renLedger.matches[id], nil
}

func (renLedger *mockRenLedger) Trader(id order.ID) (string, error) {
	if _, ok := renLedger.buyOrders[id]; ok {
		return GenesisBuyer, nil
	}

	if _, ok := renLedger.sellOrders[id]; ok {
		return GenesisSeller, nil
	}

	return "", errors.New("order doesn't exist in the ledger ")
}

func randomComputaion() Computation {
	buy, sell := newRandomOrder().ID, newRandomOrder().ID
	return Computation{
		Buy:      buy,
		Sell:     sell,
		Priority: Priority(0),
	}
}

func newRandomOrder() order.Order {
	parity := []order.Parity{order.ParityBuy, order.ParitySell}[rand.Intn(2)]
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensBTCDGX,
		order.TokensBTCREN,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]

	ord := order.NewOrder(order.TypeLimit, parity, time.Now().Add(1*time.Hour), tokens, randomCoExp(), randomCoExp(), randomCoExp(), rand.Int63())
	return ord
}

func randomCoExp() order.CoExp {
	co := uint64(rand.Intn(1999) + 1)
	exp := uint64(rand.Intn(25))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}

func ganacheConfig() ethereum.Config {
	return ethereum.Config{
		Network:                 ethereum.NetworkGanache,
		URI:                     "http://localhost:8545",
		RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
		DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
		RenLedgerAddress:        ethereum.RenLedgerAddressOnGanache.String(),
	}
}
