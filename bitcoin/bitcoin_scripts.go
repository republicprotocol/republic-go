package bitcoin

import (
	"github.com/btcsuite/btcd/txscript"
	"golang.org/x/crypto/ripemd160"
)

const (
	// redeemAtomicSwapSigScriptSize is the worst case (largest) serialize size
	// of a transaction input script to redeem the atomic swap contract.  This
	// does not include final push for the contract itself.
	//
	//   - OP_DATA_73
	//   - 72 bytes DER signature + 1 byte sighash
	//   - OP_DATA_33
	//   - 33 bytes serialized compressed pubkey
	//   - OP_DATA_32
	//   - 32 bytes secret
	//   - OP_TRUE
	redeemAtomicSwapSigScriptSize = 1 + 73 + 1 + 33 + 1 + 32 + 1

	// refundAtomicSwapSigScriptSize is the worst case (largest) serialize size
	// of a transaction input script that refunds a P2SH atomic swap output.
	// This does not include final push for the contract itself.
	//
	//   - OP_DATA_73
	//   - 72 bytes DER signature + 1 byte sighash
	//   - OP_DATA_33
	//   - 33 bytes serialized compressed pubkey
	//   - OP_FALSE
	refundAtomicSwapSigScriptSize = 1 + 73 + 1 + 33 + 1
)

/*
Bitcoin AtomicSwap Script: Alice is trying to do an atomic swap with bob.

OP_IF
	OP_SHA256
	<secret_hash>
	OP_EQUALVERIFY
	OP_DUP
	OP_HASH160
	<pubkey_hash_bob>
OP_ELSE
	<lock_time>
	OP_CHECKLOCKTIMEVERIFY
	OP_DROP
	OP_HASH160
	<pubKey_hash_alice>
OP_ENDIF
OP_EQUALVERIFY
OP_CHECKSIG

*/

func atomicSwapContract(pkhMe, pkhThem *[ripemd160.Size]byte, locktime int64, secretHash []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()

	b.AddOp(txscript.OP_IF)
	{
		b.AddOp(txscript.OP_SHA256)
		b.AddData(secretHash)
		b.AddOp(txscript.OP_EQUALVERIFY)
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhThem[:])
	}
	b.AddOp(txscript.OP_ELSE)
	{
		b.AddInt64(locktime)
		b.AddOp(txscript.OP_CHECKLOCKTIMEVERIFY)
		b.AddOp(txscript.OP_DROP)
		b.AddOp(txscript.OP_DUP)
		b.AddOp(txscript.OP_HASH160)
		b.AddData(pkhMe[:])
	}
	b.AddOp(txscript.OP_ENDIF)
	b.AddOp(txscript.OP_EQUALVERIFY)
	b.AddOp(txscript.OP_CHECKSIG)

	return b.Script()
}

/*
Bitcoin Refund Script: Alice is trying to get refunded

<Signature>
<PublicKey>
<False>(Int 0)
<Contract>
*/
func refundP2SHContract(contract, sig, pubkey []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddInt64(0)
	b.AddData(contract)
	return b.Script()
}

/*
Bitcoin Refund Script: Bob is trying to redeem and get his bitcoins.

<Signature>
<PublicKey>
<Secret>
<True>(Int 1)
<Contract>
*/

func redeemP2SHContract(contract, sig, pubkey, secret []byte) ([]byte, error) {
	b := txscript.NewScriptBuilder()
	b.AddData(sig)
	b.AddData(pubkey)
	b.AddData(secret)
	b.AddInt64(1)
	b.AddData(contract)
	return b.Script()
}
