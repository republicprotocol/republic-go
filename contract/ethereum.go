package contract

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
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// Network is used to represent an Ethereum chain
type Network string

const (
	// NetworkMainnet represents the Ethereum mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkRopsten represents the Ethereum Ropsten testnet
	NetworkRopsten Network = "ropsten"
	// NetworkKovan represents the Ethereum Kovan testnet
	NetworkKovan Network = "kovan"
	// NetworkGanache represents a Ganache testrpc server
	NetworkGanache Network = "ganache"
)

// Config defines the different settings for connecting the Darknode
// to an Ethereum network, and the Republic Protocol smart contracts deployed
// on Ethereum.
type Config struct {
	Network                 Network `json:"network"` // One of "ganache", "ropsten", or "mainnet" ("mainnet" is not current supported)
	URI                     string  `json:"uri"`
	RepublicTokenAddress    string  `json:"republicTokenAddress"`
	DarknodeRegistryAddress string  `json:"darknodeRegistryAddress"`
	OrderbookAddress        string  `json:"orderbookAddress"`
	RewardVaultAddress      string  `json:"rewardVaultAddress"`
	RenExBalancesAddress    string  `json:"renExBalancesAddress"`
	RenExSettlementAddress  string  `json:"renExSettlementAddress"`
}

// Conn contains the client and the contracts deployed to it
type Conn struct {
	RawClient *ethrpc.Client
	Client    *ethclient.Client
	Config    Config
}

// Connect to a URI.
func Connect(config Config) (Conn, error) {
	if config.URI == "" {
		switch config.Network {
		case NetworkGanache:
			config.URI = "http://localhost:8545"
		case NetworkKovan:
			config.URI = "https://kovan.infura.io"
		case NetworkRopsten:
			config.URI = "https://ropsten.infura.io"
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", config.Network)
		}
	}

	ethclient, err := ethclient.Dial(config.URI)
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		Client: ethclient,
		Config: config,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (conn *Conn) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch conn.Config.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return nil, nil
	default:
		return bind.WaitMined(ctx, conn.Client, tx)
	}
}

// PatchedWaitDeployed waits for a contract deployment transaction and returns the on-chain
// contract address when it is mined. It stops waiting when ctx is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (conn *Conn) PatchedWaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	switch conn.Config.Network {
	case NetworkGanache:
		time.Sleep(100 * time.Millisecond)
		return common.Address{}, nil
	default:
		return bind.WaitDeployed(ctx, conn.Client, tx)
	}
}

// TransferEth is a helper function for sending ETH to an address
func (conn *Conn) TransferEth(ctx context.Context, from *bind.TransactOpts, to common.Address, value *big.Int) error {
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
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, conn.Client, nil)
	tx, err := bound.Transfer(transactor)
	if err != nil {
		return err
	}
	_, err = conn.PatchedWaitMined(ctx, tx)
	return err
}
