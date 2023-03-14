package experiments

import (
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/logs"
)

type History struct {
	log *logs.Log
}

type Register interface {
	models.Experiment | models.ExperimentResult | models.ExperimentDeclaration
}

func Build(l *logs.Log) *History {
	return &History{
		log: l,
	}
}

func (h *History) Add() {
}

func (h *History) Remove() {

}

func (h *History) Update() {

}

func (h *History) List() {

}

func (h *History) Get() {

}
