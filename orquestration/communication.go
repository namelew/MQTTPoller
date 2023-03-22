package orquestration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/namelew/mqtt-bm-latency/databases"
	"github.com/namelew/mqtt-bm-latency/databases/filters"
	"github.com/namelew/mqtt-bm-latency/databases/models"
	"github.com/namelew/mqtt-bm-latency/databases/services/experiments"
	seworkers "github.com/namelew/mqtt-bm-latency/databases/services/workers"
	"github.com/namelew/mqtt-bm-latency/input"
	"github.com/namelew/mqtt-bm-latency/logs"
	"github.com/namelew/mqtt-bm-latency/messages"
	nmqtt "github.com/namelew/mqtt-bm-latency/network/mqtt"
	"github.com/namelew/mqtt-bm-latency/output"
)

var infos = make([]output.Info, 0, 10)
var ws = make([]messages.Worker, 1, 10)

type queue struct {
	items []output.ExperimentResult
	m *sync.Mutex
}

type Orquestrator struct {
	log         *logs.Log
	workers     *seworkers.Workers
	experiments *experiments.Experiments
	client      *nmqtt.Client
	waitGroup   *sync.WaitGroup
	response 	queue
	repress 	queue
	tolerance   int
}

func Build(c *nmqtt.Client, t int) *Orquestrator {
	return &Orquestrator{
		log:         c.Log,
		workers:     seworkers.Build(c.Log),
		experiments: experiments.Build(c.Log),
		waitGroup: &sync.WaitGroup{},
		response: queue{
			items: []output.ExperimentResult{},
			m: &sync.Mutex{},
		},
		repress: queue{
			items: []output.ExperimentResult{},
			m: &sync.Mutex{},
		},
		client:      c,
		tolerance:   t,
	}
}

func (o Orquestrator) ListWorkers(filter *filters.Worker) []models.Worker {
	return o.workers.List(filter)
}

func (o Orquestrator) GetWorker(id int) *models.Worker {
	return o.workers.Get(id)
}

func (o Orquestrator) timeout(token string, login bool) {
	timer, cancel := context.WithCancel(context.Background())

	o.client.Register(token+"/KeepAlive", 1, true, func(c mqtt.Client, m mqtt.Message) {
		cancel()
		go o.timeout(string(m.Payload()), false)
	})

	tolerance := func(t int) int {
		if login {
			return t * 2
		}
		return t
	}(o.tolerance)

	go func(t context.Context, tk string, tl int) {
		select {
		case <-t.Done():
			return
		case <-time.After(time.Second * time.Duration(tl)):
			o.client.Unregister(tk + "/KeepAlive")
			o.log.Register("lost connection with worker " + tk)
			o.workers.ChangeStatus(&filters.Worker{Token: tk, Online: false})
		}
	}(timer, token, tolerance)
}

func (o Orquestrator) Init() error {
	o.client.Create()

	o.log.Register("Starting database")

	databases.Connect(o.log)

	o.client.Register("Orquestrator/Register", 1, false, func(c mqtt.Client, m mqtt.Message) {
		go func(messagePayload []byte) {
			var clientID string = ""
			worker := string(messagePayload)

			for i := 0; i < 10; i++ {
				seed := rand.NewSource(time.Now().UnixNano())
				random := rand.New(seed)
				clientID += fmt.Sprintf("%d", random.Int()%10)
			}

			o.workers.Add(models.Worker{Token: clientID, KeepAliveDeadline: 1, Online: true, Experiments: nil})

			o.setMessageHandler(&clientID)

			o.log.Register("worker " + worker + " registed as " + clientID)

			o.client.Send("Orquestrator/Register/Log", worker+"-"+clientID)

			go o.timeout(clientID, true)
		}(m.Payload())
	})

	o.client.Register("Orquestrator/Login", 1, false, func(c mqtt.Client, m mqtt.Message) {
		go func(messagePayload []byte) {
			token := string(messagePayload)

			o.workers.ChangeStatus(&filters.Worker{Token: token, Online: true})

			o.log.Register("worker " + token + " loged")

			o.setMessageHandler(&token)

			o.client.Send(token+"/Login/Log", "true")

			go o.timeout(token, true)
		}(m.Payload())
	})

	return nil
}

