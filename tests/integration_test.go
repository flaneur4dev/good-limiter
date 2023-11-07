package integration_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	es "github.com/flaneur4dev/good-limiter/internal/mistakes"
	"github.com/flaneur4dev/good-limiter/internal/server/grpc/pb"
)

type RateLimiterSuite struct {
	suite.Suite

	db             *redis.Client
	conn           *grpc.ClientConn
	grpcClient     pb.RateLimiterClient
	requestTimeout time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
}

func (s *RateLimiterSuite) SetupSuite() {
	s.setupDelay()
	s.setupDB()
	s.setupGRPC()
	s.setupTimeout()
}

func (s *RateLimiterSuite) TearDownSuite() {
	err := s.db.Close()
	s.Require().NoError(err)

	err = s.conn.Close()
	s.Require().NoError(err)
}

func (s *RateLimiterSuite) SetupTest() {
	s.ctx, s.cancel = context.WithTimeout(context.Background(), s.requestTimeout)
}

func (s *RateLimiterSuite) TearDownTest() {
	s.cancel()
}

func (s *RateLimiterSuite) TestRateLimiter() {
	s.T().Log(`Конфигурации для тестирования:
	логин - 10 запросов в минуту;
	пароль - 100 запросов в минуту;
	ip - 1000 запросов в минуту`)

	s.Run("case with adding subnet in white list", func() {
		alReq := &pb.AllowRequest{Login: "login_1", Password: "password_1", Ip: "16.168.0.100"}
		adReq := &pb.AddRequest{SubNet: "16.168.0.0/24", List: "white"}
		iters := 10

		for i := 0; i <= iters; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().True(alRes.GetOk())
		}

		for i := 0; i <= iters*2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().False(alRes.GetOk())
		}

		time.Sleep(6 * time.Second)
		alRes, err := s.grpcClient.Allow(context.Background(), alReq)
		s.Require().NoError(err)
		s.Require().True(alRes.GetOk())

		for i := 0; i <= iters*2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().False(alRes.GetOk())
		}

		adRes, err := s.grpcClient.AddNet(context.Background(), adReq)
		s.Require().NoError(err)
		s.Require().Equal("added", adRes.GetMessage())

		for i := 0; i <= iters*2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().True(alRes.GetOk())
		}
	})

	s.Run("case with adding subnet in black list", func() {
		alReq := &pb.AllowRequest{Login: "login_2", Password: "password_2", Ip: "152.218.0.100"}
		adReq := &pb.AddRequest{SubNet: "152.218.0.0/24", List: "black"}
		iters := 10

		for i := 0; i <= iters/2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().True(alRes.GetOk())
		}

		adRes, err := s.grpcClient.AddNet(context.Background(), adReq)
		s.Require().NoError(err)
		s.Require().Equal("added", adRes.GetMessage())

		for i := 0; i <= iters*2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().False(alRes.GetOk())
		}
	})

	s.Run("case with droping bucket", func() {
		alReq := &pb.AllowRequest{Login: "login_3", Password: "password_3", Ip: "255.10.1.42"}
		drReq := &pb.DropRequest{Login: "login_3", Ip: "255.10.1.42"}
		iters := 10

		for i := 0; i <= iters; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().True(alRes.GetOk())
		}

		for i := 0; i <= iters/2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().False(alRes.GetOk())
		}

		adRes, err := s.grpcClient.DropBucket(context.Background(), drReq)
		s.Require().NoError(err)
		s.Require().Equal("droped", adRes.GetMessage())

		for i := 0; i <= iters; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().True(alRes.GetOk())
		}

		for i := 0; i <= iters*2; i++ {
			alRes, err := s.grpcClient.Allow(context.Background(), alReq)
			s.Require().NoError(err)
			s.Require().False(alRes.GetOk())
		}
	})
}

func (s *RateLimiterSuite) TestAddNet() {
	tests := [...]struct {
		name      string
		subNet    string
		list      string
		respError error
	}{
		{
			name:      "add new subnet",
			subNet:    "225.101.12.0/24",
			list:      "black",
			respError: nil,
		},
		{
			name:      "add existing subnet 1",
			subNet:    "225.101.12.0/24",
			list:      "black",
			respError: es.ErrNetExist,
		},
		{
			name:      "add existing subnet 2",
			subNet:    "225.101.12.0/24",
			list:      "white",
			respError: es.ErrNetAnotherExist,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.grpcClient.AddNet(s.ctx, &pb.AddRequest{SubNet: tt.subNet, List: tt.list})
			if err != nil {
				s.Require().Contains(err.Error(), tt.respError.Error())
			}

			if res != nil {
				s.Require().Equal("added", res.GetMessage())

				// проверка, что приложение произвело изменения именно в нужной базе данных, а не где-либо ещё
				l, err := s.db.Get(context.Background(), tt.subNet).Result()
				s.Require().NoError(err)
				s.Require().Equal(tt.list, l)
			}
		})
	}
}

func (s *RateLimiterSuite) TestDeleteNet() {
	tests := [...]struct {
		name      string
		subNet    string
		list      string
		respError error
	}{
		{
			name:      "delete unknown subnet",
			subNet:    "155.101.101.0/24",
			list:      "black",
			respError: es.ErrNetNotFound,
		},
		{
			name:      "delete existing subnet 1",
			subNet:    "192.168.0.0/24",
			list:      "black",
			respError: es.ErrNetAnotherExist,
		},
		{
			name:      "delete existing subnet 2",
			subNet:    "192.168.0.0/24",
			list:      "white",
			respError: nil,
		},
	}

	err := s.db.Set(context.Background(), "192.168.0.0/24", "white", 0).Err()
	s.Require().NoError(err)

	for _, tt := range tests {
		s.Run(tt.name, func() {
			res, err := s.grpcClient.DeleteNet(s.ctx, &pb.DeleteRequest{SubNet: tt.subNet, List: tt.list})
			if err != nil {
				s.Require().Contains(err.Error(), tt.respError.Error())
			}

			if res != nil {
				s.Require().Equal("deleted", res.GetMessage())

				// проверка, что приложение произвело изменения именно в нужной базе данных, а не где-либо ещё
				_, err := s.db.Get(context.Background(), tt.subNet).Result()
				s.Require().Equal(redis.Nil, err)
			}
		})
	}
}

func (s *RateLimiterSuite) setupDelay() {
	delay := os.Getenv("TEST_DELAY")
	if delay == "" {
		return
	}

	d, err := time.ParseDuration(delay)
	s.Require().NoError(err)

	s.T().Logf("wait %s for service availability...", delay)
	time.Sleep(d)
}

func (s *RateLimiterSuite) setupDB() {
	a := os.Getenv("TEST_DB")
	if a == "" {
		a = "localhost:6379"
	}

	rc := redis.NewClient(&redis.Options{
		Addr:     a,
		Password: "",
		DB:       0,
	})

	sc := rc.Ping(context.TODO())
	err := sc.Err()
	s.Require().NoError(err)

	s.db = rc
}

func (s *RateLimiterSuite) setupGRPC() {
	host := os.Getenv("TEST_GRPC")
	if host == "" {
		host = ":50051"
	}

	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.conn = conn
	s.grpcClient = pb.NewRateLimiterClient(conn)
}

func (s *RateLimiterSuite) setupTimeout() {
	d := os.Getenv("TEST_REQUEST_TIMEOUT")
	if d == "" {
		d = "100ms"
	}

	t, err := time.ParseDuration(d)
	s.Require().NoError(err)

	s.requestTimeout = t
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterSuite))
}
