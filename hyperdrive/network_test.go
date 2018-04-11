package hyper_test

import (
	"context"
	"sync"

	. "github.com/republicprotocol/republic-go/hyperdrive"
)

type TestNetwork struct {
	commanderCount uint8
	Ingress        []ChannelSet
	Egress         []ChannelSet
}

func NewTestNetwork(commanderCount uint8) *TestNetwork {
	return &TestNetwork{
		commanderCount: commanderCount,
		Ingress:        make([]ChannelSet, commanderCount),
		Egress:         make([]ChannelSet, commanderCount),
	}
}

func (t *TestNetwork) init() {
	for i := uint8(0); i < t.commanderCount; i++ {
		t.Egress[i] = EmptyChannelSet(t.commanderCount)
		t.Ingress[i] = EmptyChannelSet(t.commanderCount)
	}
}

func (t *TestNetwork) propose(p Proposal) {
	for i := uint8(0); i < t.commanderCount; i++ {
		go func(j uint8) {
			t.Ingress[j].Proposal <- p
		}(i)
	}
}

func (t *TestNetwork) proposeMultiple(proposals []Proposal) {
	var wg sync.WaitGroup
	for i := uint8(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()
			for _, proposal := range proposals {
				t.Ingress[i].Proposal <- proposal
			}
		}(i)
	}
	wg.Wait()
}

func (t *TestNetwork) run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint8(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()
			t.Egress[i].Split(t.Ingress)
		}(i)
	}
	wg.Wait()
}
