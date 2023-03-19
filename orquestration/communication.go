package orquestration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

var oLog = logs.Build("orquestrator.log")
var serviceExperiments = experiments.Build(oLog)
var infos = make([]output.Info, 0, 10)
var workers = make([]messages.Worker, 1, 10)
var rexp []output.ExperimentResult
var rexpMutex sync.Mutex
var waitQueueMutex sync.Mutex
var waitQueue []output.ExperimentResult
var expWG sync.WaitGroup

type Orquestrator struct {
	log         *logs.Log
	workers     *seworkers.Workers
	experiments *experiments.Experiments
	client      *nmqtt.Client
	tolerance   int
}

func Build(c *nmqtt.Client, t int) *Orquestrator {
	return &Orquestrator{
		log:         c.Log,
		workers:     seworkers.Build(c.Log),
		experiments: experiments.Build(c.Log),
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
		messageHandlerExperiment(m)
	})

	o.client.Register(*t+"/Experiments/Status", 1, false, func(c mqtt.Client, m mqtt.Message) {
		messageHandlerExperimentStatus(m)
	})
}

func (o Orquestrator) StartExperiment(arg input.Start) ([]output.ExperimentResult, error) {
	expid := time.Now().Unix()

	rexpMutex.Lock()
	rexp = nil
	rexpMutex.Unlock()

	o.experiments.Add(
		models.Experiment{
			Finish: false,
		},
		arg.Description,
		arg.Id...,
	)

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
		return rexp, err
	}

	msg, err := json.Marshal(cmd)

	if err != nil {
		return rexp, err
	}

	if arg.Id[0] == -1 {
		nw := len(workers)
		expWG.Add(nw)
		for i := 0; i < nw; i++ {
			if !workers[i].Status {
				o.log.Register("Worker " + strconv.Itoa(i) + " is off, skipping")
				continue
			}

			workers[i].Historic.Add(expid, cmd, arg.Description.Attempts)

			o.client.Send(workers[i].Id+"/Command", msg)

			// go receiveControl(i, arg.Description.ExecTime*5)

			o.log.Register("Requesting experiment in worker " + strconv.Itoa(i))
		}
	} else {
		argTam := len(arg.Id)
		expWG.Add(argTam)
		for i := 0; i < argTam; i++ {
			if !workers[arg.Id[i]].Status {
				if argTam > 1 {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, skipping")
					continue
				} else {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, aborting experiment")
					break
				}
			}

			workers[i].Historic.Add(expid, cmd, arg.Description.Attempts)

			o.client.Send(workers[arg.Id[i]].Id+"/Command", msg)

			// go receiveControl(arg.Id[i], arg.Description.ExecTime*5)

			o.log.Register("Requesting experiment in worker " + strconv.Itoa(arg.Id[i]))
		}
	}

	expWG.Wait()

	rexpMutex.Lock()
	waitQueueMutex.Lock()
	rexp = append(rexp, waitQueue...)
	rexpMutex.Unlock()

	waitQueue = nil
	waitQueueMutex.Unlock()

	log.Println(o.experiments.List())

	return rexp, nil
}

func (o Orquestrator) CancelExperiment(id int, expid int64) error {
	arg := make(map[string]interface{})
	arg["id"] = expid
	cmd := messages.Command{Name: "cancel", CommandType: "moderation command", Arguments: arg}
	msg, err := json.Marshal(cmd)

	if err != nil {
		return err
	}

	o.client.Send(workers[id].Id+"/Command", msg)

	workers[id].ReceiveConfirmation = true
	exp := workers[id].Historic.Search(expid)
	exp.Finished = true

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
		for i := 0; i < len(workers); i++ {
			if !workers[i].Status {
				o.log.Register("Worker " + strconv.Itoa(i) + " isn't report, skipping")
				continue
			}
			o.client.Register(workers[i].Id+"/Info", 1, true, func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, i)
			})

			o.client.Send(workers[i].Id+"/Command", msg)

			for !workers[i].ReceiveConfirmation {
				if !workers[i].Status {
					o.log.Register("Worker " + strconv.Itoa(i) + " isn't report, skipping")
					break
				}
				time.Sleep(time.Second)
			}
			workers[i].ReceiveConfirmation = false
			o.client.Unregister(workers[i].Id + "/Info")
		}
	} else {
		argTam := len(arg.Id)

		for i := 0; i < argTam; i++ {
			if !workers[arg.Id[i]].Status {
				if argTam > 1 {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, skipping")
					continue
				} else {
					o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, aborting request")
					break
				}
			}

			o.client.Register(workers[arg.Id[i]].Id+"/Info", 1, true, func(c mqtt.Client, m mqtt.Message) {
				messageHandlerInfos(m, arg.Id[i])
			})

			o.client.Send(workers[arg.Id[i]].Id+"/Command", msg)

			for !workers[arg.Id[i]].ReceiveConfirmation {
				if !workers[arg.Id[i]].Status {
					if argTam > 1 {
						o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " isn't report, skipping")
					} else {
						o.log.Register("Worker " + strconv.Itoa(arg.Id[i]) + " is off, aborting request")
					}
					break
				}
				time.Sleep(time.Second)
			}
			workers[arg.Id[i]].ReceiveConfirmation = false
			o.client.Unregister(workers[i].Id + "/Info")
		}
	}

	return infos, nil
}
