package hd_test

import (
	"context"
	"log"

	"github.com/jbenet/go-base58"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/blockchain/ethereum/hd"
	"github.com/republicprotocol/republic-go/blockchain/test"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
)

var _ = Describe("hyperdrive", func() {
	test.SkipCIContext("Local test with ganache", func() {

		It("should be able to send txs which have no conflicts", func() {
			// Connect to local ganache blockchain
			conn, err := ganache.Connect("http://localhost:8545")
			Ω(err).ShouldNot(HaveOccurred())

			// Create new transactor
			ethereumPair, err := crypto.GenerateKey()
			ethereumKey := &keystore.Key{
				Address:    crypto.PubkeyToAddress(ethereumPair.PublicKey),
				PrivateKey: ethereumPair,
			}
			auth := bind.NewKeyedTransactor(ethereumKey.PrivateKey)

			// Distribute ren and eth to the address
			err = ganache.DistributeREN(conn, auth.From)
			Ω(err).ShouldNot(HaveOccurred())

			err = ganache.DistributeEth(conn, auth.From)
			Ω(err).ShouldNot(HaveOccurred())

			// Register the account
			darknodeRegistry, err := dnr.NewDarknodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			darknodeRegistry.SetGasLimit(1000000)
			minimumBond, err := darknodeRegistry.MinimumBond()
			Ω(err).ShouldNot(HaveOccurred())

			isRegistered, err := darknodeRegistry.IsRegistered(auth.From.Bytes())
			Ω(err).ShouldNot(HaveOccurred())
			if !isRegistered {
				transaction, err := darknodeRegistry.Register(auth.From.Bytes(), []byte{}, &minimumBond)
				Ω(err).ShouldNot(HaveOccurred())
				_, err = conn.PatchedWaitMined(context.Background(), transaction)
				Ω(err).ShouldNot(HaveOccurred())
				err = darknodeRegistry.WaitForEpoch()
				Ω(err).ShouldNot(HaveOccurred())
			}

			//Create newHyperdriveContract for sending Txs
			hyper, err := NewHyperdriveContract(context.Background(), conn, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())
			hyper.SetGasLimit(1000000)

			tx := Tx{
				Nonces: [][]byte{
					{0, 1},
				},
			}
			transaction, err := hyper.SendTx(tx)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(transaction).ShouldNot(BeNil())
		})
	})

	test.SkipCIContext("Tests with Ropsten", func() {

		BeforeEach(func() {

		})

		It("should be able to send txs with no conflicts", func() {

			config := ethereum.Config{
				Network:                 ethereum.NetworkRopsten,
				URI:                     "https://ropsten.infura.io",
				RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.Hex(),
				DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.Hex(),
				HyperdriveAddress:       ethereum.HyperdriveAddressOnRopsten.Hex(),
				ArcAddress:              ethereum.ArcAddressOnRopsten.Hex(),
			}
			conn, err := ethereum.Connect(config)
			Ω(err).ShouldNot(HaveOccurred())

			//conn, err := ethereum.Connect("https://ropsten.infura.io", ethereum.NetworkRopsten, renContractAddress, dnrContractAddress, hyperdriveAddress)
			//Ω(err).ShouldNot(HaveOccurred())

			testKey, err := crypto.HexToECDSA("f8421de8bcddfd340346f979f22547cccd01fc436becf542c7ea971866a58949")
			Ω(err).ShouldNot(HaveOccurred())
			auth := bind.NewKeyedTransactor(testKey)

			//Create newHyperdriveContract for sending Txs
			hyper, err := NewHyperdriveContract(context.Background(), conn, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			nonce := base58.Decode("FEMmnsQwwNsmV1MArbFgdSwQqQq6VwJMdqc9VaRgLjUA")
			blockNumber, err := hyper.CheckOrders(nonce)
			Ω(blockNumber).Should(BeZero())

			tx := Tx{
				Nonces: [][]byte{
					{7}, // Make sure you increment this number before running the test
				},
			}
			transaction, err := hyper.SendTx(tx)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(transaction).ShouldNot(BeNil())
			log.Println("transaction is ", transaction)

			//blockNumber, err := hyper.GetBlockNumberOfTx(transaction.Hash())
			//Ω(err).ShouldNot(HaveOccurred())
			//Ω(blockNumber).Should(BeNumerically(">", 3097093))
			// log.Println("blockNumber is ", blockNumber)
		})

	})
})
