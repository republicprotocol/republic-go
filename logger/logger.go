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

// NewLogger returns a new Logger that will start and stop a set of plugins.
func NewLogger(plugins ...Plugin) *Logger {
	return &Logger{
		GuardedObject: do.NewGuardedObject(),
		Plugins:       plugins,
	}
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

func (logger Logger) Error(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Error(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger Logger) Info(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Info(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger Logger) Warn(tag, message string) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Warn(tag, message); err != nil {
			log.Println(err)
		}
	}
}

func (logger Logger) Usage(cpu float32, memory, network int32) {
	for _, plugin := range logger.Plugins {
		if err := plugin.Usage(cpu, memory, network); err != nil {
			log.Println(err)
		}
	}
}

type Request struct {
	Type string      `json:"type"`
	Data RequestData `json:"data"`
}

type RequestData struct {
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Interval int       `json:"interval"`
}

type Usage struct {
	Type string    `json:"type"`
	Time time.Time `json:"timestamp"`
	Data UsageData `json:"data"`
}

type UsageData struct {
	Cpu     float32 `json:"cpu"`
	Memory  int32   `json:"memory"`
	network int32   `json:"network"`
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
