package ome_test

//
//import (
//	"context"
//	"math/rand"
//	"sync"
//	"sync/atomic"
//	"time"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	. "github.com/republicprotocol/republic-go/ome"
//
//	"github.com/republicprotocol/republic-go/crypto"
//	"github.com/republicprotocol/republic-go/order"
//	"github.com/republicprotocol/republic-go/orderbook"
//	"github.com/republicprotocol/republic-go/shamir"
//	"github.com/republicprotocol/republic-go/smpc"
//)
//
//var _ = Describe("Computer", func() {
//
//	Context("when given computations", func() {
//
//		It("should successfully complete all computations", func(done Done) {
//			defer close(done)
//
//			numberOfComputations := 10
//
//			// Start mockSmpcer
//			smpcer := newMockSmpcer(PassAll)
//			err := smpcer.Start()
//			Expect(err).ShouldNot(HaveOccurred())
//
//			// Create mockOrderbook
//			orderbook := newMockOrderbook()
//
//			// Create Computer
//			computer := NewComputer(&orderbook, &smpcer)
//
//			quit := make(chan struct{})
//			defer close(quit)
//			go computer.ComputeResults(quit)
//
//			// Generate computations for newly created buy and sell orders
//			computations := make(Computations, numberOfComputations)
//			for i := 0; i < numberOfComputations; i++ {
//				buyOrder := newOrder(true)
//				orderbook.AddOrder(buyOrder)
//
//				sellOrder := newOrder(false)
//				orderbook.AddOrder(sellOrder)
//
//				computations[i] = Computation{
//					Buy:      buyOrder.ID,
//					Sell:     sellOrder.ID,
//					Priority: Priority(i),
//				}
//			}
//
//			computer.Compute([32]byte{1}, computations)
//
//			// Constantly call Compute() until all orders have been matched
//			computations = Computations{}
//			for {
//				time.Sleep(20 * time.Millisecond)
//				computer.Compute([32]byte{1}, computations)
//
//				if atomic.LoadInt64(&orderbook.numberOfOrderMatches) == int64(numberOfComputations) {
//					break
//				}
//			}
//			Expect(atomic.LoadInt64(&orderbook.numberOfOrderMatches)).Should(Equal(int64(numberOfComputations)))
//		}, 2 /* 2 second timeout */)
//	})
//})
//
//// mockStorer is a mock implementation of the storer it stores the order and
//// orderFragments in the cache.
//type mockStorer struct {
//	mu             *sync.Mutex
//	orders         map[[32]byte]order.Order
//	orderFragments map[[32]byte]order.Fragment
//}
//
//func NewMockStorer() *mockStorer {
//	return &mockStorer{
//		mu:             new(sync.Mutex),
//		orders:         map[[32]byte]order.Order{},
//		orderFragments: map[[32]byte]order.Fragment{},
//	}
//}
//
//func (storer mockStorer) InsertOrderFragment(fragment order.Fragment) error {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	storer.orderFragments[fragment.OrderID] = fragment
//	return nil
//}
//
//func (storer mockStorer) InsertOrder(order order.Order) error {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	storer.orders[order.ID] = order
//	return nil
//}
//
//func (storer mockStorer) OrderFragment(id order.ID) (order.Fragment, error) {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	return storer.orderFragments[id], nil
//}
//
//func (storer mockStorer) Order(id order.ID) (order.Order, error) {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	return storer.orders[id], nil
//}
//
//func (storer mockStorer) RemoveOrderFragment(id order.ID) error {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	delete(storer.orderFragments, id)
//	return nil
//}
//
////
//
//func (storer mockStorer) RemoveOrder(id order.ID) error {
//	storer.mu.Lock()
//	defer storer.mu.Unlock()
//
//	delete(storer.orders, id)
//	return nil
//}
//
//type mockOrderbook struct {
//	numberOfOrderMatches int64
//	hasSynced            bool
//	rsaKey               crypto.RsaKey
//	orderFragments       map[order.ID]order.Fragment
//	orders               map[order.ID]order.Order
//}
//
//func (book *mockOrderbook) OrderFragment(orderID order.ID) (order.Fragment, error) {
//	return book.orderFragments[orderID], nil
//}
//
//func (book *mockOrderbook) Order(orderID order.ID) (order.Order, error) {
//	return book.orders[orderID], nil
//}
//
//func (book *mockOrderbook) ConfirmOrderMatch(buy order.ID, sell order.ID) error {
//	atomic.AddInt64(&book.numberOfOrderMatches, 1)
//	return nil
//}
//
//func (book *mockOrderbook) OpenOrder(ctx context.Context, orderFragment order.EncryptedFragment) error {
//	var err error
//	book.orderFragments[orderFragment.OrderID], err = orderFragment.Decrypt(*book.rsaKey.PrivateKey)
//	return err
//}
//
//func (book *mockOrderbook) Sync() (orderbook.ChangeSet, error) {
//	if !book.hasSynced {
//		changes := make(orderbook.ChangeSet, 5)
//		i := 0
//		for _, orderFragment := range book.orderFragments {
//			changes[i] = orderbook.Change{
//				OrderID:       orderFragment.OrderID,
//				OrderParity:   orderFragment.OrderParity,
//				OrderPriority: uint64(i),
//				OrderStatus:   order.Open,
//			}
//			i++
//		}
//		book.hasSynced = true
//		return changes, nil
//	}
//	return orderbook.ChangeSet{}, nil
//}
//
//func (book *mockOrderbook) AddOrder(ord order.Order) {
//	if _, ok := book.orders[ord.ID]; !ok {
//		book.orders[ord.ID] = ord
//	}
//}
//
//func newMockOrderbook() mockOrderbook {
//	// Generate new RSA key
//	rsaKey, err := crypto.RandomRsaKey()
//	Expect(err).ShouldNot(HaveOccurred())
//
//	return mockOrderbook{
//		numberOfOrderMatches: int64(0),
//		hasSynced:            false,
//		rsaKey:               rsaKey,
//		orderFragments:       make(map[order.ID]order.Fragment),
//		orders:               make(map[order.ID]order.Order),
//	}
//}
//
//type mockSmpcer struct {
//	result       uint64
//	instructions chan smpc.Inst
//	results      chan smpc.Result
//}
//
//func (smpcer *mockSmpcer) Start() error {
//	smpcer.instructions = make(chan smpc.Inst)
//	smpcer.results = make(chan smpc.Result)
//
//	go func() {
//		for {
//			select {
//			case inst, ok := <-smpcer.instructions:
//				if !ok {
//					return
//				}
//				if inst.InstJ != nil {
//					smpcer.results <- smpc.Result{
//						InstID:    inst.InstID,
//						NetworkID: inst.NetworkID,
//						ResultJ: &smpc.ResultJ{
//							Value: smpcer.result,
//						},
//					}
//				}
//			}
//		}
//	}()
//	return nil
//}
//
//func (smpcer *mockSmpcer) Shutdown() error {
//	close(smpcer.instructions)
//	close(smpcer.results)
//	return nil
//}
//
//func (smpcer *mockSmpcer) Instructions() chan<- smpc.Inst {
//	return smpcer.instructions
//}
//
//func (smpcer *mockSmpcer) Results() <-chan smpc.Result {
//	return smpcer.results
//}
//
//type ResultStatus uint8
//
//// ResultStatus values.
//const (
//	FailAll     ResultStatus = 0
//	PassAll     ResultStatus = 1
//	PassPartial ResultStatus = 2
//)
//
//func newMockSmpcer(resultStatus ResultStatus) mockSmpcer {
//	result := uint64(0)
//
//	switch resultStatus {
//	case FailAll:
//		result = uint64(shamir.Prime - 1)
//	case PassAll:
//		result = uint64(0)
//	case PassPartial:
//		result = uint64(1)
//	}
//
//	return mockSmpcer{
//		result: result,
//	}
//}
//
//func newOrder(isBuy bool) order.Order {
//	price := uint64(rand.Intn(2000))
//	volume := uint64(rand.Intn(2000))
//	nonce := int64(rand.Intn(1000000000))
//	parity := order.ParityBuy
//	if !isBuy {
//		parity = order.ParitySell
//	}
//	return order.NewOrder(order.TypeLimit, parity, time.Now().Add(time.Hour), order.TokensETHREN, order.NewCoExp(price, 26), order.NewCoExp(volume, 26), order.NewCoExp(volume, 26), nonce)
//}
