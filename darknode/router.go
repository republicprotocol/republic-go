package darknode

import (
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/rpc"
)

type Router struct {
}

func (router *Router) OrderFragments(epoch [32]byte) <-chan order.Fragment {
	orderFragmentReceiver := make(chan *rpc.Computation)

	go func() {
		defer close(orderFragmentReceiver)

		// TODO: Connect to a stream of order.Fragments
	}()

	return orderFragmentReceiver
}

func (router *Router) Compute(epoch [32]byte, addr identity.Address, computationSender <-chan *rpc.Computation) (<-chan *rpc.Computation, <-chan error) {
	computationReceiver := make(chan *rpc.Computation)
	errs := make(chan error)

	go func() {
		defer close(computationsOut)
		defer close(errs)

		// TODO: Connect to a bidirectional stream of rpc.Computations
	}()

	return computationReceiver, errs
}
