package node

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// const PATH = "/home/ubuntu/"
const PATH = ""

// Plugin
type Plugin struct {
	io.Writer

	Name     string
	Host     string
	Port     string
	Username string
	Password string
}

// DefaultPlugin returns a default plugin which turns the logs
// into the front-end UI
func DefaultPlugin(username, password string) *Plugin {
	// Setup output log file
	f, err := os.OpenFile(fmt.Sprintf("%sdarknode.log", PATH), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	return &Plugin{
		Name:     "default",
		Username: username,
		Password: password,
		Writer:   f,
	}
}

// NewPlugin creates a new plugin which can be added to the logger
func NewPlugin(name, username, password, host, port string, writer io.Writer) *Plugin {
	return &Plugin{
		Name:     name,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Writer:   writer,
	}
}

type Logger struct {
	Plugins []*Plugin
}

// NewLogger returns a new logger
func NewLogger(plugins ...*Plugin) Logger {
	plugins = append([]*Plugin{DefaultPlugin("", "")}, plugins...)
	return Logger{
		Plugins: plugins,
	}
}

func (logger Logger) Start(address, port string) error {
	http.HandleFunc("/log", logger.logHandler)
	logger.Info(fmt.Sprintf("Websocket starts on %s:%s", address, port))
	return http.ListenAndServe(fmt.Sprintf("%s:%s", address, port), nil)
}

func (logger Logger) logHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	//requestType := r.URL.Query()["type"]
	//if len(requestType) == 1 {
	//
	//}

	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(err)
		return
	}
	defer c.Close()
	for {
		// todo : handle request
		request := new(Request)
		err := c.ReadJSON(request)
		if err != nil {
			logger.Error(err)
			return
		}
		log.Println(request)
		err = c.WriteJSON(request)
		if err != nil {
			logger.Error(err)
			return
		}
	}
}

func (logger Logger) Error(err error) {
	for _, plugin := range logger.Plugins {
		plugin.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
		plugin.Write([]byte("ERROR : "))
		plugin.Write([]byte(err.Error() + "\n"))
	}
}

func (logger Logger) Info(info string) {
	for _, plugin := range logger.Plugins {
		plugin.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
		plugin.Write([]byte("INFO : "))
		plugin.Write([]byte(info + "\n"))
	}
}

func (logger Logger) Debug(debug string) {
	for _, plugin := range logger.Plugins {
		plugin.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
		plugin.Write([]byte("DEBUG : "))
		plugin.Write([]byte(debug + "\n"))
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
	Memory  int     `json:"memory"`
	network int     `json:"network"`
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
