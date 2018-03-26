package smpc

import "github.com/republicprotocol/republic-go/network/rpc"

// A Worker receives messages from a Dispatcher until the Dispatcher is
// shutdown. It is primarily responsible for decoding the message and
// delegating work to the appropriate component.
type Worker struct{}

// Run until the Dispatcher is shutdown. The worker will read a message from
// the Dispatcher, delegate work to the appropriate component, wait for the
// component to complete, and then read the next message from the Dispatcher.
// This function blocks until the Dispatcher is shutdown.
func (worker *Worker) Run(dispatcher *Dispatcher) {
	for {
		message, ok := dispatcher.Recv()
		if !ok {
			break
		}
		processMessage(message)
	}
}

func processMessage(message *rpc.TauMessage) {
	if message.GenerateRandomShares != nil {
		processGenerateRandomShares(message.GenerateRandomShares)
	}
	if message.GenerateXiShares != nil {
		processGenerateXiShares(message.GenerateXiShares)
	}
	if message.GenerateXiFragment != nil {
		processGenerateXiFragment(message.GenerateXiFragment)
	}
	if message.BroadcastRhoSigmaFragment != nil {
		processBroadcastRhoSigmaFragment(message.BroadcastRhoSigmaFragment)
	}
	if message.BroadcastDeltaFragment != nil {
		processBroadcastDeltaFragment(message.BroadcastDeltaFragment)
	}
}

func processGenerateRandomShares(request *rpc.GenerateRandomShares) {
	panic("unimplemented")
}

func processGenerateXiShares(request *rpc.GenerateXiShares) {
	panic("unimplemented")
}

func processGenerateXiFragment(request *rpc.GenerateXiFragment) {
	panic("unimplemented")
}

func processBroadcastRhoSigmaFragment(request *rpc.BroadcastRhoSigmaFragment) {
	panic("unimplemented")
}

func processBroadcastDeltaFragment(request *rpc.BroadcastDeltaFragment) {
	panic("unimplemented")
}
