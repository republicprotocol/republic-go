package network_test

import (
	"log"
	"net"
	"sync"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"google.golang.org/grpc"
)

var _ = Describe("Swarm service", func() {
	var err error
	var mu = new(sync.Mutex)
	var swarms []*network.SwarmService
	var servers []*grpc.Server

	for _, numberOfSwarms := range []int{8 ,16 /* , 32, 64, 128*/} {
		for _, connectivity := range []int{100 /*, 80, 40*/} {
			func(numberOfNodes, connectivity int) {
				Context("when bootstrapping", func() {
					BeforeEach(func() {
						mu.Lock()

						swarms, servers, err = generateSwarmServices(numberOfNodes)
						Ω(err).ShouldNot(HaveOccurred())

						err = startSwarmServices(servers, swarms)
						Ω(err).ShouldNot(HaveOccurred())

						err := connectSwarms(swarms, connectivity)
						Ω(err).ShouldNot(HaveOccurred())

						time.Sleep(1 * time.Second)
					})

					AfterEach(func() {
						stopSwarmServices(servers, swarms)
						mu.Unlock()
					})

					It("should be able to find most of the nodes in the network after bootstrapping ", func() {
						do.CoForAll(swarms, func(i int) {
							swarms[i].Bootstrap()
						})

						for _, swarm := range swarms {
							Ω(len(swarm.DHT.MultiAddresses())).Should(BeNumerically(">=", len(swarms)*1/2))
						}
					})

					It("should not error when we turn on the concurrent option for bootstrapping", func() {
						for _, swarm := range swarms {
							swarm.Concurrent = true
						}

						do.CoForAll(swarms, func(i int) {
							swarms[i].Bootstrap()
						})

						for _, swarm := range swarms {
							Ω(len(swarm.DHT.MultiAddresses())).Should(BeNumerically(">=", len(swarms)*1/2))
						}
					})
				})
			}(numberOfSwarms, connectivity)
		}
	}
})

func startSwarmServices(servers []*grpc.Server, swarms []*network.SwarmService) error {
	for i, nd := range swarms {
		server := grpc.NewServer(grpc.ConnectionTimeout(time.Minute))
		servers[i] = server

		nd.Register(server)
		host, err := nd.MultiAddress().ValueForProtocol(identity.IP4Code)
		if err != nil {
			return err
		}
		port, err := nd.MultiAddress().ValueForProtocol(identity.TCPCode)
		if err != nil {
			return err
		}
		listener, err := net.Listen("tcp", host+":"+port)

		if err != nil {
			nd.Logger.Error(logger.TagNetwork, err.Error())
		}
		go func() {
			if err := server.Serve(listener); err != nil {
				log.Println("fail to start the server")
			}
		}()
	}

	return nil
}

func stopSwarmServices(servers []*grpc.Server, swarms []*network.SwarmService) {
	for _, server := range servers {
		server.Stop()
	}
	for _, swarm := range swarms {
		swarm.Logger.Stop()
	}
}
