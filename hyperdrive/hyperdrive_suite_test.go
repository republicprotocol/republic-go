package hyper_test

import (
	"context"
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
	commanderCount uint8
	network        TestNetwork
	replicas       []Replica
}

func NewHyperDrive(commanderCount uint8) *HyperDrive {
	replicas := make([]Replica, commanderCount)
	network := NewTestNetwork(commanderCount)
	return &HyperDrive{
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
		h.replicas[i] = NewReplica(validator, h.network.Ingress[i])
	}
}

func (h *HyperDrive) run(ctx context.Context) {
	for i := uint8(0); i < h.commanderCount; i++ {
		go func(i uint8) {
			h.network.Egress[i].Copy(h.replicas[i].Run(ctx))
		}(i)
	}
	go h.network.run(ctx)
}
