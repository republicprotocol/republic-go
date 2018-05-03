# New RPC endpoint and Configurable tags for logger. 

## RPC endpoint for getting status

We want to have an endpoint in our grpc service for getting status of the node.

> E.g. Is the node registered? Is the node bootstrapped? Number of nodes connected? ...
  
The status will has its own service and the protubuf file will be looked like this:

```$protobuf

syntax = "proto3";

package status;

service status {
  rpc Status (StatusRequest) returns (StatusResponse);
}

message StatusRequest {
}

message StatusResponse {
    string address = 1;
    bool bootstrapped = 2;
    int64 peers = 3;
}
```

## Command-line tool for rpc calls 

Since `status` endpoint is created for the need of falconry devOp tools, we also 
need to provide a command-line tool to call this rpc. For now, it should only 
support `status`, and will support all rpc calls in the future if needed.

```bash
$ rpc status 192.168.0.1:8080
8MGfbzAMS59Gb4cSjpm34soGNYsM2ft true 32
```

## Add command for getting pool hash in the registrar Command-line tool

We also want to know which pool the node is in. Since the data can be get from the 
darknode registry smart contract, we'll add a sub-command in the `registrar` cmd.

Get index of the pool
```bash
$ registrar pool 8MGfbzAMS59Gb4cSjpm34soGNYsM2ft
1 

$ registrar pool dsahfjdshjfdsjhfjsahd
-1
```

Check registration of node
```bash
$ registrar checkreg 0x009FbB7Aafee69081EF24c14A55A19c3cDd25eF4
true
```

## Configurable tags for the logger

In order to filter the logs for falconry tests easily, we want to have configurable tags 
in the logger. Tags are a list of key-value pairs which be stored as a map inside the logger.
Logger and logger option will have an extra filed which contains all the tags.
All logs will have tags attached to it for quick filtering and all plugins show support tags.
 




 

  
