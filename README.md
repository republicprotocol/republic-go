# Identity

An official reference implementation for Republic Identities written in Golang. It supports the generation of a Republic Identity, a Republic Address, and a Republic MultiAddress.

## Republic ID

A Republic ID an ECDSA key pair.

## Republic Address

A Republic Address is the Base58 encoding of the ECDSA public key of the Republic ID of the first 20 bytes.

## Republic MultiAddress

A Republic MultiAddress is a MultiAddress holding an IPv4/6 address and a Republic Address. Future implementations will use a I2P (or Kovri) address instead of an IPv4/6 address.