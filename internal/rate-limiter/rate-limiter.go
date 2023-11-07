package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	cs "github.com/flaneur4dev/good-limiter/internal/contracts"
	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
)

type logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
}

type bucketStorage interface {
	Bucket(ctx context.Context, key string) (cs.Bucket, bool)
	AddBucket(ctx context.Context, key string, b cs.Bucket) error
	DeleteBucket(ctx context.Context, key string) error
}

type netStorage interface {
	Net(ctx context.Context, ip string) (string, error)
	AddNet(ctx context.Context, ip, list string) error
	DeleteNet(ctx context.Context, ip string) error
}

type constraint struct {
	limit  int
	period time.Duration
}

type RateLimiter struct {
	log         logger
	bStore      bucketStorage
	nStore      netStorage
	loginCnt    constraint
	passwordCnt constraint
	ipCnt       constraint
}

const maxPrefixBits = 32

func New(log logger, bStore bucketStorage, nStore netStorage, a, b, c string) (*RateLimiter, error) {
	lc, err := parseConstraint(a)
	if err != nil {
		return nil, err
	}

	pc, err := parseConstraint(b)
	if err != nil {
		return nil, err
	}

	ic, err := parseConstraint(c)
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		log:         log,
		bStore:      bStore,
		nStore:      nStore,
		loginCnt:    lc,
		passwordCnt: pc,
		ipCnt:       ic,
	}, nil
}

func (rl *RateLimiter) Allow(ctx context.Context, login, password, ip string) bool {
	l := rl.checkNetList(ctx, ip)
	switch l {
	case "white":
		return true
	case "black":
		return false
	}

	if !rl.checkLogin(ctx, login) || !rl.checkPassword(ctx, password) || !rl.checkIP(ctx, ip) {
		return false
	}

	return true
}

func (rl *RateLimiter) AddNet(ctx context.Context, subNet, list string) error {
	l, err := rl.nStore.Net(ctx, subNet)
	switch {
	case errors.Is(err, es.ErrNetNotFound):
	case err != nil:
		rl.log.Error("net storage error: " + err.Error())
		return err
	case l == list:
		return es.ErrNetExist
	case l != list:
		return es.ErrNetAnotherExist
	}

	err = rl.nStore.AddNet(ctx, subNet, list)
	if err != nil {
		rl.log.Error("failed to add new subnet: " + err.Error())
		return err
	}

	return nil
}

func (rl *RateLimiter) DeleteNet(ctx context.Context, subNet, list string) error {
	l, err := rl.nStore.Net(ctx, subNet)
	switch {
	case errors.Is(err, es.ErrNetNotFound):
		return err
	case err != nil:
		rl.log.Error("net storage error: " + err.Error())
		return err
	case l != list:
		return es.ErrNetAnotherExist
	}

	err = rl.nStore.DeleteNet(ctx, subNet)
	if err != nil {
		rl.log.Error("failed to delete subnet: " + err.Error())
		return err
	}

	return nil
}

func (rl *RateLimiter) DropBucket(ctx context.Context, login, ip string) error {
	err := rl.bStore.DeleteBucket(ctx, login)
	if err != nil {
		rl.log.Error("failed to drop login bucket: " + err.Error())
		return err
	}

	err = rl.bStore.DeleteBucket(ctx, ip)
	if err != nil {
		rl.log.Error("failed to drop ip bucket: " + err.Error())
		return err
	}

	return nil
}

func (rl *RateLimiter) checkNetList(ctx context.Context, ip string) string {
	for i := maxPrefixBits; i >= 0; i-- {
		_, ipNet, err := net.ParseCIDR(fmt.Sprintf("%s/%d", ip, i))
		if err != nil {
			rl.log.Error("parse CIDR error: " + err.Error())
			continue
		}

		l, err := rl.nStore.Net(ctx, ipNet.String())
		switch {
		case errors.Is(err, es.ErrNetNotFound):
			continue
		case err != nil:
			rl.log.Error("net storage error: " + err.Error())
			continue
		}

		return l
	}

	return ""
}

func (rl *RateLimiter) checkLogin(ctx context.Context, login string) bool {
	return rl.checkOrSet(ctx, login, rl.loginCnt.limit, rl.loginCnt.period)
}

func (rl *RateLimiter) checkPassword(ctx context.Context, password string) bool {
	return rl.checkOrSet(ctx, password, rl.passwordCnt.limit, rl.passwordCnt.period)
}

func (rl *RateLimiter) checkIP(ctx context.Context, ip string) bool {
	return rl.checkOrSet(ctx, ip, rl.ipCnt.limit, rl.ipCnt.period)
}

func (rl *RateLimiter) checkOrSet(ctx context.Context, key string, limit int, period time.Duration) bool {
	b, ok := rl.bStore.Bucket(ctx, key)
	if !ok {
		rl.bStore.AddBucket(ctx, key, newLeakyBucket(limit, period))
		return true
	}

	return b.Allow()
}

func parseConstraint(s string) (constraint, error) {
	var c constraint

	l, err := strconv.Atoi(s[:len(s)-2])
	if err != nil {
		return c, err
	}

	p, err := time.ParseDuration("1" + s[len(s)-1:])
	if err != nil {
		return c, err
	}

	c.limit, c.period = l, p
	return c, nil
}
