package transact

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	// NetworkMainnet is the public mainnet. Mainnet is not ready to be used.
	// NetworkMainnet = defaultNetwork("mainnet")

	// NetworkTestnet is the public F∅ testnet.
	NetworkTestnet = defaultNetwork("testnet")

	// NetworkLocal is a local development network.
	NetworkLocal = defaultNetwork("local")
)

// Network is a Republic Protocol network. It stores a human-readable string
// identifier and the Ethereum contract addresses that make up the blockchain
// component of Republic Protocol.
type Network struct {
	ID string

	// Republic Protocol
	RepublicToken    common.Address
	DarknodeRegistry common.Address
	Orderbook        common.Address
	RewardVault      common.Address

	// RenEx
	RenExBalances   common.Address
	RenExSettlement common.Address
}

func defaultNetwork(id string) (network Network) {
	switch id {
	case "mainnet":
		// Mainnet
		network = Network{
			ID:               "mainnet",
			RepublicToken:    common.HexToAddress("0x21C482f153D0317fe85C60bE1F7fa079019fcEbD"),
			DarknodeRegistry: common.HexToAddress("0x3799006a87fde3ccfc7666b3e6553b03ed341c2f"),
			Orderbook:        common.HexToAddress("0x6b8bb175c092de7d81860b18db360b734a2598e0"),
			RewardVault:      common.HexToAddress("0x880407c9cd119bef48b1821cdfc434e3ca3cd588"),
			RenExBalances:    common.HexToAddress("0x9636f9ac371ca0965b7c2b4ad13c4cc64d0ff2dc"),
			RenExSettlement:  common.HexToAddress("0x908262de0366e42d029b0518d5276762c92b21e1"),
		}
	case "testnet":
		// F∅ Testnet
		network = Network{
			ID:               "testnet",
			RepublicToken:    common.HexToAddress("0x6f429121a3bd3e6c1c17edbc676eec44cf117faf"),
			DarknodeRegistry: common.HexToAddress("0x372b6204263c6867f81e2a9e11057ff43efea14b"),
			Orderbook:        common.HexToAddress("0xa7caa4780a39d8b8acd6a0bdfb5b906210bc76cd"),
			RewardVault:      common.HexToAddress("0x5d62ccc1086f38286dc152962a4f3e337eec1ec1"),
			RenExBalances:    common.HexToAddress("0xc5b98949AB0dfa0A7d4c07Bb29B002D6d6DA3e25"),
			RenExSettlement:  common.HexToAddress("0xc4f1420de7efbd76e973fe8c99294fe482819f9a"),
		}

	case "local":
		// TODO: Setup local addresses.

	default:
		panic(fmt.Sprintf("unexpected network %v", network))
	}
	return network
}
