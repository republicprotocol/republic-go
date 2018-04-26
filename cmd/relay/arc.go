package main

func processAtomicSwaps(swaps <-chan) {

	conn, err := client.Connect(uri, network, republicTokenAddress, darknodeRegistryAddr)
	if err != nil {
		return err
	}
	transOps := bind.NewKeyedTransactor(privateKey)
	arc.NewArc(context.Background(), conn, transOps)

	return nil
	
}

//	func atomicSwap(entries <-chan orderbook.Entry, privateKey *ecdsa.PrivateKey) error {