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

	Context("when using a file plugin", func() {

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

		Context("when the filter level is set to 0", func() {
			BeforeEach(func() {
				SetFilterLevel(0)
			})

			It("should filter error messages", func() {
				Error(msg)
				checkNilLog()
			})

			It("should filter warn messages", func() {
				Warn(msg)
				checkNilLog()
			})

			It("should filter info messages", func() {
				Info(msg)
				checkNilLog()
			})

			It("should filter high-level debug messages", func() {
				DebugHigh(msg)
				checkNilLog()
			})

			It("should filter regular debug messages", func() {
				Debug(msg)
				checkNilLog()
			})

			It("should filter low-level debug messages", func() {
				DebugLow(msg)
				checkNilLog()
			})
		})

		Context("when the filter level is set to 1", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelError)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should filter warn messages", func() {
				testGenericWarn(false)
			})

			It("should filter info messages", func() {
				testGenericInfo(false)
			})

			It("should filter high-level debug messages", func() {
				testGenericDebugHigh(false)
			})

			It("should filter regular debug messages", func() {
				testGenericDebug(false)
			})

			It("should filter low-level debug messages", func() {
				testGenericDebugLow(false)
			})
		})

		Context("when the filter level is set to 2", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelWarn)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should show warn messages", func() {
				testGenericWarn(true)
			})

			It("should filter info messages", func() {
				testGenericInfo(false)
			})

			It("should filter high-level debug messages", func() {
				testGenericDebugHigh(false)
			})

			It("should filter regular debug messages", func() {
				testGenericDebug(false)
			})

			It("should filter low-level debug messages", func() {
				testGenericDebugLow(false)
			})
		})

		Context("when the filter level is set to 3", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelInfo)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should show warn messages", func() {
				testGenericWarn(true)
			})

			It("should show info messages", func() {
				testGenericInfo(true)
			})

			It("should filter high-level debug messages", func() {
				testGenericDebugHigh(false)
			})

			It("should filter regular debug messages", func() {
				testGenericDebug(false)
			})

			It("should filter low-level debug messages", func() {
				testGenericDebugLow(false)
			})
		})

		Context("when the filter level is set to 4", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebugHigh)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should show warn messages", func() {
				testGenericWarn(true)
			})

			It("should show info messages", func() {
				testGenericInfo(true)
			})

			It("should show high-level debug messages", func() {
				testGenericDebugHigh(true)
			})

			It("should filter regular debug messages", func() {
				testGenericDebug(false)
			})

			It("should filter low-level debug messages", func() {
				testGenericDebugLow(false)
			})

		})

		Context("when the filter level is set to 5", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebug)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should show warn messages", func() {
				testGenericWarn(true)
			})

			It("should show info messages", func() {
				testGenericInfo(true)
			})

			It("should show high-level debug messages", func() {
				testGenericDebugHigh(true)
			})

			It("should show debug messages", func() {
				testGenericDebug(true)
			})

			It("should filter low-level debug messages", func() {
				testGenericDebugLow(false)
			})

		})

		Context("when the filter level is set to 6", func() {
			BeforeEach(func() {
				SetFilterLevel(LevelDebugLow)
			})

			It("should show error messages", func() {
				testGenericError(true)
			})

			It("should show warn messages", func() {
				testGenericWarn(true)
			})

			It("should show info messages", func() {
				testGenericInfo(true)
			})

			It("should show high-level debug messages", func() {
				testGenericDebugHigh(true)
			})

			It("should show regular debug messages", func() {
				testGenericDebug(true)
			})

			It("should show low-level debug messages", func() {
				testGenericDebugLow(true)
			})

		})

		Context("when we set event filter to whitelist generic events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeGeneric})
				SetFilterLevel(6)
			})

			It("should filter non-generic messages", func() {
				testEpoch(false)
				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

			It("should show all generic messages", func() {
				testGenericError(true)
			})

			It("should show generic warn messages", func() {
				testGenericWarn(true)
			})

			It("should show generic info messages", func() {
				testGenericInfo(true)
			})

			It("should show generic high-level debug messages", func() {
				testGenericDebugHigh(true)
			})

			It("should show generic regular debug messages", func() {
				testGenericDebug(true)
			})

			It("should show generic low-level debug messages", func() {
				testGenericDebugLow(true)
			})

		})

		Context("when we set event filter to whitelist epoch events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeEpoch})
				SetFilterLevel(6)
			})

			It("should show epoch messages", func() {
				testEpoch(true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist usage events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeUsage})
				SetFilterLevel(6)
			})

			It("should show usage messages", func() {
				testUsage(true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist order confirmed events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeOrderConfirmed})
				SetFilterLevel(6)
			})

			It("should show order confirmed messages", func() {
				testOrderConfirmed(LevelError, true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testUsage(false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist order matched events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeOrderMatch})
				SetFilterLevel(6)
			})

			It("should show order confirmed messages", func() {
				testOrderMatch(LevelWarn, true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist order received events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeOrderReceived})
				SetFilterLevel(6)
			})

			It("should show buy order confirmed messages", func() {
				testBuyOrderReceived(LevelInfo, true)
			})

			It("should show sell order confirmed messages", func() {
				testSellOrderReceived(LevelDebugHigh, true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testNetwork(LevelDebug, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist network events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeNetwork})
				SetFilterLevel(6)
			})

			It("should show network messages", func() {
				testNetwork(LevelDebug, true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testCompute(LevelDebugLow, false)
			})

		})

		Context("when we set event filter to whitelist compute events", func() {
			BeforeEach(func() {
				SetFilterEvents([]EventType{TypeCompute})
				SetFilterLevel(6)
			})

			It("should show compute messages", func() {
				testCompute(LevelDebugLow, true)
			})

			It("should filter non-usage messages", func() {
				testGenericError(false)
				testGenericWarn(false)
				testGenericInfo(false)
				testGenericDebugHigh(false)
				testGenericDebug(false)
				testGenericDebugLow(false)

				testEpoch(false)
				testUsage(false)
				testOrderConfirmed(LevelError, false)
				testOrderMatch(LevelWarn, false)
				testBuyOrderReceived(LevelInfo, false)
				testSellOrderReceived(LevelDebugHigh, false)
				testNetwork(LevelDebug, false)
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

func initFileLogger() (*Logger, error) {
	logger, err := NewLogger(Options{
		Plugins: []PluginOptions{
			PluginOptions{File: &FilePluginOptions{Path: tmpFile}},
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

func testGenericError(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Error(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelError))
	Expect(log.Level.String()).Should(BeEquivalentTo("error"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericWarn(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Warn(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelWarn))
	Expect(log.Level.String()).Should(BeEquivalentTo("warn"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericInfo(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Info(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelInfo))
	Expect(log.Level.String()).Should(BeEquivalentTo("info"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebugHigh(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	DebugHigh(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebugHigh))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebug(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Debug(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebug))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testGenericDebugLow(shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	DebugLow(msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())

	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelDebugLow))
	Expect(log.Level.String()).Should(BeEquivalentTo("debug"))
	Expect(log.EventType).Should(Equal(TypeGeneric))
	Expect(log.Event.(GenericEvent).Message).Should(Equal(msg))
}

func testEpoch(shouldLog bool) {
	var hash [32]byte
	testHash := []byte("Here is a string....")
	copy(hash[:], testHash)

	start := time.Now()
	Epoch(hash)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelInfo))
	Expect(log.Level.String()).Should(BeEquivalentTo("info"))
	Expect(log.EventType).Should(Equal(TypeEpoch))
	Expect(log.Event.(EpochEvent).Hash).Should(Equal(hash))
}

func testUsage(shouldLog bool) {
	cpu := 3.14
	memory := 1.23
	var network uint64
	network = 37
	start := time.Now()
	Usage(cpu, memory, network)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(LevelInfo))
	Expect(log.Level.String()).Should(BeEquivalentTo("info"))
	Expect(log.EventType).Should(Equal(TypeUsage))
	Expect(log.Event.(UsageEvent).CPU).Should(Equal(cpu))
	Expect(log.Event.(UsageEvent).Memory).Should(Equal(memory))
	Expect(log.Event.(UsageEvent).Network).Should(Equal(network))
}

func testOrderConfirmed(l Level, shouldLog bool) {
	orderID := "someOrderId"
	start := time.Now()
	OrderConfirmed(l, orderID)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeOrderConfirmed))
	Expect(log.Event.(OrderConfirmedEvent).OrderID).Should(Equal(orderID))
}

func testOrderMatch(l Level, shouldLog bool) {
	id := "someid"
	buyID := "someBuyId"
	sellID := "someSellId"

	start := time.Now()
	OrderMatch(l, id, buyID, sellID)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeOrderMatch))
	Expect(log.Event.(OrderMatchEvent).ID).Should(Equal(id))
	Expect(log.Event.(OrderMatchEvent).BuyID).Should(Equal(buyID))
	Expect(log.Event.(OrderMatchEvent).SellID).Should(Equal(sellID))
}

func testBuyOrderReceived(l Level, shouldLog bool) {
	buyID := "someid"
	fragmentID := "someFragId"
	start := time.Now()
	BuyOrderReceived(l, buyID, fragmentID)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeOrderReceived))
	Expect(*log.Event.(OrderReceivedEvent).BuyID).Should(Equal(buyID))
	Expect(log.Event.(OrderReceivedEvent).FragmentID).Should(Equal(fragmentID))
}

func testSellOrderReceived(l Level, shouldLog bool) {
	sellID := "someid"
	fragmentID := "someFragId"
	start := time.Now()
	SellOrderReceived(l, sellID, fragmentID)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeOrderReceived))
	Expect(*log.Event.(OrderReceivedEvent).SellID).Should(Equal(sellID))
	Expect(log.Event.(OrderReceivedEvent).FragmentID).Should(Equal(fragmentID))
}

func testNetwork(l Level, shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Network(l, msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeNetwork))
	Expect(log.Event.(NetworkEvent).Message).Should(Equal(msg))
}

func testCompute(l Level, shouldLog bool) {
	start := time.Now()
	msg := "Some information"
	Compute(l, msg)
	end := time.Now()

	if !shouldLog {
		checkNilLog()
		return
	}

	log, err := readTmp()
	Expect(err).ShouldNot(HaveOccurred())
	Expect(start.Before(log.Timestamp)).Should(BeTrue())
	Expect(log.Timestamp.Before(end)).Should(BeTrue())
	Expect(log.Level).Should(Equal(l))
	Expect(log.EventType).Should(Equal(TypeCompute))
	Expect(log.Event.(ComputeEvent).Message).Should(Equal(msg))
}
