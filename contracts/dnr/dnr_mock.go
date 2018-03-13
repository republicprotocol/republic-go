package dnr

import (
	"crypto/rand"
	"errors"
	"math/big"
	"time"

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
	mockDnr.epoch = Epoch{
		Blockhash: b32,
		Timestamp: big.NewInt(time.Now().Unix()),
	}
	return mockDnr, nil
}

func (mockDnr *MockDarkNodeRegistrar) Register(_darkNodeID []byte, _publicKey []byte) (*types.Transaction, error) {
	isRegistered, _ := mockDnr.IsDarkNodeRegistered(_darkNodeID)
	isPending, _ := mockDnr.IsDarkNodePendingRegistration(_darkNodeID)
	if isRegistered || isPending {
		return nil, errors.New("Must not be registered to register")
	}
	mockDnr.toRegister = append(mockDnr.toRegister, _darkNodeID)
	return nil, nil
}

func (mockDnr *MockDarkNodeRegistrar) Deregister(_darkNodeID []byte) (*types.Transaction, error) {
	for i, id := range mockDnr.toRegister {
		if string(_darkNodeID) == string(id) {
			mockDnr.toDeregister[i] = mockDnr.toDeregister[len(mockDnr.toDeregister)-1]
			mockDnr.toDeregister = mockDnr.toDeregister[:len(mockDnr.toDeregister)-1]
			return nil, nil
		}
	}
	if isRegistered, _ := mockDnr.IsDarkNodeRegistered(_darkNodeID); !isRegistered {
		return nil, errors.New("Must be registered to deregister")
	}
	mockDnr.toDeregister = append(mockDnr.toRegister, _darkNodeID)
	return nil, nil
}

func (mockDnr *MockDarkNodeRegistrar) GetBond(_darkNodeID []byte) (*big.Int, error) {
	return big.NewInt(86000), nil
}

func (mockDnr *MockDarkNodeRegistrar) IsDarkNodeRegistered(_darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.registered {
		if string(_darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}

func (mockDnr *MockDarkNodeRegistrar) IsDarkNodePendingRegistration(_darkNodeID []byte) (bool, error) {
	for _, id := range mockDnr.toRegister {
		if string(_darkNodeID) == string(id) {
			return true, nil
		}
	}
	return false, nil
}

func (mockDnr *MockDarkNodeRegistrar) CurrentEpoch() (Epoch, error) {
	return mockDnr.epoch, nil
}

func (mockDnr *MockDarkNodeRegistrar) Epoch() (*types.Transaction, error) {
	var b32 [32]byte

	_, err := rand.Read(b32[:])
	if err != nil {
		return nil, err
	}

	mockDnr.epoch = Epoch{
		Blockhash: b32,
		Timestamp: big.NewInt(time.Now().Unix()),
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

func (mockDnr *MockDarkNodeRegistrar) GetCommitment(_darkNodeID []byte) ([32]byte, error) {
	var nil32 [32]byte
	return nil32, nil
}

func (mockDnr *MockDarkNodeRegistrar) GetOwner(_darkNodeID []byte) (common.Address, error) {
	var nil20 [20]byte
	return nil20, nil
}

func (mockDnr *MockDarkNodeRegistrar) GetPublicKey(_darkNodeID []byte) ([]byte, error) {
	return nil, nil
}

func (mockDnr *MockDarkNodeRegistrar) GetAllNodes() ([][]byte, error) {
	return mockDnr.registered, nil
}

func (mockDnr *MockDarkNodeRegistrar) MinimumBond() (*big.Int, error) {
	return big.NewInt(86000), nil
}

func (mockDnr *MockDarkNodeRegistrar) MinimumEpochInterval() (*big.Int, error) {
	return big.NewInt(0), nil
}

func (mockDnr *MockDarkNodeRegistrar) Refund(_darkNodeID []byte) (*types.Transaction, error) {
	return nil, nil
}

func (mockDnr *MockDarkNodeRegistrar) WaitUntilRegistration(_darkNodeID []byte) error {
	for {
		isRegistered, err := mockDnr.IsDarkNodeRegistered(_darkNodeID)
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
// 	timestamp *big.Int
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

// func (darkNodeRegistrar *MockDarkNodeRegistrar) GetBond(nodeID []byte) (*big.Int, error) {
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

// func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumBond() (*big.Int, error) {
// 	darkNodeRegistrar.EnterReadOnly(nil)
// 	defer darkNodeRegistrar.ExitReadOnly()
// 	return big.NewInt(86000), nil
// }

// func (darkNodeRegistrar *MockDarkNodeRegistrar) MinimumEpochInterval() (*big.Int, error) {
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
