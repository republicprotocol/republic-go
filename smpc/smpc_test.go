package smpc_test

import (
	"context"
	"log"
	"math/rand"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/stream"
)

var channelHub stream.ChannelHub

var _ = Describe("Smpc", func() {
	BeforeEach(func() {
		channelHub = stream.NewChannelHub()
	})

	Context("when starting", func() {

		It("should return error if smpcer has already been started", func() {
			smpcer, _, err := createSMPCer()
			Expect(err).ShouldNot(HaveOccurred())
			err = smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())

			// On starting Smpcer again, should throw an error
			err = smpcer.Start()
			Expect(err).Should(HaveOccurred())
			Expect(err).To(Equal(ErrSmpcerIsAlreadyRunning))
		})
	})

	Context("when shutting down", func() {

		It("should return error if smpcer is not running", func() {
			smpcer, _, err := createSMPCer()
			Expect(err).ShouldNot(HaveOccurred())

			err = smpcer.Shutdown()
			Expect(err).Should(HaveOccurred())
			Expect(err).To(Equal(ErrSmpcerIsNotRunning))
		})

		It("should not return error if smpcer is running", func() {
			smpcer, _, err := createSMPCer()
			Expect(err).ShouldNot(HaveOccurred())
			err = smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())

			err = smpcer.Shutdown()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when joining shares", func() {

		FIt("should join shares to obtain final values", func() {
			// Create 16 smpcers and issue 5 random secrets
			count, err := runSmpcers(24, 5, 0)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(5))

		})

		It("should join shares when faults are below the threshold", func() {
			// Create 24 smpcers and issue 5 random secrets
			// Do not start 1/3 of the smpcers (for example: 7)
			count, err := runSmpcers(24, 5, 7)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(5))
		})

		It("should not join shares when faults are above the threshold", func() {
			timer := time.NewTimer(time.Second * 4)
			count := int32(0)

			go func() {
				defer GinkgoRecover()

				// Create 24 smpcers and do not start 10 smpcers
				c, err := runSmpcers(24, 5, 10)
				Expect(err).ShouldNot(HaveOccurred())
				atomic.StoreInt32(&count, int32(c))
			}()

			// Wait until the timer goes off
			<-timer.C
			Expect(atomic.LoadInt32(&count)).To(Equal(int32(0)))
		})

		It("should join when nodes are in multiple non-overlapping networks", func() {
			// Run 12 smpcers in 2 networks such that each network has 6 smpcers which
			// in such a way the that each network has seperate smpcers.
			count, err := runSmpcersInTwoNetworks(12, 6)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("should join when nodes are in multiple overlapping networks", func() {
			// Run 9 smpcers in 2 networks such that each network has 6 smpcers
			// (with replacement) and issue unique random secrets to each of the networks
			count, err := runSmpcersInTwoNetworks(9, 6)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(2))
		})

	})
})

type mockSwarmer struct {
}

func (swarmer mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	return nil
}

func (swarmer mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	return query.MultiAddress()
}

func createSMPCer() (Smpcer, identity.Address, error) {
	swarmer := &mockSwarmer{}

	// Generate multiaddress
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, "", err
	}
	multiaddr, err := identity.Address(ecdsaKey.Address()).MultiAddress()
	if err != nil {
		return nil, "", err
	}

	// Create client and server channels
	client := stream.NewChannelClient(multiaddr.Address(), &channelHub)
	server := stream.NewChannelServer(multiaddr.Address(), &channelHub)

	// Create a new channel streamer
	streamer := stream.NewStreamRecycler(stream.NewStreamer(multiaddr.Address(), client, server))

	return NewSmpcer(swarmer, streamer, 10), multiaddr.Address(), nil
}

func runSmpcers(numberOfSmpcers, numberOfJoins, numberOfDeadNodes int) (int, error) {
	var err error
	smpcers := make(map[int]Smpcer, numberOfSmpcers)
	addrs := make(identity.Addresses, numberOfSmpcers)

	// Generate Smpcers with addresses
	for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
		smpcers[i], addrs[i], err = createSMPCer()
		if err != nil {
			return 0, err
		}
		err = smpcers[i].Start()
		if err != nil {
			return 0, err
		}
	}

	// Send connect instructions to all smpcers
	for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
		// Make a list of addresses after removing the current smpcer's address
		addresses := make(identity.Addresses, numberOfSmpcers)
		copy(addresses, addrs)
		addresses = append(addresses[:i], addresses[i+1:]...)

		instConnect := InstConnect{
			K:     int64(2 * (numberOfSmpcers + 1) / 3),
			N:     int64(numberOfSmpcers),
			Nodes: addresses,
		}
		message := Inst{
			InstID:      [32]byte{1},
			NetworkID:   [32]byte{1},
			InstConnect: &instConnect,
		}
		smpcers[i].Instructions() <- message
	}

	time.Sleep(2 * time.Second)

	go func() {

		// Send shares of multiple instructions
		for j := 0; j < numberOfJoins; j++ {

			// Create a secret and split it
			secret := uint64(rand.Intn(100))
			shares, err := shamir.Split(int64(numberOfSmpcers), int64(2*(numberOfSmpcers+1)/3), secret)
			if err != nil {
				log.Printf("cannot split secret: %v", err)
				return
			}

			// Send instJ instructions to join all the secret shares
			for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
				instJ := InstJ{
					Share: shares[i],
				}
				message := Inst{
					InstID:    [32]byte{byte(2 + j)},
					NetworkID: [32]byte{1},
					InstJ:     &instJ,
				}
				smpcers[i].Instructions() <- message
			}
		}
	}()

	log.Println("waiting for results...")

	count := 0
	for count < numberOfJoins {
		for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
			<-smpcers[i].Results()
		}
		count++
	}

	return count, nil
}

