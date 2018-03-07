package logger

import (
	"time"
)

type Logger struct {
	Plugins []Plugin
}

// NewLogger returns a new Logger that will start and stop a set of plugins.
func NewLogger(plugins ...Plugin) Logger {
	return Logger{
		Plugins: plugins,
	}
}

func (logger Logger) Start() error {
	for _, plugin := range logger.Plugins {
		var err error
		switch plugin.(type){
		case FilePlugin:
			err = plugin.(FilePlugin).Start()
		case WebSocketPlugin:
			err = plugin.(WebSocketPlugin).Start()
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (logger Logger) Stop() error {
	for _, plugin := range logger.Plugins {
		var err error
		switch plugin.(type){
		case FilePlugin:
			err = plugin.(FilePlugin).Stop()
		case WebSocketPlugin:
			err = plugin.(WebSocketPlugin).Stop()
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (logger Logger) Error(err error) {
	for _, plugin := range logger.Plugins {
		plugin.Error(err)
	}
}

func (logger Logger) Info(info string) {
	for _, plugin := range logger.Plugins {
		switch plugin.(type){
		case FilePlugin:
			plugin.(FilePlugin).Info(info)
		case WebSocketPlugin:
			plugin.(WebSocketPlugin).Info(info)
		}
	}
}

func (logger Logger) Warning(warning string) {
	for _, plugin := range logger.Plugins {
		switch plugin.(type){
		case FilePlugin:
			plugin.(FilePlugin).Warning(warning)
		case WebSocketPlugin:
			plugin.(WebSocketPlugin).Warning(warning)
		}
	}
}

type Request struct {
	Type string `json:"type"`
	Data RequestData `json:"data"`
}

type RequestData struct {
	Start time.Time `json:"start"`
	End time.Time `json:"end"`
	Interval int `json:"interval"`
}

type Usage struct {
	Type string `json:"type"`
	Time time.Time `json:"timestamp"`
	Data UsageData `json:"data"`
}

type UsageData struct {
	Cpu float32 `json:"cpu"`
	Memory int `json:"memory"`
	network int `json:"network"`
}

type Event struct {
	Type string `json:"type"`
	Time time.Time `json:"timestamp"`
	Data EventData `json:"data"`
}

type EventData struct {
	Tag string `json:"tag"`
	Level string `json:"level"`
	Message string `json:"message"`
}
