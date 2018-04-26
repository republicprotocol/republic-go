package darknodetest

import (
	"crypto/rand"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/darknode"
	"github.com/republicprotocol/republic-go/identity"
)

// NewDarknodes configured for a local test environment.
func NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (darknode.Darknodes, error) {
	var err error

	darknodes := make(darknode.Darknodes, numberOfDarknodes)
	multiAddrs := make([]identity.MultiAddress, numberOfDarknodes)
	configs := make([]darknode.Config, numberOfDarknodes)
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
		darknodes[i], err = darknode.NewDarknode(multiAddrs[i], &configs[i])
		if err != nil {
			return nil, err
		}
	}

	return darknodes, nil
}

// RegisterDarknodes using the minimum required bond and wait until the next
// epoch. This must only be used in local test environments.
func RegisterDarknodes(darknodes darknode.Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {

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
func DeregisterDarknodes(darknodes darknode.Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {
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
func RefundDarknodes(darknodes darknode.Darknodes, conn ethereum.Conn, darknodeRegistry dnr.DarknodeRegistry) error {
	for i := range darknodes {
		darknodeID := darknodes[i].ID()
		_, err := darknodeRegistry.Refund(darknodeID)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewLocalConfig(ecdsaKey keystore.Key, host, port string) (identity.MultiAddress, darknode.Config, error) {
	keyPair, err := identity.NewKeyPairFromPrivateKey(ecdsaKey.PrivateKey)
	if err != nil {
		return identity.MultiAddress{}, darknode.Config{}, err
	}

	rsaKey, err := crypto.NewRsaKeyPair()
	if err != nil {
		return identity.MultiAddress{}, darknode.Config{}, err
	}

	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", host, port, keyPair.Address()))
	if err != nil {
		return identity.MultiAddress{}, darknode.Config{}, err
	}
	return multiAddr, darknode.Config{
		EcdsaKey: ecdsaKey,
		RsaKey:   rsaKey,
		Host:     host,
		Port:     port,
		Ethereum: darknode.EthereumConfig{
			Network:                 ethereum.NetworkGanache,
			URI:                     "http://localhost:8545",
			RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
			DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
		},
	}, nil
}

func DistributeEth(darknodes darknode.Darknodes, conn client.Connection, darknodeRegistry contracts.DarknodeRegistry) {

}

func NewFalconConfig() darknode.Config {
	return darknode.Config{}
}

var FalconBootstrapMultis = []string{
	"/ip4/52.79.194.108/tcp/18514/republic/8MGBUdoFFd8VsfAG5bQSAptyjKuutE",
	"/ip4/52.21.44.236/tcp/18514/republic/8MGzXN7M1ucxvtumVjQ7Ybb7xQ8TUw",
	"/ip4/52.41.118.171/tcp/18514/republic/8MHmrykz65HimBPYaVgm8bTSpRUoXA",
	"/ip4/52.59.176.141/tcp/18514/republic/8MKFT9CDQQru1hYqnaojXqCQU2Mmuk",
	"/ip4/52.77.88.84/tcp/18514/republic/8MGb8k337pp2GSh6yG8iv2GK6FbNHN",
}
