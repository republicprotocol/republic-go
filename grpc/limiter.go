package grpc

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

// RateLimiter wraps the time/rate.Limiter. It first limits the request
// according to the ip address.Then do a global limiting on all the
// request received.
type RateLimiter struct {
	mu     *sync.Mutex
	limit  rate.Limit
	burst  int
	global *rate.Limiter
	local  map[string]*rate.Limiter
}

// NewRateLimiter returns a new RateLimiter
func NewRateLimiter(limiter *rate.Limiter, limit float64, burst int) *RateLimiter {
	return &RateLimiter{
		mu:     new(sync.Mutex),
		limit:  rate.Limit(limit),
		burst:  burst,
		global: limiter,
		local:  map[string]*rate.Limiter{},
	}
}

func (limiter *RateLimiter) Allow(addr string) bool {
	limiter.mu.Lock()
	if _, ok := limiter.local[addr]; !ok {
		limiter.local[addr] = rate.NewLimiter(limiter.limit, limiter.burst)
	}
	addrLimiter := limiter.local[addr]
	limiter.mu.Unlock()

	if !addrLimiter.Allow() {
		return false
	}

	return limiter.global.Allow()
}

func (limiter *RateLimiter) Wait(ctx context.Context, addr string) error {
	limiter.mu.Lock()
	if _, ok := limiter.local[addr]; !ok {
		limiter.local[addr] = rate.NewLimiter(5, 20)
	}
	addrLimiter := limiter.local[addr]
	limiter.mu.Unlock()

	if err := addrLimiter.Wait(ctx); err != nil {
		return err
	}

	return limiter.global.Wait(ctx)
}

func (limiter *RateLimiter) Reserve(addr string) *rate.Reservation {
	limiter.mu.Lock()
	if _, ok := limiter.local[addr]; !ok {
		limiter.local[addr] = rate.NewLimiter(5, 20)
	}
	addrLimiter := limiter.local[addr]
	limiter.mu.Unlock()

	if reservation := addrLimiter.Reserve(); reservation != nil {
		return reservation
	}

	return limiter.global.Reserve()
}

func (limiter *RateLimiter) SetLimit(limit float64) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	limiter.limit = rate.Limit(limit)
}

func (limiter *RateLimiter) SetBurst(burst int) {
	limiter.mu.Lock()
	defer limiter.mu.Unlock()

	limiter.burst = burst
}
