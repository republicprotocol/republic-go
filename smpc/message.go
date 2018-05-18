package smpc

import (
	"encoding/json"

	"github.com/republicprotocol/republic-go/smpc/delta"
)

type DeltaFragmentMessage struct {
	PeersID       []byte
	DeltaFragment delta.Fragment
}

func (message *DeltaFragmentMessage) IsMessage() {
}

func (message *DeltaFragmentMessage) MarshalBinary() ([]byte, error) {
	return json.Marshal(message)
}

func (message *DeltaFragmentMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, message)
}

//type DeltaFragmentMessage struct {
//	DeltaFragment delta.Fragment
//	PeersID       [32]byte
//}
//
//func (message DeltaFragmentMessage) MarshalBinary() (data []byte, err error) {
//	buf := new(bytes.Buffer)
//	binary.Write(buf, binary.BigEndian, message.DeltaFragment)
//	binary.Write(buf, binary.BigEndian, message.PeersID)
//	return buf.Bytes(), nil
//}
//
//func (message DeltaFragmentMessage) UnmarshalBinary(data []byte) error {
//	if data == nil || len(data) == 0 {
//		return ErrUnmarshalNilBytes
//	}
//	buf := bytes.NewBuffer(data)
//	binary.Read(buf, binary.BigEndian, &message.DeltaFragment)
//	binary.Read(buf, binary.BigEndian, &message.PeersID)
//	return nil
//}
//
//func (message DeltaFragmentMessage) IsMessage() {
//}
