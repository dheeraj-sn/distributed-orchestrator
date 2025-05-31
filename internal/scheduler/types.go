package scheduler

import "time"

type WorkerInfo struct {
	ID       string
	Host     string
	LastSeen time.Time
}

type Job struct {
	ID     string
	Task   string
	Args   []string
	Status string
	Result string
}
