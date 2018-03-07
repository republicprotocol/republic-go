package dnr

import (
	"bytes"
	"errors"

	"github.com/republicprotocol/go-do"
	"github.com/republicprotocol/go-identity"
)

// DarkPool is a list of node multiaddresses
type DarkPool struct {
	do.GuardedObject

	Nodes identity.MultiAddresses
}

// IDDarkPool is a list of node ids
type IDDarkPool []identity.ID

// DarkOcean contains a list of dark pools
type DarkOcean struct {
	Pools []IDDarkPool
}

// FindDarkPool returns the pool containing a prticular ID
func (ocean DarkOcean) FindDarkPool(id identity.ID) (IDDarkPool, error) {

	for _, pool := range ocean.Pools {
		for _, node := range pool {
			if bytes.Compare(node, id) == 0 {
				return pool, nil
			}
		}
	}

	return nil, errors.New("Node is not a part of a pool")
}
