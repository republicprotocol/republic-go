package contracts

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/hyperdrive"
)

// HyperdriveContract defines the golang interface for interacting with the
// hyperdrive smart contract on Ethereum.
type HyperdriveContract struct {
	network           client.Network
	context           context.Context
	client            client.Connection
	callOpts          *bind.CallOpts
	transactOpts      *bind.TransactOpts
	binding           *bindings.Hyperdrive
	hyperdriveAddress common.Address
}

// NewHyperdriveContract creates a new HyperdriveContract with given parameters
func NewHyperdriveContract(ctx context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (HyperdriveContract, error) {
	contract, err := bindings.NewHyperdrive(clientDetails.HDAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return HyperdriveContract{}, err
	}

	return HyperdriveContract{
		network:           clientDetails.Network,
		context:           ctx,
		client:            clientDetails,
		callOpts:          callOpts,
		transactOpts:      transactOpts,
		binding:           contract,
		hyperdriveAddress: clientDetails.HDAddress,
	}, nil
}

// SendTx sends an tx to the contract, it returns an error if there is an
// conflict with previous orders. You need to register with the darkNodeRegistry
// to send the tx.
func (hyper HyperdriveContract) SendTx(tx hyperdrive.Tx) (*types.Transaction, error) {
	transaction, err := hyper.binding.SendTx(hyper.transactOpts, tx.Nonces)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (hyper HyperdriveContract) GetBlockNumberOfTx(tx *types.Transaction) (uint64, error) {
	receipt, err := hyper.client.Client.(*ethclient.Client).TransactionReceipt(hyper.context, tx.Hash())
	log.Println("err here  is ", err)
	if err != nil {
		return 0, err
	}

	log.Println("receipt is ", receipt)
	return 1, nil
}

func (hyper HyperdriveContract) CurrentBlock() (*types.Block, error) {
	return hyper.client.Client.(*ethclient.Client).BlockByNumber(hyper.context, nil)
}

// SetGasLimit sets the gas limit to use for transactions
func (hyper *HyperdriveContract) SetGasLimit(limit uint64) {
	hyper.transactOpts.GasLimit = limit
}

func (hyper *HyperdriveContract) BlockByNumber(blockNumber *big.Int) (*types.Block, error) {
	return hyper.client.Client.(*ethclient.Client).BlockByNumber(hyper.context, blockNumber)
}
