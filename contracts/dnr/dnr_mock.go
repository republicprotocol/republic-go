package dnr

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/republicprotocol/republic-go/contracts/connection"
)

// // MockDarkNodeRegistry implemented to DarkNodeRegistry interface for
// // testing.
// type MockDarkNodeRegistry struct {
// 	do.GuardedObject

// 	registered   [][]byte
// 	toRegister   [][]byte
// 	toDeregister [][]byte
// 	epoch        Epoch
// }

// TestnetDNR returns a new MockDarkNodeRegistry
func TestnetDNR(auth *bind.TransactOpts) (DarkNodeRegistry, error) {

	conn, err := connection.ConnectToTestnet()
	if err != nil {
		return DarkNodeRegistry{}, err
	}

	if auth == nil {
		auth = connection.GenesisAuth
	} else {
		err := connection.DistributeEth(conn, auth.From)
		if err != nil {
			return DarkNodeRegistry{}, err
		}
		err = connection.DistributeRen(conn, auth.From)
		if err != nil {
			return DarkNodeRegistry{}, err
		}
	}
	return NewDarkNodeRegistry(context.Background(), &conn, auth, &bind.CallOpts{})
}

// // Register a new dark node
// func (mockDnr *MockDarkNodeRegistry) Register(darkNodeID []byte, publicKey []byte, bond *stackint.Int1024) (*types.Transaction, error) {
// 	mockDnr.Enter(nil)
// 	defer mockDnr.Exit()
// 	isRegistered, _ := mockDnr.isDarkNodeRegistered(darkNodeID)
// 	isPending, _ := mockDnr.isDarkNodePendingRegistration(darkNodeID)
// 	if isRegistered || isPending {
// 		return nil, errors.New("Must not be registered to register")
// 	}
// 	mockDnr.toRegister = append(mockDnr.toRegister, darkNodeID)
// 	return nil, nil
// }

// // Deregister an existing dark node
// func (mockDnr *MockDarkNodeRegistry) Deregister(darkNodeID []byte) (*types.Transaction, error) {
// 	mockDnr.Enter(nil)
// 	defer mockDnr.Exit()
// 	for i, id := range mockDnr.toRegister {
// 		if string(darkNodeID) == string(id) {
// 			mockDnr.toDeregister[i] = mockDnr.toDeregister[len(mockDnr.toDeregister)-1]
// 			mockDnr.toDeregister = mockDnr.toDeregister[:len(mockDnr.toDeregister)-1]
// 			return nil, nil
// 		}
// 	}
// 	if isRegistered, _ := mockDnr.isDarkNodeRegistered(darkNodeID); !isRegistered {
// 		return nil, errors.New("Must be registered to deregister")
// 	}
// 	mockDnr.toDeregister = append(mockDnr.toRegister, darkNodeID)
// 	return nil, nil
// }

// // GetBond retrieves the bond of an existing dark node
// func (mockDnr *MockDarkNodeRegistry) GetBond(darkNodeID []byte) (stackint.Int1024, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return stackint.FromUint(86000), nil
// }

// // IsRegistered returns true if the node is registered
// func (mockDnr *MockDarkNodeRegistry) IsRegistered(darkNodeID []byte) (bool, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return mockDnr.isDarkNodeRegistered(darkNodeID)
// }

// func (mockDnr *MockDarkNodeRegistry) isDarkNodeRegistered(darkNodeID []byte) (bool, error) {
// 	for _, id := range mockDnr.registered {
// 		if string(darkNodeID) == string(id) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// // IsDarkNodePendingRegistration returns true if the node will become registered at the next epoch
// func (mockDnr *MockDarkNodeRegistry) IsDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return mockDnr.isDarkNodePendingRegistration(darkNodeID)
// }

// func (mockDnr *MockDarkNodeRegistry) isDarkNodePendingRegistration(darkNodeID []byte) (bool, error) {
// 	for _, id := range mockDnr.toRegister {
// 		if string(darkNodeID) == string(id) {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// // CurrentEpoch returns the current epoch
// func (mockDnr *MockDarkNodeRegistry) CurrentEpoch() (Epoch, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return mockDnr.epoch, nil
// }

// // Epoch updates the current Epoch if the Minimum Epoch Interval has passed since the previous Epoch
// func (mockDnr *MockDarkNodeRegistry) Epoch() (*types.Transaction, error) {

// 	mockDnr.Enter(nil)
// 	defer mockDnr.Exit()

// 	var b32 [32]byte

// 	_, err := rand.Read(b32[:])
// 	if err != nil {
// 		return nil, err
// 	}

// 	now := stackint.FromUint(uint(time.Now().Unix()))
// 	mockDnr.epoch = Epoch{
// 		Blockhash: b32,
// 		Timestamp: &now,
// 	}

// 	// Remove toRegister nodes
// 	for _, deregNode := range mockDnr.toDeregister {
// 		for i := 0; i < len(mockDnr.registered); i++ {
// 			if string(mockDnr.registered[i]) == string(deregNode) {
// 				mockDnr.registered[i] = mockDnr.registered[len(mockDnr.registered)-1]
// 				mockDnr.registered = mockDnr.registered[:len(mockDnr.registered)-1]
// 				break
// 			}
// 		}
// 	}

// 	mockDnr.registered = append(mockDnr.registered, mockDnr.toRegister...)

// 	mockDnr.toDeregister = make([][]byte, 0)
// 	mockDnr.toRegister = make([][]byte, 0)

// 	return nil, nil
// }

// // GetCommitment gets the signed commitment (not implemented)
// func (mockDnr *MockDarkNodeRegistry) GetCommitment(darkNodeID []byte) ([32]byte, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return [32]byte{}, nil
// }

