package middleware

import (
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

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

	rateEnabled bool
	rateValue   int
	burstValue  int
)

func init() {
	rateEnabled = os.Getenv("RATE_LIMIT_ENABLED") == "true"

	rateValue = 50
	if v, err := strconv.Atoi(os.Getenv("RATE_LIMIT")); err == nil {
		rateValue = v
	}

	burstValue = 500
	if v, err := strconv.Atoi(os.Getenv("RATE_BURST")); err == nil {
		burstValue = v
	}

	go cleanupRateLimiters()
}

// getLimiter returns a rate limiter for a given requester's IP address, creating a new one if necessary.
func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	lim, exists := rateLimiters[ip]
	if !exists {
		lim = &ipLimiter{
			limiter:  rate.NewLimiter(rate.Limit(rateValue), burstValue),
			lastSeen: time.Now(),
		}
		rateLimiters[ip] = lim
	} else {
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

// clientIP extracts the client IP from the request.
func clientIP(r *http.Request) string {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return forwarded
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// RateLimitMiddleware wraps an http.Handler with per-IP rate limiting.
func RateLimitMiddleware(next http.Handler) http.Handler {
	if !rateEnabled {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		limiter := getLimiter(clientIP(r))

		if !limiter.Allow() {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
