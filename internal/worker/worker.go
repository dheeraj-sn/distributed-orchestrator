package worker

import (
	"context"
	"fmt"
	"time"

	pb "distributed-orchestrator/proto"

	"go.uber.org/zap"

	"google.golang.org/grpc"
)

type Worker struct {
	ID       string
	Host     string
	Client   pb.OrchestratorClient
	Logger   *zap.Logger
	stopChan chan struct{}
}

func NewWorker(id, host string, conn *grpc.ClientConn, logger *zap.Logger) *Worker {
	client := pb.NewOrchestratorClient(conn)
	return &Worker{
		ID:       id,
		Host:     host,
		Client:   client,
		Logger:   logger,
		stopChan: make(chan struct{}),
	}
}

func (w *Worker) Register() error {
	_, err := w.Client.RegisterWorker(context.Background(), &pb.RegisterWorkerRequest{
		Worker_id: w.ID,
		Host:      w.Host,
	})
	if err == nil {
		w.Logger.Info("Worker registered", zap.String("id", w.ID))
	}
	return err
}

func (w *Worker) StartHeartbeat(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				w.Client.SendHeartbeat(context.Background(), &pb.HeartbeatRequest{
					Worker_id: w.ID,
				})
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) StartMockExecutor() {
	go func() {
		for {
			// Simulate pulling a job (to be replaced with real stream/pull API)
			time.Sleep(5 * time.Second)
			task := "echo"
			args := []string{"hello from worker " + w.ID}

			w.Logger.Info("Executing job", zap.String("task", task))
			output := fmt.Sprintf("Executed task: %s with args: %v", task, args)
			w.Logger.Info("Job completed", zap.String("output", output))
		}
	}()
}

func (w *Worker) Stop() {
	close(w.stopChan)
}
