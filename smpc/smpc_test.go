package smpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
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
	NumberOfdarknodes      = 24
	NumberOfBootstrapNodes = 6
	BuckSize               = 50
)

var _ = Describe("Smpcer", func() {
	var nodes []*mockNode
	var addresses []identity.Address

	Context("Connect and disconnect", func() {
		BeforeEach(func() {
			By("Setting up test environment")
			nodes, addresses = generateMocknodes(NumberOfdarknodes)
			bootstraps := make([]identity.MultiAddress, NumberOfBootstrapNodes)
			for i := 0; i < NumberOfBootstrapNodes; i++ {
				bootstraps[i] = nodes[i].Multiaddress
			}

			By("Start serving and bootstrapping")
			for i := 0; i < NumberOfdarknodes; i++ {
				go func(i int) {
					Ω(nodes[i].Start()).ShouldNot(HaveOccurred())
				}(i)
			}

			dispatch.CoForAll(nodes, func(i int) {
				Ω(nodes[i].Swarmer.Bootstrap(context.Background(), bootstraps)).ShouldNot(HaveOccurred())

			})

			time.Sleep(3 * time.Second)
		})

		AfterEach(func() {
			for i := 0; i < NumberOfdarknodes; i++ {
				nodes[i].Stop()
			}
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

		It("should be able to join Joins and call the callback ", func() {
			var networkID [32]byte
			copy(networkID[:], crypto.Keccak256([]byte{1}))
			ord, joins := generateJoins()
			var called int64
			callback := generateCallback(&called, ord)

			for i := range nodes {
				nodes[i].Smpcer.Connect(networkID, addresses)
			}

			dispatch.CoForAll(nodes, func(i int) {
				err := nodes[i].Smpcer.Join(networkID, joins[i], callback)
				Ω(err).ShouldNot(HaveOccurred())
			})

			//Ω(atomic.LoadInt64(&called)).Should(Equal(NumberOfdarknodes))

			for i := range nodes {
				nodes[i].Smpcer.Disconnect(networkID)
			}
		})

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
	StreamClient stream.Client
	Streamer     stream.Streamer
	Smpcer       Smpcer
}

func (node *mockNode) Start() error {
	log.Printf("listening on %v:%v...", node.Host, node.Port)
	return node.Server.Serve(node.Listener)
}

func (node *mockNode) Stop() {
	node.Server.Stop()
}

func generateMocknodes(n int) ([]*mockNode, []identity.Address) {
	nodes := make([]*mockNode, n)
	addresses := make([]identity.Address, n)

	for i := range nodes {
		nodes[i] = new(mockNode)
		keystore, err := crypto.RandomKeystore()
		Ω(err).ShouldNot(HaveOccurred())

		addr := identity.Address(keystore.Address())
		dht := dht.NewDHT(addr, BuckSize)
		multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%v", 4000+i, addr))
		Ω(err).ShouldNot(HaveOccurred())
		connPool := grpc.NewConnPool(n)
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 4000+i))
		Ω(err).ShouldNot(HaveOccurred())

		swarmClient := grpc.NewSwarmClient(multiAddr, &connPool)
		swarmer := swarm.NewSwarmer(swarmClient, &dht)
		swarmServices := grpc.NewSwarmService(swarm.NewServer(swarmClient, &dht))

		streamClient := grpc.NewStreamClient(testutils.NewMockCrypter(), addr)
		streamService := grpc.NewStreamService(testutils.NewMockCrypter(), addr)
		streamer := stream.NewStreamRecycler(stream.NewStreamer(addr, streamClient, &streamService))

		smpcer := NewSmpcer(swarmer, streamer)

		addresses[i] = addr
		nodes[i].Address = addr
		nodes[i].Multiaddress = multiAddr
		nodes[i].Host = "127.0.0.1"
		nodes[i].Port = fmt.Sprintf("%d", 3000+i)
		nodes[i].Listener = listener
		nodes[i].Server = grpc.NewServer()
		nodes[i].Swarmer = swarmer
		nodes[i].SwarmService = swarmServices
		nodes[i].StreamClient = streamClient
		nodes[i].Streamer = streamer
		nodes[i].Smpcer = smpcer

		swarmServices.Register(nodes[i].Server)
		streamService.Register(nodes[i].Server)
	}

	return nodes, addresses
}
