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
	msg := "Some information"

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

		Context("and when FilterLevel is set to 0", func() {
			BeforeEach(func() {
				SetFilterLevel(0)
			})

			It("should filter Errors", func() {
				Error(msg)
				checkNilLog()
			})

			It("should filter Warns", func() {
				Warn(msg)
				checkNilLog()
			})

			It("should filter Info", func() {
				Info(msg)
				checkNilLog()
			})

			It("should filter DebugHigh", func() {
				DebugHigh(msg)
				checkNilLog()
			})

			It("should filter Debug", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})
		})

		Context("and when FilterLevel is set to Error", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelError)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should filter Warns", func() {
				Warn(msg)
				checkNilLog()
			})

			It("should filter Info", func() {
				Info(msg)
				checkNilLog()
			})

			It("should filter DebugHigh", func() {
				DebugHigh(msg)
				checkNilLog()
			})

			It("should filter Debug", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})
		})

		Context("and when FilterLevel is set to Warn", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelWarn)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should show Warns", func() {
				testGenericWarn()
			})

			It("should filter Info", func() {
				Info(msg)
				checkNilLog()
			})

			It("should filter DebugHigh", func() {
				DebugHigh(msg)
				checkNilLog()
			})

			It("should filter Debug", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})
		})

		Context("and when FilterLevel is set to Info", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelInfo)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should show Warns", func() {
				testGenericWarn()
			})

			It("should show Info", func() {
				testGenericInfo()
			})

			It("should filter DebugHigh", func() {
				DebugHigh(msg)
				checkNilLog()
			})

			It("should filter Debug", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})
		})

		Context("and when FilterLevel is set to DebugHigh", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebugHigh)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should show Warns", func() {
				testGenericWarn()
			})

			It("should show Info", func() {
				testGenericInfo()
			})

			It("should show DebugHigh", func() {
				testGenericDebugHigh()
			})

			It("should filter Debug", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})

		})

		Context("and when FilterLevel is set to Debug", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebug)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should show Warns", func() {
				testGenericWarn()
			})

			It("should show Info", func() {
				testGenericInfo()
			})

			It("should show DebugHigh", func() {
				testGenericDebugHigh()
			})

			It("should show Debug", func() {
				testGenericDebug()
			})

			It("should filter DebugLow", func() {
				DebugLow(msg)
				checkNilLog()
			})

		})

		Context("and when FilterLevel is set to DebugLow", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebugLow)
			})

			It("should show Errors", func() {
				testGenericError()
			})

			It("should show Warns", func() {
				testGenericWarn()
			})

			It("should show Info", func() {
				testGenericInfo()
			})

			It("should show DebugHigh", func() {
				testGenericDebugHigh()
			})

			It("should show Debug", func() {
				testGenericDebug()
			})

			It("should show DebugLow", func() {
				testGenericDebugLow()
			})

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
