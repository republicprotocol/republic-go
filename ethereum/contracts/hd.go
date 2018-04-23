package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/hyperdrive"
)

// HyperdriveContract defines the golang interface for interacting with the
// hyperdrive smart contract on Ethereum.
type HyperdriveContract struct {
	network                   client.Network
	context                   context.Context
	client                    client.Connection
	callOpts                  *bind.CallOpts
	transactOpts              *bind.TransactOpts
	binding                   *bindings.Hyperdrive
	hyperdriveRegistryAddress common.Address
}

// NewHyperdriveContract creates a new HyperdriveContract with given parameters
func NewHyperdriveContract(ctx context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (HyperdriveContract, error) {
	contract, err := bindings.NewHyperdrive(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return HyperdriveContract{}, err
	}

	return HyperdriveContract{
		network:      clientDetails.Network,
		context:      ctx,
		callOpts:     callOpts,
		transactOpts: transactOpts,
		binding:      contract,
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
