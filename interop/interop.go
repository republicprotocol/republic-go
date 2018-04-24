package interop

import (
	"os"
	"strconv"

	"github.com/onsi/ginkgo"
)

// LocalContext allows you to mark a ginkgo context as being local-only.
// It won't run if the CI environment variable is true.
func LocalContext(description string, f func()) {
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
