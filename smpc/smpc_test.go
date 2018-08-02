package smpc_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/leveldb"
	"github.com/republicprotocol/republic-go/swarm"
	"github.com/republicprotocol/republic-go/testutils"
)

var (
	numDarknodes = 6
	numBootstrap = 2
	α            = 4
)

var _ = Describe("Smpcer", func() {

	var nodes []*mockNode
	var stores []swarm.MultiAddressStorer
	var addresses []identity.Address

	XContext("when connecting and disconnecting", func() {
		BeforeEach(func() {
			var err error

			By("generating nodes")
			nodes, addresses, stores, err = generateMocknodes(numDarknodes, α)
			Expect(err).ShouldNot(HaveOccurred())
			bootstraps := make(identity.MultiAddresses, numBootstrap)
			for i := 0; i < numBootstrap; i++ {
				bootstraps[i] = nodes[i].Multiaddress
			}

			By("serving")
			for i := 0; i < numDarknodes; i++ {
				go func(i int) {
					Expect(nodes[i].Start()).ShouldNot(HaveOccurred())
				}(i)
			}
			time.Sleep(time.Second)

			By("bootstrapping")
			dispatch.CoForAll(nodes, func(i int) {
				defer GinkgoRecover()
				for j := 0; j < numBootstrap; j++ {
					stores[i].PutMultiAddress(bootstraps[j])
				}
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := nodes[i].Swarmer.Ping(ctx); err != nil {
					log.Println(err)
				}
			})
			time.Sleep(time.Second)
		})

		AfterEach(func() {
			os.RemoveAll("./tmp")
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
	Streamer     grpc.ConnectorListener
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

func generateMocknodes(n, α int) ([]*mockNode, []identity.Address, []swarm.MultiAddressStorer, error) {
	nodes := make([]*mockNode, n)
	addresses := make([]identity.Address, n)
	stores := make([]swarm.MultiAddressStorer, n)

	for i := range nodes {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return nil, nil, nil, err
		}

		addr := identity.Address(keystore.Address())
		multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d/republic/%v", 3000+i, addr))
		if err != nil {
			return nil, nil, nil, err
		}
		// Create leveldb store and store own multiaddress.
		db, err := leveldb.NewStore(fmt.Sprintf("./tmp/node.%v.out", i+1), 72*time.Hour)
		Expect(err).ShouldNot(HaveOccurred())
		stores[i] = db.SwarmMultiAddressStore()
		listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", 3000+i))
		if err != nil {
			return nil, nil, nil, err
		}
		err = stores[i].PutMultiAddress(multiAddr)
		if err != nil {
			return nil, nil, nil, err
		}

		swarmClient := grpc.NewSwarmClient(stores[i], multiAddr.Address())

		key, err := crypto.RandomEcdsaKey()
		if err != nil {
			return nil, nil, nil, err
		}
		swarmer := swarm.NewSwarmer(swarmClient, stores[i], α, &key)

		swarmService := grpc.NewSwarmService(swarm.NewServer(swarmer, stores[i], α), time.Microsecond)

		streamer := grpc.NewConnectorListener(addr, testutils.NewCrypter(), testutils.NewCrypter())
		streamerService := grpc.NewStreamerService(addr, testutils.NewCrypter(), testutils.NewCrypter(), streamer.Listener)

		smpcer := NewSmpcer(streamer, swarmer)

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

	return nodes, addresses, stores, nil
}
