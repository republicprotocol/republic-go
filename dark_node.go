package rpc

import (
	"context"
	"io"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"google.golang.org/grpc"
)

// SyncWithTarget using a new grpc.ClientConn to make Sync RPC to the target.
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
			chunks <- do.Ok(chunk)
		}
	}()
	return chunks, quit
}

// StartPropose using a new grpc.ClientConn to make a Propose RPC call
// to a target identity.MultiAddress.
func StartPropose(target identity.MultiAddress, chunkRequest *ChunkRequest, timeout time.Duration) (*ChunkResponse, error) {
	conn, err := Dial(target, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	chunkResponse, err := client.Propose(ctx, chunkRequest, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}
	return chunkResponse, nil
}

// AskToCompute using a new grpc.ClientConn to make a Compute RPC call
// to a target identity.MultiAddress.
func AskToCompute(target identity.MultiAddress, chunkRequest *ChunkRequest, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.Compute(ctx, chunkRequest, grpc.FailFast(false))
	if err != nil {
		return err
	}
	return nil
}
