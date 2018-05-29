package cal

import (
	"bytes"
	"crypto/rsa"
	"errors"

	"github.com/republicprotocol/republic-go/identity"
)

// ErrPodNotFound is returned when a Pod cannot be found for a given
// identity.Address. This happens when an identity.Address is not registered in
// the current Epoch.
var ErrPodNotFound = errors.New("pod not found")

// ErrPublicKeyNotFound is returned when an rsa.PublicKey cannot be found for a
// given identity.Address. This happens when an identity.Address is not registered in
// the current Epoch.
var ErrPublicKeyNotFound = errors.New("public key not found")

// Darkpool is an interface for interacting with the Darkpool. Its core purpose
// is to expose to configuration of Darknodes into Pods for the different
// Epochs.
type Darkpool interface {

	// Darknodes registered in the Darkpool.
	Darknodes() (identity.Addresses, error)

	// Epoch returns the current Epoch which includes the Pod configuration.
	Epoch() (Epoch, error)

	// Pods returns the Pod configuration for the current Epoch.
	Pods() ([]Pod, error)

	// Pod returns the Pod that contains the given identity.Address in the
	// current Epoch. It returns ErrPodNotFound if the identity.Address is not
	// registered in the current Epoch.
	Pod(addr identity.Address) (Pod, error)

	// PublicKey returns the RSA public key of the Darknode registered with the
	// given identity.Address.
	PublicKey(addr identity.Address) (rsa.PublicKey, error)

	// IsRegistered returns true if the identity.Address is a current
	// registered Darknode. Otherwise, it returns false.
	IsRegistered(addr identity.Address) (bool, error)
}

// An Epoch represents the state of an epoch in the Darkpool. It stores the
// epoch hash, an ordered list of Pods for the epoch, and all Darknode
// identity.Addresses that are registered for the epoch.
type Epoch struct {
	Hash      [32]byte
	Pods      []Pod
	Darknodes []identity.Address
}

// Equal returns true if the hash of two Epochs is equal. Otherwise it returns
// false.
func (ξ *Epoch) Equal(arg *Epoch) bool {
	return bytes.Equal(ξ.Hash[:], arg.Hash[:])
}

// Pod returns the Pod that contains the given identity.Address in this Epoch.
// It returns ErrPodNotFound if the identity.Address is not registered in this
// Epoch.
func (ξ *Epoch) Pod(addr identity.Address) (Pod, error) {
	for _, pod := range ξ.Pods {
		for _, darknodeAddr := range pod.Darknodes {
			if addr == darknodeAddr {
				return pod, nil
			}
		}
	}
	return Pod{}, ErrPodNotFound
}

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

// EpochListener is an interface that can receive updates whenever the Epoch is
// changed in the Darkpool.
type EpochListener interface {
	OnChangeEpoch(ξ Epoch)
}
