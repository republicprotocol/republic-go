package dht

import (
	"testing"
)

const IDLengthBase58 = 30

func TestNewRoutingTable(t *testing.T) {
	// Generating new routing table
	id, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}
	rt := NewRoutingTable(id)
	if rt == nil {
		t.Fatal("new routing table is nil")
	}

	// Check the fields of the routing table
	if len(rt.ID) != IDLengthBase58 {
		t.Fatal("routing table has wrong id length")
	}
	if len(rt.Buckets) != IDLengthInBits {
		t.Fatal("routing table has wrong bucket length")
	}
}

func TestRoutingTable_Update(t *testing.T) {
	// Generating test IDs
	id := ID("ikuMqMkI9GcY7viwBiXnh3s7Mn4=")
	id2 := ID("ikuI8Kni4fKfArBrFH4umRU4czA=")
	id3 := ID("ikufowkI9GcY7viwBiXnh3s7Mn4=")

	multi2, err := ma.NewMultiaddr("/republic/" + string(id2))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}

	multi3, err := ma.NewMultiaddr("/republic/" + string(id3))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}

	rt := NewRoutingTable(id)
	if rt == nil {
		t.Fatal("new routing table is nil")
	}

	// Update id2 into the routing table
	err = rt.Update(multi2)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id2 is correctly inserted into the bucket
	list := rt.Buckets[138]
	if list.Front() == nil || list.Front().Value != "/republic/"+string(id2) {
		t.Fatal("fail to store the right multiaddress of a id")
	}

	// Check if all other buckets are empty
	for i, j := range rt.Buckets {
		if i != 138 {
			if j.Front() != nil {
				t.Fatal("bucket["+string(i)+"] should be empty, but got:", list.Front().Value)
			}
		}
	}

	// update id3 into the routing table which should be updated into
	// bucket[1], same bucket as id2
	err = rt.Update(multi3)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id3 is correctly inserted into the bucket
	list = rt.Buckets[140]
	if list.Front() == nil || list.Front().Value != "/republic/"+string(id3) {
		t.Fatal("fail to store the right multiaddress of a id")
	}

	// Check if all other buckets are empty
	for i, j := range rt.Buckets {
		if i != 138 && i != 140 {
			if j.Front() != nil {
				t.Fatal("bucket["+string(i)+"] should be empty, but got:", list.Front().Value)
			}
		}
	}
}

func TestNewRoutingTable2(t *testing.T) {
	// Generating test IDs
	// id2 ,id3 will be inserted into the same bucket[138]
	id := ID("ikuMqMkI9GcY7viwBiXnh3s7Mn4=")
	id2 := ID("ikuI8Kni4fKfArBrFH4umRU4czA=")
	id3 := ID("ikuIowkI9GcY7viwBiXnh3s7Mn4=")

	multi2, err := ma.NewMultiaddr("/republic/" + string(id2))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}

	multi3, err := ma.NewMultiaddr("/republic/" + string(id3))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}

	rt := NewRoutingTable(id)
	if rt == nil {
		t.Fatal("new routing table is nil")
	}

	// Update id2 into the routing table
	err = rt.Update(multi2)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id2 is correctly inserted into the bucket
	list := rt.Buckets[138]
	if list.Front() == nil || list.Front().Value != "/republic/"+string(id2) {
		t.Fatal("fail to store the right multiaddress of a id")
	}

	// Check if all other buckets are empty
	for i, j := range rt.Buckets {
		if i != 138 {
			if j.Front() != nil {
				t.Fatal("bucket["+string(i)+"] should be empty, but got:", list.Front().Value)
			}
		}
	}

	// update id3 into the routing table which should be updated into
	// bucket[1], same bucket as id2
	err = rt.Update(multi3)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id3 is correctly inserted into the bucket
	list = rt.Buckets[138]
	if list.Front() == nil || list.Front().Value != "/republic/"+string(id3) || list.Front().Next().Value != "/republic/"+string(id2) {
		t.Fatal("fail to store the right multiaddress of a id")
	}

	// Update id2 again
	err = rt.Update(multi2)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id2 is been pushed into the front of the list
	list = rt.Buckets[138]
	if list.Front() == nil || list.Front().Value != "/republic/"+string(id2) || list.Front().Next().Value != "/republic/"+string(id3) {
		t.Fatal("fail to store the right multiaddress of a id")
	}
}

func TestRoutingTable_All(t *testing.T) {
	id := ID("ikuMqMkI9GcY7viwBiXnh3s7Mn4=")
	id2 := ID("ikuI8Kni4fKfArBrFH4umRU4czA=")
	id3 := ID("ikuIowkI9GcY7viwBiXnh3s7Mn4=")

	multi2, err := ma.NewMultiaddr("/republic/" + string(id2))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}

	multi3, err := ma.NewMultiaddr("/republic/" + string(id3))
	if err != nil {
		t.Fatal("error in creating new multi address:", err)
	}
	rt := NewRoutingTable(id)
	if rt == nil {
		t.Fatal("new routing table is nil")
	}

	// Update id2 into the routing table
	err = rt.Update(multi2)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id2 is correctly inserted into the bucket
	list := rt.MultiAddresses()
	if list[0] != "/republic/"+string(id2) {
		t.Fatal("fail to store the right multiaddress of a id")
	}

	// update id3 into the routing table which should be updated into
	// bucket[1], same bucket as id2
	err = rt.Update(multi3)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	// Check if id3 is correctly inserted into the bucket
	list = rt.MultiAddresses()
	if list[0] != "/republic/"+string(id3) || list[1] != "/republic/"+string(id2) {
		t.Fatal("fail to store the right multiaddress of a id")
	}
}

// Test function of finding closer nodes of a target
func TestRoutingTable_FindClosest(t *testing.T) {

}
