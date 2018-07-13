package transact

import (
	"context"
	"math/big"
	"net/rpc"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Conn is an interface to the Ethereum blockchain.
type Conn interface {

	// WaitMined will wait until a transaction has been mined, or until the
	// context.Context is canceled.
	WaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error)

	// WaitDeployed will wait until a contract deployment transaction has been
	// mined, or until the context.Context is canceled.
	WaitDeployed(ctx context.Context, tx *types.Transaction) (*types.Receipt, error)

	// PendingNonce returns the pending nonce of an address.
	PendingNonce(ctx context.Context, addr common.Address) (*big.Int, error)
}

type conn struct {
	rpcClient *rpc.Client
	ethClient *ethclient.Client
}

// WaitMined implements the Conn interface. This does not work with Parity,
// because Parity sends a transaction receipt upon receiving the transaction
// instead of after the transaction is mined.
func (conn *conn) WaitMined(ctx context.Context, tx *types.Transaction) (*types.Receipt, error) {
	time.Sleep(100 * time.Millisecond)
	return bind.WaitMined(ctx, conn.ethClient, tx)
}

// WaitDeployed implements the Conn interface. This does not work with Parity,
// because Parity sends a transaction receipt upon receiving the transaction
// instead of after the transaction is mined.
func (conn *conn) WaitDeployed(ctx context.Context, tx *types.Transaction) (common.Address, error) {
	time.Sleep(100 * time.Millisecond)
	return bind.WaitDeployed(ctx, conn.ethClient, tx)
}

// PendingNonce implements the Conn interface.
func (conn *conn) PendingNonce(ctx context.Context, addr common.Address) (*big.Int, error) {
	nonce, err := conn.ethClient.PendingNonceAt(ctx, addr)
	return big.NewInt(0).SetUint64(nonce), err
}
