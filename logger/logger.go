package logger

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

// Level defines the different levels of Log messages that can be sent.
type Level uint8

// Values for the LogType.
const (
	LevelError     = Level(1)
	LevelWarn      = Level(2)
	LevelInfo      = Level(3)
	LevelDebugHigh = Level(4)
	LevelDebug     = Level(5)
	LevelDebugLow  = Level(6)
)

// EventType defines the different types of Event messages that can be sent in a
// Log.
type EventType string

// Values for the EventType.
const (
	TypeGeneric        = EventType("generic")
	TypeEpoch          = EventType("epoch")
	TypeUsage          = EventType("usage")
	TypeOrderConfirmed = EventType("orderConfirmed")
	TypeOrderMatch     = EventType("orderMatch")
	TypeOrderReceived  = EventType("orderReceived")
	TypeNetwork        = EventType("network")
	TypeCompute        = EventType("compute")
)

// Log an Event.
func (logger *Logger) Log(l Log) {
	if _, ok := logger.FilterEvents[l.EventType]; !ok && len(logger.FilterEvents) > 0 {
		return
	}
	if l.Level <= logger.FilterLevel {
		l.Tags = logger.Tags
		for _, plugin := range logger.Plugins {
			if err := plugin.Log(l); err != nil {
				log.Println(err)
			}
		}
	}
}

// StdoutLogger logs to standard output and is used as the logger by default
var StdoutLogger = func() *Logger {
	logger, err := NewLogger(Options{
		Plugins: []PluginOptions{
			PluginOptions{File: &FilePluginOptions{Path: "stdout"}},
		},
		FilterLevel: LevelWarn,
	})
	if err != nil {
		panic(fmt.Sprintf("cannot init StdoutLogger: %v", err))
	}
	logger.Start()
	return logger
}()

var defaultLoggerMu = new(sync.RWMutex)
var defaultLogger = StdoutLogger

// SetFilterLevel changes the logging level of the defaultLogger.
func SetFilterLevel(l Level) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.FilterLevel = l
}

// SetFilterEvents sets the defaultLogger to filter specific EventTypes.
func SetFilterEvents(events []EventType) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.FilterEvents = eventListToMap(events)
}

// SetDefaultLogger to a specific Logger object. The previous DefaultLogger
// will be stopped, and the new DefaultLogger will be started.
func SetDefaultLogger(logger *Logger) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Stop()
	defaultLogger = logger
	defaultLogger.Start()
}

// ResetDefaultLogger stops the existing logger and resets the defaultLogger.
func ResetDefaultLogger() {
	SetDefaultLogger(StdoutLogger)
}

// Error logs an error Log using a GenericEvent using the DefaultLogger.
func Error(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Error(message)
}

// Warn logs a warn Log using a GenericEvent using the DefaultLogger.
func Warn(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Warn(message)
}

// Info logs an info Log using a GenericEvent using the DefaultLogger.
func Info(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Info(message)
}

// DebugHigh logs an info Log using a GenericEvent using the DefaultLogger.
func DebugHigh(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.DebugHigh(message)
}

// Debug logs an info Log using a GenericEvent using the DefaultLogger.
func Debug(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Debug(message)
}

// DebugLow logs an info Log using a GenericEvent using the DefaultLogger.
func DebugLow(message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.DebugLow(message)
}

// Epoch logs an info Log using a UsageEvent using the DefaultLogger.
func Epoch(hash [32]byte) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Epoch(hash)
}

// Usage logs an info Log using a UsageEvent using the DefaultLogger.
func Usage(cpu, memory float64, network uint64) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Usage(cpu, memory, network)
}

// OrderConfirmed logs an OrderConfirmedEvent using the DefaultLogger.
func OrderConfirmed(ty Level, orderID string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.OrderConfirmed(ty, orderID)
}

// OrderMatch logs an OrderMatchEvent using the DefaultLogger.
func OrderMatch(ty Level, id, buyID, sellID string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.OrderMatch(ty, id, buyID, sellID)
}

// BuyOrderReceived logs an OrderReceivedEvent using the DefaultLogger.
func BuyOrderReceived(ty Level, id, fragmentID string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.BuyOrderReceived(ty, id, fragmentID)
}

// SellOrderReceived logs an OrderReceivedEvent using the DefaultLogger.
func SellOrderReceived(ty Level, id, fragmentID string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.SellOrderReceived(ty, id, fragmentID)
}

