package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrLockAlreadyAcquired = fmt.Errorf("lock already acquired")
var ErrLockNotAcquired = fmt.Errorf("lock not acquired")

type Lock interface {
	Lock(context.Context) error
	TryLock(context.Context) (bool, error)
	Unlock(context.Context) error
}

type lock struct {
	client      *redis.Client
	key         string
	acquired    bool
	acquirewait time.Duration
	renew       time.Duration
	cancel      context.CancelFunc
}

func NewLock(client *redis.Client, key string, acquirewait time.Duration, renew time.Duration) Lock {
	return &lock{
		client:      client,
		key:         key,
		acquirewait: acquirewait,
		renew:       renew,
		cancel:      func() {},
	}
}

func (l *lock) Lock(ctx context.Context) error {
	if l.acquired {
		return ErrLockAlreadyAcquired
	}

	for {
		result, err := l.TryLock(ctx)
		if err != nil {
			return err
		}

		if result {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(l.acquirewait):
		}
	}
}

func (l *lock) TryLock(ctx context.Context) (bool, error) {
	if l.acquired {
		return false, ErrLockAlreadyAcquired
	}

	result, err := l.client.SetNX(ctx, l.key, "", l.renew*3/2).Result()
	if err != nil {
		return false, err
	}

	if result {
		l.acquired = true

		renewCtx, cancel := context.WithCancel(context.Background())
		l.cancel = cancel
		go l.renewLock(renewCtx)

		return true, nil
	}

	return false, nil
}

func (l *lock) Unlock(ctx context.Context) error {
	defer l.cancel()

	if !l.acquired {
		return ErrLockNotAcquired
	}

	_, err := l.client.Del(ctx, l.key).Result()
	if err != nil {
		return err
	}

	l.acquired = false

	return nil
}

func (l *lock) renewLock(ctx context.Context) {
	ticker := time.NewTicker(l.renew)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			l.client.Set(ctx, l.key, "", l.renew*3/2)
		}
	}
}
