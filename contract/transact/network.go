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
			DarknodeRegistry: common.HexToAddress("0x5d09eb34ce084bece690651f147dbd8ff41007bf"),
			Orderbook:        common.HexToAddress("0xB01219Cf49e92ffcd48fecC96241dBd1372B8Bb1"),
			RewardVault:      common.HexToAddress("0x5d62ccc1086f38286dc152962a4f3e337eec1ec1"),
			RenExBalances:    common.HexToAddress("0xc5b98949AB0dfa0A7d4c07Bb29B002D6d6DA3e25"),
			RenExSettlement:  common.HexToAddress("0xc53abbc5713e606a86533088707e80fcae33eff8"),
		}

	case "falcon":
		// Falcon Testnet
		network = Network{
			ID:               "falcon",
			RepublicToken:    common.HexToAddress("0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"),
			DarknodeRegistry: common.HexToAddress("0x7352e7244899b7cb5d803cc02741c8910d3b75de"),
			Orderbook:        common.HexToAddress("0x044b08eec761c39ac32aee1d6ef0583812f21699"),
			RewardVault:      common.HexToAddress("0x0e6bbbb35835cc3624a000e1698b7b68e9eec7df"),
			RenExBalances:    common.HexToAddress("0x3083e5ba36c6b42ca93c22c803013a4539eedc7f"),
			RenExSettlement:  common.HexToAddress("0x038b63c120a7e60946d6ebaa6dcfc3a475108cc9"),
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
			RenExSettlement:  common.HexToAddress("0xaa14f741b8dc7fc9166807144d881b016f0ee92f"),
		}

	case "local":
		// TODO: Setup local addresses.

	default:
		panic(fmt.Sprintf("unexpected network %v", network))
	}
	return network
}
