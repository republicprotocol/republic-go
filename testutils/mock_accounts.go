package testutils

import (
	"sync"

	"github.com/republicprotocol/republic-go/order"
)

// NewDarkpoolAccounts is a mock implementation of the cal.DarkpoolAccounts.
type DarkpoolAccounts struct {
	mu    *sync.Mutex
	comps int
	buys  map[order.ID]struct{}
	sells map[order.ID]struct{}
}

// NewDarkpoolAccounts creates a new DarkpoolAccounts
func NewDarkpoolAccounts() *DarkpoolAccounts {
	return &DarkpoolAccounts{
		mu:    new(sync.Mutex),
		comps: 0,
		buys:  map[order.ID]struct{}{},
		sells: map[order.ID]struct{}{},
	}
}

// Settle implements the Settle function of cal.DarkpoolAccounts interface
func (accounts *DarkpoolAccounts) Settle(buy order.Order, sell order.Order) error {
	accounts.mu.Lock()
	defer accounts.mu.Unlock()

	accounts.comps++
	accounts.buys[buy.ID] = struct{}{}
	accounts.sells[buy.ID] = struct{}{}

	return nil
}

func (accounts *DarkpoolAccounts) SettleCounts() int {
	accounts.mu.Lock()
	defer accounts.mu.Unlock()

	return accounts.comps
}

//// ErrInsufficientBalance is returned when the buyer or seller doesn't have
//// enough balance in the contract to settle the match.
//var ErrInsufficientBalance = errors.New("insufficient balance")
//
//// ErrTraderNotFound is returned when the trader is found in the accounts.
//var ErrTraderNotFound = errors.New("trader not found")
//
//// initialBalance returns a map representing the balance of a trader.
//var initialBalance = map[order.Token]float64{
//	order.TokenBTC: 1000.0,
//	order.TokenETH: 1000.0,
//	order.TokenDGX: 1000.0,
//	order.TokenREN: 1000.0,
//}
//
//// DarkpoolAccounts is a mock implementation of the cal.DarkpoolAccounts which
//// assume there are only one buyer and one seller. Each of them has 1k balance
//// for each token.
//type DarkpoolAccounts struct {
//	buyer  map[order.Token]float64
//	seller map[order.Token]float64
//}
//
//// NewDarkpoolAccounts creates a new mock cal.DarkpoolAccounts.
//func NewDarkpoolAccounts() cal.DarkpoolAccounts {
//	return &DarkpoolAccounts{
//		buyer:  initialBalance,
//		seller: initialBalance,
//	}
//}
//
//// Settle the order pair which gets confirmed.
//func (accounts *DarkpoolAccounts) Settle(buy order.Order, sell order.Order) error {
//	// Calculate price and volume
//	price := (float64(buy.Price.Co)*0.005*math.Pow(10, float64(buy.Price.Exp-26)) + (float64(sell.Price.Co) * 0.005 * math.Pow(10, float64(sell.Price.Exp-26)))) / 2
//	buyVolume := float64(buy.Volume.Co) * 0.2 * math.Pow(10, float64(buy.Volume.Exp))    // in 10^-12  btc
//	sellVolume := float64(sell.Volume.Co) * 0.2 * math.Pow(10, float64(sell.Volume.Exp)) // in 10^-12 ren
//	volume := math.Min(buyVolume/price, sellVolume)                                      // in 10^-12 ren
//
//	// Update the balance
//	accounts.buyer[buy.Tokens.NonPriorityToken()] -= volume * math.Pow(10, 12)
//	if accounts.buyer[buy.Tokens.NonPriorityToken()] < 0 {
//		return ErrInsufficientBalance
//	}
//	accounts.buyer[buy.Tokens.PriorityToken()] += volume * math.Pow(10, 12) * price
//	accounts.seller[sell.Tokens.PriorityToken()] -= volume * math.Pow(10, 12)
//	if accounts.seller[sell.Tokens.PriorityToken()] < 0 {
//		return ErrInsufficientBalance
//	}
//	accounts.seller[sell.Tokens.NonPriorityToken()] += volume * math.Pow(10, 12) * price
//
//	return nil
//}
//
//// Balance returns the balance of a trader for a particular order.Token.
//func (accounts *DarkpoolAccounts) Balance(trader string, token order.Token) (float64, error) {
//	if trader == "buyer" {
//		return accounts.buyer[token], nil
//	} else if trader == "seller" {
//		return accounts.seller[token], nil
//	} else {
//		return 0, ErrTraderNotFound
//	}
//}
