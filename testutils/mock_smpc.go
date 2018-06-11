package testutils

import (
	"math/rand"

	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/shamir"
	"github.com/republicprotocol/republic-go/smpc"
)

// Smpc is a mock implementation of the smpc.Smpcer interface.
type Smpc struct {
	value          uint64
	useRandomValue bool
}

// NewAlwaysMatchSmpc returns a mock implementation of an smpc.Smpcer interface
// that will call smpc.Callbacks with values that will always result in a
// match.
func NewAlwaysMatchSmpc() *Smpc {
	return &Smpc{
		value:          0,
		useRandomValue: false,
	}
}

// NewAlwaysMismatchSmpc returns a mock implementation of an smpc.Smpcer
// interface that will call smpc.Callbacks with values that will never result
// in a match.
func NewAlwaysMismatchSmpc() *Smpc {
	return &Smpc{
		value:          shamir.Prime - 1,
		useRandomValue: false,
	}
}

// NewSmpc returns a mock implementation of an smpc.Smpcer interface that will
// call smpc.Callbacks with values that will result in a match half of the
// time.
func NewSmpc() *Smpc {
	return &Smpc{
		value:          0,
		useRandomValue: true,
	}
}

// Connect implements smpc.Smpcer.
func (smpc *Smpc) Connect(networkID smpc.NetworkID, nodes identity.Addresses) {
}

// Disconnect implements smpc.Smpcer.
func (smpc *Smpc) Disconnect(networkID smpc.NetworkID) {
}

// Join implements smpc.Smpcer.
func (smpc *Smpc) Join(networkID smpc.NetworkID, join smpc.Join, callback smpc.Callback) error {
	values := make([]uint64, len(join.Shares))
	for i := range values {
		if smpc.useRandomValue {
			if rand.Uint64()%2 == 0 {
				values[i] = 0
			} else {
				values[i] = shamir.Prime - 1
			}
			continue
		}
		values[i] = smpc.value
	}
	callback(join.ID, values)
	return nil
}
