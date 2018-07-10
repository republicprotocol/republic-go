package transact

import (
	"context"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/republicprotocol/republic-go/contract"
)

type TxSender interface {
	Send(f func() (*types.Transaction, error)) (*types.Transaction, error)
}

type txSender struct {
	conn           contract.Conn
	transactOptsMu *sync.Mutex
	transactOpts   bind.TransactOpts
}

func NewTxSender(conn contract.Conn, transactOpts bind.TransactOpts) (TxSender, error) {
	nonce, err := conn.Client.PendingNonceAt(context.Background(), transactOpts.From)
	if err != nil {
		return nil, err
	}
	transactOpts.Nonce = big.NewInt(int64(nonce))
	return &txSender{
		conn:           conn,
		transactOptsMu: new(sync.Mutex),
		transactOpts:   transactOpts,
	}, nil
}

// Send locks TxSender resources to execute function f (handling nonces explicitly)
// and will wait until the block has been mined on the blockchain. This will allow
// parallel requests to the blockchain since the sender will be unlocked before
// waiting for transaction to complete execution on the blockchain.
func (txSender *txSender) Send(f func() (*types.Transaction, error)) (*types.Transaction, error) {
	txSender.transactOptsMu.Lock()
	defer txSender.transactOptsMu.Unlock()
	return txSender.send(f)
}

func (txSender *txSender) send(f func() (*types.Transaction, error)) (*types.Transaction, error) {
	tx, err := f()
	if err == nil {
		txSender.transactOpts.Nonce.Add(txSender.transactOpts.Nonce, big.NewInt(1))
		return tx, nil
	}
	if err == core.ErrNonceTooLow || err == core.ErrReplaceUnderpriced || strings.Contains(err.Error(), "nonce is too low") {
		log.Info("[tx error] nonce too low = %v", err)
		txSender.transactOpts.Nonce.Add(txSender.transactOpts.Nonce, big.NewInt(1))
		return txSender.send(f)
	}
	if err == core.ErrNonceTooHigh {
		log.Info("[tx error] nonce too high = %v", err)
		txSender.transactOpts.Nonce.Sub(txSender.transactOpts.Nonce, big.NewInt(1))
		return txSender.send(f)
	}

	// If any other type of nonce error occurs we will refresh the nonce and
	// try again for up to 1 minute
	var nonce uint64
	for try := 0; try < 60 && strings.Contains(err.Error(), "nonce"); try++ {
		log.Errorf("[tx error] unknown = %v", err)
		time.Sleep(time.Second)
		nonce, err = txSender.conn.Client.PendingNonceAt(context.Background(), txSender.transactOpts.From)
		if err != nil {
			continue
		}
		txSender.transactOpts.Nonce = big.NewInt(int64(nonce))
		if tx, err = f(); err == nil {
			txSender.transactOpts.Nonce.Add(txSender.transactOpts.Nonce, big.NewInt(1))
			return tx, nil
		}
	}

	return tx, err
}
