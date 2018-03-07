package logger

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Plugin
type Plugin interface {
	Start() error
	Stop() error

	Info(info string)
	Warning(warning string)
	Error(err error)
}

// A FilePlugin implements the Plugin interface by logging all events to an
// output file.
type FilePlugin struct {
	Path string
	File *os.File
}

func NewFilePlugin(path string) Plugin {
	return FilePlugin{
		Path: path,
	}
}

func (plugin FilePlugin) Start() error {
	var err error
	plugin.File, err = os.OpenFile(plugin.Path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return err
}

func (plugin FilePlugin) Stop() error {
	return plugin.File.Close()
}

func (plugin FilePlugin) Info(info string) {
	plugin.File.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
	plugin.File.Write([]byte("INFO : "))
	plugin.File.Write([]byte(info + "\n"))
}

func (plugin FilePlugin) Warning(warning string) {
	plugin.File.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
	plugin.File.Write([]byte("WARNING : "))
	plugin.File.Write([]byte(warning + "\n"))
}

func (plugin FilePlugin) Error(err error) {
	plugin.File.Write([]byte(time.Now().Format("2006/01/02 15:04:05 ")))
	plugin.File.Write([]byte("ERROR : "))
	plugin.File.Write([]byte(err.Error() + "\n"))
}

type WebSocketPlugin struct {
	Srv      *http.Server
	Host     string
	Port     string
	Username string
	Password string
	Handler  func(w http.ResponseWriter, r *http.Request)
}

func NewWebSocketPlugin(host, port, username, password string) Plugin {
	plugin := WebSocketPlugin{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	plugin.Handler = plugin.logHandler
	return plugin
}

func (plugin WebSocketPlugin)logHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		plugin.Error(err)
		return
	}

	defer c.Close()
	for {
		// todo : handle request
		request := new(Request)
		err := c.ReadJSON(request)
		if err != nil {
			plugin.Error(err)
			return
		}

		err = c.WriteJSON(request)
		if err != nil {
			plugin.Error(err)
			return
		}
	}
}

func (plugin WebSocketPlugin) Start() error {
	plugin.Srv = &http.Server{
		Addr:":8080",
	}
	http.HandleFunc("/logs", plugin.Handler)
	go func() {
		plugin.Info(fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.Host, plugin.Port))
		plugin.Srv.ListenAndServe()
	}()
	return nil
}

func (plugin WebSocketPlugin) Stop() error {
	return plugin.Srv.Shutdown(nil)
}

func (plugin WebSocketPlugin) Info(info string) {

}

func (plugin WebSocketPlugin) Error(err error) {

}

func (plugin WebSocketPlugin) Warning(warning string) {

}
