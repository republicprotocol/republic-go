package grpc_test

import (
	"fmt"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"
	"golang.org/x/net/context"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/oracle"
	"github.com/republicprotocol/republic-go/swarm"
)

var _ = Describe("Oracle", func() {

	var server *Server
	var service OracleService
	var serviceMultiAddr identity.MultiAddress
	var client oracle.Client
	var multiAddrStorer swarm.MultiAddressStorer
	var midpointPriceStorer oracle.MidpointPriceStorer
	var oracler oracle.Oracler
	var ecdsaKey crypto.EcdsaKey

	BeforeEach(func() {
		var err error

		db, err := leveldb.NewStore("./tmp/oracle.1.out", 10*time.Hour)
		Expect(err).ShouldNot(HaveOccurred())
		multiAddrStorer = db.SwarmMultiAddressStore()
		client, ecdsaKey, err = newOracleClient(multiAddrStorer)
		Expect(err).ShouldNot(HaveOccurred())

		oracler = oracle.NewOracler(client, &ecdsaKey, multiAddrStorer, 10)
		midpointPriceStorer = oracle.NewMidpointPriceStorer()
		service = NewOracleService(oracle.NewServer(oracler, identity.Address(ecdsaKey.Address()), multiAddrStorer, midpointPriceStorer, 10), time.Microsecond)
		serviceMultiAddr = client.MultiAddress()
		server = NewServer()
		service.Register(server)
	})

	AfterEach(func() {
		os.RemoveAll("./tmp")
		server.Stop()
	})

	Context("when storing a midpoint price", func() {
		It("should add the midpoint price to the storer", func() {
			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			var err error
			midpointPrice := &oracle.MidpointPrice{TokenPairs: []uint64{0, 1}, Prices: []uint64{0, 1}, Nonce: 1}
			midpointPrice.Signature, err = ecdsaKey.Sign(midpointPrice.Hash())
			Expect(err).ShouldNot(HaveOccurred())
			err = client.UpdateMidpoint(context.Background(), serviceMultiAddr, *midpointPrice)
			Expect(err).ShouldNot(HaveOccurred())

			prices, err := midpointPriceStorer.MidpointPrice()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(prices.TokenPairs).Should(HaveLen(2))
			Expect(prices.TokenPairs).Should(Equal(midpointPrice.TokenPairs))
			Expect(prices.Prices).Should(HaveLen(2))
			Expect(prices.Prices).Should(Equal(midpointPrice.Prices))
			Expect(prices.Nonce).Should(Equal(midpointPrice.Nonce))
		})

		It("should overwrite existing midpoint prices with new information", func() {
			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			var err error
			midpointPrice := &oracle.MidpointPrice{TokenPairs: []uint64{0, 1}, Prices: []uint64{0, 1}, Nonce: 1}
			midpointPrice.Signature, err = ecdsaKey.Sign(midpointPrice.Hash())
			Expect(err).ShouldNot(HaveOccurred())
			err = client.UpdateMidpoint(context.Background(), serviceMultiAddr, *midpointPrice)
			Expect(err).ShouldNot(HaveOccurred())

			midpointPrice = &oracle.MidpointPrice{TokenPairs: []uint64{1, 2, 3}, Prices: []uint64{2, 1, 5}, Nonce: 2}
			midpointPrice.Signature, err = ecdsaKey.Sign(midpointPrice.Hash())
			Expect(err).ShouldNot(HaveOccurred())
			err = client.UpdateMidpoint(context.Background(), serviceMultiAddr, *midpointPrice)
			Expect(err).ShouldNot(HaveOccurred())

			prices, err := midpointPriceStorer.MidpointPrice()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(prices.TokenPairs).Should(HaveLen(3))
			Expect(prices.TokenPairs).Should(Equal(midpointPrice.TokenPairs))
			Expect(prices.Prices).Should(HaveLen(3))
			Expect(prices.Prices).Should(Equal(midpointPrice.Prices))
			Expect(prices.Nonce).Should(Equal(midpointPrice.Nonce))
		})

		It("should not update the midpoint price when it receives a lower nonce", func() {
			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			var err error
			midpointPrice := &oracle.MidpointPrice{TokenPairs: []uint64{0, 1}, Prices: []uint64{0, 1}, Nonce: 2}
			midpointPrice.Signature, err = ecdsaKey.Sign(midpointPrice.Hash())
			Expect(err).ShouldNot(HaveOccurred())
			err = client.UpdateMidpoint(context.Background(), serviceMultiAddr, *midpointPrice)
			Expect(err).ShouldNot(HaveOccurred())

			// This midpoint price should not be stored as it has a lower
			// nonce.
			oldMidpointPrice := &oracle.MidpointPrice{TokenPairs: []uint64{1, 2, 3}, Prices: []uint64{2, 1, 5}, Nonce: 1}
			oldMidpointPrice.Signature, err = ecdsaKey.Sign(oldMidpointPrice.Hash())
			Expect(err).ShouldNot(HaveOccurred())
			err = client.UpdateMidpoint(context.Background(), serviceMultiAddr, *oldMidpointPrice)
			Expect(err).ShouldNot(HaveOccurred())

			prices, err := midpointPriceStorer.MidpointPrice()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(prices.TokenPairs).Should(HaveLen(2))
			Expect(prices.TokenPairs).Should(Equal(midpointPrice.TokenPairs))
			Expect(prices.Prices).Should(HaveLen(2))
			Expect(prices.Prices).Should(Equal(midpointPrice.Prices))
			Expect(prices.Nonce).Should(Equal(midpointPrice.Nonce))
		})
	})
})

func newOracleClient(db swarm.MultiAddressStorer) (oracle.Client, crypto.EcdsaKey, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, crypto.EcdsaKey{}, err
	}
	addr := identity.Address(ecdsaKey.Address())
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", addr))
	if err != nil {
		return nil, crypto.EcdsaKey{}, err
	}
	db.InsertMultiAddress(multiAddr)
	client := NewOracleClient(multiAddr.Address(), db)
	return client, ecdsaKey, nil
}
