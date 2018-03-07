package dnr

import (
	"context"
	"log"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
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

// FromURI will connect to a provided RPC uri
func FromURI(uri string) (Client, error) {
	return ethclient.Dial(uri)
}

// Simulated will create a simulated client
func Simulated(auth1 *bind.TransactOpts, auth2 *bind.TransactOpts) Client {
	sim := backends.NewSimulatedBackend(core.GenesisAlloc{auth1.From: {Balance: big.NewInt(9000000000000000000)}, auth2.From: {Balance: big.NewInt(9000000000000000000)}})
	return sim
}

// DeployERC20 deploys and ERC20 contract
func DeployERC20(context context.Context, connection Client, auth *bind.TransactOpts) (*types.Transaction, common.Address) {
	// Deploy a token contract on the simulated blockchain
	address, tx, _, err := contracts.DeployAtomicSwapERC20(auth, connection)
	if err != nil {
		log.Fatalf("Failed to deploy: %v", err)
	}
	PatchedWaitDeployed(context, connection, tx)
	return tx, address
}

// PatchedWaitMined waits for tx to be mined on the blockchain.
// It stops waiting when the context is canceled.
// (Go-ethereum's WaitMined is not compatible with Parity's getTransactionReceipt)
// NOTE: If something goes wrong, this will hang!
func PatchedWaitMined(ctx context.Context, b Client, tx *types.Transaction) (*types.Receipt, error) {
	return bind.WaitMined(ctx, b, tx)

	// sim, ok := b.(*backends.SimulatedBackend)
	// if ok {
	// 	sim.Commit()
	// 	sim.AdjustTime(10 * time.Second)
	// }

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
// (Go-ethereum's WaitMined is not compatible with Parity's getTransactionReceipt)
// NOTE: If something goes wrong, this will hang!
func PatchedWaitDeployed(ctx context.Context, b Client, tx *types.Transaction) (common.Address, error) {
	return bind.WaitDeployed(ctx, b, tx)

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