func runSmpcersInTwoNetworks(numberOfSmpcers, minimumNumberOfSmpcers int) (int, error) {
	var err error
	smpcers := make(map[int]Smpcer, numberOfSmpcers)
	addrs := make(identity.Addresses, numberOfSmpcers)

	// Generate Smpcers with addresses
	for i := 0; i < numberOfSmpcers; i++ {
		smpcers[i], addrs[i], err = createSMPCer()
		if err != nil {
			return 0, err
		}
		err = smpcers[i].Start()
		if err != nil {
			return 0, err
		}
	}

	// Create 2 lists that will store the smpcer addresses of each network
	addrsNetwork1 := addrs[:minimumNumberOfSmpcers]
	addrsNetwork2 := addrs[numberOfSmpcers-minimumNumberOfSmpcers:]

	// Send connect instructions to all smpcers in Network 1
	for i := 0; i < minimumNumberOfSmpcers; i++ {
		// Make a list of addresses after removing the current smpcer's address
		addresses := make(identity.Addresses, minimumNumberOfSmpcers)
		copy(addresses, addrsNetwork1)
		addresses = append(addresses[:i], addresses[i+1:]...)

		instConnect := InstConnect{
			K:     int64(2 * (minimumNumberOfSmpcers + 1) / 3),
			N:     int64(minimumNumberOfSmpcers),
			Nodes: addresses,
		}
		message := Inst{
			InstID:      [32]byte{1},
			NetworkID:   [32]byte{1},
			InstConnect: &instConnect,
		}
		smpcers[i].Instructions() <- message
	}

	// Send connect instructions to all smpcers in Network 2
	for i := numberOfSmpcers - minimumNumberOfSmpcers; i < numberOfSmpcers; i++ {
		// Make a list of addresses after removing the current smpcer's address
		addresses := make(identity.Addresses, minimumNumberOfSmpcers)
		copy(addresses, addrsNetwork2)
		addresses = append(addresses[:i-numberOfSmpcers+minimumNumberOfSmpcers], addresses[i-numberOfSmpcers+minimumNumberOfSmpcers+1:]...)

		instConnect := InstConnect{
			K:     int64(2 * (minimumNumberOfSmpcers + 1) / 3),
			N:     int64(minimumNumberOfSmpcers),
			Nodes: addresses,
		}
		message := Inst{
			InstID:      [32]byte{1},
			NetworkID:   [32]byte{2},
			InstConnect: &instConnect,
		}
		smpcers[i].Instructions() <- message
	}

	time.Sleep(2 * time.Second)

	go func() {
		// Create a secret and split it
		secret := uint64(rand.Intn(100))
		shares, err := shamir.Split(int64(minimumNumberOfSmpcers), int64(2*(minimumNumberOfSmpcers+1)/3), secret)
		if err != nil {
			log.Printf("cannot split secret: %v", err)
			return
		}

		// Send instJ instructions to Network 1 join all the secret shares
		for i := 0; i < minimumNumberOfSmpcers; i++ {
			instJ := InstJ{
				Share: shares[i],
			}
			message := Inst{
				InstID:    [32]byte{byte(2)},
				NetworkID: [32]byte{byte(1)},
				InstJ:     &instJ,
			}
			smpcers[i].Instructions() <- message
		}

		// Create a new random secret for network2 and split it
		secret = uint64(rand.Intn(100))
		shares, err = shamir.Split(int64(minimumNumberOfSmpcers), int64(2*(minimumNumberOfSmpcers+1)/3), secret)
		if err != nil {
			log.Printf("cannot split secret: %v", err)
			return
		}

		// Send instJ instructions to Network 2 join all the secret shares
		for i := numberOfSmpcers - minimumNumberOfSmpcers; i < numberOfSmpcers; i++ {
			instJ := InstJ{
				Share: shares[i-numberOfSmpcers+minimumNumberOfSmpcers],
			}
			message := Inst{
				InstID:    [32]byte{byte(2)},
				NetworkID: [32]byte{byte(2)},
				InstJ:     &instJ,
			}
			smpcers[i].Instructions() <- message
		}
	}()

	count := 0
	// Ensure results from network1 has arrived
	for i := 0; i < minimumNumberOfSmpcers; i++ {
		<-smpcers[i].Results()
	}
	count++

	// Ensure results from network2 has arrived
	for i := numberOfSmpcers - minimumNumberOfSmpcers; i < numberOfSmpcers; i++ {
		<-smpcers[i].Results()
	}
	count++

	return count, nil
}
