package darknode

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/blockchain/ethereum/dnr"
	"github.com/republicprotocol/republic-go/blockchain/test/ganache"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/dispatch"
	"github.com/republicprotocol/republic-go/identity"
)

// The TestnetEnv will use the ethConn to handle Ethereum and the
// darknodeRegistry to acccess the darknodeRegistry of the local
// testnet. TestnetEnv also stores all darknodes of the testnet in
// an exported Darknodes field and the bootstrapMultiAddrs field to
// store the multiaddresses of the bootstrap nodes.
type TestnetEnv struct {
	// Ethereum
	EthConn          ethereum.Conn
	DarknodeRegistry dnr.DarknodeRegistry

	// Darknodes
	BootstrapMultiAddrs identity.MultiAddresses
	Darknodes           Darknodes

	done chan struct{}
}

// NewTestnet will create a testnet that will register newly created darknodes
// to a darknode registry. It will connect to a local ganache server. A call to
// this method must always be folllowed by a call to TearDown after the use of
// the testnet is completed.
func NewTestnet(numberOfDarknodes, numberOfBootstrapDarknodes int) (TestnetEnv, error) {

	darknodes, bootstrapMultiAddrs, err := NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes)
	if err != nil {
		return TestnetEnv{}, fmt.Errorf("cannot create new darknodes: %v", err)
	}

	conn, err := ganache.StartAndConnect()
	if err != nil {
		return TestnetEnv{}, fmt.Errorf("cannot connect to ganache: %v", err)
	}

	// Connect to Ganache
	transactor := ganache.GenesisTransactor()
	darknodeRegistry, err := dnr.NewDarknodeRegistry(context.Background(), conn, &transactor, &bind.CallOpts{})
	if err != nil {
		return TestnetEnv{}, fmt.Errorf("cannot create a new darknode registry: %v", err)
	}

	darknodeRegistry.SetGasLimit(3000000)

	// Register the Darknodes and trigger an epoch to accept their
	// registrations
	err = RegisterDarknodes(darknodes, conn, darknodeRegistry)
	if err != nil {
		return TestnetEnv{}, fmt.Errorf("cannot register darknodes: %v", err)
	}

	// Distribute eth to each node
	for _, node := range darknodes {
		err = ganache.DistributeEth(conn, common.HexToAddress("0x"+hex.EncodeToString(node.ID())))
		if err != nil {
			return TestnetEnv{}, fmt.Errorf("cannot distribute ether to darknode %v: %v", node.ID(), err)
		}
	}

	return TestnetEnv{
		EthConn:             conn,
		DarknodeRegistry:    darknodeRegistry,
		BootstrapMultiAddrs: bootstrapMultiAddrs,
		Darknodes:           darknodes,
		done:                make(chan struct{}),
	}, nil
}

// Teardown must be called in all tests that call NewTestnet to ensure that
// the ganache server is stopped and the darknodes are deregistered and refunded.
func (env *TestnetEnv) Teardown() error {
	defer ganache.Stop()

	close(env.done)

	// Deregister the DarkNodes
	err := DeregisterDarknodes(env.Darknodes, env.EthConn, env.DarknodeRegistry)
	if err != nil {
		return fmt.Errorf("could not deregister darknodes: %v", err)
	}

	// Refund the DarkNodes
	err = RefundDarknodes(env.Darknodes, env.EthConn, env.DarknodeRegistry)
	if err != nil {
		return fmt.Errorf("could not refund darknodes: %v", err)
	}

	return nil
}

// Run all Darknodes. This will start the gRPC services and bootstrap the
// Darknodes. Calls to TestnetEnv.Run are blocking, until a call to
// TestnetEnv.Teardown is made. Errors returned by the Darknodes while running
// are ignored.
// FIXME: Store the errors in a buffer that can be inspected after the test.
func (env *TestnetEnv) Run() {
	dispatch.CoForAll(env.Darknodes, func(i int) {
		env.Darknodes[i].Run(env.done)
		// Ignoring errors as this is a local testnet
		// for err := range env.Darknodes[i].Run(env.done) {
		// 	log.Printf("darknode run-time error: %v", err)
		// }
	})
}

// NewDarknodes configured for a local test environment. This method will also return
// multiaddresses of bootstrap nodes in the testnet.
func NewDarknodes(numberOfDarknodes, numberOfBootstrapDarknodes int) (Darknodes, identity.MultiAddresses, error) {
	var err error

	darknodes := make(Darknodes, numberOfDarknodes)
	bootstrapMultiAddrs := make(identity.MultiAddresses, numberOfBootstrapDarknodes)
	multiAddrs := make([]identity.MultiAddress, numberOfDarknodes)
	configs := make([]Config, numberOfDarknodes)
	for i := 0; i < numberOfDarknodes; i++ {
		multiAddrs[i], configs[i], err = NewLocalConfig("127.0.0.1", fmt.Sprintf("%d", 3000+i))
		if err != nil {
			return nil, nil, err
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

	for j := 0; j < numberOfBootstrapDarknodes; j++ {
		bootstrapMultiAddrs = append(bootstrapMultiAddrs, multiAddrs[j])
	}

	for i := 0; i < numberOfDarknodes; i++ {
		darknodes[i], err = NewDarknode(multiAddrs[i], &configs[i])
		if err != nil {
			return nil, nil, err
		}
	}

	return darknodes, bootstrapMultiAddrs, nil
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

// NewLocalConfig will return newly generated multiaddress and config that are
// constructed from an EcdsaKey, host and port that are passed as arguments to
// the method.
func NewLocalConfig(host, port string) (identity.MultiAddress, Config, error) {
	keystore, err := crypto.RandomKeystore()
	if err != nil {
		return identity.MultiAddress{}, Config{}, err
	}
	multiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%s/tcp/%s/republic/%s", host, port, keystore.Address()))
	if err != nil {
		return identity.MultiAddress{}, Config{}, err
	}
	return multiAddr, Config{
		Keystore: keystore,
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
