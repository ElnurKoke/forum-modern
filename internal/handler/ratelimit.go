package handler

import (
	"net"
	"net/http"
	"sync"
)

type RateLimiter struct {
	mutex         sync.Mutex
	requestsCount map[string]int
	next          http.Handler
}

func NewRateLimiter(next http.Handler) *RateLimiter {
	return &RateLimiter{
		requestsCount: make(map[string]int),
		next:          next,
	}
}

func (rl *RateLimiter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	rl.requestsCount[clientIP]++
	if rl.requestsCount[clientIP] > 10 {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	rl.next.ServeHTTP(w, r)
}
