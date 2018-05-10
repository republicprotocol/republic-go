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

type OrderFragmentMapping map[string][]order.Fragment

type RelayAdapter struct {
	relay.Relay
}

func (adapter *RelayAdapter) OpenOrder(signatureIn string, orderFragmentMappingIn OrderFragmentMapping) error {

	orderFragment, err := adapter.firstOrderFragmentInMapping(orderFragmentMappingIn)
	if err != nil {
		return err
	}

	signature, err := adapter.adaptSignature(signatureIn)
	if err != nil {
		return err
	}

	orderFragmentMapping, err := adapter.adaptOrderFragmentMapping(orderFragmentMappingIn)
	if err != nil {
		return err
	}

	return adapter.Relay.OpenOrder(
		signature,
		orderFragment.OrderID,
		orderFragment.OrderType,
		orderFragment.OrderParity,
		orderFragment.OrderExpiry.Unix(),
		orderFragmentMapping,
	)
}

func (adapter *RelayAdapter) firstOrderFragmentInMapping(orderFragmentMapping OrderFragmentMapping) (*order.Fragment, error) {
	// Adapt the order data
	var orderFragment *order.Fragment
	for _, orderFragments := range orderFragmentMapping {
		if len(orderFragments) > 0 {
			orderFragment = &orderFragments[0]
			break
		}
	}
	if orderFragment == nil {
		return nil, ErrEmptyOrderFragmentMapping
	}
	return orderFragment, nil
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

func (adapter *RelayAdapter) adaptOrderFragmentMapping(orderFragmentMappingIn OrderFragmentMapping) (relay.OrderFragmentMapping, error) {
	orderFragmentMapping := relay.OrderFragmentMapping{}
	for key, value := range orderFragmentMappingIn {
		hashBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return orderFragmentMapping, fmt.Errorf("cannot decode pool hash %v: %v", key, err)
		}
		hash := [32]byte{}
		if len(hashBytes) != 32 {
			return orderFragmentMapping, ErrInvalidPoolHashLength
		}
		copy(hash[:], hashBytes)
		orderFragmentMapping[hash] = value
	}
	return orderFragmentMapping, nil
}
