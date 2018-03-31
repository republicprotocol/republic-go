package smpc_test

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	. "github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

func TestWorker(t *testing.T) {
	runtime.GOMAXPROCS(1)
	queueLimit := 0
	numOrders := 10
	numWorkers := 1
	numBroadcasters := 1
	n := 72
	k := 48
	prime, err := stackint.FromString("10007")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("configuring a secure multiparty computation")

	deltaQueue := NewDeltaQueue(queueLimit)
	messageQueues := make(dispatch.MessageQueues, n)
	for i := 0; i < n; i++ {
		messageQueues[i] = dispatch.NewChannelQueue(queueLimit)
	}
	multiplexers := make([]dispatch.Multiplexer, n)
	for i := 0; i < n; i++ {
		multiplexers[i] = dispatch.NewMultiplexer(queueLimit)
	}
	workers := make(Workers, n*numWorkers)
	broadcasters := make(Broadcasters, n*numBroadcasters)
	for i := 0; i < n; i++ {
		// go func(i int) {
		deltaFragmentMatrix := smpc.NewDeltaFragmentMatrix(prime)
		deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)

		// Run the MessageQueue that will feed data into the Multiplexer
		err := multiplexers[i].RunMessageQueue(fmt.Sprintf("%d", i), messageQueues[i])
		if err != nil {
			t.Fatal(err)
		}

		// Create a Worker that is connected to all other parties
		workerPeerQueues := make(dispatch.MessageQueues, 0, n-1)
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			workerPeerQueues = append(workerPeerQueues, messageQueues[j])
		}

		for j := 0; j < numWorkers; j++ {
			workers[i*numWorkers+j] = NewWorker(i == 0, nil, workerPeerQueues, &multiplexers[i], deltaFragmentMatrix, deltaBuilder, &deltaQueue)
			go func(i, j int) {
				workers[i*numWorkers+j].Run()
			}(i, j)
		}
		for j := 0; j < numBroadcasters; j++ {
			broadcasters[i*numBroadcasters+j] = NewBroadcaster(i == 0, nil, workerPeerQueues, deltaFragmentMatrix, deltaBuilder, &deltaQueue)
			go func(i, j int) {
				broadcasters[i*numBroadcasters+j].Run()
			}(i, j)
		}
		// }(i)
	}

	t.Log("sending order fragments")
	for i := 0; i < numOrders; i++ {
		go func(i int) {
			buyOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range multiplexers {
				multiplexers[j].Send(Message{
					OrderFragment: buyOrderFragments[j],
				})
			}
		}(i)
		go func(i int) {
			sellOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range multiplexers {
				multiplexers[j].Send(Message{
					OrderFragment: sellOrderFragments[j],
				})
			}
		}(i)
	}

	t.Log("waiting for deltas")
	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Printf("%d living goroutines", runtime.NumGoroutine())
		}
	}()

	for i := 0; i < n*numOrders*numOrders; i++ {
		delta, ok := deltaQueue.Recv()
		if !ok {
			t.Fatal("not ok!")
		}
		switch delta := delta.(type) {
		case Delta:
			if i%(n*numOrders) == 0 {
				log.Println(delta.ID.String())
			}
			if !delta.IsMatch(prime) {
				t.Fatal("no match!")
			}
		default:
			t.Fatal("unrecognized type %T: expected smpc.Delta", delta)
		}
	}

	t.Log("shutting down")
	for i := range multiplexers {
		if err := multiplexers[i].Shutdown(); err != nil {
			t.Fatal(err)
		}
	}
	for i := range workers {
		workers[i].Shutdown()
	}
	for i := range broadcasters {
		broadcasters[i].Shutdown()
	}
}
