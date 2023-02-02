package utils

import "github.com/namelew/mqtt-bm-latency/messages"

func IsIn(workers []messages.Worker, clientID string) bool {
	for _, w := range workers {
		if clientID == w.Id {
			return true
		}
	}
	return false
}
