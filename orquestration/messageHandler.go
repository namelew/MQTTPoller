package orquestration

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
)

func messageHandlerExperimentStatus(msg mqtt.Message, id int) {
	var exps messages.Status
	json.Unmarshal(msg.Payload(), &exps)

	tokens := strings.Split(exps.Type, " ")
	expid, _ := strconv.Atoi(tokens[2])

	exp := workers[id].Historic.Search(int64(expid))

	switch exps.Status {
	case "start":
		return
	case "finish":
		if exp != nil {
			exp.Finished = true
		}
	default:
		if exp != nil {
			exp.Finished = true
			exp.Err = true
			redoExperiment(id, exp)
		}
	}
}

func messageHandlerExperiment(m mqtt.Message, id int) {
	var output output.ExperimentResult

	worker := workers[id].Historic.FindLarger()

	if worker != nil {
		worker.Finished = true
	}

	err := json.Unmarshal(m.Payload(), &output)

	if err != nil {
		log.Fatal(err.Error())
	}

	rexp = append(rexp, output)

	if output.Meta.LogFile.Name != "" {
		ioutil.WriteFile(workers[id].Id+output.Meta.LogFile.Name, output.Meta.LogFile.Data, 0644)
	}

	workers[id].ReceiveConfirmation = true

	log.Printf("Experiment %d in worker %d return\n", output.Meta.ID, id)
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
