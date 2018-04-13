package rpc

import (
	"fmt"
	"io"
	"runtime"
	"sync"

	"github.com/republicprotocol/republic-go/identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// A Client is used to create and manage a gRPC connection. It provides methods
// for all RPCs and handles all timeouts and retries.
type Client struct {
	SmpcClient
	SwarmClient
	RelayClient
	SyncerClient

	Connection *grpc.ClientConn
	To         *MultiAddress
	From       *MultiAddress
	Options    ClientOptions
}

// NewClient returns a Client that is connected to the given MultiAddress and
// will always identify itself from the given MultiAddress. The connection will
// be closed when the Client is garbage collected.
func NewClient(ctx context.Context, to, from identity.MultiAddress) (*Client, error) {
	host, err := to.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := to.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Options: DefaultClientOptions(),
		To:      MarshalMultiAddress(&to),
		From:    MarshalMultiAddress(&from),
	}

	connection, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", host, port), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client.Connection = connection
	runtime.SetFinalizer(client, func(client *Client) {
		client.Connection.Close()
	})

	client.SmpcClient = NewSmpcClient(client.Connection)
	client.SwarmClient = NewSwarmClient(client.Connection)
	client.RelayClient = NewRelayClient(client.Connection)
	client.SyncerClient = NewSyncerClient(client.Connection)

	return client, nil
}

func (client *Client) Ping(ctx context.Context) error {
	_, err := client.SwarmClient.Ping(ctx, client.To, grpc.FailFast(false))
	return err
}

func (client *Client) QueryPeers(ctx context.Context, target *Address) (<-chan *MultiAddress, <-chan error) {
	multiAddressCh := make(chan *MultiAddress)
	errCh := make(chan error)

	go func() {
		defer close(multiAddressCh)
		defer close(errCh)

		stream, err := client.SwarmClient.QueryPeers(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			errCh <- err
			return
		}

		for {
			multiAddress, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case multiAddressCh <- multiAddress:
			}
		}
	}()

	return multiAddressCh, errCh
}

func (client *Client) QueryPeersDeep(ctx context.Context, target *Address) (<-chan *MultiAddress, <-chan error) {
	multiAddressCh := make(chan *MultiAddress)
	errCh := make(chan error)

	go func() {
		defer close(multiAddressCh)
		defer close(errCh)

		stream, err := client.SwarmClient.QueryPeersDeep(ctx, &Query{
			From:   client.From,
			Target: target,
		}, grpc.FailFast(false))
		if err != nil {
			errCh <- err
			return
		}

		for {
			multiAddress, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case multiAddressCh <- multiAddress:
			}
		}
	}()

	return multiAddressCh, errCh
}

func (client *Client) Sync(ctx context.Context) (<-chan *SyncBlock, <-chan error) {
	syncBlockCh := make(chan *SyncBlock)
	errCh := make(chan error)

	go func() {
		defer close(syncBlockCh)
		defer close(errCh)

		stream, err := client.SyncerClient.Sync(ctx, &SyncRequest{
			From: client.From,
		}, grpc.FailFast(false))
		if err != nil {
			errCh <- err
			return
		}

		for {
			syncBlock, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					errCh <- err
				}
				return
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case syncBlockCh <- syncBlock:
			}
		}
	}()

	return syncBlockCh, errCh
}

func (client *Client) SignOrderFragment(ctx context.Context, orderFragmentSignature *OrderFragmentId) (*OrderFragmentId, error) {
	return client.RelayClient.SignOrderFragment(ctx, &OrderFragmentId{
		Signature:       orderFragmentSignature.Signature,
		OrderFragmentId: orderFragmentSignature.OrderFragmentId,
	}, grpc.FailFast(false))
}

func (client *Client) OpenOrder(ctx context.Context, openOrderRequest *OpenOrderRequest) error {
	_, err := client.RelayClient.OpenOrder(ctx, &OpenOrderRequest{
		From:          client.From,
		OrderFragment: openOrderRequest.OrderFragment,
	}, grpc.FailFast(false))

	return err
}

func (client *Client) CancelOrder(ctx context.Context, cancelOrderRequest *CancelOrderRequest) error {
	_, err := client.RelayClient.CancelOrder(ctx, &CancelOrderRequest{
		From:            client.From,
		OrderFragmentId: cancelOrderRequest.OrderFragmentId,
	}, grpc.FailFast(false))

	return err
}

func (client *Client) Compute(ctx context.Context, messageChIn <-chan *SmpcMessage) (<-chan *SmpcMessage, <-chan error) {
	messageCh := make(chan *SmpcMessage, 1)
	errCh := make(chan error, 1)

	stream, err := client.SmpcClient.Compute(ctx, grpc.FailFast(false))
	if err != nil {
		errCh <- err
		return messageCh, errCh
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		defer close(messageCh)

		for {
			message, err := stream.Recv()
			if err != nil {
				errCh <- err
				return
			}
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
				return
			case messageCh <- message:
			}
		}
	}()
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case message, ok := <-messageChIn:
				if !ok {
					return
				}
				if err := stream.Send(message); err != nil {
					s, _ := status.FromError(err)
					if s.Code() != codes.Canceled && s.Code() != codes.DeadlineExceeded {
						errCh <- err
					}
					return
				}
			}
		}
	}()
	go func() {
		defer close(errCh)
		wg.Wait()
	}()

	return messageCh, errCh
}
