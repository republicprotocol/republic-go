# Registrar
abigen --sol ./eth-republic/contracts/DarkNodeRegistrar.sol -pkg bindings --out dnr.go

# Atomic Swap
abigen --sol ./eth-atomic-swap/contracts/AtomicSwapEther.sol -pkg bindings --out AtomicSwapEth.go
abigen --sol ./eth-atomic-swap/contracts/AtomicSwapERC20.sol -pkg bindings --out AtomicSwapERC20.go
abigen --sol ./eth-atomic-swap/contracts/TestERC20.sol -pkg bindings --out TestERC20.go