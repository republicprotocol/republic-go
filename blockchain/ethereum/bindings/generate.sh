#!/usr/bin/env bash

# Setup
sed -i "" -e 's/"zeppelin-solidity\/contracts\//".\/zeppelin-solidity\/contracts\//' contracts/*.sol contracts/*/*.sol
mkdir ./contracts/zeppelin-solidity
cp -r ./node_modules/zeppelin-solidity/contracts ./contracts/zeppelin-solidity/contracts

### GENERATE BINDINGS HERE
abigen --sol ./contracts/RenLedger.sol -pkg bindings --out bindings.go

# Revert setup
sed -i "" -e 's/".\/zeppelin-solidity\/contracts\//"zeppelin-solidity\/contracts\//' contracts/*.sol contracts/*/*.sol
rm -r ./contracts/zeppelin-solidity


