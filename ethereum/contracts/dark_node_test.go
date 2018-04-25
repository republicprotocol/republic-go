package contracts_test

import (
	"context"
	"fmt"
	"log"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
	"github.com/republicprotocol/republic-go/hyperdrive"
	"github.com/republicprotocol/republic-go/identity"
)

const (
	renContractAddress = "0xad6ab5ccbd2d761d11ba7e976ba7a93a6e3dd41a"
	dnrContractAddress = "0x429b5ba768e58f1a26b58742975aaeee417f3211"
	hyperdriveAddress  = "0x348496ad820f2ee256268f9f9d0b9f5bacdc26cd"
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
			txInput := make(chan hyperdrive.TxWithTimestamp)

			go func() {
				// Quit the test after 5 minutes
				time.Sleep(10 * time.Minute)
				close(done)
			}()

			go WatchForHyperdriveContract(done, txInput, depth, hyper)

			i := uint8(80)
			for {
				t := time.NewTimer(time.Second)
				select {
				case <-done:
					return
				case <-t.C:
					delta := Delta{
						BuyOrderID:  []byte{i},
						SellOrderID: []byte{i + 1},
					}
					log.Printf("Found order match! Sending it to hyperdrive...")
					err = OrderMatchToHyperdrive(delta, hyper, txInput)
					Ω(err).ShouldNot(HaveOccurred())
				}
				i += 2
			}
		})
	})
})

// OrderMatchToHyperdrive converts an order match into a hyperdrive.Tx and
// forwards it to the Hyperdrive.
func OrderMatchToHyperdrive(delta Delta, hyper contracts.HyperdriveContract, txInput chan hyperdrive.TxWithTimestamp) error {

	// Convert an order match into a Tx
	tx := hyperdrive.NewTxFromByteSlices(delta.SellOrderID, delta.BuyOrderID)

	_, err := hyper.SendTx(tx)

	if err != nil {
		return fmt.Errorf("fail to send tx to hyperdrive contract , %s", err)
	}

	txInput <- hyperdrive.NewTxWithTimestamp(tx, time.Now())

	return nil
}

// Decouple the WatchForHyperdriveContract from the darknode so that we can
// do unit testing on it .
func WatchForHyperdriveContract(done <-chan struct{}, txInput chan hyperdrive.TxWithTimestamp, depth uint64, hyper contracts.HyperdriveContract) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		watchingList := map[identity.Hash]hyperdrive.TxWithTimestamp{}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case tx := <-txInput:
				log.Printf("receive tx with %v  at  ", tx.Hash)
				watchingList[tx.Tx.Hash] = tx
			case <-ticker.C:
				log.Println("tik tok ......")
				for key, tx := range watchingList {
					if time.Now().Before(tx.Timestamp.Add(5 * time.Minute)) {
						finalized := true
						for _, nonce := range tx.Nonces {
							dep, err := hyper.GetDepth(nonce)
							if err != nil {
								log.Println("fail to get the depth of the transaction. ")
								finalized = false
								delete(watchingList, key)
								break
							}
							if dep < depth {
								finalized = false
								break
							}
						}

						if finalized {
							delete(watchingList, key)
							log.Println(tx.Hash, "has been finalized in block ")
						}
					} else {
						log.Println("time expire")
						delete(watchingList, key)
					}
				}
				log.Println("number of elements in the map ", len(watchingList))
			}
		}
	}()

	return errs
}