// Network logs a NetworkEvent using the DefaultLogger.
func Network(ty Level, message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Network(ty, message)
}

// Compute logs a ComputeEvent using the DefaultLogger.
func Compute(ty Level, message string) {
	defaultLoggerMu.Lock()
	defer defaultLoggerMu.Unlock()
	defaultLogger.Compute(ty, message)
}

// Logger handles distributing logs to plugins registered with it
type Logger struct {
	Plugins      []Plugin
	Tags         map[string]string
	FilterLevel  Level
	FilterEvents map[EventType]struct{}
}

// Options are used to Unmarshal a Logger from JSON.
type Options struct {
	Plugins      []PluginOptions   `json:"plugins"`
	Tags         map[string]string `json:"tags"`
	FilterLevel  Level             `json:"filterLevel"`
	FilterEvents []EventType       `json:"filterEvents"`
}

// The Plugin interface describes a worker that consumes logs
type Plugin interface {
	Start() error
	Stop() error
	Log(log Log) error
}

// PluginOptions are used to Unmarshal plugins from JSON.
type PluginOptions struct {
	File *FilePluginOptions `json:"file,omitempty"`
}

func eventListToMap(events []EventType) map[EventType]struct{} {
	// Convert the events from an array into a map for O(1) lookup
	eventMap := make(map[EventType]struct{})
	for _, event := range events {
		eventMap[event] = struct{}{}
	}
	return eventMap
}

// NewLogger returns a new Logger that will start and stop a set of plugins.
func NewLogger(options Options) (*Logger, error) {
	logger := &Logger{
		Plugins:      make([]Plugin, 0, len(options.Plugins)),
		Tags:         options.Tags,
		FilterLevel:  options.FilterLevel,
		FilterEvents: eventListToMap(options.FilterEvents),
	}
	for i := range options.Plugins {
		if options.Plugins[i].File != nil {
			plugin := NewFilePlugin(*options.Plugins[i].File)
			logger.Plugins = append(logger.Plugins, plugin)
		}
	}
	return logger, nil
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

// Error logs an error Log using a GenericEvent.
func (logger *Logger) Error(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelError,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Warn logs a warn Log using a GenericEvent.
func (logger *Logger) Warn(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelWarn,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Info logs an info Log using a GenericEvent.
func (logger *Logger) Info(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// DebugHigh logs a debug Log using a GenericEvent.
func (logger *Logger) DebugHigh(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelDebugHigh,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Debug logs a debug Log using a GenericEvent.
func (logger *Logger) Debug(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelDebug,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// DebugLow logs a debug Log using a GenericEvent.
func (logger *Logger) DebugLow(message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelDebugLow,
		EventType: TypeGeneric,
		Event: GenericEvent{
			Message: message,
		},
	})
}

// Epoch logs a EpochEvent.
func (logger *Logger) Epoch(hash [32]byte) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		EventType: TypeEpoch,
		Event: EpochEvent{
			Hash: hash,
		},
	})
}

// Usage logs an info Log using a UsageEvent.
func (logger *Logger) Usage(cpu, memory float64, network uint64) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     LevelInfo,
		EventType: TypeUsage,
		Event: UsageEvent{
			CPU:     cpu,
			Memory:  memory,
			Network: network,
		},
	})
}

// OrderConfirmed logs an OrderConfirmedEvent.
func (logger *Logger) OrderConfirmed(ty Level, orderID string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeOrderConfirmed,
		Event: OrderConfirmedEvent{
			OrderID: orderID,
		},
	})
}

// OrderMatch logs an OrderMatchEvent.
func (logger *Logger) OrderMatch(ty Level, id, buyID, sellID string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeOrderMatch,
		Event: OrderMatchEvent{
			ID:     id,
			BuyID:  buyID,
			SellID: sellID,
		},
	})
}

// BuyOrderReceived logs an OrderReceivedEvent.
func (logger *Logger) BuyOrderReceived(ty Level, id, fragmentID string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeOrderReceived,
		Event: OrderReceivedEvent{
			BuyID:      &id,
			FragmentID: fragmentID,
		},
	})
}

// SellOrderReceived logs an OrderReceivedEvent.
func (logger *Logger) SellOrderReceived(ty Level, id, fragmentID string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeOrderReceived,
		Event: OrderReceivedEvent{
			SellID:     &id,
			FragmentID: fragmentID,
		},
	})
}

