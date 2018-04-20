package darknode

import (
	"bytes"
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/dispatch"
	"google.golang.org/grpc"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
)

type Router struct {
	maxConnections int
	address        identity.Address
	multiAddress   identity.MultiAddress

	orderFragmentSplitterCh chan order.Fragment
	orderFragmentSplitter   *dispatch.Splitter

	mu               *sync.Mutex
	computeSenders   map[identity.Address]chan *rpc.Computation
	computeReceivers map[identity.Address]*dispatch.Splitter
	computeErrs      map[identity.Address]*dispatch.Splitter
	computeArcs      map[identity.Address]int

	dht        *dht.DHT
	clientPool *rpc.ClientPool
	swarmer    rpc.SwarmService
	relayer    rpc.RelayService
	syncer     rpc.SyncerService
	smpcer     rpc.ComputerService
}

func NewRouter(maxConnections int, multiAddress identity.MultiAddress, options rpc.Options, key identity.KeyPair, privateKey *rsa.PrivateKey, orderbook *orderbook.Orderbook) *Router {
	router := &Router{
		maxConnections: maxConnections,
		address:        multiAddress.Address(),
		multiAddress:   multiAddress,

		orderFragmentSplitterCh: make(chan order.Fragment, 100),
		orderFragmentSplitter:   dispatch.NewSplitter(maxConnections),

		mu:               new(sync.Mutex),
		computeSenders:   map[identity.Address]chan *rpc.Computation{},
		computeReceivers: map[identity.Address]*dispatch.Splitter{},
		computeErrs:      map[identity.Address]*dispatch.Splitter{},
		computeArcs:      map[identity.Address]int{},
	}
	router.dht = dht.NewDHT(multiAddress.Address(), 100)
	router.clientPool = rpc.NewClientPool(multiAddress, key)
	router.swarmer = rpc.NewSwarmService(options, router.clientPool, router.dht, logger.StdoutLogger)
	router.relayer = rpc.NewRelayService(options, router, privateKey, logger.StdoutLogger)
	router.smpcer = rpc.NewComputerService()
	router.syncer = rpc.NewSyncerService(options, logger.StdoutLogger, orderbook)
	return router
}

func (router *Router) Run(done <-chan struct{}, host, port string) <-chan error {
	errs := make(chan error, 1)

	go func() {
		defer close(errs)

		server := grpc.NewServer()
		router.swarmer.Register(server)
		router.relayer.Register(server)
		router.syncer.Register(server)
		router.smpcer.Register(server)
		listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", host, port))
		if err != nil {
			errs <- err
			return
		}

		dispatch.CoBegin(func() {
			if err := server.Serve(listener); err != nil {
				errs <- err
				return
			}
		}, func() {
			router.swarmer.Bootstrap()
		}, func() {
			router.orderFragmentSplitter.Split(router.orderFragmentSplitterCh)
		})
	}()

	return errs
}

func (router *Router) OrderFragments(done <-chan struct{}) (<-chan order.Fragment, <-chan error) {
	orderFragmentReceiver := make(chan order.Fragment)
	errs := make(chan error, 1)

	go func() {
		defer close(orderFragmentReceiver)
		defer close(errs)

		if err := router.orderFragmentSplitter.Subscribe(orderFragmentReceiver); err != nil {
			errs <- err
			return
		}

		<-done
		router.orderFragmentSplitter.Unsubscribe(orderFragmentReceiver)
	}()

	return orderFragmentReceiver, errs
}

func (router *Router) Compute(done <-chan struct{}, addr identity.Address, computationSender <-chan *rpc.Computation) (<-chan *rpc.Computation, <-chan error) {
	computationReceiver := make(chan *rpc.Computation)
	errs := make(chan error, 1)

	go func() {
		defer close(computationReceiver)
		defer close(errs)

		var err error
		router.mu.Lock()
		if _, ok := router.computeSenders[addr]; !ok {
			err = router.setupCompute(addr)
		}
		router.computeArcs[addr]++
		router.mu.Unlock()
		if err != nil {
			log.Printf("err: %v", err)
			errs <- err
			return
		}

		router.computeReceivers[addr].Subscribe(computationReceiver)
		router.computeErrs[addr].Subscribe(errs)

		defer func() {

			router.computeReceivers[addr].Unsubscribe(computationReceiver)
			router.computeErrs[addr].Unsubscribe(errs)

			router.mu.Lock()
			if router.computeArcs[addr]--; router.computeArcs[addr] == 0 {
				router.teardownCompute(addr)
			}
			router.mu.Unlock()
		}()

		for {
			select {
			case <-done:
				return
			case computation, ok := <-computationSender:
				if !ok {
					return
				}
				select {
				case <-done:
				case router.computeSenders[addr] <- computation:
				}
			}
		}
	}()

	return computationReceiver, errs
}

func (router *Router) setupCompute(addr identity.Address) error {

	var receiver <-chan *rpc.Computation
	var errs <-chan error

	sender := make(chan *rpc.Computation)
	router.computeSenders[addr] = sender

	if bytes.Compare(router.address.ID()[:], addr.ID()[:]) < 0 {

		multiAddr, err := router.swarmer.FindNode(addr.ID())
		if err != nil {
			return err
		}
		if multiAddr == nil {
			return fmt.Errorf("multiaddress not found: %v", addr)
		}

		router.mu.Unlock()
		receiver, errs = router.clientPool.Compute(context.Background(), *multiAddr, sender)
		sender <- &rpc.Computation{MultiAddress: rpc.MarshalMultiAddress(&router.multiAddress)}
		println("client rendezvous connected!")
		router.mu.Lock()

	} else {
		router.mu.Unlock()
		receiver, errs = router.smpcer.WaitForCompute(addr, sender)
		println("service rendezvous connected!")
		router.mu.Lock()
	}

	router.computeReceivers[addr] = dispatch.NewSplitter(router.maxConnections)
	go router.computeReceivers[addr].Split(receiver)
	router.computeErrs[addr] = dispatch.NewSplitter(router.maxConnections)
	go router.computeErrs[addr].Split(errs)
	return nil
}

func (router *Router) teardownCompute(addr identity.Address) {
	close(router.computeSenders[addr])
	delete(router.computeSenders, addr)
	delete(router.computeReceivers, addr)
	delete(router.computeErrs, addr)
}

func (router *Router) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
	router.orderFragmentSplitterCh <- *orderFragment
}
