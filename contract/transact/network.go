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

	// NetworkFalcon is the internal Falcon testnet.
	NetworkFalcon = defaultNetwork("falcon")

	// NetworkNightly is the internal Nightly testnet.
	NetworkNightly = defaultNetwork("nightly")

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

	case "falcon":
		// Falcon Testnet
		network = Network{
			ID:               "falcon",
			RepublicToken:    common.HexToAddress("0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"),
			DarknodeRegistry: common.HexToAddress("0xfafd5c83d1e21763b79418c4ecb5d62b4970df8e"),
			Orderbook:        common.HexToAddress("0x044b08eec761c39ac32aee1d6ef0583812f21699"),
			RewardVault:      common.HexToAddress("0x0e6bbbb35835cc3624a000e1698b7b68e9eec7df"),
			RenExBalances:    common.HexToAddress("0x3083e5ba36c6b42ca93c22c803013a4539eedc7f"),
			RenExSettlement:  common.HexToAddress("0x8617dcd709bb8660602ef70ade78626b7408a210"),
		}

	case "nightly":
		// Nightly Testnet
		network = Network{
			ID:               "nightly",
			RepublicToken:    common.HexToAddress("0x15f692d6b9ba8cec643c7d16909e8acdec431bf6"),
			DarknodeRegistry: common.HexToAddress("0xb3972e45d16b0942ed34943fdde413190cf5b12a"),
			Orderbook:        common.HexToAddress("0x8356e57aa32547685149a859293ad83c144b800c"),
			RewardVault:      common.HexToAddress("0x7214c4584ab01e61355244e2325ab3f40aca4d85"),
			RenExBalances:    common.HexToAddress("0xc2c126e1eb32e6ad50c611fb92d009b4b4518b00"),
			RenExSettlement:  common.HexToAddress("0x65712325c41fb39b9205e08483b43142d919cc42"),
		}

	case "local":
		// TODO: Setup local addresses.

	default:
		panic(fmt.Sprintf("unexpected network %v", network))
	}
	return network
}
