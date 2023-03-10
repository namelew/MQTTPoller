package models

import "gorm.io/gorm"

type Worker struct {
	gorm.Model
	Token             string
	KeepAliveDeadline uint64
	Experiments       []*Experiment `gorm:"many2many:experiments_workers;"`
}

type WorkerStatus struct {
	WorkerID uint64
	Online   bool
	Error    string
}
