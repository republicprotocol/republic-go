package dispatch_test

import (
	"fmt"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/republic-go/dispatch"
)

var _ = Describe("Concurrency", func() {

	Context("when using cobegin", func() {

		It("should block until all goroutines terminate", func() {
			start := time.Now()
			CoBegin(func() {
				time.Sleep(1 * time.Second)
			})
			end := time.Now()
			Expect(end.Sub(start).Seconds() >= 1).Should(BeTrue())
		})

		It("should wait for multiple functions to end", func() {
			start := time.Now()
			CoBegin(func() {
				time.Sleep(250 * time.Millisecond)
			}, func() {
				time.Sleep(500 * time.Millisecond)
			}, func() {
				time.Sleep(750 * time.Millisecond)
			}, func() {
				time.Sleep(1 * time.Second)
			})
			end := time.Now()
			Expect(end.Sub(start).Seconds() >= 1).Should(BeTrue())
		})

	})

	Context("when using coforall loops", func() {

		It("should iterate over arrays", func() {
			num := int64(0)
			xs := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			CoForAll(xs, func(i int) {
				Expect(xs[i]).Should(Equal(i + 1))
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(10)))
		})

		It("should iterate over slices", func() {
			num := int64(0)
			xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			CoForAll(xs[:], func(i int) {
				Expect(xs[i]).Should(Equal(i + 1))
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(10)))
		})

		It("should over maps", func() {
			num := int64(0)
			xs := map[string]int{
				"1": 1,
				"2": 2,
				"3": 3,
			}
			CoForAll(xs, func(key string) {
				switch key {
				case "1":
					Expect(xs[key]).Should(Equal(1))
				case "2":
					Expect(xs[key]).Should(Equal(2))
				case "3":
					Expect(xs[key]).Should(Equal(3))
				default:
					panic(fmt.Sprintf("coforall error: invalid key found in map %v", key))
				}
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(3)))
		})

		It("should iterate over ints", func(done Done) {
			defer close(done)

			n := int64(0)
			xs := make(chan int, 10)
			CoForAll(10, func(i int) {
				xs <- (i + 1)
				atomic.AddInt64(&n, 1)
			})
			close(xs)

			for i := range xs {
				Expect(i).Should(BeNumerically(">=", 1))
				Expect(i).Should(BeNumerically("<=", 10))
			}
			Expect(n).Should(Equal(int64(10)))
		})

	})

	Context("when using forall loops", func() {

		It("should iterate over arrays", func() {
			num := int64(0)
			xs := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			ForAll(xs, func(i int) {
				Expect(xs[i]).Should(Equal(i + 1))
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(10)))
		})

		It("should iterate over slices", func() {
			num := int64(0)
			xs := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			ForAll(xs[:], func(i int) {
				Expect(xs[i]).Should(Equal(i + 1))
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(10)))
		})

		It("should over maps", func() {
			num := int64(0)
			xs := map[string]int{
				"1": 1,
				"2": 2,
				"3": 3,
			}
			ForAll(xs, func(key string) {
				switch key {
				case "1":
					Expect(xs[key]).Should(Equal(1))
				case "2":
					Expect(xs[key]).Should(Equal(2))
				case "3":
					Expect(xs[key]).Should(Equal(3))
				default:
					panic(fmt.Sprintf("coforall error: invalid key found in map %v", key))
				}
				atomic.AddInt64(&num, 1)
			})
			Expect(num).Should(Equal(int64(3)))
		})

		It("should iterate over ints", func(done Done) {
			defer close(done)

			n := int64(0)
			xs := make(chan int, 10)
			ForAll(10, func(i int) {
				xs <- (i + 1)
				atomic.AddInt64(&n, 1)
			})
			close(xs)

			for i := range xs {
				Expect(i).Should(BeNumerically(">=", 1))
				Expect(i).Should(BeNumerically("<=", 10))
			}
			Expect(n).Should(Equal(int64(10)))
		})

	})

})
