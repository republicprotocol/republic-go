.PHONY: all clean build

all: clean build

clean:
	rm -f ./*.pb.go

build: ./rpc.proto
	protoc -I ./ ./*.proto --go_out=plugins=grpc:.