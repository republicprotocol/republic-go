package logger

import (
	"log"
	"time"

	"github.com/republicprotocol/go-do"
)

// Constant strings for tagging logs.
const (
	TagNetwork   = "net"
	TagCompute   = "cmp"
	TagRegister  = "reg"
	TagUsage     = "usg"
	TagGeneral   = "gen"
	TagEthereum  = "eth"
	TagConsensus = "con"
)

type Logger struct {
	do.GuardedObject
	Plugins []Plugin
}

// Options are used to Unmarshal a Logger from JSON.
type Options struct {
	Plugins []PluginOptions `json:"plugins"`
}

type Plugin interface {
	Start() error
	Stop() error
	Info(tag, message string) error
	Warn(tag, message string) error
	Error(tag, message string) error
	Usage(cpu float32, memory, network int32) error
}

// PluginOptions are used to Unmarshal plugins from JSON.
type PluginOptions struct {
	File      *FilePluginOptions      `json:"file"`
	WebSocket *WebSocketPluginOptions `json:"websocket"`
}

// NewLogger returns a new Logger that will start and stop a set of plugins.
func NewLogger(options Options) (*Logger, error) {
	plugins := make([]Plugin, 0, len(options.Plugins))
	for i := range options.Plugins {
		if options.Plugins[i].File != nil {
			plugin, err := NewFilePlugin(*options.Plugins[i].File)
			if err != nil {
				return nil, err
			}
			plugins = append(plugins, plugin)
		}
		if options.Plugins[i].WebSocket != nil {
			plugin := NewWebSocketPlugin(*options.Plugins[i].WebSocket)
			plugins = append(plugins, plugin)
		}
	}
	return &Logger{
		GuardedObject: do.NewGuardedObject(),
		Plugins:       plugins,
	}, nil
}

// Start starts all the plugins of the logger
func (logger *Logger) Start() {
	for _, plugin := range logger.Plugins {
		if err := plugin.Start(); err != nil {
			log.Println(err)
		}
	}
}

// Stop stops all the plugins of the logger
func (logger Logger) Stop() {
	for _, plugin := range logger.Plugins {
		if err := plugin.Stop(); err != nil {
			log.Println(err)
		}
	}
}

func (logger *Logger) Error(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Error(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger *Logger) Info(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Info(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger *Logger) Warn(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Warn(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger *Logger) Usage(cpu float32, memory, network int32) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Usage(cpu, memory, network); err != nil {
			log.Println(err)
		}
	}
}

type Usage struct {
	Type string    `json:"type"`
	Time time.Time `json:"timestamp"`
	Data UsageData `json:"data"`
}

type UsageData struct {
	CPU     float32 `json:"cpu"`
	Memory  int32   `json:"memory"`
	Network int32   `json:"network"`
}

type Event struct {
	Type string    `json:"type"`
	Time time.Time `json:"timestamp"`
	Data EventData `json:"data"`
}

type EventData struct {
	Tag     string `json:"tag"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

type Error struct {
	Tag     string
	Message string
}

type Registration struct {
	NodeID     string `json:"nodeID"`
	PublicKey  string `json:"publicKey""`
	Address    string `json:"address"`
	RepublicID string `json:"republicID"`
}
