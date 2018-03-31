package smpc_test

import (
	"context"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/smpc"
	. "github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/stackint"
)

func TestWorker(t *testing.T) {

	numOrders := 100
	n := 48
	k := 32
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

	results := make(chan Delta, n*numOrders*numOrders)
	for i := 0; i < n; i++ {
		go func(i int) {
			deltaFragmentMatrix := smpc.NewDeltaFragmentMatrix(prime)
			deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)

			go func() {
				if err := OrderFragmentReceiver(contexts[i], orderFragmentReceivers[i], deltaFragmentMatrix); err != nil {
					t.Fatal(err)
				}
			}()
			go func() {
				if err := DeltaFragmentReceiver(contexts[i], deltaFragmentReceivers[i], deltaBuilder); err != nil {
					t.Fatal(err)
				}
			}()
			go func() {
				deltaFragments, errors := DeltaFragmentWaiter(contexts[i], deltaFragmentMatrix, deltaBuilder)
				for {
					select {
					case deltaFragment := <-deltaFragments:
						for j := 0; j < n; j++ {
							if i == j {
								continue
							}
							deltaFragmentReceivers[j] <- deltaFragment
						}
					case <-errors:
						log.Printf("[%d] delta fragment waiter shutting down", i)
						return
					}
				}
			}()
			go func() {
				deltas, errors := DeltaBroadcaster(contexts[i], deltaBuilder)
				for {
					select {
					case delta := <-deltas:
						log.Printf("[%d] %s", i, delta.ID.String())
						results <- delta
					case <-errors:
						log.Printf("[%d] delta waiter shutting down", i)
						return
					}
				}
			}()
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
}
