package redis

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"

	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
)

type NetStore struct {
	client *redis.Client
}

func New(a, p string) (*NetStore, error) {
	rc := redis.NewClient(&redis.Options{
		Addr:     a,
		Password: p,
		DB:       0,
	})

	sc := rc.Ping(context.TODO())
	if err := sc.Err(); err != nil {
		return nil, err
	}

	ns := &NetStore{rc}
	return ns, nil
}

func (ns *NetStore) Net(ctx context.Context, ip string) (string, error) {
	res, err := ns.client.Get(ctx, ip).Result()
	if errors.Is(err, redis.Nil) {
		return "", es.ErrNetNotFound
	}

	return res, err
}

func (ns *NetStore) AddNet(ctx context.Context, ip, list string) error {
	return ns.client.Set(ctx, ip, list, 0).Err()
}

func (ns *NetStore) DeleteNet(ctx context.Context, ip string) error {
	return ns.client.GetDel(ctx, ip).Err()
}

func (ns *NetStore) Close() error {
	return ns.client.Close()
}
