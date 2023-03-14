package experiments

import (
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/logs"
)

type Experiments struct {
	log *logs.Log
}

type Register interface {
	models.Experiment | models.ExperimentResult | models.ExperimentDeclaration
}

func Build(l *logs.Log) *Experiments {
	return &Experiments{
		log: l,
	}
}

func (h *Experiments) Add() {
}

func (h *Experiments) Remove() {

}

func (h *Experiments) Update(key filters.Experiment, new models.Experiment) {

}

func (h *Experiments) List() {

}

func (h *Experiments) Get(filter filters.Experiment) models.Experiment {
	return models.Experiment{}
}
