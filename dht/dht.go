package dht

import (
	"container/list"
	"fmt"
	"github.com/pkg/errors"
	"github.com/republicprotocol/republic/crypto"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = 20
	IDLengthInBits = IDLength * 8
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
type ID string


// NewID creates a new set of public/private SECP256K1 key pair and
// returns the public address as the ID string
func NewID() (ID,error) {
	secp, err:= crypto.NewSECP256K1()
	if err != nil {
		return "",err
	}
	return ID(secp.PublicAddress()),nil
}

// Get distance of two ID
func (id ID) Xor(other ID) ID {
	idByte ,otherByte:= []byte(id), []byte(other)
	var xor [IDLength]byte
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i] ^ otherByte[i]
	}
	return ID(xor[:])
}

// Similar postfix bits length with another ID
func (id ID) SimilarPostfixLen(other ID) int {
	diff := []byte(id.Xor(other))
	ret := 0
	for i:= len(diff)-1;i>=0;i--{
		if diff[i] == uint8(0){
			ret+=8
		}else{
			bit:= fmt.Sprintf("%08b", diff[i])
			for j:=len(bit)-1;j>=0;j--{
				if bit[j]=='1'{
					return ret
				}
				ret++
			}
		}
	}
	return ret
}

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their ID.
type RoutingTable struct {
	ID ID
	Buckets [IDLengthInBits]list.List
}

// Createing new routing table
func NewRoutingTable(id ID) *RoutingTable {
	return &RoutingTable{ID: id, Buckets: [IDLengthInBits]list.List{}}
}


// Updating the new id in the routing table
func (rt *RoutingTable) Update(id ID) error {
	index:= rt.ID.SimilarPostfixLen(id)
	if index == IDLength {
		return errors.New("Can not updating node itself")
	}

	list := rt.Buckets[index]
	for e := list.Front(); e != nil; e = e.Next() {
		if id == e.Value{
			rt.Buckets[index].MoveToFront(e)
			break
		}
		if e.Next() == nil {
			rt.Buckets[index].InsertBefore(id,rt.Buckets[index].Front())
		}
	}
	return nil
}