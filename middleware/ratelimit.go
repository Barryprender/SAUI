package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type tokenBucket struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	refillPS float64
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

// NewRateLimit returns a per-IP token bucket rate limiter (60 req burst, 1 req/s refill).
// trustedProxy: if non-empty, X-Real-IP / X-Forwarded-For headers are trusted only when
// the direct TCP peer IP matches this value.
func NewRateLimit(trustedProxy string, logger *slog.Logger) func(http.Handler) http.Handler {
	rl := newRateLimiter()
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := clientIP(r, trustedProxy)
			if !rl.bucket(key).allow() {
				logger.Warn("rate limit exceeded", "ip", key, "path", r.URL.Path)
				http.Error(w, "too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// clientIP extracts the real client IP, honouring proxy headers only when the
// direct TCP peer matches trustedProxy.
func clientIP(r *http.Request, trustedProxy string) string {
	if trustedProxy != "" {
		peerHost, _, _ := net.SplitHostPort(r.RemoteAddr)
		if peerHost == trustedProxy {
			if ip := r.Header.Get("X-Real-IP"); ip != "" {
				return strings.TrimSpace(ip)
			}
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				if i := strings.Index(xff, ","); i != -1 {
					return strings.TrimSpace(xff[:i])
				}
				return strings.TrimSpace(xff)
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
