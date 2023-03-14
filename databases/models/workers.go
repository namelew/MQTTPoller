package models

import "gorm.io/gorm"

type Worker struct {
	gorm.Model
	Token             string
	KeepAliveDeadline uint64
	WorkerID          uint64
	Online            bool
	Error             string
	Experiments       []*Experiment `gorm:"many2many:experiments_workers;"`
}
