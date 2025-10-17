package ratelimiter

import "time"

// Implementing the RateLimiter as an interface will allow us to easily swap out rate limiters in the future if needed
// as long as they conform to the Limiter interface. The Ad duration-specific methods can be implemented as empty
// interface methods for other rate limiter types

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
	GetCurrentAdDurationServed(ip string) int
	UpdateCurrentAdDurationServed(ip string, newAdDuration int)
}

type Config struct {
	AdDurationLimit int
	TimeFrame       time.Duration
}
