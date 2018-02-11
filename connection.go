package go_eth

import (
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/go-eth/contracts"
)

// Connection ...
type Connection interface {
	Open(_swapID [32]byte, ethAddr common.Address, ethAmount uint64, secretHash [32]byte, amountInWei *big.Int) (*types.Transaction, error)
	Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error)
	RetrieveSecretKey(_swapID [32]byte) ([]byte, error)
	Expire(_swapID [32]byte) (*types.Transaction, error)
	Validate()
	// GetState(_swapID [32]byte) (uint8, error)
	Check(id [32]byte) (struct {
		TimeRemaining  *big.Int
		Value          *big.Int
		WithdrawTrader common.Address
		SecretLock     [32]byte
	}, error)
}

// EtherConnection ...
type EtherConnection struct {
	client   bind.ContractBackend
	auth     *bind.TransactOpts
	contract *contracts.AtomicSwapEther
}

func randomAuth() *bind.TransactOpts {
	// Generate a new random account
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	return auth
}

// Ropsten ...
func Ropsten() *ethclient.Client {
	// Create an IPC based RPC connection to a remote node and an authorized transactor
	conn, err := ethclient.Dial("http://13.54.129.55:8545")
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

func existing(connection bind.ContractBackend, address common.Address) *contracts.AtomicSwapEther {
	contract, err := contracts.NewAtomicSwapEther(address, connection)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return contract
}

// NewEtherConnection ...
func NewEtherConnection(client bind.ContractBackend, auth1 *bind.TransactOpts, address common.Address) Connection {
	contract := existing(client, address)

	return EtherConnection{
		client:   client,
		auth:     auth1,
		contract: contract,
	}
}
