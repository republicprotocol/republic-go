package main

import (
	"flag"
	"log"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/relay"
)

func main() {
	paramType := flag.String("type", "limit", "")
	paramParity := flag.String("parity", "buy", "")
	// paramExpiry := flag.String("expiry", "1h", "")
	// paramTokens := flag.String("tokens", "RENETH", "")
	paramPrice := flag.Int("price", 0, "")
	paramVolume := flag.Int("volume", 0, "")
	paramMinimumVolume := flag.Int("minimumVolume", 0, "")

	flag.Parse()

	ty := order.TypeLimit
	switch *paramType {
	case "limit":
		ty = order.TypeLimit
	case "midpoint":
		ty = order.TypeIBBO
	}
	parity := order.ParityBuy
	switch *paramParity {
	case "buy":
		parity = order.ParityBuy
	case "sell":
		parity = order.ParityBuy
	}
	expiry := time.Now().Add(time.Hour * 24)
	// TODO: Parse `paramExpiry` and use it.

	fstCode := order.CurrencyCodeREN
	sndCode := order.CurrencyCodeETH

	url := "http://localhost:18515"
	api := relay.NewAPI(url, relay.Insecure())
	order := relay.Order{
		Type:      ty,
		Parity:    parity,
		Expiry:    expiry,
		FstCode:   fstCode,
		SndCode:   sndCode,
		Price:     uint(*paramPrice),
		MaxVolume: uint(*paramVolume),
		MinVolume: uint(*paramMinimumVolume),
	}
	orderID, err := api.OpenOrder(order)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("opened:", orderID)

}
