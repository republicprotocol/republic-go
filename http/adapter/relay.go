package adapter

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/republicprotocol/republic-go/order"
	"github.com/republicprotocol/republic-go/relay"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidPoolHashLength = errors.New("invalid pool hash length")
var ErrEmptyOrderFragmentMapping = errors.New("empty order fragment mapping")

type RelayAdapter struct {
	relay.Relayer
}

func (adapter *RelayAdapter) OpenOrder(signatureIn string, orderFragmentMappingIn OrderFragmentMapping) error {
	signature, err := adapter.adaptSignature(signatureIn)
	if err != nil {
		return err
	}

	orderID, orderFragmentMapping, err := adapter.adaptOrderFragmentMapping(orderFragmentMappingIn)
	if err != nil {
		return err
	}

	return adapter.Relayer.OpenOrder(
		signature,
		orderID,
		orderFragmentMapping,
	)
}

func (adapter *RelayAdapter) CancelOrder(signatureIn string, orderID order.ID) error {
	signature, err := adapter.adaptSignature(signatureIn)
	if err != nil {
		return err
	}

	return adapter.Relayer.CancelOrder(signature, orderID)
}

func (adapter *RelayAdapter) adaptSignature(signatureIn string) ([65]byte, error) {
	signature := [65]byte{}
	signatureBytes, err := base64.StdEncoding.DecodeString(signatureIn)
	if err != nil {
		return signature, fmt.Errorf("cannot decode signature %v: %v", signatureIn, err)
	}
	if len(signatureBytes) != 65 {
		return signature, ErrInvalidSignatureLength
	}
	copy(signature[:], signatureBytes)
	return signature, nil
}

func (adapter *RelayAdapter) adaptOrderFragmentMapping(orderFragmentMappingIn OrderFragmentMapping) (order.ID, relay.OrderFragmentMapping, error) {
	orderID := order.ID{}
	orderFragmentMapping := relay.OrderFragmentMapping{}
	for key, value := range orderFragmentMappingIn {
		if len(orderID) == 0 && len(value) > 0 {
			orderID = value[0].OrderID
		}
		hashBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return orderID, orderFragmentMapping, fmt.Errorf("cannot decode pool hash %v: %v", key, err)
		}
		hash := [32]byte{}
		if len(hashBytes) != 32 {
			return orderID, orderFragmentMapping, ErrInvalidPoolHashLength
		}
		copy(hash[:], hashBytes)
		orderFragmentMapping[hash] = value
	}
	return orderID, orderFragmentMapping, nil
}
