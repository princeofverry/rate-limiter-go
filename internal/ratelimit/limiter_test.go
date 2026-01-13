package ratelimit

import (
	"testing"
	"time"
)

func TestBucketAllowAndExhaust(t *testing.T) {
	b := newBucket(3, 0)

	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	b.now = func() time.Time { return base }
	b.last = base

	if !b.allow(1) || !b.allow(1) || !b.allow(1) {
		t.Fatalf("expected first 3 requests to succeed")
	}
	if b.allow(1) {
		t.Fatalf("expected 4th request to fail")
	}
}

func TestBucketRefill(t *testing.T) {
	b := newBucket(2, 60) 
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	cur := base
	b.now = func() time.Time { return cur }
	b.last = base
	b.tokens = 0
  
	// after 1 sec -> 1 token
	cur = base.Add(1 * time.Second)
	if !b.allow(1) {
	  t.Fatalf("expected to allow after refill")
	}
  }