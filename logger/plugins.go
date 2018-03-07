package  logger

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
)

// Plugin
type Plugin interface {
	Start() error
	Stop() error

	Info(info string) error
	Warning(warning string) error
	Error(err error) error
}

// A FilePlugin implements the Plugin interface by logging all events to an
// output file.
type FilePlugin struct {
	Path string
	File *os.File
}

func NewFilePlugin(path string) Plugin {
	return &FilePlugin{
		Path: path,
	}
}

func (plugin *FilePlugin) Start() error {
	var err error
	plugin.File, err = os.OpenFile(plugin.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return err
}

func (plugin *FilePlugin) Stop() error {
	return plugin.File.Close()
}

func (plugin *FilePlugin) Info(info string) error {
	if plugin.File == nil {
		return errors.New("start the file plugin to use Info")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("INFO : ")
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString(info + "\n")
	return err
}

func (plugin *FilePlugin) Warning(warning string) error {
	if plugin.File == nil {
		return errors.New("start the file plugin to use Info")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("WARNING : ")
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString(warning + "\n")
	return err
}

func (plugin *FilePlugin) Error(e error) error {
	if plugin.File == nil {
		return errors.New("start the file plugin to use Info")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("ERROR : ")
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString(err.Error() + "\n")
	return err
}

type WebSocketPlugin struct {
	Srv        *http.Server
	Connection *websocket.Conn
	Host       string
	Port       string
	Username   string
	Password   string
}

func NewWebSocketPlugin(host, port, username, password string) Plugin {
	plugin := &WebSocketPlugin{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
	return plugin
}

func (plugin *WebSocketPlugin) logHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	upgrader := websocket.Upgrader{}
	plugin.Connection, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		plugin.Error(err)
		return
	}

	defer plugin.Connection.Close()
	for {
		request := new(Request)
		err := plugin.Connection.ReadJSON(request)
		if err != nil {
			plugin.Error(err)
			return
		}

		switch request.Type {
		case "usage":
			err = plugin.Connection.WriteJSON(request) // todo
		case "event":
			err = plugin.Connection.WriteJSON(request) // todo
		}

		if err != nil {
			plugin.Error(err)
			return
		}
	}
}

func (plugin *WebSocketPlugin) Start() error {
	plugin.Srv = &http.Server{
		Addr: ":8080",
	}
	http.HandleFunc("/logs", plugin.logHandler)
	go func() {
		err := plugin.Info(fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.Host, plugin.Port))
		if err != nil {
			log.Println(err)
		}
		err = plugin.Srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (plugin *WebSocketPlugin) Stop() error {
	return plugin.Srv.Shutdown(nil)
}

type Message struct {
	Time    string
	Type    string
	Message string
}

func (plugin *WebSocketPlugin) Info(info string) error {
	if plugin.Connection == nil {
		return errors.New("nil connection")
	}

	err := plugin.Connection.WriteJSON(Message{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Type:    "INFO",
		Message: info,
	})
	return err
}

func (plugin *WebSocketPlugin) Error(e error) error {
	if plugin.Connection == nil {
		return errors.New("nil connection")
	}
	err := plugin.Connection.WriteJSON(Message{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Type:    "ERROR",
		Message: e.Error(),
	})
	return err
}

func (plugin *WebSocketPlugin) Warning(warning string) error {
	if plugin.Connection == nil {
		return errors.New("nil connection")
	}
	err := plugin.Connection.WriteJSON(Message{
		Time:    time.Now().Format("2006/01/02 15:04:05"),
		Type:    "warning",
		Message: warning,
	})
	return err
}
