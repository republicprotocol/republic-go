# Registry
abigen --sol ./republic-sol/contracts/DarkNodeRegistry.sol -pkg bindings --out dnr.go
abigen --sol ./republic-sol/contracts/Arc.sol -pkg bindings --out arc.go

# Hyperdrive registry
#abigen --sol ./republic-sol/contracts/HyperdriveEpoch.sol -pkg bindings --out hde.go

# Hyperdrive contract
abigen --sol ./republic-sol/contracts/Hyperdrive.sol -pkg bindings --out hyperdrive.go

