package testutils

import (
	"fmt"

	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/contract"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
)

// RandomConfigs will generate n configs and b of them are bootstrap node.
func RandomConfigs(n int, b int) ([]config.Config, error) {
	configs := []config.Config{}

	for i := 0; i < n; i++ {
		keystore, err := crypto.RandomKeystore()
		if err != nil {
			return configs, err
		}

		addr := identity.Address(keystore.Address())
		configs = append(configs, config.Config{
			Keystore:                keystore,
			Host:                    "0.0.0.0",
			Port:                    fmt.Sprintf("%d", 18514+i),
			Address:                 addr,
			BootstrapMultiAddresses: identity.MultiAddresses{},
			Logs: logger.Options{
				Plugins: []logger.PluginOptions{
					{
						File: &logger.FilePluginOptions{
							Path: fmt.Sprintf("%v.out", addr),
						},
					},
				},
			},
			Ethereum: contract.Config{
				Network: contract.NetworkLocal,
				URI:     "http://localhost:8545",
			},
		})
	}

	for i := 0; i < n; i++ {
		for j := 0; j < b; j++ {
			if i == j {
				continue
			}
			bootstrapMultiAddr, err := identity.NewMultiAddressFromString(fmt.Sprintf("/ip4/%v/tcp/%v/republic/%v", configs[j].Host, configs[j].Port, configs[j].Address))
			if err != nil {
				return configs, err
			}
			configs[i].BootstrapMultiAddresses = append(configs[i].BootstrapMultiAddresses, bootstrapMultiAddr)
		}
	}

	return configs, nil
}
