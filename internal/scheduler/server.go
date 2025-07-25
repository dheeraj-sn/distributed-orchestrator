package scheduler

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	pb "github.com/dheeraj-sn/distributed-orchestrator/proto"
)

type SchedulerServer struct {
	pb.UnimplementedOrchestratorServer

	Jobs       *JobManager
	Workers    *WorkerManager
	Dispatcher *Dispatcher
	Logger     *zap.Logger
}

func NewSchedulerServer(logger *zap.Logger) *SchedulerServer {
	jobs := NewJobManager()
	dispatcher := NewDispatcher(jobs, logger)
	workers := NewWorkerManager()

	return &SchedulerServer{
		Jobs:       jobs,
		Workers:    workers,
		Dispatcher: dispatcher,
		Logger:     logger,
	}
}

func (s *SchedulerServer) SubmitJob(ctx context.Context, req *pb.JobRequest) (*pb.JobResponse, error) {
	jobID := s.Jobs.Submit(req.Task, req.Args)
	s.Dispatcher.JobQueue <- jobID

	s.Logger.Info("Job submitted", zap.String("job_id", jobID), zap.String("task", req.Task))
	return &pb.JobResponse{JobId: jobID}, nil
}

func (s *SchedulerServer) GetJobStatus(ctx context.Context, req *pb.JobStatusRequest) (*pb.JobStatusResponse, error) {
	job, ok := s.Jobs.Get(req.JobId)
	if !ok {
		return &pb.JobStatusResponse{Status: "not_found"}, nil
	}
	return &pb.JobStatusResponse{Status: job.Status, Result: job.Result}, nil
}

func (s *SchedulerServer) RegisterWorker(ctx context.Context, req *pb.RegisterWorkerRequest) (*pb.RegisterWorkerResponse, error) {
	s.Workers.Register(req.WorkerId, req.Host)
	s.Logger.Info("Worker registered", zap.String("worker_id", req.WorkerId), zap.String("host", req.Host))
	return &pb.RegisterWorkerResponse{Success: true}, nil
}

func (s *SchedulerServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	alive := s.Workers.Heartbeat(req.WorkerId)
	return &pb.HeartbeatResponse{Alive: alive}, nil
}

func (s *SchedulerServer) PullJob(ctx context.Context, req *pb.PullJobRequest) (*pb.PullJobResponse, error) {
	job := s.Dispatcher.NextJob()
	if job == nil {
		return &pb.PullJobResponse{Found: false}, nil
	}
	s.Logger.Info("Job pulled", zap.String("job_id", job.ID), zap.String("worker", req.WorkerId))
	return &pb.PullJobResponse{
		Found: true,
		JobId: job.ID,
		Task:  job.Task,
		Args:  job.Args,
	}, nil
}

func (s *SchedulerServer) CompleteJob(ctx context.Context, req *pb.CompleteJobRequest) (*pb.CompleteJobResponse, error) {
	ok := s.Jobs.Complete(req.JobId, req.Result)
	if ok {
		s.Logger.Info("Job completed", zap.String("job_id", req.JobId))
	}
	return &pb.CompleteJobResponse{Success: ok}, nil
}

func (s *SchedulerServer) ListJobs(ctx context.Context, req *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	s.Jobs.mu.RLock()
	defer s.Jobs.mu.RUnlock()

	var jobStatuses []*pb.JobStatus
	for _, j := range s.Jobs.jobs {
		jobStatuses = append(jobStatuses, &pb.JobStatus{
			JobId:  j.ID,
			Status: j.Status,
			Result: j.Result,
		})
	}

	return &pb.ListJobsResponse{Jobs: jobStatuses}, nil
}

func (s *SchedulerServer) StreamLogs(stream pb.Orchestrator_StreamLogsServer) error {
	for {
		entry, err := stream.Recv()
		if err != nil {
			return err
		}
		fmt.Printf("[Job %s | Worker %s | %s] %s\n", entry.JobId, entry.WorkerId, entry.Timestamp, entry.Message)
		if err := stream.Send(&pb.LogAck{Received: true}); err != nil {
			return err
		}
	}
}
