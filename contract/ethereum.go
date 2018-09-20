package contract

import (
	"context"
	"errors"
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
		case NetworkMainnet:
			config.URI = "https://mainnet.infura.io"
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

	// Get default contract addresses if they haven't been set.
	if config.RepublicTokenAddress == "" {
		config.RepublicTokenAddress = RepublicTokenAddress(config.Network)
	}
	if config.DarknodeRegistryAddress == "" {
		config.DarknodeRegistryAddress = DarknodeRegistryAddress(config.Network)
	}
	if config.DarknodeRewardVaultAddress == "" {
		config.DarknodeRewardVaultAddress = DarknodeRewardVaultAddress(config.Network)
	}
	if config.DarknodeSlasherAddress == "" {
		config.DarknodeSlasherAddress = DarknodeSlasherAddress(config.Network)
	}
	if config.OrderbookAddress == "" {
		config.OrderbookAddress = OrderbookAddress(config.Network)
	}
	if config.SettlementRegistryAddress == "" {
		config.SettlementRegistryAddress = SettlementRegistryAddress(config.Network)
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
		receipt, err := bind.WaitMined(ctx, conn.Client, tx)
		if err != nil {
			return nil, err
		}
		if receipt.Status != types.ReceiptStatusSuccessful {
			return receipt, errors.New("transaction reverted")
		}
		return receipt, nil
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
		GasLimit: 21000,
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

// SendEth is a helper function for sending ETH to an address
func (conn *Conn) SendEth(ctx context.Context, from *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	transactor := &bind.TransactOpts{
		From:     from.From,
		Nonce:    from.Nonce,
		Signer:   from.Signer,
		Value:    value,
		GasPrice: from.GasPrice,
		GasLimit: 21000,
		Context:  from.Context,
	}

	// Why is there no ethclient.Transfer?
	bound := bind.NewBoundContract(to, abi.ABI{}, nil, conn.Client, nil)
	return bound.Transfer(transactor)
}
