package ledger_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/republicprotocol/republic-go/blockchain/ethereum"
)

var _ = Describe("Ren Ledger Contract", func() {
	Context("interacting with deployment on Ropsten ", func() {
		It("should be able to open order ", func() {
			config := ethereum.Config{
				Network:                 ethereum.NetworkRopsten,
				URI:                     "https://ropsten.infura.io",
				RepublicTokenAddress:    ethereum.RepublicTokenAddressOnRopsten.Hex(),
				DarknodeRegistryAddress: ethereum.DarknodeRegistryAddressOnRopsten.Hex(),
				HyperdriveAddress:       ethereum.HyperdriveAddressOnRopsten.Hex(),
				ArcAddress:              ethereum.ArcAddressOnRopsten.Hex(),
			}
		})
	})
})
