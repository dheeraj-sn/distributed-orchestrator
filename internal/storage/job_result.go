package storage

import (
	"time"
)

type JobResult struct {
	JobID     string `gorm:"primaryKey;type:uuid"`
	Output    string `gorm:"type:text"`
	Logs      string `gorm:"type:text"`
	CreatedAt time.Time
}

func SaveJobResult(result *JobResult) error {
	return DB.Save(result).Error
}

func GetJobResult(jobID string) (*JobResult, error) {
	var res JobResult
	if err := DB.First(&res, "job_id = ?", jobID).Error; err != nil {
		return nil, err
	}
	return &res, nil
}
