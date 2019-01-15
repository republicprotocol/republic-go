package testutils

import (
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/order"
)

// Tokens values.
const (
	TokensBTCETH  = order.Tokens((uint64(order.TokenBTC) << 32) | uint64(order.TokenETH))
	TokensBTCTUSD = order.Tokens((uint64(order.TokenBTC) << 32) | uint64(order.TokenTUSD))
	TokensBTCREN  = order.Tokens((uint64(order.TokenBTC) << 32) | uint64(order.TokenREN))
	TokensETHDGX  = order.Tokens((uint64(order.TokenETH) << 32) | uint64(order.TokenDGX))
	TokensETHTUSD = order.Tokens((uint64(order.TokenETH) << 32) | uint64(order.TokenTUSD))
	TokensETHREN  = order.Tokens((uint64(order.TokenETH) << 32) | uint64(order.TokenREN))
	TokensETHZRX  = order.Tokens((uint64(order.TokenETH) << 32) | uint64(order.TokenZRX))
	TokensETHOMG  = order.Tokens((uint64(order.TokenETH) << 32) | uint64(order.TokenOMG))
)

var tokens = []order.Tokens{
	TokensBTCETH,
	TokensETHDGX,
	TokensETHREN,
	TokensETHTUSD,
	TokensETHZRX,
	TokensETHOMG,
}

// RandomOrder will generate a random order.
func RandomOrder() order.Order {
	parity := []order.Parity{order.ParityBuy, order.ParitySell}[rand.Intn(2)]
	tokens := []order.Tokens{TokensBTCETH,
		TokensETHDGX,
		TokensETHREN,
		TokensETHTUSD,
		TokensETHZRX,
		TokensETHOMG,
	}
	token := tokens[rand.Intn(len(tokens))]

	ord := order.NewOrder(parity, order.TypeLimit, time.Now().Add(1*time.Hour), order.SettlementRenEx, token, rand.Uint64(), rand.Uint64(), 0, uint64(rand.Int63()))
	return ord
}

// RandomBuyOrder will generate a random buy order.
func RandomBuyOrder() order.Order {
	tokens := []order.Tokens{TokensBTCETH,
		TokensETHDGX,
		TokensETHREN,
		TokensETHTUSD,
		TokensETHZRX,
		TokensETHOMG,
	}
	token := tokens[rand.Intn(len(tokens))]

	ord := order.NewOrder(order.ParityBuy, order.TypeLimit, time.Now().Add(1*time.Hour), order.SettlementRenEx, token, rand.Uint64(), rand.Uint64(), 0, uint64(rand.Int63()))
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
	tokens := []order.Tokens{TokensBTCETH,
		TokensETHDGX,
		TokensETHREN,
		TokensETHTUSD,
		TokensETHZRX,
		TokensETHOMG,
	}
	token := tokens[rand.Intn(len(tokens))]

	ord := order.NewOrder(order.ParitySell, order.TypeLimit, time.Now().Add(1*time.Hour), order.SettlementRenEx, token, rand.Uint64(), rand.Uint64(), 0, uint64(rand.Int63()))
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
	tokens := []order.Tokens{TokensBTCETH,
		TokensETHDGX,
		TokensETHREN,
		TokensETHTUSD,
		TokensETHZRX,
		TokensETHOMG,
	}
	token := tokens[rand.Intn(len(tokens))]

	price := rand.Uint64()
	volume := rand.Uint64()
	buy := order.NewOrder(order.ParityBuy, order.TypeLimit, time.Now().Add(1*time.Hour), order.SettlementRenEx, token, price, volume, 0, uint64(rand.Int63()))
	sell := order.NewOrder(order.ParitySell, order.TypeLimit, time.Now().Add(1*time.Hour), order.SettlementRenEx, token, price, volume, 0, uint64(rand.Int63()))
	return buy, sell
}

// RandomCoExp will generate a random number represented in CoExp format.
func RandomCoExp() order.CoExp {
	co := uint64(rand.Intn(1999) + 1)
	exp := uint64(rand.Intn(27))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}
