package main

import (
	"net"
	"os"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"distributed-orchestrator/internal/scheduler"
	pb "distributed-orchestrator/proto"
)

func main() {
	// Logger setup
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Read port from env or use default
	port := os.Getenv("SCHEDULER_PORT")
	if port == "" {
		port = "50051"
	}
	addr := ":" + port

	// Start TCP listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}
	logger.Info("Scheduler listening", zap.String("addr", addr))

	// gRPC server
	grpcServer := grpc.NewServer()

	// Initialize scheduler logic
	srv := scheduler.NewSchedulerServer(logger)

	// Register gRPC service
	pb.RegisterOrchestratorServer(grpcServer, srv)

	// Start job dispatcher in background
	srv.Dispatcher.Run()

	// Serve
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("Failed to serve", zap.Error(err))
	}
}
