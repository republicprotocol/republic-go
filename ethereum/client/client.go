package client

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client combines the interfaces for bind.ContractBackend and bind.DeployBackend
// type Client interface {
// 	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
// 	CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error)
// 	CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
// 	PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error)
// 	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
// 	SuggestGasPrice(ctx context.Context) (*big.Int, error)
// 	EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error)
// 	SendTransaction(ctx context.Context, tx *types.Transaction) error
// 	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
// 	SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
// 	BalanceAt(ctx context.Context, contract common.Address, blockNumber *big.Int) (*big.Int, error)
// }

// Network is used to represent an Ethereum chain
type Network string

const (
	// NetworkMainnet represents the Ethereum mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkRopsten represents the Ethereum Ropsten testnet
	NetworkRopsten Network = "ropsten"
	// NetworkGanache represents a Ganache testrpc server
	NetworkGanache Network = "ganache"
)

// Connection contains the client and the contracts deployed to it
type Connection struct {
	Client     *ethclient.Client
	RenAddress common.Address
	DNRAddress common.Address
	HDAddress  common.Address
	Network    Network
}

// Connect to a URI.
func Connect(uri string, network Network, republicTokenAddress, darkNodeRegistryAddress, hyperdriveAddress string) (Connection, error) {
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
	if hyperdriveAddress == "" {
		switch network {
		case NetworkGanache:
			hyperdriveAddress = HyperDriveAddressOnGanache.String()
		case NetworkRopsten:
			// fixme
			hyperdriveAddress = HyperDriveAddressOnGanache.String()
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
		HDAddress:  common.HexToAddress(hyperdriveAddress),
		Network:    network,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Connection) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	default:
		return bind.WaitMined(ctx, b.Client, tx)
	}
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (b *Connection) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	switch b.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	default:
		return bind.WaitDeployed(ctx, b.Client, tx)
	}
}

// TransferEth is a helper function for sending ETH to an address
func (b *Connection) TransferEth(ctx context.Context, from *bind.TransactOpts, to common.Address, value *big.Int) error {
	transactor := &bind.TransactOpts{
		From:     from.From,
		Nonce:    from.Nonce,
		Signer:   from.Signer,
		Value:    value,
		GasPrice: from.GasPrice,
		GasLimit: 30000,
		Context:  from.Context,
	}

	// Why is there no ethclient.Transfer?
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, b.Client, nil)
	tx, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	_, err = b.PatchedWaitMined(ctx, tx)
	return err
}
