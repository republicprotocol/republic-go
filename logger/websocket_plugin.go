package logger

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/republicprotocol/go-do"
)

type WebSocketPlugin struct {
	do.GuardedObject

	server   *http.Server
	host     string
	port     string
	username string
	password string

	info  chan interface{}
	err   chan Message
	warn  chan interface{}
	usage chan Usage
}

type WebSocketPluginOptions struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewWebSocketPlugin(webSocketPluginOptions WebSocketPluginOptions) Plugin {
	plugin := &WebSocketPlugin{
		GuardedObject: do.NewGuardedObject(),
		host:          webSocketPluginOptions.Host,
		port:          webSocketPluginOptions.Port,
		username:      webSocketPluginOptions.Username,
		password:      webSocketPluginOptions.Password,
		info:          make(chan interface{}, 1),
		err:           make(chan Message, 1),
		warn:          make(chan interface{}, 1),
		usage:         make(chan Usage, 1),
	}
	return plugin
}

func (plugin *WebSocketPlugin) handler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// The deadlines and intervals for messaging over the socket.
	writeDeadline := 30 * time.Second // We must write notifications within 10 seconds
	pingInterval := 30 * time.Second  // We must ping every 30 seconds
	pongInterval := 60 * time.Second  // We expect a pong every 60 seconds

	// Start the pinger.
	ping := time.NewTicker(pingInterval)
	defer func() {
		ping.Stop()
		conn.Close()
	}()

	// Start the ponger.
	conn.SetReadDeadline(time.Now().Add(pongInterval))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	// Broadcast messages to the WebSocket
	for {
		var val interface{}
		select {
		case val = <-plugin.usage:
		case val = <-plugin.err:
		case val = <-plugin.info:
		case val = <-plugin.warn:
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
		conn.SetWriteDeadline(time.Now().Add(writeDeadline))
		conn.WriteJSON(val)
	}
}

// Start implements the Plugin interface. It starts a WebSocket server.
func (plugin *WebSocketPlugin) Start() error {
	plugin.server = &http.Server{
		Addr: fmt.Sprintf("%s:%s", plugin.host, plugin.port),
	}
	http.HandleFunc("/logs", plugin.handler)
	go func() {
		log.Println(fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.host, plugin.port))
		if err := plugin.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	return nil
}

// Stop implements the Plugin interface. It stops the WebSocket server.
func (plugin *WebSocketPlugin) Stop() error {
	return plugin.server.Shutdown(nil)
}

type Message struct {
	Time    string
	Type    string
	Message string
}

func (plugin *WebSocketPlugin) Info(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	event := Event{
		Type: "event",
		Time: time.Now(),
		Data: EventData{
			Tag:     tag,
			Level:   "INFO",
			Message: message,
		},
	}
	if len(plugin.info) == 1 {
		<-plugin.info
	}
	plugin.info <- event

	return nil
}

func (plugin *WebSocketPlugin) Error(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	msg := Message{
		time.Now().Format("2006/01/02 15:04:05 "), tag, message,
	}
	if len(plugin.err) == 1 {
		<-plugin.err
	}
	plugin.err <- msg
	return nil
}

func (plugin *WebSocketPlugin) Warn(tag, message string) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	event := Event{
		Type: "event",
		Time: time.Now(),
		Data: EventData{
			Tag:     tag,
			Level:   "WARN",
			Message: message,
		},
	}
	if len(plugin.warn) == 1 {
		<-plugin.warn
	}
	plugin.warn <- event

	return nil
}

func (plugin *WebSocketPlugin) Usage(cpu float32, memory, network int32) error {
	plugin.Enter(nil)
	defer plugin.Exit()

	usage := Usage{
		Type: "usage",
		Time: time.Now(),
		Data: UsageData{
			CPU:     cpu,
			Memory:  memory,
			Network: network,
		},
	}
	if len(plugin.usage) == 1 {
		<-plugin.usage
	}
	plugin.usage <- usage
	return nil
}
