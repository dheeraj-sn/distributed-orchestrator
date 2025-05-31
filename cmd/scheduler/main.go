package main

import (
	"log"
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/dheeraj-sn/distributed-orchestrator/internal/scheduler"
	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"
	"github.com/spf13/viper"
)

func main() {
	loadConfig()

	// Logger setup
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Start TCP listener
	schedulerHost := viper.GetString("scheduler.host")
	lis, err := net.Listen("tcp", schedulerHost)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	// gRPC server
	grpcServer := grpc.NewServer()

	// Initialize scheduler logic
	srv := scheduler.NewSchedulerServer(logger)

	// Register gRPC service
	pb.RegisterOrchestratorServer(grpcServer, srv)

	// Start job dispatcher in background
	srv.Dispatcher.Run()
	logger.Info("Scheduler listening", zap.String("addr", schedulerHost))

	// Serve
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}

func loadConfig() {
	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		path = "config/dev.yaml"
	}
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
}
