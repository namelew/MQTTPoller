package main

import (
	"flag"
	"github.com/namelew/mqtt-bm-latency/internal/worker"
)

func main() {
	var (
		loginTimeout = flag.Int("login_t", 30, "login timeout in seconds")
		broker = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		isUnix = flag.Bool("isunix", true, "define if worker will run a Unix system or not")
		tool = flag.String("tool", "./tools/mqttloader/bin/mqttloader", "benckmark tool for the simulations")
	)
	flag.Parse()

	worker.Init(*broker, *tool, *loginTimeout, *isUnix)
}
