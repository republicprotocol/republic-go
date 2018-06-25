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
			config.RepublicTokenAddress = "0xf1da6f4a594553335edea6b1203a4b590c752e32"
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
			config.DarknodeRegistryAddress = "0x5e3c8b0f7229f1f1873267b6811465fef73d53ca"
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
			config.OrderbookAddress = "0x5E976b687D902f13F6af33a5CB097440DDCB149e"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.RewardVaultAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.RewardVaultAddress = "0x98e32f470978915aE0C2a11B2f696A125AB7fB5b"
		case NetworkFalcon:
			config.RewardVaultAddress = "0xdd0b6eae1bee54bac330886c5a8c93a661d5a43a"
		case NetworkNightly:
			config.RewardVaultAddress = "0x880ad65dc5b3f33123382416351eef98b4aad7f1"
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
			config.RenExBalancesAddress = "0xD5B5b26521665Cb37623DCA0E49c553b41dbF076"
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
			config.RenExSettlementAddress = "0xef685d1d44ea983927d9f8d67f77894faec92fcf"
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
