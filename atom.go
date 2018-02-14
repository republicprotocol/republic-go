package atom

import "math/big"

type Ledger int64

const (
	LedgerBitcoin  Ledger = 1
	LedgerEthereum Ledger = 2

	LedgerBitcoinTestnet  Ledger = -1
	LedgerEthereumRopsten Ledger = -2
)

type LedgerData []byte

type Atom struct {
	ID     []byte
	Ledger Ledger
	Data   LedgerData
}

type AtomContract interface {
	Initiate(hash, to, from []byte, value *big.Int, expiry int64) (err error)
	Read() (hash, to, from []byte, value *big.Int, expiry int64, err error)
	ReadSecret() (secret []byte, err error)
	Redeem() error
	Refund() error
}
