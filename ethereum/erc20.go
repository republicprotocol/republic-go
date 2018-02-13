package ethereum

// ERC20AtomContract ...
type ERC20AtomContract struct {
}

// NewERC20AtomContract returns a new NewERC20Atom instance
func NewERC20AtomContract() *ERC20AtomContract {
	return &ERC20AtomContract{}
}

// Initiate starts or reciprocates an atomic swap
func (contract *ERC20AtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	panic("unimplemented")
}

// Read returns details about an atomic swap
func (contract *ERC20AtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	panic("unimplemented")
}

// ReadSecret returns the secret of an atomic swap if it's available
func (contract *ERC20AtomContract) ReadSecret() (secret []byte, err error) {
	panic("unimplemented")
}

// Redeem closes an atomic swap by revealing the secret
func (contract *ERC20AtomContract) Redeem() error {
	panic("unimplemented")
}

// Refund will return the funds of an atomic swap, if the expiry period has passed
func (contract *ERC20AtomContract) Refund() error {
	panic("unimplemented")
}

/*

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
*/
