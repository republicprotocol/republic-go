package bitcoin

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func ExtractSecret(redemptionTxBytes, secretHash []byte, rpcUser string, rpcPass string) (Error error, secret []byte) {
	var redemptionTx wire.MsgTx
	err := redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return fmt.Errorf("failed to decode redemption transaction: %v", err), []byte{}
	}

	if len(secretHash) != sha256.Size {
		return errors.New("secret hash has wrong size"), []byte{}
	}

	for _, in := range redemptionTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return err, []byte{}
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), secretHash) {
				return nil, push
			}
		}
	}
	return errors.New("transaction does not contain the secret"), []byte{}
}
