package go_eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// The Delegate is used as a callback interface to inject logic into the
// different RPCs.
type Delegate interface {
}

// type interface interface {
// 	//
// }

// Open opens an order
func (connection EtherConnection) Open(ethAddr common.Address, ethAmount uint64, secretHash [32]byte) ([32]byte, error) {
	id := [32]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	_, err := connection.contract.Open(connection.auth, id, ethAddr, secretHash)
	return id, err
}

// Close closes an order
func (connection EtherConnection) Close() {
}

func (connection EtherConnection) Check(id [32]byte) (struct {
	TimeRemaining  *big.Int
	Value          *big.Int
	WithdrawTrader common.Address
	SecretLock     [32]byte
}, error) {
	return connection.contract.Check(&bind.CallOpts{}, id)
}

// Expire expires an order
func (connection EtherConnection) Expire() {
}

// Validate validates an order
func (connection EtherConnection) Validate() {
}

// RetrieveSecretKey retrieves a secret key from an order
func (connection EtherConnection) RetrieveSecretKey() {
}
