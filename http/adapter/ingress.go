package adapter

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/republicprotocol/republic-go/ingress"
	"github.com/republicprotocol/republic-go/order"
)

var ErrInvalidSignatureLength = errors.New("invalid signature length")
var ErrInvalidOrderIDLength = errors.New("invalid order id length")
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
	orderIDBytes, err := base64.StdEncoding.DecodeString(orderIDIn)
	if err != nil {
		return order.ID{}, fmt.Errorf("cannot decode order id %v: %v", orderIDIn, err)
	}
	if len(orderIDBytes) != 32 {
		return order.ID{}, ErrInvalidOrderIDLength
	}
	return order.ID(orderIDBytes), nil
}

func (adapter *IngressAdapter) adaptOrderFragmentMapping(orderFragmentMappingIn OrderFragmentMapping) (order.ID, ingress.OrderFragmentMapping, error) {
	orderID := order.ID{}
	orderFragmentMapping := ingress.OrderFragmentMapping{}
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
