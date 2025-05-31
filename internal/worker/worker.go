package worker

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Worker struct {
	ID          string
	Host        string
	Client      pb.OrchestratorClient
	Logger      *zap.Logger
	stopChan    chan struct{}
	concurrency int
	sem         chan struct{}
	wg          sync.WaitGroup
}

func NewWorker(id, host string, conn *grpc.ClientConn, logger *zap.Logger, concurrency int) *Worker {
	client := pb.NewOrchestratorClient(conn)
	return &Worker{
		ID:          id,
		Host:        host,
		Client:      client,
		Logger:      logger,
		stopChan:    make(chan struct{}),
		concurrency: concurrency,
		sem:         make(chan struct{}, concurrency),
	}
}

func (w *Worker) Register() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := w.Client.RegisterWorker(ctx, &pb.RegisterWorkerRequest{
		WorkerId: w.ID,
		Host:     w.Host,
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
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				w.Client.SendHeartbeat(ctx, &pb.HeartbeatRequest{
					WorkerId: w.ID,
				})
				cancel()
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	close(w.stopChan)
	w.Logger.Info("Worker shutting down", zap.String("id", w.ID))
	w.wg.Wait()
}

func (w *Worker) StartExecutorLoop(pollInterval time.Duration) {
	go func() {
		for {
			select {
			case <-w.stopChan:
				return
			default:
				time.Sleep(pollInterval)

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				resp, err := w.Client.PullJob(ctx, &pb.PullJobRequest{
					WorkerId: w.ID,
				})
				cancel()
				if err != nil || !resp.Found {
					continue
				}

				w.sem <- struct{}{}
				w.wg.Add(1)

				go func(jobID, task string, args []string) {
					defer func() {
						<-w.sem
						w.wg.Done()
					}()

					w.Logger.Info("Pulled job", zap.String("id", jobID), zap.String("task", task))
					result := fmt.Sprintf("Executed task %s with args %v", task, args)
					time.Sleep(2 * time.Second)

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					_, err := w.Client.CompleteJob(ctx, &pb.CompleteJobRequest{
						JobId:  jobID,
						Result: result,
					})
					cancel()

					if err != nil {
						w.Logger.Error("Failed to complete job", zap.Error(err))
					} else {
						w.Logger.Info("Reported job completion", zap.String("job_id", jobID))
					}
				}(resp.JobId, resp.Task, resp.Args)
			}
		}
	}()
}

func WaitForShutdown(w *Worker) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	w.Stop()
	time.Sleep(1 * time.Second)
}
