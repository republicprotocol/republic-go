package smpc_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/smpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
	"github.com/republicprotocol/republic-go/swarm"
)

var _ = Describe("smpc ", func() {

	Context("on starting an smpcer", func() {

		It("should return error if smpcer has already been started", func() {
			swarmer := &mockSwarmer{}
			streamer := &mockStreamer{}
			smpcer := NewSmpcer(swarmer, streamer, 10)

			err := smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())
			err = smpcer.Start()
			Expect(err).Should(HaveOccurred())
			Expect(err).To(Equal(ErrSmpcerIsAlreadyRunning))
		})
	})

	Context("on shutting down an smpcer", func() {

		It("should return error if smpcer is not running", func() {
			swarmer := &mockSwarmer{}
			streamer := &mockStreamer{}
			smpcer := NewSmpcer(swarmer, streamer, 10)

			err := smpcer.Shutdown()
			Expect(err).Should(HaveOccurred())
			Expect(err).To(Equal(ErrSmpcerIsNotRunning))
		})

		It("should not return error if smpcer is running", func() {
			swarmer := &mockSwarmer{}
			streamer := &mockStreamer{}
			smpcer := NewSmpcer(swarmer, streamer, 10)

			err := smpcer.Start()
			Expect(err).ShouldNot(HaveOccurred())
			err = smpcer.Shutdown()
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

type mockSwarmer struct {
}

func (swarmer mockSwarmer) Bootstrap(ctx context.Context, multiAddrs identity.MultiAddresses) error {
	return nil
}

func (swarmer mockSwarmer) Query(ctx context.Context, query identity.Address, depth int) (identity.MultiAddress, error) {
	n := rand.Intn(1000)
	if n <= 200 {
		return identity.MultiAddress{}, swarm.ErrMultiAddressNotFound
	}
	return query.MultiAddress()
}

type mockStreamer struct {
}

func (streamer *mockStreamer) Open(ctx context.Context, multiAddr identity.MultiAddress) (stream.Stream, error) {
	return nil, nil
}

func (streamer *mockStreamer) Close(addr identity.Address) error {
	return nil
}

type mockMessage struct {
	value string
}

func (message *mockMessage) IsMessage() {
	return
}

func (message *mockMessage) MarshalBinary() (data []byte, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, message.value)
	return buf.Bytes(), nil
}

func (message *mockMessage) UnmarshalBinary(data []byte) error {
	if data == nil || len(data) == 0 {
		return ErrUnmarshalNilBytes
	}
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.BigEndian, &message.value)
	return nil
}

func createSMPCer() (Smpcer, identity.MultiAddress) {
	swarmer := &mockSwarmer{}
	channelHub := stream.NewChannelHub()
	keystore, err := crypto.RandomKeystore()
	Expect(err).ShouldNot(HaveOccurred())
	multiaddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/80/republic/%s", keystore.Address()))
	Expect(err).ShouldNot(HaveOccurred())
	client := stream.NewChannelClient(multiaddr.Address(), &channelHub)
	server := stream.NewChannelServer(multiaddr.Address(), &channelHub)
	streamer := stream.NewStreamer(multiaddr.Address(), client, server)
	return NewSmpcer(swarmer, streamer, 10), multiaddr
}
