package orquestration

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/messages"
	"github.com/namelew/mqtt-bm-latency/output"
)

func Register(c mqtt.Client, m mqtt.Message) {
	var clientID string = ""
	var seed rand.Source
	var random *rand.Rand

	for i := 0; i < 10; i++ {
		seed = rand.NewSource(time.Now().UnixNano())
		random = rand.New(seed)
		clientID += fmt.Sprintf("%d", random.Int()%10)
	}

	go serviceWorkers.Add(models.Worker{Token: clientID, KeepAliveDeadline: 1, Online: true, Experiments: nil})

	setMessageHandler(&clientID)

	oLog.Register("worker " + clientID + " registed")

	t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
	t.Wait()
}

func Login(c mqtt.Client, m mqtt.Message) {
	token := string(m.Payload())
	go serviceWorkers.ChangeStatus(&filters.Worker{Token: token, Online: true})

	oLog.Register("worker " + token + " loged")

	setMessageHandler(&token)

	t := client.Publish(token+"/Login/Log", byte(1), true, "true")
	t.Wait()
}

func Ping(c mqtt.Client, m mqtt.Message, t int) {
	id := 0
	for i := 0; i < len(workers); i++ {
		if workers[i].Id == string(m.Payload()) {
			workers[i].Status = true
			id = i
			break
		}
	}
	if workers[id].TestPing {
		workers[id].TestPing = false
		go watcher(id, t)
	} else {
		workers[id].TestPing = true
		go watcher(id, t)
		workers[id].TestPing = false
	}
}

func messageHandlerExperimentStatus(msg mqtt.Message) {
	var exps messages.Status
	json.Unmarshal(msg.Payload(), &exps)

	tokens := strings.Split(exps.Type, " ")
	expid, _ := strconv.Atoi(tokens[2])

	exp := serviceExperiments.Get(filters.Experiment{ExperimentID: uint64(expid)})

	switch exps.Status {
	case "start":
		return
	case "finish":
		if exp.ID != 0 {
			exp.Finish = true
			go serviceExperiments.Update(filters.Experiment{ExperimentID: uint64(expid)}, exp)
		}
	default:
		if exp.ID != 0 {
			exp.Finish = true
			exp.Error = exps.Status
			go serviceExperiments.Update(filters.Experiment{ExperimentID: uint64(expid)}, exp)
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
