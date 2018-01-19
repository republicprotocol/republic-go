package main

import (
	"github.com/republicprotocol/go-identity"
	"log"

	"github.com/republicprotocol/go-miner"
	"github.com/republicprotocol/go-network"
	"github.com/republicprotocol/go-order-compute"
)

func main() {
	buyOrder := compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderBuy,
		compute.CurrencyCodeETH,
		compute.CurrencyCodeREN,
		17200,
		1000,
		100,
		0,
	)

	sellOrder := compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderSell,
		compute.CurrencyCodeETH,
		compute.CurrencyCodeREN,
		17200,
		1000,
		100,
		0,
	)
	buyOrderFragments, err := buyOrder.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}
	sellOrderFragments, err := sellOrder.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}

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

	network.SendOrderFragment(multiAddress1, buyOrderFragments[0])
	network.SendOrderFragment(multiAddress2, buyOrderFragments[1])
	network.SendOrderFragment(multiAddress3, buyOrderFragments[2])

	network.SendOrderFragment(multiAddress1, sellOrderFragments[0])
	network.SendOrderFragment(multiAddress2, sellOrderFragments[1])
	network.SendOrderFragment(multiAddress3, sellOrderFragments[2])
}