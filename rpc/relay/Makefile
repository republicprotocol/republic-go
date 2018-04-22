.PHONY: all clean build

all: clean build

clean:
	rm -f ./*.pb.go

build: ./relay.proto
	protoc -I ./ ./*.proto --go_out=plugins=grpc:.