package middleware

import (
	"net/http"
	"sync"
	"time"
)

type tokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	refillPS float64 // tokens per second
	lastSeen time.Time
}

func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.lastSeen).Seconds()
	b.lastSeen = now
	b.tokens = min(b.max, b.tokens+elapsed*b.refillPS)
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}

type rateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*tokenBucket
}

func newRateLimiter() *rateLimiter {
	rl := &rateLimiter{buckets: make(map[string]*tokenBucket)}
	go rl.evict()
	return rl
}

func (rl *rateLimiter) bucket(key string) *tokenBucket {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	b, ok := rl.buckets[key]
	if !ok {
		b = &tokenBucket{tokens: 60, max: 60, refillPS: 1, lastSeen: time.Now()}
		rl.buckets[key] = b
	}
	return b
}

func (rl *rateLimiter) evict() {
	t := time.NewTicker(5 * time.Minute)
	for range t.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-10 * time.Minute)
		for k, b := range rl.buckets {
			b.mu.Lock()
			if b.lastSeen.Before(cutoff) {
				delete(rl.buckets, k)
			}
			b.mu.Unlock()
		}
		rl.mu.Unlock()
	}
}

var globalLimiter = newRateLimiter()

// RateLimit applies a per-IP token bucket (60 req/min burst, 1 req/s refill).
func RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.RemoteAddr
		if !globalLimiter.bucket(key).allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
