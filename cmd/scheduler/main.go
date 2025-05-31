package main

import (
	"log"
	"net"

	"github.com/dheeraj-sn/distributed-orchestrator/internal/config"
	"github.com/dheeraj-sn/distributed-orchestrator/internal/scheduler"
	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

func main() {
	// Load nested config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set up structured logger using config
	var zapCfg zap.Config
	if cfg.Logging.Format == "json" {
		zapCfg = zap.NewProductionConfig()
	} else {
		zapCfg = zap.NewDevelopmentConfig()
	}
	if err := zapCfg.Level.UnmarshalText([]byte(cfg.Logging.Level)); err != nil {
		log.Fatalf("Invalid log level: %v", err)
	}
	zapCfg.EncoderConfig.TimeKey = "timestamp"
	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zapCfg.Build()
	if err != nil {
		log.Fatalf("Failed to build logger: %v", err)
	}
	defer logger.Sync()

	// Start TCP listener
	lis, err := net.Listen("tcp", cfg.Scheduler.Host)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// gRPC server setup
	grpcServer := grpc.NewServer()

	// Initialize scheduler logic and register gRPC service
	srv := scheduler.NewSchedulerServer(logger)
	pb.RegisterOrchestratorServer(grpcServer, srv)

	// Start job dispatcher loop in background
	srv.Dispatcher.Run()

	logger.Info("Scheduler listening", zap.String("addr", cfg.Scheduler.Host))

	// Serve gRPC
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