// Network logs a NetworkEvent.
func (logger *Logger) Network(ty Level, message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeNetwork,
		Event: NetworkEvent{
			Message: message,
		},
	})
}

// Compute logs a ComputeEvent.
func (logger *Logger) Compute(ty Level, message string) {
	logger.Log(Log{
		Timestamp: time.Now(),
		Level:     ty,
		EventType: TypeCompute,
		Event: ComputeEvent{
			Message: message,
		},
	})
}

func (level Level) String() string {
	switch level {
	case LevelError:
		return "error"
	case LevelWarn:
		return "warn"
	case LevelInfo:
		return "info"
	case LevelDebugHigh:
		fallthrough
	case LevelDebug:
		fallthrough
	case LevelDebugLow:
		return "debug"
	default:
		return ""
	}
}

// A Log is logged by the Logger using all available Plugins.
type Log struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	EventType EventType         `json:"eventType"`
	Event     Event             `json:"event"`
	Tags      map[string]string `json:"tags"`
}

type rawLog struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	EventType EventType         `json:"eventType"`
	Event     json.RawMessage   `json:"event"`
	Tags      map[string]string `json:"tags"`
}

func (log *Log) UnmarshalJSON(data []byte) error {
	rawLog := rawLog{}
	if err := json.Unmarshal(data, &rawLog); err != nil {
		return err
	}

	switch rawLog.EventType {
	case TypeGeneric:
		ev := GenericEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeEpoch:
		ev := EpochEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeUsage:
		ev := UsageEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeOrderConfirmed:
		ev := OrderConfirmedEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeOrderMatch:
		ev := OrderMatchEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeOrderReceived:
		ev := OrderReceivedEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeNetwork:
		ev := NetworkEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	case TypeCompute:
		ev := ComputeEvent{}
		if err := json.Unmarshal(rawLog.Event, &ev); err != nil {
			return err
		}
		log.Event = ev
	}

	log.Timestamp = rawLog.Timestamp
	log.Level = rawLog.Level
	log.EventType = rawLog.EventType
	return nil
}

// The Event interface describes a log event
type Event interface {
	String() string
}

// A GenericEvent logs a string
type GenericEvent struct {
	Message string `json:"message"`
}

func (event GenericEvent) String() string {
	return event.Message
}

// An EpochEvent logs that an epoch transition has been observed
type EpochEvent struct {
	Hash [32]byte `json:"hash"`
}

func (event EpochEvent) String() string {
	return base64.StdEncoding.EncodeToString(event.Hash[:])
}

// UsageEvent logs CPU, Memory and Network usage
type UsageEvent struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory"`
	Network uint64  `json:"network"`
}

func (event UsageEvent) String() string {
	return fmt.Sprintf("cpu = %v; memory = %v; network = %v", event.CPU, event.Memory, event.Network)
}

// OrderConfirmedEvent logs two confirmed orders
type OrderConfirmedEvent struct {
	OrderID string `json:"orderId"`
}

func (event OrderConfirmedEvent) String() string {
	return fmt.Sprintf("confirmation = %s", event.OrderID)
}

// OrderMatchEvent logs two matched orders
type OrderMatchEvent struct {
	ID     string `json:"id"`
	BuyID  string `json:"buyId"`
	SellID string `json:"sellId"`
}

func (event OrderMatchEvent) String() string {
	return fmt.Sprintf("buy = %s; sell = %s", event.BuyID, event.SellID)
}

// OrderReceivedEvent logs a newly received buy or sell fragment
type OrderReceivedEvent struct {
	BuyID      *string `json:"buyId,omitempty"`
	SellID     *string `json:"sellId,omitempty"`
	FragmentID string  `json:"fragmentId"`
}

func (event OrderReceivedEvent) String() string {
	if event.BuyID != nil {
		return "buy = " + *event.BuyID
	}
	if event.SellID != nil {
		return "sell = " + *event.SellID
	}
	return ""
}

// NetworkEvent logs a generic network-related message
type NetworkEvent struct {
	Message string `json:"message"`
}

func (event NetworkEvent) String() string {
	return event.Message
}

// ComputeEvent logs a generic compute-related message
type ComputeEvent struct {
	Message string `json:"message"`
}

func (event ComputeEvent) String() string {
	return event.Message
}
