package testutils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

// SkipCIContext can be used instead of Context to skip tests when they are
// being run in a CI environment (to avoid getting flagged for running Bitcoin
// mining software).
func SkipCIContext(description string, f func()) bool {
	if GetCIEnv() {
		return ginkgo.PContext(description, func() {
			ginkgo.It("SKIPPING LOCAL TESTS", func() {})
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

func GanacheBeforeSuite(body interface{}, timeout ...float64) bool {
	fmt.Printf("Ganache is listening on %shttp://localhost:8545%s...\n", green, reset)

	ganache.Start()
	time.Sleep(time.Duration(10) * time.Second)

	conn, err := ganache.Connect("http://localhost:8545")
	if err != nil {
		log.Fatalf("cannot connect to ganache: %v", err)
	}

	if err := ganache.DeployContracts(conn); err != nil {
		log.Fatalf("cannot deploy contracts to ganache: %v", err)
	}
	return ginkgo.BeforeSuite(body, timeout)
}
