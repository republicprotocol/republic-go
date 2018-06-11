# Secure Multiparty Computation

## Overview

Writing instructions in and asynchronously receive result back. Correlate results to instructions by providing a unique ID (generating a unique ID that is the same one across the decentralized network of SMPCers is the responsibility of the user for most instructions, some high level functions will generate them for you, and is usually best done by making it dependent on the semantics of the operation so that all SMPCers will generate the same InstID without needing to communicate). Results for an InstID are *never* returned unless the local user has sent the revalent instruction.

Can maintain connections to different networks identified by ID. Separate state is kept for each network.

Each network has a dedicated Goroutine for processing the required computations. All streams have dedicated child Goroutines for sending/receiving messages with other nodes. Connection attempts are run using contexts that do not timeout but can be canceled (SMPCer is happy to wait indefinitely to successfully connect as it is more resilient to do this).

Router is used to forward instructions to the correct network Goroutine, and to forward messages across networks to the correct state (it is possible that one stream, in one network, will receive a message that was erroneously meant for another network and so it should be correctly handled).

## Interface

**Start**

No instructions processed until start.

**Shutdown**

No instructions processed after shutdown.

**Connect**

Issue an `InstConnect` to the SMPCer.

**Disconnect**

Issue an `InstDisconnect` to the SMPCer.

**Join**

Issue an `InstJ` to the SMPCer. Asynchronously join shares and eventually output a `ResultJ`.

## Example

Naive order matching.

1. How are InstIDs generated
2. Do sub outside of SMPC
3. Issue an InstJ
4. Wait for result

## Reconnecting

SMPC is not concerned about maintaining the lifetime of a connection. This is an issue that is addressed by the `stream` package and its implementations.

## Reputations

SMPC is responsible for tracking which nodes are only consuming data without producing it. Should, in future, implement a tit-for-tat model when interacting with other nodes. This ensure game theoretically optimal behaviour for punishing misbehaving SMPC nodes.