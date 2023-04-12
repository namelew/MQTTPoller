package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/namelew/mqtt-bm-latency/internal/orquestrator"
	"github.com/namelew/mqtt-bm-latency/internal/orquestrator/router"
	"github.com/namelew/mqtt-bm-latency/packages/logs"
	"github.com/namelew/mqtt-bm-latency/packages/network"
)

func main() {
	var (
		port        = flag.String("port", "8000", "api default port")
		broker      = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		hk_interval = flag.Int("hk-interval", 1, "housekeeper executions interval in hours")
		t_interval  = flag.Int("tl", 5, "orquestrator tolerance interval in seconds")
	)
	flag.Parse()

	var oLog = logs.Build("orquestrator.log")
	oLog.Create()

	o := orquestrator.Build(&network.Client{
		Broker: *broker,
		ID:     "Orquestrator",
		KA:     time.Second * 1000,
		Log:    oLog,
	}, *t_interval, *hk_interval)

	err := o.Init()

	if err != nil {
		oLog.Fatal(err.Error())
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		o.End()
		os.Exit(1)
	}()

	router.Route(o, *port)

	o.End()
}
