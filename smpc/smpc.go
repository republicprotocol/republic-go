package smpc

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/republicprotocol/republic-go/identity"
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

// ErrUnmarshalNilBytes is returned when a call to UnmarshalBinary happens on
// an empty list of bytes.
var ErrUnmarshalNilBytes = errors.New("unmarshall nil bytes")

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
	mu             *sync.Mutex
	buffer         int
	swarmer        swarm.Swarmer
	streamer       stream.Streamer
	instructions   chan Inst
	results        chan Result
	router         map[string]*deltaBuilder
	multiAddresses map[identity.Address]identity.MultiAddress

	shutdownMu        *sync.Mutex
	shutdown          chan struct{}
	shutdownDone      chan struct{}
	shutdownInitiated bool
}

func NewSmpc(buffer int, swarmer swarm.Swarmer, streamer stream.Streamer) Smpcer {
	return &smpc{
		mu:             new(sync.Mutex),
		buffer:         buffer,
		swarmer:        swarmer,
		streamer:       streamer,
		instructions:   make(chan Inst, buffer),
		results:        make(chan Result, buffer),
		router:         map[string]*deltaBuilder{},
		multiAddresses: map[identity.Address]identity.MultiAddress{},

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
	peers := [][]identity.Address{}

	for {
		select {
		case <-smpc.shutdown:
			close(smpc.shutdownDone)
			//todo : shall we clean up the stream connections ?
			return
		case inst := <-smpc.instructions:
			if inst.InstConnect != nil {
				peers = append(peers, inst.Peers)
				if len(peers) >= 3 {
					smpc.receivingMessages(inst, peers[0])
				}
				peers = peers[1:]
			}

			if inst.InstCompute != nil {
				deltaFragment := delta.NewDeltaFragment(&inst.InstCompute.Buy, &inst.InstCompute.Sell)
				smpc.sendDeltaFragment(deltaFragment, inst)

				smpc.mu.Lock()
				builder, ok := smpc.router[string(inst.InstCompute.PeersID)]
				if !ok {
					builder = NewDeltaBuilder(inst.InstCompute.PeersID, (len(inst.Peers)+1)*2/3)
					smpc.router[string(inst.InstCompute.PeersID)] = builder
				}
				smpc.mu.Unlock()

				dlt, err := builder.InsertDeltaFragment(deltaFragment)
				if err != nil {
					log.Println("can't insert delta to deltabuilder")
					continue
				}
				if dlt != nil {
					smpc.results <- Result{
						ResultCompute: &ResultCompute{
							Delta: *dlt,
							Err:   nil,
						},
					}
				}
			}
		}
	}
}

// receivingMessages handles received deltaFragments
func (smpc *smpc) receivingMessages(inst Inst, peers []identity.Address) {
	for _, node := range inst.InstConnect.Peers {
		multi, err := smpc.findMultiaddress(node)
		if err != nil {
			log.Printf("can't find node %v", node)
			continue
		}

		stm, err := smpc.streamer.Open(context.Background(), multi)
		go func(stm stream.Stream, multi identity.MultiAddress, inst Inst) {
			var message *DeltaFragmentMessage
			for {
				err := stm.Recv(message)
				if err != nil {
					stm.Close()
					return
				}

				smpc.mu.Lock()
				builder, ok := smpc.router[string(message.PeersID)]
				if !ok {
					builder = NewDeltaBuilder(message.PeersID, (len(inst.Peers)+1)*2/3)
					smpc.router[string(message.PeersID)] = builder
				}
				smpc.mu.Unlock()

				dlt, err := builder.InsertDeltaFragment(message.DeltaFragment)
				if err != nil {
					log.Println("can't insert node to delta builder")
					return
				}
				if dlt != nil {
					smpc.results <- Result{
						ResultCompute: &ResultCompute{
							Delta: *dlt,
							Err:   nil,
						},
					}
				}
			}
		}(stm, multi, inst)
	}

	// Close node which are two epochs ago
	for _, node := range peers {
		multi, err := smpc.findMultiaddress(node)
		if err != nil {
			log.Printf("can't find node %v", node)
			continue
		}
		err = smpc.streamer.Close(multi)
		if err != nil {
			log.Printf("can't close stream with node %v", node)
			continue
		}
	}
}

// sendDeltaFragment to nodes in the same pod.
func (smpc *smpc) sendDeltaFragment(fragment delta.Fragment, inst Inst) {
	for _, node := range inst.Peers {
		multi, err := smpc.findMultiaddress(node)
		if err != nil {
			log.Printf("can't find node %v", node)
			continue
		}

		stream, err := smpc.streamer.Open(context.Background(), multi)
		msg := &DeltaFragmentMessage{
			PeersID:       inst.InstCompute.PeersID,
			DeltaFragment: fragment,
		}
		err = stream.Send(msg)
		if err != nil {
			log.Printf("can't send message to %v", node)
			continue
		}
	}
}

// findMultiaddress returns the multiaddress of the address. It will look up the mapping first.
// If it's not there, it will start a query for the target address.
func (smpc *smpc) findMultiaddress(addr identity.Address) (identity.MultiAddress, error) {
	if multi, ok := smpc.multiAddresses[addr]; ok {
		return multi, nil
	}

	queryCtx, queryCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer queryCancel()
	multi, err := smpc.swarmer.Query(queryCtx, addr, 3)
	if err != nil {
		return identity.MultiAddress{}, err
	}
	smpc.multiAddresses[addr] = multi

	return multi, nil
}

type deltaBuilder struct {
	PeersID []byte
	k       int
	mu      *sync.Mutex

	deltas                 map[delta.ID][]delta.Fragment
	tokenSharesCache       []shamir.Share
	priceCoShareCache      []shamir.Share
	priceExpShareCache     []shamir.Share
	volumeCoShareCache     []shamir.Share
	volumeExpShareCache    []shamir.Share
	minVolumeCoShareCache  []shamir.Share
	minVolumeExpShareCache []shamir.Share
}

func NewDeltaBuilder(peersID []byte, k int) *deltaBuilder {
	builder := new(deltaBuilder)
	builder.PeersID = peersID
	builder.k = k
	builder.mu = new(sync.Mutex)
	builder.deltas = map[delta.ID][]delta.Fragment{}
	builder.tokenSharesCache = make([]shamir.Share, k)
	builder.priceCoShareCache = make([]shamir.Share, k)
	builder.priceExpShareCache = make([]shamir.Share, k)
	builder.volumeCoShareCache = make([]shamir.Share, k)
	builder.volumeExpShareCache = make([]shamir.Share, k)
	builder.minVolumeCoShareCache = make([]shamir.Share, k)
	builder.minVolumeCoShareCache = make([]shamir.Share, k)

	return builder
}

func (builder *deltaBuilder) InsertDeltaFragment(fragment delta.Fragment) (*delta.Delta, error) {
	builder.mu.Lock()
	defer builder.mu.Unlock()

	builder.deltas[fragment.DeltaID] = append(builder.deltas[fragment.DeltaID], fragment)
	if len(builder.deltas[fragment.DeltaID]) > builder.k {
		// join the shares to a delta
		dlt, err := builder.Join(builder.deltas[fragment.DeltaID]...)
		if err != nil {
			return nil, err
		}
		delete(builder.deltas, fragment.DeltaID)

		return dlt, nil
	}

	return nil, nil
}

func (builder *deltaBuilder) Join(deltaFragments ...delta.Fragment) (*delta.Delta, error) {
	if !delta.IsCompatible(deltaFragments) {
		return nil, errors.New("delta fragment are not compatible with each other ")
	}

	for i := 0; i < builder.k; i++ {
		builder.tokenSharesCache[i] = deltaFragments[i].TokenShare
		builder.priceCoShareCache[i] = deltaFragments[i].PriceShare.Co
		builder.priceExpShareCache[i] = deltaFragments[i].PriceShare.Exp
		builder.volumeCoShareCache[i] = deltaFragments[i].VolumeShare.Co
		builder.volumeExpShareCache[i] = deltaFragments[i].VolumeShare.Exp
		builder.minVolumeCoShareCache[i] = deltaFragments[i].MinVolumeShare.Co
		builder.minVolumeExpShareCache[i] = deltaFragments[i].MinVolumeShare.Exp
	}

	return delta.NewDeltaFromShares(
		deltaFragments[0].BuyOrderID,
		deltaFragments[0].SellOrderID,
		builder.tokenSharesCache,
		builder.priceCoShareCache,
		builder.priceExpShareCache,
		builder.volumeCoShareCache,
		builder.volumeExpShareCache,
		builder.minVolumeCoShareCache,
		builder.minVolumeExpShareCache,
	), nil
}
