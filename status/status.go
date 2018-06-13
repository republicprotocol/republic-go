package status

import (
	"sync"

	"github.com/republicprotocol/republic-go/identity"
)

// Writer will write the address
type Writer interface {
	WriteAddress(address string) error
	WriteMultiAddress(multiAddress identity.MultiAddress) error
}

// Reader the address
type Reader interface {
	Address() (string, error)
	MultiAddress() (identity.MultiAddress, error)
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
	mu           *sync.Mutex
	address      string
	multiAddress identity.MultiAddress
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
	sp.address = addr
	return nil
}

func (sp *provider) Address() (string, error) {
	return sp.address, nil
}

func (sp *provider) WriteMultiAddress(multiAddr identity.MultiAddress) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.multiAddress = multiAddr
	return nil
}

func (sp *provider) MultiAddress() (identity.MultiAddress, error) {
	return sp.multiAddress, nil
}
