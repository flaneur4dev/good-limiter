package memory

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	"github.com/flaneur4dev/good-limiter/mocks"
)

const (
	interval = 3000 * time.Millisecond
	delay    = 10 * time.Millisecond
)

func TestMemStore(t *testing.T) {
	defer goleak.VerifyNone(t)

	tests := [...]struct {
		key     string
		lastUse time.Time
		period  time.Duration
		expired bool
	}{
		{
			key:     "login_1",
			lastUse: time.Now().Add(-time.Minute * 10),
			period:  time.Minute,
			expired: true,
		},
		{
			key:     "password_1234",
			lastUse: time.Now().Add(-time.Millisecond * 100),
			period:  time.Minute,
			expired: false,
		},
		{
			key:     "login_2",
			lastUse: time.Now().Add(-time.Minute),
			period:  time.Minute,
			expired: false,
		},
		{
			key:     "192.168.0.100",
			lastUse: time.Now().Add(-time.Minute * 2),
			period:  time.Minute,
			expired: true,
		},
		{
			key:     "password_42",
			lastUse: time.Now().Add(-time.Minute * 5),
			period:  time.Minute,
			expired: true,
		},
		{
			key:     "login_3",
			lastUse: time.Now().Add(-time.Millisecond * 420),
			period:  time.Minute,
			expired: false,
		},
		{
			key:     "10.1.0.1",
			lastUse: time.Now().Add(-time.Second * 45),
			period:  time.Minute,
			expired: false,
		},
		{
			key:     "login_4",
			lastUse: time.Now().Add(-time.Minute * 10),
			period:  time.Minute,
			expired: true,
		},
		{
			key:     "password_efkjrnwrj34",
			lastUse: time.Now().Add(-time.Minute * 10),
			period:  time.Minute,
			expired: true,
		},
		{
			key:     "255.1.0.0",
			lastUse: time.Now().Add(-time.Second * 5),
			period:  time.Minute,
			expired: false,
		},
	}

	ms := New(interval)
	defer ms.Close()

	t.Run("concurrently addition", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(len(tests))

		for _, tt := range tests {
			tt := tt
			go func() {
				defer wg.Done()

				b := mocks.NewBucket(t)
				b.EXPECT().LastUse().Return(tt.lastUse).Once()
				b.EXPECT().Period().Return(tt.period).Once()
				b.EXPECT().Stop().Maybe()

				err := ms.AddBucket(context.TODO(), tt.key, b)
				require.NoError(t, err)
			}()
		}

		wg.Wait()

		// дожидаемся пока сработает механизм очистки старых бакетов
		time.Sleep(interval + delay)
	})

	t.Run("concurrently receiving", func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(len(tests))

		for _, tt := range tests {
			tt := tt
			go func() {
				defer wg.Done()

				b, ok := ms.Bucket(context.TODO(), tt.key)
				if tt.expired {
					require.Nil(t, b)
					require.False(t, ok)
				} else {
					require.NotNil(t, b)
					require.True(t, ok)
				}
			}()
		}

		wg.Wait()
	})
}
