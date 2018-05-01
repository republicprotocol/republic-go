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
	OrderID   [32]byte
	From      []byte
	To        []byte
	Value     *big.Int
	Expiry    int64
	Arc       arc.Arc
	GoesFirst bool
}
type Swap struct {
	Hash    [32]byte
	Secret  [32]byte
	FstItem SwapItem
	SndItem SwapItem
}

func NewSwap(fstItem, sndItem SwapItem) Swap {
	// Generating secret using random [32]byte

	secret := [32]byte{}
	hash := [32]byte{}
	if fstItem.GoesFirst {
		rand.Read(secret[:])
		hash = sha256.Sum256(secret[:])
	}

	return Swap{
		Hash:    hash,
		Secret:  secret,
		FstItem: fstItem,
		SndItem: sndItem,
	}
}

func (swap *Swap) Execute(ctx context.Context) error {

	if swap.FstItem.GoesFirst {
		if err := swap.FstItem.Arc.Initiate(swap.Hash, swap.FstItem.From, swap.FstItem.To, swap.FstItem.Value, swap.FstItem.Expiry); err != nil {
			return err
		}
		status := waitForTheOtherTraderToInitiate(swap)
		if status {
			return swap.FstItem.Arc.Redeem(swap.SndItem.OrderID, swap.Secret)
		} else {
			return swap.FstItem.Arc.Refund(swap.SndItem.OrderID)
		}
	} else {
		hash, from, to, _, expiry, err := swap.SndItem.Arc.Audit(swap.SndItem.OrderID)
		if err := swap.FstItem.Arc.Initiate(hash, to, from, swap.FstItem.Value, expiry-(24*60*60)); err != nil {
			return err
		}
		status := waitForOtherTraderToRedeem(swap, expiry-(24*60*60))
		secret, err := swap.FstItem.Arc.AuditSecret(swap.FstItem.OrderID)
		if err != nil {
			return err
		}
		if status {
			// This should never return an error, if it happens it means someone is going to lose their funds.
			return swap.FstItem.Arc.Redeem(swap.FstItem.OrderID, secret)
		} else {
			return swap.FstItem.Arc.Refund(swap.FstItem.OrderID)
		}
	}
}

func waitForTheOtherTraderToInitiate(swap *Swap) bool {
	for {
		t := time.Now().Unix()
		_, _, _, _, expiry, err := swap.SndItem.Arc.Audit(swap.SndItem.OrderID)
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
		_, err := swap.SndItem.Arc.AuditSecret(swap.SndItem.OrderID)
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
