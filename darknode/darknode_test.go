package darknode_test

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

// NewDarknodes configured for a local test environment.
func NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (darknode.Darknodes, []context.Context, []context.CancelFunc, error) {
	darknodes := make(darknode.Darknodes, numberOfDarknodes)
	ctxs := make([]context.Context, numberOfDarknodes)
	cancels := make([]context.CancelFunc, numberOfDarknodes)

	configs := make([]darknode.Config, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		configs[i] = darknode.NewLocalConfig(*key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
	}
	// FIXME: Load bootstrap nodes into the config file
	// for i := 0; i < numberOfDarknodes; i++ {
	// 	for j := 0; j < numberOfDarknodes; j++ {
	// 		configs[i].NetworkOption.BootstrapMultiAddresses = append(configs[i].NetworkOption.BootstrapMultiAddresses, configs[j].NetworkOption.MultiAddress)
	// 	}
	// }
	var err error
	for i := 0; i < numberOfDarknodes; i++ {
		darknodes[i], err = darknode.NewDarknode(configs[i])
		if err != nil {
			return nil, nil, nil, err
		}
		ctxs[i], cancels[i] = context.WithCancel(context.Background())
	}
	return darknodes, ctxs, cancels, nil
}

// RegisterDarknodes using the minimum required bond and wait until the next
// epoch. This must only be used in local test environments.
func RegisterDarknodes(darknodes darknode.Darknodes, conn client.Connection, darknodeRegistry contracts.DarkNodeRegistry) error {

	minimumBond, err := darknodeRegistry.MinimumBond()
	if err != nil {
		return err
	}

	for i := range darknodes {
		darknodeID := darknodes[i].ID()
		tx, err := darknodeRegistry.Register(darknodeID, []byte{}, &minimumBond)
		if err != nil {
			return err
		}
		if _, err := conn.PatchedWaitMined(context.Background(), tx); err != nil {
			return err
		}
	}
	return nil
}

// DeregisterDarknodes and wait until the next epoch. This must only be used
// in local test environments.
func DeregisterDarknodes(darknodes darknode.Darknodes, conn client.Connection, darknodeRegistry contracts.DarkNodeRegistry) error {
	for i := range darknodes {
		darknodeID := darknodes[i].ID()
		tx, err := darknodeRegistry.Deregister(darknodeID)
		if err != nil {
			return err
		}
		if _, err := conn.PatchedWaitMined(context.Background(), tx); err != nil {
			return err
		}
	}
	return nil
}

// RefundDarknodes after they have been deregistered. This must only be used
// in local test environments.
func RefundDarknodes(darknodes darknode.Darknodes, conn client.Connection, darknodeRegistry contracts.DarkNodeRegistry) error {
	for i := range darknodes {
		darknodeID := darknodes[i].ID()
		tx, err := darknodeRegistry.Refund(darknodeID)
		if err != nil {
			return err
		}
		if _, err := conn.PatchedWaitMined(context.Background(), tx); err != nil {
			return err
		}
	}
	return nil
}
