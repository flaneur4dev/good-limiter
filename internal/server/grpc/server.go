package grpcserver

import (
	"net"

	"google.golang.org/grpc"

	"github.com/flaneur4dev/good-limiter/internal/server/grpc/pb"
)

type Server struct {
	port string
	srv  *grpc.Server
}

func New(app app, port string) *Server {
	s := &Server{port: port}

	srv := grpc.NewServer()
	pb.RegisterRateLimiterServer(srv, newRateLimiterServer(app))

	s.srv = srv
	return s
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.port)
	if err != nil {
		return err
	}

	return s.srv.Serve(l)
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
