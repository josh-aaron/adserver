package main

import (
	"testing"

	"github.com/josh-aaron/adserver/internal/model"
	"github.com/josh-aaron/adserver/internal/ratelimiter"
)

func newTestApplication(t *testing.T, config config) *application {
	t.Helper()

	mockRepo := model.NewMockRepo()

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.AdDurationLimit,
		config.rateLimiter.TimeFrame,
	)

	return &application{
		repository:  mockRepo,
		config:      config,
		rateLimiter: rateLimiter,
	}
}
