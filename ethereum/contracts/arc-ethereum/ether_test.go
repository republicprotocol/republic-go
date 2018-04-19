package arc_ethereum_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts/arc-ethereum"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
)

var _ = Describe("ether", func() {
	var alice, bob *bind.TransactOpts
	var aliceAddr, bobAddr common.Address
	var conn client.Connection
	var swapID [32]byte
	var value *big.Int
	var validity int64

	BeforeEach(func() {
		// Setup...

		var err error
		conn, err = ganache.Connect("http://localhost:8545")
		Expect(err).ShouldNot(HaveOccurred())

		alice, aliceAddr, err = ganache.NewAccount(conn, big.NewInt(10))
		Expect(err).ShouldNot(HaveOccurred())
		alice.GasLimit = 3000000
		bob, bobAddr, err = ganache.NewAccount(conn, big.NewInt(10))
		Expect(err).ShouldNot(HaveOccurred())

		fmt.Printf("Alice: %v\n", aliceAddr)
		fmt.Printf("Bobsa: %v\n", bobAddr)
		Expect(common.BytesToAddress(aliceAddr.Bytes())).Should(Equal(aliceAddr))

		value = big.NewInt(10)
		validity = int64(time.Hour * 24)
	})

	It("can perform ETH-ETH arc swap", func() {

		var aliceArcAddress, bobArcAddress common.Address

		var aliceSecret []byte
		var secretHash [32]byte

		{ // Alice can initiate swap
			aliceArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, alice, common.Address{}, swapID)
			Expect(err).ShouldNot(HaveOccurred())

			aliceSecret := []byte("the secret")
			secretHash = sha256.Sum256(aliceSecret)
			fmt.Printf("secret: %v\n", secretHash)
			err = aliceArc.Initiate(secretHash, bobAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			aliceArcAddress = aliceArc.EthereumArcData.ContractAddress
		}

		{ // Bob can audit Alice's contract and upload his own
			aliceArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, bob, aliceArcAddress, swapID)
			Expect(err).ShouldNot(HaveOccurred())

			_secretHash, _, to, _value, _, err := aliceArc.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(to).Should(Equal(bobAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			// Expect(_expiry.Int64()).Should(Equal(expiry))

			bobArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, bob, common.Address{}, swapID)
			Expect(err).ShouldNot(HaveOccurred())
			err = bobArc.Initiate(_secretHash, aliceAddr.Bytes(), value, validity)
			Expect(err).ShouldNot(HaveOccurred())
			bobArcAddress = bobArc.EthereumArcData.ContractAddress
		}

		{ // Alice can audit Bob's contract and reveal the secret
			bobArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, alice, bobArcAddress, swapID)
			Expect(err).ShouldNot(HaveOccurred())

			secretHash := sha256.Sum256(aliceSecret[:])

			_secretHash, _, to, _value, _, err := bobArc.Audit()
			Expect(err).ShouldNot(HaveOccurred())
			Expect(_secretHash).Should(Equal(secretHash))
			Expect(to).Should(Equal(aliceAddr.Bytes()))
			Expect(to).Should(Equal(aliceAddr.Bytes()))
			Expect(_value).Should(Equal(value))
			// Expect(_expiry.Int64()).Should(Equal(expiry))

			err = bobArc.Redeem(aliceSecret)
			Expect(err).ShouldNot(HaveOccurred())
		}

		{ // Bob can retrieve the secret from his contract and complete the swap
			bobArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, bob, bobArcAddress, swapID)
			Expect(err).ShouldNot(HaveOccurred())

			secret, err := bobArc.AuditSecret()
			Expect(err).ShouldNot(HaveOccurred())

			aliceArc, err := arc_ethereum.NewEthereumArc(context.Background(), conn, bob, aliceArcAddress, swapID)
			Expect(err).ShouldNot(HaveOccurred())

			err = aliceArc.Redeem(secret)
			Expect(err).ShouldNot(HaveOccurred())
		}
	})
})
