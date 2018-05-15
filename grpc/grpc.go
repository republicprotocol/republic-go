package grpc

import (
	"context"

	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/grpc/client"
	"github.com/republicprotocol/republic-go/grpc/dht"
	"github.com/republicprotocol/republic-go/grpc/relayer"
	"github.com/republicprotocol/republic-go/grpc/smpcer"
	"github.com/republicprotocol/republic-go/grpc/status"
	"github.com/republicprotocol/republic-go/grpc/swarmer"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"google.golang.org/grpc"
)

type RPC struct {
	crypter  crypto.Crypter
	dht      dht.DHT
	connPool client.ConnPool

	relayerClient relayer.Client
	relayer       relayer.Relayer

	smpcerClient smpcer.Client
	smpcer       smpcer.Smpcer

	swarmerClient swarmer.Client
	swarmer       swarmer.Swarmer

	onOpenOrder   func([]byte, order.Fragment) error
	onCancelOrder func([]byte, order.ID) error
}

func NewRPC(crypter crypto.Crypter, multiAddress identity.MultiAddress, orderbook *orderbook.Orderbook) *RPC {
	rpc := new(RPC)

	rpc.crypter = crypter
	rpc.dht = dht.NewDHT(multiAddress.Address(), 128)
	rpc.connPool = client.NewConnPool(256)

	rpc.relayerClient = relayer.NewClient(rpc.crypter, &rpc.dht, &rpc.connPool)
	rpc.relayer = relayer.NewRelayer(&rpc.relayerClient, orderbook)

	rpc.smpcerClient = smpcer.NewClient(rpc.crypter, multiAddress, &rpc.connPool)
	rpc.smpcer = smpcer.NewSmpcer(&rpc.smpcerClient, rpc)

	rpc.swarmerClient = swarmer.NewClient(rpc.crypter, multiAddress, &rpc.dht, &rpc.connPool)
	rpc.swarmer = swarmer.NewSwarmer(&rpc.swarmerClient)

	return rpc
}

// OpenOrder implements the smpcer.Delegate interface.
func (rpc *RPC) OpenOrder(signature []byte, orderFragment order.Fragment) error {
	if rpc.onOpenOrder != nil {
		return rpc.onOpenOrder(signature, orderFragment)
	}
	return nil
}

// CancelOrder implements the smpcer.Delegate interface.
func (rpc *RPC) CancelOrder(signature []byte, orderID order.ID) error {
	if rpc.onCancelOrder != nil {
		return rpc.onCancelOrder(signature, orderID)
	}
	return nil
}

// OnOpenOrder call the delegate method.
func (rpc *RPC) OnOpenOrder(delegate func([]byte, order.Fragment) error) {
	rpc.onOpenOrder = delegate
}

// OnCancelOrder call the delegate method.
func (rpc *RPC) OnCancelOrder(delegate func([]byte, order.ID) error) {
	rpc.onCancelOrder = delegate
}

// RelayerClient used by the RPC.
func (rpc *RPC) RelayerClient() *relayer.Client {
	return &rpc.relayerClient
}

// Relayer used by the RPC.
func (rpc *RPC) Relayer() *relayer.Relayer {
	return &rpc.relayer
}

// SmpcerClient used by the RPC.
func (rpc *RPC) SmpcerClient() *smpcer.Client {
	return &rpc.smpcerClient
}

// Smpcer used by the RPC.
func (rpc *RPC) Smpcer() *smpcer.Smpcer {
	return &rpc.smpcer
}

// SwarmerClient used by the RPC.
func (rpc *RPC) SwarmerClient() *swarmer.Client {
	return &rpc.swarmerClient
}

// Swarmer used by the RPC.
func (rpc *RPC) Swarmer() *swarmer.Swarmer {
	return &rpc.swarmer
}

// Status will return the status needed by the falconry tool
func (rpc *RPC) Status(ctx context.Context, request *status.StatusRequest) (*status.StatusResponse, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		return &status.StatusResponse{
			Address:      string(rpc.dht.Address),
			Bootstrapped: rpc.swarmerClient.Bootstrapped(),
			Peers:        int64(len(rpc.dht.MultiAddresses())),
		}, nil
	}
}

// RegisterStatus will register the rpc with a grpc server
func (rpc *RPC) RegisterStatus(server *grpc.Server) {
	status.RegisterStatusServer(server, rpc)
}
