package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dheeraj-sn/distributed-orchestrator/internal/worker"
)

func main() {
	loadConfig()

	logger, _ := zap.NewDevelopment()

	workerID := generateWorkerID()
	logger.Info("Generated random worker ID", zap.String("worker_id", workerID))

	schedulerAddr := viper.GetString("scheduler.host")

	conn, err := grpc.NewClient(schedulerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to create client", zap.Error(err))
	}
	defer conn.Close()

	w := worker.NewWorker(workerID, viper.GetString("worker.host"), conn, logger)
	if err := w.Register(); err != nil {
		logger.Fatal("Registration failed", zap.Error(err))
	}

	w.StartHeartbeat(10)
	w.StartExecutorLoop(3)

	select {}
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

func generateWorkerID() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return fmt.Sprintf("%s-%s", hostname, string(b))
}
