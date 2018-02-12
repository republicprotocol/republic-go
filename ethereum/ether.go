package ethereum

type ETHAtomContract struct {
}

func NewETHAtomContract() *ETHAtomContract {
	return &ETHAtomContract{}
}

func (contract *ETHAtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	panic("unimplemented")
}

func (contract *ETHAtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	panic("unimplemented")
}

func (contract *ETHAtomContract) ReadSecret() (secret []byte, err error) {
	panic("unimplemented")
}

func (contract *ETHAtomContract) Redeem() error {
	panic("unimplemented")
}

func (contract *ETHAtomContract) Refund() error {
	panic("unimplemented")
}
