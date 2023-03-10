package databases

import (
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/logs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Connect(l *logs.Log) {
	DB, err = gorm.Open(sqlite.Open("./orquestrator.db"), &gorm.Config{})

	if err != nil {
		l.Fatal("Database open error: " + err.Error())
	}

	err = DB.AutoMigrate(
		&models.Experiment{},
		&models.ExperimentResult{},
		&models.ExperimentResultPerSecondThrouput{},
		&models.Info{},
		&models.ExperimentDeclaration{},
		&models.ExperimentStatus{},
		&models.Worker{},
		&models.WorkerStatus{},
	)

	if err != nil {
		l.Fatal("Database migrate error: " + err.Error())
	}
}
