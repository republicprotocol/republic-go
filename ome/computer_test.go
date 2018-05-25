package ome_test

import (
	"context"
	"math/rand"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

var _ = Describe("OME Computer", func() {

	Context("when opening new orders", func() {

		It("should not return an error", func() {
			numberOfComputations := 1

			smpcer := newMockSmpcer(PassAll)
			err := smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())

			orderbook := newMockOrderbook()

			computer := NewComputer(&orderbook, &smpcer)

			// Generate computations for newly created buy and sell orders
			computations := make(Computations, numberOfComputations)
			for i := 0; i < numberOfComputations; i++ {
				buyOrder := newOrder(true)
				fragments, err := buyOrder.Split(5, 4)
				Expect(err).ShouldNot(HaveOccurred())
				orderbook.orderFragments[buyOrder.ID] = fragments[0]

				sellOrder := newOrder(false)
				fragments, err = sellOrder.Split(5, 4)
				Expect(err).ShouldNot(HaveOccurred())
				orderbook.orderFragments[sellOrder.ID] = fragments[0]

				computations[i] = Computation{
					Buy:      buyOrder.ID,
					Sell:     sellOrder.ID,
					Priority: Priority(i),
				}
			}

			computer.Compute([32]byte{1}, computations)
		})
	})
})

type mockOrderbook struct {
	numberOfOrderMatches int
	synced               bool
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
	book.numberOfOrderMatches++
	return nil
}

func (book *mockOrderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
	var err error
	book.orderFragments[orderFragment.OrderID], err = orderFragment.Decrypt(*book.rsaKey.PrivateKey)
	return err
}

func (book *mockOrderbook) Sync() (orderbook.ChangeSet, error) {
	if !book.synced {
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
		book.synced = true
		return changes, nil
	}
	return orderbook.ChangeSet{}, nil
}

func newMockOrderbook() mockOrderbook {
	// Generate new RSA key
	rsaKey, err := crypto.RandomRsaKey()
	Expect(err).ShouldNot(HaveOccurred())

	return mockOrderbook{
		numberOfOrderMatches: 0,
		synced:               false,
		rsaKey:               rsaKey,
		orderFragments:       make(map[order.ID]order.Fragment),
		orders:               make(map[order.ID]order.Order),
	}
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

func newMockSmpcer(resultStatus ResultStatus) mockSmpcer {
	result := uint64(0)

	switch resultStatus {
	case FailAll:
		result = uint64(shamir.Prime - 1)
	case PassAll:
		result = uint64(0)
	case PassPartial:
		result = uint64(1)
	}

	return mockSmpcer{
		result: result,
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
