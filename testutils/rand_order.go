package testutils

import (
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// RandomOrder will generate a random order.
func RandomOrder() order.Order {
	parity := []order.Parity{order.ParityBuy, order.ParitySell}[rand.Intn(2)]
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(4)]
	volume := RandomCoExp()

	ord := order.NewOrder(order.TypeLimit, parity, order.SettlementRenEx, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Uint64())
	return ord
}

// RandomBuyOrder will generate a random buy order.
func RandomBuyOrder() order.Order {
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(4)]
	volume := RandomCoExp()

	ord := order.NewOrder(order.TypeLimit, order.ParityBuy, order.SettlementRenEx, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Uint64())
	return ord
}

// RandomBuyOrderFragments will generate order fragments for a random buy
// order.
func RandomBuyOrderFragments(n, k int64) ([]order.Fragment, error) {
	ord := RandomBuyOrder()
	frags, err := ord.Split(n, k)
	return frags, err
}

// RandomSellOrder will generate a random sell order.
func RandomSellOrder() order.Order {
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(4)]
	volume := RandomCoExp()

	ord := order.NewOrder(order.TypeLimit, order.ParitySell, order.SettlementRenEx, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Uint64())
	return ord
}

// RandomSellOrderFragments will generate order fragments for a random buy
// order.
func RandomSellOrderFragments(n, k int64) ([]order.Fragment, error) {
	ord := RandomSellOrder()
	frags, err := ord.Split(n, k)
	return frags, err
}

// RandomOrderMatch will generate a random order and its match.
func RandomOrderMatch() (order.Order, order.Order) {
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(4)]
	volume := RandomCoExp()

	buy := order.NewOrder(order.TypeLimit, order.ParityBuy, order.SettlementRenEx, time.Now().Add(24*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Uint64())
	sell := order.NewOrder(order.TypeLimit, order.ParitySell, order.SettlementRenEx, time.Now().Add(24*time.Hour), tokens, buy.Price, buy.Volume, buy.MinimumVolume, buy.Nonce)
	return buy, sell
}

// RandomCoExp will generate a random number represented in CoExp format.
func RandomCoExp() order.CoExp {
	co := uint64(rand.Intn(1999) + 1)
	exp := uint64(rand.Intn(25))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}

// LessRandomCoExp will generate a random CoExp that is no more than the given CoExp.
func LessRandomCoExp(coExp order.CoExp) order.CoExp {
	co := uint64(rand.Intn(int(coExp.Co)) + 1)
	exp := uint64(rand.Intn(int(coExp.Exp + 1)))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}
