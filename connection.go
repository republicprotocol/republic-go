package go_eth

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/republicprotocol/go-eth/atomic_swap"
)

// Connection ...
type Connection interface {
	Open(ethAddr common.Address, ethAmount uint64, secretHash [32]byte) (id [32]byte, err error)
	Check(id [32]byte) (struct {
		TimeRemaining  *big.Int
		Value          *big.Int
		WithdrawTrader common.Address
		SecretLock     [32]byte
	}, error)
}

// EtherConnection ...
type EtherConnection struct {
	client   *backends.SimulatedBackend
	auth     *bind.TransactOpts
	contract *atomic_swap.AtomicSwapEther
}

func randomAuth() *bind.TransactOpts {
	// Generate a new random account
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	return auth
}

// Simulated ...
func Simulated(auth1 *bind.TransactOpts, auth2 *bind.TransactOpts) *backends.SimulatedBackend {
	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth1.From: {Balance: big.NewInt(10000000000)}})
	return sim
}

// func ropsten() (*ethclient.Client, *bind.TransactOpts) {
// 	// Create an IPC based RPC connection to a remote node and an authorized transactor
// 	conn, err := ethclient.Dial("http://127.0.0.1:8545")
// 	if err != nil {
// 		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
// 	}
// 	auth, err := bind.NewTransactor(strings.NewReader(key), "M#171949298ac29dcce331029254cb77ce234d2613")
// 	if err != nil {
// 		log.Fatalf("Failed to create authorized transactor: %v", err)
// 	}
// 	return conn, auth
// }

// DeployEther ...
func DeployEther(connection *backends.SimulatedBackend, auth *bind.TransactOpts) common.Address {
	// Deploy a token contract on the simulated blockchain
	address, _, _, err := atomic_swap.DeployAtomicSwapEther(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node
	return address
}

func existing(connection *backends.SimulatedBackend, address common.Address) *atomic_swap.AtomicSwapEther {
	contract, err := atomic_swap.NewAtomicSwapEther(address, connection)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return contract
}

// NewEtherConnection ...
func NewEtherConnection(client *backends.SimulatedBackend, auth1 *bind.TransactOpts, address common.Address) Connection {
	contract := existing(client, address)

	return EtherConnection{
		client:   client,
		auth:     auth1,
		contract: contract,
	}
}
