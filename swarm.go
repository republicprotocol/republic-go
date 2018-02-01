package rpc

import (
	"time"

	"github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// PingTarget using a new grpc.ClientConn to make a Ping RPC to a target
// identity.MultiAddress.
func PingTarget(to identity.MultiAddress, from identity.MultiAddress, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
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
func QueryCloserPeersFromTarget(to identity.MultiAddress, from identity.MultiAddress, query identity.Address, timeout time.Duration) (identity.MultiAddresses, error) {
	conn, err := Dial(to, timeout)
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

// QueryCloserPeersFromTarget using a new grpc.ClientConn to make a
// QueryCloserPeers RPC to a targetMultiAddress.
func QueryCloserPeersOnFrontierFromTarget(to identity.MultiAddress, from identity.MultiAddress, query identity.Address, timeout time.Duration) (identity.MultiAddresses, error) {
	conn, err := Dial(to, timeout)
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
	multiAddresses, err := client.QueryCloserPeersOnFrontier(ctx, rpcQuery, grpc.FailFast(false))
	if err != nil {
		return identity.MultiAddresses{}, err
	}
	return DeserializeMultiAddresses(multiAddresses)
}