package stream_test

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

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
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					_, err := streamers[i].Open(ctx, multiAddrs[j])
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
						ctx, cancel := context.WithCancel(context.Background())
						defer cancel()
						_, err := streamers[i].Open(ctx, multiAddrs[j])
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

			cancelConns := map[int]map[int]map[int]context.CancelFunc{}
			for i := 0; i < 128; i++ {
				cancelConns[i] = map[int]map[int]context.CancelFunc{}
				for j := 0; j < 128; j++ {
					if i == j {
						continue
					}
					cancelConns[i][j] = map[int]context.CancelFunc{}
					for k := 0; k < 4; k++ {
						ctx, cancel := context.WithCancel(context.Background())
						cancelConns[i][j][k] = cancel
						_, err := streamers[i].Open(ctx, multiAddrs[j])
						Expect(err).ShouldNot(HaveOccurred())
					}
				}
			}

			// Cancel all but the last connection
			for i := 0; i < 128; i++ {
				for j := 0; j < 128; j++ {
					if i == j {
						continue
					}
					for k := 0; k < 3; k++ {
						cancelConns[i][j][k]()
						Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(127))
					}
				}
			}

			// Cancel the last connection
			for i := 0; i < 128; i++ {
				for j := 0; j < 128; j++ {
					if i == j {
						continue
					}
					cancelConns[i][j][3]()
				}
			}

			// Expect shutdown of streams
			time.Sleep(time.Millisecond)
			for i := 0; i < 128; i++ {
				for j := 0; j < 128; j++ {
					if i == j {
						continue
					}
					Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(0))
				}
			}
		})
	})
})

type mockMessage []byte

func (message *mockMessage) MarshalBinary() ([]byte, error) {
	return *message, nil
}

func (message *mockMessage) UnmarshalBinary(data []byte) error {
	copy(*message, data)
	return nil
}

func (message mockMessage) IsMessage() {}

type mockStream struct {
	addr  identity.Address
	sends *int64
	recvs *int64
}

func (stream mockStream) Send(message Message) error {
	atomic.AddInt64(stream.sends, 1)
	return nil
}

func (stream mockStream) Recv(message Message) error {
	atomic.AddInt64(stream.recvs, 1)
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

func (client *mockClient) Connect(ctx context.Context, multiAddr identity.MultiAddress) (Stream, error) {
	client.streamsMu.Lock()
	defer client.streamsMu.Unlock()

	i := client.streamsCounter
	client.streamsCounter++

	client.streams[i] = mockStream{
		addr:  multiAddr.Address(),
		sends: &client.sends,
		recvs: &client.recvs,
	}

	go func() {
		<-ctx.Done()
		client.streamsMu.Lock()
		delete(client.streams, i)
		client.streamsMu.Unlock()
	}()

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

func (server *mockServer) Listen(ctx context.Context, addr identity.Address) (Stream, error) {
	server.streamsMu.Lock()
	defer server.streamsMu.Unlock()

	i := server.streamsCounter
	server.streamsCounter++

	server.streams[i] = mockStream{
		addr:  addr,
		sends: &server.sends,
		recvs: &server.recvs,
	}

	go func() {
		<-ctx.Done()
		server.streamsMu.Lock()
		delete(server.streams, i)
		server.streamsMu.Unlock()
	}()

	return server.streams[i], nil
}
