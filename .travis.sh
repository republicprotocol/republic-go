# Simulate the environment used by Travis CI so that we can run local tests to
# find and resolve issues that are consistent with the Travis CI environment.
# This is helpful because Travis CI often finds issues that our own local tests
# do not.

# go vet ./...
# golint -set_exit_status `go list ./... | grep -Ev "(stackint/asm|vendor)"`
# golint `go list ./... | grep -Ev "(stackint/asm|vendor)"`

go build ./...
GOMAXPROCS=1 CI=true ginkgo -v --race --cover --coverprofile coverprofile.out ./...
covermerge crypto/coverprofile.out dht/coverprofile.out dispatch/coverprofile.out grpc/coverprofile.out http/coverprofile.out http/adapter/coverprofile.out identity/coverprofile.out leveldb/coverprofile.out logger/coverprofile.out ome/coverprofile.out order/coverprofile.out orderbook/coverprofile.out shamir/coverprofile.out smpc/coverprofile.out stackint/coverprofile.out stream/coverprofile.out swarm/coverprofile.out > coverprofile.out

sed -i '/.pb.go/d' coverprofile.out
sed -i '/bindings/d' coverprofile.out
sed -i '/cmd/d' coverprofile.out