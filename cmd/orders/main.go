package main

import (
	"log"

	base58 "github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-identity"

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

	log.Println("placing buy order...")

	if _, err := network.SendOrderFragment(multiAddress1, buyOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(buyOrderFragments[0].ID))

	if _, err := network.SendOrderFragment(multiAddress2, buyOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(buyOrderFragments[1].ID))

	if _, err := network.SendOrderFragment(multiAddress3, buyOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(buyOrderFragments[2].ID))

	log.Println("placing sell order...")

	if _, err := network.SendOrderFragment(multiAddress1, sellOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(sellOrderFragments[0].ID))

	if _, err := network.SendOrderFragment(multiAddress2, sellOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(sellOrderFragments[1].ID))

	if _, err := network.SendOrderFragment(multiAddress3, sellOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("sent", base58.Encode(sellOrderFragments[2].ID))
}
