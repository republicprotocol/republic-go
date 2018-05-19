package stream_test

import (
	"context"
	"sync"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stream"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

var _ = Describe("Streaming", func() {

	Context("when using streamers", func() {

		It("should abstract connecting to servers and listening for client", func() {

			multiAddrs := [128]identity.MultiAddress{}
			clients := [128]mockClient{}
			servers := [128]mockServer{}
			streamers := [128]Streamer{}

			for i := 0; i < 128; i++ {
				ecdsaKey, err := crypto.RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				multiAddrs[i], err = identity.Address(ecdsaKey.Address()).MultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i])
			}

			for i := 0; i < 128; i++ {
				for j := 0; j < 128; j++ {
					if i == j {
						continue
					}
					_, err := streamers[i].Open(context.Background(), multiAddrs[j])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			for i := 0; i < 128; i++ {
				Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(127))
			}
		})

	})

	Context("when using stream recyclers", func() {

		It("should recycle streams for multiple connections", func() {

			multiAddrs := [128]identity.MultiAddress{}
			clients := [128]mockClient{}
			servers := [128]mockServer{}
			streamers := [128]Streamer{}

			for i := 0; i < 128; i++ {
				ecdsaKey, err := crypto.RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				multiAddrs[i], err = identity.Address(ecdsaKey.Address()).MultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			for conns := 0; conns < 4; conns++ {
				for i := 0; i < 128; i++ {
					for j := 0; j < 128; j++ {
						if i == j {
							continue
						}
						_, err := streamers[i].Open(context.Background(), multiAddrs[j])
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
			}

			for i := 0; i < 128; i++ {
				Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(127))
			}

		})

		It("should not close streams until all references have closed", func() {

			multiAddrs := [128]identity.MultiAddress{}
			clients := [128]mockClient{}
			servers := [128]mockServer{}
			streamers := [128]Streamer{}

			for i := 0; i < 128; i++ {
				ecdsaKey, err := crypto.RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				multiAddrs[i], err = identity.Address(ecdsaKey.Address()).MultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			for conns := 0; conns < 4; conns++ {
				for i := 0; i < 128; i++ {
					for j := 0; j < 128; j++ {
						if i == j {
							continue
						}
						_, err := streamers[i].Open(context.Background(), multiAddrs[j])
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
			}

			for conns := 0; conns < 4; conns++ {
				for i := 0; i < 128; i++ {
					for j := 0; j < 128; j++ {
						if i == j {
							continue
						}
						err := streamers[i].Close(multiAddrs[j].Address())
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
				for i := 0; i < 128; i++ {
					if conns == 3 {
						Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(0))
					} else {
						Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(127))
					}
				}
			}
		})

		It("should return an error when closing streams that are not open", func() {

			ecdsaKey, err := crypto.RandomEcdsaKey()
			Expect(err).ShouldNot(HaveOccurred())

			multiAddr, err := identity.Address(ecdsaKey.Address()).MultiAddress()
			Expect(err).ShouldNot(HaveOccurred())

			client := newMockClient()
			server := newMockServer()
			streamer := NewStreamRecycler(NewStreamer(multiAddr.Address(), &client, &server))

			for i := 0; i < 100; i++ {
				ecdsaKeyOther, err := crypto.RandomEcdsaKey()
				Expect(err).ShouldNot(HaveOccurred())

				err = streamer.Close(identity.Address(ecdsaKeyOther.Address()))
				Expect(err).Should(Equal(ErrCloseOnClosedStream))
			}
		})
	})
})

type mockMessage []byte

func (message mockMessage) MarshalBinary() ([]byte, error) {
	return message, nil
}

func (message mockMessage) UnmarshalBinary(data []byte) error {
	copy(message, data)
	return nil
}

func (message mockMessage) IsMessage() {}

type mockStream struct {
	addr   identity.Address
	sends  *int64
	recvs  *int64
	closer func()
}

func (stream mockStream) Send(message Message) error {
	atomic.AddInt64(stream.sends, 1)
	return nil
}

func (stream mockStream) Recv(message Message) error {
	atomic.AddInt64(stream.recvs, 1)
	return nil
}

func (stream mockStream) Close() error {
	stream.closer()
	return nil
}

type mockClient struct {
	streamsMu      *sync.Mutex
	streams        map[int]mockStream
	streamsCounter int

	sends int64
	recvs int64
}

func newMockClient() mockClient {
	return mockClient{
		streamsMu: new(sync.Mutex),
		streams:   map[int]mockStream{},
	}
}

func (client *mockClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (CloseStream, error) {
	client.streamsMu.Lock()
	defer client.streamsMu.Unlock()

	i := client.streamsCounter
	client.streamsCounter++

	client.streams[i] = mockStream{
		addr:  multiAddr.Address(),
		sends: &client.sends,
		recvs: &client.recvs,
		closer: func() {
			client.streamsMu.Lock()
			defer client.streamsMu.Unlock()
			delete(client.streams, i)
		},
	}
	return client.streams[i], nil
}

type mockServer struct {
	streamsMu      *sync.Mutex
	streams        map[int]mockStream
	streamsCounter int

	sends int64
	recvs int64
}

func newMockServer() mockServer {
	return mockServer{
		streamsMu: new(sync.Mutex),
		streams:   map[int]mockStream{},
	}
}

func (server *mockServer) Listen(ctx context.Context, addr identity.Address) (CloseStream, error) {
	server.streamsMu.Lock()
	defer server.streamsMu.Unlock()

	i := server.streamsCounter
	server.streamsCounter++

	server.streams[i] = mockStream{
		addr:  addr,
		sends: &server.sends,
		recvs: &server.recvs,
		closer: func() {
			server.streamsMu.Lock()
			defer server.streamsMu.Unlock()
			delete(server.streams, i)
		},
	}
	return server.streams[i], nil
}
