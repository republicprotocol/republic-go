package dnr

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/republicprotocol/republic-go/stackint"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	do "github.com/republicprotocol/go-do"
)

// MockDarkNodeRegistrar implemented to DarkNodeRegistrar interface for
// testing.
type MockDarkNodeRegistrar struct {
	do.GuardedObject

	registered   [][]byte
	toRegister   [][]byte
	toDeregister [][]byte
	epoch        Epoch
}

// NewMockDarkNodeRegistrar returns a new MockDarkNodeRegistrar
func NewMockDarkNodeRegistrar() (DarkNodeRegistrar, error) {
	mockDnr := new(MockDarkNodeRegistrar)
	mockDnr.GuardedObject = do.NewGuardedObject()
	mockDnr.registered = make([][]byte, 0)
	mockDnr.toRegister = make([][]byte, 0)
	mockDnr.toDeregister = make([][]byte, 0)
	var b32 [32]byte
	_, err := rand.Read(b32[:])
	if err != nil {
		return nil, err
	}
	now := stackint.FromUint(uint(time.Now().Unix()))
	mockDnr.epoch = Epoch{
		Blockhash: b32,
		Timestamp: &now,
	}
	return mockDnr, nil
}

// Register a new dark node
func (mockDnr *MockDarkNodeRegistrar) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
	mockDnr.Enter(nil)
	defer mockDnr.Exit()
	isRegistered, _ := mockDnr.isDarkNodeRegistered(darkNodeID)
	isPending, _ := mockDnr.isDarkNodePendingRegistration(darkNodeID)
	if isRegistered || isPending {
		return nil, errors.New("Must not be registered to register")
	}
	mockDnr.toRegister = append(mockDnr.toRegister, darkNodeID)
	return nil, nil
}

// Deregister an existing dark node
func (mockDnr *MockDarkNodeRegistrar) Deregister(darkNodeID []byte) (*types.Transaction, error) {
	mockDnr.Enter(nil)
	defer mockDnr.Exit()
	for i, id := range mockDnr.toRegister {
		if string(darkNodeID) == string(id) {
			mockDnr.toDeregister[i] = mockDnr.toDeregister[len(mockDnr.toDeregister)-1]
			mockDnr.toDeregister = mockDnr.toDeregister[:len(mockDnr.toDeregister)-1]
			return nil, nil
		}
	}
	if isRegistered, _ := mockDnr.isDarkNodeRegistered(darkNodeID); !isRegistered {
		return nil, errors.New("Must be registered to deregister")
	}
	mockDnr.toDeregister = append(mockDnr.toRegister, darkNodeID)
	return nil, nil
}

// GetBond retrieves the bond of an existing dark node
func (mockDnr *MockDarkNodeRegistrar) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return stackint.FromUint(86000), nil
}

// IsDarkNodeRegistered returns true if the node is registered
func (mockDnr *MockDarkNodeRegistrar) IsDarkNodeRegistered(darkNodeID []byte) (bool, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return mockDnr.isDarkNodeRegistered(darkNodeID)
}

