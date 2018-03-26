package syncer

import "github.com/republicprotocol/republic-go/network/rpc"

type Broadcaster interface {
	Subscribe(id string, listener chan *rpc.SyncBlock) error
	Unsubscribe(id string)
}

