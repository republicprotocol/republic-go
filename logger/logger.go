package logger

import (
	"time"
)

type Logger struct {
	Plugins []Plugin
}

// NewLogger returns a new Logger that will start and stop a set of plugins.
func NewLogger(plugins ...Plugin) *Logger {
	return &Logger{
		Plugins: plugins,
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
	for _, plugin := range logger.Plugins {
		plugin.Stop()
	}
	return nil
}

// Error outputs the error though each plugin
func (logger Logger) Error(tag, message string) {
	for _, plugin := range logger.Plugins {
		plugin.Error(tag, message)
	}
}

// Info outputs the info though each plugin
func (logger Logger) Info(tag, message string) {
	for _, plugin := range logger.Plugins {
		plugin.Info(tag, message)
	}
}

// Warning outputs the warning though each plugin
func (logger Logger) Warning(tag, message string) {
	for _, plugin := range logger.Plugins {
		plugin.Warn(tag, message)
	}
}

// Warning outputs the warning though each plugin
func (logger Logger) Usage(cpu float32, memory, network int32) {
	for _, plugin := range logger.Plugins {
		plugin.Usage(cpu, memory, network)
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
