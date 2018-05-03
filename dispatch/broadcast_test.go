package dispatch_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Broadcaster", func() {

	Context("when shutting down", func() {

		It("should not block existing broadcasts after shutting down", func() {

		})

		It("should not block new broadcasts after shutting down", func() {

		})

		It("should not block existing listeners after shutting down", func() {

		})

		It("should not block new listeners after shutting down", func() {

		})

		It("should not block when shutting down under heavy usage", func() {

		})

	})

	Context("when broadcasting", func() {

		It("should send message from one broadcast to many listeners", func() {

		})

		It("should send messages from many broadcasts to one listener", func() {

		})

		It("should send messages from many broadcasts to many listeners", func() {

		})

	})

})
