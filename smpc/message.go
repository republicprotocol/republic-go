package smpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"

	"github.com/republicprotocol/republic-go/shamir"
)

// ErrUnexpectedMessageType is returned when a message has an unexpected
// message type that cannot be marshaled/unmarshaled to/from binary.
var ErrUnexpectedMessageType = errors.New("unexpected message type")

// MessageID is used to associate asynchronous messages with each other. The
// value used usually correlates to an InstructionID.
type MessageID [32]byte

// MessageType distinguishes between the different messages that can be sent
// between SMPC nodes.
type MessageType byte

// MessageType values for messages passed between SMPC nodes.
const (
	MessageTypeJoinComponents         = MessageType(1)
	MessageTypeJoinComponentsResponse = MessageType(2)
)

// Message sent between SMPC nodes. These are for internal use by the SMPCs and
// are not needed by the users of an SMPC node.
type Message struct {
	MessageType

	*MessageJoinComponents
	*MessageJoinComponentsResponse
}

// MarshalBinary implements the stream.Message interface.
func (message *Message) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.MessageType); err != nil {
		return nil, err
	}

	var err error
	switch message.MessageType {
	case MessageTypeJoinComponents:
		log.Println("message writing:", message.MessageJoinComponents)
		if err := binary.Write(buf, binary.BigEndian, message.MessageJoinComponents); err != nil {
			return nil, err
		}
	case MessageTypeJoinComponentsResponse:
		if err := binary.Write(buf, binary.BigEndian, message.MessageJoinComponentsResponse); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnexpectedMessageType
	}

	return buf.Bytes(), err
}

// UnmarshalBinary implements the stream.Message interface.
func (message *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.MessageType); err != nil {
		return err
	}

	switch message.MessageType {
	case MessageTypeJoinComponents:
		message.MessageJoinComponents = new(MessageJoinComponents)
		return binary.Read(buf, binary.BigEndian, message.MessageJoinComponents)
	case MessageTypeJoinComponentsResponse:
		message.MessageJoinComponentsResponse = new(MessageJoinComponentsResponse)
		return binary.Read(buf, binary.BigEndian, message.MessageJoinComponentsResponse)
	default:
		return ErrUnexpectedMessageType
	}
}

// IsMessage implements the stream.Message interface.
func (message *Message) IsMessage() {}

// MessageJoinComponents declares the intent to reconstruct components by
// joining the shamir.Shares of those components. Support for multiple
// components in a single message reduces the number of messages, but each
// component should be handled separately and are not necessarily related.
type MessageJoinComponents struct {
	NetworkID
	Components
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (message *MessageJoinComponents) MarshalBinary() ([]byte, error) {
	log.Println("messageJoinCompontents writing:", message)

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.NetworkID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint64(len(message.Components))); err != nil {
		return nil, err
	}
	for _, component := range message.Components {
		if err := binary.Write(buf, binary.BigEndian, component.ComponentID); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, component.Share); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (message *MessageJoinComponents) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.NetworkID); err != nil {
		return err
	}
	numComponents := uint64(0)
	if err := binary.Read(buf, binary.BigEndian, &numComponents); err != nil {
		return err
	}
	message.Components = make(Components, numComponents)
	for i := uint64(0); i < numComponents; i++ {
		if err := binary.Read(buf, binary.BigEndian, &message.Components[i].ComponentID); err != nil {
			return err
		}
		if err := binary.Read(buf, binary.BigEndian, &message.Components[i].Share); err != nil {
			return err
		}
	}
	return nil
}

// MessageJoinComponentsResponse is sent in response to a MessageJoinComponents
// assuming that responding SMPC node agrees that the join should happen.
type MessageJoinComponentsResponse struct {
	NetworkID
	Components
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (message *MessageJoinComponentsResponse) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.NetworkID); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, uint64(len(message.Components))); err != nil {
		return nil, err
	}
	for _, component := range message.Components {
		if err := binary.Write(buf, binary.BigEndian, component.ComponentID); err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, component.Share); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (message *MessageJoinComponentsResponse) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.NetworkID); err != nil {
		return err
	}
	numComponents := uint64(0)
	if err := binary.Read(buf, binary.BigEndian, &numComponents); err != nil {
		return err
	}
	message.Components = make(Components, numComponents)
	for i := uint64(0); i < numComponents; i++ {
		if err := binary.Read(buf, binary.BigEndian, &message.Components[i].ComponentID); err != nil {
			return err
		}
		if err := binary.Read(buf, binary.BigEndian, &message.Components[i].Share); err != nil {
			return err
		}
	}
	return nil
}

// ComponentID for a shamir.Share to identify the component of which this
// shamir.Share is a part.
type ComponentID [32]byte

// Components is a set of distinct components.
type Components []Component

// Component is an ID associated with a shamir.Share.
type Component struct {
	ComponentID
	shamir.Share
}
