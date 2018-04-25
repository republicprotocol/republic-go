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
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
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

func (crypter *WeakCrypter) Encrypt(plaintext []byte) ([]byte, error) {
	return plaintext, nil
}

func (crypter *WeakCrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	return ciphertext, nil
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

func (crypter *ErrCrypter) Encrypt(plaintext []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}

func (crypter *ErrCrypter) Decrypt(ciphertext []byte) ([]byte, error) {
	return []byte{}, ErrUnimplemented
}
