package grpc_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/smpc"
	"github.com/republicprotocol/republic-go/testutils"
	"golang.org/x/net/context"
)

var _ = Describe("Streaming", func() {

	var server *Server
	var connectorListener ConnectorListener
	// var serviceMultiAddr identity.MultiAddress
	BeforeEach(func() {
		server = NewServer()
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when sending a message to a service", func() {

		It("should not error", func(done Done) {
			defer close(done)

			var addr identity.Address
			var err error

			connectorListener, addr, err = newStreamer()
			Expect(err).ShouldNot(HaveOccurred())

			// server := NewServer()
			service, _, serviceAddr, err := newStreamerService(addr)
			Expect(err).ShouldNot(HaveOccurred())
			service.Register(server)

			serviceMultiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", serviceAddr))
			Expect(err).ShouldNot(HaveOccurred())

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			sender, err := connectorListener.Connect(context.Background(), smpc.NetworkID([32]byte{}), serviceMultiAddr, testutils.NewSmpcReceiver())
			Expect(err).Should(HaveOccurred())

			sender, err = connectorListener.Connect(context.Background(), smpc.NetworkID(testutils.Random32Bytes()), serviceMultiAddr, testutils.NewSmpcReceiver())
			Expect(err).ShouldNot(HaveOccurred())

			err = sender.Send(smpc.Message{
				MessageJoin: &smpc.MessageJoin{
					Join:      smpc.Join{},
					NetworkID: smpc.NetworkID(testutils.Random32Bytes()),
				},
				MessageType: smpc.MessageTypeJoin,
			})
			Expect(err).ShouldNot(HaveOccurred())
		}, 60 /* 60 second timeout */)

		It("should connect when the service is started after the connection request", func(done Done) {
			defer close(done)

		}, 60 /* 60 second timeout */)

	})

})

func newStreamer() (ConnectorListener, identity.Address, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return ConnectorListener{}, identity.Address(""), err
	}
	addr := identity.Address(ecdsaKey.Address())
	return NewConnectorListener(addr, &ecdsaKey, testutils.NewCrypter()), addr, nil
}

func newStreamerService(clientAddr identity.Address) (*StreamerService, *ConnectorListener, identity.Address, error) {
	var streamer ConnectorListener
	var addr identity.Address
	var err error
	for {
		streamer, addr, err = newStreamer()
		if err != nil {
			return nil, &streamer, addr, err
		}
		if addr > clientAddr {
			break
		}
	}
	service := NewStreamerService(addr, crypto.NewEcdsaVerifier(clientAddr.String()), testutils.NewCrypter(), streamer.Listener)
	return &service, &streamer, addr, nil
}

type mockStreamMessage struct {
	i int64
}

func (message *mockStreamMessage) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.i); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (message *mockStreamMessage) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.i); err != nil {
		return err
	}
	return nil
}

func (message *mockStreamMessage) IsMessage() {
}
