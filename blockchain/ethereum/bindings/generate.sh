
# Setup
cd ./republic-sol
npm install
sed -i -e 's/"zeppelin-solidity\/contracts\//".\/zeppelin-solidity\/contracts\//' contracts/*.sol contracts/*/*.sol
mkdir ./contracts/zeppelin-solidity
cp -r ./node_modules/zeppelin-solidity/contracts ./contracts/zeppelin-solidity/contracts
cd ../





### GENERATE BINDINGS HERE

# Registry
abigen --sol ./republic-sol/contracts/DarknodeRegistry.sol -pkg bindings --out dnr.go
abigen --sol ./republic-sol/contracts/Arc.sol -pkg bindings --out arc.go




# Revert setup
cd ./republic-sol
sed -i -e 's/".\/zeppelin-solidity\/contracts\//"zeppelin-solidity\/contracts\//' contracts/*.sol contracts/*/*.sol
rm -r ./contracts/zeppelin-solidity
cd ../


