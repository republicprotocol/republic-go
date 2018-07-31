package testutils_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/testutils"
)

var _ = Describe("Testutils", func() {

	// 1. All helpers functions (cannot use non-argument variables other than
	// the constants defined above). Any code that is re-used in this file must
	// be moved to a helper function.
	//
	// 2. Any code that could be re-used in multiple test files must go in the
	// testutils package. For instance; mocks of services and random generation
	// of addresses. The filenames of the mocks and random generators in testutils
	// package should be prepended with 'mock_' and 'rand_', respectively.

	/****************************************************************************/
	/* 				Examples for helper functions								*/
	/****************************************************************************/

	// setupCI sets the CI OS parameter to true.
	setupCI := func() {
		os.Setenv("CI", "true")
	}

	// setupLocal sets the CI OS parameter to false.
	setupLocal := func() {
		os.Setenv("CI", "false")
	}

	// 3. Think about variables in this space as global variables. It is very
	// rare to need to declare anything in this space (other than constants
	// that parameterise the test — e.g. number of nodes in the test network).

	/****************************************************************************/
	/* 				Examples for global constants								*/
	/****************************************************************************/

	const (
		DebugMode          = false
		NumberOfIterations = 420
	)

	// 4. Contexts blocks are used to test the component in
	// different situations. The description of Context blocks must begin with
	// 'when' followed by the scenario being tested. Each Context block can have
	// separate BeforeEach and AfterEach blocks that would be responsible for
	// common set-up and tear-down operations for all tests in that situation.

	Context("when tests are run in CI", func() {

		// Code that must be run for setting up tests goes here.
		BeforeEach(func() {
			setupCI()
		})

		// Any clean-up after every test will be present here.
		AfterEach(func() {
		})

		Context("when using SkipCI prefix", func() {

			It("should skip BeforeSuite", func(done Done) {
				defer close(done)

				Expect(SkipCIBeforeSuite(func() {})).To(BeFalse())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)

			It("should skip AfterSuite", func(done Done) {
				defer close(done)

				Expect(SkipCIAfterSuite(func() {})).To(BeFalse())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)

			It("should skip Describe", func(done Done) {
				defer close(done)

				Expect(SkipCIDescribe("", func() {})).To(BeFalse())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)
		})
	})

	Context("when tests are run locally", func() {
		BeforeEach(func() {
			setupLocal()
		})

		AfterEach(func() {
		})

		Context("when using SkipCI prefix", func() {

			It("should not skip BeforeSuite", func(done Done) {
				defer close(done)

				Expect(SkipCIBeforeSuite(func() {})).To(BeTrue())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)

			It("should not skip AfterSuite", func(done Done) {
				defer close(done)

				Expect(SkipCIAfterSuite(func() {})).To(BeTrue())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)

			It("should not skip Describe", func(done Done) {
				defer close(done)

				Expect(SkipCIDescribe("", func() {})).To(BeTrue())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)

			It("should not skip ganache tests", func(done Done) {
				defer close(done)

				Expect(GanacheContext("test description", func() {})).To(BeTrue())
			}, 1.0 /* 3. Timeout in seconds, defaults to 1.0 */)
		})
	})
})
