package grpcserver

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
	"github.com/flaneur4dev/good-limiter/internal/server/grpc/pb"
)

type app interface {
	Allow(ctx context.Context, login, password, ip string) bool
	AddNet(ctx context.Context, subNet, list string) error
	DeleteNet(ctx context.Context, subNet, list string) error
	DropBucket(ctx context.Context, login, ip string) error
}

type rateLimiterServer struct {
	pb.UnimplementedRateLimiterServer
	app app
}

func newRateLimiterServer(app app) *rateLimiterServer {
	return &rateLimiterServer{app: app}
}

func (s *rateLimiterServer) Allow(ctx context.Context, req *pb.AllowRequest) (*pb.AllowResponse, error) {
	ok := s.app.Allow(ctx, req.GetLogin(), req.GetPassword(), req.GetIp())

	return &pb.AllowResponse{Ok: ok}, nil
}

func (s *rateLimiterServer) AddNet(ctx context.Context, req *pb.AddRequest) (*pb.Response, error) {
	err := s.app.AddNet(ctx, req.GetSubNet(), req.GetList())
	if err != nil {
		sc := codes.Internal
		if errors.Is(err, es.ErrNetExist) {
			sc = codes.InvalidArgument
		}

		return nil, status.Error(sc, err.Error())
	}

	return &pb.Response{Message: "added"}, nil
}

func (s *rateLimiterServer) DeleteNet(ctx context.Context, req *pb.DeleteRequest) (*pb.Response, error) {
	err := s.app.DeleteNet(ctx, req.GetSubNet(), req.GetList())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Response{Message: "deleted"}, nil
}

func (s *rateLimiterServer) DropBucket(ctx context.Context, req *pb.DropRequest) (*pb.Response, error) {
	err := s.app.DropBucket(ctx, req.GetLogin(), req.GetIp())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Response{Message: "droped"}, nil
}
