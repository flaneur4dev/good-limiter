package ratelimiter

import (
	"sync"
	"time"
)

type leakyBucket struct {
	leakyCh chan struct{}
	doneCh  chan struct{}
	p       time.Duration
	mu      sync.Mutex
	t       time.Time
}

func newLeakyBucket(limit int, period time.Duration) *leakyBucket {
	lb := &leakyBucket{
		t:       time.Now(),
		p:       period,
		leakyCh: make(chan struct{}, limit),
		doneCh:  make(chan struct{}),
	}

	leakInterval := period.Nanoseconds() / int64(limit)
	go lb.start(time.Duration(leakInterval))

	return lb
}

func (lb *leakyBucket) start(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-lb.doneCh:
			return

		case <-timer.C:
			select {
			case <-lb.leakyCh:
			default:
			}
		}
	}
}

func (lb *leakyBucket) Allow() bool {
	select {
	case lb.leakyCh <- struct{}{}:
		lb.mu.Lock()
		lb.t = time.Now()
		lb.mu.Unlock()
		return true
	default:
		return false
	}
}

func (lb *leakyBucket) LastUse() time.Time {
	return lb.t
}

func (lb *leakyBucket) Period() time.Duration {
	return lb.p
}

func (lb *leakyBucket) Stop() {
	close(lb.doneCh)
}
