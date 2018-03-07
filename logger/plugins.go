package logger

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

	Info(tag, message string) error
	Warn(tag, message string) error
	Error(tag, message string) error
	Usage(cpu float32, memory, network int32) error
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

func (plugin *FilePlugin) Info(tag, message string) error {
	if plugin.File == nil {
		return errors.New("start the file plugin first")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("INFO : " + tag + message + "\n")
	return err
}

func (plugin *FilePlugin) Warn(tag, message string) error {
	if plugin.File == nil {
		return errors.New("start the file plugin first")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("WARN: " + tag + message + "\n")
	return err
}

func (plugin *FilePlugin) Error(tag, message string) error {
	if plugin.File == nil {
		return errors.New("start the file plugin first")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString("ERROR : " + tag + message + "\n")
	return err
}

func (plugin *FilePlugin) Usage(cpu float32, memory, network int32) error {
	if plugin.File == nil {
		return errors.New("start the file plugin first")
	}
	_, err := plugin.File.WriteString(time.Now().Format("2006/01/02 15:04:05 "))
	if err != nil {
		return err
	}
	_, err = plugin.File.WriteString(fmt.Sprintf("USAGE : cpu=%.3f, memory=%d, network=%d \n", cpu, memory, network))
	return err
}

type WebSocketPlugin struct {
	Srv      *http.Server
	Host     string
	Port     string
	Username string
	Password string

	info  chan interface{}
	error chan Error
	warn  chan interface{}
	usage chan Usage
}

func NewWebSocketPlugin(host, port, username, password string) Plugin {
	plugin := &WebSocketPlugin{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		info:     make(chan interface{}, 1),
		error:    make(chan Error, 1),
		warn:     make(chan interface{}, 1),
		usage:    make(chan Usage, 1),
	}
	return plugin
}

func (plugin *WebSocketPlugin) logHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	defer c.Close()
	for {
		request := new(Request)
		err := c.ReadJSON(request)
		if err != nil {
			return
		}

		switch request.Type {
		case "usage":
			go func() {
				for {
					u := <-plugin.usage
					err := c.WriteJSON(u)
					if err != nil {
						break
					}
				}
			}()
		case "event":
			go func() {
				for {
					i := <-plugin.info
					err := c.WriteJSON(i)
					if err != nil {
						break
					}
				}
			}()
		}

		// Broadcast errors
		go func() {
			e := <-plugin.error
			err := c.WriteJSON(e)
			if err != nil {
				return
			}
		}()

		// Broadcast warnings
		go func() {
			warning := <-plugin.warn
			err := c.WriteJSON(warning)
			if err != nil {
				return
			}
		}()

		if err != nil {
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
		err := plugin.Info("INFO", fmt.Sprintf("WebSocket logger listening on %s:%s", plugin.Host, plugin.Port))
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

func (plugin *WebSocketPlugin) Info(tag, message string) error {
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
	err := struct {
		Tag     string
		Message string
	}{
		tag, message,
	}
	if len(plugin.error) == 1 {
		<-plugin.error
	}
	plugin.info <- err
	return nil
}

func (plugin *WebSocketPlugin) Warn(tag, message string) error {
	event := Event{
		Type: "event",
		Time: time.Now(),
		Data: EventData{
			Tag:     tag,
			Level:   "INFO",
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
	usage := Usage{
		Type: "usage",
		Time: time.Now(),
		Data: UsageData{
			Cpu:     cpu,
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
