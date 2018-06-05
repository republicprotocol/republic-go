package logger_test

import (
	"encoding/json"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/logger"
)

var _ = Describe("Logger", func() {

	Context("When using a file plugin", func() {

		BeforeEach(func() {
			// Make the temp directory and initialise the logger
			err := makeTmp()
			Expect(err).ShouldNot(HaveOccurred())
			_, err = initFileLogger()
			Expect(err).ShouldNot(HaveOccurred())
		})

		AfterEach(func() {
			ResetDefaultLogger()
			err := removeTmp()
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("by default it should show errors", func() {
			testGenericError()
		})

		It("by default it should show warnings", func() {
			testGenericWarn()
		})

		It("by default it should not show Info or Debug messages", func() {
			msg := "Some information"
			// Check Info doesn't show
			Info(msg)
			checkNilLog()
			// Check DebugHigh doesn't show
			DebugHigh(msg)
			checkNilLog()
			// Check Debug doesn't show
			Debug(msg)
			checkNilLog()
			// Check DebugLow doesn't show
			DebugLow(msg)
			checkNilLog()
		})

		It("SetFilterLevel(0) should filter all messages", func() {
			msg := "Some information"
			// Filter everything
			SetFilterLevel(0)
			Error(msg)
			checkNilLog()
			Warn(msg)
			checkNilLog()
			Info(msg)
			checkNilLog()
			DebugHigh(msg)
			checkNilLog()
			Debug(msg)
			checkNilLog()
			DebugLow(msg)
			checkNilLog()
		})

		It("SetFilterLevel(LevelError) should only show Error logs", func() {
			msg := "Some information"
			// Filter all but Error
			SetFilterLevel(LevelError)
			Warn(msg)
			checkNilLog()
			Info(msg)
			checkNilLog()
			DebugHigh(msg)
			checkNilLog()
			Debug(msg)
			checkNilLog()
			DebugLow(msg)
			checkNilLog()
			// We should get these logs
			testGenericError()
		})

		It("SetFilterLevel(LevelWarn) should only show Error and Warn logs", func() {
			msg := "Some information"
			// Filter all but Error and Warn
			SetFilterLevel(LevelWarn)
			Info(msg)
			checkNilLog()
			DebugHigh(msg)
			checkNilLog()
			Debug(msg)
			checkNilLog()
			DebugLow(msg)
			checkNilLog()
			// We should get these logs
			testGenericError()
			resetTmp()
			testGenericWarn()
		})

		It("SetFilterLevel(LevelInfo) should only show Error, Warn, and Info logs", func() {
			setup := func() {
				err := resetTmp()
				Expect(err).ShouldNot(HaveOccurred())
				SetFilterLevel(LevelInfo)
			}
			setup()
			msg := "Some information"
			// Filter all but Error, Warn, and Info
			DebugHigh(msg)
			checkNilLog()
			Debug(msg)
			checkNilLog()
			DebugLow(msg)
			checkNilLog()
			// We should get these logs
			testGenericError()
			setup()
			testGenericWarn()
			setup()
			testGenericInfo()
		})

		It("SetFilterLevel(LevelDebugHigh) should filter Debug and DebugLow logs", func() {
			setup := func() {
				err := resetTmp()
				Expect(err).ShouldNot(HaveOccurred())
				SetFilterLevel(LevelDebugHigh)
			}
			setup()
			msg := "Some information"
			// Filter all but Error, Warn, Info, and DebugHigh
			Debug(msg)
			checkNilLog()
			DebugLow(msg)
			checkNilLog()
			// We should get these logs
			testGenericError()
			setup()
			testGenericWarn()
			setup()
			testGenericInfo()
			setup()
			testGenericDebugHigh()
		})

		It("SetFilterLevel(LevelDebug) should filter DebugLow logs", func() {
			setup := func() {
				err := resetTmp()
				Expect(err).ShouldNot(HaveOccurred())
				SetFilterLevel(LevelDebug)
			}
			setup()
			msg := "Some information"
			// Filter DebugLow
			DebugLow(msg)
			checkNilLog()
			// We should get these logs
			testGenericError()
			setup()
			testGenericWarn()
			setup()
			testGenericInfo()
			setup()
			testGenericDebugHigh()
			setup()
			testGenericDebug()
		})

		It("SetFilterLevel(LevelDebugLow) should show all logs", func() {
			setup := func() {
				err := resetTmp()
				Expect(err).ShouldNot(HaveOccurred())
				SetFilterLevel(LevelDebugLow)
			}
			setup()
			// We should get these logs
			testGenericError()
			setup()
			testGenericWarn()
			setup()
			testGenericInfo()
			setup()
			testGenericDebugHigh()
			setup()
			testGenericDebug()
			setup()
			testGenericDebugLow()
		})

	})

})

const tmpFolder = "./tmp/"
const tmpFile = tmpFolder + "test.log"

func makeTmp() error {
	return os.MkdirAll(tmpFolder, os.ModePerm)
}

func removeTmp() error {
	return os.RemoveAll(tmpFolder)
}

func readTmp() (Log, error) {
	file, err := os.OpenFile(tmpFile, os.O_RDONLY, 0640)
	//defer file.Close()
	log := Log{}
	if err != nil {
		return log, err
	}
	// Marshal the file back to a log
	json.NewDecoder(file).Decode(&log)
	return log, nil
}

func resetTmp() error {
	ResetDefaultLogger()
	err := removeTmp()
	if err != nil {
		return err
	}
	err = makeTmp()
	if err != nil {
		return err
	}
	_, err = initFileLogger()
	if err != nil {
		return err
	}
	return nil
}

func initFileLogger() (*Logger, error) {
	logger, err := NewLogger(Options{
		Plugins: []PluginOptions{
			PluginOptions{File: &FilePluginOptions{Path: tmpFile}, WebSocket: nil},
		},
		FilterLevel: LevelWarn,
	})
	if err != nil {
		return logger, err
	}
	SetDefaultLogger(logger)
	return logger, nil
}

func checkNilLog() {
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(log.Timestamp.IsZero()).Should(BeTrue())
	Expect(log.EventType).Should(Equal(EventType("")))
	Expect(log.Level).Should(Equal(Level(0)))
	Expect(log.Event).Should(BeNil())
}

func testGenericError() {
	start := time.Now()
	msg := "Some information"
	Error(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelError))
	Expect(log.Level.String()).Should(BeEquivalentTo("error"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericWarn() {
	start := time.Now()
	msg := "Some information"
	Warn(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelWarn))
	Expect(log.Level.String()).Should(BeEquivalentTo("warn"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericInfo() {
	start := time.Now()
	msg := "Some information"
	Info(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelInfo))
	Expect(log.Level.String()).Should(BeEquivalentTo("info"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebugHigh() {
	start := time.Now()
	msg := "Some information"
	DebugHigh(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebugHigh))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebug() {
	start := time.Now()
	msg := "Some information"
	Debug(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebug))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebugLow() {
	start := time.Now()
	msg := "Some information"
	DebugLow(msg)
	end := time.Now()
	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebugLow))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}
