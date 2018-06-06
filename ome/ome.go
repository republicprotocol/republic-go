package ome

import (
	"encoding/base64"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
)

type Ome interface {
	cal.EpochListener

	// Run starts running the ome, it syncs with the orderbook to discover new
	// orders, purge confirmed orders and reprioritize order matching computations.
	Run(done <-chan struct{}) <-chan error
}

type ome struct {
	ranker    Ranker
	computer  Computer
	orderbook orderbook.Orderbook
	smpcer    smpc.Smpcer

	ξMu *sync.RWMutex
	ξ   cal.Epoch
}

func NewOme(ranker Ranker, computer Computer, orderbook orderbook.Orderbook, smpcer smpc.Smpcer) Ome {
	return &ome{
		ranker:    ranker,
		computer:  computer,
		orderbook: orderbook,
		smpcer:    smpcer,

		ξMu: new(sync.RWMutex),
		ξ:   cal.Epoch{},
	}
}

// OnChangeEpoch implements the cal.EpochListener interface.
func (ome *ome) OnChangeEpoch(ξ cal.Epoch) {
	ome.ξMu.Lock()
	defer ome.ξMu.Unlock()

	ome.smpcer.Disconnect(ome.ξ.Hash)

	ome.ξ = ξ
	log.Printf("connecting to peers:\n  %v", ome.ξ.Darknodes)

	ome.smpcer.Connect(ome.ξ.Hash, ome.ξ.Darknodes, int64(2*(len(ome.ξ.Darknodes)+1)/3))
}

func (ome *ome) Run(done <-chan struct{}) <-chan error {
	computations := make(chan ComputationEpoch)
	errs := make(chan error, 3)
	wg := new(sync.WaitGroup)

	// Sync with the orderbook
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			syncStart := time.Now()

			select {
			case <-done:
				return
			default:
			}

			log.Println("orderbook sync")
			changeset, err := ome.orderbook.Sync()
			if err != nil {
				errs <- fmt.Errorf("cannot sync orderbook: %v", err)
				continue
			}
			if err := ome.syncRanker(changeset); err != nil {
				errs <- fmt.Errorf("cannot sync ranker: %v", err)
				continue
			}

			if time.Now().After(syncStart.Add(14 * time.Second)) {
				continue
			}
			time.Sleep(syncStart.Add(14 * time.Second).Sub(time.Now()))
		}
	}()

	// Sync with the ranker
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(computations)

		for {
			syncStart := time.Now()

			select {
			case <-done:
				return
			default:
			}

			buffer := [128]Computation{}
			n := ome.ranker.Computations(buffer[:])
			for i := 0; i < n; i++ {
				ome.ξMu.RLock()
				computation := ComputationEpoch{
					Computation: buffer[i],
					Epoch:       ome.ξ.Hash,
				}
				ome.ξMu.RUnlock()
				select {
				case <-done:
					return
				case computations <- computation:
					log.Printf("new computation: buy = %v; sell = %v", base64.StdEncoding.EncodeToString(computation.Buy[:8]), base64.StdEncoding.EncodeToString(computation.Sell[:8]))
				}
			}
			if n == 128 {
				continue
			}

			if time.Now().After(syncStart.Add(14 * time.Second)) {
				continue
			}
			time.Sleep(syncStart.Add(14 * time.Second).Sub(time.Now()))
		}
	}()

	// Sync with the computer
	wg.Add(1)
	go func() {
		defer wg.Done()

		computationErrs := ome.computer.Compute(done, computations)
		for err := range computationErrs {
			select {
			case <-done:
				return
			case errs <- err:
			}
		}
	}()

	go func() {
		defer close(errs)

		wg.Wait()
	}()

	return errs
}

func (ome *ome) syncRanker(changeset orderbook.ChangeSet) error {
	for _, change := range changeset {
		switch change.OrderStatus {
		case order.Open:
			if change.OrderParity == order.ParityBuy {
				ome.ranker.InsertBuy(PriorityOrder{
					ID:       change.OrderID,
					Priority: Priority(change.OrderPriority),
				})
			} else {
				ome.ranker.InsertSell(PriorityOrder{
					ID:       change.OrderID,
					Priority: Priority(change.OrderPriority),
				})
			}
		case order.Canceled, order.Confirmed:
			ome.ranker.Remove(change.OrderID)
		}
	}
	return nil
}
