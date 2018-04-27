package main

import (
	"math/big"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/arc"
)

/** 
type SwapItem struct {
	orderID   [32]byte
	from      []byte
	to        []byte
	value     *big.Int
	expiry    int64
	arc       arc.Arc
	goesFirst bool
}
*/
func initSwap(fstOrderID, sndOrderID [32]byte, fstOrder []byte, fstValue, sndValue *big.Int, from, to []byte, fstTokenAddress, sndTokenAddress common.Address) Swap {
	fee := 2
	if fstTokenAddress != swap.ETHEREUM {
		fstTokenAddress 
	}

	fstItem = arc.SwapItem{
		orderID: fstOrderID,
		from: from,
		to: to,
		value: fstValue,
		expiry: time.Now().Unix() + 48*60*60,
		arc: arc.NewArc(context.Background(), conn, auth, fstOrder, fstTokenAddress, fee),
		goesFirst: true,
	}
	sndItem = arc.SwapItem{
		orderID: sndOrderID,
		from: to,
		to: from,
		value: sndValue,
		expiry: time.Now().Unix() + 24*60*60,
		arc: arc.NewArc(context.Background(), conn, auth, byte[]{}, sndTokenAddress, fee),
		goesFirst: false,
	}

	return arc.NewSwap(fstItem, sndItem)
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
