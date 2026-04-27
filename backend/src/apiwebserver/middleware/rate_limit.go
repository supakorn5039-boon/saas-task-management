package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit is a per-IP fixed-window limiter. Designed for the auth endpoints
// where we want to slow down credential stuffing without standing up Redis.
//
// Trade-offs (intentional):
//   - Memory grows with unique IPs in the window. The reaper goroutine prunes
//     entries older than 2× window so this stays bounded for any sane traffic.
//   - State is per-process — running multiple replicas weakens the limit
//     proportionally. For a single-binary deploy that's fine; behind a load
//     balancer with multiple replicas, swap to Redis (see ADR notes).
//
// 5 requests per minute is enough headroom for a real user fat-fingering a
// password while still meaningfully slowing a brute-force attempt.
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	type bucket struct {
		count       int
		windowStart time.Time
	}

	var (
		mu      sync.Mutex
		buckets = make(map[string]*bucket)
	)

	// Reap stale buckets to keep memory bounded under churn.
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		for now := range ticker.C {
			mu.Lock()
			for ip, b := range buckets {
				if now.Sub(b.windowStart) > 2*window {
					delete(buckets, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		b, ok := buckets[ip]
		if !ok || now.Sub(b.windowStart) >= window {
			buckets[ip] = &bucket{count: 1, windowStart: now}
			mu.Unlock()
			c.Next()
			return
		}
		b.count++
		exceeded := b.count > maxRequests
		mu.Unlock()

		if exceeded {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "too many requests, please try again later",
			})
			return
		}
		c.Next()
	}
}
