package rpc

import (
	"io"
	"time"

	"github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PingTarget using a new grpc.ClientConn to make a Ping RPC to a target
// identity.MultiAddress.
func PingTarget(target identity.MultiAddress, from identity.MultiAddress, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewSwarmNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	_, err = client.Ping(ctx, SerializeMultiAddress(from), grpc.FailFast(false))
	return err
}

// QueryCloserPeersFromTarget using a new grpc.ClientConn to make a QueryCloserPeers
// RPC to a target identity.MultiAddress.
func QueryCloserPeersFromTarget(target identity.MultiAddress, from identity.MultiAddress, query identity.Address, timeout time.Duration) (identity.MultiAddresses, error) {
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

// QueryCloserPeersOnFrontierFromTarget using a new grpc.ClientConn to make a
// QueryCloserPeersOnFrontier RPC to a targetMultiAddress.
func QueryCloserPeersOnFrontierFromTarget(target identity.MultiAddress, from identity.MultiAddress, query identity.Address, timeout time.Duration) (identity.MultiAddresses, error) {
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
	stream, err := client.QueryCloserPeersOnFrontier(ctx, rpcQuery, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}

	multiAddresses := make(identity.MultiAddresses, 0)
	for {
		multiAddress, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return multiAddresses, err
		}
		deserializedMultiAddress, err := DeserializeMultiAddress(multiAddress)
		if err != nil {
			return multiAddresses, err
		}
		multiAddresses = append(multiAddresses, deserializedMultiAddress)
	}
	return multiAddresses, nil
}
