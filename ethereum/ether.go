package ethereum

// ETHAtomContract ...
type ETHAtomContract struct {
}

// NewETHAtomContract returns a new NewETHAtom instance
func NewETHAtomContract() *ETHAtomContract {
	return &ETHAtomContract{}
}

// Initiate starts or reciprocates an atomic swap
func (contract *ETHAtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	panic("unimplemented")
}

// Read returns details about an atomic swap
func (contract *ETHAtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	panic("unimplemented")
}

// ReadSecret returns the secret of an atomic swap if it's available
func (contract *ETHAtomContract) ReadSecret() (secret []byte, err error) {
	panic("unimplemented")
}

// Redeem closes an atomic swap by revealing the secret
func (contract *ETHAtomContract) Redeem() error {
	panic("unimplemented")
}

// Refund will return the funds of an atomic swap, if the expiry period has passed
func (contract *ETHAtomContract) Refund() error {
	panic("unimplemented")
}
