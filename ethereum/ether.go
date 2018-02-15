package ethereum

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
)

// BytesTo32Bytes ...
func BytesTo32Bytes(bytes []byte) ([32]byte, error) {
	var bytes32 [32]byte
	if len(bytes) != 32 {
		return bytes32, errors.New("Expected 32 bytes")
	}
	for i := 0; i < 32; i++ {
		bytes32[i] = bytes[i]
	}

	return bytes32, nil
}

// ETHAtomContract ...
type ETHAtomContract struct {
	context context.Context
	client  Client
	auth    *bind.TransactOpts
	binding *contracts.AtomicSwapEther
	swapID  [32]byte
	chainID int8
}

// NewETHAtomContract returns a new NewETHAtom instance
func NewETHAtomContract(context context.Context, client Client, auth1 *bind.TransactOpts, address common.Address, data []byte) (*ETHAtomContract, error) {
	contract, err := contracts.NewAtomicSwapEther(address, bind.ContractBackend(client))
	if err != nil {
		return nil, err
	}

	var swapID [32]byte
	if data == nil {
		swapID = [32]byte{}
		_, err = rand.Read(swapID[:])
		if err != nil {
			return nil, err
		}
	} else {
		swapID, err = BytesTo32Bytes(data)
	}

	return &ETHAtomContract{
		context: context,
		client:  client,
		auth:    auth1,
		binding: contract,
		swapID:  swapID,
	}, nil
}

// Initiate starts or reciprocates an atomic swap
func (contract *ETHAtomContract) Initiate(hash, to, from []byte, value *big.Int, expiry int64) (err error) {
	hash32, err := BytesTo32Bytes(hash)
	if err != nil {
		log.Fatalf("Expected 32 bytes: %v", err)
	}
	authWithValue := contract.auth
	authWithValue.Value = value
	toAddress := common.BytesToAddress(to)
	tx, err := contract.binding.Open(authWithValue, contract.swapID, toAddress, hash32, big.NewInt(expiry))
	if err != nil {
		return err
	}
	_, err = PatchedWaitMined(contract.context, contract.client, tx)
	return err
}

// Read returns details about an atomic swap
func (contract *ETHAtomContract) Read() (hash, to, from []byte, value *big.Int, expiry int64, err error) {
	ret, err := contract.binding.Check(&bind.CallOpts{}, contract.swapID)
	return ret.SecretLock[:],
		ret.WithdrawTrader.Bytes(),
		nil,
		ret.Value,
		ret.Timelock.Int64(),
		err
}

// ReadSecret returns the secret of an atomic swap if it has been revealed
func (contract *ETHAtomContract) ReadSecret() (secret []byte, err error) {
	return contract.binding.CheckSecretKey(&bind.CallOpts{}, contract.swapID)
}

// func (connection EtherConnection) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {

// Redeem closes an atomic swap by revealing the secret
func (contract *ETHAtomContract) Redeem(secret []byte) error {
	tx, err := contract.binding.Close(contract.auth, contract.swapID, secret)
	if err != nil {
		return err
	}
	_, err = PatchedWaitMined(contract.context, contract.client, tx)
	return err
}

// Refund will return the funds of an atomic swap, provided the expiry period has passed
func (contract *ETHAtomContract) Refund() error {
	tx, err := contract.binding.Expire(contract.auth, contract.swapID)
	if err != nil {
		return err
	}
	_, err = PatchedWaitMined(contract.context, contract.client, tx)
	return err
}

// GetData returns the data required for another party to participate in an atomic swap
func (contract *ETHAtomContract) GetData() []byte {
	return contract.swapID[:]
}
