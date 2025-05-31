package scheduler

import (
	"context"

	"go.uber.org/zap"

	pb "distributed-orchestrator/proto"
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
	return &pb.JobResponse{Job_id: jobID}, nil
}

func (s *SchedulerServer) GetJobStatus(ctx context.Context, req *pb.JobStatusRequest) (*pb.JobStatusResponse, error) {
	job, ok := s.Jobs.Get(req.Job_id)
	if !ok {
		return &pb.JobStatusResponse{Status: "not_found"}, nil
	}
	return &pb.JobStatusResponse{Status: job.Status}, nil
}

func (s *SchedulerServer) RegisterWorker(ctx context.Context, req *pb.RegisterWorkerRequest) (*pb.RegisterWorkerResponse, error) {
	s.Workers.Register(req.Worker_id, req.Host)
	s.Logger.Info("Worker registered", zap.String("worker_id", req.Worker_id), zap.String("host", req.Host))
	return &pb.RegisterWorkerResponse{Success: true}, nil
}

func (s *SchedulerServer) SendHeartbeat(ctx context.Context, req *pb.HeartbeatRequest) (*pb.HeartbeatResponse, error) {
	alive := s.Workers.Heartbeat(req.Worker_id)
	return &pb.HeartbeatResponse{Alive: alive}, nil
}
