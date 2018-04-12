package hyper_test

import (
	"context"
	"sync"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
		t.Ingress[i] = EmptyChannelSet(t.commanderCount)
		t.Egress[i] = EmptyChannelSet(t.commanderCount)
	}
}

func (t *TestNetwork) run(ctx context.Context) {
	var wg sync.WaitGroup
	for i := uint8(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()
			t.Egress[i].Split(ctx, t.Ingress)
		}(i)
	}
	wg.Wait()
}

var _ = Describe("Network", func() {

	Context("Test Network", func() {
		testnet := NewTestNetwork(100)
		It("Start and stop a network", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			testnet.init()
			wg.Add(1)
			go func() {
				defer wg.Done()
				testnet.run(ctx)
			}()
			cancel()
			wg.Wait()
		})

		It("Broadcasting a proposal on the network", func() {
			var wg sync.WaitGroup
			ctx, cancel := context.WithCancel(context.Background())
			testnet.init()
			wg.Add(1)
			go func() {
				defer wg.Done()
				testnet.run(ctx)
			}()

			go func() {
				defer cancel()
				testnet.Egress[0].Proposal <- Proposal{
					Signature: Signature("Hello"),
				}
				for i := uint8(0); i < testnet.commanderCount; i++ {
					proposal := <-testnet.Ingress[i].Proposal
					Ω(proposal.Signature).Should(Equal(Signature("Hello")))
				}
			}()

			wg.Wait()
		})
	})
})

// Helper functions
func (t *TestNetwork) propose(p Proposal) {
	var wg sync.WaitGroup
	for i := uint8(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(j uint8) {
			defer wg.Done()
			t.Ingress[j].Proposal <- p
		}(i)
	}
	wg.Wait()
}

func (t *TestNetwork) proposeMultiple(proposals []Proposal) {
	var wg sync.WaitGroup
	for i := uint8(0); i < t.commanderCount; i++ {
		wg.Add(1)
		go func(i uint8) {
			defer wg.Done()

			var wgInner sync.WaitGroup
			for _, proposal := range proposals {
				wgInner.Add(1)
				go func(proposal Proposal) {
					defer wgInner.Done()
					t.Ingress[i].Proposal <- proposal
				}(proposal)
			}
			wgInner.Wait()
		}(i)
	}
	wg.Wait()
}
