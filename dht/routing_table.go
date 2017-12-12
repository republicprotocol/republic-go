package dht

import (
	"container/list"
	"errors"
)

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their ID.
type RoutingTable struct {
	ID      ID
	Buckets [IDLengthInBits]list.List
}

// Create new routing table
func NewRoutingTable(id ID) *RoutingTable {
	return &RoutingTable{ID: id, Buckets: [IDLengthInBits]list.List{}}
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
		return errors.New("Can not updating node itself")
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

// Return all multi-addresses in the routing table
func (rt *RoutingTable) All() *list.List {
	all := list.New()
	for _, lt := range rt.Buckets {
		if lt.Front() != nil {
			all.PushBackList(&lt)
		}
	}
	return all
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
		if res.Len() >= Alpha {
			return SortNode(res, id), nil
		}

		if index-i >= 0 {
			res.PushBackList(&rt.Buckets[index-i])
		}

		if index+i < IDLengthInBits {
			res.PushBackList(&rt.Buckets[index+i])
		}
	}

	return SortNode(res, id), nil
}

// Sort the node list and return the closets 20 nodes to the target
func SortNode(lt *list.List, target ID) *list.List {
	if lt.Len() == 0 {
		return lt
	}
	ret := list.New()

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
func (rt *RoutingTable) CheckAvailability(id ID) (string,error) {

	same, err := rt.ID.SamePrefixLen(id)
	if err != nil {
		return "", err
	}

	// The more same prefix-bit, the closer they are
	index := IDLengthInBits - 1 - same
	if index < 0 {
		return "", err
	}
	if  rt.Buckets[index].Len() < IDLength {
		return "", nil
	}else{
		return rt.Buckets[index].Back().Value.(string), nil
	}
}

// Kick the node from the routing table
func (rt *RoutingTable) Kick(id string){
	for i:= 0;i< IDLengthInBits;i++{
		for e:= rt.Buckets[i].Front();e!= nil; e= e.Next(){
			if e.Value == id {
				rt.Buckets[i].Remove(e)
				return
			}
		}
	}
}
