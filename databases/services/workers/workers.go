package workers

import (
	"github.com/namelew/mqtt-bm-latency/databases"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/logs"
)

type Workers struct {
	log *logs.Log
}

func Build(l *logs.Log) *Workers {
	return &Workers{
		log: l,
	}
}

func (h *Workers) Add(w models.Worker) {
	cerr := make(chan error, 1)

	go func(worker *models.Worker) {
		err := (databases.DB.Create(worker)).Error

		if err != nil {
			cerr <- err
		}

		err = (databases.DB.Create(&models.WorkerStatus{WorkerID: uint64(worker.ID), Online: true, Error: ""})).Error

		cerr <- err
	}(&w)

	err := <-cerr

	if err != nil {
		h.log.Fatal("Unable to add worker in database")
	}
}

func (h *Workers) Remove(id uint) {
	cerr := make(chan error, 1)

	go func() {
		cerr <- (databases.DB.Delete(&models.Worker{}, id)).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Fatal("Unable to remove worker data")
	}
}

func (h *Workers) Update(id uint, new models.Worker) {
	cerr := make(chan error, 1)

	go func() {
		var worker models.Worker

		if err := (databases.DB.Find(&worker, id)).Error; err != nil {
			cerr <- err
		}

		cerr <- (databases.DB.Model(&worker).UpdateColumns(new)).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Fatal("Unable to update worker data")
	}
}

func (h *Workers) ChangeStatus(id uint64, new models.WorkerStatus) {
	cerr := make(chan error, 1)

	go func() {
		cerr <- (databases.DB.Where(&models.WorkerStatus{WorkerID: id}).UpdateColumns(new)).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Fatal("Unable to update worker status")
	}
}

func (h *Workers) List(filter *models.Worker) []models.Worker {
	var workers []models.Worker

	if filter != nil {
		if err := (databases.DB.Where(filter).Find(&workers)).Error; err != nil {
			h.log.Fatal("Unable to find matched workers")
		}
	} else {
		if err := (databases.DB.Find(&workers)).Error; err != nil {
			h.log.Fatal("Unable to get workers registers on database")
		}
	}

	return workers
}

func (h *Workers) Get(id uint) *models.Worker {
	var worker models.Worker

	if err := (databases.DB.Find(&worker, id)); err != nil {
		h.log.Fatal("Unable to find worker")
	}

	return &worker
}
