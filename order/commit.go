package order

import (
	"crypto/rsa"
	"math/big"

	"github.com/republicprotocol/republic-go/crypto"
)

type PedersenS struct {
	*big.Int
}

func (s PedersenS) Encrypt(pubKey rsa.PublicKey) ([]byte, error) {
	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{PublicKey: pubKey}}
	data := s.Bytes()
	return rsaKey.Encrypt(data)
}

type EncryptedPedersenS []byte

func (s EncryptedPedersenS) Decrypt(privKey *rsa.PrivateKey) (PedersenS, error) {
	rsaKey := crypto.RsaKey{PrivateKey: privKey}
	bs, err := rsaKey.Decrypt([]byte(s))
	if err != nil {
		return PedersenS{nil}, err
	}
	return PedersenS{big.NewInt(0).SetBytes(bs)}, nil
}

type PedersenCommitment struct {
	Index uint64 `json:"index"`

	Price         PedersenCoExpCommitment `json:"priceCommit"`
	Volume        PedersenCoExpCommitment `json:"volumeCommit"`
	MinimumVolume PedersenCoExpCommitment `json:"minimumVolumeCommit"`
}

type PedersenCoExpCommitment struct {
	Co  *big.Int `json:"co"`
	Exp *big.Int `json:"exp"`
}
