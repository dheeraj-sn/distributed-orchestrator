package scheduler

import (
	"sync"

	"github.com/google/uuid"
)

type JobManager struct {
	mu   sync.RWMutex
	jobs map[string]*Job
}

func NewJobManager() *JobManager {
	return &JobManager{
		jobs: make(map[string]*Job),
	}
}

func (jm *JobManager) Submit(task string, args []string) string {
	jm.mu.Lock()
	defer jm.mu.Unlock()

	id := uuid.New().String()
	jm.jobs[id] = &Job{
		ID:     id,
		Task:   task,
		Args:   args,
		Status: "queued",
	}
	return id
}

func (jm *JobManager) Get(id string) (*Job, bool) {
	jm.mu.RLock()
	defer jm.mu.RUnlock()
	job, ok := jm.jobs[id]
	return job, ok
}

func (jm *JobManager) SetStatus(id, status string) {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	if job, ok := jm.jobs[id]; ok {
		job.Status = status
	}
}

func (jm *JobManager) Complete(id string, result string) bool {
	jm.mu.Lock()
	defer jm.mu.Unlock()
	if job, ok := jm.jobs[id]; ok {
		job.Status = "completed"
		job.Result = result
		return true
	}
	return false
}
