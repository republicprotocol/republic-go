package smpc

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/cal"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc/delta"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
)

// ErrSmpcerIsAlreadyRunning is returned when a call to Smpcer.Start happens
// on an Smpcer that has already been started.
var ErrSmpcerIsAlreadyRunning = errors.New("smpcer is already running")

// ErrSmpcerIsNotRunning is returned when a call to Smpcer.Shutdown happens on
// an Smpcer that has not been started yet.
var ErrSmpcerIsNotRunning = errors.New("smpcer is not running")

// Smpcer is an interface for a secure multi-party computer. It asynchronously
// consumes computation instructions and produces computation results.
type Smpcer interface {

	// Start the Smpcer. Until a call to Smpcer.Start, no computation
	// instruction will be processed.
	Start() error

	// Shutdown the Smpcer. After a call to Smpcer.Shutdown, no computation
	// instruction will be processed.
	Shutdown() error

	// Instructions channel for sending computation instructions to the Smpcer.
	Instructions() chan<- Inst

	// Results channel for receiving computation results from the Smpcer.
	Results() <-chan Result
}

type smpc struct {
	mu       *sync.Mutex
	buffer       int
	epoch        cal.Epoch
	instructions chan Inst
	results      chan Result
	layers      map[string]deltaBuilder
	swarmer     swarm.Swarmer
	client      stream.Client

	shutdownMu        *sync.Mutex
	shutdown          chan struct{}
	shutdownDone      chan struct{}
	shutdownInitiated bool
}

func NewSmpc(buffer int ,swarmer swarm.Swarmer, client stream.Client) Smpcer {
	return &smpc{
		mu    :new(sync.Mutex),
		buffer:       buffer,
		instructions: make(chan Inst, buffer),
		results:      make(chan Result, buffer),
		layers :      map[string]deltaBuilder{},
		swarmer: swarmer,
		client : client,

		shutdownMu:        new(sync.Mutex),
		shutdown:          nil,
		shutdownDone:      nil,
		shutdownInitiated: true,
	}
}

// Start implements the Smpcer interface.
func (smpc *smpc) Start() error {
	smpc.shutdownMu.Lock()
	defer smpc.shutdownMu.Unlock()

	if smpc.shutdown != nil {
		return ErrSmpcerIsAlreadyRunning
	}
	smpc.shutdown = make(chan struct{})
	smpc.shutdownDone = make(chan struct{})
	smpc.shutdownInitiated = false

	go smpc.run()
	return nil
}

// Shutdown implements the Smpcer interface.
func (smpc *smpc) Shutdown() error {
	smpc.shutdownMu.Lock()
	defer smpc.shutdownMu.Unlock()

	if smpc.shutdownInitiated {
		return ErrSmpcerIsNotRunning
	}
	smpc.shutdownInitiated = true

	close(smpc.shutdown)
	<-smpc.shutdownDone

	smpc.shutdown = nil
	smpc.shutdownDone = nil

	close(smpc.instructions)
	close(smpc.results)
	smpc.instructions = make(chan Inst, smpc.buffer)
	smpc.results = make(chan Result, smpc.buffer)

	return nil
}

// Instructions implements the Smpcer interface.
func (smpc *smpc) Instructions() chan<- Inst {
	return smpc.instructions
}

// Results implements the Smpcer interface.
func (smpc *smpc) Results() <-chan Result {
	return smpc.results
}

func (smpc *smpc) run() {
	for {
		select {
		case <-smpc.shutdown:
			close(smpc.shutdownDone)
			return
		case inst := <- smpc.instructions:
			if inst.InstConnect != nil {
				// todo :epoch happens , update connections

			}

			if inst.InstCompute != nil {
				deltaFragment := LessThan(inst.InstCompute.Buy, inst.InstCompute.Sell)
				// spread the deltafragments
				err := smpc.sendDeltaFragment(deltaFragment)

				builder := smpc.layers[string(inst.InstCompute.PeersID)]
				dlt, err :=builder.InsertDeltaFragment(deltaFragment)
				if err != nil{
					continue
				}
				if dlt != nil {
					smpc.results <- Result{
						ResultCompute : &ResultCompute{
							Delta: *dlt,
							Err:  nil ,
						},
					}
				}
			}
		}
	}
}

func (smpc *smpc) sendDeltaFragment (fragment delta.Fragment) error {
	smpc.mu.Lock()
	defer smpc.mu.Unlock()

	dispatch.CoForAll(smpc.epoch.Darknodes, func(i int) {
		queryCtx, queryCancel  := context.WithTimeout(context.Background(), 5 *time.Second)
		defer queryCancel()

		multi ,err  := smpc.swarmer.Query(queryCtx, smpc.epoch.Darknodes[i], 3 )
		if err != nil{
			log.Printf("can't find node %v",smpc.epoch.Darknodes[i] )
			return
		}

		computeCtx, computeCancel := context.WithTimeout(context.Background(), 5 *time.Second)
		defer computeCancel()
		client := stream.NewClientRecycler(smpc.client)
		stm ,err  := client.Connect(computeCtx, multi )
		if err != nil {
			log.Printf("can't connect to node %v",smpc.epoch.Darknodes[i] )
			return
		}

		//  todo : define the type we want to send
		err = stm.Send(fragment)
		if err != nil {
			log.Printf("can't send fragmetn to node %v",smpc.epoch.Darknodes[i] )
			return
		}
	})

	return nil
}

