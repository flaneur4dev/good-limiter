package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/flaneur4dev/good-limiter/internal/config"
	"github.com/flaneur4dev/good-limiter/internal/logger"
	limiter "github.com/flaneur4dev/good-limiter/internal/rate-limiter"
	grpcserver "github.com/flaneur4dev/good-limiter/internal/server/grpc"
	"github.com/flaneur4dev/good-limiter/internal/storage/memory"
	"github.com/flaneur4dev/good-limiter/internal/storage/redis"
)

var configFile = flag.String("config", "./configs/limiter.yaml", "path to configuration file")

func main() {
	flag.Parse()

	switch flag.Arg(0) {
	case "help":
		printHelp()
		return
	case "version":
		printVersion()
		return
	}

	cfg, err := config.New(*configFile)
	if err != nil {
		fmt.Printf("invalid config: %s", err)
		os.Exit(1)
	}

	logg := logger.New(cfg.Logger.Level)

	bs := memory.New(cfg.Storage.Memory.CleanUpInterval)
	defer bs.Close()

	ns, err := redis.New(cfg.Storage.DB.Host, cfg.Storage.DB.Password)
	if err != nil {
		logg.Error("failed connection to storage: " + err.Error())
		bs.Close()
		os.Exit(1) //nolint:gocritic
	}
	defer ns.Close()

	logg.Info("database connected...")

	rl, err := limiter.New(logg, bs, ns, cfg.Constraints.Login, cfg.Constraints.Password, cfg.Constraints.IP)
	if err != nil {
		logg.Error("failed to create rate limiter: " + err.Error())
		bs.Close()
		ns.Close()
		os.Exit(1) //nolint:gocritic
	}

	gsrv := grpcserver.New(rl, cfg.Server.GRPC.Port)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		<-ctx.Done()
		gsrv.Stop()
	}()

	logg.Info("grpc server is running...")
	if err := gsrv.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
		cancel()
		bs.Close()
		ns.Close()
		os.Exit(1) //nolint:gocritic
	}

	wg.Wait()
	logg.Info("grpc server is stopped")
}
