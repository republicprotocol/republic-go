package dnr

import (
	"bytes"
	"errors"

	"github.com/republicprotocol/go-identity"
)

// DarkPool is a list of IDs
type DarkPool []identity.ID

// DarkOcean contains a list of dark pools
type DarkOcean struct {
	Pools []DarkPool
}

// FindDarkPool returns the pool containing a prticular ID
func (ocean DarkOcean) FindDarkPool(id identity.ID) (DarkPool, error) {

	for _, pool := range ocean.Pools {
		for _, node := range pool {
			if bytes.Compare(node, id) == 0 {
				return pool, nil
			}
		}
	}

	return nil, errors.New("Node is not a part of a pool")
}
