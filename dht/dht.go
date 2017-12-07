package dht

import (
	"container/list"
	"fmt"
)

// IDLength is the number of bytes needed to store an ID.
const (
	IDLength       = 20
	IDLengthInBits = IDLength * 8
)

// ID is the public address used to identify Nodes, and other entities, in the
// overlay network. It is generated from the public key of a key pair.
type ID string

// RoutingTable is a k-bucket routing table, where each bucket is a list of
// multiaddress strings, identifying peers by their network address as well as
// their ID.
type RoutingTable struct {
	Buckets [IDLengthInBits]list.List
}

// Createing new routing table
func NewRoutingTable() *RoutingTable {
	var buckets [IDLengthInBits]list.List
	return &RoutingTable{Buckets:buckets}
}

// Get distance of two ID
func (node ID) Xor(other ID) ID {
	nodeByte ,otherByte:= []byte(node), []byte(other)
	var xor [IDLength]byte
	for i := 0; i < IDLength; i++ {
		xor[i] = nodeByte[i] ^ otherByte[i]
	}
	return ID(xor[:])
}

// Similar postfix bits length with another ID
func (node ID) SimilarPostfixLen(other ID) int {
	diff := []byte(node.Xor(other))
	fmt.Println(diff)
	ret := 0
	for i:= len(diff)-1;i>=0;i--{
		if diff[i] == uint8(0){
			ret+=8
		}else{
			bit:= fmt.Sprintf("%08b", diff[i])
			fmt.Println(bit)
			for j:=len(bit)-1;j>=0;j--{
				if bit[j]=='1'{
					fmt.Println("ret=",ret)
					return ret
				}
				ret++
			}
		}
	}
	return ret
}