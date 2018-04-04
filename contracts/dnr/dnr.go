package dnr

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
	"github.com/republicprotocol/republic-go/stackint"
)

// DarkNodeRegistry is the interface defining the Dark Node Registrar
// type DarkNodeRegistry interface {
// 	Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error)
// 	Deregister(darkNodeID []byte) (*types.Transaction, error)
// 	Refund(darkNodeID []byte) (*types.Transaction, error)
// 	Epoch() (*types.Transaction, error)

// 	CurrentEpoch() (Epoch, error)
// 	GetBond(darkNodeID []byte) (stackint.Int1024, error)
// 	GetCommitment(darkNodeID []byte) ([32]byte, error)
// 	GetOwner(darkNodeID []byte) (common.Address, error)
// 	GetPublicKey(darkNodeID []byte) ([]byte, error)
// 	GetAllNodes() ([][]byte, error)

// 	MinimumBond() (stackint.Int1024, error)
// 	MinimumEpochInterval() (stackint.Int1024, error)

// 	IsRegistered(darkNodeID []byte) (bool, error)
// 	IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error)
// 	WaitUntilRegistration(darkNodeID []byte) error
// }

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash [32]byte
	Timestamp *stackint.Int1024
}

// DarkNodeRegistry is the dark node interface
type DarkNodeRegistry struct {
	context                 context.Context
	client                  *connection.ClientDetails
	auth1                   *bind.TransactOpts
	auth2                   *bind.CallOpts
	binding                 *bindings.DarkNodeRegistry
	tokenBinding            *bindings.RepublicToken
	darkNodeRegistryAddress common.Address
}

// NewDarkNodeRegistry returns a Dark node registrar
func NewDarkNodeRegistry(context context.Context, clientDetails *connection.ClientDetails, auth1 *bind.TransactOpts, auth2 *bind.CallOpts) (DarkNodeRegistry, error) {
	contract, err := bindings.NewDarkNodeRegistry(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	renContract, err := bindings.NewRepublicToken(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return DarkNodeRegistry{}, err
	}
	return DarkNodeRegistry{
		context:                 context,
		client:                  clientDetails,
		auth1:                   auth1,
		auth2:                   auth2,
		binding:                 contract,
		tokenBinding:            renContract,
		darkNodeRegistryAddress: clientDetails.DNRAddress,
	}, nil
}

// Register a new dark node
func (darkNodeRegistry *DarkNodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	// allowanceBig, err := darkNodeRegistry.tokenBinding.Allowance(darkNodeRegistry.auth2, darkNodeRegistry.auth1.From, darkNodeRegistry.darkNodeRegistryAddress)
	// if err != nil {
	// 	return &types.Transaction{}, err
	// }
	// allowance, err := stackint.FromBigInt(allowanceBig)
	// if err != nil {
	// 	return &types.Transaction{}, err
	// }
	// if allowance.Cmp(bond) < 0 {
	// 	return &types.Transaction{}, errors.New("Not enough allowance to register a node")
	// }

	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := darkNodeRegistry.binding.Register(darkNodeRegistry.auth1, darkNodeIDByte, publicKey, bond.ToBigInt())
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return txn, err
	// return txn, err
}

// Deregister an existing dark node
func (darkNodeRegistry *DarkNodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistry.binding.Deregister(darkNodeRegistry.auth1, darkNodeIDByte)
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
	tx, err := darkNodeRegistry.binding.Refund(darkNodeRegistry.auth1, darkNodeIDByte)
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
	bond, err := darkNodeRegistry.binding.GetBond(darkNodeRegistry.auth2, darkNodeIDByte)
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
	return darkNodeRegistry.binding.IsRegistered(darkNodeRegistry.auth2, darkNodeIDByte)
}

// IsDeregistered returns true if the node is deregistered
func (darkNodeRegistry *DarkNodeRegistry) IsDeregistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistry.binding.IsDeregistered(darkNodeRegistry.auth2, darkNodeIDByte)
}

