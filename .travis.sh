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

fuser -k 8545/tcp

# Combine coverage reports
# Could use pattern matching
gocovmerge stackint/coverprofile.out identity/coverprofile.out shamir/coverprofile.out order/coverprofile.out logger/coverprofile.out contracts/coverprofile.out darknode/coverprofile.out > coverprofile.out

sed -i '/rpc.pb.go/d' coverprofile.out
sed -i '/bindings/d' coverprofile.out