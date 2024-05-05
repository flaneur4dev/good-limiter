package memory

import (
	"context"
	"sync"
	"time"

	cs "github.com/flaneur4dev/good-limiter/internal/contracts"
	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
)

type BucketStore struct {
	done    chan struct{}
	mu      sync.RWMutex
	buckets map[string]cs.Bucket
}

func New(interval time.Duration) *BucketStore {
	ms := &BucketStore{
		done:    make(chan struct{}),
		buckets: make(map[string]cs.Bucket),
	}

	go ms.cleanUp(interval)

	return ms
}

func (ms *BucketStore) Bucket(ctx context.Context, key string) (cs.Bucket, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	default:
	}

	ms.mu.RLock()
	defer ms.mu.RUnlock()

	b, ok := ms.buckets[key]
	return b, ok
}

func (ms *BucketStore) AddBucket(ctx context.Context, key string, b cs.Bucket) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.buckets[key] = b
	return nil
}

func (ms *BucketStore) DeleteBucket(ctx context.Context, key string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	b, ok := ms.buckets[key]
	if !ok {
		return es.ErrBucketNotFound
	}

	b.Stop()
	delete(ms.buckets, key)

	return nil
}

func (ms *BucketStore) Close() {
	close(ms.done)

	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, b := range ms.buckets {
		b.Stop()
	}

	ms.buckets = nil
}

func (ms *BucketStore) cleanUp(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-ms.done:
			return

		case <-timer.C:
			ms.mu.Lock()

			for k, b := range ms.buckets {
				if time.Since(b.LastUse()) >= b.Period()*2 {
					b.Stop()
					delete(ms.buckets, k)
				}
			}

			ms.mu.Unlock()
		}
	}
}
