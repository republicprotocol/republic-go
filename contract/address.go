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
		tokens["TUSD"] = "0xD09A65Fd6DF182dBf9cC853697eFA520369015E4"
		tokens["DGX"] = "0x7583D3019b078037B8607487cc4c161e19C50869"
		tokens["REN"] = "0x81793734c6Cf6961B5D0D2d8a30dD7DF1E1803f1"
		tokens["ZRX"] = "0x932d170Cd254Db4c7321C6A89D7722714d82a69f"
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
		return "0x81793734c6Cf6961B5D0D2d8a30dD7DF1E1803f1"
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
		return "0xf7daA0Baf257547A6Ad3CE7FFF71D55cb7426F76"
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
		return "0x0C03998EdF2fC7F29801A05CDeCA2289dD51A158"
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
		return "0x9c65050f9Bc35De0B25f99643797667Ee300aeDa"
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
		return "0xA53Da4093c682a4259DE38302341BFEf7e9f7a4f"
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
		return "0x762D83C9e39882b86cfdbd113a3B95804c1d6A31"
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
