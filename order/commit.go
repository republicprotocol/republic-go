package order

import (
	"crypto/rsa"
	"math/big"

	"github.com/republicprotocol/republic-go/crypto"
)

var G = big.NewInt(0)
var H = big.NewInt(0)
var P = big.NewInt(0)

type Commitments map[uint64]Commitment

type Commitment struct {
	*big.Int
}

func NewCommitment(secret *big.Int, blinding *big.Int) Commitment {

	Gx := big.NewInt(0).Exp(G, secret, P)
	Hs := big.NewInt(0).Exp(H, blinding, P)
	C := big.NewInt(0).Mul(Gx, Hs)
	C.Mod(C, P)

	return Commitment{Int: C}
}

type CommitmentSet struct {
	PriceCo          Commitment `json:"priceCo"`
	PriceExp         Commitment `json:"priceExp"`
	VolumeCo         Commitment `json:"volumeCo"`
	VolumeExp        Commitment `json:"volumeExp"`
	MinimumVolumeCo  Commitment `json:"minimumVolumeCo"`
	MinimumVolumeExp Commitment `json:"minimumVolumeExp"`
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
