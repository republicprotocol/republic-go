package rpc

import (
	"context"
	"io"
	"time"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-order-compute"
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
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		stream, err := client.Sync(ctx, syncRequest, grpc.FailFast(false))
		if err != nil {
			shards <- do.Err(err)
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

// StartElectShard uses a new grpc.ClientConn to make a ElectShard RPC call
// to a target identity.MultiAddress.
func StartElectShard(target, from identity.MultiAddress, shard compute.Shard, timeout time.Duration) (*Shard, error) {
	conn, err := Dial(target, timeout)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &ElectShardRequest{
		From:  SerializeMultiAddress(from),
		Shard: SerializeShard(shard),
	}

	rsp, err := client.ElectShard(ctx, request, grpc.FailFast(false))
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

// AskToComputeShard using a new grpc.ClientConn to make a Compute RPC call
// to a target identity.MultiAddress.
func AskToComputeShard(target, from identity.MultiAddress, shard compute.Shard, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &ComputeShardRequest{
		From:  SerializeMultiAddress(from),
		Shard: SerializeShard(shard),
	}

	_, err = client.ComputeShard(ctx, request, grpc.FailFast(false))
	if err != nil {
		return err
	}
	return nil
}

// FinalizeShard using a new grpc.ClientConn to make a Compute RPC call
// to a target identity.MultiAddress.
func FinalizeShard(target, from identity.MultiAddress, shard compute.DeltaShard, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewDarkNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	request := &FinalizeShardRequest{
		From:  SerializeMultiAddress(from),
		Shard: SerializeFinalShard(shard),
	}

	_, err = client.FinalizeShard(ctx, request, grpc.FailFast(false))
	if err != nil {
		return err
	}
	return nil
}
