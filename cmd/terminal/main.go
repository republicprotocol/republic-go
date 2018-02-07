package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jbenet/go-base58"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-miner"
	"github.com/republicprotocol/go-order-compute"
	rpc "github.com/republicprotocol/go-rpc"
)

var orders = []*compute.Order{
	// Buy ETH for REN
	compute.NewOrder(
		compute.OrderTypeLimit,
		compute.OrderParityBuy,
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
		compute.OrderParitySell,
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
		compute.OrderParityBuy,
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
		compute.OrderParitySell,
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
		compute.OrderParityBuy,
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
		compute.OrderParitySell,
		compute.CurrencyCodeREN,
		compute.CurrencyCodeBTC,
		172000,
		1000,
		1000,
		0,
	),
}

func main() {

	traderMultiAddress, err := identity.NewMultiAddressFromString("/ip4/127.0.0.1/tcp/3000/republic/8MGg76n7RfC6tuw23PYf85VFyAbCd")
	if err != nil {
		log.Fatal(err)
	}

	minerMultiAddresses := make(identity.MultiAddresses, miner.N)
	for i := range minerMultiAddresses {
		multiAddress, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/8MGg76n7RfC6tuw23PYf85VFyM8Zto", 4000+i))
		if err != nil {
			log.Fatal(err)
		}
		minerMultiAddresses[i] = multiAddress
	}

	for _, order := range orders {
		shares, err := order.Split(miner.N, miner.K, miner.Prime)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("placing order", base58.Encode(order.ID))

		for i := range shares {
			if err := rpc.SendOrderFragmentToTarget(minerMultiAddresses[i], minerMultiAddresses[i].Address(), traderMultiAddress, shares[i], time.Minute); err != nil {
				log.Fatal(err)
			}
		}
	}
}
