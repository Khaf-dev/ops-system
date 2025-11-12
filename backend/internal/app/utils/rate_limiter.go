package utils

import (
	"sync"
	"time"
)

type attemptInfo struct {
	Count     int
	FirstAt   time.Time
	BlockedAt time.Time
}

type LoginLimiter struct {
	mu        sync.Mutex
	attempts  map[string]*attemptInfo
	limit     int
	window    time.Duration
	blockTime time.Duration
}

func NewLoginLimiter(limit int, window, blockTime time.Duration) *LoginLimiter {
	return &LoginLimiter{
		attempts:  make(map[string]*attemptInfo),
		limit:     limit,
		window:    window,
		blockTime: blockTime,
	}
}

func (l *LoginLimiter) TooManyAttempts(key string) (blocked bool, remain time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	a, ok := l.attempts[key]
	if !ok {
		return false, 0
	}
	now := time.Now()
	if !a.BlockedAt.IsZero() && now.Before(a.BlockedAt.Add(l.blockTime)) {
		return true, a.BlockedAt.Add(l.blockTime).Sub(now)
	}
	if now.Sub(a.FirstAt) > l.window {
		// reset window
		delete(l.attempts, key)
		return false, 0
	}
	if a.Count >= l.limit {
		// block
		a.BlockedAt = now
		return true, l.blockTime
	}
	return false, 0
}

func (l *LoginLimiter) RegisterFailure(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	a, ok := l.attempts[key]
	if !ok {
		l.attempts[key] = &attemptInfo{
			Count:   1,
			FirstAt: now,
		}
		return
	}
	// reset window if expired
	if now.Sub(a.FirstAt) > l.window {
		a.Count = 1
		a.FirstAt = now
		a.BlockedAt = time.Time{}
		return
	}
	a.Count++
}

func (l *LoginLimiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.attempts, key)
}

func (l *LoginLimiter) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			l.mu.Lock()
			now := time.Now()
			for key, a := range l.attempts {
				// fungsi untuk hapus entry jika windows nya udah expired
				if now.Sub(a.FirstAt) > l.window && now.After(a.BlockedAt.Add(l.blockTime)) {
					delete(l.attempts, key)
				}
			}
			l.mu.Unlock()
		}
	}()
}
