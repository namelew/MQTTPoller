package main

import (
	"flag"

	"github.com/namelew/mqtt-bm-latency/internal/worker"
)

func main() {
	var (
		loginTimeout = flag.Int("login_t", 30, "login timeout in seconds")
		loginAttempts = flag.Int("login_th", -1, "login attempts threshold")
		broker       = flag.String("broker", "tcp://localhost:1883", "broker url to worker/orquestrator communication")
		tool         = flag.String("tool", "./tools/mqttloader/bin/mqttloader", "benckmark tool for the simulations")
	)
	flag.Parse()

	worker.Init(*broker, *tool, *loginTimeout, *loginAttempts)
}
