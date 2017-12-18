.PHONY: all clean rpc test install

TAG = $(shell git log --pretty=format:'%h' -n 1)

all: clean rpc test install

clean:
	rm -f rpc/*.pb.go

rpc: rpc/rpc.proto
	protoc -I rpc/ rpc/*.proto --go_out=plugins=grpc:rpc

test: rpc
	go test -v
	go test -v ./dht

install: rpc
	go install