package arc_bitcoin

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/republicprotocol/republic-go/blockchain/bitcoin"
)

func readSecret(connection client.Connection, redemptionTxBytes, secretHash []byte) ([32]byte, error) {
	var redemptionTx wire.MsgTx
	err := redemptionTx.Deserialize(bytes.NewReader(redemptionTxBytes))
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to decode redemption transaction: %v", err)
	}

	if len(secretHash) != sha256.Size {
		return [32]byte{}, errors.New("secret hash has wrong size")
	}

	for _, in := range redemptionTx.TxIn {
		pushes, err := txscript.PushedData(in.SignatureScript)
		if err != nil {
			return [32]byte{}, err
		}
		for _, push := range pushes {
			if bytes.Equal(sha256Hash(push), secretHash) {
				var secret [32]byte
				for i := 0; i < 32; i++ {
					secret[i] = push[i]
				}
				return secret, nil
			}
		}
	}
	return [32]byte{}, errors.New("transaction does not contain the secret")
}
