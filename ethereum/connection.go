package ethereum

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
)

// EtherConnection ...
type EtherConnection struct {
	client   bind.ContractBackend
	auth     *bind.TransactOpts
	contract *contracts.AtomicSwapEther
}

// ERC20Connection ...
type ERC20Connection struct {
	client   bind.ContractBackend
	auth     *bind.TransactOpts
	contract *contracts.AtomicSwapERC20
}

func randomAuth() *bind.TransactOpts {
	// Generate a new random account
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	return auth
}

// Ropsten ...
func Ropsten(uri string) *ethclient.Client {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial(uri)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	return conn
}

// Simulated ...
func Simulated(auth1 *bind.TransactOpts, auth2 *bind.TransactOpts) *backends.SimulatedBackend {
	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth1.From: {Balance: big.NewInt(10000000000)}, auth2.From: {Balance: big.NewInt(10000000000)}})
	return sim
}

// DeployEther ...
func DeployEther(connection *backends.SimulatedBackend, auth *bind.TransactOpts) common.Address {
	// Deploy a token contract on the simulated blockchain
	address, _, _, err := contracts.DeployAtomicSwapEther(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node
	return address
}

func existingEth(connection bind.ContractBackend, address common.Address) *contracts.AtomicSwapEther {
	contract, err := contracts.NewAtomicSwapEther(address, connection)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return contract
}

func existingERC20(connection bind.ContractBackend, address common.Address) *contracts.AtomicSwapERC20 {
	contract, err := contracts.NewAtomicSwapERC20(address, connection)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return contract
}