func (o Orquestrator) End() {
	o.log.Register("disconnect mqtt.client")

	o.client.Disconnect(0)

	o.log.Register("shutdown")
}

func (o Orquestrator) setMessageHandler(t *string) {
	o.client.Register(*t+"/Experiments/Results", 1, false, func(c mqtt.Client, m mqtt.Message) {
		var output output.ExperimentResult

		err := json.Unmarshal(m.Payload(), &output)

		if err != nil {
			log.Fatal(err.Error())
		}

		exp :=  o.experiments.Get(output.Meta.ID)

		if exp.Finish {
			o.repress.m.Lock()
			o.repress.items = append(o.repress.items, output)
			o.repress.m.Unlock()
		} else {
			o.experiments.Update(output.Meta.ID, models.Experiment{Finish: true})
			o.response.m.Lock()
			o.response.items = append(o.response.items, output)
			o.response.m.Unlock()
		}
	})

	o.client.Register(*t+"/Experiments/Status", 1, false, func(c mqtt.Client, m mqtt.Message) {
		var exps messages.Status
		json.Unmarshal(m.Payload(), &exps)

		tokens := strings.Split(exps.Type, " ")
		expid, _ := strconv.Atoi(tokens[2])

		exp := o.experiments.Get(uint64(expid))

		switch exps.Status {
		case "start":
			return
		case "finish":
			if exp.ID != 0 {
				go o.experiments.Update(uint64(expid), models.Experiment{Finish: true})
			}
		default:
			if exp.ID != 0 {
				go o.experiments.Update(uint64(expid), models.Experiment{Finish: true, Error: exps.Status})
				//redoExperiment(id, exp)
			}
		}
	})
}

func (o *Orquestrator) StartExperiment(arg input.Start) ([]output.ExperimentResult, error) {
	expid := time.Now().Unix()

	o.response.m.Lock()
	o.response.items = nil
	o.response.m.Unlock()

	var cmd messages.Command
	var experiment messages.CommandExperiment

	cmd.Name = "start"
	cmd.CommandType = "experiment command"

	experiment.Expid = expid
	experiment.Attempts = arg.Description.Attempts
	experiment.Broker = arg.Description.Broker
	experiment.ExecTime = int(arg.Description.ExecTime)
	experiment.Interval = int(arg.Description.Interval)
	experiment.LogLevel = arg.Description.LogLevel
	experiment.MqttVersion = int(arg.Description.MqttVersion)
	experiment.Ntp = arg.Description.Ntp
	experiment.NumMessages = int(arg.Description.NumMessages)
	experiment.NumPublishers = int(arg.Description.NumPublishers)
	experiment.NumSubscriber = int(arg.Description.NumSubscriber)
	experiment.Output = arg.Description.Output
	experiment.Password = arg.Description.Password
	experiment.Payload = int(arg.Description.Payload)
	experiment.Port = int(arg.Description.Port)
	experiment.QosPublisher = int(arg.Description.QosPublisher)
	experiment.QosSubscriber = int(arg.Description.QosSubscriber)
	experiment.RampDown = arg.Description.RampDown
	experiment.RampUp = arg.Description.RampUp
	experiment.Retain = arg.Description.Retain
	experiment.SharedSubscrition = arg.Description.SharedSubscrition
	experiment.SubscriberTimeout = int(arg.Description.SubscriberTimeout)
	experiment.TlsKeystore = arg.Description.TlsKeystore
	experiment.TlsKeystorePassword = arg.Description.TlsKeystorePassword
	experiment.TlsTrustsore = arg.Description.TlsTrustsore
	experiment.TlsTruststorePassword = arg.Description.TlsTruststorePassword
	experiment.Tool = arg.Description.Tool
	experiment.Topic = arg.Description.Topic
	experiment.User = arg.Description.User

	err := experiment.Attach(&cmd)

	if err != nil {
		return o.response.items, err
	}

	msg, err := json.Marshal(cmd)

	if err != nil {
		return o.response.items, err
	}

	o.experiments.Add(
		models.Experiment{
			ID: uint64(expid),
			Finish: false,
		},
		arg.Description,
		arg.Id...,
	)

	if arg.Id[0] == -1 {
		workers := o.workers.List(nil)

		o.waitGroup.Add(len(workers))

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].Token + " is off, skipping")
				continue
			}

			o.client.Send(workers[i].Token+"/Command", msg)

			// go receiveControl(i, arg.Description.ExecTime*5)

			o.log.Register("Requesting experiment in worker " + workers[i].Token)
		}
	} else {
		workers := make([]*models.Worker, 10)

		for _,i := range arg.Id {
			workers = append(workers, o.workers.Get(i))
		}

		o.waitGroup.Add(len(workers))

		for i := range workers {
			if !workers[i].Online {
				o.log.Register("Worker " + workers[i].Token + " is off, skipping")
			}

			o.client.Send(workers[i].Token+"/Command", msg)

			// go receiveControl(arg.Id[i], arg.Description.ExecTime*5)

			o.log.Register("Requesting experiment in worker " + workers[i].Token)
		}
	}

	o.waitGroup.Wait()

	o.response.m.Lock()
	o.repress.m.Lock()
	o.response.items = append(o.response.items, o.repress.items...)
	o.response.m.Unlock()

	o.repress.items = nil
	o.repress.m.Unlock()

	return o.response.items, nil
}

