package main

import (
	"context"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/go-eth/atomic_swap"
)

// The Delegate is used as a callback interface to inject logic into the
// different RPCs.
type Delegate interface {
}

// Node implements the gRPC Node service.
type Node struct {
	Delegate
}

// Open opens an order
func Open(ethAddr common.Address, btcAddr string, ethAmount uint64) {
	conn, auth := simulated()
	contract := existing(conn, auth)

	secret := [20]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	tx, _ := contract.Initiate(auth, big.NewInt(10), secret, ethAddr)
	println(tx.Hash)
}

// Close closes an order
func Close() {

}

// Expire expires an order
func Expire() {

}

// Validate validates an order
func Validate() {

}

// RetrieveSecretKey retrieves a secret key from an order
func RetrieveSecretKey() {

}

func simulated() (*backends.SimulatedBackend, *bind.TransactOpts) {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth.From: {Balance: big.NewInt(10000000000)}})
	return sim, auth
}

const key = `***REMOVED***`

func ropsten() (*ethclient.Client, *bind.TransactOpts) {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	auth, err := bind.NewTransactor(strings.NewReader(key), "***REMOVED***")
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	return conn, auth
}

func deploy(connection bind.ContractBackend, auth *bind.TransactOpts) *atomic_swap.Atomic_swap {
	// Deploy a token contract on the simulated blockchain
	_, _, contract, err := atomic_swap.DeployAtomic_swap(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node :P
	return contract
}

func existing(connection bind.ContractBackend, auth *bind.TransactOpts) *atomic_swap.Atomic_swap {
	contract, err := atomic_swap.NewAtomic_swap(common.StringToAddress("0x8a4d144fbced56e47c1d3a5a0d5f9c35bec354f8"), connection)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return contract
}

func main() {
	conn, auth := ropsten()
	// deploy(conn, auth)

	ctx := context.Background()
	bal, err := conn.BalanceAt(ctx, auth.From, nil)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	println(auth.From.String())
	println(bal.String())

	Open(auth.From, ``, 0)
}
