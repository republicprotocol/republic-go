package testutils

import (
	"log"
	"os"
	"strconv"

	"github.com/republicprotocol/republic-go/contract"

	"github.com/onsi/ginkgo"
	"github.com/republicprotocol/republic-go/testutils/ganache"
)

const reset = "\x1b[0m"
const green = "\x1b[32;1m"
const yellow = "\x1b[33;1m"

// GetCIEnv returns true if the CI environment variable is set
func GetCIEnv() bool {
	ciEnv := os.Getenv("CI")
	ci, err := strconv.ParseBool(ciEnv)
	if err != nil {
		ci = false
	}
	return ci
}

// GanacheContext can be used instead of Context to skip tests when they are
// being run in a CI environment (to avoid getting flagged for running Bitcoin
// mining software, and Ganache software).
func GanacheContext(description string, f func()) bool {
	if GetCIEnv() {
		return ginkgo.PContext(description, func() {
			ginkgo.It("Skipping ganache tests...", func() {})
		})
	}
	return ginkgo.Context(description, f)
}

// SkipCIBeforeSuite skips the BeforeSuite, which runs even if there are no tests
func SkipCIBeforeSuite(f func()) bool {
	if !GetCIEnv() {
		return ginkgo.BeforeSuite(f)
	}
	return false
}

// SkipCIAfterSuite skips the AfterSuite, which runs even if there are no tests
func SkipCIAfterSuite(f func()) bool {
	if !GetCIEnv() {
		return ginkgo.AfterSuite(f)
	}
	return false
}

func SkipCIDescribe(d string, f func()) bool {
	if !GetCIEnv() {
		return ginkgo.Describe(d, f)
	}
	return false
}

func GanacheBeforeSuite(body interface{}, timeout ...float64) (contract.Conn, contract.Binder, bool) {
	if !GetCIEnv() {
		conn, err := ganache.StartAndConnect()
		if err != nil {
			log.Fatalf("cannot connect to ganache: %v", err)
		}

		auth := ganache.GenesisTransactor()
		// GasLimit must not be set to 0 to avoid "Out Of Gas" errors
		auth.GasLimit = 3000000
		binder, err := contract.NewBinder(&auth, conn)

		return conn, binder, ginkgo.BeforeSuite(body, timeout...)
	}
	return contract.Conn{}, contract.Binder{}, ginkgo.BeforeSuite(body, timeout...)
}

func GanacheAfterSuite(body interface{}, timeout ...float64) bool {
	if !GetCIEnv() {
		ganache.Stop()
	}
	return ginkgo.AfterSuite(body, timeout...)
}

func GanacheBeforeEach(body interface{}, timeout ...float64) bool {
	_, err := ganache.StartAndConnect()
	if err != nil {
		log.Fatalf("cannot connect to ganache: %v", err)
	}

	return ginkgo.BeforeEach(body, timeout...)
}

func GanacheAfterEach(body interface{}, timeout ...float64) bool {
	return ginkgo.AfterEach(body, timeout...)
}
