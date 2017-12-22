# Identity

An official reference implementation for Republic Identities written in Golang. 

It supports the generation of a Republic ID, a Republic Address, and a Republic MultiAddress.

## Republic ID

A Republic ID the first 20 bytes of the Keccak256 hash of the public key of an ECDSA key pair.

```go
	keyPair, _ := identity.NewKeyPair()
	republicID := keyPair.PublicID()
	fmt.Println(republicID) // [247 107 192 130 208 33 69 182 183 127 237 24 22 83 157 99 207 221 9 243]
```

## Republic Address

A Republic Address is the Base58 encoding of the MultiHash of the Republic ID.


```go
	keyPair, _ := identity.NewKeyPair()
	republicAddress := keyPair.PublicAddress()
	fmt.Println(republicAddress) // 8MGHyHWxQ5zUzXXnK5b6LkWGFCWSyX
```

## Republic MultiAddress

A Republic MultiAddress is a MultiAddress holding an IPv4/6 address and a Republic Address.

For example : `/ip4/127.0.0.1/udp/1234/republic/8MGfbzAMS59Gb4cSjpm34soGNYsM2f` 


```go
	ipMulti,_:= idendity.NewMultiaddr("/ip4/127.0.0.1/tcp/80")
	t,_ := identity.NewKeyPair()
	republicMulti, _  := t.MultiAddress()
	republicMulti = multiaddr.Join(republicMulti,ipMulti)

	fmt.Println(republicMulti.ValueForProtocol(identity.P_REPUBLIC)) // 8MGfbzAMS59Gb4cSjpm34soGNYsM2f <nil>
	fmt.Println(republicMulti.ValueForProtocol(identity.P_IP4)) // 127.0.0.1 <nil>
	fmt.Println(republicMulti.ValueForProtocol(identity.P_TCP)) // 80 <nil>
```
 
Future implementations will use a I2P (or Kovri) address instead of an IPv4/6 address.
