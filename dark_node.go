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
// This function returns two channels. The first is used to read shards received
// in the synchronization. The second is used by the caller to quit when he no
// longer wants to receive dark.Chunk.
func SyncWithTarget(target, from identity.MultiAddress, timeout time.Duration) (chan do.Option, chan struct{}) {
	shards := make(chan do.Option, 1)
	quit := make(chan struct{}, 1)
	syncRequest := &SyncRequest{
		From: SerializeMultiAddress(from),
	}

	go func() {
		defer close(shards)
		conn, err := Dial(target, timeout)
		if err != nil {
			shards <- do.Err(err)
			return
		}
		defer conn.Close()

		client := NewDarkNodeClient(conn)
		ctx := context.Background()

		stream, err := client.Sync(ctx, syncRequest, grpc.FailFast(false))
		if err != nil {
			shards <- do.Err(err)
			return
		}

		for {
			select {
			case _, ok := <-quit:
				if ok {
					return
				}
			default:

			}

			shard, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				shards <- do.Err(err)
				continue
			}
			shards <- do.Ok(shard)
		}
	}()
	return shards, quit
}

// Logs ...
// This function returns two channels. The first is used to read shards received
// in the synchronization. The second is used by the caller to quit when he no
// longer wants to receive dark.Chunk.
func Logs(target identity.MultiAddress, timeout time.Duration) (chan do.Option, chan struct{}) {
	logEvents := make(chan do.Option, 1)
	quit := make(chan struct{}, 1)
	logRequest := &LogRequest{}

	go func() {
		defer close(logEvents)
		conn, err := Dial(target, timeout)
		if err != nil {
			logEvents <- do.Err(err)
			return
		}
		defer conn.Close()

		client := NewDarkNodeClient(conn)
		ctx := context.Background()

		stream, err := client.Logs(ctx, logRequest, grpc.FailFast(false))
		if err != nil {
			logEvents <- do.Err(err)
			return
		}

		for {
			select {
			case _, ok := <-quit:
				if ok {
					return
				}
			default:

			}

			shard, err := stream.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				logEvents <- do.Err(err)
				continue
			}
			logEvents <- do.Ok(shard)
		}
	}()
	return logEvents, quit
}
