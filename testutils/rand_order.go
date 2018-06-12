package testutils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/republicprotocol/republic-go/blockchain/ethereum"
	"github.com/republicprotocol/republic-go/cmd/darknode/config"
	"github.com/republicprotocol/republic-go/crypto"
	"github.com/republicprotocol/republic-go/identity"
	"github.com/republicprotocol/republic-go/logger"
	"github.com/republicprotocol/republic-go/order"
)

// RandomOrder will generate a random order.
func RandomOrder() order.Order {
	parity := []order.Parity{order.ParityBuy, order.ParitySell}[rand.Intn(2)]
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]
	volume := RandomCoExp()

	ord := order.NewOrder(order.TypeLimit, parity, time.Now().Add(1*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Int63())
	return ord
}

// RandomOrderMatch will generate a random order and its match.
func RandomOrderMatch() (order.Order, order.Order) {
	tokens := []order.Tokens{order.TokensBTCETH,
		order.TokensETHDGX,
		order.TokensETHREN,
		order.TokensDGXREN,
	}[rand.Intn(6)]
	volume := RandomCoExp()

	buy := order.NewOrder(order.TypeLimit, order.ParityBuy, time.Now().Add(24*time.Hour), tokens, RandomCoExp(), volume, LessRandomCoExp(volume), rand.Int63())
	sell := order.NewOrder(order.TypeLimit, order.ParitySell, time.Now().Add(24*time.Hour), tokens, buy.Price, buy.Volume, buy.MinimumVolume, buy.Nonce)
	return buy, sell
}

// RandomCoExp will generate a random number represented in CoExp format.
func RandomCoExp() order.CoExp {
	co := uint64(rand.Intn(1999) + 1)
	exp := uint64(rand.Intn(25))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}

// LessRandomCoExp will generate a random CoExp that is no more than the given CoExp.
func LessRandomCoExp(coExp order.CoExp) order.CoExp {
	co := uint64(rand.Intn(int(coExp.Co)) + 1)
	exp := uint64(rand.Intn(int(coExp.Exp + 1)))
	return order.CoExp{
		Co:  co,
		Exp: exp,
	}
}

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
			Ethereum: ethereum.Config{
				Network:                 ethereum.NetworkGanache,
				URI:                     "http://localhost:8545",
				RepublicTokenAddress:    ethereum.RepublicTokenAddressOnGanache.String(),
				DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnGanache.String(),
				RenLedgerAddress:        ethereum.RenLedgerAddressOnGanache.String(),
				RenExAccountsAddress:    ethereum.RenExAccountsAddressOnGanache.String(),
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

// RandomNetworkID will generate a random [32]byte array
func RandomNetworkID() [32]byte {
	var networkID [32]byte
	i := fmt.Sprintf("%d", rand.Int())
	hash := crypto.Keccak256([]byte(i))
	copy(networkID[:], hash)

	return networkID
}
