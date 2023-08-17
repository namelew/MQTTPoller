package workers

import (
	"time"

	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/filters"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/databases/models"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
)

type Workers struct {
	log *logs.Log
}

func Build(l *logs.Log) *Workers {
	return &Workers{
		log: l,
	}
}

func (h *Workers) Add(w models.Worker) error{
	cerr := make(chan error, 1)

	go func(worker *models.Worker) {
		cerr <- (databases.DB.Create(worker)).Error
	}(&w)

	err := <-cerr

	if err != nil {
		h.log.Register("Unable to add worker in database")
	}

	return err
}

func (h *Workers) Remove(id uint) error{
	cerr := make(chan error, 1)

	go func() {
		cerr <- (databases.DB.Model(&models.Worker{}).Where("id = ?", id).Delete(&models.Worker{})).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Register("Unable to remove worker data")
	}

	return err
}

func (h *Workers) Update(id uint, new models.Worker) error{
	cerr := make(chan error, 1)

	go func() {
		var worker models.Worker

		if err := (databases.DB.Model(&models.Worker{}).Where("id = ?", id).Find(&worker)).Error; err != nil {
			cerr <- err
		}

		cerr <- (databases.DB.Model(&worker).UpdateColumns(new)).Error
	}()

	err := <-cerr

	if err != nil {
		h.log.Register("Unable to update worker data")
	}

	return err
}

func (h *Workers) ChangeStatus(new *filters.Worker) {
	var worker models.Worker

	if (databases.DB.Model(&models.Worker{}).Where("token = ?", new.Token).Find(&worker)).Error != nil || worker.ID == 0 {
		h.log.Register("Update error, unable to find worker")
	}

	if worker.Online != new.Online {
		if (databases.DB.Model(&models.Worker{}).Where("id = ?", worker.ID).Update("online", new.Online)).Error != nil {
			h.log.Register("Unable to update worker status")
		}
	}
}

func (h *Workers) List(filter *filters.Worker) ([]models.Worker, error) {
	var workers []models.Worker

	if filter != nil {
		err := (databases.DB.Where(&models.Worker{
			ID:     filter.WorkerID,
			Token:  filter.Token,
			Online: filter.Online,
			Error:  filter.Error,
		}).Find(&workers)).Error;

		if err != nil {
			h.log.Register("Unable to find matched workers")
			return workers, err
		}
	} else {
		if err := (databases.DB.Find(&workers)).Error; err != nil {
			h.log.Register("Unable to get workers registers on database")
			return workers, err
		}
	}

	return workers, nil
}

func (h *Workers) Get(id int) (*models.Worker, error) {
	var worker models.Worker

	err := (databases.DB.Model(&models.Worker{}).Where("id = ?", id).Find(&worker)).Error
	
	if err != nil {
		h.log.Register("Unable to find worker")
		return nil, err
	}

	return &worker, nil
}

func (h *Workers) TrashOut(i time.Time) {
	if (databases.DB.Model(&models.Worker{}).Where("online = ? AND (created_at < ? OR updated_at < ?)", false, i, i).Delete(&models.Worker{})).Error != nil {
		h.log.Fatal("Unable to remove retired worker")
	}
}

func (h *Workers) ID() int {
	return 1
}
