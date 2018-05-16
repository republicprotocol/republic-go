package cal

import "github.com/republicprotocol/republic-go/identity"

type Pod struct {
	Hash      [32]byte
	Darknodes []identity.Address
}

func (pod *Pod) Threshold() int {
	return (2 * (len(pod.Darknodes) + 1)) / 3
}

type Darkpool interface {
	Pods() ([]Pod, error)
	Pool(id identity.ID) (int, Pod, error)
}
