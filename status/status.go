package status

import (
	"sync"
)

// Writer will write the address
type Writer interface {
	WriteAddress(addr string) error
}

// Reader the address
type Reader interface {
	Address() (string, error)
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
	addr string
}

// NewProvider returns a new provider
func NewProvider() Provider {
	return &provider{
		mu: new(sync.Mutex),
	}
}

func (sp *provider) WriteAddress(addr string) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.addr = addr
	return nil
}

func (sp *provider) Address() (string, error) {
	return sp.addr, nil
}
