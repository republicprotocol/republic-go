package interop

import (
	"math/big"
)

/*
 * Steps in an Arc swap:
 *
 * 0. (A) and (B) share addresses, (A) creates HASH from SECRET
 * 1. (A) calls Initiate(HASH, details) to ADDR1, gives (HASH, ADDR1) to (B)
 * 2. (B) calls ADDR1.Audit(HASH, details)
 * 3. (B) calls Initiate(HASH) to ADDR2, gives (ADDR2) to (A)
 * 4. (A) calls ADDR2.Audit(HASH, details)
 * 5. (A) calls ADDR2.Redeem(SECRET)
 * 6. (B) calls ADDR2.AuditSecret(), retrieving SECRET
 * 7. (B) calls ADDR1.Redeem(SECRET)
 */

// Arc is the interface defining the Atomic Swap Interface
type Arc interface {
	Initiate(hash [32]byte, from, to []byte, value *big.Int, expiry int64) error
	Redeem(secret []byte) error
	Refund() error
	Audit() (hash [32]byte, from, to []byte, value *big.Int, expiry int64, err error)
	AuditSecret() (secret []byte, err error)
	Serialize() ([]byte, error)
	Deserialize([]byte) error
}
