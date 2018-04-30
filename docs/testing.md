# Testing

How to install Ginkgo and run Ginkgo tests

## Darknode Local Testnet

How to use the Darknode local testnet in Ginkgo

What is the Darknode local testnet?
What does it setup?
    Darknodes
    Bootstraps them
    Turns on their RPCs
    Connects to ganache
    Deploys smart contracts
What does it *not* setup?
    Bitcoin
    Cannot control smart contract variables
How do you use it?
    It prevents multiple tests running simultaneously because it needs to control local ports
    It uses ganache snapshots to keep a clean local ethereum testnet between tests
    See the darknode/testnet_test.go for an example