package test

import (
	"os"
	"strconv"

	"github.com/onsi/ginkgo"
)

// SkipCiContext can be used instead of Context to skip tests when they are
// being run in a CI environment (to avoid getting flagged for running Bitcoin
// mining software).
func SkipCiContext(description string, f func()) {
	var local bool

	ciEnv := os.Getenv("CI")
	ci, err := strconv.ParseBool(ciEnv)
	if err != nil {
		ci = false
	}

	// Assume tests are running locally if CI environment variable is not defined
	local = !ci

	if local {
		ginkgo.Context(description, f)
	} else {
		ginkgo.PContext(description, func() {
			ginkgo.It("SKIPPING LOCAL TESTS", func() {})
		})
	}
}
