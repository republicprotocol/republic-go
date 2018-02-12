package ethereum

type ERC20AtomContract struct {
}

func NewERC20AtomContract() *ERC20AtomContract {
	return &ERC20AtomContract{}
}

func (contract *ERC20AtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	panic("unimplemented")
}

func (contract *ERC20AtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	panic("unimplemented")
}

func (contract *ERC20AtomContract) ReadSecret() (secret []byte, err error) {
	panic("unimplemented")
}

func (contract *ERC20AtomContract) Redeem() error {
	panic("unimplemented")
}

func (contract *ERC20AtomContract) Refund() error {
	panic("unimplemented")
}
