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
		order.TokensBTCDGX,
		order.TokensBTCREN,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]
	volume := RandomCoExp()

	ord := order.NewOrder(order.TypeLimit, parity, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Int63())
	return ord
}

// RandomOrder will generate a random order and its match.
func RandomOrderMatch() (order.Order, order.Order) {
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensBTCDGX,
		order.TokensBTCREN,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]
	volume := RandomCoExp()

	buy := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Int63())
	sell := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(1*time.Hour), tokens, buy.Price, buy.Volume, buy.MinimumVolume, buy.Nonce)
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

func LessRandomCoExp(coexp order.CoExp) order.CoExp {
	co := uint64(rand.Intn(int(coexp.Co)) + 1)
	exp := uint64(rand.Intn(int(coexp.Exp + 1)))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}
