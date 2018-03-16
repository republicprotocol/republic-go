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

	logger   *Logger
	server   *http.Server
	host     string
	port     string
	username string
	password string

	logs chan Log
}

type WebSocketPluginOptions struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewWebSocketPlugin(logger *Logger, webSocketPluginOptions WebSocketPluginOptions) Plugin {
	plugin := &WebSocketPlugin{
		GuardedObject: do.NewGuardedObject(),
		logger:        logger,
		host:          webSocketPluginOptions.Host,
		port:          webSocketPluginOptions.Port,
		username:      webSocketPluginOptions.Username,
		password:      webSocketPluginOptions.Password,
		logs:          make(chan Log, 100),
	}
	return plugin
}

func (plugin *WebSocketPlugin) handler(w http.ResponseWriter, r *http.Request) {
	defer func() {
		go plugin.logger.Network(Info, "WebSocket logger connection released")
	}()
	go plugin.logger.Network(Info, "WebSocket logger connection acquired")

	upgrader := websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	// The deadlines and intervals for messaging over the socket
	writeDeadline := 10 * time.Second // We must write notifications within 10 seconds
	pingInterval := 30 * time.Second  // We must ping every 30 seconds
	pongInterval := 60 * time.Second  // We expect a pong every 60 seconds

	// Start the pinger
	ping := time.NewTicker(pingInterval)
	defer func() {
		ping.Stop()
		conn.Close()
	}()

	// Start the ponger
	conn.SetReadDeadline(time.Now().Add(pongInterval))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongInterval))
		return nil
	})

	// Broadcast logs to the WebSocket
	for {
		select {
		case val := <-plugin.logs:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			conn.WriteJSON(val)
		case <-ping.C:
			conn.SetWriteDeadline(time.Now().Add(writeDeadline))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Start implements the Plugin interface. It starts a WebSocket server.
func (plugin *WebSocketPlugin) Start() error {
	plugin.server = &http.Server{
		Addr: fmt.Sprintf("%s:%s", plugin.host, plugin.port),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/logs", plugin.handler)
	plugin.server.Handler = mux
	go func() {
		if err := plugin.server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()
	go plugin.logger.Network(Info, fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.host, plugin.port))
	return nil
}

// Stop implements the Plugin interface. It stops the WebSocket server.
func (plugin *WebSocketPlugin) Stop() error {
	return plugin.server.Shutdown(nil)
}

// Log implements the Plugin interface.
func (plugin *WebSocketPlugin) Log(l Log) error {
	written := false
	for !written {
		select {
		case plugin.logs <- l:
			// Write was successful
			written = true
		default:
			// Logging queue was full, drop the oldest message and loop
			<-plugin.logs
		}
	}
	return nil
}
