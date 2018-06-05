package ethereum

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

// Contract addresses on ganache
var (
	RepublicTokenAddressOnGanache    = common.HexToAddress("0x8DE2a0D285cd6fDB47ABAe34024a6EED79ef0E92")
	DarknodeRegistryAddressOnGanache = common.HexToAddress("0xbF195E17802736Ff4E19275b961bb1c2D45f2c8D")
	RenLedgerAddressOnGanache        = common.HexToAddress("0x01cbe20EA5A49649F5615A59FaA30E88584634a2")
	RenExAccountsAddressOnGanache    = common.HexToAddress("0x9D174894dEa6470d25C0F6D847B94801EaE17Bc3")
)

// Contract addresses on Ropsten
var (
	RepublicTokenAddressOnRopsten    = common.HexToAddress("0x65d54eda5f032f2275caa557e50c029cfbccbb54")
	DarknodeRegistryAddressOnRopsten = common.HexToAddress("0x69eb8d26157b9e12f959ea9f189A5D75991b59e3")
	RenLedgerAddressOnRopsten        = common.HexToAddress("0x6235E09103bC7f205837237e4eAD855bC196E4D3")
	RenExAccountsAddressOnRopsten    = common.HexToAddress("0x0111111111111111111111111111111111111111") //fixme
)

// Contract addresses on Kovan
var (
	RepublicTokenAddressOnKovan    = common.HexToAddress("0x596F8c39aEc9fb72D0F591DEe4408516f4C9DdA4")
	DarknodeRegistryAddressOnKovan = common.HexToAddress("0x7b0e5c2945020996408fdf73ef1846d7c0dcac78")
	RenLedgerAddressOnKovan        = common.HexToAddress("0x9ac38a5f17aae6d473b0f87bd6e42e8958043c70")
	RenExAccountsAddressOnKovan    = common.HexToAddress("0x42b481e2dafb19f9df6f0cdc6bd6e45034e38529")
)

// Config defines the different settings for connecting the Darknode
// to an Ethereum network, and the Republic Protocol smart contracts deployed
// on Ethereum.
type Config struct {
	Network                 Network `json:"network"` // One of "ganache", "ropsten", or "mainnet" ("mainnet" is not current supported)
	URI                     string  `json:"uri"`
	RepublicTokenAddress    string  `json:"republicTokenAddress"`
	DarknodeRegistryAddress string  `json:"darknodeRegistryAddress"`
	RenLedgerAddress        string  `json:"renLedgerAddress"`
	RenExAccountsAddress    string  `json:"renExAccountsAddress"`
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
	if config.RepublicTokenAddress == "" {
		switch config.Network {
		case NetworkGanache:
			config.RepublicTokenAddress = RepublicTokenAddressOnGanache.String()
		case NetworkKovan:
			config.RepublicTokenAddress = RepublicTokenAddressOnKovan.String()
		case NetworkRopsten:
			config.RepublicTokenAddress = RepublicTokenAddressOnRopsten.String()
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", config.Network)
		}
	}
	if config.DarknodeRegistryAddress == "" {
		switch config.Network {
		case NetworkGanache:
			config.DarknodeRegistryAddress = DarknodeRegistryAddressOnGanache.String()
		case NetworkKovan:
			config.DarknodeRegistryAddress = DarknodeRegistryAddressOnKovan.String()
		case NetworkRopsten:
			config.DarknodeRegistryAddress = DarknodeRegistryAddressOnRopsten.String()
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", config.Network)
		}
	}
	if config.RenLedgerAddress == "" {
		switch config.Network {
		case NetworkGanache:
			config.RenLedgerAddress = RenLedgerAddressOnGanache.String()
		case NetworkKovan:
			config.RenLedgerAddress = RenLedgerAddressOnKovan.String()
		case NetworkRopsten:
			config.RenLedgerAddress = RenLedgerAddressOnRopsten.String()
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", config.Network)
		}
	}

	if config.RenExAccountsAddress == "" {
		switch config.Network {
		case NetworkGanache:
			config.RenExAccountsAddress = RenExAccountsAddressOnGanache.String()
		case NetworkKovan:
			config.RenExAccountsAddress = RenExAccountsAddressOnKovan.String()
		case NetworkRopsten:
			panic("not deployed yet")
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
