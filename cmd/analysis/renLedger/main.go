package main

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/ledger"
	"github.com/republicprotocol/republic-go/order"
)

func main() {
	key, err := LoadKey()
	if err != nil {
		log.Fatal(err)
	}
	auth := bind.NewKeyedTransactor(key.PrivateKey)
	config := ethereum.Config{
		Network:                 ethereum.NetworkKovan,
		URI:                     "https://kovan.infura.io",
		RepublicTokenAddress:    ethereum.RepublicTokenAddressOnKovan.String(),
		DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnKovan.String(),
	}
	conn, err := ethereum.Connect(config)
	if err != nil {
		log.Fatalf("cannot connect to ethereum: %v", err)
	}
	renLedger, err := ledger.NewRenLedgerContract(context.Background(), conn, auth, &bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	counts, err := renLedger.OrderCounts()
	if err != nil {
		log.Fatal(err)
	}
	openCounts, confirmedCount := 0, 0
	for i := 0; uint64(i) < counts; i++ {
		orderID, err := renLedger.OrderID(i)
		if err != nil {
			log.Fatal(err)
		}
		hexID := hexutil.Encode(orderID[:])
		status, err := renLedger.Status(orderID)
		switch status {
		case order.Open:
			openCounts++
		case order.Confirmed:
			confirmedCount++
		default:

		}
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%d | %v | %v | %v", i, hexID, status, base64.StdEncoding.EncodeToString(orderID[:]))
	}

	log.Printf("Open Orders : %v, Confirmed Orders: %v, base64 : ", openCounts, confirmedCount)

}

func LoadKey() (*keystore.Key, error) {
	var keyJSON string = `{"address":"90e6572ef66a11690b09dd594a18f36cf76055c8",
  					"privatekey":"dc3f937b4aa1fc7bbf7643f1dead1faf37594ad2f1edcd6b56bf6719f85fa406",
  					"id":"ddd54c1c-6c2e-42a9-a224-6532a90fd4e9", "version":3}`
	key := new(keystore.Key)
	err := key.UnmarshalJSON([]byte(keyJSON))

	return key, err
}
