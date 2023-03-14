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
	"github.com/namelew/mqtt-bm-latency/utils"
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

	oLog.Register("worker " + clientID + " registed")

	serviceWorkers.Add(models.Worker{Token: clientID, KeepAliveDeadline: 1, Experiments: nil})

	if workers[0].Id == "" {
		workers[0] = messages.Worker{Id: clientID, Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}}

		t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
		t.Wait()

		setMessageHandler(0)
		return
	}
	workers = append(workers, messages.Worker{Id: clientID, Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}})

	t := client.Publish("Orquestrator/Register/Log", byte(1), false, string(m.Payload())+"-"+clientID)
	t.Wait()

	setMessageHandler(len(workers) - 1)
}

func Login(c mqtt.Client, m mqtt.Message) {
	serviceWorkers.ChangeStatus(&filters.Worker{Token: string(m.Payload()), Online: true})

	oLog.Register("worker " + string(m.Payload()) + " loged")
	if !utils.IsIn(workers, string(m.Payload())) {
		if workers[0].Id == "" {
			workers[0] = messages.Worker{Id: string(m.Payload()), Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}}

			t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
			t.Wait()

			setMessageHandler(0)
			return
		}
		workers = append(workers, messages.Worker{Id: string(m.Payload()), Status: true, ReceiveConfirmation: false, TestPing: true, Historic: messages.ExperimentHistory{}})

		t := client.Publish(string(m.Payload())+"/Login/Log", byte(1), true, "true")
		t.Wait()

		setMessageHandler(len(workers) - 1)
	} else {
		id := 0
		for i := 0; i < len(workers); i++ {
			if workers[i].Id == string(m.Payload()) {
				id = i
				break
			}
		}
		workers[id].Status = true
		setMessageHandler(id)
	}
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

	err := json.Unmarshal(m.Payload(), &output)

	if err != nil {
		log.Fatal(err.Error())
	}

	worker := workers[id].Historic.FindLarger()

	if worker == nil {
		waitQueueMutex.Lock()
		waitQueue = append(waitQueue, output)
		waitQueueMutex.Unlock()
	} else {
		worker.Finished = true

		rexpMutex.Lock()
		rexp = append(rexp, output)
		rexpMutex.Unlock()
		workers[id].ReceiveConfirmation = true
	}

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
