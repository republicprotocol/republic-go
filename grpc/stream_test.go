package grpc_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/grpc"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/stream"
)

var _ = Describe("Streaming", func() {

	var server *Server
	var service *StreamerService
	var serviceStreamer *Streamer
	var serviceAddr identity.Address
	var serviceMultiAddr identity.MultiAddress
	var clientStreamer *Streamer
	var clientAddr identity.Address
	var clientMultiAddr identity.MultiAddress

	BeforeEach(func() {
		var err error

		clientStreamer, clientAddr, err = newStreamer()
		Expect(err).ShouldNot(HaveOccurred())

		server = NewServer()
		service, serviceStreamer, serviceAddr, err = newStreamerService(clientAddr)
		Expect(err).ShouldNot(HaveOccurred())
		service.Register(server)

		serviceMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18514/republic/%v", serviceAddr))
		Expect(err).ShouldNot(HaveOccurred())

		clientMultiAddr, err = identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/0.0.0.0/tcp/18515/republic/%v", clientAddr))
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

			_, err := clientStreamer.Open(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
		}, 60 /* 60 second timeout */)

		It("should connect when the service is started after the connection request", func(done Done) {
			defer close(done)

			go func() {
				defer GinkgoRecover()

				time.Sleep(time.Second)
				err := server.Start("0.0.0.0:18514")
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := clientStreamer.Open(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

		}, 60 /* 60 second timeout */)

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
			_, err := clientStreamer.Open(context.Background(), serviceMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())

			_, err = serviceStreamer.Open(context.Background(), clientMultiAddr)
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should connect when the client sends the connection request after the service is listening", func() {

			doneListening := make(chan struct{})
			go func() {
				defer GinkgoRecover()
				defer close(doneListening)

				_, err := serviceStreamer.Open(context.Background(), clientMultiAddr)
				Expect(err).ShouldNot(HaveOccurred())
			}()
			time.Sleep(time.Millisecond)

			_, err := clientStreamer.Open(context.Background(), serviceMultiAddr)
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
			clientStream, err = clientStreamer.Open(ctx, serviceMultiAddr)
			clientStreamCancel = cancel
			Expect(err).ShouldNot(HaveOccurred())

			ctx, cancel = context.WithCancel(context.Background())
			serviceStream, err = serviceStreamer.Open(ctx, clientMultiAddr)
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

		Context("when the client disconnects and reconnects", func() {
			It("should send messages to the server without the server opening a new stream", func() {
				var err error

				// Disconnect
				clientStreamCancel()
				time.Sleep(10 * time.Millisecond)

				// Confirm that receiving returns an error
				message := mockStreamMessage{}
				err = serviceStream.Recv(&message)
				Expect(err).Should(HaveOccurred())

				// Reconnect
				ctx, cancel := context.WithCancel(context.Background())
				clientStream, err = clientStreamer.Open(ctx, serviceMultiAddr)
				clientStreamCancel = cancel
				Expect(err).ShouldNot(HaveOccurred())
				time.Sleep(10 * time.Millisecond)

				// Send a message from the client
				go func() {
					defer GinkgoRecover()
					err := clientStream.Send(&mockStreamMessage{int64(420)})
					Expect(err).ShouldNot(HaveOccurred())
				}()

				// Receive a message from the service without opening a new
				// stream
				err = serviceStream.Recv(&message)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message.i).Should(Equal(int64(420)))
			})

			It("should receive messages from the server without the server opening a new stream", func() {
				var err error

				// Disconnect
				clientStreamCancel()
				time.Sleep(10 * time.Millisecond)

				// Confirm that sending and receiving returns an error
				err = serviceStream.Send(&mockStreamMessage{int64(420)})
				Expect(err).Should(HaveOccurred())

				// Reconnect
				ctx, cancel := context.WithCancel(context.Background())
				clientStream, err = clientStreamer.Open(ctx, serviceMultiAddr)
				clientStreamCancel = cancel
				Expect(err).ShouldNot(HaveOccurred())
				time.Sleep(10 * time.Millisecond)

				// Send a message from the client
				go func() {
					defer GinkgoRecover()
					err := serviceStream.Send(&mockStreamMessage{int64(420)})
					Expect(err).ShouldNot(HaveOccurred())
				}()

				// Receive a message from the service without opening a new
				// stream
				message := mockStreamMessage{}
				err = clientStream.Recv(&message)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message.i).Should(Equal(int64(420)))
			})
		})

		Context("when the server disconnects and reconnects", func() {
			It("should send messages to the client without the client opening a new stream", func(done Done) {
				defer close(done)

				// Disconnect
				serviceStreamCancel()
				time.Sleep(10 * time.Millisecond)

				// Send a message from the service
				go func() {
					defer GinkgoRecover()

					// Reconnect
					ctx, cancel := context.WithCancel(context.Background())
					serviceStream, err := serviceStreamer.Open(ctx, clientMultiAddr)
					serviceStreamCancel = cancel
					Expect(err).ShouldNot(HaveOccurred())
					time.Sleep(10 * time.Millisecond)

					err = serviceStream.Send(&mockStreamMessage{int64(420)})
					Expect(err).ShouldNot(HaveOccurred())
				}()

				// Receive a message from the client without opening a new
				// stream
				message := mockStreamMessage{}
				err := clientStream.Recv(&message)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message.i).Should(Equal(int64(420)))

			}, 60 /* 60s timeout */)

			It("should receive messages from the client without the client opening a new stream", func(done Done) {
				defer close(done)

				// Disconnect
				serviceStreamCancel()
				time.Sleep(10 * time.Millisecond)

				// Send a message from the service
				go func() {
					defer GinkgoRecover()

					err := clientStream.Send(&mockStreamMessage{int64(420)})
					Expect(err).ShouldNot(HaveOccurred())
				}()

				// Reconnect
				ctx, cancel := context.WithCancel(context.Background())
				serviceStream, err := serviceStreamer.Open(ctx, clientMultiAddr)
				serviceStreamCancel = cancel
				Expect(err).ShouldNot(HaveOccurred())
				time.Sleep(10 * time.Millisecond)

				// Receive a message from the client without opening a new
				// stream
				message := mockStreamMessage{}
				err = serviceStream.Recv(&message)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(message.i).Should(Equal(int64(420)))
			}, 60 /* 60s timeout */)
		})
	})
})

func newStreamer() (*Streamer, identity.Address, error) {
	ecdsaKey, err := crypto.RandomEcdsaKey()
	if err != nil {
		return nil, identity.Address(""), err
	}
	addr := identity.Address(ecdsaKey.Address())
	return NewStreamer(&ecdsaKey, testutils.NewCrypter(), addr), addr, nil
}

func newStreamerService(clientAddr identity.Address) (*StreamerService, *Streamer, identity.Address, error) {
	var streamer *Streamer
	var addr identity.Address
	var err error
	for {
		streamer, addr, err = newStreamer()
		if err != nil {
			return nil, streamer, addr, err
		}
		if addr > clientAddr {
			break
		}
	}
	service := NewStreamerService(crypto.NewEcdsaVerifier(clientAddr.String()), testutils.NewCrypter(), streamer)
	return &service, streamer, addr, nil
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
