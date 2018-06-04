package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/republicprotocol/go-do"
)

// A FilePlugin implements the Plugin interface by logging all events to a
// File. The Stdout File can be used to create a plugin that logs to Stdout.
type FilePlugin struct {
	do.GuardedObject

	file     *os.File
	filePath string
}

// FilePluginOptions are used to Unmarshal a FilePlugin from JSON. If the Path
// is set to stdout, or stderr, the respective output stream will be used
// instead of opening a File.
type FilePluginOptions struct {
	Path string `json:"path"`
}

// NewFilePlugin uses the FilePluginOptions to create a new FilePlugin.
func NewFilePlugin(filePluginOptions FilePluginOptions) (Plugin, error) {
	plugin := new(FilePlugin)
	plugin.GuardedObject = do.NewGuardedObject()
	plugin.filePath = filePluginOptions.Path
	return plugin, nil
}

// Start implements the Plugin interface. It opens the log file which will
// be opened as appendable and will be closed when the plugin is stopped.
func (plugin *FilePlugin) Start() error {
	var err error
	plugin.Enter(nil)
	defer plugin.Exit()
	// Initialise the file based on path
	switch plugin.filePath {
	case "stdout":
		plugin.file = os.Stdout
	case "stderr":
		plugin.file = os.Stderr
	default:
		plugin.file, err = os.OpenFile(plugin.filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0640)
	}
	return err
}

// Stop implements the Plugin interface. If the filePath is stdout or stderr
// it does nothing, otherwise it closes the open log file.
func (plugin *FilePlugin) Stop() error {
	plugin.Enter(nil)
	defer plugin.Exit()
	if plugin.file == os.Stdout || plugin.file == os.Stderr {
		return nil
	}
	return plugin.file.Close()
}

// Log implements the Plugin interface.
func (plugin *FilePlugin) Log(l Log) error {
	plugin.Enter(nil)
	defer plugin.Exit()
	if plugin.file == nil {
		return fmt.Errorf("cannot write log to file plugin: nil file")
	}
	if plugin.file == os.Stdout || plugin.file == os.Stderr {
		// format the tags to a string
		tags := []string{}
		for key, value := range l.Tags {
			tags = append(tags, fmt.Sprintf("%s:%s,", key, value))
		}
		tag := ""
		if len(tags) > 0 {
			tag = "{" + strings.Join(tags, ",") + "} "
		}

		_, err := plugin.file.WriteString(fmt.Sprintf("%s [%s] (%s) %s%s\n", l.Timestamp.Format("2006/01/02 15:04:05"), l.Level, l.EventType, tag, l.Event.String()))
		return err
	}
	return json.NewEncoder(plugin.file).Encode(l)
}
