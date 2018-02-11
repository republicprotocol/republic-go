package main

type extractSecretCmd struct {
	redemptionTx *wire.MsgTx
	secretHash   []byte
}

func ExtractSecret(redemptionTransaction string, secretHash []byte, chain string, rpcUser string, rpcPass string) (Error error, showUsage bool){
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}

	redemptionTxBytes, err := hex.DecodeString(redemptionTransaction)
		if err != nil {
			return fmt.Errorf("failed to decode redemption transaction: %v", err), true
		}
		var redemptionTx wire.MsgTx
		err = redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode redemption transaction: %v", err), true
		}

		if len(secretHash) != sha256.Size {
			return errors.New("secret hash has wrong size"), true
		}

		cmd := &extractSecretCmd{redemptionTx: &redemptionTx, secretHash: secretHash}

		err = cmd.runCommand()
		return err, false
}

func (cmd *extractSecretCmd) runCommand() error {
	// Loop over all pushed data from all inputs, searching for one that hashes
	// to the expected hash.  By searching through all data pushes, we avoid any
	// issues that could be caused by the initiator redeeming the participant's
	// contract with some "nonstandard" or unrecognized transaction or script
	// type.
	for _, in := range cmd.redemptionTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), cmd.secretHash) {
				fmt.Printf("Secret: %x\n", push)
				return nil
			}
		}
	}
	return errors.New("transaction does not contain the secret")
}