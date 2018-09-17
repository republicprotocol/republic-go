package contract

// TokenAddresses returns the tokens for the provided network
func TokenAddresses(network Network) map[string]string {
	tokens := map[string]string{}
	switch network {
	case NetworkTestnet:
		tokens["ABC"] = "0xfc42491547f1837e2d0f9a0e6b12b1d883fb8bd0"
		tokens["DGX"] = "0x842F0Db4943174EC458b790868e330444c18c9F2"
		tokens["REN"] = "0x99806D107eda625516d954621dF175a002D223e6"
		tokens["PQR"] = "0x724c964a614Eb0748b48dF79eD5D93C108E361c4"
		tokens["XYZ"] = "0xC9382f7b2C683e08AaDe773EB97BcE4a0d6461A0"
	case NetworkFalcon:
		tokens["ABC"] = "0xc96884276D70a1176b2fe102469348d224B0A1fa"
		tokens["DGX"] = "0xF4FAf1b22CEe0a024ad6b12Bb29EC0E13F5827c2"
		tokens["REN"] = "0x87e83f957a2F3A2E5Fe16d5C6B22e38FD28bdc06"
		tokens["PQR"] = "0x295a3894fc98b021735a760dbc7aed265663ca42"
		tokens["XYZ"] = "0x8a4a68DB5Ad08C215c6078111BE8793843A53302"
	case NetworkNightly:
		tokens["ABC"] = "0x49fa7a3B9705Fa8DEb135B7bA64C2Ab00Ab915a1"
		tokens["DGX"] = "0x092eCE29781777604aFAc04887Af30042c3bC5dF"
		tokens["REN"] = "0x15f692D6B9Ba8CEC643C7d16909e8acdEc431bF6"
		tokens["PQR"] = "0xeb5a7335e850176b44ca1990730d1a2433e195f3"
		tokens["XYZ"] = "0x6662449d05312Afe0Ca147Db6Eb155641077883F"
	default:
		panic("unknown network")
	}
	return tokens
}

// RepublicTokenAddress returns the REN contract address.
func RepublicTokenAddress(network Network) string {
	switch network {
	case NetworkTestnet:
		return "0x99806d107eda625516d954621df175a002d223e6"
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
	case NetworkTestnet:
		return "0xd1c3b5f2fe4eec6c262a5e1b161e5e099fd8325e"
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
	case NetworkTestnet:
		return "0xceac6b255ccdd901fefcdb874db092e6f682fee0"
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
	case NetworkTestnet:
		return "0x6c52b2fd5b6c3e6baf47e05af880fc95b9c8079c"
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
	case NetworkTestnet:
		return "0x9a016649d97d44a055c26cbcadbc45a1ac563c89"
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
	case NetworkTestnet:
		return "0xc07780d6e1f24434b1766068f0e44b10a5ff5755"
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
