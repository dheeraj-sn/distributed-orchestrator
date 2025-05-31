package storage

import (
	"time"

	"gorm.io/gorm"
)

type Worker struct {
	ID            string    `gorm:"primaryKey"`
	Host          string    `gorm:"not null"`
	LastHeartbeat time.Time `gorm:"not null"`
}

func RegisterWorker(worker *Worker) error {
	return DB.Save(worker).Error
}

func UpdateHeartbeat(workerID string) error {
	return DB.Model(&Worker{}).Where("id = ?", workerID).Update("last_heartbeat", gorm.Expr("NOW()")).Error
}

func GetWorker(id string) (*Worker, error) {
	var w Worker
	if err := DB.First(&w, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &w, nil
}
