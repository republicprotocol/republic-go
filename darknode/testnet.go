package darknode

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
)

type TestnetEnv struct {
	// Ethereum
	ethConn          ethereum.Conn
	darknodeRegistry dnr.DarknodeRegistry

	// Darknodes
	bootstrapMultiAddrs identity.MultiAddresses
	darknodes           Darknodes
}

func NewTestnet(done <-chan struct{}, numberOfDarknodes, numberOfBootstrapDarknodes int) (TestnetEnv, error) {
	// TODO:
	// 1. Call NewDarknodes
	// 2. Call ganache.StartAndConnect
	// 3. Call dnr.NewDarknodeRegistry
	// 4. Collect all bootstrap Darknode multi-addresses

	panic("unimplemented")
}

// NewDarknodes configured for a local test environment.
func NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (Darknodes, error) {
	var err error

	darknodes := make(Darknodes, numberOfDarknodes)
	multiAddrs := make([]identity.MultiAddress, numberOfDarknodes)
	configs := make([]Config, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		key := keystore.NewKeyForDirectICAP(rand.Reader)
		multiAddrs[i], configs[i], err = NewLocalConfig(*key, "127.0.0.1", fmt.Sprintf("%d", 3000+i))
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < numberOfDarknodes; i++ {
		for j := 0; j < numberOfBootstrapDarknodes; j++ {
			if i == j {
				continue
			}
			configs[i].BootstrapMultiAddresses = append(configs[i].BootstrapMultiAddresses, multiAddrs[j])
		}
	}
	for i := 0; i < numberOfDarknodes; i++ {
		darknodes[i], err = NewDarknode(multiAddrs[i], &configs[i])
		if err != nil {
			return nil, err
		}
	}

	return darknodes, nil
}

// RegisterDarknodes using the minimum required bond and wait until the next
// epoch. This must only be used in local test environments.
func RegisterDarknodes(darknodes Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {

	minimumBond, err := darknodeRegistry.MinimumBond()
	if err != nil {
		return err
	}

	for i := range darknodes {
		darknodeID := darknodes[i].ID()

		_, err := darknodeRegistry.ApproveRen(&minimumBond)
		if err != nil {
			return err
		}
		_, err = darknodeRegistry.Register(darknodeID, []byte{}, &minimumBond)
		if err != nil {
			return err
		}
	}

	// Turn the epoch to approve registrations
	return darknodeRegistry.WaitForEpoch()
}

// DeregisterDarknodes and wait until the next epoch. This must only be used
// in local test environments.
func DeregisterDarknodes(darknodes Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {
	for i := range darknodes {
		darknode := darknodes[i]
		_, err := darknodeRegistry.Deregister(darknode.ID())
		if err != nil {
			return err
		}
	}
	return darknodeRegistry.WaitForEpoch()
}

// RefundDarknodes after they have been deregistered. This must only be used
// in local test environments.
func RefundDarknodes(darknodes Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {
	for i := range darknodes {
		darknodeID := darknodes[i].ID()
		_, err := darknodeRegistry.Refund(darknodeID)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewLocalConfig(ecdsaKey keystore.Key, host, port string) (identity.MultiAddress, Config, error) {
	keyPair, err := identity.NewKeyPairFromPrivateKey(ecdsaKey.PrivateKey)
	if err != nil {
		return identity.MultiAddress{}, Config{}, err
	}

	rsaKey, err := crypto.NewRsaKeyPair()
	if err != nil {
		return identity.MultiAddress{}, Config{}, err
	}

	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", host, port, keyPair.Address()))
	if err != nil {
		return identity.MultiAddress{}, Config{}, err
	}
	return multiAddr, Config{
		EcdsaKey: ecdsaKey,
		RsaKey:   rsaKey,
		Host:     host,
		Port:     port,
		Ethereum: ethereum.Config{
			Network:                 ethereum.NetworkGanache,
			URI:                     "http://localhost:8545",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
		},
	}, nil
}
