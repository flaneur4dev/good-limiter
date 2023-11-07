package ratelimiter

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cs "github.com/flaneur4dev/good-limiter/internal/contracts"
	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
	"github.com/flaneur4dev/good-limiter/internal/rate-limiter/mocks"
	bucketMock "github.com/flaneur4dev/good-limiter/mocks"
)

func TestRateLimiter(t *testing.T) {
	type args struct {
		ctx      context.Context
		login    string
		password string
		ip       string
		ipNet    string
		netList  string
	}

	type mocksBehavior func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args)

	tests := [...]struct {
		name     string
		args     args
		bfunc    mocksBehavior
		expected bool
	}{
		{
			name: "ipNet from black list",
			args: args{
				ctx:      context.Background(),
				login:    "login_1",
				password: "password_1",
				ip:       "192.168.0.100",
				ipNet:    "192.168.0.0/24",
				netList:  "black",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().
					Net(args.ctx, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, s string) (string, error) {
						if s == args.ipNet {
							return args.netList, nil
						}
						return "", es.ErrNetNotFound
					}).
					Times(9)
			},
			expected: false,
		},
		{
			name: "ipNet from white list",
			args: args{
				ctx:      context.Background(),
				login:    "login_2",
				password: "password_2",
				ip:       "192.168.0.100",
				ipNet:    "192.168.0.96/27",
				netList:  "white",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().
					Net(args.ctx, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, s string) (string, error) {
						if s == args.ipNet {
							return args.netList, nil
						}
						return "", es.ErrNetNotFound
					}).
					Times(6)
			},
			expected: true,
		},
		{
			name: "false case due to login",
			args: args{
				ctx:      context.Background(),
				login:    "login_3",
				password: "password_3",
				ip:       "10.1.0.1",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, mock.AnythingOfType("string")).Return("", es.ErrNetNotFound).Times(33)

				bs.EXPECT().
					Bucket(args.ctx, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, s string) (cs.Bucket, bool) {
						res := true
						if s == args.login {
							res = false
						}

						b.EXPECT().Allow().Return(res).Once()
						return b, true
					}).
					Times(1)
			},
			expected: false,
		},
		{
			name: "false case due to password",
			args: args{
				ctx:      context.Background(),
				login:    "login_4",
				password: "password_4",
				ip:       "127.111.1.234",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, mock.AnythingOfType("string")).Return("", es.ErrNetNotFound).Times(33)

				bs.EXPECT().
					Bucket(args.ctx, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, s string) (cs.Bucket, bool) {
						res := true
						if s == args.password {
							res = false
						}

						b.EXPECT().Allow().Return(res).Once()
						return b, true
					}).
					Times(2)
			},
			expected: false,
		},
		{
			name: "false case due to ip",
			args: args{
				ctx:      context.Background(),
				login:    "login_5",
				password: "password_5",
				ip:       "127.111.1.234",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, mock.AnythingOfType("string")).Return("", es.ErrNetNotFound).Times(33)

				bs.EXPECT().
					Bucket(args.ctx, mock.AnythingOfType("string")).
					RunAndReturn(func(ctx context.Context, s string) (cs.Bucket, bool) {
						res := true
						if s == args.ip {
							res = false
						}

						b.EXPECT().Allow().Return(res).Once()
						return b, true
					}).
					Times(3)
			},
			expected: false,
		},
		{
			name: "true case",
			args: args{
				ctx:      context.Background(),
				login:    "login_6",
				password: "password_6",
				ip:       "255.1.0.0",
			},
			bfunc: func(b *bucketMock.Bucket, bs *mocks.BucketStorage, ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, mock.AnythingOfType("string")).Return("", es.ErrNetNotFound).Times(33)
				bs.EXPECT().Bucket(args.ctx, mock.AnythingOfType("string")).Return(b, true).Times(3)
				b.EXPECT().Allow().Return(true).Times(3)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := bucketMock.NewBucket(t)
			bs := mocks.NewBucketStorage(t)
			ns := mocks.NewNetStorage(t)

			tt.bfunc(b, bs, ns, tt.args)

			rl, err := New(slog.Default(), bs, ns, "10/m", "10/m", "10/m")
			require.NoError(t, err)

			res := rl.Allow(tt.args.ctx, tt.args.login, tt.args.password, tt.args.ip)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestAddNetToStorage(t *testing.T) {
	type args struct {
		ctx    context.Context
		subNet string
		list   string
	}

	type mockBehavior func(ns *mocks.NetStorage, args args)

	errSome := errors.New("some error")

	tests := [...]struct {
		name     string
		args     args
		bfunc    mockBehavior
		expected error
	}{
		{
			name: "add existing subnet with error",
			args: args{
				ctx:    context.Background(),
				subNet: "225.101.12.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("", errSome).Once()
			},
			expected: errSome,
		},
		{
			name: "add existing subnet 1",
			args: args{
				ctx:    context.Background(),
				subNet: "10.1.12.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return(args.list, nil).Once()
			},
			expected: es.ErrNetExist,
		},
		{
			name: "add existing subnet 2",
			args: args{
				ctx:    context.Background(),
				subNet: "101.101.101.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("white", nil).Once()
			},
			expected: es.ErrNetAnotherExist,
		},
		{
			name: "add new subnet with error",
			args: args{
				ctx:    context.Background(),
				subNet: "172.16.0.0/24",
				list:   "white",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("", es.ErrNetNotFound).Once()
				ns.EXPECT().AddNet(args.ctx, args.subNet, args.list).Return(errSome).Once()
			},
			expected: errSome,
		},
		{
			name: "add new subnet",
			args: args{
				ctx:    context.Background(),
				subNet: "192.168.0.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("", es.ErrNetNotFound).Once()
				ns.EXPECT().AddNet(args.ctx, args.subNet, args.list).Return(nil).Once()
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := mocks.NewNetStorage(t)

			tt.bfunc(ns, tt.args)

			rl, err := New(slog.Default(), nil, ns, "10/m", "10/m", "10/m")
			require.NoError(t, err)

			res := rl.AddNet(tt.args.ctx, tt.args.subNet, tt.args.list)
			require.Equal(t, tt.expected, res)
		})
	}
}

func TestDeleteNetFromStorage(t *testing.T) {
	type args struct {
		ctx    context.Context
		subNet string
		list   string
	}

	type mockBehavior func(ns *mocks.NetStorage, args args)

	errSome := errors.New("some error")

	tests := [...]struct {
		name     string
		args     args
		bfunc    mockBehavior
		expected error
	}{
		{
			name: "delete unknown subnet",
			args: args{
				ctx:    context.Background(),
				subNet: "225.101.12.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("", es.ErrNetNotFound).Once()
			},
			expected: es.ErrNetNotFound,
		},
		{
			name: "delete unknown subnet with error",
			args: args{
				ctx:    context.Background(),
				subNet: "10.1.12.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("", errSome).Once()
			},
			expected: errSome,
		},
		{
			name: "delete subnet from another list",
			args: args{
				ctx:    context.Background(),
				subNet: "101.101.101.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return("white", nil).Once()
			},
			expected: es.ErrNetAnotherExist,
		},
		{
			name: "delete subnet with error",
			args: args{
				ctx:    context.Background(),
				subNet: "172.16.0.0/24",
				list:   "white",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return(args.list, nil).Once()
				ns.EXPECT().DeleteNet(args.ctx, args.subNet).Return(errSome).Once()
			},
			expected: errSome,
		},
		{
			name: "delete subnet",
			args: args{
				ctx:    context.Background(),
				subNet: "192.168.0.0/24",
				list:   "black",
			},
			bfunc: func(ns *mocks.NetStorage, args args) {
				ns.EXPECT().Net(args.ctx, args.subNet).Return(args.list, nil).Once()
				ns.EXPECT().DeleteNet(args.ctx, args.subNet).Return(nil).Once()
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ns := mocks.NewNetStorage(t)

			tt.bfunc(ns, tt.args)

			rl, err := New(slog.Default(), nil, ns, "10/m", "10/m", "10/m")
			require.NoError(t, err)

			res := rl.DeleteNet(tt.args.ctx, tt.args.subNet, tt.args.list)
			require.Equal(t, tt.expected, res)
		})
	}
}
