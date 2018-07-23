package smpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dht"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	numDarknodes = 6
	numBootstrap = 2
)

var _ = Describe("Smpcer", func() {

	var nodes []*mockNode
	var addresses []identity.Address

	Context("when connecting and disconnecting", func() {
		BeforeEach(func() {
			var err error

			By("generating nodes")
			nodes, addresses, err = generateMocknodes(numDarknodes)
			Expect(err).ShouldNot(HaveOccurred())
			bootstraps := make(identity.MultiAddresses, numBootstrap)
			for i := 0; i < numBootstrap; i++ {
				bootstraps[i] = nodes[i].Multiaddress
			}

			By("serving")
			for i := 0; i < numDarknodes; i++ {
				go func(i int) {
					Ω(nodes[i].Start()).ShouldNot(HaveOccurred())
				}(i)
			}
			time.Sleep(time.Second)

			By("bootstrapping")
			dispatch.CoForAll(nodes, func(i int) {
				defer GinkgoRecover()
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := nodes[i].Swarmer.Bootstrap(ctx, bootstraps); err != nil {
					log.Println(err)
				}
			})
			time.Sleep(time.Second)
		})

		AfterEach(func() {
			dispatch.CoForAll(nodes, func(i int) {
				nodes[i].Stop()
			})
			time.Sleep(time.Second)
		})

		It("should be able to connect and disconnect to the smpcer", func() {
			networkID := NetworkID{0}
			for i := range nodes {
				nodes[i].Smpcer.Connect(networkID, addresses)
			}
			for i := range nodes {
				nodes[i].Smpcer.Disconnect(networkID)
			}
		})

		Context("when connecting before joining", func() {
			It("should be able to join values", func() {
				networkID := NetworkID{0}
				for i := range nodes {
					nodes[i].Smpcer.Connect(networkID, addresses)
				}
				time.Sleep(time.Second)

				var called int64
				ord, joins := generateJoins(int64(numDarknodes), int64(2*(numDarknodes+1)/3))
				callback := generateCallback(&called, ord)

				dispatch.CoForAll(nodes, func(i int) {
					err := nodes[i].Smpcer.Join(networkID, joins[i], callback)
					Expect(err).ShouldNot(HaveOccurred())
				})
				for atomic.LoadInt64(&called) < int64(numDarknodes) {
					time.Sleep(time.Second)
				}

				for i := range nodes {
					nodes[i].Smpcer.Disconnect(networkID)
				}
			})

			// It("should be able to join values in multiple networks", func() {
			// })
		})

		// Context("when connecting after joining", func() {
		// 	It("should be able to join values", func() {
		// 	})
		// 	It("should be able to join values in multiple networks", func() {
		// 	})
		// })

	})
})

type mockNode struct {
	Address      identity.Address
	Multiaddress identity.MultiAddress

	Host     string
	Port     string
	Listener net.Listener
	Server   *grpc.Server

	Swarmer      swarm.Swarmer
	SwarmService grpc.SwarmService
	Streamer     stream.Streamer
	Smpcer       Smpcer
}

func (node *mockNode) Start() error {
	log.Printf("listening on %v:%v...", node.Host, node.Port)
	return node.Server.Serve(node.Listener)
}

func (node *mockNode) Stop() {
	node.Server.Stop()
	node.Listener.Close()
}

func generateMocknodes(n int) ([]*mockNode, []identity.Address, error) {
	nodes := make([]*mockNode, n)
	addresses := make([]identity.Address, n)

	for i := range nodes {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return nil, nil, err
		}

		addr := identity.Address(keystore.Address())
		dht := dht.NewDHT(addr, n)
		multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%v", 3000+i, addr))
		if err != nil {
			return nil, nil, err
		}
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 3000+i))
		if err != nil {
			return nil, nil, err
		}

		swarmClient := grpc.NewSwarmClient(multiAddr)
		swarmer := swarm.NewSwarmer(swarmClient, &dht)
		swarmService := grpc.NewSwarmService(swarm.NewServer(testutils.NewCrypter(), swarmClient, &dht))

		streamer := grpc.NewStreamer(testutils.NewCrypter(), testutils.NewCrypter(), addr)
		streamerService := grpc.NewStreamerService(testutils.NewCrypter(), testutils.NewCrypter(), streamer)

		smpcer := NewSmpcer(swarmer, streamer)

		addresses[i] = addr
		nodes[i] = new(mockNode)
		nodes[i].Address = addr
		nodes[i].Multiaddress = multiAddr
		nodes[i].Host = "127.0.0.1"
		nodes[i].Port = fmt.Sprintf("%d", 3000+i)
		nodes[i].Listener = listener
		nodes[i].Server = grpc.NewServer()
		nodes[i].Swarmer = swarmer
		nodes[i].SwarmService = swarmService
		nodes[i].Streamer = streamer
		nodes[i].Smpcer = smpcer

		swarmService.Register(nodes[i].Server)
		streamerService.Register(nodes[i].Server)
	}

	return nodes, addresses, nil
}
