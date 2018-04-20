package darknodetest

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/ethereum/client"
	"github.com/republicprotocol/republic-go/ethereum/contracts"
)

// NewDarknodes configured for a local test environment.
func NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (darknode.Darknodes, error) {
	var err error

	darknodes := make(darknode.Darknodes, numberOfDarknodes)
	configs := make([]darknode.Config, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		configs[i], err = darknode.NewLocalConfig(*key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < numberOfDarknodes; i++ {
		for j := 0; j < numberOfDarknodes; j++ {
			configs[i].Network.BootstrapMultiAddresses = append(configs[i].Network.BootstrapMultiAddresses, configs[j].Network.MultiAddress)
		}
	}
	for i := 0; i < numberOfDarknodes; i++ {
		darknodes[i], err = darknode.NewDarknode(configs[i])
		if err != nil {
			return nil, err
		}
	}

	return darknodes, nil
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

	// Turn the epoch to approve registrations
	time.Sleep(time.Second)
	tx, err := darknodeRegistry.Epoch()
	if err != nil {
		return err
	}
	if _, err := conn.PatchedWaitMined(context.Background(), tx); err != nil {
		return err
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
