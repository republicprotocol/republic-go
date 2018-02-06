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
func SendAtomToTarget(target identity.MultiAddress, to identity.Address, from identity.MultiAddress, a atom.Atom, timeout time.Duration) error {
	conn, err := Dial(target, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	client := NewTerminalNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	serializedAtom := SerializeAtom(a)
	serializedAtom.To = SerializeAddress(to)
	serializedAtom.From = SerializeMultiAddress(from)
	_, err = client.SendAtom(ctx, serializedAtom, grpc.FailFast(false))
	return err
}
