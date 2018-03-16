package logger

import (
	"fmt"
	"log"
	"time"

	"github.com/republicprotocol/go-do"
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
	Log(log Log) error
}

// PluginOptions are used to Unmarshal plugins from JSON.
type PluginOptions struct {
	File      *FilePluginOptions      `json:"file,omitempty"`
	WebSocket *WebSocketPluginOptions `json:"websocket,omitempty"`
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

// Log an Event.
func (logger *Logger) Log(l Log) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Log(l); err != nil {
			log.Println(err)
		}
	}
}

// Info logs an info Log using a GenericEvent.
func (logger *Logger) Info(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Type:      Info,
		EventType: Generic,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Warn logs a warn Log using a GenericEvent.
func (logger *Logger) Warn(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Type:      Warn,
		EventType: Generic,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Error logs an error Log using a GenericEvent.

func (logger *Logger) Error(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Type:      Error,
		EventType: Generic,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Usage logs an info Log using a UsageEvent.
func (logger *Logger) Usage(cpu, memory float64, network uint64) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Type:      Info,
		EventType: Usage,
		Event: UsageEvent{
			CPU:     cpu,
			Memory:  memory,
			Network: network,
		},
	})
}

// Type defines the different types of Log messages that can be sent.
type Type string

// Values for the LogType.
const (
	Info  = Type("info")
	Warn  = Type("warn")
	Error = Type("error")
)

// EventType defines the different types of Event messages that can be sent in a
// Log.
type EventType string

// Values for the EventType.
const (
	Generic       = EventType("generic")
	Network       = EventType("network")
	Usage         = EventType("usage")
	Ethereum      = EventType("ethereum")
	OrderMatch    = EventType("orderMatch")
	OrderReceived = EventType("orderReceived")
)

// A Log is logged by the Logger using all available Plugins.
type Log struct {
	Timestamp time.Time `json:"timestamp"`
	Type      Type      `json:"type"`
	EventType EventType `json:"eventType"`
	Event     Event     `json:"event"`
}

type Event interface {
	String() string
}

type GenericEvent struct {
	Message string `json:"message"`
}

func (event GenericEvent) String() string {
	return event.Message
}

type UsageEvent struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Network uint64  `json:"network"`
}

func (event UsageEvent) String() string {
	return fmt.Sprintf("cpu = %v; memory = %v; network = %v", event.CPU, event.Memory, event.Network)
}

type OrderMatchEvent struct {
	ID     string `json:"id"`
	BuyID  string `json:"sellId"`
	SellID string `json:"buyId"`
}

func (event OrderMatchEvent) String() string {
	return fmt.Sprintf("order match = (%v, %v)", event.BuyID, event.SellID)
}

type OrderReceivedEvent struct {
	ID         string `json:"id"`
	FragmentID string `json:"fragmentId"`
}

func (event OrderReceivedEvent) String() string {
	return fmt.Sprintf("order recevied = (%v)", event.ID)
}
