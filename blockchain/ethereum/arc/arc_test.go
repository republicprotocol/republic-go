package arc_test

import (
	"context"
	"crypto/sha256"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
	"github.com/republicprotocol/republic-go/blockchain/test"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

var _ = Describe("ether", func() {

	var conn ethereum.Conn

	var alice, bob *bind.TransactOpts
	var aliceAddr, bobAddr common.Address
	var orderID [32]byte
	var value *big.Int
	var validity int64

	BeforeEach(func() {
		var err error

		conn, err = ganache.StartAndConnect()
		Expect(err).ShouldNot(HaveOccurred())

		alice, aliceAddr, err = ganache.NewAccount(conn, big.NewInt(10))
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000

		bob, bobAddr, err = ganache.NewAccount(conn, big.NewInt(10))
		Expect(err).ShouldNot(HaveOccurred())
		bob.GasLimit = 3000000

		value = big.NewInt(10)
		validity = int64(time.Hour * 24)

		orderID[0] = 0x13
	})

	AfterEach(func() {
		ganache.Stop()
	})

	test.SkipCIContext("when using Ethereum", func() {

		It("should successfully perform Ether to Ether atomic swaps", func() {

			var aliceArcData, bobArcData []byte

			var aliceSecret [32]byte
			var secretHash [32]byte

			{ // Alice can initiate swap
				aliceArc, err := NewArc(context.Background(), conn, alice, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())

				aliceSecret = [32]byte{1, 3, 3, 7}
				secretHash = sha256.Sum256(aliceSecret[:])
				err = aliceArc.Initiate(secretHash, aliceAddr.Bytes(), bobAddr.Bytes(), value, validity)
				Expect(err).ShouldNot(HaveOccurred())
				aliceArcData, err = aliceArc.Serialize()
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Bob can audit Alice's contract and upload his own
				aliceArc, err := NewArc(context.Background(), conn, bob, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())
				err = aliceArc.Deserialize(aliceArcData)
				Expect(err).ShouldNot(HaveOccurred())

				_secretHash, _, to, _value, _, err := aliceArc.Audit(orderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(to).Should(Equal(bobAddr.Bytes()))
				Expect(_value).Should(Equal(value))
				// Expect(_expiry.Int64()).Should(Equal(expiry))

				bobArc, err := NewArc(context.Background(), conn, bob, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())
				err = bobArc.Initiate(_secretHash, bobAddr.Bytes(), aliceAddr.Bytes(), value, validity)
				Expect(err).ShouldNot(HaveOccurred())
				bobArcData, err = bobArc.Serialize()
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Alice can audit Bob's contract and reveal the secret
				bobArc, err := NewArc(context.Background(), conn, alice, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())
				err = bobArc.Deserialize(bobArcData)
				Expect(err).ShouldNot(HaveOccurred())

				_secretHash, _, to, _value, _, err := bobArc.Audit(orderID)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(_secretHash).Should(Equal(secretHash))
				Expect(to).Should(Equal(aliceAddr.Bytes()))
				Expect(to).Should(Equal(aliceAddr.Bytes()))
				Expect(_value).Should(Equal(value))

				err = bobArc.Redeem(orderID, aliceSecret)
				Expect(err).ShouldNot(HaveOccurred())
			}

			{ // Bob can retrieve the secret from his contract and complete the swap
				bobArc, err := NewArc(context.Background(), conn, bob, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())
				err = bobArc.Deserialize(bobArcData)
				Expect(err).ShouldNot(HaveOccurred())

				secret, err := bobArc.AuditSecret(orderID)
				Expect(err).ShouldNot(HaveOccurred())

				aliceArc, err := NewArc(context.Background(), conn, bob, orderID[:], ETHEREUM, big.NewInt(0))
				Expect(err).ShouldNot(HaveOccurred())
				err = aliceArc.Deserialize(aliceArcData)
				Expect(err).ShouldNot(HaveOccurred())

				err = aliceArc.Redeem(orderID, secret)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
	})
})
