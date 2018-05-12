package ledger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLedger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ledger Suite")
}
