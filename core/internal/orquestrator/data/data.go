package data

import (
	"github.com/namelew/mqtt-poller/core/internal/orquestrator/data/models"
	"github.com/namelew/mqtt-poller/core/internal/orquestrator/data/services"
	"github.com/namelew/mqtt-poller/core/packages/logs"
)

const EXPERIMENTS_DATA string = "experiments.data"
const WORKERS_DATA string = "workers.data"

var (
	ExperimentTable *services.Table[uint64, models.Experiment]
	WorkersTable    *services.Table[string, models.Worker]
)

func Init(l *logs.Log) {
	ExperimentTable = services.Build[uint64, models.Experiment](EXPERIMENTS_DATA, l)
	WorkersTable = services.Build[string, models.Worker](WORKERS_DATA, l)
}
