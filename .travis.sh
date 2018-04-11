
##############################
### Runs all tests locally ###
##############################

sleep_time=3

set -e
go vet ./...
golint -set_exit_status `go list ./... | grep -Ev "(stackint/asm|vendor)"`
go run cmd/testnetwork/main.go -sleep $sleep_time &
ganache_pid=$!

sleep $sleep_time

# Ensure that ganache is still running:
if ! kill -0 $ganache_pid >/dev/null 2>/dev/null; then
    RED='\033[1;31m'
    NOC='\033[0m'
    echo "${RED}Couldn't start ganache-cli!${NOC}"
    exit 1
fi

# go run cmd/epoch/main.go -testrpc &
ginkgo -v --cover --coverprofile coverprofile.out ./...
ginkgo -v --trace --race --cover --coverprofile coverprofile.out stackint
ginkgo -v --trace --race --cover --coverprofile coverprofile.out identity
ginkgo -v --trace --race --cover --coverprofile coverprofile.out shamir
ginkgo -v --trace --race --cover --coverprofile coverprofile.out order
ginkgo -v --trace --race --cover --coverprofile coverprofile.out compute
ginkgo -v --trace --race --cover --coverprofile coverprofile.out logger
ginkgo -v --trace --race --cover --coverprofile coverprofile.out contracts/connection
ginkgo -v --trace --race --cover --coverprofile coverprofile.out contracts/dnr
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network/dht
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network/rpc
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network
ginkgo -v --trace --race --cover --coverprofile coverprofile.out dark
ginkgo -v --trace --race --cover --coverprofile coverprofile.out hyperdrive
ginkgo -v --trace --cover --coverprofile coverprofile.out dark-node

fuser -k 8545/tcp

# Combine coverage reports
# Could use pattern matching
gocovmerge stackint/coverprofile.out identity/coverprofile.out shamir/coverprofile.out order/coverprofile.out compute/coverprofile.out logger/coverprofile.out contracts/connection/coverprofile.out contracts/dnr/coverprofile.out network/dht/coverprofile.out network/rpc/coverprofile.out network/coverprofile.out dark/coverprofile.out hyperdrive/coverprofile.out dark-node/coverprofile.out > coverprofile.out

sed -i '/rpc.pb.go/d' coverprofile.out
set -i '/bindings/d' coverprofile.out