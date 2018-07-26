package registry

import (
	"bytes"
	"math/big"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/order"
)

// A Pod stores its hash, the combined hash of all Darknodes, and an ordered
// list of Darknode identity.Addresses.
type Pod struct {
	Position  int
	Hash      [32]byte
	Darknodes []identity.Address
}

// Size returns the number of Darknodes in the Pod.
func (pod *Pod) Size() int {
	return len(pod.Darknodes)
}

// Threshold returns the minimum number of Darknodes needed to run the order
// matching engine. It is the ceiling of 2/3rds of the Pod size.
func (pod *Pod) Threshold() int {
	return (2 * (len(pod.Darknodes) + 1)) / 3
}

type PodHeap ([]Pod)

func (heap PodHeap) PathOfOrder(orderID order.ID) PodPath {
	if len(heap) <= 1 {
		return PodPath(heap)
	}

	threshold := len(heap) / 2
	numberOfLeaves := len(heap) - threshold

	orderIDAsBigNum := big.NewInt(0).SetBytes(orderID[:])
	leafAsBigNum := big.NewInt(0).Mod(orderIDAsBigNum, big.NewInt(int64(numberOfLeaves)))
	leaf := int64(threshold) + leafAsBigNum.Int64()

	i := leaf
	path := PodPath{heap[i]}
	for i > 0 {
		i = (i - 1) / 2
		path = append(PodPath{heap[i]}, path...)
	}
	return path
}

// PodPath is an ordered list of Pods that is used to build hierarchies. The
// first Pod is the root and the last Pod is the leaf.
type PodPath ([]Pod)

// IndexOfPod returns the position of a Pod in the path. The returned boolean
// is true if the Pod is in the path, and false otherwise.
func (path PodPath) IndexOfPod(pod Pod) (int, bool) {
	for i := range path {
		if bytes.Equal(path[i].Hash[:], pod.Hash[:]) {
			return i, true
		}
	}
	return -1, false
}

// IndexOfAddress returns the position of the Pod in the path that contains a
// Darknode with the given address. The returned boolean is true if such a Pod
// is found, and false otherwise.
func (path PodPath) IndexOfAddress(addr identity.Address) (int, bool) {
	for i := range path {
		for j := range path[i].Darknodes {
			if path[i].Darknodes[j] == addr {
				return i, true
			}
		}
	}
	return -1, false
}

// Ancestor returns a path that is common to two paths. It begins at the root,
// and walks to the leaves, terminating as soon as a different pod is found.
// The returned path is guaranteed to be a common ancestor.
func (path PodPath) Ancestor(other PodPath) PodPath {
	common := PodPath{}
	for i := 0; i < len(path) && i < len(other); i++ {
		if !bytes.Equal(path[i].Hash[:], other[i].Hash[:]) {
			break
		}
		common = append(common, path[i])
	}
	return common
}
