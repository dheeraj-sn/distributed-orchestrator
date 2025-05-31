package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dheeraj-sn/distributed-orchestrator/internal/config"
	"github.com/dheeraj-sn/distributed-orchestrator/internal/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Setup logger
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

	// Generate worker ID if not set
	workerID := cfg.Worker.WorkerID
	if workerID == "" || workerID == "worker-default" {
		workerID = generateWorkerID()
		logger.Info("Generated random worker ID", zap.String("worker_id", workerID))
	}

	// Connect to scheduler
	schedulerAddr := cfg.Client.SchedulerAddr
	conn, err := grpc.NewClient(schedulerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("Failed to connect to scheduler", zap.Error(err))
	}
	defer conn.Close()

	// Create and start worker
	w := worker.NewWorker(
		workerID,
		cfg.Worker.Host,
		conn,
		logger,
		cfg.Worker.Concurrency,
	)

	if err := w.Register(); err != nil {
		logger.Fatal("Worker registration failed", zap.Error(err))
	}

	w.StartHeartbeat(10)
	w.StartExecutorLoop(3)

	select {}
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
