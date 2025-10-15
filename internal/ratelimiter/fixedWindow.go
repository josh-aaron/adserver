package ratelimiter

import (
	"log"
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex
	clients         map[string]int
	adDurationLimit int
	window          time.Duration
}

func NewFixedWindowLimiter(adDurationLimit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients:         make(map[string]int),
		adDurationLimit: adDurationLimit,
		window:          window,
	}
}

func (rl *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	log.Printf("FixedWindowRateLimiter Allow() for ip %v", ip)

	// Use locks to prevent race conditions
	rl.RLock()
	currentAdDurationServed, exists := rl.clients[ip]
	rl.RUnlock()

	// If we haven't seen the IP before, OR if we have but the IP hasn't reached its ad duration limit,
	// then allow the request to proceed to the getVastHandler. If it's a new IP, start a
	// Go routine to manage resetting the duration served in one hour.
	if !exists || currentAdDurationServed < rl.adDurationLimit {
		rl.Lock()
		if !exists {
			go rl.resetCount(ip)
		}
		rl.Unlock()
		return true, 0
	}

	return false, rl.window
}

// After an IP address is
func (rl *FixedWindowRateLimiter) resetCount(ip string) {
	log.Printf("FixedWindowRateLimiter resetCount() for ip %v", ip)
	time.Sleep(rl.window)
	log.Printf("window: %v", rl.window)
	rl.Lock()
	delete(rl.clients, ip)
	rl.Unlock()
}

func (rl *FixedWindowRateLimiter) GetCurrentAdDurationServed(ip string) int {
	log.Printf("GetCurrentAdDurationServed for ip %v", ip)
	rl.RLock()
	currentAdDurationServed := rl.clients[ip]
	rl.RUnlock()
	return currentAdDurationServed
}

// Update the amount of ad duraiton we've served to an IP address with the duraiton from the latest returned VAST
func (rl *FixedWindowRateLimiter) UpdateCurrentAdDurationServed(ip string, newAdDuration int) {
	log.Printf("UpdateCurrentAdDurationServed for ip %v", ip)
	rl.Lock()
	currentAdDurationServed := rl.clients[ip]
	rl.clients[ip] = currentAdDurationServed + newAdDuration
	rl.Unlock()
	log.Printf("UpdateCurrentAdDurationServed new currentAdDurationServed: %v", rl.clients[ip])
}
