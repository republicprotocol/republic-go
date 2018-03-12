package dnr

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MockDarkNodeRegistrar implemented to DarkNodeRegistrar interface for
// testing.
type MockDarkNodeRegistrar struct {
	hash      [32]byte
	timestamp *big.Int
	nodeIDs   map[string]bool
}

func NewMockDarkNodeRegistrar(nodeIDs [][]byte) DarkNodeRegistrar {
	darkNodeRegistrar := new(MockDarkNodeRegistrar)
	darkNodeRegistrar.hash = [32]byte{1}
	darkNodeRegistrar.timestamp = big.NewInt(1)
	darkNodeRegistrar.nodeIDs = map[string]bool{}
	for _, nodeID := range nodeIDs {
		darkNodeRegistrar.nodeIDs[string(nodeID)] = true
	}
	return darkNodeRegistrar
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) Register(nodeID []byte, publicKey []byte) (*types.Transaction, error) {
	darkNodeRegistrar.nodeIDs[string(nodeID)] = true
	return nil, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) Deregister(nodeID []byte) (*types.Transaction, error) {
	delete(darkNodeRegistrar.nodeIDs, string(nodeID))
	return nil, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) GetBond(nodeID []byte) (*big.Int, error) {
	if _, ok := darkNodeRegistrar.nodeIDs[string(nodeID)]; ok {
		return big.NewInt(86000), nil
	}
	return big.NewInt(0), nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) IsDarkNodeRegistered(nodeID []byte) (bool, error) {
	_, ok := darkNodeRegistrar.nodeIDs[string(nodeID)]
	return ok, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) IsDarkNodePendingRegistration(nodeID []byte) (bool, error) {
	return false, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) CurrentEpoch() (Epoch, error) {
	return Epoch{
		Blockhash: darkNodeRegistrar.hash,
		Timestamp: darkNodeRegistrar.timestamp,
	}, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) Epoch() (*types.Transaction, error) {
	darkNodeRegistrar.timestamp.Add(darkNodeRegistrar.timestamp, big.NewInt(1))
	darkNodeRegistrar.hash[0]++
	return nil, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) GetCommitment(nodeID []byte) ([32]byte, error) {
	return [32]byte{}, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) GetOwner(nodeID []byte) (common.Address, error) {
	return common.Address{}, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) GetPublicKey(nodeID []byte) ([]byte, error) {
	return []byte{}, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) GetAllNodes() ([][]byte, error) {
	allNodes := make([][]byte, 0, len(darkNodeRegistrar.nodeIDs))
	for nodeID := range darkNodeRegistrar.nodeIDs {
		allNodes = append(allNodes, []byte(nodeID))
	}
	return allNodes, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumBond() (*big.Int, error) {
	return big.NewInt(86000), nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(1), nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) Refund(nodeID []byte) (*types.Transaction, error) {
	return nil, nil
}

func (darkNodeRegistrar *MockDarkNodeRegistrar) WaitUntilRegistration(nodeID []byte) error {
	for {
		if registered, err := darkNodeRegistrar.IsDarkNodeRegistered(nodeID); err == nil && registered {
			return nil
		}
		time.Sleep(time.Minute)
	}
}
