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

		It("should show warnings by default", func() {
			start := time.Now()
			msg := "Some information"
			Warn(msg)
			end := time.Now()
			log, err := readTmp()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(start.Before(log.Timestamp)).Should(BeTrue())
			Expect(log.Timestamp.Before(end)).Should(BeTrue())
			Expect(log.Level).Should(Equal(LevelWarn))
			Expect(log.EventType).Should(Equal(TypeGeneric))
			Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
		})

		It("should show errors by default", func() {
			start := time.Now()
			msg := "Some information"
			Error(msg)
			end := time.Now()
			log, err := readTmp()
			Expect(err).ShouldNot(HaveOccurred())

			Expect(start.Before(log.Timestamp)).Should(BeTrue())
			Expect(log.Timestamp.Before(end)).Should(BeTrue())
			Expect(log.Level).Should(Equal(LevelError))
			Expect(log.EventType).Should(Equal(TypeGeneric))
			Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
		})

		It("should not show info or debug messages by default", func() {
			msg := "Some information"
			// Check Info doesn't show
			Info(msg)
			log, err := readTmp()
			Expect(err).ShouldNot(HaveOccurred())
			checkNilLog(log)

			// Check DebugHigh doesn't show
			DebugHigh(msg)
			log, err = readTmp()
			Expect(err).ShouldNot(HaveOccurred())
			checkNilLog(log)

			// Check Debug doesn't show
			Debug(msg)
			log, err = readTmp()
			Expect(err).ShouldNot(HaveOccurred())
			checkNilLog(log)

			// Check DebugLow doesn't show
			DebugLow(msg)
			log, err = readTmp()
			Expect(err).ShouldNot(HaveOccurred())
			checkNilLog(log)
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

func checkNilLog(l Log) {
	Expect(l.Timestamp.IsZero()).Should(BeTrue())
	Expect(l.EventType).Should(Equal(EventType("")))
	Expect(l.Level).Should(Equal(Level(0)))
	Expect(l.Event).Should(BeNil())
}
