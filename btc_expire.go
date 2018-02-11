
package main

type expireCmd struct {
	contract   []byte
	contractTx *wire.MsgTx
}

func Expire(contractString string, contractTransaction string, chain string, rpcUser string, rpcPass string) (Error error, showUsage bool) {
	var chainParams *chaincfg.Params ;
	if (chain == "testnet"){
		chainParams = &chaincfg.TestNet3Params
	} else {
		chainParams = &chaincfg.MainNetParams
	}
	
	contract, err := hex.DecodeString(contractString)
		if err != nil {
			return fmt.Errorf("failed to decode contract: %v", err), true
		}

		contractTxBytes, err := hex.DecodeString(contractTransaction)
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}
		var contractTx wire.MsgTx
		err = contractTx.Deserialize(bytes.NewReader(contractTxBytes))
		if err != nil {
			return fmt.Errorf("failed to decode contract transaction: %v", err), true
		}

		cmd := &refundCmd{contract: contract, contractTx: &contractTx}

		connect, err := normalizeAddress("localhost", walletPort(chainParams))
	if err != nil {
		return fmt.Errorf("wallet server address: %v", err), true
	}

	connConfig := &rpc.ConnConfig{
		Host:         connect,
		User:         rpcUser,
		Pass:         rpcPass,
		DisableTLS:   true,
		HTTPPostMode: true,
	}

	client, err := rpc.New(connConfig, nil)
	if err != nil {
		return fmt.Errorf("rpc connect: %v", err), false
	}
	defer func() {
		client.Shutdown()
		client.WaitForShutdown()
	}()

	err = cmd.runCommand(client)
	return err, false
}


func (cmd *refundCmd) runCommand(c *rpc.Client) error {
	pushes, err := txscript.ExtractAtomicSwapDataPushes(cmd.contract)
	if err != nil {
		return err
	}
	if pushes == nil {
		return errors.New("contract is not an atomic swap script recognized by this tool")
	}

	feePerKb, minFeePerKb, err := getFeePerKb(c)
	if err != nil {
		return err
	}

	refundTx, refundFee, err := buildRefund(c, cmd.contract, cmd.contractTx, feePerKb, minFeePerKb)
	if err != nil {
		return err
	}
	refundTxHash := refundTx.TxHash()
	var buf bytes.Buffer
	buf.Grow(refundTx.SerializeSize())
	refundTx.Serialize(&buf)

	refundFeePerKb := calcFeePerKb(refundFee, refundTx.SerializeSize())

	fmt.Printf("Refund fee: %v (%0.8f BTC/kB)\n\n", refundFee, refundFeePerKb)
	fmt.Printf("Refund transaction (%v):\n", &refundTxHash)
	fmt.Printf("%x\n\n", buf.Bytes())

	return promptPublishTx(c, refundTx, "refund")
}


func buildRefund(c *rpc.Client, contract []byte, contractTx *wire.MsgTx, feePerKb, minFeePerKb btcutil.Amount, chainParams *chaincfg.Params) (
	refundTx *wire.MsgTx, refundFee btcutil.Amount, err error) {

	contractP2SH, err := btcutil.NewAddressScriptHash(contract, chainParams)
	if err != nil {
		return nil, 0, err
	}
	contractP2SHPkScript, err := txscript.PayToAddrScript(contractP2SH)
	if err != nil {
		return nil, 0, err
	}

	contractTxHash := contractTx.TxHash()
	contractOutPoint := wire.OutPoint{Hash: contractTxHash, Index: ^uint32(0)}
	for i, o := range contractTx.TxOut {
		if bytes.Equal(o.PkScript, contractP2SHPkScript) {
			contractOutPoint.Index = uint32(i)
			break
		}
	}
	if contractOutPoint.Index == ^uint32(0) {
		return nil, 0, errors.New("contract tx does not contain a P2SH contract payment")
	}

	refundAddress, err := getRawChangeAddress(c)
	if err != nil {
		return nil, 0, fmt.Errorf("getrawchangeaddress: %v", err)
	}
	refundOutScript, err := txscript.PayToAddrScript(refundAddress)
	if err != nil {
		return nil, 0, err
	}

	pushes, err := txscript.ExtractAtomicSwapDataPushes(contract)
	if err != nil {
		// expected to only be called with good input
		panic(err)
	}

	refundAddr, err := btcutil.NewAddressPubKeyHash(pushes.RefundHash160[:], chainParams)
	if err != nil {
		return nil, 0, err
	}

	refundTx = wire.NewMsgTx(txVersion)
	refundTx.LockTime = uint32(pushes.LockTime)
	refundTx.AddTxOut(wire.NewTxOut(0, refundOutScript)) // amount set below
	refundSize := estimateRefundSerializeSize(contract, refundTx.TxOut)
	refundFee = txrules.FeeForSerializeSize(feePerKb, refundSize)
	refundTx.TxOut[0].Value = contractTx.TxOut[contractOutPoint.Index].Value - int64(refundFee)
	if txrules.IsDustOutput(refundTx.TxOut[0], minFeePerKb) {
		return nil, 0, fmt.Errorf("refund output value of %v is dust", btcutil.Amount(refundTx.TxOut[0].Value))
	}

	txIn := wire.NewTxIn(&contractOutPoint, nil, nil)
	txIn.Sequence = 0
	refundTx.AddTxIn(txIn)

	refundSig, refundPubKey, err := createSig(refundTx, 0, contract, refundAddr, c)
	if err != nil {
		return nil, 0, err
	}
	refundSigScript, err := refundP2SHContract(contract, refundSig, refundPubKey)
	if err != nil {
		return nil, 0, err
	}
	refundTx.TxIn[0].SignatureScript = refundSigScript

	if verify {
		e, err := txscript.NewEngine(contractTx.TxOut[contractOutPoint.Index].PkScript,
			refundTx, 0, txscript.StandardVerifyFlags, txscript.NewSigCache(10),
			txscript.NewTxSigHashes(refundTx), contractTx.TxOut[contractOutPoint.Index].Value)
		if err != nil {
			panic(err)
		}
		err = e.Execute()
		if err != nil {
			panic(err)
		}
	}

	return refundTx, refundFee, nil
}

// refundP2SHContract returns the signature script to refund a contract output
// using the contract author's signature after the locktime has been reached.
// This function assumes P2SH and appends the contract as the final data push.
func refundP2SHContract(contract, sig, pubkey []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddInt64(0)
	b.AddData(contract)
	return b.Script()
}