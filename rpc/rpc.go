package rpc

import (
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/orderbook"
	"github.com/republicprotocol/republic-go/rpc/client"
	"github.com/republicprotocol/republic-go/rpc/dht"
	"github.com/republicprotocol/republic-go/rpc/relayer"
	"github.com/republicprotocol/republic-go/rpc/smpcer"
	"github.com/republicprotocol/republic-go/rpc/swarmer"
)

type RPC struct {
	crypter  crypto.Crypter
	dht      dht.DHT
	connPool client.ConnPool

	swarmerClient swarmer.Client
	swarmer       swarmer.Swarmer

	smpcerClient smpcer.Client
	smpcer       smpcer.Smpcer

	relayerClient relayer.Client
	relayer       relayer.Relayer

	onOpenOrder   func(order.Fragment) error
	onCancelOrder func(order.ID) error
}

func NewRPC(crypter crypto.Crypter, multiAddress identity.MultiAddress, orderbook *orderbook.Orderbook) *RPC {
	rpc := new(RPC)

	rpc.crypter = crypter
	rpc.dht = dht.NewDHT(multiAddress.Address(), 128)
	rpc.connPool = client.NewConnPool(256)

	rpc.swarmerClient = swarmer.NewClient(rpc.crypter, multiAddress, &rpc.dht, &rpc.connPool)
	rpc.swarmer = swarmer.NewSwarmer(&rpc.swarmerClient)

	rpc.smpcerClient = smpcer.NewClient(rpc.crypter, multiAddress, &rpc.connPool)
	rpc.smpcer = smpcer.NewSmpcer(rpc, &rpc.smpcerClient)

	rpc.relayerClient = relayer.NewClient(&rpc.dht, &rpc.connPool)
	rpc.relayer = relayer.NewRelayer(orderbook)

	return rpc
}

func (rpc *RPC) OpenOrder(signature []byte, orderFragment *smpcer.OrderFragment) error {
	panic("unimplemented")
}

func (rpc *RPC) CancelOrder(signature []byte, orderID []byte) error {
	panic("unimplemented")
}

func (rpc *RPC) OnOpenOrder(handler func()) {
	panic("unimplemented")
}

func (rpc *RPC) OnCancelOrder(signature []byte, orderID []byte) error {
	panic("unimplemented")
}

// SwarmerClient used by the RPC.
func (rpc *RPC) SwarmerClient() *swarmer.Client {
	return &rpc.swarmerClient
}

// Swarmer used by the RPC.
func (rpc *RPC) Swarmer() *swarmer.Swarmer {
	return &rpc.swarmer
}

// SmpcerClient used by the RPC.
func (rpc *RPC) SmpcerClient() *smpcer.Client {
	return &rpc.smpcerClient
}

// Smpcer used by the RPC.
func (rpc *RPC) Smpcer() *smpcer.Smpcer {
	return &rpc.smpcer
}

// RelayerClient used by the RPC.
func (rpc *RPC) RelayerClient() *relayer.Client {
	return &rpc.relayerClient
}

// Relayer used by the RPC.
func (rpc *RPC) Relayer() *relayer.Relayer {
	return &rpc.relayer
}
