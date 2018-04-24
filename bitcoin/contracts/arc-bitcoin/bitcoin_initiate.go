package arc_bitcoin

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/btcsuite/btcutil"
	"github.com/republicprotocol/republic-go/bitcoin/client"
)

func initiate(connection client.Connection, participantAddress string, value int64, hash []byte, lockTime int64) (BitcoinData, error) {

	cp2Addr, err := btcutil.DecodeAddress(participantAddress, connection.ChainParams)
	if err != nil {
		return BitcoinData{}, fmt.Errorf("failed to decode participant address: %v", err)
	}
	if !cp2Addr.IsForNet(connection.ChainParams) {
		return BitcoinData{}, fmt.Errorf("participant address is not "+
			"intended for use on %v", connection.ChainParams.Name)
	}
	cp2AddrP2PKH, ok := cp2Addr.(*btcutil.AddressPubKeyHash)
	if !ok {
		return BitcoinData{}, errors.New("participant address is not P2PKH")
	}

	b, err := buildContract(connection, &contractArgs{
		them:       cp2AddrP2PKH,
		amount:     value,
		locktime:   lockTime,
		secretHash: hash,
	})
	if err != nil {
		return BitcoinData{}, err
	}

	var contractBuf bytes.Buffer
	contractBuf.Grow(b.contractTx.SerializeSize())
	b.contractTx.Serialize(&contractBuf)

	var refundBuf bytes.Buffer
	refundBuf.Grow(b.refundTx.SerializeSize())
	b.refundTx.Serialize(&refundBuf)

	txHash, err := connection.PromptPublishTx(b.contractTx, "contract")

	if err != nil {
		return BitcoinData{}, err
	}

	connection.WaitForConfirmations(txHash, 1)

	refundTx := *b.refundTx
	return BitcoinData{
		Contract:       b.contract,
		ContractHash:   b.contractP2SH.EncodeAddress(),
		ContractTx:     contractBuf.Bytes(),
		ContractTxHash: b.contractTxHash.CloneBytes(),
		RefundTx:       refundBuf.Bytes(),
		RefundTxHash:   refundTx.TxHash(),
	}, nil
}
