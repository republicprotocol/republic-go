package dht

import (
	"container/list"
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"github.com/republicprotocol/republic/crypto"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = crypto.PublicAddressLength
	IDLengthBase64 = 28
	IDLengthInBits = IDLength * 8
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
// ID is a string in a Base64 encoding
type ID string

// NewID creates a new set of public/private SECP256K1 key pair and
// returns the public address as the ID string
func NewID() (ID, error) {
	secp, err := crypto.NewSECP256K1()
	if err != nil {
		return "", err
	}
	return ID(secp.PublicAddress()), nil
}

// Get distance of two ID
func (id ID) Xor(other ID) ([]byte, error) {
	// Decode both the IDs into bytes
	idByte, err := base64.StdEncoding.DecodeString(string(id))
	if err != nil {
		return nil, err
	}
	otherByte, err := base64.StdEncoding.DecodeString(string(other))
	if err != nil {
		return nil, err
	}

	xor := make([]byte, 20)
	for i := 0; i < IDLength; i++ {
		xor[i] = idByte[i] ^ otherByte[i]
	}
	return xor, nil
}

// Similar postfix bits length with another ID
func (id ID) SamePrefixLen(other ID) (int, error) {
	diff, err := id.Xor(other)
	if err != nil {
		return 0, err
	}

	ret := 0
	for i := 0; i < IDLength; i++ {
		if diff[i] == uint8(0) {
			ret += 8
		} else {
			bit := fmt.Sprintf("%08b", diff[i])
			for j := 0; j < len(bit); j++ {
				if bit[j] == '1' {
					return ret, nil
				}
				ret++
			}
		}
	}

	return ret, nil
}

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their ID.
type RoutingTable struct {
	ID      ID
	Buckets [IDLengthInBits]list.List
}

// Createing new routing table
func NewRoutingTable(id ID) *RoutingTable {
	return &RoutingTable{ID: id, Buckets: [IDLengthInBits]list.List{}}
}

// Updating the new id in the routing table
func (rt *RoutingTable) Update(id ID) error {

	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return err
	}

	// The more same prefix-bit, the closer they are
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return errors.New("Can not updating node itself")
	}

	// If not exist, insert into the front of the bucket
	IdAddress := multiAddress(id)
	if rt.Buckets[index].Front() == nil {
		rt.Buckets[index].PushFront(IdAddress)
	}
	// Otherwise, move it to the front
	for e := rt.Buckets[index].Front(); e != nil; e = e.Next() {
		if IdAddress == e.Value {
			rt.Buckets[index].MoveToFront(e)
			return nil
		}
		if e.Next() == nil {
			rt.Buckets[index].InsertBefore(IdAddress, rt.Buckets[index].Front())
		}
	}

	return nil
}

// Return the addresses in the closest bucket
func (rt *RoutingTable) FindClosest(id ID) (*list.List, error) {
	// Find the bucket holding the target id
	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return nil, err
	}
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return nil, errors.New("Can not updating node itself")
	}
	res := list.New()
	res.PushBackList(&rt.Buckets[index])

	// Keep adding nodes adjacent to the target bucket until we get enough node
	for i := 1; i < IDLengthInBits; i++ {
		if res.Len() >= IDLength {
			return sortNode(res, id), nil
		}

		if index-i >= 0 {
			res.PushBackList(&rt.Buckets[index-i])
		}

		if index+i < IDLengthInBits {
			res.PushBackList(&rt.Buckets[index+i])
		}
	}

	return sortNode(res, id), nil
}

// Return all multiaddresses in the routing table
func (rt *RoutingTable) All() *list.List {
	all := list.New()
	for _, lt := range rt.Buckets {
		if lt.Front() != nil {
			all.PushBackList(&lt)
		}
	}
	return all
}

// todo: to be decided
func multiAddress(id ID) string {
	return "/republic/" + string(id)
}

// Sort the node list and return the closets 20 nodes to the target
func sortNode(lt *list.List, target ID) *list.List {
	if lt.Len() == 0 {
		return lt
	}
	ret := list.New()

	// Define less function between IDs
	less := func(add1, add2 string ) bool {
		// todo : need to update when we decied the format of multi-address
		id1, id2 := ID(add1[10:]) ,ID(add2[10:])
		xor1, _ := id1.Xor(target)
		xor2, _ := id2.Xor(target)

		for i := 0; i < IDLength; i++ {
			if xor1[i] < xor2[i] {
				return true
			} else if xor1[i] > xor2[i] {
				return false
			}
		}
		return false
	}

	// Select sort the list
	for i := 0; i < IDLength; i++ {
		if lt.Len() == 0 {
			return ret
		}
		min := lt.Front()
		for e := lt.Front(); e != nil; e = e.Next() {
			if !less(min.Value.(string), e.Value.(string)) {
				min = e
			}
		}
		ret.PushBack(lt.Remove(min))
	}

	return ret
}
