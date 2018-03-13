package network_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/compute"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/network"
	"github.com/republicprotocol/republic-go/network/dht"
	"github.com/republicprotocol/republic-go/network/rpc"
	"github.com/republicprotocol/republic-go/order"
	"google.golang.org/grpc"
)

const (
	NumberOfBootstrapNodes = 5
)

// MockDelegate for testing purpose
type MockDelegate struct {
}

func (mockDelegate *MockDelegate) OnOpenOrder(from identity.MultiAddress, orderFragment *order.Fragment) {
	return
}

func (mockDelegate *MockDelegate) OnBroadcastDeltaFragment(from identity.MultiAddress, deltaFragment *compute.DeltaFragment) {
	return
}

func generateSwarmServices(numberOfSwarms int) ([]*network.SwarmService, []*grpc.Server, error) {
	// Initialize bootstrap nodes and swarm nodes.
	swarms := make([]*network.SwarmService, NumberOfBootstrapNodes+numberOfSwarms)
	bootstrapNodes := make([]identity.MultiAddress, NumberOfBootstrapNodes)

	for i := 0; i < len(swarms); i++ {
		address, _, err := identity.NewAddress()
		if err != nil {
			return nil, nil, err
		}
		options := network.Options{}

		if i < NumberOfBootstrapNodes {
			multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000+i, address))
			if err != nil {
				return nil, nil, err
			}
			options.MultiAddress = multi
			bootstrapNodes[i] = multi
		} else {
			multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000+i, address))
			if err != nil {
				return nil, nil, err
			}
			options.MultiAddress = multi
		}
		options.BootstrapMultiAddresses = []identity.MultiAddress{}
		options.Debug = network.DebugHigh
		options.Alpha = 3
		options.MaxBucketLength = 20
		options.ClientPoolCacheLimit = 20
		options.TimeoutBackoff = 5 * time.Second
		options.Timeout = 30 * time.Second
		options.TimeoutRetries = 2
		options.Concurrent = false

		l, err := logger.NewLogger(logger.Options{
			Plugins: []logger.PluginOptions{
				{File: &logger.FilePluginOptions{Path: "stdout"}},
			},
		})
		if err != nil {
			return nil, nil, err
		}
		l.Start()

		swarms[i] = network.NewSwarmService(MockDelegate{}, options,
			&logger.Logger{}, rpc.NewClientPool(options.MultiAddress),
			dht.NewDHT(options.MultiAddress.Address(), options.MaxBucketLength))

	}
	for i := 0; i < len(swarms); i++ {
		for j := range bootstrapNodes {
			if i == j {
				continue
			}
			swarms[i].BootstrapMultiAddresses = append(swarms[i].BootstrapMultiAddresses, bootstrapNodes[j])
		}
	}

	return swarms, make([]*grpc.Server, len(swarms)), nil
}

func connectSwarms(nodes []*network.SwarmService, connectivity int) error {
	for i, from := range nodes {
		for j, to := range nodes {
			if i == j {
				continue
			}
			// Connect bootstrap nodes in a fully connected topology
			if i < NumberOfBootstrapNodes {
				if j < NumberOfBootstrapNodes {
					err := from.ClientPool.Ping(to.MultiAddress())
					if err != nil {
						return err
					}
				}
				continue
			}
			// Connect standard nodes randomly
			isConnected := rand.Intn(100) < connectivity
			if isConnected {
				err := from.ClientPool.Ping(to.MultiAddress())
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func generateDarkServices(numberOfDarkService int) ([]*network.DarkService, []*grpc.Server, error) {
	// Initialize bootstrap nodes and dark services.
	nodes := make([]*network.DarkService, NumberOfBootstrapNodes+numberOfDarkService)
	bootstrapNodes := make([]identity.MultiAddress, NumberOfBootstrapNodes)

	for i := 0; i < len(nodes); i++ {
		address, _, err := identity.NewAddress()
		if err != nil {
			return nil, nil, err
		}
		options := network.Options{}

		if i < NumberOfBootstrapNodes {
			multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 3000+i, address))
			if err != nil {
				return nil, nil, err
			}
			options.MultiAddress = multi
			bootstrapNodes[i] = multi
		} else {
			multi, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d/republic/%s", 4000+i, address))
			if err != nil {
				return nil, nil, err
			}
			options.MultiAddress = multi
		}
		options.BootstrapMultiAddresses = []identity.MultiAddress{}
		options.Debug = network.DebugHigh
		options.Alpha = 3
		options.MaxBucketLength = 20
		options.ClientPoolCacheLimit = 20
		options.TimeoutBackoff = 5 * time.Second
		options.Timeout = 30 * time.Second
		options.TimeoutRetries = 2
		options.Concurrent = false

		l, err := logger.NewLogger(logger.Options{
			Plugins: []logger.PluginOptions{
				{File: &logger.FilePluginOptions{Path: "stdout"}},
			},
		})
		if err != nil {
			return nil, nil, err
		}
		l.Start()
		nodes[i] = network.NewDarkService(&MockDelegate{}, options, &logger.Logger{})

	}
	for i := 0; i < len(nodes); i++ {
		for j := range bootstrapNodes {
			if i == j {
				continue
			}
			nodes[i].BootstrapMultiAddresses = append(nodes[i].BootstrapMultiAddresses, bootstrapNodes[j])
		}
	}

	return nodes, make([]*grpc.Server, len(nodes)), nil
}
