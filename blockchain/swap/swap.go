package swap

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/arc"
)

type SwapItem struct {
	orderID   [32]byte
	from      []byte
	to        []byte
	value     *big.Int
	expiry    int64
	arc       arc.Arc
	goesFirst bool
}
type Swap struct {
	hash    [32]byte
	secret  [32]byte
	fstItem SwapItem
	sndItem SwapItem
}

func NewSwap(fstItem, sndItem SwapItem) Swap {
	// Generating secret using random [32]byte

	secret := [32]byte{}
	hash := [32]byte{}
	if fstItem.goesFirst {
		rand.Read(secret[:])
		hash = sha256.Sum256(secret[:])
	}

	return Swap{
		hash:    hash,
		secret:  secret,
		fstItem: fstItem,
		sndItem: sndItem,
	}
}

func (swap *Swap) Execute(ctx context.Context) error {

	if swap.fstItem.goesFirst {
		if err := swap.fstItem.arc.Initiate(swap.hash, swap.fstItem.from, swap.fstItem.to, swap.fstItem.value, swap.fstItem.expiry); err != nil {
			return err
		}
		status := waitForTheOtherTraderToInitiate(swap)
		if status {
			return swap.fstItem.arc.Redeem(swap.sndItem.orderID, swap.secret)
		} else {
			return swap.fstItem.arc.Refund(swap.sndItem.orderID)
		}
	} else {
		hash, from, to, _, expiry, err := swap.sndItem.arc.Audit(swap.sndItem.orderID)
		if err := swap.fstItem.arc.Initiate(hash, to, from, swap.fstItem.value, expiry-(24*60*60)); err != nil {
			return err
		}
		status := waitForOtherTraderToRedeem(swap, expiry-(24*60*60))
		secret, err := swap.fstItem.arc.AuditSecret(swap.fstItem.orderID)
		if err != nil {
			return err
		}
		if status {
			// This should never return an error, if it happens it means someone is going to lose their funds.
			return swap.fstItem.arc.Redeem(swap.fstItem.orderID, secret)
		} else {
			return swap.fstItem.arc.Refund(swap.fstItem.orderID)
		}
	}
}

func waitForTheOtherTraderToInitiate(swap *Swap) bool {
	for {
		t := time.Now().Unix()
		_, _, _, _, expiry, err := swap.sndItem.arc.Audit(swap.sndItem.orderID)
		if err != nil {
			if expiry < t {
				return false
			}
			time.Sleep(time.Second)
			continue
		}
		return true
	}
}

func waitForOtherTraderToRedeem(swap *Swap, expiry int64) bool {
	for {
		t := time.Now().Unix()
		_, err := swap.sndItem.arc.AuditSecret(swap.sndItem.orderID)
		if err != nil {
			if expiry < t {
				return false
			}
			time.Sleep(time.Second)
			continue
		}
		return true
	}
}
