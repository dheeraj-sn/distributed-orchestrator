package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"
)

func main() {
	// CLI flags
	mode := flag.String("mode", "submit", "Mode: submit or status")
	task := flag.String("task", "echo", "Task name")
	args := flag.String("args", "hello", "Comma-separated arguments")
	jobID := flag.String("id", "", "Job ID to check status")
	schedulerAddr := flag.String("addr", "localhost:50051", "Scheduler gRPC address")

	flag.Parse()

	// Connect to scheduler
	conn, err := grpc.NewClient(*schedulerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
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
		res, err := client.GetJobStatus(ctx, &pb.JobStatusRequest{
			JobId: *jobID,
		})
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
	var result []string
	result = append(result, splitAndTrim(raw, ",")...)
	return result
}

func splitAndTrim(s, sep string) []string {
	var result []string
	for _, part := range split(s, sep) {
		result = append(result, trim(part))
	}
	return result
}

func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
