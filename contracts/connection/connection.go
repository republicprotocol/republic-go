package connection

import (
	"context"
	"math/big"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client combines the interfaces for bind.ContractBackend and bind.DeployBackend
type Client interface {
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
}

// ClientDetails contains the simulated client and the contracts deployed to it
type ClientDetails struct {
	Client     Client
	RenAddress common.Address
	DNRAddress common.Address
}

// FromURI will connect to a provided RPC uri
func FromURI(uri string) (ClientDetails, error) {
	if uri == "" {
		uri = "https://ropsten.infura.io/"
	}
	client, err := ethclient.Dial(uri)
	if err != nil {
		return ClientDetails{}, err
	}
	return ClientDetails{
		Client:     client,
		RenAddress: common.HexToAddress("0x889debfe1478971bcff387f652559ae1e0b6d34a"),
		DNRAddress: common.HexToAddress("0x6e48bdd8949d0c929e9b5935841f6ff18de0e613"),
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// If the client is a simulated backend, it will commit the pending transactions to a block
func PatchedWaitMined(ctx context.Context, b Client, tx *types.Transaction) (*types.Receipt, error) {

	sim, ok := b.(*backends.SimulatedBackend)
	if ok {
		sim.Commit()
		sim.AdjustTime(10 * time.Second)
	}

	return bind.WaitMined(ctx, b, tx)

	// (Go-ethereum's WaitMined is not compatible with Parity's getTransactionReceipt)
	// NOTE: If something goes wrong, this will hang!

	// queryTicker := time.NewTicker(time.Second)
	// defer queryTicker.Stop()

	// logger := log15.New("hash", tx.Hash())
	// for {
	// 	receipt, err := b.TransactionReceipt(ctx, tx.Hash())
	// 	if receipt != nil && receipt.Status != 0 {
	// 		return receipt, nil
	// 	}
	// 	if err != nil {
	// 		logger.Trace("Receipt retrieval failed", "err", err)
	// 	} else {
	// 		logger.Trace("Transaction not yet mined")
	// 	}
	// 	// Wait for the next round.
	// 	select {
	// 	case <-ctx.Done():
	// 		return nil, ctx.Err()
	// 	case <-queryTicker.C:
	// 	}
	// }
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
// If the client is a simulated backend, it will commit the pending transactions to a block
func PatchedWaitDeployed(ctx context.Context, b Client, tx *types.Transaction) (common.Address, error) {

	sim, ok := b.(*backends.SimulatedBackend)
	if ok {
		sim.Commit()
		sim.AdjustTime(10 * time.Second)
	}

	return bind.WaitDeployed(ctx, b, tx)

	// (Go-ethereum's WaitMined is not compatible with Parity's getTransactionReceipt)
	// NOTE: If something goes wrong, this will hang!

	// if tx.To() != nil {
	// 	return common.Address{}, fmt.Errorf("tx is not contract creation")
	// }
	// receipt, err := PatchedWaitMined(ctx, b, tx)
	// if err != nil {
	// 	return common.Address{}, err
	// }
	// if receipt.ContractAddress == (common.Address{}) {
	// 	return common.Address{}, fmt.Errorf("zero address")
	// }
	// // Check that code has indeed been deployed at the address.
	// // This matters on pre-Homestead chains: OOG in the constructor
	// // could leave an empty account behind.
	// code, err := b.CodeAt(ctx, receipt.ContractAddress, nil)
	// if err == nil && len(code) == 0 {
	// 	err = errors.New("no contract code after deployment")
	// }
	// return receipt.ContractAddress, err
}
