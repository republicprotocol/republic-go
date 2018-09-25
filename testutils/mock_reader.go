package testutils

import (
	"errors"

	"github.com/republicprotocol/republic-go/identity"
)

var alwaysFailError = errors.New("Error")

// Reader is a mock implementation of status.Reader interface.
type Reader struct {
	err error
}

func NewMockReader(alwaysFail bool) Reader {
	var err error
	if alwaysFail {
		err = alwaysFailError
	}
	return Reader{
		err: err,
	}
}

func (reader *Reader) Network() (string, error) {
	return "", reader.err
}

func (reader *Reader) MultiAddress() (identity.MultiAddress, error) {
	return identity.MultiAddress{}, reader.err
}

func (reader *Reader) PublicKey() ([]byte, error) {
	return []byte{}, reader.err
}

func (reader *Reader) Peers() (int, error) {
	return 0, reader.err
}

func (reader *Reader) EthereumNetwork() (string, error) {
	return "", reader.err
}

func (reader *Reader) EthereumAddress() (string, error) {
	return "", reader.err
}

func (reader *Reader) DarknodeRegistryAddress() (string, error) {
	return "", reader.err
}

func (reader *Reader) RewardVaultAddress() (string, error) {
	return "", reader.err
}

func (reader *Reader) InfuraURL() (string, error) {
	return "", reader.err
}

func (reader *Reader) Tokens() (map[string]string, error) {
	return map[string]string{}, reader.err
}
