package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
	GetCurrentAdDurationServed(ip string) int
	UpdateCurrentAdDurationServed(ip string, newAdDuration int)
}

type Config struct {
	AdDurationLimit int
	TimeFrame       time.Duration
}
