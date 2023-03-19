package experiments

import (
	"github.com/namelew/mqtt-bm-latency/databases"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/logs"
)

type Experiments struct {
	log *logs.Log
}

func Build(l *logs.Log) *Experiments {
	return &Experiments{
		log: l,
	}
}

func (h *Experiments) Add(e models.Experiment, d models.ExperimentDeclaration, wid ...int) {
	var wkrs []models.Worker
	var declaration models.ExperimentDeclaration

	if (databases.DB.Where("id IN ?", wid).Find(&wkrs)).Error != nil {
		h.log.Fatal("Unable to find workers")
	}

	if (databases.DB.Model(d).Find(&declaration)).Error != nil {
		h.log.Fatal("Unable to query experiment declaration")
	}

	if declaration != d {
		if (databases.DB.Create(d)).Error != nil {
			h.log.Fatal("Unable to create experiment declaration")
		}
	}

	for i := range wkrs {
		e.Workers = append(e.Workers, &wkrs[i])
	}

	e.ExperimentDeclaration = declaration

	if (databases.DB.Create(e)).Error != nil {
		h.log.Fatal("Unable to register experiment")
	}
}

func (h *Experiments) Remove(id uint64) {
	if (databases.DB.Model(&models.Experiment{}).Where("id = ?", id).Delete(&models.Experiment{})).Error != nil {
		h.log.Fatal("Unable to remove experiment")
	}
}

func (h *Experiments) Update(key filters.Experiment, new models.Experiment) {

}

func (h *Experiments) List() []models.Experiment {
	var experiments []models.Experiment

	if err := (databases.DB.Find(&experiments)).Error; err != nil {
		h.log.Fatal("Unable to get experiments registers on database")
	}

	return experiments
}

func (h *Experiments) Get(id uint64) models.Experiment {
	var experiment models.Experiment

	if err := (databases.DB.Where("id = ?", id).Find(&experiment)).Error; err != nil {
		h.log.Fatal("Unable to get experiment register on database")
	}

	return experiment
}
