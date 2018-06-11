# RPC Command-line tool

The rpc command-line tool is used for sending gRPC requests to Darknodes. For now, it only supports the `Status` service, which exposes the `status` RPC. It is primarily used by the `falconry` command-line tool.
                            
# Install

Ensure that you have `$GOBIN` configured and included in your `$PATH` environment variable.

```sh
$ go get -u github.com/republicprotocol/republic-go/cmd/rpc
```

> Note: This will not work until the current `develop` branch has been merged into `master`.

# Usage

```bash
$  rpc [global options] command [command options] [arguments...]
```


**COMMANDS**

- `status {ipAddress}` 
   It take one argument which is the network address in the format of `0.0.0.0:8080`.
   The values are separated by space and in the order of *address*, *isBootstrapped*, *connectNodes*  
- `help`        
   Shows a list of commands or help for one command

The `rpc` supports shorthand subcommands for quick access. You can use the first letter of the subcommand instead of the entire subcommand. For example, `rpc s` has same effect as `rpc status`.

**GLOBAL OPTIONS**

- `--help, -h`     
    show help
- `--version, -v`  
    print the version

