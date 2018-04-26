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
	// NetworkGanache represents a Ganache testrpc server
	NetworkGanache Network = "ganache"
)

// Contract addresses on ganache
var (
	RepublicTokenAddressOnGanache    = common.HexToAddress("0x8DE2a0D285cd6fDB47ABAe34024a6EED79ef0E92")
	DarknodeRegistryAddressOnGanache = common.HexToAddress("0xbF195E17802736Ff4E19275b961bb1c2D45f2c8D")
)

// Contract addresses on Ropsten
var (
	RepublicTokenAddressOnRopsten    = common.HexToAddress("0x65d54eda5f032f2275caa557e50c029cfbccbb54")
	DarknodeRegistryAddressOnRopsten = common.HexToAddress("0x69eb8d26157b9e12f959ea9f189A5D75991b59e3")
)

// Conn contains the client and the contracts deployed to it
type Conn struct {
	RawClient               *ethrpc.Client
	Client                  *ethclient.Client
	Network                 Network
	RepublicTokenAddress    common.Address
	DarknodeRegistryAddress common.Address
	TraderRegistryAddress   common.Address
	HyperdriveAddress       common.Address
}

// Connect to a URI.
func Connect(uri string, network Network, republicTokenAddress, darknodeRegistryAddr string) (Conn, error) {
	if uri == "" {
		switch network {
		case NetworkGanache:
			uri = "http://localhost:8545"
		case NetworkRopsten:
			uri = "https://ropsten.infura.io"
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}
	if republicTokenAddress == "" {
		switch network {
		case NetworkGanache:
			republicTokenAddress = RepublicTokenAddressOnGanache.String()
		case NetworkRopsten:
			republicTokenAddress = RepublicTokenAddressOnRopsten.String()
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}
	if darknodeRegistryAddr == "" {
		switch network {
		case NetworkGanache:
			darknodeRegistryAddr = DarknodeRegistryAddressOnGanache.String()
		case NetworkRopsten:
			darknodeRegistryAddr = DarknodeRegistryAddressOnRopsten.String()
		default:
			return Conn{}, fmt.Errorf("cannot connect to %s: unsupported", network)
		}
	}

	ethclient, err := ethclient.Dial(uri)
	if err != nil {
		return Conn{}, err
	}

	return Conn{
		Client:                  ethclient,
		RepublicTokenAddress:    common.HexToAddress(republicTokenAddress),
		DarknodeRegistryAddress: common.HexToAddress(darknodeRegistryAddr),
		Network:                 network,
	}, nil
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
//
// TODO: THIS DOES NOT WORK WITH PARITY, WHICH SENDS A TRANSACTION RECEIPT UPON
// RECEIVING A TX, NOT AFTER IT'S MINED
func (conn *Conn) PatchedWaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	switch conn.Network {
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
	switch conn.Network {
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
