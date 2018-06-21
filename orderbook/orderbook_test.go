package orderbook_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/leveldb"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	numberOfOrders = 20
)

var _ = Describe("Orderbook", func() {

	Context("when opening new orders", func() {

		It("should not return an error and must add fragment to storer", func() {
			// Generate new RSA key
			rsaKey, err := crypto.RandomRsaKey()
			Ω(err).ShouldNot(HaveOccurred())

			// Create mock syncer and storer
			syncer := testutils.NewSyncer(numberOfOrders)
			storer, err := leveldb.NewStore("./data.out")
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, syncer, storer)

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

			iter, err := storer.OrderFragments()
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
			syncer := testutils.NewSyncer(numberOfOrders)
			storer, err := leveldb.NewStore("./data.out")
			Expect(err).ShouldNot(HaveOccurred())
			defer func() {
				os.RemoveAll("./data.out")
			}()

			// Create orderbook
			orderbook := NewOrderbook(rsaKey, syncer, storer)

			Ω(syncer.HasSynced()).Should(BeFalse())
			changeset, err := orderbook.Sync()
			Ω(err).ShouldNot(HaveOccurred())
			Ω(len(changeset)).Should(BeZero())
			Ω(syncer.HasSynced()).Should(BeTrue())
		})
	})
})
