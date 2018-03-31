package smpc_test

import (
	"context"
	"log"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	. "github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

func TestWorker(t *testing.T) {

	numOrders := 10
	n := 24
	k := 16
	prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	if err != nil {
		t.Fatal(err)
	}

	log.Println("configuring a secure multiparty computation")

	contexts := make([]context.Context, n)
	cancelFuncs := make([]context.CancelFunc, n)
	orderFragmentReceivers := make([]chan order.Fragment, n)
	deltaFragmentReceivers := make([]chan DeltaFragment, n)

	for i := 0; i < n; i++ {
		contexts[i], cancelFuncs[i] = context.WithCancel(context.Background())
		orderFragmentReceivers[i] = make(chan order.Fragment, 2*n)
		deltaFragmentReceivers[i] = make(chan DeltaFragment, 2*n)
	}

	var nodesWg sync.WaitGroup
	results := make(chan Delta, n*numOrders*numOrders)
	for i := 0; i < n; i++ {
		nodesWg.Add(1)
		go func(i int) {
			defer nodesWg.Done()

			log.Printf("[%d] booting", i)

			var nodeWg sync.WaitGroup
			deltaFragmentMatrix := smpc.NewDeltaFragmentMatrix(prime)
			deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				if err := OrderFragmentReceiver(contexts[i], orderFragmentReceivers[i], deltaFragmentMatrix); err != nil && err != context.Canceled {
					t.Fatal(err)
				}
			}()

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				if err := DeltaFragmentReceiver(contexts[i], deltaFragmentReceivers[i], deltaBuilder); err != nil && err != context.Canceled {
					t.Fatal(err)
				}
			}()

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				deltaFragments, errors := DeltaFragmentBroadcaster(contexts[i], deltaFragmentMatrix, 2*n)
				for {
					select {
					case deltaFragment := <-deltaFragments:
						for j := 0; j < n; j++ {
							deltaFragmentReceivers[j] <- deltaFragment
						}
					case err := <-errors:
						if err != nil && err != context.Canceled {
							t.Fatal(err)
						}
						return
					}
				}
			}()

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				deltas, errors := DeltaBroadcaster(contexts[i], deltaBuilder, 2*n)
				for {
					select {
					case delta := <-deltas:
						log.Printf("[%d] buy = %s, sell = %s", i, delta.BuyOrderID.String(), delta.SellOrderID.String())
						results <- delta
					case err := <-errors:
						if err != nil && err != context.Canceled {
							t.Fatal(err)
						}
						return
					}
				}
			}()

			nodeWg.Wait()
			log.Printf("[%d] shutting down", i)
		}(i)
	}

	log.Println("sending order fragments")
	for i := 0; i < numOrders; i++ {
		go func(i int) {
			buyOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range orderFragmentReceivers {
				orderFragmentReceivers[j] <- *buyOrderFragments[j]
			}
		}(i)
		go func(i int) {
			sellOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range orderFragmentReceivers {
				orderFragmentReceivers[j] <- *sellOrderFragments[j]
			}
		}(i)
	}

	log.Println("waiting for deltas")
	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Printf("[main] %d active goroutines", runtime.NumGoroutine())
		}
	}()

	for i := 0; i < n*numOrders*numOrders; i++ {
		delta := <-results
		if !delta.IsMatch(prime) {
			t.Fatalf("[%s] unexpected result", delta.ID.String())
		}
	}

	log.Println("shutting down")
	for i := range cancelFuncs {
		cancelFuncs[i]()
	}
	nodesWg.Wait()
}
