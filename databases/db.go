package databases

import (
	"github.com/namelew/mqtt-bm-latency/logs"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func Connect(l *logs.Log) {
	DB, err := gorm.Open(sqlite.Open("./orquestrator.db"), &gorm.Config{})

	if err != nil {
		l.Fatal("Database open error: " + err.Error())
	}

	err = DB.AutoMigrate()

	if err != nil {
		l.Fatal("Database migrate error: " + err.Error())
	}
}
