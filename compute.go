package compute

import (
	"bytes"
	"encoding/binary"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

type ComputeID []byte

type Computation struct {
	ID    ComputeID
	Left  *OrderFragment
	Right *OrderFragment
	Out   *OrderFragment
}

func NewComputation(left *OrderFragment, right *OrderFragment) *Computation {
	com := &Computation{
		Left:  left,
		Right: right,
	}
	com.ID = ComputeID(crypto.Keccak256(com.Bytes()))
	return com
}

func (com *Computation) Add(prime *big.Int) (*OrderFragment, error) {
	out, err := com.Left.Add(com.Right, prime)
	if err != nil {
		return nil, err
	}
	com.Out = out
	return com.Out, nil
}

func (com *Computation) Sub(prime *big.Int) (*OrderFragment, error) {
	out, err := com.Left.Sub(com.Right, prime)
	if err != nil {
		return nil, err
	}
	com.Out = out
	return com.Out, nil
}

// Bytes returns an Order serialized into a bytes.
func (com *Computation) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, com.Left.ID)
	binary.Write(buf, binary.LittleEndian, com.Right.ID)
	return buf.Bytes()
}

type State struct {
	orderFragmentsMu *sync.Mutex
	orderFragments   []*OrderFragment

	computationsMu *sync.RWMutex
	computations   []*Computation

	computationStatesMu *sync.Mutex
	computationStates   map[string]struct{}
}

func (state *State) AddOrderFragment(orderFragment *OrderFragment) {
	state.orderFragmentsMu.Lock()
	defer state.orderFragmentsMu.Unlock()

	// Do nothing if the OrderFragment has been added previously.
	for _, rhs := range state.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	// Create all of the possible Computations.
	state.computationsMu.Lock()
	defer state.computationsMu.Unlock()

	for _, rhs := range state.orderFragments {
		state.computations = append(state.computations, NewComputation(orderFragment, rhs))
	}
	state.orderFragments = append(state.orderFragments, orderFragment)
}

func (state *State) CollectComputations(max int) []*Computation {
	state.computationsMu.RLock()
	defer state.computationsMu.RUnlock()
	state.computationStatesMu.Lock()
	defer state.computationStatesMu.Unlock()

	coms := make([]*Computation, 0, max)
	for _, com := range state.computations {
		if _, ok := state.computationStates[string(com.ID)]; !ok {
			state.computationStates[string(com.ID)] = struct{}{}
			coms = append(coms, com)
		}
	}
	return coms
}
