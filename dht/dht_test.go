package dht

import (
	"encoding/base64"
	"testing"
)

func TestNewID(t *testing.T) {

	// Generate new ID
	id, err := NewID()
	if err != nil {
		t.Fatal("error in generating new keys:", err)
	}

	// Check if the ID is in a base64 coding
	idBytes, err := base64.StdEncoding.DecodeString(string(id))
	if err != nil {
		t.Fatal("fail in decoding the id from base64:", err)
	}

	// Decode the ID and check the length of id in bytes
	if len(idBytes) != IDLength {
		t.Fatal("wrong id length", err)
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
	id1, id2 := ID("ikuMqMkI9GcY/viwBiXnh3s7Mn4="), ID("WPBI8Kni4fKfArBrFH4umRU4czA=")

	decode1, err := base64.StdEncoding.DecodeString(string(id1))
	if err != nil {
		t.Fatal("fail to decode the ID:", err)
	}
	decode2, err := base64.StdEncoding.DecodeString(string(id2))
	if err != nil {
		t.Fatal("fail to decode the ID:", err)
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

func TestID_SimilarPostfixLen(t *testing.T) {
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
	same, err := id1.SimilarPostfixLen(id2)

	if err != nil {
		t.Fatal("error in getting the result", err)
	}
	if same < 0 || same > IDLengthInBits {
		t.Fatal("impossilble output", err)
	}

	// Try swap the two ID and check if we still get the same difference
	commutativeSame, err := id2.SimilarPostfixLen(id1)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}
	if same != commutativeSame {
		t.Fatal("similar should be consistent in a commutative way", err)
	}
}

func TestID_SimilarPostfixLen1(t *testing.T) {
	// Create two new IDs
	id1, id2 := ID("ikuMqMkI9GcY/viwBiXnh3s7Mn4="), ID("WPBI8Kni4fKfArBrFH4umRU4czA=")

	// Get the xor difference
	same, err := id1.SimilarPostfixLen(id2)
	if err != nil {
		t.Fatal("error in running the function:", err)
	}

	if same != 1 {
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
	if len(rt.ID) != IDLengthBase64 {
		t.Fatal("routing table has wrong id length")
	}
	if len(rt.Buckets) != IDLengthInBits {
		t.Fatal("routing table has wrong bucket length")
	}
}

func TestRoutingTable_Update(t *testing.T) {
	// Generating new routing table
	id := ID("ikuMqMkI9GcY/viwBiXnh3s7Mn4=")
	id2 := ID("WPBI8Kni4fKfArBrFH4umRU4czA=")
	id3 := ID("gx4SyWyDK1I0rdBXGGBeWnkFWag=")

	rt := NewRoutingTable(id)
	if rt == nil {
		t.Fatal("new routing table is nil")
	}

	err := rt.Update(id2)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	for i, j := range rt.Buckets {
		// Check if id2 is updated in the right bucket
		if j.Front() != nil {
			if i != 1 {
				t.Fatal("peer address is stored in the wrong bucket")
			}
			if j.Front().Value != "/republic/"+id2 && j.Front().Next() == nil {
				t.Fatal("fail to store the right multiaddress of a id")
			}
			break
		}

		// if no peer is found, throw an error
		if i == len(rt.Buckets)-1 && j.Front() == nil {
			t.Fatal("fail to update the id in the routing table")
		}
	}

	// update id3 into the routing table which should be updated into
	// bucket[1], same bucket as id2
	err = rt.Update(id3)
	if err != nil {
		t.Fatal("error in updating another peer:", err)
	}

	for i, j := range rt.Buckets {
		// Check if id2 is updated in the right bucket
		if j.Front() != nil {
			if j.Front().Next() == nil {
				t.Fatal("id3 wasn't updated into the routing table ")
			}
			break
		}

		// if no peer is found, throw an error
		if i == len(rt.Buckets)-1 && j.Front() == nil {
			t.Fatal("fail to update the id in the routing table")
		}

	}

}

func TestRoutingTable_FindClosest(t *testing.T) {

}

func TestRoutingTable_All(t *testing.T) {

}
