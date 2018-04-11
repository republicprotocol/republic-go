
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
ginkgo -v --trace --race --cover --coverprofile coverprofile.out identity
ginkgo -v --trace --race --cover --coverprofile coverprofile.out shamir
ginkgo -v --trace --race --cover --coverprofile coverprofile.out order
ginkgo -v --trace --race --cover --coverprofile coverprofile.out compute
ginkgo -v --trace --race --cover --coverprofile coverprofile.out logger
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network/dht
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network/rpc
ginkgo -v --trace --race --cover --coverprofile coverprofile.out network
ginkgo -v --trace --race --cover --coverprofile coverprofile.out dark
ginkgo -v --trace --race --cover --coverprofile coverprofile.out hyperdrive
ginkgo -v --trace --cover --coverprofile coverprofile.out dark-node
# sed -i '/rpc.pb.go/d' network/rpc/coverprofile.out

fuser -k 8545/tcp
