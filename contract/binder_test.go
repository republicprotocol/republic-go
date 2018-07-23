package contract_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/ome"
	"github.com/republicprotocol/republic-go/stackint"
	"github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Contract Binder", func() {

	_, binder, _ := testutils.GanacheBeforeSuite(func() {
	})

	testutils.GanacheAfterSuite(func() {
	})

	testutils.GanacheContext("when interacting with ethereum smart contracts", func() {

		It("should not return a nonce error", func() {
			numberOfDarknodes := 24
			numberOfOrderPairs := 10

			/********************************************************/
			/* Testing registration of darknodes                    */
			/********************************************************/

			keystores := make([]crypto.Keystore, numberOfDarknodes)
			for i := 0; i < numberOfDarknodes; i++ {
				// Bond for the darknode
				bond, err := stackint.FromString("0")
				Expect(err).ShouldNot(HaveOccurred())

				// Darknode address
				keystores[i], err = crypto.RandomKeystore()
				Expect(err).ShouldNot(HaveOccurred())
				darknodeAddr := identity.Address(keystores[i].Address())

				publicKey, err := crypto.BytesFromRsaPublicKey(&(keystores[i].RsaKey.PublicKey))

				// Register darknode with the darknodeRegistry
				err = binder.Register(darknodeAddr.ID(), publicKey, &bond)
				Expect(err).ShouldNot(HaveOccurred())

				// Darknode will be waiting to be registered until a new
				// epoch
				tx, err := binder.IsRegistered(darknodeAddr)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tx).To(BeFalse())
			}

			// Trigger new Epoch. This will complete registration of
			// the darknode
			_, err := binder.NextEpoch()
			if numberOfDarknodes >= 24 {
				Expect(err).ShouldNot(HaveOccurred())
			}

			for i := 0; i < numberOfDarknodes; i++ {
				darknodeAddr := identity.Address(keystores[i].Address())

				// Darknode should be registered
				tx, err := binder.IsRegistered(darknodeAddr)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tx).To(BeTrue())
			}

			/********************************************************/
			/* Testing order matching                               */
			/********************************************************/

			orderpairs := make([]ome.Computation, numberOfOrderPairs)

			for i := 0; i < numberOfOrderPairs; i++ {
				orderpairs[i], err = testutils.RandomComputation()
				Expect(err).ShouldNot(HaveOccurred())

				err := binder.OpenBuyOrder([65]byte{}, orderpairs[i].Buy.OrderID)
				Expect(err).ShouldNot(HaveOccurred())

				err = binder.OpenSellOrder([65]byte{}, orderpairs[i].Sell.OrderID)
				Expect(err).ShouldNot(HaveOccurred())
			}

			numberOfOrders, err := binder.OrderCounts()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(numberOfOrders).To(Equal(uint64(2 * numberOfOrderPairs)))

			// // Create a transactOpts for a darknode to submit order confirmations
			// auth := bind.NewKeyedTransactor(keystores[0].EcdsaKey.PrivateKey)

			// // Transfer funds into the darknode
			// transOpts := ganache.GenesisTransactor()
			// value := big.NewInt(10)
			// value.Exp(value, big.NewInt(18), nil)
			// value.Mul(big.NewInt(10), value)
			// conn.TransferEth(context.Background(), &transOpts, auth.From, value)

			// // Get binder for the darknode
			// darknodeBinder, err := contract.NewBinder(auth, conn)
			// if err != nil {
			// 	log.Fatalf("cannot get ethereum bindings: %v", err)
			// }

			// for i := 0; i < numberOfOrderPairs; i++ {
			// 	err = darknodeBinder.ConfirmOrder(orderpairs[i].Buy, orderpairs[i].Sell)
			// 	Expect(err).ShouldNot(HaveOccurred())
			// }

			// numberOfOrders, err = binder.OrderCounts()
			// Expect(err).ShouldNot(HaveOccurred())
			// Expect(numberOfOrders).To(Equal(uint64(2 * numberOfOrderPairs)))

			/********************************************************/
			/* Testing deregistration of darknodes                  */
			/********************************************************/

			for i := 0; i < numberOfDarknodes; i++ {
				darknodeAddr := identity.Address(keystores[i].Address())
				// Deregister darknode
				err = binder.Deregister(darknodeAddr.ID())
				Expect(err).ShouldNot(HaveOccurred())
				// Darknode will be waiting to be deregistered until a
				// new epoch
				tx, err := binder.IsDeregistered(darknodeAddr.ID())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tx).To(BeFalse())
			}

			// Trigger new Epoch. This will complete deregistration of
			// the darknode
			_, err = binder.NextEpoch()
			Expect(err).Should(HaveOccurred())

			for i := 0; i < numberOfDarknodes; i++ {
				darknodeAddr := identity.Address(keystores[i].Address())

				// Darknode should be deregistered
				tx, err := binder.IsDeregistered(darknodeAddr.ID())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(tx).To(BeTrue())

				// Deregister the same darknode should return an error
				err = binder.Deregister(darknodeAddr.ID())
				Expect(err).Should(HaveOccurred())

				// Refund deregistered node should return an error since
				// the darknode had a bond amount of 0
				err = binder.Refund(darknodeAddr.ID())
				Expect(err).Should(HaveOccurred())
			}
		})
	})
})
