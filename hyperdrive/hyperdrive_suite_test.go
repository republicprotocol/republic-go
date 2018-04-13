package hyper_test

import (
	"context"
	"log"
	"sync"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/hyperdrive"
)

func TestHyperdrive(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Hyperdrive Suite")
}

type HyperDrive struct {
	ctx            context.Context
	commanderCount uint8
	network        TestNetwork
	replicas       []Replica
}

func NewHyperDrive(ctx context.Context, commanderCount uint8) *HyperDrive {
	replicas := make([]Replica, commanderCount)
	network := NewTestNetwork(ctx, commanderCount)
	return &HyperDrive{
		ctx:            ctx,
		commanderCount: commanderCount,
		network:        *network,
		replicas:       replicas,
	}
}

func (h *HyperDrive) init() {
	h.network.init()
	for i := uint8(0); i < h.commanderCount; i++ {
		blocks := NewSharedBlocks(0, 0)
		validator, _ := NewTestValidator(blocks, 2*(h.commanderCount/3))
		h.replicas[i] = NewReplica(h.ctx, validator, h.network.Ingress[i])
	}
}

func (h *HyperDrive) run() {
	defer log.Println("Stopping the hyperdrive")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		h.network.run()
	}()

	for i := uint8(0); i < h.commanderCount; i++ {
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()
			go h.network.Egress[i].Copy(h.replicas[i].Run())
		}(i)
	}

	wg.Wait()
}
