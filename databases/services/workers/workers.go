package workers

import (
	"github.com/namelew/mqtt-bm-latency/databases"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
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
		cerr <- (databases.DB.Create(worker)).Error
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

func (h *Workers) ChangeStatus(new *filters.Worker) {
	cerr := make(chan error, 1)

	go func() {
		worker := models.Worker{
			Online: new.Online,
			Error:  new.Error,
		}
		cerr <- (databases.DB.Where(&models.Worker{WorkerID: new.WorkerID, Token: new.Token}).UpdateColumns(worker)).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Fatal("Unable to update worker status")
	}
}

func (h *Workers) List(filter *filters.Worker) []models.Worker {
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