func (mockDnr *MockDarkNodeRegistrar) isDarkNodeRegistered(darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.registered {
		if string(darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}

// IsDarkNodePendingRegistration returns true if the node will become registered at the next epoch
func (mockDnr *MockDarkNodeRegistrar) IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return mockDnr.isDarkNodePendingRegistration(darkNodeID)
}

func (mockDnr *MockDarkNodeRegistrar) isDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.toRegister {
		if string(darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}

// CurrentEpoch returns the current epoch
func (mockDnr *MockDarkNodeRegistrar) CurrentEpoch() (Epoch, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return mockDnr.epoch, nil
}

// Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
func (mockDnr *MockDarkNodeRegistrar) Epoch() (*types.Transaction, error) {

	mockDnr.Enter(nil)
	defer mockDnr.Exit()

	var b32 [32]byte

	_, err := rand.Read(b32[:])
	if err != nil {
		return nil, err
	}

	now := stackint.FromUint(uint(time.Now().Unix()))
	mockDnr.epoch = Epoch{
		Blockhash: b32,
		Timestamp: &now,
	}

	// Remove toRegister nodes
	for _, deregNode := range mockDnr.toDeregister {
		for i := 0; i < len(mockDnr.registered); i++ {
			if string(mockDnr.registered[i]) == string(deregNode) {
				mockDnr.registered[i] = mockDnr.registered[len(mockDnr.registered)-1]
				mockDnr.registered = mockDnr.registered[:len(mockDnr.registered)-1]
				break
			}
		}
	}

	mockDnr.registered = append(mockDnr.registered, mockDnr.toRegister...)

	mockDnr.toDeregister = make([][]byte, 0)
	mockDnr.toRegister = make([][]byte, 0)

	return nil, nil
}

// GetCommitment gets the signed commitment (not implemented)
func (mockDnr *MockDarkNodeRegistrar) GetCommitment(darkNodeID []byte) ([32]byte, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return [32]byte{}, nil
}

// GetOwner gets the owner of the given dark node (not implemented)
func (mockDnr *MockDarkNodeRegistrar) GetOwner(darkNodeID []byte) (common.Address, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return [20]byte{}, nil
}

// GetPublicKey gets the public key of the goven dark node (not implemented)
func (mockDnr *MockDarkNodeRegistrar) GetPublicKey(darkNodeID []byte) ([]byte, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return nil, nil
}

// GetAllNodes gets all dark nodes
func (mockDnr *MockDarkNodeRegistrar) GetAllNodes() ([][]byte, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return mockDnr.registered, nil
}

// MinimumBond gets the minimum viable bond amount (hard-coded to 86000)
func (mockDnr *MockDarkNodeRegistrar) MinimumBond() (stackint.Int1024, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return stackint.FromUint(86000), nil
}

// MinimumEpochInterval gets the minimum epoch interval (hard-coded to 0)
func (mockDnr *MockDarkNodeRegistrar) MinimumEpochInterval() (stackint.Int1024, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return stackint.FromUint(0), nil
}

// Refund refunds the bond of an unregistered miner (not implemented)
func (mockDnr *MockDarkNodeRegistrar) Refund(darkNodeID []byte) (*types.Transaction, error) {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	return nil, nil
}

// WaitUntilRegistration waits until the registration is successful
func (mockDnr *MockDarkNodeRegistrar) WaitUntilRegistration(darkNodeID []byte) error {
	mockDnr.EnterReadOnly(nil)
	defer mockDnr.ExitReadOnly()
	for {
		isRegistered, err := mockDnr.isDarkNodeRegistered(darkNodeID)
		if err != nil {
			return err
		}
		if isRegistered {
			break
		}
		time.Sleep(time.Minute)
	}
	return nil
}

// type MockDarkNodeRegistrar struct {
// 	do.GuardedObject

// 	hash      [32]byte
// 	timestamp *stackint.Int1024
// 	nodeIDs   map[string]bool
// }

// func NewMockDarkNodeRegistrar(nodeIDs [][]byte) DarkNodeRegistrar {
// 	darkNodeRegistrar := new(MockDarkNodeRegistrar)
// 	darkNodeRegistrar.GuardedObject = do.NewGuardedObject()
// 	darkNodeRegistrar.hash = [32]byte{1}
// 	darkNodeRegistrar.timestamp = big.NewInt(1)
// 	darkNodeRegistrar.nodeIDs = map[string]bool{}
// 	for _, nodeID := range nodeIDs {
// 		darkNodeRegistrar.nodeIDs[string(nodeID)] = true
// 	}
// 	return darkNodeRegistrar
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) Register(nodeID []byte, publicKey []byte) (*types.Transaction, error) {
// 	darkNodeRegistrar.Enter(nil)
// 	defer darkNodeRegistrar.Exit()
// 	darkNodeRegistrar.nodeIDs[string(nodeID)] = true
// 	return nil, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) Deregister(nodeID []byte) (*types.Transaction, error) {
// 	darkNodeRegistrar.Enter(nil)
// 	defer darkNodeRegistrar.Exit()
// 	delete(darkNodeRegistrar.nodeIDs, string(nodeID))
// 	return nil, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetBond(nodeID []byte) (*stackint.Int1024, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	if _, ok := darkNodeRegistrar.nodeIDs[string(nodeID)]; ok {
// 		return big.NewInt(86000), nil
// 	}
// 	return big.NewInt(0), nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) IsDarkNodeRegistered(nodeID []byte) (bool, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	_, ok := darkNodeRegistrar.nodeIDs[string(nodeID)]
// 	return ok, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) IsDarkNodePendingRegistration(nodeID []byte) (bool, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return false, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) CurrentEpoch() (Epoch, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return Epoch{
// 		Blockhash: darkNodeRegistrar.hash,
// 		Timestamp: darkNodeRegistrar.timestamp,
// 	}, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) Epoch() (*types.Transaction, error) {
// 	darkNodeRegistrar.Enter(nil)
// 	defer darkNodeRegistrar.Exit()
// 	darkNodeRegistrar.timestamp.Add(darkNodeRegistrar.timestamp, big.NewInt(1))
// 	darkNodeRegistrar.hash[0]++
// 	return nil, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetCommitment(nodeID []byte) ([32]byte, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return [32]byte{}, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetOwner(nodeID []byte) (common.Address, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return common.Address{}, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetPublicKey(nodeID []byte) ([]byte, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return []byte{}, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetAllNodes() ([][]byte, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	allNodes := make([][]byte, 0, len(darkNodeRegistrar.nodeIDs))
// 	for nodeID := range darkNodeRegistrar.nodeIDs {
// 		allNodes = append(allNodes, []byte(nodeID))
// 	}
// 	return allNodes, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumBond() (*stackint.Int1024, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return big.NewInt(86000), nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumEpochInterval() (*stackint.Int1024, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return big.NewInt(1), nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) Refund(nodeID []byte) (*types.Transaction, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return nil, nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) WaitUntilRegistration(nodeID []byte) error {
// 	for {
// 		if registered, err := darkNodeRegistrar.IsDarkNodeRegistered(nodeID); err == nil && registered {
// 			return nil
// 		}
// 		time.Sleep(time.Minute)
// 	}
// }
