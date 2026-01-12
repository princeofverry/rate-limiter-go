package ratelimit

import (
	"sync"
	"time"
)

type bucket struct {
	capacity int
	tokens   float64
	refillPS float64 // refill per second
	last     time.Time
	mu sync.Mutex
}

func newBucket(capacity, refillPerMinute int) *bucket {
	return &bucket {
		capacity: capacity,
		tokens: float64(capacity),
		refillPS: float64(refillPerMinute) / 60,
		last: time.Now(),
	}
}

func (b *bucket) allow(n float64) bool {
	// just 1 process at a time
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	// counting the time
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	// refill
	b.tokens += elapsed * b.refillPS
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}

	// accessing the bucket
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}

type Limiter struct {
	capacity int
	refillPM int

	mu sync.Mutex
	buckets map[string]*bucket
}

func New(capacity, refillPerMinute int) *Limiter {
	return &Limiter {
		capacity: capacity,
		refillPM: refillPerMinute,
		buckets: make(map[string]*bucket),
	}
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	b, ok := l.buckets[key]
	if !ok {
		b = newBucket(l.capacity, l.refillPM)
		l.buckets[key] = b
	}
	l.mu.Unlock()

	return b.allow(1)
}