// vai precisar de um join
func (o Orquestrator) CancelExperiment(id int, expid int64) error {
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: make(map[string]interface{})}
	msg, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	worker := o.workers.Get(id)

	o.client.Send(worker.Token+"/Command", msg)

	o.experiments.Update(uint64(expid), models.Experiment{Finish: true})

	return nil
}

func (o Orquestrator) GetInfo(arg input.Info) ([]output.Info, error) {
	var infoCommand messages.Command
	infos = nil

	infoCommand.Name = "info"
	infoCommand.CommandType = "command moderation"
	infoCommand.Arguments = map[string]interface{}{"cpuDisplay": arg.CpuDisplay, "discDisplay": arg.DiscDisplay, "memoryDisplay": arg.MemoryDisplay}

	msg, err := json.Marshal(&infoCommand)

	if err != nil {
		return infos, err
	}

	if arg.Id[0] == -1 {
		for i := 0; i < len(ws); i++ {
			if !ws[i].Status {
				o.log.Register("Worker " + strconv.Itoa(i) + " isn't report, skipping")
				continue
			}
			o.client.Register(ws[i].Id+"/Info", 1, true, func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, i)
			})

			o.client.Send(ws[i].Id+"/Command", msg)

			for !ws[i].ReceiveConfirmation {
				if !ws[i].Status {
					o.log.Register("Worker " + strconv.Itoa(i) + " isn't report, skipping")
					break
				}
				time.Sleep(time.Second)
			}
			ws[i].ReceiveConfirmation = false
			o.client.Unregister(ws[i].Id + "/Info")
		}
	} else {
		argTam := len(arg.Id)

		for i := 0; i < argTam; i++ {
			if !ws[arg.Id[i]].Status {
				if argTam > 1 {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, skipping")
					continue
				} else {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, aborting request")
					break
				}
			}

			o.client.Register(ws[arg.Id[i]].Id+"/Info", 1, true, func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, arg.Id[i])
			})

			o.client.Send(ws[arg.Id[i]].Id+"/Command", msg)

			for !ws[arg.Id[i]].ReceiveConfirmation {
				if !ws[arg.Id[i]].Status {
					if argTam > 1 {
						o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " isn't report, skipping")
					} else {
						o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, aborting request")
					}
					break
				}
				time.Sleep(time.Second)
			}
			ws[arg.Id[i]].ReceiveConfirmation = false
			o.client.Unregister(ws[i].Id + "/Info")
		}
	}

	return infos, nil
}
