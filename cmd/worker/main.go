package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "distributed-orchestrator/proto"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	viper.SetConfigName("dev")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}

	addr := viper.GetString("worker.host")
	logger.Info("Starting worker", zap.String("address", addr))

	// TODO: connect to scheduler gRPC endpoint and register worker
	schedulerAddr := viper.GetString("scheduler.host")

	conn, err := grpc.Dial(schedulerAddr, grpc.WithInsecure())
	if err != nil {
		logger.Fatal("Failed to connect to scheduler", zap.Error(err))
	}
	defer conn.Close()

	client := pb.NewOrchestratorClient(conn)

	// TODO: implement worker registration and heartbeat logic here

	// Wait for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down worker")
}
