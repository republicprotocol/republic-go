package Syncer

type Broadcaster interface {
	Subscribe(id string, listener chan OrderStatusEvent)
	Unsubscribe(id string)
}

