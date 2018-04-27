package main

import (
	"context"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin/arc"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/swap"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/stackint"
)

func initSwap(ctx context.Context, conn ethereum.Conn, auth *bind.TransactOpts, entry orderbook.Entry, fstOrderID, sndOrderID [32]byte) swap.Swap {
	fstTokenAddress := getTokenAddress(entry.FstCode)
	sndTokenAddress := getTokenAddress(entry.SndCode)
	from := common.HexToAddress("").Bytes()
	to := common.HexToAddress("").Bytes()
	fstValue := getFstValue(entry.Price, entry.MaxVolume) // Assuming max and min volume to be same
	sndValue := getSndValue(entry.Price, entry.MaxVolume) // Assuming max and min volume to be same
	order := entry.Order.Bytes()
	fee := 2
	fstItem := swap.SwapItem{
		orderID:   fstOrderID,
		from:      from,
		to:        to,
		value:     fstValue,
		expiry:    time.Now().Unix() + 48*60*60,
		arc:       arc.NewArc(ctx, conn, auth, fstOrder, fstTokenAddress, fee),
		goesFirst: entry.FstCode > entry.SndCode,
	}
	sndItem := swap.SwapItem{
		orderID:   sndOrderID,
		from:      to,
		to:        from,
		value:     sndValue,
		expiry:    time.Now().Unix() + 24*60*60,
		arc:       arc.NewArc(ctx, conn, auth, []byte{}, sndTokenAddress, fee),
		goesFirst: entry.FstCode < entry.SndCode,
	}
	return swap.NewSwap(fstItem, sndItem)
}

func processAtomicSwaps(swaps <-chan swap.Swap) {
	go func() {
		for {
			select {
			case swap, ok := <-swaps:
				if !ok {
					return
				}
				go func() {
					err := swap.Execute(context.Background())
					if err != nil {
						log.Fatalf("failed to execute the atomic swap: %v", err)
						return
					}
				}()
			}
		}
	}()
}

func getTokenAddress(currencyCode order.CurrencyCode) common.Address {
	switch currencyCode {
	case order.CurrencyCodeREN:
		return common.HexToAddress("")
	case order.CurrencyCodeETH:
		return common.HexToAddress("")
	default:
		panic("unimplemented")
	}
}

func getFstValue(price, volume stackint.Int1024) *big.Int {
	value := volume.Mul(&price)
	return value.ToBigInt()
}

func getSndValue(price, volume stackint.Int1024) *big.Int {
	one := stackint.One()
	price = one.Div(&price)
	value := volume.Mul(&price)
	return value.ToBigInt()
}
