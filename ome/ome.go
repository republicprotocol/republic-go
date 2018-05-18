package ome

import (
	"errors"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/smpc/delta"
)

var ErrNotFoundInStore = errors.New("not found in store")

type Omer interface {
	cal.EpochListener

	// Sync the Omer with the orderbook.Orderbooker so that it can discover new
	// orders, purge confirmed orders, and reprioritize order matching
	// computations.
	Sync() error
}

type ome struct {
	ranker Ranker
	deltas map[delta.ID][]delta.Fragment

	ξs                chan cal.Epoch
	orderbooker       orderbook.Orderbooker
	orderbookListener orderbook.Listener
	smpcer            smpc.Smpcer
}

func NewOme(orderbooker orderbook.Orderbooker, orderbookListener orderbook.Listener, smpcer smpc.Smpcer) Omer {
	return &ome{
		ranker: NewPriorityQueue(),
		deltas: map[delta.ID][]delta.Fragment{},

		ξs:                make(chan cal.Epoch, 1),
		orderbooker:       orderbooker,
		orderbookListener: orderbookListener,
		smpcer:            smpcer,
	}
}

// Sync implements the Omer interface.
func (ome *ome) Sync() error {
	updates, err := ome.orderbooker.Sync()
	if err != nil {
		return err
	}
	for _, update := range updates {
		switch update.Status {
		case order.Open:
			ome.ranker.Insert(update.ID, update.Parity, update.Priority)
		case order.Canceled, order.Settled, order.Confirmed:
			ome.ranker.Remove(update.ID)
		}
	}
	return nil
}

func (ome *ome) Compute(done <-chan struct{}) <-chan error {
	errs := make(chan error, 1)
	var currentEpoch cal.Epoch

	if err := ome.smpcer.Start(); err != nil {
		errs <- err
		close(errs)
		return errs
	}

	orderPairs := ome.ranker.OrderPairs(done)
	input, output := ome.smpcer.Instructions(), ome.smpcer.Results()

	go func() {
		defer close(errs)
		defer ome.smpcer.Shutdown()

		orderPairWaitingList := []OrderPair{}
		deltaWaitingList := []delta.Delta{}

		for {
			select {
			case <-done:
				return

			case ξ, ok := <-ome.ξs:
				if !ok {
					return
				}
				currentEpoch = ξ
				input <- smpc.Inst{
					InstConnect: &smpc.InstConnect{
						PeersID: ξ.Hash[:],
						Peers:   ξ.Darknodes,
						N:       len(ξ.Darknodes),
						K:       (len(ξ.Darknodes) + 1) * 2 / 3,
					},
				}

			case pair, ok := <-orderPairs:
				if !ok {
					return
				}
				for i := 0; i < len(orderPairWaitingList); i++ {
					//todo : need to have a way for element to expire
					inst, err := generateInstruction(orderPairWaitingList[i], ome.orderbooker, currentEpoch)
					if err != nil {
						select {
						case <-done:
							return
						case errs <- err:
						}
						continue
					}

					input <- inst
					orderPairWaitingList = append(orderPairWaitingList[:i], orderPairWaitingList[i+1:]...)
					i--
				}
				inst, err := generateInstruction(pair, ome.orderbooker, currentEpoch)
				if err != nil {
					orderPairWaitingList = append(orderPairWaitingList, pair)
					select {
					case <-done:
						return
					case errs <- err:
					}
					continue
				}
				select {
				case <-done:
					return
				case input <- inst:
				}

			case result, ok := <-output:
				if !ok {
					return
				}

				for i := 0; i < len(deltaWaitingList); i++ {
					buyOrd, sellOrd, err := getOrderDetails(deltaWaitingList[i], ome.orderbooker)
					if err != nil {
						continue
					}
					if err := ome.orderbooker.ConfirmOrderMatch(result.Delta.BuyOrderID, result.Delta.SellOrderID); err != nil {
						continue
					}
					// FIXME: Reconstruct orders after the confirmation has been accepted
					ome.orderbookListener.OnConfirmOrderMatch(buyOrd, sellOrd)
					deltaWaitingList = append(deltaWaitingList[:i], deltaWaitingList[i+1:]...)
					i--
				}

				if result.Delta.IsMatch() {
					buyOrd, sellOrd, err := getOrderDetails(result.Delta, ome.orderbook)
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
					err = ome.ledger.ConfirmOrder(result.Delta.BuyOrderID, []order.ID{result.Delta.SellOrderID})
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
					err = ome.delegate.OnConfirmOrder(buyOrd, []order.Order{sellOrd})
					if err != nil {
						deltaWaitingList = append(deltaWaitingList, result.Delta)
						continue
					}
				}
			}
		}
	}()

	return errs
}

// OnChangeEpoch implements the cal.EpochListener interface.
func (ome *ome) OnChangeEpoch(ξ cal.Epoch) {
	ome.ξs <- ξ
}

func generateInstruction(pair OrderPair, orderbook orderbook.Orderbooker, epoch cal.Epoch) (smpc.Inst, error) {
	// Get orderFragment from the orderbook
	buyOrderFragment, err := orderbook.OrderFragment(pair.BuyOrder)
	if err != nil {
		return smpc.Inst{}, err
	}
	sellOrderFragment, err := orderbook.OrderFragment(pair.SellOrder)
	if err != nil {
		return smpc.Inst{}, err
	}

	// generate instruction message for smpc
	return smpc.Inst{
		InstCompute: &smpc.InstCompute{
			PeersID: epoch.Hash[:],
			Buy:     buyOrderFragment,
			Sell:    sellOrderFragment,
		},
	}, nil
}

func getOrderDetails(delta delta.Delta, orderbooker orderbook.Orderbooker) (order.Order, order.Order, error) {
	buyOrd, err := orderbooker.Order(delta.BuyOrderID)
	if err != nil {
		return order.Order{}, order.Order{}, err
	}
	sellOrd, err := orderbooker.Order(delta.SellOrderID)
	if err != nil {
		return order.Order{}, order.Order{}, err
	}

	return buyOrd, sellOrd, nil
}