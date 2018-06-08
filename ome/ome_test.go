package ome_test

// TODO : FIX THIS TEST
//import (
//	"context"
//	"sync/atomic"
//	"time"
//
//	. "github.com/onsi/ginkgo"
//	. "github.com/onsi/gomega"
//	. "github.com/republicprotocol/republic-go/ome"
//
//	"github.com/republicprotocol/republic-go/cal"
//	"github.com/republicprotocol/republic-go/crypto"
//	"github.com/republicprotocol/republic-go/order"
//	"github.com/republicprotocol/republic-go/orderbook"
//	"github.com/republicprotocol/republic-go/smpc"
//)
//
//var _ = Describe("Ome", func() {
//	Context("ome should manage everything about order matching ", func() {
//		var computer Computer
//		var confirmer Confirmer
//		var ranker Ranker
//
//		var book orderbook.Orderbook
//		var ledger cal.RenLedger
//		var accounts cal.DarkpoolAccounts
//		var smpcer smpc.Smpcer
//		var storer Storer
//
//		BeforeEach(func() {
//			// Generate mock instance for all the parts we need
//			book = newMockOrderbook()
//			ledger = newMockRenLedger()
//			smpcer = newMockSmpcer(PassAll)
//			storer = NewMockStorer()
//			accounts = newMockAccounts()
//
//			confirmer = NewConfirmer(0, 2*time.Second, ledger, storer)
//			computer = NewComputer(storer, smpcer, confirmer, ledger, accounts)
//			ranker = NewRanker(1, 0)
//		})
//
//		It("should be able to sync with the order book ", func() {
//			done := make(chan struct{})
//			go func() {
//				defer close(done)
//				time.Sleep(10 * time.Second)
//			}()
//			ome := NewOme(ranker, computer, book, smpcer)
//			errs := ome.Run(done)
//			for err := range errs {
//				Ω(err).ShouldNot(HaveOccurred())
//			}
//		})
//
//		It("should be able to listen for epoch change event", func() {
//			done := make(chan struct{})
//			ome := NewOme(ranker, computer, book, smpcer)
//			go func() {
//				defer close(done)
//
//				time.Sleep(3 * time.Second)
//				epoch := cal.Epoch{}
//				ome.OnChangeEpoch(epoch)
//				time.Sleep(3 * time.Second)
//			}()
//			errs := ome.Run(done)
//			for err := range errs {
//				Ω(err).ShouldNot(HaveOccurred())
//			}
//		})
//
//	})
//})
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
//func newMockOrderbook() *mockOrderbook {
//	// Generate new RSA key
//	rsaKey, err := crypto.RandomRsaKey()
//	Expect(err).ShouldNot(HaveOccurred())
//
//	return &mockOrderbook{
//		numberOfOrderMatches: int64(0),
//		hasSynced:            false,
//		rsaKey:               rsaKey,
//		orderFragments:       make(map[order.ID]order.Fragment),
//		orders:               make(map[order.ID]order.Order),
//	}
//}
