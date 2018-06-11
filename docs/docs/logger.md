# Logger

A Logger is an logging object that generates lines of output to an logger plugin.
Each plugin will have its own implementation of outputing the logs. Please refer 
[Documentation](https://godoc.org/github.com/republicprotocol/republic-go/logger) 
for detail usage.

#### Plugin

Currently we support two kinds of plugin, file plugin which will output the logs 
to a system file (or stdout) and websocket plugin which pumps logs through
websocket. 

#### Option

Logger will generated with a set of options, usually loaded from a config file.
Here's an example of the `StdoutLogger` with a pair of tags.

```json
"logger" : {
    "plugins": [
      {
        "file": {
          "path": "stdout"
        }
      }
    ] , 
    "tags" : {
      "falconry": "true"
    },
    "filterLevel": 2,
    "filterEvents": [
      "network",
      "usage",
    ]
  },
```

#### Tags

You can define a set of tags as key-value pairs in the logger so that each message 
generate by the logger will have all the tags attached. In this case, you can
easily filter the logs by using the tags.


#### Log Levels

We define 6 different logging levels, in order from most critical to least:

1. `LevelError`       something critically wrong occurred
2. `LevelWarn`        something bad occurred but it was not critical
3. `LevelInfo`        some helpful information 
4. `LevelDebugHigh`   important debugging information
5. `LevelDebug`       general debugging information
6. `LevelDebugLow`    unimportant debugging information

The logs levels can be filtered using `SetFilterLevel(Level)`, which defines the
highest possible level message that will be logged. The `FilterLevel` by default is
set to 2 which will only show `LevelError` and `LevelWarn` messages.
`SetFilterLevel(0)` will disable all log messages (not recommended), and
`SetFilterLevel(6)` will show all messages.


#### Event Types

Each log message can also be classified into Event Types. The Event Types we define
are:

- `TypeGeneric` 
- `TypeEpoch`   
- `TypeUsage` 
- `TypeOrderConfirmed`
- `TypeOrderMatch`
- `TypeOrderReceived`
- `TypeNetwork`
- `TypeCompute`

The `SetFilterEvents([]EventType)` acts as a whitelist for specific types of events.
For example, if you were only interested in viewing Network and Usage logs, you
could use `SetFilterEvents([]EventType{TypeNetwork, TypeUsage})` which would
hide all other types of events.
