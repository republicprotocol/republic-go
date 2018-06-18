package ome_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/ome"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/orderbook"
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
		epoch    cal.Epoch
		storer   Storer
		book     orderbook.Orderbook
		smpcer   smpc.Smpcer
		ledger   cal.RenLedger
		accounts cal.DarkpoolAccounts

		// Ome components
		ranker    Ranker
		matcher   Matcher
		confirmer Confirmer
		settler   Settler
	)

	Context("ome should manage everything about order matching ", func() {

		BeforeEach(func() {
			done = make(chan struct{})
			addr, err = testutils.RandomAddress()
			Ω(err).ShouldNot(HaveOccurred())
			epoch = newEpoch(0, addr)
			storer = testutils.NewStorer()
			book, err = testutils.NewOrderbook()
			Ω(err).ShouldNot(HaveOccurred())
			smpcer = testutils.NewAlwaysMatchSmpc()
			ledger = testutils.NewRenLedger()
			accounts = testutils.NewDarkpoolAccounts()

			ranker, err = NewRanker(done, addr, epoch)
			Ω(err).ShouldNot(HaveOccurred())
			matcher = NewMatcher(storer, smpcer)
			confirmer = NewConfirmer(storer, ledger, PollInterval, Depth)
			settler = NewSettler(storer, smpcer, accounts)
		})

		AfterEach(func() {
			close(done)
		})

		It("should be able to sync with the order book ", func() {
			ome := NewOme(addr, ranker, matcher, confirmer, settler, storer, book, smpcer, epoch)
			errs := ome.Run(done)
			go func() {
				defer GinkgoRecover()

				for err := range errs {
					Ω(err).ShouldNot(HaveOccurred())
				}
			}()
		})

		//It("should be able to listen for epoch change event", func() {
		//	done := make(chan struct{})
		//	ome := NewOme(ranker, computer, book, smpcer)
		//	go func() {
		//		defer close(done)
		//
		//		time.Sleep(3 * time.Second)
		//		epoch := cal.Epoch{}
		//		ome.OnChangeEpoch(epoch)
		//		time.Sleep(3 * time.Second)
		//	}()
		//	errs := ome.Run(done)
		//	for err := range errs {
		//		Ω(err).ShouldNot(HaveOccurred())
		//	}
		//})
	})
})
