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
			config.RepublicTokenAddress = "0x99806d107eda625516d954621df175a002d223e6"
		case NetworkFalcon:
			config.RepublicTokenAddress = "0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"
		case NetworkNightly:
			config.RepublicTokenAddress = "0x15f692d6b9ba8cec643c7d16909e8acdec431bf6"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.DarknodeRegistryAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.DarknodeRegistryAddress = "0xd1c3b5f2fe4eec6c262a5e1b161e5e099fd8325e"
		case NetworkFalcon:
			config.DarknodeRegistryAddress = "0xdaa8c30af85070506f641e456afdb84d4ba972bd"
		case NetworkNightly:
			config.DarknodeRegistryAddress = "0x8a31d477267a5af1bc5142904ef0afa31d326e03"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.DarknodeRewardVaultAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.DarknodeRewardVaultAddress = "0x870a7d5d1eb513fc422ebfcc76b598f860c7c2cf"
		case NetworkFalcon:
			config.DarknodeRewardVaultAddress = "0x401e7d7ce6f51ea1a8d4f582413e2fabda68daa8"
		case NetworkNightly:
			config.DarknodeRewardVaultAddress = "0xda43560f5fe6c6b5e062c06fee0f6fbc71bbf18a"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.DarknodeSlasherAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.DarknodeSlasherAddress = "0x6c52b2fd5b6c3e6baf47e05af880fc95b9c8079c"
		case NetworkFalcon:
			config.DarknodeSlasherAddress = "0x71ec5f4558e87d6afb5c5ff0b4bdd058d62ed3d1"
		case NetworkNightly:
			config.DarknodeSlasherAddress = "0x38458ef4a185455cba57a7594b0143c53ad057c1"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.OrderbookAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.OrderbookAddress = "0x9a016649d97d44a055c26cbcadbc45a1ac563c89"
		case NetworkFalcon:
			config.OrderbookAddress = "0x592d16f8c5fa8f1e074ab3c2cd1acd087adcdc0b"
		case NetworkNightly:
			config.OrderbookAddress = "0x376127adc18260fc238ebfb6626b2f4b59ec9b66"
		case NetworkLocal:
		default:
			return Conn{}, fmt.Errorf("no default contract address on %s", config.Network)
		}
	}
	if config.SettlementRegistryAddress == "" {
		switch config.Network {
		case NetworkTestnet:
			config.SettlementRegistryAddress = "0xc07780d6e1f24434b1766068f0e44b10a5ff5755"
		case NetworkFalcon:
			config.SettlementRegistryAddress = "0x6246ff83ddef23d9509ba80aa3ee650ab0321f0b"
		case NetworkNightly:
			config.SettlementRegistryAddress = "0x399a70ed71897836468fd74ea19138df90a78d79"
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

// TokenAddresses returns the tokens for the provided network
func TokenAddresses(network Network) map[string]string {
	tokens := map[string]string{}
	switch network {
	case NetworkTestnet:
		tokens["TUSD"] = "0x289f785d9137ecf38a46a678cf4e9e98d32a06d4"
		tokens["DGX"] = "0x0798297a11cefef7479e40e67839fee3c025691e"
		tokens["REN"] = "0x6f429121a3bd3e6c1c17edbc676eec44cf117faf"
		tokens["ZRX"] = "0x099ea44e49e34250e247a150c66c89b314216e34"
		tokens["OMG"] = "0x0f48986df7b79fbb085753dc2fefe10dde7dd232"
	case NetworkFalcon:
		tokens["TUSD"] = "0x1c428ab82c06dbe9af414e6c923862d4c3ae0579"
		tokens["DGX"] = "0xf4faf1b22cee0a024ad6b12bb29ec0e13f5827c2"
		tokens["REN"] = "0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"
		tokens["ZRX"] = "0x295a3894fc98b021735a760dbc7aed265663ca42"
		tokens["OMG"] = "0x21c1ba3ea123eb23815c689ee05a944119c7f428"
	case NetworkNightly:
		tokens["TUSD"] = "0xa86c6a3322efa371faad6a9b04708788e3592615"
		tokens["DGX"] = "0x092ece29781777604afac04887af30042c3bc5df"
		tokens["REN"] = "0x15f692d6b9ba8cec643c7d16909e8acdec431bf6"
		tokens["ZRX"] = "0xeb5a7335e850176b44ca1990730d1a2433e195f3"
		tokens["OMG"] = "0x69440b57b52e323cbd12a162a5f9870f61182918"
	default:
		panic("unknown network")
	}
	return tokens
}
