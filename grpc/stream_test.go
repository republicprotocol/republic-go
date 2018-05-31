package grpc_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	. "github.com/republicprotocol/republic-go/grpc"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
)

var _ = Describe("Streaming", func() {

	var server *Server
	var service *StreamService
	var serviceAddr identity.Address
	var serviceMultiAddr identity.MultiAddress
	var client stream.Client
	var clientAddr identity.Address

	BeforeEach(func() {
		var err error

		server = NewServer()
		service, serviceAddr, err = newStreamService()
		Expect(err).ShouldNot(HaveOccurred())
		service.Register(server)

		client, clientAddr, err = newStreamClient()
		Expect(err).ShouldNot(HaveOccurred())

		serviceMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", serviceAddr))
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		server.Stop()
	})

	Context("when connecting to a service", func() {

		It("should connect when the service is started before the connection request", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := client.Connect(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should connect when the service is started after the connection request", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				time.Sleep(time.Second)
				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := client.Connect(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

		}, 30 /* 30 second timeout */)

	})

	Context("when listening to a client", func() {

		BeforeEach(func() {
			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)
		})

		It("should connect when the client sends the connection request before the service is listening", func() {
			_, err := client.Connect(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = service.Listen(context.Background(), clientAddr)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should connect when the client sends the connection request after the service is listening", func() {

			doneListening := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(doneListening)

				_, err := service.Listen(context.Background(), clientAddr)
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := client.Connect(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			<-doneListening
		})

	})

	Context("when sending and receiving messages", func() {

		var serviceStream stream.Stream
		var serviceStreamCancel context.CancelFunc

		var clientStream stream.Stream
		var clientStreamCancel context.CancelFunc

		BeforeEach(func() {
			var err error

			go func() {
				defer GinkgoRecover()

				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			ctx, cancel := context.WithCancel(context.Background())
			clientStream, err = client.Connect(ctx, serviceMultiAddr)
			clientStreamCancel = cancel
			Expect(err).ShouldNot(HaveOccurred())

			ctx, cancel = context.WithCancel(context.Background())
			serviceStream, err = service.Listen(context.Background(), clientAddr)
			serviceStreamCancel = cancel
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			clientStreamCancel()
			serviceStreamCancel()
		})

		It("should receive messages sent by the client", func() {
			dispatch.CoBegin(func() {
				defer GinkgoRecover()

				for i := 0; i < 1000; i++ {
					err := clientStream.Send(&mockStreamMessage{int64(i)})
					Expect(err).ShouldNot(HaveOccurred())
				}
			}, func() {
				defer GinkgoRecover()

				for i := 0; i < 1000; i++ {
					message := mockStreamMessage{}
					err := serviceStream.Recv(&message)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(message.i).To(Equal(int64(i)))
				}
			})
		})

		It("should receive messages sent by the server", func() {
			dispatch.CoBegin(func() {
				defer GinkgoRecover()

				for i := 0; i < 1000; i++ {
					err := serviceStream.Send(&mockStreamMessage{int64(i)})
					Expect(err).ShouldNot(HaveOccurred())
				}
			}, func() {
				defer GinkgoRecover()

				for i := 0; i < 1000; i++ {
					message := mockStreamMessage{}
					err := clientStream.Recv(&message)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(message.i).To(Equal(int64(i)))
				}
			})
		})
	})
})

func newStreamClient() (stream.Client, identity.Address, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, identity.Address(""), err
	}
	addr := identity.Address(ecdsaKey.Address())
	client := NewStreamClient(mockSigner{}, addr)
	return client, addr, nil
}

func newStreamService() (*StreamService, identity.Address, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, identity.Address(""), err
	}
	addr := identity.Address(ecdsaKey.Address())
	service := NewStreamService(mockVerifier{}, addr)
	return &service, addr, nil
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
