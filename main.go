package main

import (
	"log"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/controllers"
	"github.com/namelew/mqtt-bm-latency/orquestration"
)

func main() {
	var (
		port     = flag.String("port", "8000", "api default port")
		broker     = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		t_interval = flag.Int("tl", 5, "orquestrator tolerance interval")
	)
	flag.Parse()

	err := orquestration.Init(*broker, *t_interval)

	if err != nil {
		log.Fatal(err.Error())
	}

	c := make(chan os.Signal)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		err = orquestration.End()

		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(1)
	}()

	api := echo.New()
	
	api.GET("/orquestrator/worker", controllers.GetWorker)
	api.GET("/orquestrator/info", controllers.GetInfo)
	api.POST("/orquestrator/experiment/start", controllers.StartExperiment)
	api.DELETE("/orquestrator/experiment/cancel/:id/:expid", controllers.CancelExperiment)

	api.Logger.Fatal(api.Start(":"+*port))

	err = orquestration.End()

	if err != nil {
		log.Fatal(err.Error())
	}
}
