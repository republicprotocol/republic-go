package bitcoin

type BTCAtomContract struct {
}

func NewBTCAtomContract() *BTCAtomContract {
	return &BTCAtomContract{}
}

func (contract *BTCAtomContract) Initiate(hash, to, from []byte, value, expiry int64) (err error) {
	panic("unimplemented")
}

func (contract *BTCAtomContract) Read() (hash, to, from []byte, value, expiry int64, err error) {
	panic("unimplemented")
}

func (contract *BTCAtomContract) ReadSecret() (secret []byte, err error) {
	panic("unimplemented")
}

func (contract *BTCAtomContract) Redeem(secret []byte) error {
	panic("unimplemented")
}

func (contract *BTCAtomContract) Refund() error {
	panic("unimplemented")
}
