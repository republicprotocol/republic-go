package rpc

import (
	"time"

	"github.com/republicprotocol/go-atom"
	"github.com/republicprotocol/go-identity"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// SendAtomToTarget using a new grpc.ClientConn to make a SendAtom RPC to a
// target identity.MultiAddress.
func SendAtomToTarget(to identity.MultiAddress, a atom.Atom, timeout time.Duration) error {
	conn, err := Dial(to, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewTerminalNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err = client.SendAtom(ctx, SerializeAtom(a), grpc.FailFast(false))
	return err
}
