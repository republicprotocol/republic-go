package transact

import (
	"net/rpc"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Client interface {
	RPC() *rpc.Client
	Ethereum() *ethclient.Client
}

type client struct {
	rpc      *rpc.Client
	ethereum *ethclient.Client
}

func (client *client) RPC() *rpc.Client {
	return client.rpc
}

func (client *client) Ethereum() *ethclient.Client {
	return client.ethereum
}
