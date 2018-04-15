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
	"github.com/republicprotocol/republic-go/network"
	"google.golang.org/grpc"
)

var _ = Describe("Swarm service", func() {
	var err error
	var mu = &sync.Mutex{}
	var swarms []*network.SwarmService
	var servers []*grpc.Server

	for _, numberOfSwarms := range []int{8, 16 /* , 32, 64, 128*/} {
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
						// Bootstrap twice

						do.CoForAll(swarms, func(i int) {
							swarms[i].Bootstrap()
						})

						do.CoForAll(swarms, func(i int) {
							swarms[i].Bootstrap()
						})

						for _, swarm := range swarms {
							// Decreased from 1/2 to 1/3
							Ω(len(swarm.DHT.MultiAddresses())).Should(BeNumerically(">=", len(swarms)*1/3))
						}
					})

					It("should not error when we turn on the concurrent option for bootstrapping", func() {
						for _, swarm := range swarms {
							swarm.Concurrent = true
						}

						// Bootstrap twice

						do.CoForAll(swarms, func(i int) {
							swarms[i].Bootstrap()
						})

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

	Context("when finding nodes", func() {
		BeforeEach(func() {

			mu.Lock()

			swarms, servers, err = generateSwarmServices(5)
			Ω(err).ShouldNot(HaveOccurred())

			err = startSwarmServices(servers, swarms)
			Ω(err).ShouldNot(HaveOccurred())

			err := connectSwarms(swarms, 100)
			Ω(err).ShouldNot(HaveOccurred())

			do.CoForAll(swarms, func(i int) {
				swarms[i].Bootstrap()
			})

			time.Sleep(1 * time.Second)
		})

		AfterEach(func() {
			stopSwarmServices(servers, swarms)
			mu.Unlock()
		})

		It("should be able to find nodes", func() {

			do.CoForAll(swarms, func(i int) {
				for j := range swarms {
					if i == j {
						continue
					}

					multiaddress, err := swarms[i].FindNode(swarms[j].MultiAddress().ID())
					Ω(multiaddress.String()).Should(Equal(swarms[j].MultiAddress().String()))
					Ω(err).ShouldNot(HaveOccurred())
				}
				swarms[i].Bootstrap()
			})
		})

		It("should return nil if unable to find node", func() {
			id, _, err := identity.NewID()
			Ω(err).ShouldNot(HaveOccurred())

			multi, err := swarms[0].FindNode(id)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(multi).Should(BeNil())

		})
	})

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
			nd.Logger.Error(err.Error())
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
