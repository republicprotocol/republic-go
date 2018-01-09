# Swarm

The P2P networking layer for the Republic Protocol.

## Install Proto3

Install Protobuf.

```
curl -OL https://github.com/google/protobuf/releases/download/v3.2.0/protoc-3.2.0-linux-x86_64.zip
unzip protoc-3.2.0-linux-x86_64.zip -d protoc3

sudo mv protoc3/bin/* /usr/local/bin/
sudo mv protoc3/include/* /usr/local/include/

sudo chown $USER /usr/local/bin/protoc
sudo chown -R $USER /usr/local/include/google
```

Install the Go plugin.

```
go get -u github.com/golang/protobuf/protoc-gen-go
```

## Install gRPC

Install the gRPC plugin.

```
go get -u google.golang.org/grpc
```