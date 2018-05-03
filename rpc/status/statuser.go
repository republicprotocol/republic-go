package status

import (
	"context"

	"github.com/republicprotocol/go-identity"
)

type Statuser struct {
	address      identity.Address
	bootstrapped bool
	registered   bool
	connectPeers int
}

func NewStatuser(address identity.Address) Statuser {
	return Statuser{
		address:      address,
		bootstrapped: false,
		registered:   false,
		connectPeers: 0,
	}
}

func (statuser *Statuser) AfterBootstrap() {
	statuser.bootstrapped = true
}

func (statuser *Statuser) Register() {
	statuser.registered = true
}

func (statuser *Statuser) Deregister() {
	statuser.registered = false
}

func (statuser *Statuser) SetConnectNodes(number int) {
	statuser.connectPeers = number
}

func (statuser *Statuser) Status(ctx context.Context, request *StatusRequest) (*StatusResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &StatusResponse{
			Address:      statuser.address.String(),
			Bootstrapped: statuser.bootstrapped,
			Registered:   statuser.registered,
			Peers:        int64(statuser.connectPeers),
		}, nil
	}
}
