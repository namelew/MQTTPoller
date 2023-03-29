package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/namelew/mqtt-bm-latency/internal/controllers"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/network"
	"github.com/namelew/mqtt-bm-latency/internal/orquestration"
)

func main() {
	var (
		port       = flag.String("port", "8000", "api default port")
		broker     = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		t_interval = flag.Int("tl", 5, "orquestrator tolerance interval")
	)
	flag.Parse()

	var oLog = logs.Build("orquestrator.log")
	oLog.Create()

	orquestrator := orquestration.Build(&network.Client{
		Broker: *broker,
		ID:     "Orquestrator",
		KA:     time.Second * 1000,
		Log:    oLog,
	}, *t_interval)

	err := orquestrator.Init()

	if err != nil {
		oLog.Fatal(err.Error())
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		orquestrator.End()
		os.Exit(1)
	}()

	api := echo.New()

	controller := controllers.Build(orquestrator)

	api.GET("/orquestrator/worker", controller.List)
	api.GET("/orquestrator/worker/:id", controller.Get)
	api.GET("/orquestrator/info", controller.List)
	api.POST("/orquestrator/experiment/start", controller.Procedure)
	api.DELETE("/orquestrator/experiment/cancel/:id/:expid", controller.Procedure)

	api.Logger.Fatal(api.Start(":" + *port))

	orquestrator.End()
}
