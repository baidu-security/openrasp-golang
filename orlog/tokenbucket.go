package orlog

import "time"

type TokenBucket struct {
	refillInterval     time.Duration
	capacity           uint64
	currentTokenAmount uint64
	lastConsumedTime   time.Time
}

func NewTokenBucket(capacity uint64, refillInterval time.Duration) *TokenBucket {
	tb := &TokenBucket{
		refillInterval:     refillInterval,
		capacity:           capacity,
		currentTokenAmount: capacity,
		lastConsumedTime:   time.Now(),
	}
	return tb
}

func (tb *TokenBucket) Consume() bool {
	tb.refill()
	isEmpty := tb.currentTokenAmount <= 0
	if !isEmpty {
		tb.currentTokenAmount--
		tb.lastConsumedTime = time.Now()
	}
	return isEmpty
}

func (tb *TokenBucket) refill() {
	current := time.Now()
	elapsedTimeFromLastConsumed := current.Sub(tb.lastConsumedTime)
	if int64(elapsedTimeFromLastConsumed) > int64(tb.refillInterval) {
		tb.currentTokenAmount = tb.capacity
	}
}
