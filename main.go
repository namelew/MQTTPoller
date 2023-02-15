package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/controllers"
	"github.com/namelew/mqtt-bm-latency/orquestration"
)

func main() {
	err := orquestration.Init()

	if err != nil {
		log.Fatal(err.Error())
	}

	api := echo.New()
	api.GET("/orquestrator/worker", controllers.GetWorker)
	api.GET("/orquestrator/info", controllers.GetInfo)
	api.POST("/orquestrator/experiment/start", controllers.StartExperiment)
	api.DELETE("/orquestrator/experiment/cancel/:id/:expid", controllers.CancelExperiment)
	api.Logger.Fatal(api.Start(":8080"))

	err = orquestration.End()

	if err != nil {
		log.Fatal(err.Error())
	}
}
