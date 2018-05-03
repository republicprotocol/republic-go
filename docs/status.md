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
need to provide a command-line tool to make this rpc call. For now , it will only 
support `status`, but will support all rpc calls in the future if needed.

todo : details of the cmd

## Add command for getting pool hash in the registrar Command-line tool

We also want to know which pool the node is in. Since the data can be get from the 
darknode registry smart contract, we'll add a sub-command in the `registrar` cmd.

## Configurable tags for the logger

In order to filter the logs for falconry tests easily, we want to have configurable tags 
in the logger. 
 




 

  
