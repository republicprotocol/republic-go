# Registrar Command-line tool

The **registrar** is the command-line tool for interacting with the darknode registry smart contract on ropsten.
It supports registering/deregistering nodes, checking registration , calling epoch and getting pool index.

# Install

Make you have included `$GOROOT` in the `$PATH` variable.
```bash
$ cd cmd/registrar
$ go install
```

Alternatively, after we merge this into master, you can run 
```bash
$ go get github.com/republicprotocol/republic-go/cmd/registrar
```

# Usage

```bash
$ registrar [global options] command [command options] [arguments...]
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

> You can use first letter of each command for quick access. e.g. `registrar e` has same effect with `registrar epoch`  

**GLOBAL OPTIONS**

- `--ren {value}`    
    republic token contract address (default: "0x65d54eda5f032f2275caa557e50c029cfbccbb54")
- `--dnr {value}`    
    dark node registry address (default: "0x69eb8d26157b9e12f959ea9f189A5D75991b59e3")
- `--help, -h`     
    show help
- `--version, -v`  
    print the version
