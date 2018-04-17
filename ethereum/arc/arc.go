package arc

import (
	"math/big"
)

// Arc is the interface defining the Atomic Swap Interface
type Arc interface {
	Initiate(hash, to, from []byte, value *big.Int, expiry int64) error
	Redeem(secret []byte) error
	Refund() error
	Audit() (hash, to, from []byte, value *big.Int, expiry int64, err error)
	AuditSecret() (secret []byte, err error)
}
