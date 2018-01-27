package main

import (
	"log"
	"time"

	"github.com/jbenet/go-base58"
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
	buyOrderFragments, err := buyOrder.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}

	badBuyOrder := compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderBuy,
		compute.CurrencyCodeBTC,
		compute.CurrencyCodeREN,
		17200,
		1000,
		100,
		0,
	)
	badBuyOrderFragments, err := badBuyOrder.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}

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
	sellOrderFragments, err := sellOrder.Split(miner.N, miner.K, miner.Prime)
	if err != nil {
		log.Fatal(err)
	}

	badSellOrder := compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderSell,
		compute.CurrencyCodeBTC,
		compute.CurrencyCodeETH,
		17200,
		1000,
		100,
		0,
	)
	badSellOrderFragments, err := badSellOrder.Split(miner.N, miner.K, miner.Prime)
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

	log.Println("placing bad buy order", base58.Encode(badBuyOrder.ID))
	if _, err := network.SendOrderFragment(multiAddress1, badBuyOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badBuyOrderFragments[0].ID))
	if _, err := network.SendOrderFragment(multiAddress2, badBuyOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badBuyOrderFragments[1].ID))
	if _, err := network.SendOrderFragment(multiAddress3, badBuyOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badBuyOrderFragments[2].ID))

	log.Println("placing bad sell order", base58.Encode(badSellOrder.ID))
	if _, err := network.SendOrderFragment(multiAddress1, badSellOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badSellOrderFragments[0].ID))
	if _, err := network.SendOrderFragment(multiAddress2, badSellOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badSellOrderFragments[1].ID))
	if _, err := network.SendOrderFragment(multiAddress3, badSellOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(badSellOrderFragments[2].ID))

	time.Sleep(4 * time.Second)

	log.Println("placing good buy order", base58.Encode(buyOrder.ID))
	if _, err := network.SendOrderFragment(multiAddress1, buyOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(buyOrderFragments[0].ID))
	if _, err := network.SendOrderFragment(multiAddress2, buyOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(buyOrderFragments[1].ID))
	if _, err := network.SendOrderFragment(multiAddress3, buyOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(buyOrderFragments[2].ID))

	log.Println("placing good sell order", base58.Encode(sellOrder.ID))
	if _, err := network.SendOrderFragment(multiAddress1, sellOrderFragments[0]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(sellOrderFragments[0].ID))
	if _, err := network.SendOrderFragment(multiAddress2, sellOrderFragments[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(sellOrderFragments[1].ID))
	if _, err := network.SendOrderFragment(multiAddress3, sellOrderFragments[2]); err != nil {
		log.Fatal(err)
	}
	log.Println("  sent", base58.Encode(sellOrderFragments[2].ID))
}
