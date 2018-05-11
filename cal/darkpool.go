package cal

import "github.com/republicprotocol/republic-go/identity"

type Pod struct {
	Hash      [32]byte
	Darknodes []identity.Address
}

type Darkpool interface {
	Pods() ([]Pod, error)
}
