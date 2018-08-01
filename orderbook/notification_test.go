package orderbook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/orderbook"

	"github.com/republicprotocol/republic-go/order"
)

var _ = Describe("Notification", func() {

	FContext("when testing compatibility", func() {

		BeforeEach(func() {

		})

		Context("when traders are the same", func() {
			It("should return false", func() {
				notification := NotificationOpenOrder{
					OrderID:       [32]byte{},
					OrderFragment: order.Fragment{},
					Trader:        "trader1",
					Priority:      0,
				}
				Expect(notification.IsCompatible(order.Fragment{}, "trader1", 1)).To(BeFalse())
			})
		})

		Context("when both orders are fill-or-kill", func() {
			It("should return false", func() {
				notification := NotificationOpenOrder{
					OrderID: [32]byte{},
					OrderFragment: order.Fragment{
						OrderType: order.TypeLimitFOK,
					},
					Trader:   "trader1",
					Priority: 0,
				}
				orderFragment := order.Fragment{
					OrderID:   [32]byte{12},
					OrderType: order.TypeMidpointFOK,
				}
				Expect(notification.IsCompatible(orderFragment, "trader2", 1)).To(BeFalse())
			})
		})

		Context("when none of the orders are fill-or-kill", func() {
			It("should return true", func() {
				notification := NotificationOpenOrder{
					OrderID: [32]byte{},
					OrderFragment: order.Fragment{
						OrderType: order.TypeLimit,
					},
					Trader:   "trader1",
					Priority: 0,
				}
				orderFragment := order.Fragment{
					OrderID:   [32]byte{12},
					OrderType: order.TypeLimit,
				}
				Expect(notification.IsCompatible(orderFragment, "trader2", 1)).To(BeTrue())
			})
		})

		Context("when one of the orders are fill-or-kill", func() {
			It("should return true if priority of the FOK order is higher that the second order's priority", func() {
				notification := NotificationOpenOrder{
					OrderID: [32]byte{},
					OrderFragment: order.Fragment{
						OrderType: order.TypeLimitFOK,
					},
					Trader:   "trader1",
					Priority: 2,
				}
				orderFragment := order.Fragment{
					OrderID:   [32]byte{12},
					OrderType: order.TypeLimit,
				}
				Expect(notification.IsCompatible(orderFragment, "trader2", 1)).To(BeTrue())
			})

			It("should return false if priority of the FOK order is lower that the second order's priority", func() {
				notification := NotificationOpenOrder{
					OrderID: [32]byte{},
					OrderFragment: order.Fragment{
						OrderType: order.TypeFOK,
					},
					Trader:   "trader1",
					Priority: 0,
				}
				orderFragment := order.Fragment{
					OrderID:   [32]byte{12},
					OrderType: order.TypeMidpoint,
				}
				Expect(notification.IsCompatible(orderFragment, "trader2", 1)).To(BeFalse())
			})
		})
	})
})
