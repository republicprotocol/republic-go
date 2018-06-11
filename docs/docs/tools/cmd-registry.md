# Registry Command-line tool

The **registry** is the command-line tool for interacting with the Farknode Registry smart contract on the Ethereum Ropsten network. It supports the registration/deregistration of Darknodes, checking registration status of Darknodes, triggering the epoch, and getting the pool index of a Darknode.

# Install

Ensure that you have `$GOBIN` configured and included in your `$PATH` environment variable.

```sh
$ go get -u github.com/republicprotocol/republic-go/cmd/registry
```

> Note: This will not work until the current `develop` branch has been merged into `master`.

# Usage

```bash
$ registry [global options] command [command options] [arguments...]
```

**COMMANDS**

- `epoch`
   Calling epoch
- `checkreg {address}`
   It takes one argument which is the nods's ethereum address.
   Check if the node with given `address` is registered or not
- `register {address1, address2 ...}`
   It takes a list of ethereum address.
   Register nodes with given address in the dark node registry
- `deregister {address1, address2 ...}`
   It takes a list of ethereum address.
   Deregister nodes with given address in the dark node registry
- `pool {address}`
   It takes one argument which is the node's ethereum address.
   Get the index of the pool the node is in, return -1 if no pool found
- `help`        
   Shows a list of commands or help for one command

The `registry` supports shorthand subcommands for quick access. You can use the first letter of the subcommand instead of the entire subcommand. For example, `registry e` has same effect as `registry epoch`.

**GLOBAL OPTIONS**

- `--ren {value}`    
    republic token contract address (default: "0x65d54eda5f032f2275caa557e50c029cfbccbb54")
- `--dnr {value}`    
    dark node registry address (default: "0x69eb8d26157b9e12f959ea9f189A5D75991b59e3")
- `--help, -h`     
    show help
- `--version, -v`  
    print the version