func LessThan(lhs, rhs order.Fragment) delta.Fragment{
	var buyOrderFragment, sellOrderFragment *order.Fragment
	if lhs.OrderParity == order.ParityBuy {
		buyOrderFragment = &lhs
		sellOrderFragment = &rhs
	} else {
		buyOrderFragment = &rhs
		sellOrderFragment = &lhs
	}
	token := lhs.Tokens.Sub( &rhs.Tokens)
	price := lhs.Price( &rhs.Price)

	panic("unimplemented")


	//
	//fstCodeShare := shamir.Share{
	//	Index:  buyOrderFragment.Tokens.Index,
	//	Value: buyOrderFragment.Tokens.Value.(&sellOrderFragment.FstCodeShare.Value, prime),
	//}
	//sndCodeShare := shamir.Share{
	//	Index:   buyOrderFragment.SndCodeShare.Index,
	//	Value: buyOrderFragment.SndCodeShare.Value.SubModulo(&sellOrderFragment.SndCodeShare.Value, prime),
	//}
	//priceShare := shamir.Share{
	//	Index:   buyOrderFragment.PriceShare.Index,
	//	Value: buyOrderFragment.PriceShare.Value.SubModulo(&sellOrderFragment.PriceShare.Value, prime),
	//}
	//maxVolumeShare := shamir.Share{
	//	Index:   buyOrderFragment.MaxVolumeShare.Index,
	//	Value: buyOrderFragment.MaxVolumeShare.Value.SubModulo(&sellOrderFragment.MinVolumeShare.Value, prime),
	//}
	//minVolumeShare := shamir.Share{
	//	Index:   buyOrderFragment.MinVolumeShare.Index,
	//	Value: sellOrderFragment.MaxVolumeShare.Value.SubModulo(&buyOrderFragment.MinVolumeShare.Value, prime),
	//}
	//
	//return Fragment{
	//	ID:                  FragmentID(crypto.Keccak256([]byte(buyOrderFragment.ID), []byte(sellOrderFragment.ID))),
	//	DeltaID:             ID(crypto.Keccak256([]byte(buyOrderFragment.OrderID), []byte(sellOrderFragment.OrderID))),
	//	BuyOrderID:          buyOrderFragment.OrderID,
	//	SellOrderID:         sellOrderFragment.OrderID,
	//	BuyOrderFragmentID:  buyOrderFragment.ID,
	//	SellOrderFragmentID: sellOrderFragment.ID,
	//	FstCodeShare:        fstCodeShare,
	//	SndCodeShare:        sndCodeShare,
	//	PriceShare:          priceShare,
	//	MaxVolumeShare:      maxVolumeShare,
	//	MinVolumeShare:      minVolumeShare,
	//}
}


type deltaBuilder struct {
	PeersID   []byte
	n   int
	k  int
	mu *sync.Mutex

	deltas             map[delta.ID][]delta.Fragment
	tokenSharesCache   []shamir.Share
	priceCoShareCache  []shamir.Share
	priceExpShareCache []shamir.Share
	volumeCoShareCache  []shamir.Share
	volumeExpShareCache []shamir.Share
	minVolumeCoShareCache []shamir.Share
	minVolumeExpShareCache []shamir.Share

}

func (builder *deltaBuilder) InsertDeltaFragment( fragment delta.Fragment) (*delta.Delta, error) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	builder.deltas[fragment.DeltaID] = append(builder.deltas[fragment.DeltaID] , fragment)
	if len(builder.deltas[fragment.DeltaID]) > builder.k {
		// join the shares to a delta
		dlt, err := builder.Join(builder.deltas[fragment.DeltaID]...)
		if err != nil {
			return nil, err
		}
		delete(builder.deltas,fragment.DeltaID )

		return dlt , nil
	}

	return nil , nil
}


func (builder *deltaBuilder) Join(deltaFragments ...delta.Fragment) (*delta.Delta, error) {
	// Build the Delta
	if len(deltaFragments) >= builder.k{
		if !delta.IsCompatible(deltaFragments) {
			return nil , errors.New("delta fragment are not compatible with each other ")
		}

		for i := 0; i < builder.k; i++ {
			builder.tokenSharesCache[i] = deltaFragments[i].
			builder.sndCodeSharesCache[i] = deltaFragments[i].SndCodeShare
			builder.priceSharesCache[i] = deltaFragments[i].PriceShare
			builder.maxVolumeSharesCache[i] = deltaFragments[i].MaxVolumeShare
			builder.minVolumeSharesCache[i] = deltaFragments[i].MinVolumeShare
		}

		CoExp{
			Co: shamir.Join(coExpShare.Co, ...),
		}

		delta := delta.NewDeltaFromShares(
			deltaFragments[0].BuyOrderID,
			deltaFragments[0].SellOrderID,
			builder.fstCodeSharesCache,
			builder.sndCodeSharesCache,
			builder.priceSharesCache,
			builder.minVolumeSharesCache,
			builder.maxVolumeSharesCache,
			builder.k,
			builder.prime)
		builder.deltas[string(delta.ID)] = delta
		builder.deltasQueue = append(builder.deltasQueue, delta)
	}
}