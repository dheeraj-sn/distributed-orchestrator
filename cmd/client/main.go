package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dheeraj-sn/distributed-orchestrator/internal/config"
	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// CLI flags
	mode := flag.String("mode", "submit", "Mode: submit or status")
	task := flag.String("task", "echo", "Task name")
	args := flag.String("args", "hello", "Comma-separated arguments")
	jobID := flag.String("id", "", "Job ID to check status")

	// Override scheduler address if passed via flag
	addr := flag.String("addr", cfg.Client.SchedulerAddr, "Scheduler gRPC address")

	flag.Parse()

	// Connect to the scheduler
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to scheduler at %s: %v", *addr, err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	switch *mode {
	case "submit":
		req := &pb.JobRequest{
			Task: *task,
			Args: splitArgs(*args),
		}
		res, err := client.SubmitJob(ctx, req)
		if err != nil {
			log.Fatalf("Submit failed: %v", err)
		}
		fmt.Printf("✅ Job submitted. ID: %s\n", res.JobId)

	case "status":
		if *jobID == "" {
			log.Fatal("Please provide a job ID using -id flag")
		}
		res, err := client.GetJobStatus(ctx, &pb.JobStatusRequest{JobId: *jobID})
		if err != nil {
			log.Fatalf("Status check failed: %v", err)
		}
		fmt.Printf("📦 Job Status: %s\n", res.Status)
		if res.Result != "" {
			fmt.Printf("📝 Result: %s\n", res.Result)
		}

	default:
		log.Fatalf("Unknown mode: %s", *mode)
	}
}

func splitArgs(raw string) []string {
	if raw == "" {
		return []string{}
	}
	return splitAndTrim(raw, ",")
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
