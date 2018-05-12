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
    }
  },
```

#### Tags

You can define a set of tags as key-value pairs in the logger so that each message 
generate by the logger will have all the tags attached. In this case, you can easily 
filter the logs by using the tags.


#### Log type and event type


Log type shows the type of the log message.

- `info` normal information or message
- `warn`  warning 
- `error` something wrong happens

Event type shows the event relating to the message which have below categories.

- `generic` 
- `epoch`   
- `usage` 
- `ethereum`
- `orderMatch`
- `orderReceived`
- `network`
- `compute`
 