package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
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

// Network is used to represent an Ethereum chain
type Network string

const (
	// NetworkMainnet represents the Ethereum mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkRopsten represents the Ethereum Ropsten testnet
	NetworkRopsten Network = "ropsten"
	// NetworkGanache represents a Ganache testrpc server
	NetworkGanache Network = "ganache"
	// NetworkSimulated represents a go-ethereum simulated backend
	NetworkSimulated Network = "simulated"
)

// Connection contains the simulated client and the contracts deployed to it
type Connection struct {
	Client     Client
	RenAddress common.Address
	DNRAddress common.Address
	HDEAddress common.Address
	Network    Network
}

// Connect to a URI.
func Connect(uri string, network Network, republicTokenAddress, darkNodeRegistryAddress, hyperdriveRegistry string) (Connection, error) {
	if uri == "" {
		switch network {
		case NetworkGanache:
			uri = "http://localhost:8545"
		case NetworkRopsten:
			uri = "https://ropsten.infura.io"
		default:
			return Connection{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}
	if republicTokenAddress == "" {
		switch network {
		case NetworkGanache:
			republicTokenAddress = RepublicTokenAddressOnGanache.String()
		case NetworkRopsten:
			republicTokenAddress = RepublicTokenAddressOnRopsten.String()
		default:
			return Connection{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}
	if darkNodeRegistryAddress == "" {
		switch network {
		case NetworkGanache:
			darkNodeRegistryAddress = DarkNodeRegistryAddressOnGanache.String()
		case NetworkRopsten:
			darkNodeRegistryAddress = DarkNodeRegistryAddressOnRopsten.String()
		default:
			return Connection{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}
	if hyperdriveRegistry == "" {
		switch network {
		case NetworkGanache:
			// fixme
			darkNodeRegistryAddress = DarkNodeRegistryAddressOnGanache.String()
		case NetworkRopsten:
			// fixme
			darkNodeRegistryAddress = DarkNodeRegistryAddressOnRopsten.String()
		default:
			return Connection{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}

	ethclient, err := ethclient.Dial(uri)
	if err != nil {
		return Connection{}, err
	}

	return Connection{
		Client:     ethclient,
		RenAddress: common.HexToAddress(republicTokenAddress),
		DNRAddress: common.HexToAddress(darkNodeRegistryAddress),
		Network:    network,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// If the client is a simulated backend, it will commit the pending transactions to a block
func (b *Connection) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {

	// sim, ok := b.(*backends.SimulatedBackend)
	// if ok {
	// 	sim.Commit()
	// 	sim.AdjustTime(10 * time.Second)
	// }

	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	case NetworkSimulated:
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
func (b *Connection) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {

	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	case NetworkSimulated:
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
