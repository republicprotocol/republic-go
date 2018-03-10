package logger

import (
	"time"

	"github.com/republicprotocol/go-do"
)

// Constant strings for tagging logs.
const (
	TagNetwork   = "network"
	TagCompute   = "compute"
	TagConsensus = "consensus"
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
func (logger *Logger) Start() error {

	for _, plugin := range logger.Plugins {
		err := plugin.Start()
		if err != nil {
			return err
		}
	}
	return nil
}

// Stop stops all the plugins of the logger
func (logger Logger) Stop() error {
	var err error
	for _, plugin := range logger.Plugins {
		if e := plugin.Stop(); e != nil {
			err = e
		}
	}
	return err
}

func (logger Logger) Error(tag, message string) error {
	var err error
	for _, plugin := range logger.Plugins {
		if e := plugin.Error(tag, message); e != nil {
			err = e
		}
	}
	return err
}

func (logger Logger) Info(tag, message string) error {
	var err error
	for _, plugin := range logger.Plugins {
		if e := plugin.Info(tag, message); e != nil {
			err = e
		}
	}
	return err
}

func (logger Logger) Warn(tag, message string) error {
	var err error
	for _, plugin := range logger.Plugins {
		if e := plugin.Warn(tag, message); e != nil {
			err = e
		}
	}
	return err
}

func (logger Logger) Usage(cpu float32, memory, network int32) error {
	var err error
	for _, plugin := range logger.Plugins {
		if e := plugin.Usage(cpu, memory, network); e != nil {
			err = e
		}
	}
	return err
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
