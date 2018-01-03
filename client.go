package swarm

import (
	"fmt"

	"github.com/republicprotocol/go-identity"
	"github.com/republicprotocol/go-swarm/rpc"
	"google.golang.org/grpc"
)

// NewNodeClient returns a new rpc.NodeClient that is connected to the given
// target identity.MultiAddress, or an error. It uses a background
// context.Context.
func NewNodeClient(target identity.MultiAddress) (rpc.NodeClient, error) {
	host, err := target.ValueForProtocol(identity.IP4Code)
	if err != nil {
		return nil, err
	}
	port, err := target.ValueForProtocol(identity.TCPCode)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		return nil, err
	}
	return rpc.NewNodeClient(conn), nil
}
