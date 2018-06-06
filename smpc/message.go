package smpc

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
)

// ErrUnexpectedMessageType is returned when a message has an unexpected
// message type that cannot be marshaled/unmarshaled to/from binary.
var ErrUnexpectedMessageType = errors.New("unexpected message type")

// MessageID is an abstract identity that is used to associate asynchronous
// messages with each other. The value of of the MessageID is dependent on the
// context of the Message. For example, when sending a MessageJoin the
// MessageID is set to the JoinID.
type MessageID [32]byte

// MessageType distinguishes between the different messages that can be sent
// between SMPC nodes.
type MessageType byte

// MessageType values for messages passed between SMPC nodes.
const (
	MessageTypeJoin         = MessageType(1)
	MessageTypeJoinResponse = MessageType(2)
)

// A Message is sent internally between nodes. It is not intended for direct
// use when interacting with an Smpcer.
type Message struct {
	MessageType

	MessageJoin         *MessageJoin
	MessageJoinResponse *MessageJoinResponse
}

// MarshalBinary implements the stream.Message interface.
func (message *Message) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.MessageType); err != nil {
		return nil, err
	}

	switch message.MessageType {
	case MessageTypeJoin:
		bytes, err := message.MessageJoin.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, bytes); err != nil {
			return nil, err
		}
	case MessageTypeJoinResponse:
		bytes, err := message.MessageJoinResponse.MarshalBinary()
		if err != nil {
			return nil, err
		}
		if err := binary.Write(buf, binary.BigEndian, bytes); err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnexpectedMessageType
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the stream.Message interface.
func (message *Message) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.MessageType); err != nil {
		return err
	}

	switch message.MessageType {
	case MessageTypeJoin:
		bytes, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		message.MessageJoin = new(MessageJoin)
		return message.MessageJoin.UnmarshalBinary(bytes)
	case MessageTypeJoinResponse:
		bytes, err := ioutil.ReadAll(buf)
		if err != nil {
			return err
		}
		message.MessageJoinResponse = new(MessageJoinResponse)
		return message.MessageJoinResponse.UnmarshalBinary(bytes)
	default:
		return ErrUnexpectedMessageType
	}
}

// IsMessage implements the stream.Message interface.
func (message *Message) IsMessage() {}

// A MessageJoin is used to broadcast a Join between nodes in the same network.
type MessageJoin struct {
	NetworkID
	Join Join
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (message *MessageJoin) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.NetworkID); err != nil {
		return nil, err
	}
	joinData, err := message.Join.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, joinData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (message *MessageJoin) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.NetworkID); err != nil {
		return err
	}
	joinData, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	return message.Join.UnmarshalBinary(joinData)
}

// A MessageJoinResponse is sent in response to a MessageJoin to return the
// Join owned by the responder, if such a Join exists.
type MessageJoinResponse struct {
	NetworkID
	Join Join
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (message *MessageJoinResponse) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, message.NetworkID); err != nil {
		return nil, err
	}
	joinData, err := message.Join.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, joinData); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (message *MessageJoinResponse) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.BigEndian, &message.NetworkID); err != nil {
		return err
	}
	joinData, err := ioutil.ReadAll(buf)
	if err != nil {
		return err
	}
	return message.Join.UnmarshalBinary(joinData)
}
