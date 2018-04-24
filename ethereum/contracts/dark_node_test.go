package contracts_test

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/hyperdrive"
)

const (
	renContractAddress = "0xad6ab5ccbd2d761d11ba7e976ba7a93a6e3dd41a"
	dnrContractAddress = "0x429b5ba768e58f1a26b58742975aaeee417f3211"
	hyperdriveAddress  = "0x9db5820c2c5aa57cebe502727c98d952dae8e15f"
)

type Delta struct {
	BuyOrderID  []byte
	SellOrderID []byte
}

var _ = Describe("Darknode", func() {
	Context("watching for darknode contracts", func() {
		FIt("should watch the hyperdrive and finalized tx which has 16 block depth", func() {

			// create hyperdrive service
			conn, err := client.Connect("https://ropsten.infura.io", client.NetworkRopsten, renContractAddress, dnrContractAddress, hyperdriveAddress)
			Ω(err).ShouldNot(HaveOccurred())

			testKey, err := crypto.HexToECDSA("b44a49889a79983336d15385161533868644d35c1ea670854a0a0b4b784ae40c")
			Ω(err).ShouldNot(HaveOccurred())
			auth := bind.NewKeyedTransactor(testKey)

			//Create newHyperdriveContract for sending Txs
			hyper, err := contracts.NewHyperdriveContract(context.Background(), conn, auth, &bind.CallOpts{})
			Ω(err).ShouldNot(HaveOccurred())

			// Initialize other parameters for the test
			done, depth := make(chan struct{}), uint64(5)
			txInput := make(chan hyperdrive.TxWithBlockNumber)

			go func() {

				i := uint8(180)
				for {
					t := time.NewTimer(time.Second)
					select {
					case <-done:
						return
					case <-t.C:
						delta := Delta{
							BuyOrderID:  []byte{i},
							SellOrderID: []byte{i+1},
						}
						log.Printf("Found order match! Sending it to hyperdrive...")
						Ω(OrderMatchToHyperdrive(delta, hyper, txInput)).ShouldNot(HaveOccurred())
					}
					i += 2
				}
			}()

			go func() {
				// Quit the test after 5 minutes
				time.Sleep(10 * time.Minute)
				close(done)
			}()

			errs := WatchForHyperdriveContract(done, txInput, depth, hyper)
			for err := range errs {
				Ω(err).ShouldNot(HaveOccurred())
			}
		})
	})
})

// OrderMatchToHyperdrive converts an order match into a hyperdrive.Tx and
// forwards it to the Hyperdrive.
func OrderMatchToHyperdrive(delta Delta, hyper contracts.HyperdriveContract, txInput chan hyperdrive.TxWithBlockNumber) error {

	// Convert an order match into a Tx
	tx := hyperdrive.NewTxFromByteSlices(delta.SellOrderID, delta.BuyOrderID)

	transaction, err := hyper.SendTx(tx)
	if err != nil {
		return fmt.Errorf("fail to send tx to hyperdrive contract , %s", err)
	}

	blockNumber, err := hyper.GetBlockNumberOfTx(transaction.Hash())
	if err != nil {
		return fmt.Errorf("fail to get block number of the transaction , %s", err)
	}
	txInput <- hyperdrive.NewTxWithBlockNumber(transaction.Hash(), blockNumber)

	return nil
}

// Decouple the WatchForHyperdriveContract from the darknode so that we can
// do unit testing on it .
func WatchForHyperdriveContract(done <-chan struct{}, txInput chan hyperdrive.TxWithBlockNumber, depth uint64, hyper contracts.HyperdriveContract) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		watchingList := map[uint64][]common.Hash{}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case tx := <-txInput:
				if _, ok := watchingList[tx.BlockNumber]; !ok {
					watchingList[tx.BlockNumber] = []common.Hash{}
				}
				log.Printf("receive tx with block number %d", tx.BlockNumber)
				watchingList[tx.BlockNumber] = append(watchingList[tx.BlockNumber], tx.Hash)
			case <-ticker.C:
				currentBlock, err := hyper.CurrentBlock()
				log.Println("Current block is ", currentBlock.NumberU64())
				if err != nil {
					errs <- err
					return
				}

				for key, value := range watchingList {
					if key <= currentBlock.NumberU64()-depth {
						for _, hash := range value {
							// Check if there is a block shuffle
							newBlockNumber, err := hyper.GetBlockNumberOfTx(hash)
							log.Println("new block number is ", newBlockNumber)
							if err != nil {
								errs <- err
								return
							}
							if newBlockNumber != key {
								if _, ok := watchingList[newBlockNumber]; !ok {
									watchingList[newBlockNumber] = []common.Hash{}
								}
								// "If map entries are created during iteration, that entry may be produced during the iteration or may be skipped."
								watchingList[newBlockNumber] = append(watchingList[newBlockNumber], hash)
								continue
							}

							// Create a hashTable
							hashTable := map[common.Hash]struct{}{}
							block, err := hyper.BlockByNumber(big.NewInt(int64(key)))
							if err != nil {
								errs <- err
								return
							}
							for _, h := range block.Transactions() {
								hashTable[h.Hash()] = struct{}{}
							}

							if _, ok := hashTable[hash]; ok {
								log.Println(hash.Hex(), "has been finalized in block ", key)
							}
						}
						delete(watchingList, key)
					}
				}
			}
		}
	}()

	return errs
}
