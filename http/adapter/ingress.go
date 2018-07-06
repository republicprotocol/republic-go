package adapter

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/republicprotocol/republic-go/ingress"
	"github.com/republicprotocol/republic-go/order"
)

// ErrInvalidSignatureLength is returned when a signature does not have the
// required length of 65 bytes.
var ErrInvalidSignatureLength = errors.New("invalid signature length")

// ErrInvalidOrderIDLength is returned when an order ID does not have the
// required length of 32 bytes.
var ErrInvalidOrderIDLength = errors.New("invalid order id length")

// ErrInvalidOrderFragmentIDLength is returned when an order fragment ID does
// not have the required length of 32 bytes.
var ErrInvalidOrderFragmentIDLength = errors.New("invalid order fragment id length")

// ErrInvalidEncryptedCoExpShareLength is returned when an encrypted co-exp
// share does not contain exactly 2 encrypted values, an encrypted co and an
// encrypted exp.
var ErrInvalidEncryptedCoExpShareLength = errors.New("invalid encrypted co-exp share length")

// ErrInvalidPodHashLength is returned when a pod hash does not have the
// required length of 32 bytes.
var ErrInvalidPodHashLength = errors.New("invalid pod hash length")

// ErrEmptyOrderFragmentMapping is returned when an OrderFragmentMapping does
// not store any OrderFragments.
var ErrEmptyOrderFragmentMapping = errors.New("empty order fragment mapping")

type IngressAdapter struct {
	ingress.Ingress
}

func NewIngressAdapter(ingress ingress.Ingress) IngressAdapter {
	return IngressAdapter{
		Ingress: ingress,
	}
}

func (adapter *IngressAdapter) OpenOrder(signatureIn string, orderFragmentMappingsIn OrderFragmentMappings) error {
	signature, err := UnmarshalSignature(signatureIn)
	if err != nil {
		return err
	}

	orderID, orderFragmentMappings, err := UnmarshalOrderFragmentMappings(orderFragmentMappingsIn)
	if err != nil {
		return err
	}

	return adapter.Ingress.OpenOrder(
		signature,
		orderID,
		orderFragmentMappings,
	)
}

func (adapter *IngressAdapter) CancelOrder(signatureIn string, orderIDIn string) error {
	signature, err := UnmarshalSignature(signatureIn)
	if err != nil {
		return err
	}

	orderID, err := UnmarshalOrderID(orderIDIn)
	if err != nil {
		return err
	}

	return adapter.Ingress.CancelOrder(signature, orderID)
}

func MarshalSignature(signatureIn [65]byte) string {
	return base64.StdEncoding.EncodeToString(signatureIn[:])
}

