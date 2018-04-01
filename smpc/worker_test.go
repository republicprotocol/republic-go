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

	numOrders := 100
	n := 72
	k := 48
	prime, err := stackint.FromString("179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111")
	if err != nil {
		t.Fatal(err)
	}

	log.Println("configuring a secure multiparty computation")

	contexts := make([]context.Context, n)
	cancelFuncs := make([]context.CancelFunc, n)
	orderFragmentReceivers := make([]chan order.Fragment, n)
	deltaFragmentReceivers := make([]chan DeltaFragment, n)
	deltaGarbageCollectors := make([]chan Delta, n)
	deltaReceivers := make([]chan Delta, n)

	for i := 0; i < n; i++ {
		contexts[i], cancelFuncs[i] = context.WithCancel(context.Background())
		orderFragmentReceivers[i] = make(chan order.Fragment, 1)
		deltaFragmentReceivers[i] = make(chan DeltaFragment, 1)
		deltaGarbageCollectors[i] = make(chan Delta, numOrders)
		deltaReceivers[i] = make(chan Delta, numOrders)
	}

	var nodesWg sync.WaitGroup
	for i := 0; i < n; i++ {
		nodesWg.Add(1)
		go func(i int) {
			defer nodesWg.Done()

			log.Printf("[%d] booting", i)

			var nodeWg sync.WaitGroup
			computationMatrix := smpc.NewComputationMatrix()
			deltaBuilder := smpc.NewDeltaBuilder(int64(k), prime)

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				if err := OrderFragmentReceiver(contexts[i], orderFragmentReceivers[i], &computationMatrix); err != nil && err != context.Canceled {
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

				if err := DeltaFragmentGarbageCollector(contexts[i], deltaGarbageCollectors[i], &computationMatrix); err != nil && err != context.Canceled {
					t.Fatal(err)
				}
			}()

			nodeWg.Add(1)
			go func() {
				defer nodeWg.Done()

				deltaFragments, errors := DeltaFragmentComputer(contexts[i], &computationMatrix, numOrders, prime)
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
				defer close(deltaReceivers[i])
				defer close(deltaGarbageCollectors[i])

				deltas, errors := DeltaBroadcaster(contexts[i], deltaBuilder, 1)
				for {
					select {
					case delta, ok := <-deltas:
						if !ok {
							return
						}
						deltaGarbageCollectors[i] <- delta
						deltaReceivers[i] <- delta
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
	var orderFragmentWg sync.WaitGroup
	for i := 0; i < numOrders; i++ {
		orderFragmentWg.Add(2)
		go func(i int) {
			defer orderFragmentWg.Done()
			buyOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range orderFragmentReceivers {
				orderFragmentReceivers[j] <- *buyOrderFragments[j]
			}
		}(i)
		go func(i int) {
			defer orderFragmentWg.Done()
			sellOrderFragments, err := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(time.Hour), order.CurrencyCodeBTC, order.CurrencyCodeETH, stackint.FromUint(10), stackint.FromUint(1000), stackint.FromUint(100), stackint.FromUint(uint(i))).Split(int64(n), int64(k), &prime)
			if err != nil {
				t.Fatal(err)
			}
			for j := range orderFragmentReceivers {
				orderFragmentReceivers[j] <- *sellOrderFragments[j]
			}
		}(i)
	}
	go func() {
		orderFragmentWg.Wait()
		for _, ch := range orderFragmentReceivers {
			close(ch)
		}
		log.Println("closed order fragment receivers")
	}()

	log.Println("waiting for deltas")
	go func() {
		for {
			time.Sleep(10 * time.Second)
			log.Printf("[main] %d active goroutines", runtime.NumGoroutine())
		}
	}()

	deltas := make(chan Delta, n*numOrders*numOrders)
	for i, ch := range deltaReceivers {
		go func(i int, ch chan Delta) {
			for delta := range ch {
				if !delta.IsMatch(prime) {
					t.Fatalf("[%s] unexpected result", delta.ID.String())
				}
				log.Printf("[%d] buy = %s, sell = %s", i, delta.BuyOrderID.String(), delta.SellOrderID.String())
				deltas <- delta
			}
		}(i, ch)
	}

	timeout := time.NewTimer(5 * time.Minute)
	done := make(chan struct{})
	go func() {
		ordersMatched := map[string]struct{}{}
		for len(ordersMatched) < numOrders {
			delta := <-deltas
			ordersMatched[string(delta.BuyOrderID)] = struct{}{}
			ordersMatched[string(delta.SellOrderID)] = struct{}{}
		}
		close(done)
	}()

	select {
	case <-timeout.C:
	case <-done:
	}

	log.Println("shutting down")
	for i := range cancelFuncs {
		cancelFuncs[i]()
	}
	nodesWg.Wait()
}
