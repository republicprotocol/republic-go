package rpc

import (
	"context"
	"io"
	"time"

	"github.com/republicprotocol/go-dark-network"
	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"google.golang.org/grpc"
)

// SyncWithTarget using a new grpc.ClientConn to make Sync RPC to a target.
// This function returns two channels. The first is used to read chunks received
// in the synchronization. The second is used by the caller to quit when he no
// longer wants to receive dark.Chunk.
func SyncWithTarget(target, from identity.MultiAddress, timeout time.Duration) (chan do.Option, chan struct{}) {
	chunks := make(chan do.Option, 1)
	quit := make(chan struct{}, 1)

	go func() {
		defer close(chunks)
		conn, err := Dial(target, timeout)
		if err != nil {
			chunks <- do.Err(err)
			return
		}
		defer conn.Close()

		client := NewDarkNodeClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		stream, err := client.Sync(ctx, SerializeMultiAddress(from), grpc.FailFast(false))
		if err != nil {
			chunks <- do.Err(err)
			return
		}

		for {
			select {
			case _, ok := <-quit:
				if !ok {
					return
				}
			default:

			}

			chunk, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				chunks <- do.Err(err)
				continue
			}
			chunks <- do.Ok(DeserializeChunk(chunk))
		}
	}()
	return chunks, quit
}

// StartPropose
func StartPropose(target identity.MultiAddress,chunkRequest dark.ChunkRequest, timeout time.Duration) (dark.ChunkResponse, error) {
	conn, err := Dial(target, timeout)
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	multiAddresses, err := client.QueryCloserPeers(ctx, rpcQuery, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	return DeserializeMultiAddresses(multiAddresses)
}

// SendBid
func SendBid(target identity.MultiAddress, timeout time.Duration) (identity.MultiAddresses, error) {
	conn, err := Dial(target, timeout)
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	defer conn.Close()
	client := NewSwarmNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	rpcQuery := &Query{
		From:  SerializeMultiAddress(from),
		Query: SerializeAddress(query),
	}

	multiAddresses, err := client.QueryCloserPeers(ctx, rpcQuery, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	return DeserializeMultiAddresses(multiAddresses)
}