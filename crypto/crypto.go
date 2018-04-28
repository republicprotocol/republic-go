package crypto

import (
	"errors"
)

var ErrUnimplemented = errors.New("unimplemented")
var ErrInvalidSignature = errors.New("invalid signature")

type Hasher interface {
	Hash() []byte
}

type Signer interface {
	Sign(Hasher) ([]byte, error)
}

type Verifier interface {
	Verify(Hasher, []byte) error
}

type Encrypter interface {
	Encrypt(string, []byte) ([]byte, error)
	Decrypt(string, []byte) ([]byte, error)
}

type Crypter interface {
	Signer
	Verifier
	Encrypter
}

type WeakCrypter struct{}

func NewWeakCrypter() WeakCrypter {
	return WeakCrypter{}
}

func (crypter *WeakCrypter) Sign(hasher Hasher) ([]byte, error) {
	return hasher.Hash(), nil
}

func (crypter *WeakCrypter) Verify(hasher Hasher, signature []byte) error {
	return nil
}

func (crypter *WeakCrypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return plainText, nil
}

func (crypter *WeakCrypter) Decrypt(addr string, cipherText []byte) ([]byte, error) {
	return cipherText, nil
}

type ErrCrypter struct{}

func NewErrCrypter() ErrCrypter {
	return ErrCrypter{}
}

func (crypter *ErrCrypter) Signer(hasher Hasher) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

func (crypter *ErrCrypter) Verifier(hasher Hasher, signature []byte) error {
	return ErrUnimplemented
}

func (crypter *ErrCrypter) Encrypt(addr string, plainText []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

func (crypter *ErrCrypter) Decrypt(addr string, cipherText []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}
