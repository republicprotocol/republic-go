package adapter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/ingress"
	"github.com/republicprotocol/republic-go/order"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")
var ErrInvalidOrderFragmentIDLength = errors.New("invalid order fragment id length")
var ErrInvalidOrderFragmentValueLength = errors.New("invalid order fragment value length")
var ErrInvalidPoolHashLength = errors.New("invalid pool hash length")
var ErrEmptyOrderFragmentMapping = errors.New("empty order fragment mapping")

type IngressAdapter struct {
	ingress.Ingresser
}

func NewIngressAdapter(ingresser ingress.Ingresser) IngressAdapter {
	return IngressAdapter{
		Ingresser: ingresser,
	}
}

func (adapter *IngressAdapter) OpenOrder(signatureIn string, orderFragmentMappingIn OrderFragmentMapping) error {
	signature, err := adapter.adaptSignature(signatureIn)
	if err != nil {
		return err
	}

	orderID, orderFragmentMapping, err := adapter.adaptOrderFragmentMapping(orderFragmentMappingIn)
	if err != nil {
		return err
	}

	return adapter.Ingresser.OpenOrder(
		signature,
		orderID,
		orderFragmentMapping,
	)
}

func (adapter *IngressAdapter) CancelOrder(signatureIn string, orderIDIn string) error {
	signature, err := adapter.adaptSignature(signatureIn)
	if err != nil {
		return err
	}

	orderID, err := adapter.adaptOrderID(orderIDIn)
	if err != nil {
		return err
	}

	return adapter.Ingresser.CancelOrder(signature, orderID)
}

func (adapter *IngressAdapter) adaptSignature(signatureIn string) ([65]byte, error) {
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

func (adapter *IngressAdapter) adaptOrderID(orderIDIn string) (order.ID, error) {
	orderID := order.ID{}
	orderIDBytes, err := base64.StdEncoding.DecodeString(orderIDIn)
	if err != nil {
		return orderID, fmt.Errorf("cannot decode order id %v: %v", orderIDIn, err)
	}
	if len(orderIDBytes) != 32 {
		return orderID, ErrInvalidOrderIDLength
	}
	copy(orderID[:], orderIDBytes)
	return orderID, nil
}

func (adapter *IngressAdapter) adaptOrderFragmentID(orderFragmentIDIn string) (order.FragmentID, error) {
	orderFragmentID := order.FragmentID{}
	orderFragmentIDBytes, err := base64.StdEncoding.DecodeString(orderFragmentIDIn)
	if err != nil {
		return orderFragmentID, fmt.Errorf("cannot decode order fragment id %v: %v", orderFragmentIDIn, err)
	}
	if len(orderFragmentIDBytes) != 32 {
		return orderFragmentID, ErrInvalidOrderFragmentIDLength
	}
	copy(orderFragmentID[:], orderFragmentIDBytes)
	return orderFragmentID, nil
}

func (adapter *IngressAdapter) adaptOrderFragmentValue(value []string) (order.EncryptedFragmentValue, error) {
	var err error
	fragmentValue := order.EncryptedFragmentValue{}
	if len(value) != 2 {
		return fragmentValue, ErrInvalidOrderFragmentValueLength
	}
	fragmentValue.Co, err = base64.StdEncoding.DecodeString(value[0])
	if err != nil {
		return fragmentValue, err
	}
	fragmentValue.Exp, err = base64.StdEncoding.DecodeString(value[1])
	if err != nil {
		return fragmentValue, err
	}
	return fragmentValue, nil
}

func (adapter *IngressAdapter) adaptOrderFragment(orderFragmentIn OrderFragment) (ingress.OrderFragment, error) {
	var err error
	orderFragment := ingress.OrderFragment{EncryptedFragment: order.EncryptedFragment{}}
	orderFragment.Index = orderFragmentIn.Index
	orderFragment.EncryptedFragment.ID, err = adapter.adaptOrderFragmentID(orderFragmentIn.ID)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.OrderID, err = adapter.adaptOrderID(orderFragmentIn.OrderID)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.OrderType = orderFragmentIn.OrderType
	orderFragment.OrderParity = orderFragmentIn.OrderParity
	orderFragment.OrderExpiry = time.Unix(orderFragmentIn.OrderExpiry, 0)
	orderFragment.Tokens, err = base64.StdEncoding.DecodeString(orderFragmentIn.Tokens)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.Price, err = adapter.adaptOrderFragmentValue(orderFragmentIn.Price)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.Volume, err = adapter.adaptOrderFragmentValue(orderFragmentIn.Volume)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.MinimumVolume, err = adapter.adaptOrderFragmentValue(orderFragmentIn.MinimumVolume)
	if err != nil {
		return orderFragment, err
	}
	return orderFragment, nil
}

func (adapter *IngressAdapter) adaptOrderFragmentMapping(orderFragmentMappingIn OrderFragmentMapping) (order.ID, ingress.OrderFragmentMapping, error) {
	orderID := order.ID{}
	orderFragmentMapping := ingress.OrderFragmentMapping{}

	// Decode order ID
	for _, value := range orderFragmentMappingIn {
		var err error
		if orderID, err = adapter.adaptOrderID(value[0].OrderID); err != nil {
			return orderID, orderFragmentMapping, err
		}
		break
	}

	// Decode order fragments
	for key, orderFragmentsIn := range orderFragmentMappingIn {
		hashBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return orderID, orderFragmentMapping, fmt.Errorf("cannot decode pool hash %v: %v", key, err)
		}
		hash := [32]byte{}
		if len(hashBytes) != 32 {
			return orderID, orderFragmentMapping, ErrInvalidPoolHashLength
		}
		copy(hash[:], hashBytes)
		orderFragmentMapping[hash] = make([]ingress.OrderFragment, 0, len(orderFragmentsIn))
		for _, orderFragmentIn := range orderFragmentsIn {
			orderFragment, err := adapter.adaptOrderFragment(orderFragmentIn)
			if err != nil {
				return orderID, orderFragmentMapping, err
			}
			orderFragmentMapping[hash] = append(orderFragmentMapping[hash], orderFragment)
		}
	}
	return orderID, orderFragmentMapping, nil
}
