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
	ID   ComputeID
	Buy  *OrderFragment
	Sell *OrderFragment
	Out  *OrderFragment
}

func NewComputation(left *OrderFragment, right *OrderFragment) (*Computation, error) {
	if left.Buy == right.Buy {
		return nil, NewOrderComputationError(left.Buy)
	}
	com := &Computation{}
	if left.Buy != 0 {
		com.Buy = left
		com.Sell = right
	} else {
		com.Buy = right
		com.Sell = left
	}
	com.ID = ComputeID(crypto.Keccak256(com.Bytes()))
	return com, nil
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

type ComputationMatrix struct {
	orderFragments []*OrderFragment

	computationsGuard  *sync.Cond
	computationsLeft   int
	computations       []*Computation
	computationMarkers map[string]struct{}
}

func NewComputationMatrix() *ComputationMatrix {
	return &ComputationMatrix{
		orderFragments: []*OrderFragment{},

		computationsGuard:  sync.NewCond(new(sync.Mutex)),
		computationsLeft:   0,
		computations:       []*Computation{},
		computationMarkers: map[string]struct{}{},
	}
}

func (matrix *ComputationMatrix) FillComputations(orderFragment *OrderFragment) {
	matrix.computationsGuard.L.Lock()
	defer matrix.computationsGuard.L.Unlock()

	for _, rhs := range matrix.orderFragments {
		if orderFragment.ID.Equals(rhs.ID) {
			return
		}
	}

	for _, rhs := range matrix.orderFragments {
		matrix.computations = append(matrix.computations, NewComputation(orderFragment, rhs))
		matrix.computationsLeft++
	}

	matrix.orderFragments = append(matrix.orderFragments, orderFragment)
	if matrix.computationsLeft > 0 {
		matrix.computationsGuard.Signal()
	}
}

func (matrix *ComputationMatrix) WaitForComputations(max int) []*Computation {
	matrix.computationsGuard.L.Lock()
	defer matrix.computationsGuard.L.Unlock()
	for matrix.computationsLeft == 0 {
		matrix.computationsGuard.Wait()
	}

	computations := make([]*Computation, 0, max)
	for _, com := range matrix.computations {
		if _, ok := matrix.computationMarkers[string(com.ID)]; !ok {
			matrix.computationMarkers[string(com.ID)] = struct{}{}
			computations = append(computations, com)
		}
	}

	matrix.computationsLeft -= len(computations)
	if matrix.computationsLeft > 0 {
		matrix.computationsGuard.Signal()
	}
	return computations
}
