package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLeakyBucket(t *testing.T) {
	t.Parallel()

	lb := newLeakyBucket(3, time.Second)

	require.True(t, lb.Allow())
	require.True(t, lb.Allow())
	require.True(t, lb.Allow())
	require.False(t, lb.Allow())

	time.Sleep(350 * time.Millisecond)

	require.True(t, lb.Allow())
	require.False(t, lb.Allow())
}

func TestLeakyBucketWithGoroutines(t *testing.T) {
	t.Parallel()

	goroutinesNumber := 10
	lb := newLeakyBucket(goroutinesNumber, time.Second)

	wg := sync.WaitGroup{}
	wg.Add(goroutinesNumber)
	for i := 0; i < goroutinesNumber; i++ {
		go func() {
			defer wg.Done()
			require.True(t, lb.Allow())
		}()
	}

	wg.Wait()
	require.False(t, lb.Allow())
}
