package storage

import "github.com/dheeraj-sn/distributed-orchestrator/internal/models"

// Store defines the interface for task storage
type Store interface {
	SaveTask(task *models.Task) error
	GetTask(id string) (*models.Task, error)
	UpdateTask(task *models.Task) error
	ListTasks() ([]*models.Task, error)
}

// MemoryStore implements Store interface with in-memory storage
type MemoryStore struct {
	tasks map[string]*models.Task
}

// NewMemoryStore creates a new in-memory store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[string]*models.Task),
	}
} 