package dnr

import (
	"context"
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/stackint"
)

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash [32]byte
	Timestamp stackint.Int1024
}

// DarkNodeRegistry is the dark node interface
type DarkNodeRegistry struct {
	Chain                   connection.Chain
	context                 context.Context
	client                  *connection.ClientDetails
	transactOpts            *bind.TransactOpts
	callOpts                *bind.CallOpts
	binding                 *bindings.DarkNodeRegistry
	tokenBinding            *bindings.RepublicToken
	darkNodeRegistryAddress common.Address
}

// NewDarkNodeRegistry returns a Dark node registrar
func NewDarkNodeRegistry(context context.Context, clientDetails *connection.ClientDetails, transactOpts *bind.TransactOpts, callOpts *bind.CallOpts) (DarkNodeRegistry, error) {
	contract, err := bindings.NewDarkNodeRegistry(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	return DarkNodeRegistry{
		Chain:                   clientDetails.Chain,
		context:                 context,
		client:                  clientDetails,
		transactOpts:            transactOpts,
		callOpts:                callOpts,
		binding:                 contract,
		tokenBinding:            renContract,
		darkNodeRegistryAddress: clientDetails.DNRAddress,
	}, nil
}

// Register a new dark node
func (darkNodeRegistry *DarkNodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := darkNodeRegistry.binding.Register(darkNodeRegistry.transactOpts, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
}

// Deregister an existing dark node
func (darkNodeRegistry *DarkNodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Deregister(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// Refund withdraws the bond. Must be called before reregistering.
func (darkNodeRegistry *DarkNodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Refund(darkNodeRegistry.transactOpts, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// GetBond retrieves the bond of an existing dark node
func (darkNodeRegistry *DarkNodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return stackint.Int1024{}, err
	}
	bond, err := darkNodeRegistry.binding.GetBond(darkNodeRegistry.callOpts, darkNodeIDByte)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// IsRegistered returns true if the node is registered
func (darkNodeRegistry *DarkNodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsRegistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (darkNodeRegistry *DarkNodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsDeregistered(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approved Ren to it
func (darkNodeRegistry *DarkNodeRegistry) ApproveRen(value *stackint.Int1024) error {
	txn, err := darkNodeRegistry.tokenBinding.Approve(darkNodeRegistry.transactOpts, darkNodeRegistry.client.DNRAddress, value.ToBigInt())
	if err != nil {
		return err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return err
}

// CurrentEpoch returns the current epoch
func (darkNodeRegistry *DarkNodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.callOpts)
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

// Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
func (darkNodeRegistry *DarkNodeRegistry) Epoch() (*types.Transaction, error) {
	tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// WaitForEpoch guarantees that an Epoch as passed (and calls Epoch if connected to Ganache)
func (darkNodeRegistry *DarkNodeRegistry) WaitForEpoch() (*types.Transaction, error) {
	currentEpoch, err := darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return nil, err
	}
	nextEpoch := currentEpoch
	var tx *types.Transaction
	for currentEpoch.Blockhash == nextEpoch.Blockhash {
		if darkNodeRegistry.Chain == connection.ChainGanache {
			tx, err = darkNodeRegistry.binding.Epoch(darkNodeRegistry.transactOpts)
			if err != nil {
				return nil, err
			}
		}
		nextEpoch, err = darkNodeRegistry.CurrentEpoch()
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Millisecond * 10)
	}
	return tx, nil
}

// GetOwner gets the owner of the given dark node
func (darkNodeRegistry *DarkNodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return darkNodeRegistry.binding.GetOwner(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
func (darkNodeRegistry *DarkNodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return darkNodeRegistry.binding.GetPublicKey(darkNodeRegistry.callOpts, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
func (darkNodeRegistry *DarkNodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := darkNodeRegistry.binding.GetDarkNodes(darkNodeRegistry.callOpts)
	if err != nil {
		return nil, err
	}
	arr := make([][]byte, len(ret))
	for i := range ret {
		arr[i] = ret[i][:]
	}
	return arr, nil
}

// MinimumBond gets the minimum viable bond amount
func (darkNodeRegistry *DarkNodeRegistry) MinimumBond() (stackint.Int1024, error) {
	bond, err := darkNodeRegistry.binding.MinimumBond(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval gets the minimum epoch interval
func (darkNodeRegistry *DarkNodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumEpochInterval(darkNodeRegistry.callOpts)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (darkNodeRegistry *DarkNodeRegistry) SetGasLimit(limit uint64) {
	darkNodeRegistry.transactOpts.GasLimit = limit
}

// WaitUntilRegistration waits until the registration is successful
func (darkNodeRegistry *DarkNodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
	isRegistered := false
	for !isRegistered {
		var err error
		isRegistered, err = darkNodeRegistry.IsRegistered(darkNodeID)
		if err != nil {
			return err
		}
		darkNodeRegistry.WaitForEpoch()

	}
	return nil
}

func toByte(id []byte) ([20]byte, error) {
	twentyByte := [20]byte{}
	if len(id) != 20 {
		return twentyByte, errors.New("length mismatch")
	}
	for i := range id {
		twentyByte[i] = id[i]
	}
	return twentyByte, nil
}
