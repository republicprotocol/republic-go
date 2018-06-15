package stream_test

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/stream"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/testutils"
)

const numberOfNodes = 32

var _ = Describe("Streaming", func() {

	Context("when using streamers", func() {

		It("should abstract connecting to servers and listening for client", func() {

			multiAddrs := [numberOfNodes]identity.MultiAddress{}
			clients := [numberOfNodes]mockClient{}
			servers := [numberOfNodes]mockServer{}
			streamers := [numberOfNodes]Streamer{}

			var err error

			for i := 0; i < numberOfNodes; i++ {
				multiAddrs[i], err = testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i])
			}

			for i := 0; i < numberOfNodes; i++ {
				for j := 0; j < numberOfNodes; j++ {
					if i == j {
						continue
					}
					ctx, cancel := context.WithCancel(context.Background())
					defer cancel()
					_, err := streamers[i].Open(ctx, multiAddrs[j])
					Expect(err).ShouldNot(HaveOccurred())
				}
			}

			for i := 0; i < numberOfNodes; i++ {
				Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(numberOfNodes - 1))
			}
		})

	})

	Context("when using stream recyclers", func() {

		It("should recycle streams for multiple connections", func() {

			multiAddrs := [numberOfNodes]identity.MultiAddress{}
			clients := [numberOfNodes]mockClient{}
			servers := [numberOfNodes]mockServer{}
			streamers := [numberOfNodes]Streamer{}

			var err error

			for i := 0; i < numberOfNodes; i++ {
				multiAddrs[i], err = testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			for conns := 0; conns < 4; conns++ {
				for i := 0; i < numberOfNodes; i++ {
					for j := 0; j < numberOfNodes; j++ {
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

			for i := 0; i < numberOfNodes; i++ {
				Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(numberOfNodes - 1))
			}

		})

		It("should not close streams until all references have closed", func() {

			multiAddrs := [numberOfNodes]identity.MultiAddress{}
			clients := [numberOfNodes]mockClient{}
			servers := [numberOfNodes]mockServer{}
			streamers := [numberOfNodes]Streamer{}

			var err error

			for i := 0; i < numberOfNodes; i++ {
				multiAddrs[i], err = testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			cancelConns := map[int]map[int]map[int]context.CancelFunc{}
			for i := 0; i < numberOfNodes; i++ {
				cancelConns[i] = map[int]map[int]context.CancelFunc{}
				for j := 0; j < numberOfNodes; j++ {
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
			for i := 0; i < numberOfNodes; i++ {
				clients[i].streamsMu.Lock()
				servers[i].streamsMu.Lock()
				for j := 0; j < numberOfNodes; j++ {
					if i == j {
						continue
					}
					for k := 0; k < 3; k++ {
						cancelConns[i][j][k]()
						Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(numberOfNodes - 1))
					}
				}
				clients[i].streamsMu.Unlock()
				servers[i].streamsMu.Unlock()
			}

			// Cancel the last connection
			for i := 0; i < numberOfNodes; i++ {
				for j := 0; j < numberOfNodes; j++ {
					if i == j {
						continue
					}
					cancelConns[i][j][3]()
				}
			}

			// Expect shutdown of streams
			time.Sleep(time.Second)
			for i := 0; i < numberOfNodes; i++ {
				clients[i].streamsMu.Lock()
				servers[i].streamsMu.Lock()
				for j := 0; j < numberOfNodes; j++ {
					if i == j {
						continue
					}
					Expect(len(clients[i].streams) + len(servers[i].streams)).Should(Equal(0))
				}
				clients[i].streamsMu.Unlock()
				servers[i].streamsMu.Unlock()
			}
		})
	})

	Context("regression test", func() {
		It("should only have one stream for the client and server when trying connecting sequentially", func() {
			multiAddrs := [numberOfNodes]identity.MultiAddress{}
			clients := [numberOfNodes]mockClient{}
			servers := [numberOfNodes]mockServer{}
			streamers := [numberOfNodes]Streamer{}

			var err error

			for i := 0; i < numberOfNodes; i++ {
				multiAddrs[i], err = testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			for i := 0; i < 1000; i++ {
				stream, err := streamers[0].Open(context.Background(), multiAddrs[1])
				Expect(err).ShouldNot(HaveOccurred())
				Expect(stream).ShouldNot(BeNil())
				stream.Send(&mockMessage{})
			}
			if multiAddrs[0].Address() < multiAddrs[1].Address() {
				Expect(clients[0].streamsCounter).Should(Equal(1))
			} else {
				Expect(servers[0].streamsCounter).Should(Equal(1))
			}
		})

		It("should only have one stream for the client and server when trying connecting concurrently", func() {
			multiAddrs := [numberOfNodes]identity.MultiAddress{}
			clients := [numberOfNodes]mockClient{}
			servers := [numberOfNodes]mockServer{}
			streamers := [numberOfNodes]Streamer{}

			var err error

			for i := 0; i < numberOfNodes; i++ {
				multiAddrs[i], err = testutils.RandomMultiAddress()
				Expect(err).ShouldNot(HaveOccurred())

				clients[i] = newMockClient()
				servers[i] = newMockServer()
				streamers[i] = NewStreamRecycler(NewStreamer(multiAddrs[i].Address(), &clients[i], &servers[i]))
			}

			wg := new(sync.WaitGroup)
			for i := 0; i < 1000; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					stream, err := streamers[0].Open(context.Background(), multiAddrs[1])
					Expect(err).ShouldNot(HaveOccurred())
					Expect(stream).ShouldNot(BeNil())
					stream.Send(&mockMessage{})
				}()
			}
			wg.Wait()
			if multiAddrs[0].Address() < multiAddrs[1].Address() {
				Expect(clients[0].streamsCounter).Should(Equal(1))
			} else {
				Expect(servers[0].streamsCounter).Should(Equal(1))
			}
		})
	})
})

type mockMessage []byte

func (message *mockMessage) MarshalBinary() ([]byte, error) {
	return *message, nil
}

func (message *mockMessage) UnmarshalBinary(data []byte) error {
	*message = data
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
