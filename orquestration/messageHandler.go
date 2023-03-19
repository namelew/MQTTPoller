package orquestration

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
)

func messageHandlerExperimentStatus(msg mqtt.Message) {
	var exps messages.Status
	json.Unmarshal(msg.Payload(), &exps)

	tokens := strings.Split(exps.Type, " ")
	expid, _ := strconv.Atoi(tokens[2])

	exp := serviceExperiments.Get(uint64(expid))

	switch exps.Status {
	case "start":
		return
	case "finish":
		if exp.ID != 0 {
			exp.Finish = true
			go serviceExperiments.Update(uint64(expid), exp)
		}
	default:
		if exp.ID != 0 {
			exp.Finish = true
			exp.Error = exps.Status
			go serviceExperiments.Update(uint64(expid), exp)
			//redoExperiment(id, exp)
		}
	}
}

func messageHandlerExperiment(m mqtt.Message) {
	var output output.ExperimentResult

	err := json.Unmarshal(m.Payload(), &output)

	if err != nil {
		log.Fatal(err.Error())
	}

	// worker :=  workers[id].Historic.FindLarger()

	// if worker == nil {
	// 	waitQueueMutex.Lock()
	// 	waitQueue = append(waitQueue, output)
	// 	waitQueueMutex.Unlock()
	// } else {
	// 	worker.Finished = true

	// 	rexpMutex.Lock()
	// 	rexp = append(rexp, output)
	// 	rexpMutex.Unlock()
	// 	workers[id].ReceiveConfirmation = true
	// }

	// log.Printf("Experiment %d in worker %d return\n", output.Meta.ID, id)
}

func messageHandlerInfos(m mqtt.Message, id int) {
	var response messages.Info
	var output output.Info

	json.Unmarshal(m.Payload(), &response)
	workers[id].ReceiveConfirmation = true

	output.Id = id
	output.NetId = workers[id].Id
	output.Infos = response

	infos = append(infos, output)
}
