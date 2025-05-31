package scheduler

import (
	"go.uber.org/zap"
)

type Dispatcher struct {
	JobQueue   chan string
	JobManager *JobManager
	Logger     *zap.Logger
}

func NewDispatcher(jm *JobManager, logger *zap.Logger) *Dispatcher {
	return &Dispatcher{
		JobQueue:   make(chan string, 100),
		JobManager: jm,
		Logger:     logger,
	}
}

func (d *Dispatcher) Run() {
	go func() {
		for jobID := range d.JobQueue {
			job, ok := d.JobManager.Get(jobID)
			if !ok {
				continue
			}

			d.Logger.Info("Dispatching job", zap.String("job_id", jobID), zap.String("task", job.Task))

			// Simulate local execution (replace with remote worker call later)
			d.JobManager.SetStatus(jobID, "completed")

			d.Logger.Info("Job completed", zap.String("job_id", jobID))
		}
	}()
}

func (d *Dispatcher) NextJob() *Job {
	d.JobManager.mu.Lock()
	defer d.JobManager.mu.Unlock()

	for _, job := range d.JobManager.jobs {
		if job.Status == "queued" {
			job.Status = "in_progress"
			return job
		}
	}
	return nil
}
