package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	// Setup logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// Load config
	viper.SetConfigName("dev")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}

	addr := viper.GetString("scheduler.host")
	logger.Info("Starting scheduler", zap.String("address", addr))

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()

	// TODO: implement and register your scheduler service here
	// pb.RegisterOrchestratorServer(grpcServer, &SchedulerServer{})

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatal("Failed to serve gRPC", zap.Error(err))
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down scheduler")
	grpcServer.GracefulStop()
}
