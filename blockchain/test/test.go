package test

import (
	"os"
	"strconv"

	"github.com/onsi/ginkgo"
)

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
	} else {
		return ginkgo.Context(description, f)
	}
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
