# RPC Command-line tool

The rpc command-line tool is used for making grpc call to the nodes.
For now, it only supports `status` service which has a single function
for getting status of the node. It's mainly used by the falconry deveOp 
tool. We'll add support for other rpc services if needede. 
                            
# Install

Make you have included `$GOROOT` in the `$PATH` variable.
```bash
$ cd cmd/rpc
$ go install
```

Alternatively, after we merge this into master, you can run 
```bash
$ go get github.com/republicprotocol/republic-go/cmd/rpc
```

# Usage

```bash
$  rpc [global options] command [command options] [arguments...]
```


**COMMANDS**

- `status {ipAddress}` 
   It take one argument which is the network address in the format of `0.0.0.0:8080`.
   The values are seperated by space and in the order of *address*, *isBootstrapped*, *connectNodes*  
- `help`        
   Shows a list of commands or help for one command

> You can use first letter of each command for quick access. e.g. `registrar e` has same effect with `registrar epoch`  

**GLOBAL OPTIONS**

- `--help, -h`     
    show help
- `--version, -v`  
    print the version

