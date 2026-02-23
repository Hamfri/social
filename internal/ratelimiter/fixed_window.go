package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	sync.RWMutex // since we are embedding RWMutex this struct can only be used as a pointer and never as a value
	clients      map[string]int
	limit        int
	window       time.Duration
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowRateLimiter {
	return &FixedWindowRateLimiter{
		clients: make(map[string]int),
		limit:   limit,
		window:  window,
	}
}

func (l *FixedWindowRateLimiter) Allow(ip string) (bool, time.Duration) {
	// l.Lock()
	// count, exists := l.clients[ip]
	// l.Unlock()

	// if !exists || count < l.limit {
	// 	l.Lock()
	// 	if !exists {
	// 		go l.resetCount(ip)
	// 	}

	// 	l.clients[ip]++
	// 	l.Unlock()
	// 	return true, 0
	// }

	// return false, l.window

	l.Lock()
	defer l.Unlock()

	count := l.clients[ip]
	if count == 0 {
		go l.resetCount(ip)
	}

	if count < l.limit {
		l.clients[ip]++
		return true, 0
	}

	return false, l.window
}

func (l *FixedWindowRateLimiter) resetCount(ip string) {
	time.Sleep(l.window)
	l.Lock()
	delete(l.clients, ip)
	l.Unlock()
}
