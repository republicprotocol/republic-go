package main

import (
	"log"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-miner"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order-compute"
)

var orders = []*compute.Order{
	// Buy ETH for REN
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderBuy,
		compute.CurrencyCodeETH,
		compute.CurrencyCodeREN,
		17200,
		1000,
		1000,
		0,
	),
	// Sell ETH for REN
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderSell,
		compute.CurrencyCodeETH,
		compute.CurrencyCodeREN,
		17200,
		1000,
		1000,
		0,
	),

	// Buy BTC for ETH
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderBuy,
		compute.CurrencyCodeBTC,
		compute.CurrencyCodeETH,
		12,
		2000,
		200,
		0,
	),
	// Sell BTC for ETH
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderSell,
		compute.CurrencyCodeBTC,
		compute.CurrencyCodeETH,
		8,
		200,
		100,
		0,
	),

	// Buy BTC for REN
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderBuy,
		compute.CurrencyCodeBTC,
		compute.CurrencyCodeREN,
		172000,
		1000,
		1000,
		0,
	),
	// Sell REN for BTC
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderSell,
		compute.CurrencyCodeREN,
		compute.CurrencyCodeBTC,
		172000,
		1000,
		1000,
		0,
	),
}

func main() {

	multiAddress1, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/3000/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto")
	if err != nil {
		log.Fatal(err)
	}
	multiAddress2, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/3001/republic/8MHarRJdvWd7SsTJE8vRVfj2jb5cWS")
	if err != nil {
		log.Fatal(err)
	}
	multiAddress3, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/3002/republic/8MKZ8JwCU9m9affPWHZ9rxp2azXNnE")
	if err != nil {
		log.Fatal(err)
	}

	for _, order := range orders {
		shares, err := order.Split(miner.N, miner.K, miner.Prime)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("placing order", base58.Encode(order.ID))
		if _, err := network.SendOrderFragment(multiAddress1, shares[0]); err != nil {
			log.Fatal(err)
		}
		if _, err := network.SendOrderFragment(multiAddress2, shares[1]); err != nil {
			log.Fatal(err)
		}
		if _, err := network.SendOrderFragment(multiAddress3, shares[2]); err != nil {
			log.Fatal(err)
		}
	}
}
