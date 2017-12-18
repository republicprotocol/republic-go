package dht

import (
	"github.com/republicprotocol/go-identity"
	"container/list"
	"github.com/multiformats/go-multiaddr"
	"fmt"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = identity.IDLength
	AddressLength  = identity.AddressLength
	IDLengthInBits = IDLength * 8
	RepublicCode  = identity.RepublicCode
	Alpha          = 3
)

// Error when trying reach a index not in range
var ErrIndexOutOfRange = fmt.Errorf("index out of range ")

// Error when trying reach a index not in range
var ErrWrongType = fmt.Errorf("wrong variable type")

// RoutingBucket is a container List of strings.
type RoutingBucket struct {
	list.List
}

// MultiAddresses returns a string slice of multiaddresses.
func (bucket RoutingBucket) MultiAddresses() []string {
	multis := make([]string, bucket.Len())
	i := 0
	for it := bucket.Front(); it != nil; it = it.Next() {
		multis[i] = it.Value.(string)
		i++
	}
	return multis
}

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their republic address.
type RoutingTable struct {
	Address identity.Address
	Buckets [IDLengthInBits]RoutingBucket
}

// Create new routing table
func NewRoutingTable(address identity.Address) *RoutingTable {
	return &RoutingTable{Address: address, Buckets: [IDLengthInBits]RoutingBucket{}}
}

// MultiAddresses returns all multiaddresses in the table.
func (rt *RoutingTable) MultiAddresses() []string {
	// Find the total length.
	length := 0
	for _, bucket := range rt.Buckets {
		length += bucket.Len()
	}
	// Create a slice of strings and fill it with multiaddresses from all
	// buckets.
	multis := make([]string, length)
	i := 0
	for _, bucket := range rt.Buckets {
		for it := bucket.Front(); it != nil; it = it.Next() {
			multis[i] = it.Value.(string)
			i++
		}
	}
	return multis
}

// Check if we have enough space to update the node in the bucket
// Return the last node if the bucket is full, empty string otherwise
func (rt *RoutingTable) CheckAvailability(address identity.Address) (string, error) {

	same,err := rt.Address.SamePrefixLen(address)
	if err!= nil {
		return "", err
	}

	// Get the index of the bucket we want to store the ID
	index := IDLengthInBits - 1 - same
	if index < 0 || index > IDLengthInBits -1 {
		return "", ErrIndexOutOfRange
	}

	// Iterate the bucket see if we already know the node
	for e:= rt.Buckets[index].Front();e!= nil ;e = e.Next(){

		// Get the republic address from its multiaddress
		rAddress, err := e.Value.(multiaddr.Multiaddr).ValueForProtocol(RepublicCode)
		if err!= nil  {
			return "", err
		}
		if rAddress == string(address) {
			return "", nil
		}
	}

	// Check the bucket length
	if rt.Buckets[index].Len() < IDLength {
		return "", nil
	} else {
		// Return the republic address of the last node in the bucket
		rAddress, err := rt.Buckets[index].Back().Value.(multiaddr.Multiaddr).ValueForProtocol(RepublicCode)
		if err != nil {
			return "", err
		}
		return rAddress, nil
	}
}

// Update the new node in the routing table
func (rt *RoutingTable) Update(mAddress multiaddr.Multiaddr) error {

	// Get the node address from its multiaddress
	address, err := mAddress.ValueForProtocol(RepublicCode)
	if err != nil {
		return err
	}

	// Get the index of the bucket we want to store the ID
	same, err := rt.Address.SamePrefixLen(identity.Address(address))
	if err != nil {
		return err
	}
	index := IDLengthInBits - 1 - same
	if index < 0 || index > IDLengthInBits-1 {
		return ErrIndexOutOfRange
	}

	// If the node already exists, move it to the front
	for e := rt.Buckets[index].Front(); e != nil; e = e.Next() {
		// Check the type of value in the bucket
		value, ok := e.Value.(multiaddr.Multiaddr)
		if !ok {
			return ErrWrongType
		}

		// Get the republic address from its multiaddress
		rAddress, err := value.ValueForProtocol(RepublicCode)
		if err!= nil  {
			return err
		}

		// Override the multiaddress and move it to the front
		if address == rAddress {
			rt.Buckets[index].MoveToFront(e)
			return nil
		}
	}

	// Otherwise insert into the front
	rt.Buckets[index].PushFront(mAddress)

	return nil
}

// Return the addresses which is closer to the target
func (rt *RoutingTable) FindClosest(address identity.Address) (RoutingBucket, error) {
	// Find the bucket holding the target address .
	same,err := rt.Address.SamePrefixLen(address)
	if err != nil {
		return RoutingBucket{}, err
	}
	index := IDLengthInBits - 1 - same
	if index < 0 || index > IDLengthInBits-1 {
		return RoutingBucket{}, ErrIndexOutOfRange
	}

	// Initialize the returning bucket
	res := RoutingBucket{list.List{}}
	res.PushBackList(&rt.Buckets[index].List)

	// Keep adding nodes adjacent to the target bucket until we get enough node
	for i := 1; i < IDLengthInBits; i++ {
		if res.Len() >= Alpha {
			// todo : return  the sorted list and only the nodes closer than the rt itself
			return res, nil
		}

		if index-i >= 0 {
			res.PushBackList(&rt.Buckets[index-i].List)
		}

		if index+i < IDLengthInBits {
			res.PushBackList(&rt.Buckets[index+i].List)
		}
	}

	return res, nil
}

// Sort the node list and return the closets 3 nodes to the target
//func SortNode(lt RoutingBucket, target identity.Address) RoutingBucket {
//
//	if lt.Len() == 0 {
//		return lt
//	}
//	ret := RoutingBucket{list.List{}}
//
//	// Selection sort the list
//	for i := 0; i < Alpha; i++ {
//		if lt.Len() == 0 {
//			return ret
//		}
//		min := lt.Front()
//
//		for e := lt.Front(); e != nil; e = e.Next() {
//			minValue, _ := min.Value.(ma.Multiaddr).ValueForProtocol(RepublicCode)
//			value, _ := e.Value.(ma.Multiaddr).ValueForProtocol(RepublicCode)
//			if !less(minValue, value) {
//				min = e
//			}
//		}
//
//		ret.PushBack(lt.Remove(min))
//	}
//
//	return ret
//}

// Compare two lists if they are same in the first n elements
func CompareList(b1, b2 RoutingBucket, n int) bool {
	e1, e2 := b1.Front(), b2.Front()
	for i := 0; i < n; i++ {
		if e1 != nil && e2 != nil && e1 != e2 {
			return false
		} else if e1 != nil && e2 == nil {
			return false
		} else if e1 == nil && e2 != nil {
			return false
		} else if e1 == nil && e2 == nil {
			return true
		}
		e1, e2 = e1.Next(), e2.Next()
	}
	return true
}

// Kick the node from the routing table
func (rt *RoutingTable) Kick(id string) {
	for i := 0; i < IDLengthInBits; i++ {
		for e := rt.Buckets[i].Front(); e != nil; e = e.Next() {
			if e.Value == id {
				rt.Buckets[i].Remove(e)
				return
			}
		}
	}
}



