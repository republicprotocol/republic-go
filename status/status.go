package status

import (
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// Writer will write the address
type Writer interface {
	WriteAddress(addr identity.Address) error
}

// Reader the address
type Reader interface {
	Address() (identity.Address, error)
}

/*

Basic information
* address
* multiaddress
* ethereum address
* connected peers
* basic funds amounts, balance, e.g. eth
* fees it's earned
* and register or deregister a dark node (deregister is disabled metamask address is not the owner)

*/

type Provider interface {
	Writer
	Reader
}

type provider struct {
	mu   *sync.Mutex
	addr identity.Address
}

// NewProvider returns a new provider
func NewProvider() Provider {
	return &provider{
		mu: new(sync.Mutex),
	}
}

func (sp *provider) WriteAddress(addr identity.Address) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.addr = addr
	return nil
}

func (sp *provider) Address() (identity.Address, error) {
	return sp.addr, nil
}
