package utils

import (
	"time"
)

// RateLimiter wraps a ticker to limit actions
type RateLimiter struct {
	ticker *time.Ticker
}

// NewRateLimiter creates a new rate limiter that allows one event per interval
func NewRateLimiter(interval time.Duration) *RateLimiter {
	return &RateLimiter{
		ticker: time.NewTicker(interval),
	}
}

// Wait blocks until the next tick
func (rl *RateLimiter) Wait() {
	<-rl.ticker.C
}

// Stop stops the ticker
func (rl *RateLimiter) Stop() {
	rl.ticker.Stop()
}
