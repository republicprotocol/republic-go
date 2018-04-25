package crypto

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/binary"
	"math/big"
)

// RsaKey for encrypting and decrypting sensitive data that must be transported
// between actors in the network.
type RsaKey struct {
	*rsa.PrivateKey
}

// RandomRsaKey using 2048 bits, with precomputed values for improved
// performance.
func RandomRsaKey() (RsaKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	privateKey.Precompute()
	if err != nil {
		return RsaKey{}, err
	}
	return RsaKey{
		PrivateKey: privateKey,
	}, nil
}

// NewRsaKey returns an RsaKey from an existing private key. It does not verify
// that the private key was generated correctly. It precomputes values for
// improved performance
func NewRsaKey(privateKey *rsa.PrivateKey) RsaKey {
	privateKey.Precompute()
	return RsaKey{
		PrivateKey: privateKey,
	}
}

// Encrypt a plain text and return the cipher text.
func (key *RsaKey) Encrypt(plaintext []byte) ([]byte, error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &key.PublicKey, plaintext)
}

// Decrypt a cipher text and return the plain text.
func (key *RsaKey) Decrypt(ciphertext []byte) ([]byte, error) {
	return rsa.DecryptPKCS1v15(rand.Reader, key.PrivateKey, ciphertext)
}

// PublicKeyToBytes by writing E as an int64, and then writing N as a stream
// of bytes. Big endian encoding is used.
func PublicKeyToBytes(publicKey *rsa.PublicKey) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, int64(publicKey.E))
	binary.Write(buf, binary.BigEndian, publicKey.N.Bytes())
	return buf.Bytes()
}

// BytesToPublicKey converts a byte array to a RSA Public Key
// It returns pointer to a RSA public key
func BytesToPublicKey(data []byte) rsa.PublicKey {
	r := bytes.NewReader(data)
	e := int64(0)
	binary.Read(r, binary.BigEndian, &e)
	n := make([]byte, r.Len())
	binary.Read(r, binary.BigEndian, n)
	var publicKey rsa.PublicKey
	publicKey.E = int(e)
	publicKey.N = big.NewInt(0).SetBytes(n)
	return publicKey
}
