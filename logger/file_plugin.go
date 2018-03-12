package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/republicprotocol/go-do"
)

// A FilePlugin implements the Plugin interface by logging all events to a
// File. The Stdout File can be used to create a plugin that logs to Stdout.
type FilePlugin struct {
	do.GuardedObject

	file *os.File
}

// FilePluginOptions are used to Unmarshal a FilePlugin from JSON. If the Path
// is set to stdout, or stderr, the respective output stream will be used
// instead of opening a File.
type FilePluginOptions struct {
	Path string `json:"path"`
}

// NewFilePlugin uses the give File to create a new FilePlugin. The file will
// be opened as appendable and will be closed when the plugin is stopped.
func NewFilePlugin(filePluginOptions FilePluginOptions) (Plugin, error) {
	var err error
	plugin := new(FilePlugin)
	plugin.GuardedObject = do.NewGuardedObject()
	switch filePluginOptions.Path {
	case "stdout":
		plugin.file = os.Stdout
	case "stderr":
		plugin.file = os.Stderr
	default:
		plugin.file, err = os.OpenFile(filePluginOptions.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}
	return plugin, err
}

// Start implements the Plugin interface. It does nothing.
func (plugin *FilePlugin) Start() error {
	plugin.Enter(nil)
	defer plugin.Exit()

	return nil
}

// Stop implements the Plugin interface. It does nothing.
func (plugin *FilePlugin) Stop() error {
	plugin.Enter(nil)
	defer plugin.Exit()

	if plugin.file == os.Stdout {
		return nil
	}
	if plugin.file == os.Stderr {
		return nil
	}
	return plugin.file.Close()
}

// Info implements the Plugin interface.
func (plugin *FilePlugin) Info(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	if plugin.file == nil {
		return fmt.Errorf("cannot write logs to a nil file")
	}
	_, err := plugin.file.WriteString(fmt.Sprintf("%s [info] (%s) %s\n", time.Now().Format("2018/02/03 10:00:00"), tag, message))
	return err
}

// Warn implements the Plugin interface.
func (plugin *FilePlugin) Warn(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	if plugin.file == nil {
		return fmt.Errorf("cannot write logs to a nil file")
	}
	_, err := plugin.file.WriteString(fmt.Sprintf("%s [warn] (%s) %s\n", time.Now().Format("2018/02/03 10:00:00"), tag, message))
	return err
}

// Error implements the Plugin interface.
func (plugin *FilePlugin) Error(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	if plugin.file == nil {
		return fmt.Errorf("cannot write logs to a nil file")
	}
	_, err := plugin.file.WriteString(fmt.Sprintf("%s [error] (%s) %s\n", time.Now().Format("2018/02/03 10:00:00"), tag, message))
	return err
}

// Usage implements the Plugin interface.
func (plugin *FilePlugin) Usage(cpu float32, memory, network int32) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	if plugin.file == nil {
		return fmt.Errorf("cannot write logs to a nil file")
	}
	_, err := plugin.file.WriteString(fmt.Sprintf("%s [info] ("+TagUsage+") cpu = %.3f MHz; memory = %d MB; network = %d KB\n", time.Now().Format("2018/02/03 10:00:00"), cpu, memory, network))
	return err
}
