package dht

import (
	"container/list"
	"errors"
)

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
// their ID.
type RoutingTable struct {
	ID      ID
	Buckets [IDLengthInBits]RoutingBucket
}

// Create new routing table
func NewRoutingTable(id ID) *RoutingTable {
	return &RoutingTable{ID: id, Buckets: [IDLengthInBits]RoutingBucket{}}
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

// Update the new id in the routing table
func (rt *RoutingTable) Update(id ID) error {

	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return err
	}

	// The more same prefix-bit, the closer they are
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return errors.New("Can not update node itself")
	}
	IdAddress := MultiAddress(id)

	// If the node already exists, move it to the front
	for e := rt.Buckets[index].Front(); e != nil; e = e.Next() {
		if IdAddress == e.Value {
			rt.Buckets[index].MoveToFront(e)
			return nil
		}
	}

	// If we have reach the bucket limit, Ping the last node in the bucket
	if rt.Buckets[index].Len() == IDLength {
		//todo : need to redesign the structure?
	} else {
		// Otherwise simply insert the node into the front
		rt.Buckets[index].PushFront(IdAddress)
	}

	return nil
}

// Return the addresses in the closest bucket
func (rt *RoutingTable) FindClosest(id ID) (RoutingBucket, error) {
	// Find the bucket holding the target id.
	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return RoutingBucket{}, err
	}
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return RoutingBucket{}, errors.New("Can not update node itself")
	}

	res := RoutingBucket{list.List{}}
	res.PushBackList(&rt.Buckets[index].List)

	// Keep adding nodes adjacent to the target bucket until we get enough node
	for i := 1; i < IDLengthInBits; i++ {
		if res.Len() >= Alpha {
			return SortNode(res, id), nil
		}

		if index-i >= 0 {
			res.PushBackList(&rt.Buckets[index-i].List)
		}

		if index+i < IDLengthInBits {
			res.PushBackList(&rt.Buckets[index+i].List)
		}
	}

	return SortNode(res, id), nil
}

// Sort the node list and return the closets 20 nodes to the target
func SortNode(lt RoutingBucket, target ID) RoutingBucket {
	if lt.Len() == 0 {
		return lt
	}
	ret := RoutingBucket{list.List{}}

	// Define less function between IDs
	less := func(add1, add2 string) bool {
		// todo : need to update when we decided the format of multi-address
		// Currenly we assume the address will be sth like : /republic/dshfkadhfhkajdsf=
		id1, id2 := ID(add1[10:]), ID(add2[10:])
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

	// Selection sort the list
	for i := 0; i < Alpha; i++ {
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

// Compare two lists if they are same in the first n elements
func CompareList(l1, l2 *list.List, n int) bool {
	e1, e2 := l1.Front(), l2.Front()
	for i := 0; i < n; i++ {
		if e1 != nil && e2 != nil && e1 != e2 {
			return false
		}

		if e1 != nil && e2 == nil {
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

// Check if we have enough space to update the id in the bucket
// Return the last node if the bucket is full
func (rt *RoutingTable) CheckAvailability(id ID) (string, error) {

	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return "", err
	}

	// The more same prefix-bit, the closer they are
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return "", err
	}
	if rt.Buckets[index].Len() < IDLength {
		return "", nil
	} else {
		return rt.Buckets[index].Back().Value.(string), nil
	}
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
