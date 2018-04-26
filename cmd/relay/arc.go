package main

import (
	"math/big"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
)

type Swap struct {
	secret [32]byte
	hash [32]byte
	from []byte
	to []byte
	value *big.Int
	expiry int64
	arc arc.Arc
}

func NewSwap() Swap {

}

func processAtomicSwaps(swaps <-chan Swap) <-chan error {
	errs := make(chan error)
	go func() {
		for {
			select{
			case swap := <-swaps:
				err := swap.arc.Initiate(swap.hash, swap.from, swap.to, swap.value, swap.expiry)
				if err != nil {
					errs <- err
					continue
				}
				hash, to, from, value, expiry, err := swap.arc.Audit()
				if (swap.hash != hash || swap.from != to || swap.to != from) { // FIXME: Validate audit
					errs <- err
					continue
				}

				secret := swap.secret
				if len(secret) == 0 {
					secret, err := arc.AuditSecret()
					if err != nil {
						errs <- err
						continue
					}
				}
				arc.Redeem(secret)
			}
		}
	}()
	return errs
}

func handleAtomicSwap(swap arc.Arc, hash [32]byte, fro) error {
	swap.Initiate(hash, from []byte, to []byte, value *big.Int, expiry int64);
}

//	func atomicSwap(entries <-chan orderbook.Entry, privateKey *ecdsa.PrivateKey) error {
