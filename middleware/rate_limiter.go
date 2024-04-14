package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	rateLimiters = make(map[string]*ipLimiter)
	mu           sync.Mutex
	// cleanupInterval defines how frequently to clean up stale entries.
	cleanupInterval = time.Minute * 5
	// ipTimeout defines how long to keep an IP in memory since last seen.
	ipTimeout = time.Hour
)

// init starts the cleanup goroutine.
func init() {
	log.Println("Starting cleanup goroutine")
	log.Println("cleanupInterval: ", cleanupInterval)
	log.Println("ipTimeout: ", ipTimeout)

	go cleanupRateLimiters()
}

// getLimiter returns a rate limiter for a given requester's IP address, creating a new one if necessary.
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	lim, exists := rateLimiters[ip]
	if !exists {
		lim = &ipLimiter{
			limiter:  rate.NewLimiter(50, 500),
			lastSeen: time.Now(),
		}
		rateLimiters[ip] = lim
	} else {
		// Update lastSeen time.
		lim.lastSeen = time.Now()
	}

	return lim.limiter
}

// cleanupRateLimiters removes stale entries from the rateLimiters map.
func cleanupRateLimiters() {
	for {
		time.Sleep(cleanupInterval)
		mu.Lock()
		for ip, lim := range rateLimiters {
			if time.Since(lim.lastSeen) > ipTimeout {
				delete(rateLimiters, ip)
			}
		}
		mu.Unlock()
	}
}

// RateLimitMiddleware checks if the requester has exceeded their rate limit.
func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getLimiter(c.ClientIP())

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			return
		}

		c.Next()
	}
}
