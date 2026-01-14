package ratelimit

import (
	"sync"
	"time"
)

type nowfunc func() time.Time

type bucket struct {
	capacity int
	tokens   float64
	refillPS float64 // refill per second
	last     time.Time
	mu sync.Mutex
	now nowfunc
}

func newBucket(capacity, refillPerMinute int) *bucket {
	return &bucket {
		capacity: capacity,
		tokens: float64(capacity),
		refillPS: float64(refillPerMinute) / 60,
		last: time.Now(),
		now: time.Now,
	}
}

func (b *bucket) allow(n float64) bool {
	// just 1 process at a time
	b.mu.Lock()
	defer b.mu.Unlock()

	now := b.now()
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

type Status struct {
	Capacity int `json:"capacity"`
	Remaining float64 `json:"remaining"`
	RefillPM int	`json:"refill_per_minute"`
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

func (l *Limiter) Status(key string) (Status, bool) {
	l.mu.Lock()
	b, ok := l.buckets[key]
	l.mu.Unlock()
	if !ok {
		return Status{}, false
	}
	
	// sync refill before reporting
	b.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	b.tokens += elapsed * b.refillPS
	if b.tokens > float64(b.capacity) {
		b.tokens = float64(b.capacity)
	}
	
	st := Status {
		Capacity: b.capacity,
		Remaining: b.tokens,
		RefillPM: l.refillPM,
	}
	b.mu.Unlock()
	
	return st, true
}