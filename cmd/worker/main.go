package main

import (
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"distributed-orchestrator/internal/worker"
)

func main() {
	logger, _ := zap.NewDevelopment()

	workerID := os.Getenv("WORKER_ID")
	if workerID == "" {
		workerID = "worker-1"
	}

	schedulerAddr := os.Getenv("SCHEDULER_ADDR")
	if schedulerAddr == "" {
		schedulerAddr = "localhost:50051"
	}

	conn, err := grpc.Dial(schedulerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to scheduler: %v", err)
	}
	defer conn.Close()

	w := worker.NewWorker(workerID, "localhost", conn, logger)

	if err := w.Register(); err != nil {
		log.Fatalf("Registration failed: %v", err)
	}

	w.StartHeartbeat(10 * time.Second)
	w.StartMockExecutor()

	select {} // block forever
}
