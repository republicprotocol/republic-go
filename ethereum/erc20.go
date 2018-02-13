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

func (contract *ERC20AtomContract) Redeem(secret []byte) error {
	panic("unimplemented")
}

// Refund will return the funds of an atomic swap, if the expiry period has passed
func (contract *ERC20AtomContract) Refund() error {
	panic("unimplemented")
}
