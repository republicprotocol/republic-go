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

var _ = Describe("Smpc", func() {

	Context("when starting", func() {

		It("should return error if smpcer has already been started", func() {
			hub := stream.NewChannelHub()
			smpcer, _, err := createSMPCer(&hub)
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
			hub := stream.NewChannelHub()
			smpcer, _, err := createSMPCer(&hub)
			Expect(err).ShouldNot(HaveOccurred())

			err = smpcer.Shutdown()
			Expect(err).Should(HaveOccurred())
			Expect(err).To(Equal(ErrSmpcerIsNotRunning))
		})

		It("should not return error if smpcer is running", func() {
			hub := stream.NewChannelHub()
			smpcer, _, err := createSMPCer(&hub)
			Expect(err).ShouldNot(HaveOccurred())
			err = smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())

			err = smpcer.Shutdown()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	Context("when joining shares", func() {

		It("should join shares to obtain final values", func() {
			// Create 16 smpcers and issue 5 random secrets
			hub := stream.NewChannelHub()
			count, err := runSmpcers(24, 5, 0, &hub)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(5))

		})

		It("should join shares when faults are below the threshold", func() {
			// Create 24 smpcers and issue 5 random secrets
			// Do not start 1/3 of the smpcers (for example: 7)
			hub := stream.NewChannelHub()
			count, err := runSmpcers(24, 5, 8, &hub)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(5))
		})

		It("should not join shares when faults are above the threshold", func() {
			hub := stream.NewChannelHub()
			timer := time.NewTimer(time.Second * 4)
			count := int32(0)

			go func() {
				defer GinkgoRecover()

				// Create 24 smpcers and do not start 16 smpcers
				c, err := runSmpcers(24, 5, 16, &hub)
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
			hub := stream.NewChannelHub()
			count, err := runSmpcersInTwoNetworks(12, 6, &hub)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(count).To(Equal(2))
		})

		It("should join when nodes are in multiple overlapping networks", func() {
			// Run 9 smpcers in 2 networks such that each network has 6 smpcers
			// (with replacement) and issue unique random secrets to each of the networks
			hub := stream.NewChannelHub()
			count, err := runSmpcersInTwoNetworks(9, 6, &hub)
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

func createSMPCer(hub *stream.ChannelHub) (Smpcer, identity.Address, error) {
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

	// Create a new channel streamer
	streamer := stream.NewChannelStreamer(multiaddr.Address(), hub)

	return NewSmpcer(swarmer, streamer, 10), multiaddr.Address(), nil
}

func runSmpcers(numberOfSmpcers, numberOfJoins, numberOfDeadNodes int, hub *stream.ChannelHub) (int, error) {

	smpcers, addrs, err := createAddressesAndStartSmpcers(numberOfSmpcers, hub)
	if err != nil {
		return 0, err
	}

	// Send connect instructions to all smpcers
	for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
		addresses := filterAddresses(numberOfSmpcers, i, addrs)
		sendConnectInstruction(numberOfSmpcers, 1, smpcers[i], addresses)
	}

	time.Sleep(2 * time.Second)

	go func() {

		// Send shares of multiple instructions
		for j := 0; j < numberOfJoins; j++ {

			shares, err := createSecretShares(numberOfSmpcers)
			if err != nil {
				log.Printf("cannot split secret: %v", err)
				return
			}

			// Send instJ instructions to join all the secret shares
			for i := numberOfDeadNodes; i < numberOfSmpcers; i++ {
				sendJoinInstruction(shares[i], j, 1, smpcers[i])
			}
		}
	}()

	return waitForResults(numberOfJoins, numberOfDeadNodes, numberOfSmpcers, smpcers), nil
}

// Runs smpcers in two networks (either overlapped or not) and sends different secrets to both
func runSmpcersInTwoNetworks(numberOfSmpcers, minimumNumberOfSmpcers int, hub *stream.ChannelHub) (int, error) {

	smpcers, addrs, err := createAddressesAndStartSmpcers(numberOfSmpcers, hub)
	if err != nil {
		return 0, err
	}

	// Create 2 lists that will store the smpcer addresses of each network
	addrsNetwork1 := addrs[:minimumNumberOfSmpcers]
	addrsNetwork2 := addrs[numberOfSmpcers-minimumNumberOfSmpcers:]

	// Send connect instructions to all smpcers in Network 1
	for i := 0; i < minimumNumberOfSmpcers; i++ {
		addresses := filterAddresses(minimumNumberOfSmpcers, i, addrsNetwork1)
		sendConnectInstruction(minimumNumberOfSmpcers, 1, smpcers[i], addresses)
	}

	// Send connect instructions to all smpcers in Network 2
	for i := numberOfSmpcers - minimumNumberOfSmpcers; i < numberOfSmpcers; i++ {
		addresses := filterAddresses(minimumNumberOfSmpcers, i-numberOfSmpcers+minimumNumberOfSmpcers, addrsNetwork2)
		sendConnectInstruction(minimumNumberOfSmpcers, 2, smpcers[i], addresses)
	}

	time.Sleep(2 * time.Second)

	go func() {

		// Create a new random secret for network1 and split it
		shares, err := createSecretShares(minimumNumberOfSmpcers)
		if err != nil {
			log.Printf("cannot split secret: %v", err)
			return
		}

		// Send instJ instructions to Network 1 join all the secret shares
		for i := 0; i < minimumNumberOfSmpcers; i++ {
			sendJoinInstruction(shares[i], 0, 1, smpcers[i])
		}

		// Create a new random secret for network2 and split it
		shares, err = createSecretShares(minimumNumberOfSmpcers)
		if err != nil {
			log.Printf("cannot split secret: %v", err)
			return
		}

		// Send instJ instructions to Network 2 join all the secret shares
		for i := numberOfSmpcers - minimumNumberOfSmpcers; i < numberOfSmpcers; i++ {
			sendJoinInstruction(shares[i-numberOfSmpcers+minimumNumberOfSmpcers], 0, 2, smpcers[i])
		}
	}()

	// Ensure results from both networks have arrived
	count := 0
	count += waitForResults(1, 0, minimumNumberOfSmpcers, smpcers)
	count += waitForResults(1, numberOfSmpcers-minimumNumberOfSmpcers, numberOfSmpcers, smpcers)

	return count, nil
}

func createAddressesAndStartSmpcers(numberOfSmpcers int, hub *stream.ChannelHub) (map[int]Smpcer, identity.Addresses, error) {
	var err error
	smpcers := make(map[int]Smpcer, numberOfSmpcers)
	addrs := make(identity.Addresses, numberOfSmpcers)

	// Generate Smpcers with addresses and start them
	for i := 0; i < numberOfSmpcers; i++ {
		smpcers[i], addrs[i], err = createSMPCer(hub)
		if err != nil {
			return smpcers, addrs, err
		}
		err = smpcers[i].Start()
		if err != nil {
			return smpcers, addrs, err
		}
	}
	return smpcers, addrs, err
}

func filterAddresses(numberOfSmpcers, removeIndex int, addrs identity.Addresses) identity.Addresses {
	// Make a list of addresses after removing the current smpcer's address
	addresses := make(identity.Addresses, numberOfSmpcers)
	copy(addresses, addrs)
	return append(addresses[:removeIndex], addresses[removeIndex+1:]...)
}

func sendConnectInstruction(numberOfSmpcers, networkID int, smpcer Smpcer, addresses identity.Addresses) {
	instConnect := InstConnect{
		K:     int64(2 * (numberOfSmpcers + 1) / 3),
		N:     int64(numberOfSmpcers),
		Nodes: addresses,
	}
	message := Inst{
		InstID:      [32]byte{1},
		NetworkID:   [32]byte{byte(networkID)},
		InstConnect: &instConnect,
	}
	smpcer.Instructions() <- message
}

func sendJoinInstruction(share shamir.Share, joinIndex, networkID int, smpcer Smpcer) {
	instJ := InstJ{
		Share: share,
	}
	message := Inst{
		InstID:    [32]byte{byte(2 + joinIndex)},
		NetworkID: [32]byte{byte(networkID)},
		InstJ:     &instJ,
	}
	smpcer.Instructions() <- message
}

func createSecretShares(numberOfSmpcers int) (shamir.Shares, error) {
	// Create a secret and split it
	secret := uint64(rand.Intn(100))
	shares, err := shamir.Split(int64(numberOfSmpcers), int64(2*(numberOfSmpcers+1)/3), secret)
	if err != nil {
		return shamir.Shares{}, err
	}
	return shares, nil
}

func waitForResults(numberOfJoins, startIndex, endIndex int, smpcers map[int]Smpcer) int {
	count := 0
	for count < numberOfJoins {
		for i := startIndex; i < endIndex; i++ {
			<-smpcers[i].Results()
		}
		count++
	}

	return count
}
