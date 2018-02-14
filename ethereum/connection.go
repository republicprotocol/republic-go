package ethereum

import (
	"context"
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
	log15 "github.com/ethereum/go-ethereum/log"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
)

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

/** SIMULATED FOR TESTS */
func randomAuth() *bind.TransactOpts {
	// Generate a new random account
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)
	return auth
}

// Simulated ...
func Simulated(auth1 *bind.TransactOpts, auth2 *bind.TransactOpts) *backends.SimulatedBackend {
	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth1.From: {Balance: big.NewInt(10000000000)}, auth2.From: {Balance: big.NewInt(10000000000)}})
	return sim
}

// DeployETH ...
func DeployETH(connection *backends.SimulatedBackend, auth *bind.TransactOpts) common.Address {
	// Deploy a token contract on the simulated blockchain
	address, _, _, err := contracts.DeployAtomicSwapEther(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	// Don't even wait, check its presence in the local pending state
	time.Sleep(250 * time.Millisecond) // Allow it to be processed by the local node
	return address
}

// Go-ethereum's WaitMined is not compatible with Parity's getTransactionReceipt
// WaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
func PatchedWaitMined(ctx context.Context, b bind.DeployBackend, tx *types.Transaction) (*types.Receipt, error) {
	queryTicker := time.NewTicker(time.Second)
	defer queryTicker.Stop()

	logger := log15.New("hash", tx.Hash())
	for {
		receipt, err := b.TransactionReceipt(ctx, tx.Hash())
		if receipt != nil && receipt.Status != 0 {
			return receipt, nil
		}
		if err != nil {
			logger.Trace("Receipt retrieval failed", "err", err)
		} else {
			logger.Trace("Transaction not yet mined")
		}
		// Wait for the next round.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-queryTicker.C:
		}
	}
}
