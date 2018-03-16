package dnr

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/republicprotocol/republic-go/contracts/bindings"
	"github.com/republicprotocol/republic-go/contracts/connection"
)

// DarkNodeRegistrar is the interface defining the Dark Node Registrar
type DarkNodeRegistrar interface {
	Register(darkNodeID []byte, publicKey []byte, bond *big.Int) (*types.Transaction, error)
	Deregister(darkNodeID []byte) (*types.Transaction, error)
	Refund(darkNodeID []byte) (*types.Transaction, error)
	Epoch() (*types.Transaction, error)

	CurrentEpoch() (Epoch, error)
	GetBond(darkNodeID []byte) (*big.Int, error)
	GetCommitment(darkNodeID []byte) ([32]byte, error)
	GetOwner(darkNodeID []byte) (common.Address, error)
	GetPublicKey(darkNodeID []byte) ([]byte, error)
	GetAllNodes() ([][]byte, error)

	MinimumBond() (*big.Int, error)
	MinimumEpochInterval() (*big.Int, error)

	IsDarkNodeRegistered(darkNodeID []byte) (bool, error)
	IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error)
	WaitUntilRegistration(darkNodeID []byte) error
}

// Epoch contains a blockhash and a timestamp
type Epoch struct {
	Blockhash [32]byte
	Timestamp *big.Int
}

// EthereumDarkNodeRegistrar is the dark node interface
type EthereumDarkNodeRegistrar struct {
	context                  context.Context
	client                   *connection.Client
	auth1                    *bind.TransactOpts
	auth2                    *bind.CallOpts
	binding                  *bindings.DarkNodeRegistrar
	tokenBinding             *bindings.ERC20
	darkNodeRegistrarAddress common.Address
}

// NewEthereumDarkNodeRegistrar returns a Dark node registrar
func NewEthereumDarkNodeRegistrar(context context.Context, clientDetails *connection.ClientDetails, auth1 *bind.TransactOpts, auth2 *bind.CallOpts) (DarkNodeRegistrar, error) {
	contract, err := bindings.NewDarkNodeRegistrar(clientDetails.DNRAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return nil, err
	}
	renContract, err := bindings.NewERC20(clientDetails.RenAddress, bind.ContractBackend(clientDetails.Client))
	if err != nil {
		return nil, err
	}
	return &EthereumDarkNodeRegistrar{
		context:                  context,
		client:                   &clientDetails.Client,
		auth1:                    auth1,
		auth2:                    auth2,
		binding:                  contract,
		tokenBinding:             renContract,
		darkNodeRegistrarAddress: clientDetails.DNRAddress,
	}, nil
}

// Register registers a new dark node
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) Register(darkNodeID []byte, publicKey []byte, bond *big.Int) (*types.Transaction, error) {
	allowance, err := darkNodeRegistrar.tokenBinding.Allowance(darkNodeRegistrar.auth2, darkNodeRegistrar.auth1.From, darkNodeRegistrar.darkNodeRegistrarAddress)
	if err != nil {
		return &types.Transaction{}, err
	}
	if allowance.Cmp(bond) < 0 {
		return &types.Transaction{}, errors.New("Not enough allowance to register a node")
	}

	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}

	txn, err := darkNodeRegistrar.binding.Register(darkNodeRegistrar.auth1, darkNodeIDByte, publicKey, bond)
	if err == nil {
		_, err := connection.PatchedWaitMined(darkNodeRegistrar.context, *darkNodeRegistrar.client, txn)
		return txn, err
	}
	return txn, err
}

// Deregister deregisters an existing dark node
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	tx, err := darkNodeRegistrar.binding.Deregister(darkNodeRegistrar.auth1, darkNodeIDByte)
	if err != nil {
		return tx, err
	}
	_, err = connection.PatchedWaitMined(darkNodeRegistrar.context, *darkNodeRegistrar.client, tx)
	return tx, err
}

// GetBond gets the bond of an existing dark node
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) GetBond(darkNodeID []byte) (*big.Int, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &big.Int{}, err
	}
	return darkNodeRegistrar.binding.GetBond(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// IsDarkNodeRegistered check's whether a dark node is registered or not
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) IsDarkNodeRegistered(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistrar.binding.IsDarkNodeRegistered(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// IsDarkNodePendingRegistration returns true if the node will be registered in the next epoch
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return false, err
	}
	return darkNodeRegistrar.binding.IsDarkNodePendingRegistration(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// CurrentEpoch returns the current epoch
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) CurrentEpoch() (Epoch, error) {
	return darkNodeRegistrar.binding.CurrentEpoch(darkNodeRegistrar.auth2)
}

// Epoch updates the current Epoch
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) Epoch() (*types.Transaction, error) {
	tx, err := darkNodeRegistrar.binding.Epoch(darkNodeRegistrar.auth1)
	if err != nil {

		return nil, err
	}
	_, err = connection.PatchedWaitMined(darkNodeRegistrar.context, *darkNodeRegistrar.client, tx)
	return tx, err
}

// GetCommitment gets the signed commitment
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) GetCommitment(darkNodeID []byte) ([32]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return [32]byte{}, err
	}
	return darkNodeRegistrar.binding.GetCommitment(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// GetOwner gets the owner of the given dark node
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) GetOwner(darkNodeID []byte) (common.Address, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return common.Address{}, err
	}
	return darkNodeRegistrar.binding.GetOwner(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// GetPublicKey gets the public key of the goven dark node
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return []byte{}, err
	}
	return darkNodeRegistrar.binding.GetPublicKey(darkNodeRegistrar.auth2, darkNodeIDByte)
}

// GetAllNodes gets all dark nodes
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) GetAllNodes() ([][]byte, error) {
	ret, err := darkNodeRegistrar.binding.GetXingOverlay(darkNodeRegistrar.auth2)
	if err != nil {
		return nil, err
	}
	arr := make([][]byte, len(ret))
	for i, b20 := range ret {
		arr[i] = b20[:]
	}
	return arr, nil
}

// MinimumBond gets the minimum viable bonda mount
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) MinimumBond() (*big.Int, error) {
	return darkNodeRegistrar.binding.MinimumBond(darkNodeRegistrar.auth2)
}

// MinimumEpochInterval gets the minimum epoch interval
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) MinimumEpochInterval() (*big.Int, error) {
	return darkNodeRegistrar.binding.MinimumEpochInterval(darkNodeRegistrar.auth2)
}

// Refund refunds the bond of an unregistered miner
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) Refund(darkNodeID []byte) (*types.Transaction, error) {
	darkNodeIDByte, err := toByte(darkNodeID)
	if err != nil {
		return &types.Transaction{}, err
	}
	return darkNodeRegistrar.binding.Refund(darkNodeRegistrar.auth1, darkNodeIDByte)
}

// WaitUntilRegistration waits until the registration is successful.
func (darkNodeRegistrar *EthereumDarkNodeRegistrar) WaitUntilRegistration(darkNodeID []byte) error {
	isRegistered := false
	for !isRegistered {
		tx, err := darkNodeRegistrar.Epoch()
		if err != nil {
			return err
		}
		_, err = connection.PatchedWaitMined(darkNodeRegistrar.context, *darkNodeRegistrar.client, tx)
		if err != nil {
			return err
		}
		time.Sleep(time.Minute)
		isRegistered, err = darkNodeRegistrar.IsDarkNodeRegistered(darkNodeID)
		if err != nil {
			return err
		}
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