func UnmarshalSignature(signatureIn string) ([65]byte, error) {
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

func MarshalOrderID(orderIDIn order.ID) string {
	return base64.StdEncoding.EncodeToString(orderIDIn[:])
}

func UnmarshalOrderID(orderIDIn string) (order.ID, error) {
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

func MarshalOrderFragmentID(orderFragmentIDIn order.FragmentID) string {
	return base64.StdEncoding.EncodeToString(orderFragmentIDIn[:])
}

func UnmarshalOrderFragmentID(orderFragmentIDIn string) (order.FragmentID, error) {
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

func MarshalEncryptedCoExpShare(valueIn order.EncryptedCoExpShare) []string {
	return []string{
		base64.StdEncoding.EncodeToString(valueIn.Co),
		base64.StdEncoding.EncodeToString(valueIn.Exp),
	}
}

func UnmarshalEncryptedCoExpShare(valueIn []string) (order.EncryptedCoExpShare, error) {
	var err error
	value := order.EncryptedCoExpShare{}
	if len(valueIn) != 2 {
		return value, ErrInvalidEncryptedCoExpShareLength
	}
	value.Co, err = base64.StdEncoding.DecodeString(valueIn[0])
	if err != nil {
		return value, err
	}
	value.Exp, err = base64.StdEncoding.DecodeString(valueIn[1])
	if err != nil {
		return value, err
	}
	return value, nil
}

func MarshalOrderFragment(orderFragmentIn ingress.OrderFragment) OrderFragment {
	orderFragment := OrderFragment{}
	orderFragment.Index = orderFragmentIn.Index
	orderFragment.OrderID = MarshalOrderID(orderFragmentIn.OrderID)
	orderFragment.OrderType = orderFragmentIn.OrderType
	orderFragment.OrderParity = orderFragmentIn.OrderParity
	orderFragment.OrderSettlement = orderFragmentIn.OrderSettlement
	orderFragment.OrderExpiry = orderFragmentIn.OrderExpiry.Unix()
	orderFragment.ID = MarshalOrderFragmentID(orderFragmentIn.ID)
	orderFragment.Tokens = base64.StdEncoding.EncodeToString(orderFragmentIn.Tokens)
	orderFragment.Price = MarshalEncryptedCoExpShare(orderFragmentIn.Price)
	orderFragment.Volume = MarshalEncryptedCoExpShare(orderFragmentIn.Volume)
	orderFragment.MinimumVolume = MarshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume)
	orderFragment.Nonce = base64.StdEncoding.EncodeToString(orderFragmentIn.Nonce)
	return orderFragment
}

func UnmarshalOrderFragment(orderFragmentIn OrderFragment) (ingress.OrderFragment, error) {
	var err error
	orderFragment := ingress.OrderFragment{EncryptedFragment: order.EncryptedFragment{}}
	orderFragment.Index = orderFragmentIn.Index
	orderFragment.EncryptedFragment.ID, err = UnmarshalOrderFragmentID(orderFragmentIn.ID)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.OrderID, err = UnmarshalOrderID(orderFragmentIn.OrderID)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.OrderType = orderFragmentIn.OrderType
	orderFragment.OrderParity = orderFragmentIn.OrderParity
	orderFragment.OrderSettlement = orderFragmentIn.OrderSettlement
	orderFragment.OrderExpiry = time.Unix(orderFragmentIn.OrderExpiry, 0)
	orderFragment.Tokens, err = base64.StdEncoding.DecodeString(orderFragmentIn.Tokens)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.Price, err = UnmarshalEncryptedCoExpShare(orderFragmentIn.Price)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.Volume, err = UnmarshalEncryptedCoExpShare(orderFragmentIn.Volume)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.MinimumVolume, err = UnmarshalEncryptedCoExpShare(orderFragmentIn.MinimumVolume)
	if err != nil {
		return orderFragment, err
	}
	orderFragment.Nonce, err = base64.StdEncoding.DecodeString(orderFragmentIn.Nonce)
	return orderFragment, nil
}

func UnmarshalOrderFragmentMapping(orderFragmentMappingIn OrderFragmentMapping) (order.ID, ingress.OrderFragmentMapping, error) {
	orderID := order.ID{}
	orderFragmentMapping := ingress.OrderFragmentMapping{}

	// Decode order ID
	for _, values := range orderFragmentMappingIn {
		var err error
		foundOrderID := false
		for _, value := range values {
			if orderID, err = UnmarshalOrderID(value.OrderID); err != nil {
				return orderID, orderFragmentMapping, err
			}
			foundOrderID = true
			break
		}
		if foundOrderID {
			break
		}
	}

	// Decode order fragments
	for key, orderFragmentsIn := range orderFragmentMappingIn {
		hashBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return orderID, orderFragmentMapping, fmt.Errorf("cannot decode pool hash %v: %v", key, err)
		}
		hash := [32]byte{}
		if len(hashBytes) != 32 {
			return orderID, orderFragmentMapping, ErrInvalidPodHashLength
		}
		copy(hash[:], hashBytes)
		orderFragmentMapping[hash] = make([]ingress.OrderFragment, 0, len(orderFragmentsIn))
		for _, orderFragmentIn := range orderFragmentsIn {
			orderFragment, err := UnmarshalOrderFragment(orderFragmentIn)
			if err != nil {
				return orderID, orderFragmentMapping, err
			}
			orderFragmentMapping[hash] = append(orderFragmentMapping[hash], orderFragment)
		}
	}
	return orderID, orderFragmentMapping, nil
}

func UnmarshalOrderFragmentMappings(orderFragmentMappingsIn OrderFragmentMappings) (order.ID, ingress.OrderFragmentMappings, error) {

}
