package topology

import (
	"github.com/republicprotocol/go-identity"
	"sync"
)

type MockDelegate struct {
	mu            sync.Mutex
	PingCount     int
	FragmentCount int
}

func NewMockDelegate() *MockDelegate {
	return &MockDelegate{
		PingCount:     0,
		FragmentCount: 0,
	}
}

func (mockDelegate *MockDelegate) OnPingReceived(peer identity.MultiAddress) {
	mockDelegate.mu.Lock()
	mockDelegate.PingCount++
	mockDelegate.mu.Unlock()
}

func (mockDelegate *MockDelegate) OnOrderFragmentReceived() {
	mockDelegate.mu.Lock()
	mockDelegate.FragmentCount++
	mockDelegate.mu.Unlock()
}
