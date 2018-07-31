package shamir

import (
	"crypto/rsa"
	"math/big"

	"github.com/republicprotocol/republic-go/crypto"
)

type Blinds []Blind

type Blind struct {
	*big.Int
}

func (b Blind) Encrypt(pubKey rsa.PublicKey) ([]byte, error) {
	rsaKey := crypto.RsaKey{PrivateKey: &rsa.PrivateKey{PublicKey: pubKey}}
	data := b.Bytes()
	return rsaKey.Encrypt(data)
}

type EncryptedBlind []byte

func (b EncryptedBlind) Decrypt(privKey *rsa.PrivateKey) (Blind, error) {
	rsaKey := crypto.RsaKey{PrivateKey: privKey}
	bs, err := rsaKey.Decrypt([]byte(b))
	if err != nil {
		return Blind{nil}, err
	}
	return Blind{big.NewInt(0).SetBytes(bs)}, nil
}
