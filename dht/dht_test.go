package dht

import (
	"github.com/jbenet/go-base58"
	"testing"
)

const IDLengthBase58 = 30

func TestNewID(t *testing.T) {

	// Generate new ID
	id, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}

	// Check if the ID is in a base64 coding
	idBytes := base58.Decode(string(id))

	if len(idBytes) == 0 {
		t.Fatal("fail to decode the id from base58:", err)
	}

	// Check the length of decoded id in bytes
	if len(idBytes) != IDLength+2 {
		t.Fatal("wrong id length")
	}
}

func TestID_Xor(t *testing.T) {
	// Create two new IDs
	id1, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}
	id2, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}

	// Get the xor difference
	diff, err := id1.Xor(id2)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}

	// Check if the difference has valid length
	if len(diff) != IDLength {
		t.Fatal("wrong id length", err)
	}

	// Try swap the two ID and check if we still get the same difference
	commutativeDiff, err := id2.Xor(id1)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}
	if len(diff) != len(commutativeDiff) {
		t.Fatal("xor should be consistent in a commutative way", err)
	}
	for i := 0; i < len(diff); i++ {
		if diff[i] != commutativeDiff[i] {
			t.Fatal("xor should be consistent in a commutative way", err)
		}
	}
}

// Test correctness of the function
func TestID_Xor1(t *testing.T) {
	// Create two new IDs
	id1, id2 := ID("8MHUraYAumAn5dEqEnHqu8LGnR4hms"), ID("8MG913SweBDhihrsMvFYgTBRpca9UF")

	decode1 := base58.Decode(string(id1))
	if len(decode1) == 0 {
		t.Fatal("fail to decode the ID:")
	}
	decode2 := base58.Decode(string(id2))
	if len(decode2) == 0 {
		t.Fatal("fail to decode the ID:")
	}

	// Get the xor difference
	diff, err := id1.Xor(id2)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}

	for i := 0; i < IDLength; i++ {
		if diff[i] != (decode1[i] ^ decode2[i]) {
			t.Fatal("wrong result.")
		}
	}
}

func TestID_SamePrefixLen(t *testing.T) {
	// Generating two new IDs
	id1, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}

	id2, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}

	// Get the similar postfix bits length
	same, err := id1.SamePrefixLen(id2)

	if err != nil {
		t.Fatal("error in getting the result", err)
	}
	if same < 0 || same > IDLengthInBits {
		t.Fatal("impossilble output", err)
	}

	// Try swap the two ID and check if we still get the same difference
	commutativeSame, err := id2.SamePrefixLen(id1)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}
	if same != commutativeSame {
		t.Fatal("similar should be consistent in a commutative way", err)
	}
}

func TestID_SamePrefixLen2(t *testing.T) {
	// Create two new IDs
	id1, id2 := ID("8MHUraYAumAn5dEqEnHqu8LGnR4hms"), ID("8MG913SweBDhihrsMvFYgTBRpca9UF")

	// Get the same prefix of the two IDs
	same, err := id1.SamePrefixLen(id2)

	if err != nil {
		t.Fatal("error in running the function:", err)
	}

	if same != 17 {
		t.Fatal("wrong result")
	}
}

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
