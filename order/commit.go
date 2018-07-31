package order

import (
	"crypto/rsa"
	"math/big"

	"github.com/republicprotocol/republic-go/crypto"
)

type Commitment struct {
	Index uint64 `json:"index"`

	Price         CoExpCommitment `json:"priceCommit"`
	Volume        CoExpCommitment `json:"volumeCommit"`
	MinimumVolume CoExpCommitment `json:"minimumVolumeCommit"`
}

type CoExpCommitment struct {
	Co  *big.Int `json:"co"`
	Exp *big.Int `json:"exp"`
}

type Blindings []Blinding

type Blinding struct {
	*big.Int
}

func (b Blinding) Encrypt(pubKey rsa.PublicKey) ([]byte, error) {
	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{PublicKey: pubKey}}
	data := b.Bytes()
	return rsaKey.Encrypt(data)
}

type EncryptedBlinding []byte

func (b EncryptedBlinding) Decrypt(privKey *rsa.PrivateKey) (Blinding, error) {
	rsaKey := crypto.RsaKey{PrivateKey: privKey}
	bs, err := rsaKey.Decrypt([]byte(b))
	if err != nil {
		return Blinding{nil}, err
	}
	return Blinding{big.NewInt(0).SetBytes(bs)}, nil
}
