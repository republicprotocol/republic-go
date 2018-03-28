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

// Chain is used to represent an Ethereum chain
type Chain string

const (
	// ChainMainnet represents the Ethereum mainnet
	ChainMainnet Chain = "mainnet"
	// ChainRopsten represents the Ethereum Ropsten testnet
	ChainRopsten Chain = "ropsten"
	// ChainGanache represents a Ganache testrpc server
	ChainGanache Chain = "ganache"
	// ChainSimulated represents a go-ethereum simulated backend
	ChainSimulated Chain = "simulated"
)

// ClientDetails contains the simulated client and the contracts deployed to it
type ClientDetails struct {
	Client     Client
	RenAddress common.Address
	DNRAddress common.Address
	Chain      Chain
}

// FromURI will connect to a provided RPC uri
func FromURI(uri string, chain Chain) (ClientDetails, error) {
	if uri == "" && chain == ChainMainnet {
		uri = "https://mainnet.infura.io/"
	} else {
		uri = "https://ropsten.infura.io/"
	}

	client, err := ethclient.Dial(uri)
	if err != nil {
		return ClientDetails{}, err
	}

	if chain == ChainMainnet {
		panic("unimplemented")
	} else {
		return ClientDetails{
			Client:     client,
			RenAddress: common.HexToAddress("0x65d54eda5f032f2275caa557e50c029cfbccbb54"),
			DNRAddress: common.HexToAddress("0x9c06bb4e18e1aa352f99968b2984069c59ea2969"),
			Chain:      chain,
		}, nil
	}

}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// If the client is a simulated backend, it will commit the pending transactions to a block
func (b *ClientDetails) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {

	// sim, ok := b.(*backends.SimulatedBackend)
	// if ok {
	// 	sim.Commit()
	// 	sim.AdjustTime(10 * time.Second)
	// }

	switch b.Chain {
	case ChainGanache:
		return nil, nil
	case ChainSimulated:
		sim, ok := b.Client.(*backends.SimulatedBackend)
		if ok {
			sim.Commit()
			sim.AdjustTime(10 * time.Second)
			return nil, nil
		}
		fallthrough
	default:
		return bind.WaitMined(ctx, b.Client, tx)
	}

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
func (b *ClientDetails) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {

	switch b.Chain {
	case ChainGanache:
		return common.Address{}, nil
	case ChainSimulated:
		sim, ok := b.Client.(*backends.SimulatedBackend)
		if ok {
			sim.Commit()
			sim.AdjustTime(10 * time.Second)
			return common.Address{}, nil
		}
		fallthrough
	default:
		return bind.WaitDeployed(ctx, b.Client, tx)
	}

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
