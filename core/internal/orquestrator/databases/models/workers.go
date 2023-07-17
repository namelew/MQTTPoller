package models

import (
	"time"
)

type Worker struct {
	ID                uint64 `gorm:"primarykey"`
	Token             string
	KeepAliveDeadline uint64
	Online            bool
	Error             string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Experiments       []*Experiment `gorm:"many2many:experiments_workers;"`
}
