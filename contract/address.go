package contract

// TokenAddresses returns the tokens for the provided network
func TokenAddresses(network Network) map[string]string {
	tokens := map[string]string{}
	switch network {
	case NetworkMainnet:
		tokens["TUSD"] = "0x8dd5fbce2f6a956c3022ba3663759011dd51e73e"
		tokens["DGX"] = "0x4f3afec4e5a3f2a6a1a411def7d7dfe50ee057bf"
		tokens["REN"] = "0x408e41876cccdc0f92210600ef50372656052a38"
		tokens["ZRX"] = "0xe41d2489571d322189246dafa5ebde1f4699f498"
		tokens["OMG"] = "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07"
	case NetworkTestnet:
		tokens["TUSD"] = "0x525389752ffe6487d33EF53FBcD4E5D3AD7937a0"
		tokens["DGX"] = "0x932F4580B261e9781A6c3c102133C8fDd4503DFc"
		tokens["REN"] = "0x2CD647668494c1B15743AB283A0f980d90a87394"
		tokens["ZRX"] = "0x6EB628dCeFA95802899aD3A9EE0C7650Ac63d543"
		tokens["OMG"] = "0xb57b1105F41e6117F8a67170e1bd7Ec9149D7ced"
	default:
		panic("unknown network")
	}
	return tokens
}

// RepublicTokenAddress returns the REN contract address.
func RepublicTokenAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x21C482f153D0317fe85C60bE1F7fa079019fcEbD"
	case NetworkTestnet:
		return "0x2CD647668494c1B15743AB283A0f980d90a87394"
	default:
		panic("unknown network")
	}
}

// DarknodeRegistryAddress returns the darknode registry contract address.
func DarknodeRegistryAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x3799006a87fde3ccfc7666b3e6553b03ed341c2f"
	case NetworkTestnet:
		return "0x75Fa8349fc9C7C640A4e9F1A1496fBB95D2Dc3d5"
	default:
		panic("unknown network")
	}
}

// DarknodeRewardVaultAddress returns the reward vault contract address.
func DarknodeRewardVaultAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x880407c9cd119bef48b1821cdfc434e3ca3cd588"
	case NetworkTestnet:
		return "0xc08Dfa565EdB7216c3b23bBf0848B43fE9a49F0E"
	default:
		panic("unknown network")
	}
}

// DarknodeSlasherAddress returns slasher contract address.
func DarknodeSlasherAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x04ed8f5730dd4d2b2064cdb6a5bffc26a72962f2"
	case NetworkTestnet:
		return "0x1a3AbB4FfDa92894a5f1637913c031A4621aF9C0"
	default:
		panic("unknown network")
	}
}

// OrderbookAddress returns the orderbook contract address.
func OrderbookAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x6b8bb175c092de7d81860b18db360b734a2598e0"
	case NetworkTestnet:
		return "0xA9b453FC64b4766Aab8a867801d0a4eA7b1474E0"
	default:
		panic("unknown network")
	}
}

// SettlementRegistryAddress returns the settlement registry contract address.
func SettlementRegistryAddress(network Network) string {
	switch network {
	case NetworkMainnet:
		return "0x119da7a8500ade0766f758d934808179dc551036"
	case NetworkTestnet:
		return "0x6Fd909d27faDE71F475fFF50D0274939a5E4fA97"
	default:
		panic("unknown network")
	}
}
