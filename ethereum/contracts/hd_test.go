package contracts

import (
	"context"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/ethereum/ganache"
	"github.com/republicprotocol/republic-go/hyperdrive"
)

var _ = Describe("Dark Node Registrar", func() {
	Context("Send tx to hyperdrive contract", func() {
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
			darknodeRegistry, err := NewDarkNodeRegistry(context.Background(), conn, auth, &bind.CallOpts{})
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
				_, err = darknodeRegistry.WaitForEpoch()
				Ω(err).ShouldNot(HaveOccurred())
			}

			//Create newHyperdriveContract for sending Txs
			hyper, err := NewHyperdriveContract(context.Background(), conn, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			hyper.SetGasLimit(1000000)
			tx := hyperdrive.Tx{
				Nonces: [][32]byte{
					{0, 1},
				},
			}
			hyperTransaction, err := hyper.SendTx(tx)
			Ω(err).ShouldNot(HaveOccurred())

			_, err = conn.PatchedWaitMined(context.Background(), hyperTransaction)
			Ω(err).ShouldNot(HaveOccurred())

			log.Println(hyperTransaction.Hash().Hex())

			receipt, err := hyper.GetBlockNumberOfTx(hyperTransaction)
			Ω(err).ShouldNot(HaveOccurred())
			log.Println("block number of the transaction is ", receipt)
		})
	})
})
