package ethereum

import (
	"crypto/rand"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/go-atom/ethereum/contracts"
	// "github.com/ethereum/go-ethereum/common"
	// "github.com/ethereum/go-ethereum/core/types"
)

// ETHAtomContract ...
type ETHAtomContract struct {
	client     bind.ContractBackend
	auth       *bind.TransactOpts
	binding    *contracts.AtomicSwapEther
	swapID     [32]byte
	secretHash [32]byte
}

// NewETHAtomContract returns a new NewETHAtom instance
func NewETHAtomContract(client bind.ContractBackend, auth1 *bind.TransactOpts, address common.Address) *ETHAtomContract {
	contract, err := contracts.NewAtomicSwapEther(address, client)
	if err != nil {
		log.Fatalf("%v", err)
	}

	swapID := [32]byte{}
	_, err = rand.Read(swapID[:])
	if err != nil {
		panic(err)
	}

	secretHash := [32]byte{}
	_, err = rand.Read(secretHash[:])
	if err != nil {
		panic(err)
	}

	return &ETHAtomContract{
		client:     client,
		auth:       auth1,
		binding:    contract,
		swapID:     swapID,
		secretHash: secretHash,
	}
}

// Initiate starts or reciprocates an atomic swap
func (contract *ETHAtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	authWithValue := contract.auth
	authWithValue.Value = big.NewInt(value)
	toAddress := common.BytesToAddress(to)
	_, err = contract.binding.Open(authWithValue, contract.swapID, toAddress, contract.secretHash)

	return err
}

// Read returns details about an atomic swap
func (contract *ETHAtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	ret, err := contract.binding.Check(&bind.CallOpts{}, contract.swapID)
	return ret.SecretLock[:],
		ret.WithdrawTrader.Bytes(),
		nil,
		ret.Value.Int64(),
		ret.TimeRemaining.Int64(),
		err
}

// ReadSecret returns the secret of an atomic swap if it has been revealed
func (contract *ETHAtomContract) ReadSecret() (secret []byte, err error) {
	return contract.binding.CheckSecretKey(&bind.CallOpts{}, contract.swapID)
}

// func (connection EtherConnection) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {

// Redeem closes an atomic swap by revealing the secret
func (contract *ETHAtomContract) Redeem(secret []byte) error {
	_, err := contract.binding.Close(contract.auth, contract.swapID, secret)
	return err
}

// Refund will return the funds of an atomic swap, provided the expiry period has passed
func (contract *ETHAtomContract) Refund() error {
	_, err := contract.binding.Expire(contract.auth, contract.swapID)
	return err
}