// // GetOwner gets the owner of the given dark node (not implemented)
// func (mockDnr *MockDarkNodeRegistry) GetOwner(darkNodeID []byte) (common.Address, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return [20]byte{}, nil
// }

// // GetPublicKey gets the public key of the goven dark node (not implemented)
// func (mockDnr *MockDarkNodeRegistry) GetPublicKey(darkNodeID []byte) ([]byte, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return nil, nil
// }

// // GetAllNodes gets all dark nodes
// func (mockDnr *MockDarkNodeRegistry) GetAllNodes() ([][]byte, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return mockDnr.registered, nil
// }

// // MinimumBond gets the minimum viable bond amount (hard-coded to 86000)
// func (mockDnr *MockDarkNodeRegistry) MinimumBond() (stackint.Int1024, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return stackint.FromUint(86000), nil
// }

// // MinimumEpochInterval gets the minimum epoch interval (hard-coded to 0)
// func (mockDnr *MockDarkNodeRegistry) MinimumEpochInterval() (stackint.Int1024, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return stackint.FromUint(0), nil
// }

// // Refund refunds the bond of an unregistered miner (not implemented)
// func (mockDnr *MockDarkNodeRegistry) Refund(darkNodeID []byte) (*types.Transaction, error) {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	return nil, nil
// }

// // WaitUntilRegistration waits until the registration is successful
// func (mockDnr *MockDarkNodeRegistry) WaitUntilRegistration(darkNodeID []byte) error {
// 	mockDnr.EnterReadOnly(nil)
// 	defer mockDnr.ExitReadOnly()
// 	for {
// 		isRegistered, err := mockDnr.isDarkNodeRegistered(darkNodeID)
// 		if err != nil {
// 			return err
// 		}
// 		if isRegistered {
// 			break
// 		}
// 		time.Sleep(time.Minute)
// 	}
// 	return nil
// }

// // type MockDarkNodeRegistry struct {
// // 	do.GuardedObject

// // 	hash      [32]byte
// // 	timestamp *stackint.Int1024
// // 	nodeIDs   map[string]bool
// // }

// // func NewMockDarkNodeRegistry(nodeIDs [][]byte) DarkNodeRegistry {
// // 	darkNodeRegistry := new(MockDarkNodeRegistry)
// // 	darkNodeRegistry.GuardedObject = do.NewGuardedObject()
// // 	darkNodeRegistry.hash = [32]byte{1}
// // 	darkNodeRegistry.timestamp = big.NewInt(1)
// // 	darkNodeRegistry.nodeIDs = map[string]bool{}
// // 	for _, nodeID := range nodeIDs {
// // 		darkNodeRegistry.nodeIDs[string(nodeID)] = true
// // 	}
// // 	return darkNodeRegistry
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) Register(nodeID []byte, publicKey []byte) (*types.Transaction, error) {
// // 	darkNodeRegistry.Enter(nil)
// // 	defer darkNodeRegistry.Exit()
// // 	darkNodeRegistry.nodeIDs[string(nodeID)] = true
// // 	return nil, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) Deregister(nodeID []byte) (*types.Transaction, error) {
// // 	darkNodeRegistry.Enter(nil)
// // 	defer darkNodeRegistry.Exit()
// // 	delete(darkNodeRegistry.nodeIDs, string(nodeID))
// // 	return nil, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) GetBond(nodeID []byte) (*stackint.Int1024, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	if _, ok := darkNodeRegistry.nodeIDs[string(nodeID)]; ok {
// // 		return big.NewInt(86000), nil
// // 	}
// // 	return big.NewInt(0), nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) IsRegistered(nodeID []byte) (bool, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	_, ok := darkNodeRegistry.nodeIDs[string(nodeID)]
// // 	return ok, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) IsDarkNodePendingRegistration(nodeID []byte) (bool, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return false, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) CurrentEpoch() (Epoch, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return Epoch{
// // 		Blockhash: darkNodeRegistry.hash,
// // 		Timestamp: darkNodeRegistry.timestamp,
// // 	}, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) Epoch() (*types.Transaction, error) {
// // 	darkNodeRegistry.Enter(nil)
// // 	defer darkNodeRegistry.Exit()
// // 	darkNodeRegistry.timestamp.Add(darkNodeRegistry.timestamp, big.NewInt(1))
// // 	darkNodeRegistry.hash[0]++
// // 	return nil, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) GetCommitment(nodeID []byte) ([32]byte, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return [32]byte{}, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) GetOwner(nodeID []byte) (common.Address, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return common.Address{}, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) GetPublicKey(nodeID []byte) ([]byte, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return []byte{}, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) GetAllNodes() ([][]byte, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	allNodes := make([][]byte, 0, len(darkNodeRegistry.nodeIDs))
// // 	for nodeID := range darkNodeRegistry.nodeIDs {
// // 		allNodes = append(allNodes, []byte(nodeID))
// // 	}
// // 	return allNodes, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) MinimumBond() (*stackint.Int1024, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return big.NewInt(86000), nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) MinimumEpochInterval() (*stackint.Int1024, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return big.NewInt(1), nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) Refund(nodeID []byte) (*types.Transaction, error) {
// // 	darkNodeRegistry.EnterReadOnly(nil)
// // 	defer darkNodeRegistry.ExitReadOnly()
// // 	return nil, nil
// // }

// // func (darkNodeRegistry *MockDarkNodeRegistry) WaitUntilRegistration(nodeID []byte) error {
// // 	for {
// // 		if registered, err := darkNodeRegistry.IsRegistered(nodeID); err == nil && registered {
// // 			return nil
// // 		}
// // 		time.Sleep(time.Minute)
// // 	}
// // }
