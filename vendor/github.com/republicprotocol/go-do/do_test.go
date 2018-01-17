package do_test

import (
	"errors"
	"sync/atomic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/republicprotocol/go-do"
)

var _ = Describe("Concurrency", func() {

	Context("when using an optional value", func() {
		It("should correctly execute calls to then", func() {
			i := 0
			Ok(true).Then(func(ok interface{}) Option {
				// This should happen.
				i++
				Ω(ok).Should(Equal(true))
				return Ok(true)
			}).Then(func(ok interface{}) Option {
				// This should happen.
				i++
				Ω(ok).Should(Equal(true))
				return Ok(true)
			}).Catch(func(err error) Option {
				// This should not happen.
				i++
				Ω(true).Should(Equal(false))
				return Err(nil)
			}).Then(func(ok interface{}) Option {
				// This should happen.
				i++
				Ω(ok).Should(Equal(true))
				return Ok(true)
			})
			Ω(i).Should(Equal(3))
		})

		It("should correctly execute calls to catch", func() {
			i := 0
			Err(errors.New("this is an error")).Catch(func(err error) Option {
				// This should happen.
				i++
				Ω(err).Should(HaveOccurred())
				return Err(err)
			}).Catch(func(err error) Option {
				// This should happen.
				i++
				Ω(err).Should(HaveOccurred())
				return Err(err)
			}).Then(func(ok interface{}) Option {
				// This should not happen.
				i++
				Ω(true).Should(Equal(false))
				return Ok(true)
			}).Catch(func(err error) Option {
				// This should happen.
				i++
				Ω(err).Should(HaveOccurred())
				return Err(err)
			})
			Ω(i).Should(Equal(3))
		})
	})

	Context("when using a for all loop", func() {
		It("should apply the function to all items", func() {
			xs := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
			ForAll(xs, func(i int) {
				xs[i] *= 2
			})
			for i := range xs {
				Ω(xs[i]).Should(Equal(i * 2))
			}
		})
	})

	Context("when using a process", func() {
		It("should write the return value to a channel", func() {
			ret := <-Process(func() Option {
				return Ok(1 + 2)
			})
			Ω(ret.Ok).Should(Equal(3))
		})

		It("should write the error to a channel", func() {
			ret := <-Process(func() Option {
				return Err(errors.New("this is an error"))
			})
			Ω(ret.Err).Should(HaveOccurred())
		})
	})

	Context("when using cobegin", func() {
		It("should wait until all functions have terminated", func() {
			i := int64(0)
			CoBegin(func() Option {
				atomic.AddInt64(&i, 1)
				return Ok(nil)
			}, func() Option {
				atomic.AddInt64(&i, 2)
				return Ok(nil)
			}, func() Option {
				atomic.AddInt64(&i, 3)
				return Ok(nil)
			})
			Ω(i).Should(Equal(int64(6)))
		})

		It("should collect all function return values", func() {
			results := CoBegin(func() Option {
				return Ok(0)
			}, func() Option {
				return Ok(1)
			}, func() Option {
				return Ok(2)
			})
			Ω(results[0].Ok).Should(Equal(0))
			Ω(results[1].Ok).Should(Equal(1))
			Ω(results[2].Ok).Should(Equal(2))
		})
	})

	Context("when using begin", func() {
		It("should wait until all functions have terminated", func() {
			i := int64(0)
			Begin(func() Option {
				atomic.AddInt64(&i, 1)
				return Ok(nil)
			}, func() Option {
				atomic.AddInt64(&i, 2)
				return Ok(nil)
			}, func() Option {
				atomic.AddInt64(&i, 3)
				return Ok(nil)
			})
			Ω(i).Should(Equal(int64(6)))
		})

		It("should collect all function return values", func() {
			results := Begin(func() Option {
				return Ok(0)
			}, func() Option {
				return Ok(1)
			}, func() Option {
				return Ok(2)
			})
			Ω(results[0].Ok).Should(Equal(0))
			Ω(results[1].Ok).Should(Equal(1))
			Ω(results[2].Ok).Should(Equal(2))
		})
	})

})
