package storage

import (
	"time"

	"github.com/lib/pq"
)

type Job struct {
	ID        string         `gorm:"type:uuid;primaryKey"`
	Task      string         `gorm:"not null"`
	Args      pq.StringArray `gorm:"type:text[]"`
	Status    string         `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func SaveJob(job *Job) error {
	return DB.Save(job).Error
}

func GetJob(id string) (*Job, error) {
	var job Job
	if err := DB.First(&job, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

func UpdateJobStatus(id string, status string) error {
	return DB.Model(&Job{}).Where("id = ?", id).Update("status", status).Error
}
