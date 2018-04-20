package contracts

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/ethereum/bindings"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/stackint"
)

// Hyperdrive
type HyperdriveRegistry struct {
	network                   client.Network
	context                   context.Context
	client                    client.Connection
	callOpts                  *bind.CallOpts
	transactOpts              *bind.TransactOpts
	binding                   *bindings.HyperdriveEpoch
	hyperdriveRegistryAddress common.Address
}

func NewHyperdriveRegistry(ctx context.Context, clientDetails client.Connection, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (HyperdriveRegistry, error) {
	contract, err := bindings.NewHyperdriveEpoch(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return HyperdriveRegistry{}, err
	}

	return HyperdriveRegistry{
		network:      clientDetails.Network,
		context:      ctx,
		callOpts:     callOpts,
		transactOpts: transactOpts,
		binding:      contract,
	}, nil
}

func (hyper HyperdriveRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := hyper.binding.CurrentEpoch(hyper.callOpts)
	if err != nil {
		return Epoch{}, err
	}
	timestamp, err := stackint.FromBigInt(epoch.Timestamp)
	if err != nil {
		return Epoch{}, err
	}

	var blockhash [32]byte
	for i, b := range epoch.Blockhash.Bytes() {
		blockhash[i] = b
	}

	return Epoch{
		Blockhash: blockhash,
		Timestamp: timestamp,
	}, nil
}

func (hyper HyperdriveRegistry) Epoch() (*types.Transaction, error) {
	tx, err := hyper.binding.Epoch(hyper.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = hyper.client.PatchedWaitMined(hyper.context, tx)
	return tx, err
}
