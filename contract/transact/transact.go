package transact

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
)

// Transacter exposes functionality for sending transactions to the Ethereum
// blockchain. It exposes the specialised Transacter.Transfer method for
// transferring Ether, and the generic Transacter.Transact method for sending
// all other transactions.
type Transacter interface {

	// Transfer Ether to an address. This method should not wait until the
	// transaction is mined.
	Transfer(ctx context.Context, to common.Address, value *big.Int) (*types.Transaction, error)

	// Transact builds and sends a transaction using the builder function
	// provided. It ensures that the correct nonce is used for the transaction
	// and retries if an error occurs. This method should not wait until the
	// transaction is mined.
	Transact(ctx context.Context, buildTx func(context.Context, *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error)
}

type transacter struct {
	client *ethclient.Client

	transactOptsMu *sync.Mutex
	transactOpts   bind.TransactOpts
}

// NewTransacter returns a new Transacter that uses the ethclient.Client for
// interacting with the Ethereum blockchain, and uses a deep copy of the
// bind.TransactOpts for managing transaction options. The Transacter is safe
// for concurrent use.
func NewTransacter(client *ethclient.Client, transactOpts bind.TransactOpts) (Transacter, error) {
	// Build a transact opts
	nonce, err := client.PendingNonceAt(context.Background(), transactOpts.From)
	if err != nil {
		return nil, err
	}
	transactOpts.From = transactOpts.From
	transactOpts.GasLimit = transactOpts.GasLimit
	transactOpts.GasPrice = big.NewInt(0).Set(transactOpts.GasPrice)
	transactOpts.Nonce = big.NewInt(0).SetUint64(nonce)
	transactOpts.Signer = transactOpts.Signer
	transactOpts.Value = big.NewInt(0).Set(transactOpts.Value)

	return &transacter{
		client: client,

		transactOptsMu: new(sync.Mutex),
		transactOpts:   transactOpts,
	}, nil
}

// Transfer implements the Transacter interface.
func (transacter *transacter) Transfer(ctx context.Context, to common.Address, value *big.Int) (*types.Transaction, error) {
	transacter.transactOptsMu.Lock()
	defer transacter.transactOptsMu.Unlock()

	oldValue := big.NewInt(0).Set(transacter.transactOpts.Value)
	defer transacter.transactOpts.Value.Set(oldValue)

	empty := bind.NewBoundContract(to, abi.ABI{}, nil, transacter.client, nil)
	return transacter.transact(ctx, func(ctx context.Context, transactOpts *bind.TransactOpts) (*types.Transaction, error) {
		return empty.Transfer(transactOpts)
	})
}

// Transact implements the Transacter interface.
func (transacter *transacter) Transact(ctx context.Context, buildTx func(context.Context, *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	transacter.transactOptsMu.Lock()
	defer transacter.transactOptsMu.Unlock()
	return transacter.transact(ctx, buildTx)
}

func (transacter *transacter) transact(ctx context.Context, buildTx func(context.Context, *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	// Check if the context.Context is done
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	tx, err := buildTx(ctx, &transacter.transactOpts)

	// If an expected error occurs then change the nonce appropriately
	if err == nil {
		transacter.transactOpts.Nonce.Add(transacter.transactOpts.Nonce, big.NewInt(1))
		return tx, nil
	}
	if err == core.ErrNonceTooLow || err == core.ErrReplaceUnderpriced || strings.Contains(err.Error(), "nonce is too low") {
		log.Info(fmt.Sprintf("[tx error] nonce too low = %v", err))
		transacter.transactOpts.Nonce.Add(transacter.transactOpts.Nonce, big.NewInt(1))
		return transacter.transact(ctx, buildTx)
	}
	if err == core.ErrNonceTooHigh {
		log.Info(fmt.Sprintf("[tx error] nonce too high = %v", err))
		transacter.transactOpts.Nonce.Sub(transacter.transactOpts.Nonce, big.NewInt(1))
		return transacter.transact(ctx, buildTx)
	}

	// If any other type of nonce error occurs we will refresh the nonce and
	// try again for up to 1 minute
	var nonce uint64
	for try := 0; try < 60 && strings.Contains(err.Error(), "nonce"); try++ {
		log.Error(fmt.Sprintf("[tx error] unknown = %v", err))

		// Delay for a second or until the contex is done
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Second):
		}

		nonce, err = transacter.client.PendingNonceAt(ctx, transacter.transactOpts.From)
		if err != nil {
			continue
		}
		transacter.transactOpts.Nonce = big.NewInt(int64(nonce))
		if tx, err = buildTx(ctx, &transacter.transactOpts); err == nil {
			transacter.transactOpts.Nonce.Add(transacter.transactOpts.Nonce, big.NewInt(1))
			return tx, nil
		}
	}
	return tx, err
}
