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

// Network is used to represent a Republic Protocol network.
type Network string

const (
	// NetworkMainnet represents the mainnet
	NetworkMainnet Network = "mainnet"
	// NetworkTestnet represents the internal F∅ testnet
	NetworkTestnet Network = "testnet"
	// NetworkFalcon represents the internal Falcon testnet
	NetworkFalcon Network = "falcon"
	// NetworkNightly represents the internal Nightly testnet
	NetworkNightly Network = "nightly"
	// NetworkLocal represents a local network
	NetworkLocal Network = "local"
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
		case NetworkTestnet:
			config.URI = "https://kovan.infura.io"
		case NetworkFalcon:
			config.URI = "https://kovan.infura.io"
		case NetworkNightly:
			config.URI = "https://kovan.infura.io"
		case NetworkLocal:
			config.URI = "http://localhost:8545"
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", config.Network)
		}
	}

	if config.RepublicTokenAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.RepublicTokenAddress = "0x6f429121a3bd3e6c1c17edbc676eec44cf117faf"
		case NetworkFalcon:
			config.RepublicTokenAddress = "0x5e8148ab05ae724af7e6c2cbacdc65cca53ab3aa"
		case NetworkNightly:
			config.RepublicTokenAddress = "0x5e8148ab05ae724af7e6c2cbacdc65cca53ab3aa"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.DarknodeRegistryAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.DarknodeRegistryAddress = "0xf7b4360983A8fdd3E2ffb2e46cBeC65fA1b3075E"
		case NetworkFalcon:
			config.DarknodeRegistryAddress = "0x3aa3a8c5b2a4a2b0ee631650d88e9dc24f4c9254"
		case NetworkNightly:
			config.DarknodeRegistryAddress = "0x57c9b42a57b4a9f63f3bbc0685c89fc4070aa1d1"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.OrderbookAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.OrderbookAddress = "0x4782a0B10ad2EFEa1b488F53fDE2C25ceEd4a013"
		case NetworkFalcon:
			config.OrderbookAddress = "0x3DC8f53e3311750b4003BC535bea9a0bDAc172De"
		case NetworkNightly:
			config.OrderbookAddress = "0xa5f3DFaBbdDC987014dfE036Aa7e728259CeC5D9"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.RewardVaultAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.RewardVaultAddress = "0x2b82a678faB97CB3eF311680D39a2445fFd579D5"
		case NetworkFalcon:
			config.RewardVaultAddress = "0xdd0b6eae1bee54bac330886c5a8c93a661d5a43a"
		case NetworkNightly:
			config.RewardVaultAddress = "0x79b7bf5e274f804051a6916fcdf38f49ebf3526a"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.RenExBalancesAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.RenExBalancesAddress = "0xc5b98949AB0dfa0A7d4c07Bb29B002D6d6DA3e25"
		case NetworkFalcon:
			config.RenExBalancesAddress = "0x3E430B39D91F892FDFAc7D562D637183D53b5130"
		case NetworkNightly:
			config.RenExBalancesAddress = "0xFCDa8d4758CF5297fF12ceE53Eab085AfE32b799"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.RenExSettlementAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.RenExSettlementAddress = "0xd88C4f5162850B93c04EdEd90f7c552792c0B460"
		case NetworkFalcon:
			config.RenExSettlementAddress = "0x20b3cd8d1b9c7854f0efab0e774b9517e149a63b"
		case NetworkNightly:
			config.RenExSettlementAddress = "0xf204c48fe4dc196c864beb8aa8ae2b501e19653e"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
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
	case NetworkLocal:
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
	case NetworkLocal:
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
