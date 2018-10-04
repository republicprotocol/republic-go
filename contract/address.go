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
	case NetworkFalcon:
		tokens["TUSD"] = "0xc96884276D70a1176b2fe102469348d224B0A1fa"
		tokens["DGX"] = "0xF4FAf1b22CEe0a024ad6b12Bb29EC0E13F5827c2"
		tokens["REN"] = "0x87e83f957a2F3A2E5Fe16d5C6B22e38FD28bdc06"
		tokens["ZRX"] = "0x295a3894fc98b021735a760dbc7aed265663ca42"
		tokens["OMG"] = "0x8a4a68DB5Ad08C215c6078111BE8793843A53302"
	case NetworkNightly:
		tokens["TUSD"] = "0x49fa7a3B9705Fa8DEb135B7bA64C2Ab00Ab915a1"
		tokens["DGX"] = "0x092eCE29781777604aFAc04887Af30042c3bC5dF"
		tokens["REN"] = "0x15f692D6B9Ba8CEC643C7d16909e8acdEc431bF6"
		tokens["ZRX"] = "0xeb5a7335e850176b44ca1990730d1a2433e195f3"
		tokens["OMG"] = "0x6662449d05312Afe0Ca147Db6Eb155641077883F"
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
	case NetworkFalcon:
		return "0x87e83f957a2f3a2e5fe16d5c6b22e38fd28bdc06"
	case NetworkNightly:
		return "0x15f692d6b9ba8cec643c7d16909e8acdec431bf6"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
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
	case NetworkFalcon:
		return "0xdaa8c30af85070506f641e456afdb84d4ba972bd"
	case NetworkNightly:
		return "0x8a31d477267a5af1bc5142904ef0afa31d326e03"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
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
	case NetworkFalcon:
		return "0x401e7d7ce6f51ea1a8d4f582413e2fabda68daa8"
	case NetworkNightly:
		return "0xda43560f5fe6c6b5e062c06fee0f6fbc71bbf18a"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
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
	case NetworkFalcon:
		return "0x71ec5f4558e87d6afb5c5ff0b4bdd058d62ed3d1"
	case NetworkNightly:
		return "0x38458ef4a185455cba57a7594b0143c53ad057c1"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
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
	case NetworkFalcon:
		return "0x592d16f8c5fa8f1e074ab3c2cd1acd087adcdc0b"
	case NetworkNightly:
		return "0x376127adc18260fc238ebfb6626b2f4b59ec9b66"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
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
	case NetworkFalcon:
		return "0x6246ff83ddef23d9509ba80aa3ee650ab0321f0b"
	case NetworkNightly:
		return "0x399a70ed71897836468fd74ea19138df90a78d79"
	case NetworkLocal:
		return "0x0000000000000000000000000000000000000000"
	default:
		panic("unknown network")
	}
}