// ApproveRen doesn't actually talk to the DNR - instead it approved Ren to it
func (darkNodeRegistry *DarkNodeRegistry) ApproveRen(value *stackint.Int1024) error {
	txn, err := darkNodeRegistry.tokenBinding.Approve(darkNodeRegistry.auth1, darkNodeRegistry.client.DNRAddress, value.ToBigInt())
	if err != nil {
		return err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, txn)
	return err
}

// // IsDarkNodePendingRegistration returns true if the node will become registered at the next epoch
// func (darkNodeRegistry *DarkNodeRegistry) IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
// 	darkNodeIDByte, err := toByte(darkNodeID)
// 	if err != nil {
// 		return false, err
// 	}
// 	return darkNodeRegistry.binding.IsPendingRegistration(darkNodeRegistry.auth2, darkNodeIDByte)
// }

// CurrentEpoch returns the current epoch
func (darkNodeRegistry *DarkNodeRegistry) CurrentEpoch() (Epoch, error) {
	epoch, err := darkNodeRegistry.binding.CurrentEpoch(darkNodeRegistry.auth2)
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
		blockhash, &timestamp,
	}, nil
}

// Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
func (darkNodeRegistry *DarkNodeRegistry) Epoch() (*types.Transaction, error) {
	fmt.Println("Epoch called!")
	tx, err := darkNodeRegistry.binding.Epoch(darkNodeRegistry.auth1)
	if err != nil {
		return nil, err
	}
	_, err = darkNodeRegistry.client.PatchedWaitMined(darkNodeRegistry.context, tx)
	return tx, err
}

// WaitForEpoch guarantees that an Epoch as passed
func (darkNodeRegistry *DarkNodeRegistry) WaitForEpoch() (*types.Transaction, error) {
	currentEpoch, err := darkNodeRegistry.CurrentEpoch()
	if err != nil {
		return nil, err
	}
	nextEpoch := currentEpoch
	var tx *types.Transaction
	for currentEpoch.Blockhash == nextEpoch.Blockhash {
		tx, err = darkNodeRegistry.binding.Epoch(darkNodeRegistry.auth1)
		if err != nil {
			return nil, err
		}
		nextEpoch, err = darkNodeRegistry.CurrentEpoch()
		if err != nil {
			return nil, err
		}
		time.Sleep(time.Millisecond * 10)
	}
	return tx, nil
}

// // GetCommitment gets the signed commitment
// func (darkNodeRegistry *DarkNodeRegistry) GetCommitment(darkNodeID []byte) ([32]byte, error) {
// 	darkNodeIDByte, err := toByte(darkNodeID)
// 	if err != nil {
// 		return [32]byte{}, err
// 	}
// 	return darkNodeRegistry.binding.GetCommitment(darkNodeRegistry.auth2, darkNodeIDByte)
// }

// GetOwner gets the owner of the given dark node
func (darkNodeRegistry *DarkNodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return darkNodeRegistry.binding.GetOwner(darkNodeRegistry.auth2, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
func (darkNodeRegistry *DarkNodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return darkNodeRegistry.binding.GetPublicKey(darkNodeRegistry.auth2, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
func (darkNodeRegistry *DarkNodeRegistry) GetAllNodes() ([][]byte, error) {
	ret, err := darkNodeRegistry.binding.GetDarkNodes(darkNodeRegistry.auth2)
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
	bond, err := darkNodeRegistry.binding.MinimumBond(darkNodeRegistry.auth2)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(bond)
}

// MinimumEpochInterval gets the minimum epoch interval
func (darkNodeRegistry *DarkNodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
	interval, err := darkNodeRegistry.binding.MinimumEpochInterval(darkNodeRegistry.auth2)
	if err != nil {
		return stackint.Int1024{}, err
	}
	return stackint.FromBigInt(interval)
}

// SetGasLimit sets the gas limit to use for transactions
func (darkNodeRegistry *DarkNodeRegistry) SetGasLimit(limit uint64) {
	darkNodeRegistry.auth1.GasLimit = limit
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
