package stackint_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"

	. "github.com/republicprotocol/republic-go/stackint"
)

var zero = ZERO
var one = Int1024FromUint64(1)
var two = Int1024FromUint64(2)
var three = Int1024FromUint64(3)
var six = Int1024FromUint64(6)
var seven = Int1024FromUint64(7)
var eleven = Int1024FromUint64(11)
var twelve = Int1024FromUint64(12)
var oneWord = Int1024FromUint64(WORDMAX)
var twoPow1023 = TWOPOW1023
var max = MAXINT1024

var _ = Describe("Int1024", func() {
	fmt.Println(TWOPOW1023)
})
