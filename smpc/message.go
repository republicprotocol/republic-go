package smpc

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/republicprotocol/republic-go/shamir"
)

// ErrUnexpectedMessageType is returned when a message has an unexpected
// message type that cannot be marshaled/unmarshaled to/from binary.
var ErrUnexpectedMessageType = errors.New("unexpected message type")

type messageType int8

const (
	messageTypeJ = 1
)

type message struct {
	messageType

	*messageJ
}

// MarshalBinary implements the stream.Message interface.
func (msg *message) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, msg.messageType); err != nil {
		return []byte{}, err
	}

	var err error
	switch msg.messageType {
	case messageTypeJ:
		err = binary.Write(buf, binary.BigEndian, msg.messageJ)
	default:
		return []byte{}, ErrUnexpectedMessageType
	}
	return buf.Bytes(), err
}

// UnmarshalBinary implements the stream.Message interface.
func (msg *message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &msg.messageType); err != nil {
		return err
	}

	switch msg.messageType {
	case messageTypeJ:
		return binary.Read(buf, binary.BigEndian, msg.messageJ)
	default:
		return ErrUnexpectedMessageType
	}
}

// IsMessage implements the stream.Message interface.
func (msg *message) IsMessage() {}

type messageJ struct {
	InstID    [32]byte
	NetworkID [32]byte
	Share     shamir.Share
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (msg *messageJ) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, msg.InstID); err != nil {
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.BigEndian, msg.NetworkID); err != nil {
		return []byte{}, err
	}
	if err := binary.Write(buf, binary.BigEndian, msg.Share); err != nil {
		return []byte{}, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (msg *messageJ) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &msg.InstID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.NetworkID); err != nil {
		return err
	}
	if err := binary.Read(buf, binary.BigEndian, &msg.Share); err != nil {
		return err
	}
	return nil
}
