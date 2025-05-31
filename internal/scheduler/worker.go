package scheduler

import (
	"sync"
	"time"
)

type WorkerManager struct {
	mu      sync.RWMutex
	workers map[string]*WorkerInfo
}

func NewWorkerManager() *WorkerManager {
	return &WorkerManager{
		workers: make(map[string]*WorkerInfo),
	}
}

func (wm *WorkerManager) Register(id, host string) {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	wm.workers[id] = &WorkerInfo{
		ID:       id,
		Host:     host,
		LastSeen: time.Now(),
	}
}

func (wm *WorkerManager) Heartbeat(id string) bool {
	wm.mu.Lock()
	defer wm.mu.Unlock()
	if worker, ok := wm.workers[id]; ok {
		worker.LastSeen = time.Now()
		return true
	}
	return false
}
