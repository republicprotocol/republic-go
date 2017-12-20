package swarm_test

import (
	. "github.com/republicprotocol/go-swarm"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-swarm"
	"github.com/republicprotocol/go-identity"
	"strconv"
	"net"
	"google.golang.org/grpc"
	"github.com/republicprotocol/go-swarm/rpc"
)

const (
	Number_Of_Nodes = 25
	Localhost = "127.0.0.1"
)

var _ = Describe("Node", func() {

	Describe("Star topology", func() {
		var (
			nodes 		[Number_Of_Nodes]*Node
			center 		identity.Address
			testingInt  int
			)

		BeforeEach(func() {
			// Initialize the center node of the network
			secp,err  := identity.NewKeyPair()
			Ω(err).ShouldNot(HaveOccurred())
			nodes[0] = swarm.NewNode(Localhost,"7000",secp.PublicAddress())
			center = secp.PublicAddress()
			lis, err := net.Listen("tcp", Localhost+":7000")
			Ω(err).ShouldNot(HaveOccurred())
			s := grpc.NewServer()
			rpc.RegisterDHTServer(s, nodes[0])
			go s.Serve(lis)

			// Initialize other nodes
			for i:= 1;i<Number_Of_Nodes;i++{
				keyPair,err := identity.NewKeyPair()
				Ω(err).ShouldNot(HaveOccurred())
				nodes[i] = swarm.NewNode(Localhost,strconv.Itoa(7000+i),keyPair.PublicAddress())
				lis, err := net.Listen("tcp", Localhost+":"+strconv.Itoa(7000+i))
				Ω(err).ShouldNot(HaveOccurred())
				s := grpc.NewServer()
				rpc.RegisterDHTServer(s, nodes[i])
				go s.Serve(lis)

				// Ping the center node to establish connection
				pong ,err := nodes[i].PingNode(Localhost+":"+"7000")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(pong.Address).Should(Equal(string(center)))
			}
			testingInt = 1
		})

		AfterEach(func() {
			nodes = [Number_Of_Nodes]*Node{}
		})


		Context("Test connection of two random nodes", func() {
			It("should be able to find each other in the network", func() {
				target := nodes[3].DHT.Address
				multi, err := nodes[24].FindNode(string(target))
				Ω(err).ShouldNot(HaveOccurred())
				mAddress,err := nodes[3].MultiAddress()
				Ω(err).ShouldNot(HaveOccurred())
				Ω(multi.Multis[0]).Should(Equal(mAddress.String()))
			})
		})
	})


})
