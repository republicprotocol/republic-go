package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// ErrMalformedPadding is returned when padding cannot be stripped during
// decryption.
var ErrMalformedPadding = errors.New("malformed padding")

// ErrMalformedCipherText is returned when a cipher text is not a multiple of
// the block size.
var ErrMalformedCipherText = errors.New("malformed cipher text")

type AESCipher struct {
	secret []byte
}

func RandomAESCipher() (AESCipher, error) {
	secret := [16]byte{}
	if _, err := io.ReadFull(rand.Reader, secret[:]); err != nil {
		return AESCipher{}, err
	}
	return AESCipher{secret: secret[:]}, nil
}

func NewAESCipher(secret []byte) AESCipher {
	return AESCipher{secret: secret}
}

func (key *AESCipher) Encrypt(plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.secret)
	if err != nil {
		return nil, err
	}

	paddedPlainText := pad(plainText)
	cipherText := make([]byte, aes.BlockSize+len(paddedPlainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(cipherText[aes.BlockSize:], []byte(paddedPlainText))
	return cipherText, nil
}

func (key *AESCipher) Decrypt(cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.secret)
	if err != nil {
		return nil, err
	}
	if (len(cipherText) % aes.BlockSize) != 0 {
		return nil, ErrMalformedCipherText
	}

	iv := cipherText[:aes.BlockSize]
	message := cipherText[aes.BlockSize:]

	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(message, message)
	return strip(message)
}

func pad(src []byte) []byte {
	p := aes.BlockSize - len(src)%aes.BlockSize
	padding := bytes.Repeat([]byte{byte(p)}, p)
	return append(src, padding...)
}

func strip(src []byte) ([]byte, error) {
	length := len(src)
	p := int(src[length-1])
	if p > length {
		return nil, ErrMalformedPadding
	}
	return src[:(length - p)], nil
}
