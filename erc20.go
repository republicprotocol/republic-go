package go_eth

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Open opens an Atomic swap for a given match ID, with a address authorised to withdraw the amount after revealing the secret
func (connection ERC20Connection) Open(_swapID [32]byte, ethAddr common.Address, secretHash [32]byte, amountInWei *big.Int) (*types.Transaction, error) {
	return connection.contract.Open(connection.auth, _swapID, amountInWei, ethAddr, ethAddr, secretHash)
}

// Close closes an Atomic swap by revealing the secret. The locked value is sent to the address supplied to Open
func (connection ERC20Connection) Close(_swapID [32]byte, _secretKey []byte) (*types.Transaction, error) {
	return connection.contract.Close(connection.auth, _swapID, _secretKey)
}

// Check returns details about an open Atomic Swap
func (connection ERC20Connection) Check(id [32]byte) (struct {
	TimeRemaining        *big.Int
	Erc20Value           *big.Int
	Erc20ContractAddress common.Address
	WithdrawTrader       common.Address
	SecretLock           [32]byte
}, error) {
	return connection.contract.Check(&bind.CallOpts{}, id)
}

// Expire expires an Atomic Swap, provided that the required time has passed
func (connection ERC20Connection) Expire(_swapID [32]byte) (*types.Transaction, error) {
	return connection.contract.Expire(connection.auth, _swapID)
}

// Validate (not implemented) checks that there is a valid open Atomic Swap for a given _swapID
func (connection ERC20Connection) Validate() {
}

// RetrieveSecretKey retrieves the secret key from an Atomic Swap, after it has been revealed
func (connection ERC20Connection) RetrieveSecretKey(_swapID [32]byte) ([]byte, error) {
	return connection.contract.CheckSecretKey(&bind.CallOpts{}, _swapID)
}
